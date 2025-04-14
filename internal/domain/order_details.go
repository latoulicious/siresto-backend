package domain

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderDetail struct {
	gorm.Model
	ID            uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	OrderID       uuid.UUID  `gorm:"type:uuid;not null"`
	Order         *Order     `gorm:"foreignKey:OrderID"`
	ProductName   string     `gorm:"type:text;not null"`
	VariationName string     `gorm:"type:text"`
	UnitPrice     float64    `gorm:"type:numeric(10,2);not null"`
	Quantity      int        `gorm:"not null"`
	TotalPrice    float64    `gorm:"type:numeric(10,2);not null"`
	ProductID     *uuid.UUID `gorm:"type:uuid"`
	Product       *Product   `gorm:"foreignKey:ProductID"`
	VariationID   *uuid.UUID `gorm:"type:uuid"`
	Variation     *Variation `gorm:"foreignKey:VariationID"`
}
