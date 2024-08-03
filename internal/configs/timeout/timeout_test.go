package timeout_test

import (
	"context"
	"github.com/tirtahakimpambudhi/restful_api/internal/configs/timeout"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// Test NewConfig function for different scenarios
func TestNewConfig(t *testing.T) {
	t.Run("Happy Case - Valid Timeout Values", func(t *testing.T) {
		// Set up valid environment variables
		os.Setenv("CACHE_TIMEOUT", "5")
		os.Setenv("DB_TIMEOUT", "10")
		os.Setenv("DOWN_STREAM_TIMEOUT", "15")

		// Call NewConfig
		config, err := timeout.NewConfig()

		// Verify no error occurred
		require.NoError(t, err)
		// Verify the values are parsed correctly
		require.Equal(t, 5*time.Minute, config.CacheTimeout)
		require.Equal(t, 10*time.Minute, config.DatabaseTimeout)
		require.Equal(t, 15*time.Minute, config.DownstreamTimeout)
	})

	t.Run("Edge Case - Invalid CACHE_TIMEOUT", func(t *testing.T) {
		// Set up invalid CACHE_TIMEOUT
		os.Setenv("CACHE_TIMEOUT", "invalid")
		os.Setenv("DB_TIMEOUT", "10")
		os.Setenv("DOWN_STREAM_TIMEOUT", "15")

		// Call NewConfig
		_, err := timeout.NewConfig()

		// Verify an error occurred
		require.Error(t, err)
	})

	t.Run("Edge Case - Invalid DB_TIMEOUT", func(t *testing.T) {
		// Set up invalid DB_TIMEOUT
		os.Setenv("CACHE_TIMEOUT", "5")
		os.Setenv("DB_TIMEOUT", "invalid")
		os.Setenv("DOWN_STREAM_TIMEOUT", "15")

		// Call NewConfig
		_, err := timeout.NewConfig()

		// Verify an error occurred
		require.Error(t, err)
	})

	t.Run("Edge Case - Invalid DOWN_STREAM_TIMEOUT", func(t *testing.T) {
		// Set up invalid DOWN_STREAM_TIMEOUT
		os.Setenv("CACHE_TIMEOUT", "5")
		os.Setenv("DB_TIMEOUT", "10")
		os.Setenv("DOWN_STREAM_TIMEOUT", "invalid")

		// Call NewConfig
		_, err := timeout.NewConfig()

		// Verify an error occurred
		require.Error(t, err)
	})

	t.Run("Edge Case - Missing Environment Variables", func(t *testing.T) {
		// Unset all environment variables
		os.Unsetenv("CACHE_TIMEOUT")
		os.Unsetenv("DB_TIMEOUT")
		os.Unsetenv("DOWN_STREAM_TIMEOUT")

		// Call NewConfig
		_, err := timeout.NewConfig()

		// Verify an error occurred due to missing env vars
		require.Error(t, err)
	})
}

// Test CreateCacheTimeout, CreateDatabaseTimeout, CreateDownstreamTimeout
func TestConfigTimeouts(t *testing.T) {
	// Set up a valid config
	config := &timeout.Config{
		CacheTimeout:      2 * time.Minute,
		DatabaseTimeout:   3 * time.Minute,
		DownstreamTimeout: 4 * time.Minute,
	}

	t.Run("CreateCacheTimeout - Happy Case", func(t *testing.T) {
		ctx := context.Background()
		newCtx, cancel := config.CreateCacheTimeout(ctx)
		defer cancel()

		// Verify that the context has the correct deadline
		deadline, ok := newCtx.Deadline()
		require.True(t, ok)
		require.WithinDuration(t, time.Now().Add(config.CacheTimeout), deadline, time.Second)
	})

	t.Run("CreateDatabaseTimeout - Happy Case", func(t *testing.T) {
		ctx := context.Background()
		newCtx, cancel := config.CreateDatabaseTimeout(ctx)
		defer cancel()

		// Verify that the context has the correct deadline
		deadline, ok := newCtx.Deadline()
		require.True(t, ok)
		require.WithinDuration(t, time.Now().Add(config.DatabaseTimeout), deadline, time.Second)
	})

	t.Run("CreateDownstreamTimeout - Happy Case", func(t *testing.T) {
		ctx := context.Background()
		newCtx, cancel := config.CreateDownstreamTimeout(ctx)
		defer cancel()

		// Verify that the context has the correct deadline
		deadline, ok := newCtx.Deadline()
		require.True(t, ok)
		require.WithinDuration(t, time.Now().Add(config.DownstreamTimeout), deadline, time.Second)
	})

	t.Run("Edge Case - Zero Timeout Values", func(t *testing.T) {
		// Set up a config with zero timeout values
		config := &timeout.Config{
			CacheTimeout:      0,
			DatabaseTimeout:   0,
			DownstreamTimeout: 0,
		}

		ctx := context.Background()

		// Create cache timeout context
		newCtx, cancel := config.CreateCacheTimeout(ctx)
		defer cancel()

		// Verify that no deadline is set when timeout is zero
		timeDeadline, _ := newCtx.Deadline()
		require.WithinDuration(t, time.Now().Add(config.DownstreamTimeout), timeDeadline, time.Second)

		// Create database timeout context
		newCtx, cancel = config.CreateDatabaseTimeout(ctx)
		defer cancel()

		// Verify that no deadline is set when timeout is zero
		timeDeadline, _ = newCtx.Deadline()
		require.WithinDuration(t, time.Now().Add(config.DownstreamTimeout), timeDeadline, time.Second)

		// Create downstream timeout context
		newCtx, cancel = config.CreateDownstreamTimeout(ctx)
		defer cancel()

		// Verify that no deadline is set when timeout is zero
		timeDeadline, _ = newCtx.Deadline()
		require.WithinDuration(t, time.Now().Add(config.DownstreamTimeout), timeDeadline, time.Second)
	})
}
