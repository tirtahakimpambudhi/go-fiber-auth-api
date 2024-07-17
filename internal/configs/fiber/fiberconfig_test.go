// Converts FiberConfig to fiber.Config correctly
package fiber_test

import (
	"github.com/gofiber/fiber/v2"
	fiber2 "github.com/tirtahakimpambudhi/restful_api/internal/configs/fiber"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestInitialize_ConfigFiber(t *testing.T) {
	fiberConfig, err := fiber2.NewFiberConfig()
	require.NoError(t, err)
	require.NotEmpty(t, fiberConfig)
	require.NotEmpty(t, fiberConfig.ToFiberAppConfig())
}

func TestConverts_FiberConfig_To_FiberConfig_Correctly(t *testing.T) {
	fiberConfig := &fiber2.FiberConfig{
		Prefork:           true,
		StrictRouting:     true,
		CaseSensitive:     true,
		BodyLimit:         4,
		ReadTimeout:       5,
		WriteTimeout:      5,
		ReduceMemoryUsage: true,
		JSON:              "go-json",
	}

	expectedConfig := fiber.Config{
		Prefork:           true,
		StrictRouting:     true,
		CaseSensitive:     true,
		ETag:              true,
		BodyLimit:         4 * 1024 * 1024,
		ReadTimeout:       5 * time.Minute,
		WriteTimeout:      5 * time.Minute,
		ReduceMemoryUsage: true,
	}

	actualConfig := fiberConfig.ToFiberAppConfig()

	// match all config
	require.Equal(t, expectedConfig.CaseSensitive, actualConfig.CaseSensitive)
	require.Equal(t, expectedConfig.ReduceMemoryUsage, actualConfig.ReduceMemoryUsage)
	require.Equal(t, expectedConfig.BodyLimit, actualConfig.BodyLimit)
	require.Equal(t, expectedConfig.StrictRouting, actualConfig.StrictRouting)
	require.Equal(t, expectedConfig.WriteTimeout, actualConfig.WriteTimeout)
	require.Equal(t, expectedConfig.ReadTimeout, actualConfig.ReadTimeout)
}
