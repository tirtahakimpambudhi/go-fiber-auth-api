package configs_test

import (
	"os"
	"testing"
	"github.com/tirtahakimpambudhi/restful_api/internal/configs"
	"github.com/stretchr/testify/assert"
)

type EnvConfig struct {
	Port string `env:"PORT"`
}

func TestConfig_Load_WithEnvFile(t *testing.T) {
	// Setup: Create a temporary .env file
	envContent := "PORT=8080\n"
	envFile := ".env"
	err := os.WriteFile(envFile, []byte(envContent), 0644)
	assert.NoError(t, err)
	defer os.Remove(envFile)

	// Load Config
	config := configs.GetConfig()
	envConfig := &EnvConfig{}
	err = config.Load(envConfig)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, "8080", envConfig.Port)
}

func TestConfig_Load_WithoutEnvFile(t *testing.T) {
	// Ensure .env file does not exist
	envFile := ".env"
	os.Remove(envFile)

	// Set environment variable directly
	os.Setenv("PORT", "9090")
	defer os.Unsetenv("PORT")

	// Load Config
	config := configs.GetConfig()
	envConfig := &EnvConfig{}
	err := config.Load(envConfig)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, "9090", envConfig.Port)
}

func TestConfig_Load_Fail_WithoutPointerStruct(t *testing.T) {
	// Ensure .env file does not exist
	envFile := ".env"
	os.Remove(envFile)

	// Set environment variable directly
	os.Setenv("PORT", "9090")
	defer os.Unsetenv("PORT")

	// Load Config
	config := configs.GetConfig()
	envConfig := EnvConfig{}
	err := config.Load(envConfig)

	// Assertions
	assert.Error(t, err)
}

func TestConfig_Load_Fail_WithoutEnvConfig(t *testing.T) {
	// Ensure .env file does not exist
	envFile := ".env"
	os.Remove(envFile)

	// Set environment variable directly
	os.Setenv("PORT", "9090")
	defer os.Unsetenv("PORT")

	// Load Config
	config := configs.GetConfig()
	err := config.Load()

	// Assertions
	assert.Error(t, err)
}
