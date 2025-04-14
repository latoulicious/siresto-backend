package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	ID           uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID       *uuid.UUID `gorm:"type:uuid"`
	User         *User      `gorm:"foreignKey:UserID"`
	CustomerID   *uuid.UUID `gorm:"type:uuid"`
	Customer     *Customer  `gorm:"foreignKey:CustomerID"`
	QRID         *uuid.UUID `gorm:"type:uuid"`
	QRCode       *QRCode    `gorm:"foreignKey:QRID"`
	Status       string     `gorm:"type:text;not null"`
	TotalAmount  float64    `gorm:"type:numeric(10,2);default:0"`
	Notes        string     `gorm:"type:text"`
	CreatedAt    time.Time  `gorm:"default:now()"`
	PaidAt       *time.Time
	CancelledAt  *time.Time
	OrderDetails []OrderDetail `gorm:"foreignKey:OrderID"`
	Payments     []Payment     `gorm:"foreignKey:OrderID"`
	Invoice      *Invoice      `gorm:"foreignKey:OrderID"`
}
