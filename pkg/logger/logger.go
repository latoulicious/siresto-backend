package logger

import (
	"io"
	"os"

	"github.com/latoulicious/siresto-backend/pkg/core/logging"
	"github.com/sirupsen/logrus"
)

// LogrusLogger is the custom logger using logrus and persisting logs
type LogrusLogger struct {
	log       *logrus.Logger
	persister logging.LogPersister
	appName   string
	appEnv    string
	// Add a silent flag to control console output
	silent bool
}

// Constructor for LogrusLogger
func NewLogger(persister logging.LogPersister) logging.Logger {
	l := logrus.New()
	l.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})

	// Use environment variable to control console output
	silent := os.Getenv("LOG_SILENT") == "true"

	// If silent mode is enabled, discard console output
	if silent {
		l.SetOutput(io.Discard)
	} else {
		l.SetOutput(os.Stdout)
	}

	l.SetLevel(logrus.InfoLevel)

	return &LogrusLogger{
		log:       l,
		persister: persister,
		appName:   os.Getenv("APP_NAME"),
		appEnv:    os.Getenv("APP_ENV"),
		silent:    silent,
	}
}

// LogInfo logs an info message and persists it
func (l *LogrusLogger) LogInfo(msg string, fields map[string]interface{}) {
	// Only log to console if not silent
	if !l.silent {
		l.log.WithFields(fields).Info(msg)
	}
	// Always persist to DB
	l.persist("info", msg, fields)
}

// LogError logs an error message and persists it
func (l *LogrusLogger) LogError(msg string, fields map[string]interface{}) {
	// Errors are always logged to console regardless of silent mode
	l.log.WithFields(fields).Error(msg)
	l.persist("error", msg, fields)
}

// LogDebug logs a debug message and persists it
func (l *LogrusLogger) LogDebug(msg string, fields map[string]interface{}) {
	// Only log to console if not silent
	if !l.silent {
		l.log.WithFields(fields).Debug(msg)
	}
	l.persist("debug", msg, fields)
}

// persist stores the log in the database using the persister
func (l *LogrusLogger) persist(level, msg string, fields map[string]interface{}) {
	if l.persister == nil {
		return
	}

	// Pass the log to the persister to store in the database
	err := l.persister.PersistLog(level, msg, fields)
	if err != nil && !l.silent {
		// Only log persistence errors to console if not in silent mode
		l.log.WithFields(fields).Error("Failed to persist log to database")
	}
}
