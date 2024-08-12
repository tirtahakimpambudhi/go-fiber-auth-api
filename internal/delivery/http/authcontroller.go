// Package http provides HTTP handlers for authentication operations
package http

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/phuslu/log"
	"github.com/tirtahakimpambudhi/restful_api/internal/configs/token"
	errorshandler "github.com/tirtahakimpambudhi/restful_api/internal/errors"
	"github.com/tirtahakimpambudhi/restful_api/internal/model/request"
	"github.com/tirtahakimpambudhi/restful_api/internal/model/response"
	"github.com/tirtahakimpambudhi/restful_api/internal/usecase"
	reflecthelper "github.com/tirtahakimpambudhi/restful_api/pkg/helper/reflect"
	"net/url"
	"os"
	"strings"
	"time"
)

// AuthController handles requests related to authentication operations
type AuthController struct {
	usecases *usecase.AuthUsecase
	logger   *log.Logger
}

// NewAuthController creates a new AuthController
func NewAuthController(usecases *usecase.AuthUsecase, logger *log.Logger) *AuthController {
	// Initialize AuthController with provided usecases
	logger.Info().Msg("Initializing AuthController")
	return &AuthController{usecases: usecases, logger: logger}
}

// Login authenticates a user and generates a token
func (controller AuthController) Login(ctx *fiber.Ctx) error {
	controller.logger.Info().Msg("Handling login request")

	// Create a new Auth request
	req := new(request.Auth)

	// Parse the request body into the Auth struct
	if err := ctx.BodyParser(req); err != nil {
		controller.logger.Error().Msgf("Failed to parse request body: %v", err)
		return &response.StandardErrors{Errors: []*response.Error{errorshandler.NewError(errorshandler.BAD_REQUEST, fmt.Sprintf("BAD REQUEST : %s \nREQUEST BODY  \n%s", err.Error(), reflecthelper.KeyValueToString(*req)))}}
	}
	controller.logger.Info().Msgf("Request body parsed successfully: %+v", req)

	// Authenticate the user using the usecase
	res, refreshToken, errors := controller.usecases.Login(ctx.Context(), req)
	if errors != nil {
		controller.logger.Error().Msgf("Authentication failed: %v", errors)
		return errors
	}
	controller.logger.Info().Msg("Authentication successful")
	controller.logger.Info().Msg("Setting refresh token cookie")

	// Set the cookie in the response
	if err := controller.setCookies(ctx, map[string]string{"refresh_token": refreshToken}); err != nil {
		return err
	}

	// Set the response status code
	ctx.Status(res.Status)
	controller.logger.Info().Msgf("Returning response with status: %d", res.Status)

	// Return the response as JSON
	return ctx.JSON(res)
}

// Logout invalidates the user's session
func (controller AuthController) Logout(ctx *fiber.Ctx) error {
	controller.logger.Info().Msg("Handling logout request")

	// Retrieve the refresh token from cookies
	cookie := ctx.Cookies("refresh_token")
	controller.logger.Info().Msg("Retrieved refresh token from cookies")

	// Logout the user using the usecase
	res, errors := controller.usecases.Logout(ctx.Context(), cookie)
	if errors != nil {
		controller.logger.Error().Msgf("Logout failed: %v", errors)
		return errors
	}
	// Clear the cookie refresh token
	ctx.Cookie(&fiber.Cookie{
		Name:    "refresh_token",
		Value:   "",
		Expires: time.Now().Add(-time.Hour),
	})
	controller.logger.Info().Msg("Logout successful")

	// Set the response status code
	ctx.Status(res.Status)
	controller.logger.Info().Msgf("Returning response with status: %d", res.Status)

	// Return the response as JSON
	return ctx.JSON(res)
}

// RefreshToken generates a new token using the refresh token
func (controller AuthController) RefreshToken(ctx *fiber.Ctx) error {
	controller.logger.Info().Msg("Handling refresh token request")

	// Retrieve the refresh token from cookies
	cookie := utils.CopyString(ctx.Cookies("refresh_token"))
	controller.logger.Info().Msg("Retrieved refresh token from cookies")

	// Generate a new token using the usecases
	res, errors := controller.usecases.RefreshToken(ctx.Context(), cookie)
	if errors != nil {
		controller.logger.Error().Msgf("Token refresh failed: %v", errors)
		return errors
	}
	controller.logger.Info().Msg("Token refresh successful")

	// Set the response status code
	ctx.Status(res.Status)
	controller.logger.Info().Msgf("Returning response with status: %d", res.Status)

	// Return the response as JSON
	return ctx.JSON(res)
}

