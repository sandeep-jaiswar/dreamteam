package logger

import (
	"sync"

	"go.uber.org/zap"
)

var (
	instance *zap.Logger
	once     sync.Once
)

// NewProduction initializes a new production logger with default settings.
// It returns a *zap.Logger instance and an error if initialization fails.
func NewProduction() (*zap.Logger, error) {
	config := zap.NewProductionConfig()
	config.EncoderConfig.TimeKey = "timestamp" // Example of customizing the timestamp key
	logger, err := config.Build(zap.AddCallerSkip(1))
	if err != nil {
		return nil, err
	}
	return logger, nil
}

// NewLogger initializes a custom logger based on the provided configuration.
// This function allows for more flexibility compared to the predefined production and development loggers.
func NewLogger(config zap.Config, callerSkip int) (*zap.Logger, error) {
	logger, err := config.Build(zap.AddCallerSkip(callerSkip))
	if err != nil {
		return nil, err
	}
	return logger, nil
}

// GetLoggerInstance returns a singleton instance of the logger.
// It ensures that only one logger is created and reused throughout the application.
func GetLoggerInstance() (*zap.Logger, error) {
	var err error
	once.Do(func() {
		instance, err = NewProduction()
	})
	if err != nil {
		return nil, err
	}
	return instance, nil
}
