package domain

import (
	"time"

	"github.com/google/uuid"
)

type Role struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key"`
	Name        string    `gorm:"unique;not null"`
	Description string    `gorm:"type:text"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
