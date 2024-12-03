package cache

import (
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/tirtahakimpambudhi/restful_api/internal/configs"
)

// RedisConfig holds the configuration for connecting to Redis.
type RedisConfig struct {
	Name     int    `env:"CACHE_DB_NAME"`                     // Redis database name.
	Host     string `env:"CACHE_DB_HOST,required"`            // Redis server hostname.
	Port     int    `env:"CACHE_DB_PORT,required"`            // Redis server port.
	User     string `env:"CACHE_DB_USER"`                     // Redis username.
	Password string `env:"CACHE_DB_PASS"`                     // Redis password.
	MaxCon   int    `env:"CACHE_DB_MAX_CON" envDefault:"100"` // Maximum number of connections.
	MinCon   int    `env:"CACHE_DB_MIN_CON" envDefault:"10"`  // Minimum number of connections.
	MaxTime  int    `env:"CACHE_DB_MAX_TIME" envDefault:"10"` // Maximum idle connection time (in minutes).
	MinTime  int    `env:"CACHE_DB_MIN_TIME" envDefault:"2"`  // Minimum idle connection time (in minutes).
}

// NewConfig initializes a new RedisConfig by loading the configuration.
func NewConfig() (*RedisConfig, error) {
	var config RedisConfig
	// Load configuration values into RedisConfig struct.
	if err := configs.GetConfig().Load(&config); err != nil {
		return nil, err // Return error if loading configuration fails.
	}
	return &config, nil // Return the loaded configuration.
}

// NewClient creates a new Redis client using the configuration.
func (redisConfig *RedisConfig) NewClient() *redis.Client {
	addr := fmt.Sprintf("%s:%d", redisConfig.Host, redisConfig.Port)
	return redis.NewClient(&redis.Options{
		Addr:            addr,                                             // Redis server address.
		Username:        redisConfig.User,                                 // Redis username.
		Password:        redisConfig.Password,                             // Redis password.
		DB:              redisConfig.Name,                                 // Redis database number.
		MinIdleConns:    redisConfig.MinCon,                               // Minimum number of idle connections.
		MaxIdleConns:    redisConfig.MaxCon,                               // Maximum number of idle connections.
		ConnMaxIdleTime: time.Duration(redisConfig.MinTime) * time.Minute, // Maximum idle connection time.
		ConnMaxLifetime: time.Duration(redisConfig.MaxTime) * time.Minute, // Maximum connection lifetime.
	})
}
