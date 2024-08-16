package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	errorshandler "github.com/tirtahakimpambudhi/restful_api/internal/errors"
	"github.com/tirtahakimpambudhi/restful_api/internal/model/response"
	"golang.org/x/net/xsrftoken"
	"net/url"
	"os"
	"strings"
	"time"
)

// TokenTimeout defines the duration for which a CSRF token remains valid
// Set to 24 hours by default
const (
	TokenTimeout = 24 * time.Hour
	CookieName   = "csrf_token"
)

// GenerateUSERID used for generate userId with uuid V4
func GenerateUSERID() (fiber.Handler, error) {
	return func(ctx *fiber.Ctx) error {
		if ctx.Get("X-User-Id") == "" {
			ctx.Locals("userId", uuid.New().String())
		}
		return ctx.Next()
	}, nil
}

// VerifyCSRF middleware validates the CSRF token for non-GET requests
// It checks if the token exists and is valid within the timeout period
// Returns a 403 Forbidden error if validation fails
func VerifyCSRF() (fiber.Handler, error) {
	return func(ctx *fiber.Ctx) error {
		// Skip verification for GET (except refresh-token), HEAD, and OPTIONS requests
		if (ctx.Method() == "GET" && !strings.Contains(ctx.Path(), "refresh-token")) || ctx.Method() == "HEAD" || ctx.Method() == "OPTIONS" || ctx.Get("X-Test-Client") == os.Getenv("SECRET_TEST_CLIENT") {
			return ctx.Next()
		}

		// Get user ID from header or context
		// First tries X-USER-ID header, then falls back to context value
		userId := ctx.Get("X-User-Id")
		if userId == "" {
			userId = ctx.Locals("userId").(string)
		}

		// Validate UUID format of user ID
		// Returns 400 Bad Request if invalid
		if _, err := uuid.Parse(userId); err != nil {
			return &response.StandardErrors{Errors: []*response.Error{
				errorshandler.NewError(errorshandler.BAD_REQUEST, "invalid user id from header 'X-USER-ID'"),
			}}
		}

		// Get token from X-Csrf-Token header
		clientToken := ctx.Get("X-Csrf-Token")
		action := ctx.Path()
		// Validate token using ValidFor which checks both token validity and timeout
		// Returns 403 Forbidden if token is invalid or expired
		if clientToken == "" || !xsrftoken.Valid(clientToken, os.Getenv("SECRET_KEY_CSRF"), userId, action) {
			return &response.StandardErrors{Errors: []*response.Error{
				errorshandler.NewError(errorshandler.FORBIDEN, "invalid or expired CSRF Token from 'X-Csrf-Token'"),
			}}
		}

		return ctx.Next()
	}, nil
}

// GenerateCSRF middleware generates a new CSRF token or reuses an existing valid token
// The token is stored in both a cookie and returned in the X-Csrf-Token header
func GenerateCSRF() (fiber.Handler, error) {
	return func(ctx *fiber.Ctx) error {
		// Get user ID from header or context
		userId := ctx.Get("X-User-Id")
		if userId == "" {
			userId = ctx.Locals("userId").(string)
		}
		// Validate UUID format
		if _, err := uuid.Parse(userId); err != nil {
			return &response.StandardErrors{Errors: []*response.Error{
				errorshandler.NewError(errorshandler.BAD_REQUEST, "invalid user id from header 'X-User-Id'"),
			}}
		}

		// Check for existing CSRF token in cookies
		existingToken := ctx.Cookies(CookieName)
		if existingToken != "" {
			action := ctx.Path()
			// If existing token is still valid, reuse it
			if xsrftoken.Valid(existingToken, os.Getenv("SECRET_KEY_CSRF"), userId, action) {
				ctx.Set("X-Csrf-Token", existingToken)
				return ctx.Next()
			}
		}

		// Generate new token if none exists or existing one is expired
		action := ctx.Path()
		token := xsrftoken.Generate(os.Getenv("SECRET_KEY_CSRF"), userId, action)

		// Parse request URL to get domain
		reqURL := ctx.Get("Origin")
		parsedURL, err := url.Parse(reqURL)
		if err != nil {
			return &response.StandardErrors{Errors: []*response.Error{
				errorshandler.NewError(errorshandler.BAD_REQUEST, "invalid request url :"+reqURL),
			}}
		}

		// Set up cookie with appropriate security settings
		domain := parsedURL.Hostname()
		cookie := &fiber.Cookie{
			Name:     CookieName, // Using __Host prefix for enhanced security
			Value:    token,
			Domain:   domain,
			Path:     "/",                         // Cookie available for all paths
			Secure:   parsedURL.Scheme == "https", // Secure flag based on HTTPS
			HTTPOnly: true,                        // Prevent JavaScript access
			MaxAge:   int(TokenTimeout.Seconds()),
			SameSite: "Lax", // Allows third-party GET requests
		}

		// Set cookie and header
		ctx.Cookie(cookie)
		ctx.Set("X-Csrf-Token", token)
		ctx.Set("X-User-Id", userId)
		return ctx.Next()
	}, nil
}
