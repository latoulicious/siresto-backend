package domain

import (
	"time"

	"github.com/google/uuid"
)

type OrderStatus string
type FoodStatus string

const (
	OrderStatusPending   OrderStatus = "Pending"
	OrderStatusPaid      OrderStatus = "Paid"
	OrderStatusCancelled OrderStatus = "Cancelled"
)

const (
	FoodStatusReceived  FoodStatus = "Received"
	FoodStatusInProcess FoodStatus = "In Process"
	FoodStatusCompleted FoodStatus = "Completed"
	FoodStatusCancelled FoodStatus = "Cancelled"
)

type Order struct {
	ID            uuid.UUID   `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	CustomerName  string      `gorm:"type:text;not null"`
	CustomerPhone string      `gorm:"type:text;not null"`
	TableNumber   int         `gorm:"type:int;not null"`
	Status        OrderStatus `gorm:"type:text;not null;default:'PENDING'"`
	DishStatus    FoodStatus  `gorm:"type:text;not null;default:'Diterima'"`
	TotalAmount   float64     `gorm:"type:numeric(10,2);default:0" json:"total_amount"`
	Notes         string      `gorm:"type:text"`
	CreatedAt     time.Time   `gorm:"default:now()"`
	PaidAt        *time.Time
	CancelledAt   *time.Time
	OrderDetails  []OrderDetail `gorm:"foreignKey:OrderID"`
	Payments      []Payment     `gorm:"foreignKey:OrderID"`
	Invoice       *Invoice      `gorm:"foreignKey:OrderID"`
}
