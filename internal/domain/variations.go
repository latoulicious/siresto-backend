package domain

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Variation struct {
	gorm.Model
	ID            uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	ProductID     uuid.UUID `gorm:"type:uuid;not null"`
	Product       *Product  `gorm:"foreignKey:ProductID"`
	Name          string    `gorm:"type:text;not null"`
	PriceModifier *float64  `gorm:"type:numeric(10,2)"`
	PriceAbsolute *float64  `gorm:"type:numeric(10,2)"`
	IsDefault     bool      `gorm:"default:false"`
	IsAvailable   bool      `gorm:"default:true"`
}
