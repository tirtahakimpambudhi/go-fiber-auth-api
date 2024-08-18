package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
)

// HealthCheck sets up Health Check middleware.
func HealthCheck() (fiber.Handler, error) {
	return healthcheck.New(), nil
}
