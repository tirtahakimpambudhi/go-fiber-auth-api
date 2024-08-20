package middleware

import (
	"github.com/gofiber/fiber/v2"
)

// Setup initializes all middleware for the Fiber app.
func Setup(app *fiber.App) error {
	// List of middleware to set up
	middlewares := []func() (fiber.Handler, error){
		CORS, GenerateUSERID, Limiter, GenerateCSRF, VerifyCSRF, HealthCheck, ETag, Swagger,
	}
	// Loop through each middleware and apply it to the app
	for _, mw := range middlewares {
		handler, err := mw()
		if err != nil {
			return err
		}
		app.Use(handler)
	}
	return nil
}
