package domain

import (
	"time"

	"github.com/google/uuid"
)

type PaymentStatus string
type PaymentType string

const (
	PaymentStatusPending  PaymentStatus = "PENDING"
	PaymentStatusSuccess  PaymentStatus = "SUCCESS"
	PaymentStatusFailed   PaymentStatus = "FAILED"
	PaymentStatusRefunded PaymentStatus = "REFUNDED"
)

const (
	PaymentTypeTunai  PaymentType = "Tunai"
	PaymentTypeQris   PaymentType = "Qris"
	PaymentTypeDebit  PaymentType = "Debit"
	PaymentTypeKredit PaymentType = "Kredit"
)

type Payment struct {
	ID             uuid.UUID     `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	OrderID        uuid.UUID     `gorm:"type:uuid;not null" json:"order_id"`
	Order          *Order        `gorm:"foreignKey:OrderID"`
	Method         PaymentType   `gorm:"type:text;not null"`
	Amount         float64       `gorm:"type:numeric(10,2);not null"`
	Status         PaymentStatus `gorm:"type:text;not null"`
	TransactionRef string        `gorm:"type:text"`
	PaidAt         time.Time     `gorm:"default:now()"`
}
