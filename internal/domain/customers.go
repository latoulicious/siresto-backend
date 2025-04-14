package domain

import (
	"time"

	"github.com/google/uuid"
)

type Customer struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name        string    `gorm:"type:text"`
	Email       string    `gorm:"type:text"`
	Phone       string    `gorm:"type:text"`
	CreatedAt   time.Time `gorm:"default:now()"`
	LastOrderAt *time.Time
}
