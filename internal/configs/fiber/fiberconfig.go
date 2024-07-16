package fiber

import (
	"encoding/json"
	"errors"
	goJson "github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/tirtahakimpambudhi/restful_api/internal/configs"
	errorshandler "github.com/tirtahakimpambudhi/restful_api/internal/errors"
	"github.com/tirtahakimpambudhi/restful_api/internal/model/response"
	pathhelper "github.com/tirtahakimpambudhi/restful_api/pkg/helper/path"
	"net/http"
	"time"
)

// NewFiberConfig initializes and returns a new FiberConfig instance based on the configuration values loaded from environment variables.
// It sets up SSL configuration, creates necessary directories, and returns the FiberConfig instance or an error.
func NewFiberConfig() (*FiberConfig, error) {
	var (
		fiberConfig FiberConfig
		sslConfig   SSLConfig
	)

	config := configs.GetConfig()

	if err := config.Load(&sslConfig); err != nil {
		return nil, err
	}
	if err := config.Load(&fiberConfig); err != nil {
		return nil, err
	}
	fiberConfig.SSL = &sslConfig
	if err := pathhelper.MakedirFromFieldStruct(fiberConfig); err != nil {
		return nil, err
	}
	return &fiberConfig, nil
}

// SSLConfig represents the SSL configuration for the Fiber server.
type SSLConfig struct {
	Path     string `env:"FIBER_SSL_PATH" envDefault:"resource/ssl"`
	CertFile string `env:"FIBER_SSL_CERT"`
	KeyFile  string `env:"FIBER_SSL_KEY"`
}

// FiberConfig represents the configuration for the Fiber server.
type FiberConfig struct {
	Host              string `env:"FIBER_HOST,required" envDefault:"localhost"`
	Port              string `env:"FIBER_PORT,required" envDefault:"8081"`
	SSL               *SSLConfig
	Prefork           bool   `env:"FIBER_PREFORK" envDefault:"true"`
	StrictRouting     bool   `env:"FIBER_STRICT_ROUTING" envDefault:"true"`
	CaseSensitive     bool   `env:"FIBER_CASE_SENSITIVE" envDefault:"true"`
	BodyLimit         int    `env:"FIBER_BODY_LIMIT" envDefault:"4"`
	ReadTimeout       int    `env:"FIBER_READ_TIMEOUT" envDefault:"4"`
	WriteTimeout      int    `env:"FIBER_WRITE_TIMEOUT" envDefault:"5"`
	ReduceMemoryUsage bool   `env:"FIBER_REDUCE_MEMU" envDefault:"true"`
	JSON              string `env:"FIBER_JSON" envDefault:"json"`
}

// ToFiberAppConfig converts a FiberConfig instance to a fiber.Config instance for Fiber server configuration.
func (fiberConfig *FiberConfig) ToFiberAppConfig() fiber.Config {
	// Create a new fiber.Config instance based on the FiberConfig values.
	config := fiber.Config{
		Prefork:           fiberConfig.Prefork,
		StrictRouting:     fiberConfig.StrictRouting,
		CaseSensitive:     fiberConfig.CaseSensitive,
		BodyLimit:         fiberConfig.BodyLimit * 1024 * 1024,
		ReadTimeout:       time.Duration(fiberConfig.ReadTimeout) * time.Minute,
		WriteTimeout:      time.Duration(fiberConfig.WriteTimeout) * time.Minute,
		ReduceMemoryUsage: fiberConfig.ReduceMemoryUsage,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			var errFiber *fiber.Error
			if errors.As(err, &errFiber) {
				ctx.Status(errFiber.Code)
				code := errorshandler.ConvertStatusCodeToString(errFiber.Code)
				return ctx.JSON(response.StandardErrors{Errors: []*response.Error{errorshandler.NewError(errorshandler.TypeErr(code), errFiber.Message)}})
			}
			var standardErr *response.StandardErrors
			if errors.As(err, &standardErr) {
				ctx.Status(standardErr.Errors[0].Status)
				return ctx.JSON(standardErr)
			}
			ctx.Status(http.StatusInternalServerError)
			return ctx.JSON(response.StandardErrors{Errors: []*response.Error{errorshandler.NewError(errorshandler.INTERNAL_SERVER_ERROR, "Error Internal Server :"+err.Error())}})
		},
	}
	// Set up JSON encoding and decoding based on the specified JSON library.
	switch fiberConfig.JSON {
	case "go-json":
		config.JSONDecoder = goJson.Unmarshal
		config.JSONEncoder = goJson.Marshal
	default:
		config.JSONDecoder = json.Unmarshal
		config.JSONEncoder = json.Marshal
	}
	// Return the fiber.Config instance.
	return config
}
