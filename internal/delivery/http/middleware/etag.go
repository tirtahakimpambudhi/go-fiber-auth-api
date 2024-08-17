package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/etag"
)

// ETag sets up ETag middleware.
func ETag() (fiber.Handler, error) {
	return etag.New(), nil
}
