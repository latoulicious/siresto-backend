package domain

import (
	"github.com/google/uuid"
	"github.com/latoulicious/siresto-backend/pkg/db"
)

type Variation struct {
	ID            uuid.UUID           `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	ProductID     uuid.UUID           `gorm:"type:uuid;not null"`
	Product       *Product            `gorm:"foreignKey:ProductID"`
	IsDefault     bool                `gorm:"default:false"`
	IsAvailable   bool                `gorm:"default:true"`
	IsRequired    bool                `gorm:"default:false"`
	VariationType string              `gorm:"type:text;not null"`
	Options       db.VariationOptions `gorm:"type:jsonb;not null"`
}
