package security

import (
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/tirtahakimpambudhi/restful_api/internal/configs"
)

// Cors configuration for cross origins resource sharing
type Cors struct {
	AllowHeaders  string `env:"CORS_ALLOW_HEADERS,required"`
	ExposeHeaders string `env:"CORS_EXPOSE_HEADERS,required"`
	AllowMethods  string `env:"CORS_ALLOW_METHODS,required"`
	AllowOrigins  string `env:"CORS_ALLOW_ORIGINS,required"`
	Credentials   bool   `env:"CORS_ALLOW_CREDENTIALS,required"`
}

// Fiber convert config cors to cors config fiber middleware
func (c Cors) Fiber() cors.Config {
	return cors.Config{
		AllowHeaders:     c.AllowHeaders,
		AllowMethods:     c.AllowMethods,
		AllowOrigins:     c.AllowOrigins,
		AllowCredentials: c.Credentials,
		ExposeHeaders:    c.ExposeHeaders,
	}
}

func NewCors() (*Cors, error) {
	var corsConf Cors
	err := configs.GetConfig().Load(&corsConf)
	if err != nil {
		return nil, err
	}
	return &corsConf, nil
}
