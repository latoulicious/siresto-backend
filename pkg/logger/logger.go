package logger

import (
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
}

// Constructor for LogrusLogger
func NewLogger(persister logging.LogPersister) logging.Logger {
	l := logrus.New()
	l.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	l.SetOutput(os.Stdout)
	l.SetLevel(logrus.InfoLevel)

	return &LogrusLogger{
		log:       l,
		persister: persister,
		appName:   os.Getenv("APP_NAME"),
		appEnv:    os.Getenv("APP_ENV"),
	}
}

// LogInfo logs an info message and persists it
func (l *LogrusLogger) LogInfo(msg string, fields map[string]interface{}) {
	l.log.WithFields(fields).Info(msg)
	l.persist("info", msg, fields)
}

// LogError logs an error message and persists it
func (l *LogrusLogger) LogError(msg string, fields map[string]interface{}) {
	l.log.WithFields(fields).Error(msg)
	l.persist("error", msg, fields)
}

// LogDebug logs a debug message and persists it
func (l *LogrusLogger) LogDebug(msg string, fields map[string]interface{}) {
	l.log.WithFields(fields).Debug(msg)
	l.persist("debug", msg, fields)
}

// persist stores the log in the database using the persister
func (l *LogrusLogger) persist(level, msg string, fields map[string]interface{}) {
	if l.persister == nil {
		return
	}

	// Pass the log to the persister to store in the database
	err := l.persister.PersistLog(level, msg, fields)
	if err != nil {
		l.log.WithFields(fields).Error("Failed to persist log to database")
	}
}
