// Successfully loads configuration from environment variables
package loggerconfig_test

import (
	"github.com/stretchr/testify/require"
	"github.com/tirtahakimpambudhi/restful_api/internal/configs/logger"
	"os"
	"path"
	"testing"
)

func NewTestLoggerConfig(t *testing.T) (*loggerconfig.LoggerConfig, func()) {
	os.Setenv("LOG_PATH", "/var/logs")
	os.Setenv("LOG_MAX_SIZE", "10")
	os.Setenv("LOG_MAX_BACKUP", "5")
	os.Setenv("LOG_MAX_SIZE_ROTATE", "20")
	os.Setenv("LOG_TIME_FORMAT", "2006-01-02")
	os.Setenv("LOG_COLOR_OUTPUT", "true")
	os.Setenv("LOG_QUOTE_STR", "true")
	os.Setenv("LOG_END_WITH_MESSAGE", "true")
	closeFunc := func() {
		os.Unsetenv("LOG_PATH")
		os.Unsetenv("LOG_MAX_SIZE")
		os.Unsetenv("LOG_MAX_BACKUP")
		os.Unsetenv("LOG_MAX_SIZE_ROTATE")
		os.Unsetenv("LOG_TIME_FORMAT")
		os.Unsetenv("LOG_COLOR_OUTPUT")
		os.Unsetenv("LOG_QUOTE_STR")
		os.Unsetenv("LOG_END_WITH_MESSAGE")
	}

	config, err := loggerconfig.NewLoggerConfig()
	require.NoError(t, err)
	return config, closeFunc
}

func TestNewLoggerConfigLoadsEnvVars(t *testing.T) {
	//Setup
	config, close := NewTestLoggerConfig(t)
	defer close()
	defer func() {
		wd, err := os.Getwd()
		require.NoError(t, err)
		require.NoError(t, os.RemoveAll(path.Join(wd, config.LogPath, "..")))
	}()
	// Match expect to actual
	require.Equal(t, "/var/logs", config.LogPath)
	require.Equal(t, 10, config.MaxSize)
	require.Equal(t, 5, config.MaxBackup)
	require.Equal(t, 20, config.MaxSizeRotate)
	require.Equal(t, "2006-01-02", config.TimeFormat)
	require.True(t, config.ColorOutput)
	require.True(t, config.QuoteString)
	require.True(t, config.EndWithMessage)
}

// No environment variables set for LoggerConfig fields
func TestNewLoggerConfigNoEnvVars(t *testing.T) {
	os.Clearenv()

	config, err := loggerconfig.NewLoggerConfig()

	require.Error(t, err)
	require.Nil(t, config)
}

func TestNewFileWriterWithRotate_CorrectTimeFormat(t *testing.T) {
	//Setup
	config, close := NewTestLoggerConfig(t)
	defer close()
	defer func() {
		wd, err := os.Getwd()
		require.NoError(t, err)
		require.NoError(t, os.RemoveAll(path.Join(wd, config.LogPath, "..")))
	}()

	fileWriter := config.NewFileWriterWithRotate("app.log")

	require.Equal(t, "2006-01-02", fileWriter.TimeFormat)
}

func TestNewFileWriter_CorrectTimeFormat(t *testing.T) {
	//Setup
	config, close := NewTestLoggerConfig(t)
	defer close()
	defer func() {
		wd, err := os.Getwd()
		require.NoError(t, err)
		require.NoError(t, os.RemoveAll(path.Join(wd, config.LogPath, "..")))
	}()

	fileWriter := config.NewFileWriter("app.log")

	require.Equal(t, "2006-01-02", fileWriter.TimeFormat)
}

func TestNewConsoleWriter_CorrectTimeFormat(t *testing.T) {
	//Setup
	config, close := NewTestLoggerConfig(t)
	defer close()
	defer func() {
		wd, err := os.Getwd()
		require.NoError(t, err)
		require.NoError(t, os.RemoveAll(path.Join(wd, config.LogPath, "..")))
	}()

	consoleWriter := config.NewConsoleWriter()
	require.NotNil(t, consoleWriter)
}
