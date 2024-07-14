package configs

import (
	"errors"
	"log"
	"sync"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	path_helper "github.com/tirtahakimpambudhi/restful_api/pkg/helper/path"
)

// Config struct holds the application configuration.
type Config struct {
	Filename string
}

// Load loads and parses environment variables into the provided struct(s).
// values must contain at least one argument.
func (c *Config) Load(values ...any) error {
	if len(values) == 0 {
		return errors.New("Load arguments must be filled with at least one value")
	}
	for _, value := range values {
		if err := env.Parse(value); err != nil {
			return err
		}
	}
	return nil
}

var (
	ConfigFile string = ".env"
	instance   Config
	once       sync.Once
)

// GetConfig loads the .env file if it exists, and returns the Config instance.
func GetConfig() *Config {
	once.Do(func() {
		instance = Config{
			Filename: path_helper.AddWorkdirToSomePath(ConfigFile),
		}
		if err := godotenv.Load(instance.Filename); err != nil {
			log.Printf("No .env file found, using environment variables")
		}
	})
	return &instance
}
