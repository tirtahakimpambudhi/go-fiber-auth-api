package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/tirtahakimpambudhi/restful_api/internal/configs/security"
)

// CORS sets up CORS middleware.
func CORS() (fiber.Handler, error) {
	corsConf, err := security.NewCors()
	if err != nil {
		return nil, err
	}
	return cors.New(corsConf.Fiber()), nil
}
