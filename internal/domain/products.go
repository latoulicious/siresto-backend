package domain

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
	ID           uuid.UUID   `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	CategoryID   *uuid.UUID  `gorm:"type:uuid"`
	Category     *Category   `gorm:"foreignKey:CategoryID" json:"-"` // Hide full category object
	CategoryName string      `gorm:"-" json:"category_name"`         // Virtual field for JSON
	Name         string      `gorm:"type:text;not null"`
	Description  string      `gorm:"type:text"`
	ImageURL     string      `gorm:"type:text" json:"image_url"`
	BasePrice    float64     `gorm:"type:numeric(10,2)" json:"base_price"`
	IsAvailable  bool        `gorm:"default:true" json:"is_available"`
	Position     int         `gorm:"default:0"`
	Variations   []Variation `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE"`
}

// AfterFind is called by GORM after loading the entity from the database
func (p *Product) AfterFind(tx *gorm.DB) error {
	if p.Category != nil {
		p.CategoryName = p.Category.Name
	}
	return nil
}