// ResetPassword handles password reset requests
func (controller AuthController) ResetPassword(ctx *fiber.Ctx) error {
	controller.logger.Info().Msg("Handling password reset request")

	// Retrieve the user payload from context locals
	payload, ok := ctx.Locals("users").(*token.Payload)
	if !ok {
		controller.logger.Error().Msg("Failed to convert payload")
		return &response.StandardErrors{Errors: []*response.Error{errorshandler.NewError(errorshandler.INTERNAL_SERVER_ERROR, "Error Converting Payload")}}
	}
	controller.logger.Info().Msg("Payload converted successfully")

	// Create a new ResetPassword request
	req := new(request.ResetPassword)

	// Parse the request body into the ResetPassword struct
	if err := ctx.BodyParser(req); err != nil {
		controller.logger.Error().Msgf("Failed to parse request body: %v", err)
		return &response.StandardErrors{Errors: []*response.Error{errorshandler.NewError(errorshandler.BAD_REQUEST, fmt.Sprintf("BAD REQUEST : %s \nRequest Body \n%s", err.Error(), reflecthelper.KeyValueToString(*req)))}}
	}
	controller.logger.Info().Msgf("Request body parsed successfully: %+v", req)

	// Reset the user's password using the usecase
	res, errors := controller.usecases.ResetPassword(ctx.Context(), payload, req)
	if errors != nil {
		controller.logger.Error().Msgf("Password reset failed: %v", errors)
		return errors
	}
	controller.logger.Info().Msg("Password reset successful")

	// Set the response status code
	ctx.Status(res.Status)
	controller.logger.Info().Msgf("Returning response with status: %d", res.Status)

	// Return the response as JSON
	return ctx.JSON(res)
}

// UpsertRole handles update or insert role for user
func (controller AuthController) UpsertRole(ctx *fiber.Ctx) error {
	controller.logger.Info().Msg("UpsertRole Method Called")

	// Retrieve the user payload from context locals
	//payload, ok := ctx.Locals("users").(*token.Payload)
	//if !ok {
	//	controller.logger.Error().Msg("Failed to convert payload")
	//	return &response.StandardErrors{Errors: []*response.Error{errorshandler.NewError(errorshandler.INTERNAL_SERVER_ERROR, "Error Converting Payload")}}
	//}
	controller.logger.Info().Msg("Payload converted successfully")

	// Create a new ResetPassword request
	req := new(request.UpdateRole)

	// Parse the request body into the ResetPassword struct
	if err := ctx.BodyParser(req); err != nil {
		controller.logger.Error().Msgf("Failed to parse request body: %v", err)
		return &response.StandardErrors{Errors: []*response.Error{errorshandler.NewError(errorshandler.BAD_REQUEST, fmt.Sprintf("BAD REQUEST : %s \nRequest Body \n%s", err.Error(), reflecthelper.KeyValueToString(*req)))}}
	}
	controller.logger.Info().Msgf("Request body parsed successfully: %+v", req)

	// Reset the user's password using the usecase
	res, errors := controller.usecases.UpsertRole(ctx.Context(), req)
	if errors != nil {
		controller.logger.Error().Msgf("Upsert Role failed: %v", errors)
		return errors
	}
	controller.logger.Info().Msg("Upsert and Insert successful")

	// Set the response status code
	ctx.Status(res.Status)
	controller.logger.Info().Msgf("Returning response with status: %d", res.Status)

	// Return the response as JSON
	return ctx.JSON(res)
}
func (controller AuthController) setCookies(ctx *fiber.Ctx, keyValue map[string]string) error {
	// Retrieve the hostname from the client's request URL
	clientURL := ctx.Get("Origin")
	if clientURL == "" && ctx.Get("X-Test-Client") == os.Getenv("SECRET_TEST_CLIENT") {
		clientURL = ctx.BaseURL()
	}

	controller.logger.Info().Str("client_url", clientURL).Msg("Retrieving client URL from the request")

	parsedURL, err := url.Parse(clientURL)

	if err != nil {
		return &response.StandardErrors{Errors: []*response.Error{
			errorshandler.NewError(errorshandler.BAD_REQUEST, "Invalid client URL: "+clientURL),
		}}
	}
	isSecure := parsedURL.Scheme == "https"
	// Extract the domain from the parsed hostname
	domain := parsedURL.Hostname()
	domain = strings.TrimPrefix(domain, "www.") // Remove 'www.' prefix if present

	// Iterate through the key-value pairs and set cookies
	for key, value := range keyValue {
		ctx.Cookie(&fiber.Cookie{
			Name:     key,      // Cookie name
			Value:    value,    // Cookie value
			Domain:   domain,   // Set the cookie domain dynamically
			Path:     "/",      // Set the path to root
			HTTPOnly: true,     // Prevent client-side JavaScript access
			Secure:   isSecure, // Set the Secure flag based on the protocol
			SameSite: "Lax",    // Prevent cross-site request forgery
			MaxAge:   24 * 60 * 60,
		})
	}

	return nil
}
