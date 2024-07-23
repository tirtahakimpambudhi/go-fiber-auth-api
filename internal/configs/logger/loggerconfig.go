package loggerconfig

import (
	"github.com/phuslu/log"
	"github.com/tirtahakimpambudhi/restful_api/internal/configs"
	pathhelper "github.com/tirtahakimpambudhi/restful_api/pkg/helper/path"
	"os"
	"path/filepath"
)

// LoggerConfig holds configuration settings for logging.
type LoggerConfig struct {
	LogPath        string `env:"LOG_PATH,required"`
	MaxSize        int    `env:"LOG_MAX_SIZE,required"`
	MaxBackup      int    `env:"LOG_MAX_BACKUP,required"`
	MaxSizeRotate  int    `env:"LOG_MAX_SIZE_ROTATE,required"`
	TimeFormat     string `env:"LOG_TIME_FORMAT,required"`
	ColorOutput    bool   `env:"LOG_COLOR_OUTPUT" envDefault:"false"`
	QuoteString    bool   `env:"LOG_QUOTE_STR" envDefault:"false"`
	EndWithMessage bool   `env:"LOG_END_WITH_MESSAGE" envDefault:"false"`
}

// NewFileWriterWithRotate creates a file writer with rotation settings.
func (loggerConfig LoggerConfig) NewFileWriterWithRotate(filename string) *log.FileWriter {
	return &log.FileWriter{
		TimeFormat: loggerConfig.TimeFormat,
		Filename:   pathhelper.AddWorkdirToSomePath(loggerConfig.LogPath, filename),
		MaxSize:    int64(loggerConfig.MaxSize) * 1024 * 1024, // Convert size to bytes
		MaxBackups: loggerConfig.MaxBackup,
		Cleaner: func(filename string, maxBackups int, matches []os.FileInfo) {
			var dir = filepath.Dir(filename)
			var total int64
			// Iterate over log files in reverse order to manage rotation
			for i := len(matches) - 1; i >= 0; i-- {
				total += matches[i].Size()
				// Remove old files if total size exceeds the limit
				if total > int64(loggerConfig.MaxSizeRotate)*1024*1024 {
					os.Remove(filepath.Join(dir, matches[i].Name()))
				}
			}
		},
	}
}

// NewFileWriter creates a file writer without rotation settings.
func (loggerConfig LoggerConfig) NewFileWriter(filename string) *log.FileWriter {
	return &log.FileWriter{
		TimeFormat: loggerConfig.TimeFormat,
		Filename:   pathhelper.AddWorkdirToSomePath(loggerConfig.LogPath, filename),
		MaxSize:    int64(loggerConfig.MaxSize) * 1024 * 1024, // Convert size to bytes
		MaxBackups: loggerConfig.MaxBackup,
	}
}

// NewConsoleWriter creates a console writer with formatting options.
func (loggerConfig LoggerConfig) NewConsoleWriter() *log.ConsoleWriter {
	return &log.ConsoleWriter{
		ColorOutput:    loggerConfig.ColorOutput,
		QuoteString:    loggerConfig.QuoteString,
		EndWithMessage: loggerConfig.EndWithMessage,
	}
}

// NewLoggerConfig loads logger configuration from environment variables.
func NewLoggerConfig() (*LoggerConfig, error) {
	var loggerConfig LoggerConfig
	// Load configuration into LoggerConfig struct
	if err := configs.GetConfig().Load(&loggerConfig); err != nil {
		return nil, err
	}
	// Create necessary directories from configuration
	if err := pathhelper.MakedirFromFieldStruct(loggerConfig); err != nil {
		return nil, err
	}
	// Return loaded configuration
	return &loggerConfig, nil
}
