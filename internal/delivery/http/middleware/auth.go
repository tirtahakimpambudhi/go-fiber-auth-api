package middleware

import (
	"errors"
	"github.com/casbin/casbin/v2"
	"github.com/gofiber/fiber/v2"
	tokenconfig "github.com/tirtahakimpambudhi/restful_api/internal/configs/token"
	errorshandler "github.com/tirtahakimpambudhi/restful_api/internal/errors"
	"github.com/tirtahakimpambudhi/restful_api/internal/model/response"
	"net/http"
	"slices"
	"strings"
)

// NewAuthenticationToken returns a middleware handler function that verifies JWT tokens
// and handles any token-related errors. It uses the provided JWT token configuration and secret key.
func NewAuthenticationToken(jwtToken *tokenconfig.JWTToken, secretKey string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// Extract the Authorization header from the request.
		authHeader := ctx.Get("Authorization")
		if authHeader == "" {
			// If the Authorization header is missing, return a bad request error.
			ctx.Status(http.StatusBadRequest)
			return ctx.JSON(response.StandardErrors{
				Errors: []*response.Error{
					errorshandler.NewError(errorshandler.BAD_REQUEST, "Error Token: Missing or Malformed Token"),
				},
			})
		}

		// Extract the token string from the Authorization header.
		tokenStr := authHeader[len("Bearer "):]

		// Verify the token using the provided secret key.
		payload, err := jwtToken.VerifyToken(secretKey, tokenStr)
		if err != nil {
			var jwtError *tokenconfig.TokenError
			// Check if the error is a TokenError.
			if errors.As(err, &jwtError) {
				// Handle different types of token errors.
				typeErr := jwtError.TypeError()

				switch {
				case errors.Is(typeErr, tokenconfig.ErrTokenMalformed):
					// Handle malformed token error.
					ctx.Status(http.StatusBadRequest)
					return ctx.JSON(response.StandardErrors{
						Errors: []*response.Error{
							errorshandler.NewError(errorshandler.BAD_REQUEST, "Error Token: "+jwtError.Error()),
						},
					})
				case errors.Is(typeErr, tokenconfig.ErrTokenExpired):
					// Handle expired token error.
					ctx.Status(http.StatusForbidden)
					return ctx.JSON(response.StandardErrors{
						Errors: []*response.Error{
							errorshandler.NewError(errorshandler.FORBIDEN, "Error Token: "+jwtError.Error()),
						},
					})
				case errors.Is(typeErr, tokenconfig.ErrTokenSignatureInvalid):
					// Handle invalid token signature error.
					ctx.Status(http.StatusBadRequest)
					return ctx.JSON(response.StandardErrors{
						Errors: []*response.Error{
							errorshandler.NewError(errorshandler.BAD_REQUEST, "Error Token: "+jwtError.Error()),
						},
					})
				case errors.Is(typeErr, tokenconfig.ErrInvalidToken):
					// Handle generic invalid token error.
					ctx.Status(http.StatusBadRequest)
					return ctx.JSON(response.StandardErrors{
						Errors: []*response.Error{
							errorshandler.NewError(errorshandler.BAD_REQUEST, "Error Token: "+jwtError.Error()),
						},
					})
				default:
					// Handle any other server errors related to tokens.
					ctx.Status(http.StatusInternalServerError)
					return ctx.JSON(response.StandardErrors{
						Errors: []*response.Error{
							errorshandler.NewError(errorshandler.INTERNAL_SERVER_ERROR, "Error Token: "+jwtError.Error()),
						},
					})
				}
			}
		}

		// Continue to the next handler if the token is valid.
		ctx.Locals("users", payload)
		return ctx.Next()
	}
}

// NewAuthorizationById sets up Casbin authorization middleware by user ID.
func NewAuthorizationById(middleware *casbin.Enforcer, permission string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {

		// Get the ID from the URL path
		pathID := ctx.Params("id")
		// Get the user payload from the context
		payload := ctx.Locals("users").(*tokenconfig.Payload)
		// Allow access if the path ID matches the user's ID
		if pathID == payload.ID.String() {
			return ctx.Next()
		}

		var authorized bool
		var err error

		// Check if the permission contains ':'
		if strings.Contains(permission, ":") {
			// Split the permission into subject and action
			parts := strings.Split(permission, ":")
			if len(parts) != 2 {
				return &response.StandardErrors{Errors: []*response.Error{errorshandler.NewError(errorshandler.FORBIDEN, "Invalid permission format")}}
			}

			subject := parts[0]
			action := parts[1]

			// Check if the user has the required permission with subject and action
			authorized, err = middleware.Enforce(payload.Email, subject, action)
		} else {
			return &response.StandardErrors{Errors: []*response.Error{errorshandler.NewError(errorshandler.INTERNAL_SERVER_ERROR, "Permissions must be contains ':' Actual '"+permission+"'")}}
		}

		// Handle error
		if err != nil {
			return &response.StandardErrors{Errors: []*response.Error{errorshandler.NewError(errorshandler.FORBIDEN, err.Error())}}
		}

		// Check if the user is authorized
		if authorized {
			return ctx.Next()
		} else {
			return &response.StandardErrors{Errors: []*response.Error{errorshandler.NewError(errorshandler.FORBIDEN, "Access denied")}}
		}
	}
}

// // NewAuthorization sets up Casbin authorization middleware
func NewAuthorization(middleware *casbin.Enforcer, permission string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// Get the user payload from the context
		payload := ctx.Locals("users").(*tokenconfig.Payload)

		var authorized bool
		var err error

		// Check if the permission contains ':'
		if strings.Contains(permission, ":") {
			// Split the permission into subject and action
			parts := strings.Split(permission, ":")
			if len(parts) != 2 {
				return &response.StandardErrors{Errors: []*response.Error{errorshandler.NewError(errorshandler.FORBIDEN, "Invalid permission format")}}
			}

			subject := parts[0]
			action := parts[1]

			// Check if the user has the required permission with subject and action
			authorized, err = middleware.Enforce(payload.Email, subject, action)
		} else {
			roles, err := middleware.GetRolesForUser(payload.Email)
			if err != nil {
				return &response.StandardErrors{Errors: []*response.Error{errorshandler.NewError(errorshandler.INTERNAL_SERVER_ERROR, "Permissions must be contains ':' Actual '"+permission+"'")}}
			}
			authorized = slices.Contains(roles, permission)
		}
		// Handle error
		if err != nil {
			return &response.StandardErrors{Errors: []*response.Error{errorshandler.NewError(errorshandler.FORBIDEN, err.Error())}}
		}

		// Check if the user is authorized
		if authorized {
			return ctx.Next()
		} else {
			return &response.StandardErrors{Errors: []*response.Error{errorshandler.NewError(errorshandler.FORBIDEN, "Access denied")}}
		}
	}
}
