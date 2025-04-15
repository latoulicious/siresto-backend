package logging

type Logger interface {
	LogDebug(msg string, fields map[string]interface{})
	LogInfo(msg string, fields map[string]interface{})
	LogError(msg string, fields map[string]interface{})
}

type LogPersister interface {
	PersistLog(level, msg string, fields map[string]interface{}) error
}
