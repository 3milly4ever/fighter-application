package log

import (
	"io"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

// InitLogger initializes the logger and saves logs to a file specified in the .env file
func InitLogger() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		logrus.Warn("Error loading .env file, using default stderr")
	}

	// Get the log file path from the environment variable
	logFilePath := os.Getenv("LOG_FILE")
	if logFilePath == "" {
		// If LOG_FILE is not set in .env, fallback to default path
		logFilePath = filepath.Join("/Users/emilshalamberidze/FighterApplication/logs", "app.log")
	}

	// Ensure the logs directory exists
	logDir := filepath.Dir(logFilePath)
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		err := os.MkdirAll(logDir, 0755) // Use MkdirAll to create any intermediate directories
		if err != nil {
			logrus.Warn("Failed to create logs directory, using default stderr")
		}
	}

	// Set the output to the log file
	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logrus.Warn("Failed to log to file, using default stderr")
		logrus.SetOutput(os.Stdout) // Fallback to stdout if file creation fails
	} else {
		// Log to both file and terminal
		multiWriter := io.MultiWriter(file, os.Stdout)
		logrus.SetOutput(multiWriter)
	}

	// Set the log level (default is Info)
	logrus.SetLevel(logrus.InfoLevel)

	// Set log format (optional)
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
}
