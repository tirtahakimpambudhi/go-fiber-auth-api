package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
)

// Monitor sets up monitoring middleware.
func Monitor() fiber.Handler {
	return monitor.New()
}
