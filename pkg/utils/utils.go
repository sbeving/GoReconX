package utils

import (
	"crypto/rand"
	"encoding/hex"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
)

// InitLogger initializes and returns a configured logger
func InitLogger() *logrus.Logger {
	logger := logrus.New()

	// Set log level
	logger.SetLevel(logrus.InfoLevel)

	// Create logs directory
	logDir := "./logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		logger.Warnf("Failed to create logs directory: %v", err)
	}

	// Set up log file
	logFile := filepath.Join(logDir, "gorconx.log")
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		logger.SetOutput(file)
	} else {
		logger.Warnf("Failed to open log file: %v", err)
	}

	// Set formatter
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})

	return logger
}

// GenerateRandomString generates a random string of specified length
func GenerateRandomString(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		// Fallback to timestamp-based generation
		return hex.EncodeToString([]byte(time.Now().String()))[:length]
	}
	return hex.EncodeToString(bytes)[:length]
}

// GetCurrentTimestamp returns current Unix timestamp
func GetCurrentTimestamp() int64 {
	return time.Now().Unix()
}

// ValidateTarget validates if a target string is valid
func ValidateTarget(target string) bool {
	if target == "" {
		return false
	}
	// Add more validation logic here
	return true
}

// EnsureDirectory ensures a directory exists, creating it if necessary
func EnsureDirectory(path string) error {
	return os.MkdirAll(path, 0755)
}

// FileExists checks if a file exists
func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

// Contains checks if a slice contains a string
func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// FormatDuration formats a duration for human-readable display
func FormatDuration(d time.Duration) string {
	if d < time.Minute {
		return d.Truncate(time.Second).String()
	}
	if d < time.Hour {
		return d.Truncate(time.Minute).String()
	}
	return d.Truncate(time.Hour).String()
}
