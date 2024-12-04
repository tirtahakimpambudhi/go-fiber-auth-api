package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/storage/redis/v3"
	"github.com/tirtahakimpambudhi/restful_api/internal/configs/cache"
	errorshandler "github.com/tirtahakimpambudhi/restful_api/internal/errors"
	"github.com/tirtahakimpambudhi/restful_api/internal/model/response"
	"net/http"
	"time"
)

// Limiter sets up request rate limiting middleware.
func Limiter() (fiber.Handler, error) {
	var storage *redis.Storage
	// Load Redis cache configuration
	config, err := cache.NewConfig()
	if err != nil {
		return nil, err
	}

	storage = redis.New(redis.Config{
		Host:     config.Host,
		Port:     config.Port,
		Username: config.User,
		Password: config.Password,
		Database: config.Name,
		Reset:    false,
	})

	// Set up limiter middleware
	return limiter.New(limiter.Config{
		Next: func(c *fiber.Ctx) bool {
			// Allow requests from localhost without limit
			return c.IP() == "127.0.0.1"
		},
		Max:        20,               // Maximum number of requests
		Expiration: 30 * time.Second, // Time window for rate limiting
		KeyGenerator: func(c *fiber.Ctx) string {
			// Generate a key based on the client's IP address
			return c.Get("x-forwarded-for")
		},
		LimitReached: func(ctx *fiber.Ctx) error {
			// Respond with "Too Many Requests" error if limit is reached
			ctx.Status(http.StatusTooManyRequests)
			return ctx.JSON(&response.StandardErrors{
				Errors: []*response.Error{
					errorshandler.NewError(errorshandler.TO_MANY_REQUEST, "Error: Too Many Requests, wait 30 minutes"),
				},
			})
		},
		Storage: storage, // Use Redis storage
	}), nil
}
