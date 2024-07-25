package loggerconfig

import (
	"github.com/phuslu/log"
	"os"
)

// Logger manages logging configuration and loggers for application and access logs.
type Logger struct {
	Config *LoggerConfig
	Access *log.Logger
	App    *log.Logger
}

// NewLogger creates a new Logger instance with configuration and loggers.
// Returns the Logger instance or an error if configuration fails.
func NewLogger() (*Logger, error) {
	// Create logger configuration
	config, err := NewLoggerConfig()
	if err != nil {
		// Return nil logger and error if configuration fails
		return nil, err
	}

	// Initialize application logger with configuration settings
	appLog := &log.Logger{
		TimeFormat: config.TimeFormat,
		Level:      log.DebugLevel,
		Caller:     1,
	}

	// Set writer for application logger based on terminal status
	if log.IsTerminal(os.Stderr.Fd()) {
		appLog.Writer = config.NewConsoleWriter()
	} else {
		appLog.Writer = config.NewFileWriterWithRotate("app.log")
	}

	// Return new Logger instance with access and application loggers
	return &Logger{
		Config: config,
		Access: &log.Logger{
			Level:  log.InfoLevel,
			Caller: 1,
			Writer: config.NewFileWriterWithRotate("access.log"),
		},
		App: appLog,
	}, err
}
