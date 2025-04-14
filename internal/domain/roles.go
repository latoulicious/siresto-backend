package domain

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	ID          uuid.UUID `gorm:"type:uuid;primary_key"`
	Name        string    `gorm:"unique;not null"`
	Description string    `gorm:"type:text"`
}
