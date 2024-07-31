package sqlconfig

import (
	"github.com/tirtahakimpambudhi/restful_api/internal/configs"
)

// SqlConfig holds the configuration for SQL database connection.
type SqlConfig struct {
	Driver   string `env:"DB_DRIVER,required"`          // Database driver (e.g., postgres, mysql)
	Protocol string `env:"DB_PROTOCOL,required"`        // Protocol for the database connection
	Name     string `env:"DB_NAME,required"`            // Database name
	Host     string `env:"DB_HOST,required"`            // Database host address
	Port     string `env:"DB_PORT,required"`            // Port for the database connection
	User     string `env:"DB_USER,required"`            // Database user
	Password string `env:"DB_PASS,required"`            // Password for the database user
	MaxCon   int    `env:"DB_MAX_CON" envDefault:"100"` // Maximum number of connections
	MinCon   int    `env:"DB_MIN_CON" envDefault:"10"`  // Minimum number of connections
	MaxTime  int    `env:"DB_MAX_TIME" envDefault:"30"` // Maximum connection time in seconds
	MinTime  int    `env:"DB_MIN_TIME" envDefault:"5"`  // Minimum connection time in seconds
}

// NewConfig creates a new SqlConfig by loading configuration from environment variables.
func NewConfig() (*SqlConfig, error) {
	var config SqlConfig

	// Load the configuration into the config object.
	if err := configs.GetConfig().Load(&config); err != nil {
		// Return nil and the error if loading configuration fails.
		return nil, err
	}

	// Return the loaded configuration and no error.
	return &config, nil
}
