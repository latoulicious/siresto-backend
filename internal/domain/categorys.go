package domain

import (
	"github.com/google/uuid"
)

type Category struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name     string    `gorm:"type:text;not null"`
	IsActive bool      `gorm:"default:true"`
	Position int       `gorm:"default:0"`
	Products []Product `gorm:"foreignKey:CategoryID"`
}
