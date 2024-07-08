package bootstrap

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/tirtahakimpambudhi/restful_api/internal/configs/cache"
	casbinconfig "github.com/tirtahakimpambudhi/restful_api/internal/configs/casbin"
	fiberconfig "github.com/tirtahakimpambudhi/restful_api/internal/configs/fiber"
	"github.com/tirtahakimpambudhi/restful_api/internal/configs/hash"
	loggerconfig "github.com/tirtahakimpambudhi/restful_api/internal/configs/logger"
	"github.com/tirtahakimpambudhi/restful_api/internal/configs/orm"
	sqlconfig "github.com/tirtahakimpambudhi/restful_api/internal/configs/sql"
	"github.com/tirtahakimpambudhi/restful_api/internal/configs/timeout"
	tokenconfig "github.com/tirtahakimpambudhi/restful_api/internal/configs/token"
	"gorm.io/gorm"
)

// App struct holds all the application configurations and dependencies.
type App struct {
	Redis          *cache.RedisConfig     // Redis configuration
	CasbinEnforcer *casbin.Enforcer       // Casbin enforcer for policy enforcement
	Gorm           *gorm.DB               // GORM database instance
	FiberServer    *fiberconfig.Fiber     // Fiber server configuration
	Hash           *hash.Argon2           // Argon2 hashing configuration
	Logger         *loggerconfig.Logger   // Logger configuration
	SQL            *sqlconfig.SqlConfig   // SQL configuration
	Timeout        *timeout.Config        // Timeout configuration
	Token          *tokenconfig.JWTToken  // JWT token configuration
	Secret         *tokenconfig.SecretKey // Secret key for JWT
}

// configLoader is a generic function that loads a configuration using the provided function.
func configLoader[T any](configFunc func() (*T, error)) (*T, error) {
	return configFunc() // Execute the configuration function and return the result
}

// New initializes and returns a new App instance with all configurations loaded.
func New() (*App, error) {
	// Load Logger configuration
	logger, loggerErr := configLoader(loggerconfig.NewLogger)
	if loggerErr != nil {
		return nil, loggerErr // Return error if loading Logger config fails
	}

	logger.App.Info().Msg("Starting application initialization...")

	// Load Fiber configuration
	fiberConfig, fiberErr := configLoader(fiberconfig.NewFiberConfig)
	if fiberErr != nil {
		logger.App.Error().Msgs("Failed to load Fiber config:", fiberErr)
		return nil, fiberErr // Return error if loading Fiber config fails
	}
	fiberServer := fiberconfig.NewFiber(fiberConfig) // Create a new Fiber server with the config
	logger.App.Info().Msg("Successfully loaded Fiber configuration")

	// Load Redis configuration
	redisConfig, redisErr := configLoader(cache.NewConfig)
	if redisErr != nil {
		logger.App.Error().Msgs("Failed to load Redis config:", redisErr)
		return nil, redisErr // Return error if loading Redis config fails
	}
	logger.App.Info().Msg("Successfully loaded Redis configuration")

	// Load Argon2 hashing configuration
	argon2, hashErr := configLoader(hash.NewHashArgon2)
	if hashErr != nil {
		logger.App.Error().Msgs("Failed to load Argon2 hash config:", hashErr)
		return nil, hashErr // Return error if loading Argon2 config fails
	}
	logger.App.Info().Msg("Successfully loaded Argon2 hash configuration")

	// Load SQL configuration
	sqlConfig, sqlErr := configLoader(sqlconfig.NewConfig)
	if sqlErr != nil {
		logger.App.Error().Msgs("Failed to load SQL config:", sqlErr)
		return nil, sqlErr // Return error if loading SQL config fails
	}
	logger.App.Info().Msg("Successfully loaded SQL configuration")

	// Load Timeout configuration
	timeoutConfig, timeoutErr := configLoader(timeout.NewConfig)
	if timeoutErr != nil {
		logger.App.Error().Msgs("Failed to load Timeout config:", timeoutErr)
		return nil, timeoutErr // Return error if loading Timeout config fails
	}
	logger.App.Info().Msg("Successfully loaded Timeout configuration")

	// Load JWT token configuration and secret key
	jwtToken, key, tokenErr := tokenconfig.NewJWTToken()
	if tokenErr != nil {
		logger.App.Error().Msgs("Failed to load JWT token and secret key:", tokenErr)
		return nil, tokenErr // Return error if loading JWT token config fails
	}
	logger.App.Info().Msg("Successfully loaded JWT token configuration and secret key")

	// Initialize GORM (ORM) database connection
	gormDB, gormErr := orm.NewGorm()
	if gormErr != nil {
		logger.App.Error().Msgs("Failed to initialize GORM:", gormErr)
		return nil, gormErr // Return error if initializing GORM fails
	}
	logger.App.Info().Msg("Successfully initialized GORM")

	// Initialize Casbin middleware and enforcer for authorization
	enforcer, casbinErr := casbinconfig.NewCasbin(gormDB, func(ctx *fiber.Ctx) string {
		payload := ctx.Locals("users").(*tokenconfig.Payload)
		logger.App.Debug().Msg(fmt.Sprintf("Retrieved user payload: %v", payload))
		return ""
	})
	if casbinErr != nil {
		logger.App.Error().Msgs("Failed to initialize Casbin enforcer:", casbinErr)
		return nil, casbinErr // Return error if initializing Casbin fails
	}
	logger.App.Info().Msg("Successfully initialized Casbin enforcer")

	logger.App.Info().Msg("Application initialized successfully")

	// Return a new App instance with all configurations and dependencies
	return &App{
		Redis:          redisConfig,   // Assign Redis config
		FiberServer:    fiberServer,   // Assign Fiber server config
		Gorm:           gormDB,        // Assign GORM instance
		Hash:           argon2,        // Assign Argon2 hash config
		Logger:         logger,        // Assign Logger config
		SQL:            sqlConfig,     // Assign SQL config
		Timeout:        timeoutConfig, // Assign Timeout config
		Token:          jwtToken,      // Assign JWT token config
		Secret:         key,           // Assign secret key for JWT
		CasbinEnforcer: enforcer,      // Assign Casbin enforcer
	}, nil
}
