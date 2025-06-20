package logging

import (
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

// InitLogger initializes and returns a configured logger instance
func InitLogger() *logrus.Logger {
	logger := logrus.New()

	// Create logs directory if it doesn't exist
	logsDir := "logs"
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		logger.WithError(err).Warn("Failed to create logs directory")
	}

	// Set up file logging
	logFile := filepath.Join(logsDir, "goreconx.log")
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logger.WithError(err).Warn("Failed to open log file, using stdout")
		logger.SetOutput(os.Stdout)
	} else {
		logger.SetOutput(file)
	}

	// Set log level
	logger.SetLevel(logrus.InfoLevel)

	// Set formatter
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   false,
	})

	logger.Info("Logger initialized successfully")
	return logger
}
