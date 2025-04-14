package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Log = logrus.New()

// InitLogger sets up the logger with desired log level and output
func InitLogger() {
	Log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	Log.SetOutput(os.Stdout)
	Log.SetLevel(logrus.InfoLevel) // Set log level (can be adjusted)
}

// LogInfo logs general info messages
func LogInfo(msg string, fields map[string]interface{}) {
	Log.WithFields(fields).Info(msg)
}

// LogError logs error messages
func LogError(msg string, fields map[string]interface{}) {
	Log.WithFields(fields).Error(msg)
}

// LogDebug logs debug messages
func LogDebug(msg string, fields map[string]interface{}) {
	Log.WithFields(fields).Debug(msg)
}
