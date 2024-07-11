package cache_test

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tirtahakimpambudhi/restful_api/internal/configs/cache"
)

func TestNewConfig_Success(t *testing.T) {
	os.Setenv("CACHE_DB_NAME", "1")
	os.Setenv("CACHE_DB_HOST", "localhost")
	os.Setenv("CACHE_DB_PORT", "5432")
	os.Setenv("CACHE_DB_USER", "user")
	os.Setenv("CACHE_DB_PASS", "password")

	defer os.Unsetenv("CACHE_DB_NAME")
	defer os.Unsetenv("CACHE_DB_HOST")
	defer os.Unsetenv("CACHE_DB_PORT")
	defer os.Unsetenv("CACHE_DB_USER")
	defer os.Unsetenv("CACHE_DB_PASS")

	config, err := cache.NewConfig()

	require.NoError(t, err)
	require.NotNil(t, config)

	require.Equal(t, 1, config.Name)
	require.Equal(t, "localhost", config.Host)
	require.Equal(t, 5432, config.Port)
	require.Equal(t, "user", config.User)
	require.Equal(t, "password", config.Password)
}

func TestNewConfig_MissingEnvVars(t *testing.T) {
	os.Clearenv()

	config, err := cache.NewConfig()

	require.Error(t, err)
	require.Nil(t, config)
}

// Creates a new Redis client with valid configuration
func TestCreatesNewRedisClientWithValidConfig(t *testing.T) {
	os.Setenv("CACHE_DB_NAME", "1")
	os.Setenv("CACHE_DB_HOST", "localhost")
	os.Setenv("CACHE_DB_PORT", "6379")
	os.Setenv("CACHE_DB_USER", "user")
	os.Setenv("CACHE_DB_PASS", "password")

	defer os.Unsetenv("CACHE_DB_NAME")
	defer os.Unsetenv("CACHE_DB_HOST")
	defer os.Unsetenv("CACHE_DB_PORT")
	defer os.Unsetenv("CACHE_DB_USER")
	defer os.Unsetenv("CACHE_DB_PASS")
	config := &cache.RedisConfig{
		Name:     1,
		Host:     "localhost",
		Port:     6379,
		User:     "user",
		Password: "password",
		MaxCon:   100,
		MinCon:   10,
		MaxTime:  10,
		MinTime:  2,
	}

	client := config.NewClient()

	assert.NotNil(t, client)
	assert.Equal(t, "localhost:6379", client.Options().Addr)
	assert.Equal(t, "user", client.Options().Username)
	assert.Equal(t, "password", client.Options().Password)
	assert.Equal(t, 1, client.Options().DB)
	assert.Equal(t, 10, client.Options().MinIdleConns)
	assert.Equal(t, 100, client.Options().MaxIdleConns)
	// assert.Equal(t, 2*time.Minute, client.Options().ConnMaxIdleTime)
	// assert.Equal(t, 100*time.Minute, client.Options().ConnMaxLifetime)
}

// Handles missing or empty host and port values
func TestHandlesMissingOrEmptyHostAndPort(t *testing.T) {
	config := &cache.RedisConfig{
		Name:     0,
		Host:     "",
		Port:     0,
		User:     "user",
		Password: "password",
		MaxCon:   100,
		MinCon:   10,
		MaxTime:  10,
		MinTime:  2,
	}

	client := config.NewClient()

	assert.NotNil(t, client)
	assert.Equal(t, ":0", client.Options().Addr)
}
