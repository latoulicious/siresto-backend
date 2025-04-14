package domain

import (
	"time"

	"github.com/google/uuid"
)

type Invoice struct {
	ID               uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	OrderID          *uuid.UUID `gorm:"type:uuid"`
	Order            *Order     `gorm:"foreignKey:OrderID"`
	CustomerSnapshot string     `gorm:"type:jsonb"`
	ItemsSnapshot    string     `gorm:"type:jsonb"`
	Total            float64    `gorm:"type:numeric(10,2)"`
	InvoiceNumber    string     `gorm:"type:text;unique"`
	IssuedAt         time.Time  `gorm:"default:now()"`
	PdfURL           string     `gorm:"type:text"`
}
