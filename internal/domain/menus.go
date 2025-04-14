package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Menu struct {
	gorm.Model
	ID            uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name          string     `gorm:"type:text;not null"`
	IsActive      bool       `gorm:"default:true"`
	StoreID       uuid.UUID  `gorm:"type:uuid"`
	Position      int        `gorm:"default:0"`
	AvailableFrom *time.Time `gorm:"type:time"`
	AvailableTo   *time.Time `gorm:"type:time"`
	Categories    []Category `gorm:"foreignKey:MenuID"`
}
