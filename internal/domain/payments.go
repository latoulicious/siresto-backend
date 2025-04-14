package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Payment struct {
	gorm.Model
	ID             uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	OrderID        uuid.UUID `gorm:"type:uuid;not null"`
	Order          *Order    `gorm:"foreignKey:OrderID"`
	Method         string    `gorm:"type:text;not null"`
	Amount         float64   `gorm:"type:numeric(10,2);not null"`
	Status         string    `gorm:"type:text;not null"`
	TransactionRef string    `gorm:"type:text"`
	PaidAt         time.Time `gorm:"default:now()"`
}
