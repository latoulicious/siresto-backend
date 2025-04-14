package domain

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	ID       uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	MenuID   uuid.UUID `gorm:"type:uuid;not null"`
	Menu     *Menu     `gorm:"foreignKey:MenuID"`
	Name     string    `gorm:"type:text;not null"`
	IsActive bool      `gorm:"default:true"`
	Position int       `gorm:"default:0"`
	Products []Product `gorm:"foreignKey:CategoryID"`
}
