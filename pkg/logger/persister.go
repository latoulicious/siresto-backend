package logger

import (
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/latoulicious/siresto-backend/internal/domain"
	"github.com/latoulicious/siresto-backend/pkg/db"
	"gorm.io/gorm"
)

// LogServicePersister struct to persist logs
type LogServicePersister struct {
	DB *gorm.DB
}

// NewLogServicePersister creates a new LogServicePersister
func NewLogServicePersister(DB *gorm.DB) *LogServicePersister {
	return &LogServicePersister{DB: DB}
}

// PersistLog directly persists the log entry into the database
func (p *LogServicePersister) PersistLog(level, msg string, fields map[string]interface{}) error {
	log := domain.Log{
		ID:          generateUUID(), // Generate new UUID for the log
		Timestamp:   time.Now(),     // Current timestamp
		Level:       level,          // Log level (e.g., "info", "error")
		Source:      toString(fields["source"], "unknown"),
		Action:      toString(fields["action"], "unknown"),
		Entity:      toString(fields["entity"], "unknown"),
		Description: msg,
		Metadata:    mapToJSONB(fields),                // Use db.JSONB for metadata field
		RequestID:   toStringPtr(fields["request_id"]), // Convert to *string
		Environment: os.Getenv("APP_ENV"),
		Application: os.Getenv("APP_NAME"),
		Type:        "activity", // Default type for the log
	}

	// Insert the log directly into the database
	if err := p.DB.Create(&log).Error; err != nil {
		return err
	}

	return nil
}

// mapToJSONB converts fields to a db.JSONB object for storage
func mapToJSONB(fields map[string]interface{}) db.JSONB {
	if len(fields) > 0 {
		// Assuming db.JSONB is a custom type that handles JSON objects
		return db.JSONB(fields)
	}
	return db.JSONB{} // Return an empty JSONB object if fields are empty
}

// toString safely converts any value to a string, using a fallback value if nil or not a string
func toString(val interface{}, fallback string) string {
	if str, ok := val.(string); ok {
		return str
	}
	return fallback
}

// toStringPtr converts a string value to a *string (pointer to string)
func toStringPtr(val interface{}) *string {
	str := toString(val, "")
	if str == "" {
		return nil // Return nil if the string is empty
	}
	return &str // Return a pointer to the string
}

// generateUUID generates a new UUID and returns it as a uuid.UUID type
func generateUUID() uuid.UUID {
	return uuid.New() // Return uuid.UUID type
}
