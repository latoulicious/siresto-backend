package dto

import "time"

// --- Request DTOs ---
type OrderResponseDTO struct {
	ID             string             `json:"id"`
	CustomerName   string             `json:"customerName"`
	CustomerPhone  string             `json:"customerPhone,omitempty"`
	TableNumber    int                `json:"tableNumber"`
	Status         string             `json:"status"`
	DishStatus     string             `json:"dishStatus,omitempty"`
	TotalAmount    float64            `json:"totalAmount"`
	Notes          string             `json:"notes,omitempty"`
	CreatedAt      time.Time          `json:"createdAt"`
	PaidAt         *time.Time         `json:"paidAt,omitempty"`
	CancelledAt    *time.Time         `json:"cancelledAt,omitempty"`
	PaymentMethods []PaymentMethodDTO `json:"paymentMethods"`
	Items          []OrderItemDTO     `json:"items"`
}

type PaymentDTO struct {
	ID             string    `json:"id"`
	Method         string    `json:"method"`
	Amount         float64   `json:"amount"`
	Status         string    `json:"status"`
	TransactionRef string    `json:"transaction_ref,omitempty"`
	PaidAt         time.Time `json:"paid_at"`
}

// --- Response DTOs ---
type OrderItemDTO struct {
	ID          string  `json:"id"`
	ProductID   string  `json:"productId"`
	ProductName string  `json:"productName"`
	Variation   string  `json:"variation,omitempty"`
	Quantity    int     `json:"quantity"`
	UnitPrice   float64 `json:"unitPrice"`
	TotalPrice  float64 `json:"totalPrice"`
	Note        string  `json:"note,omitempty"`
	ImageURL    string  `json:"imageUrl,omitempty"`
}

type PaymentMethodDTO struct {
	ID             string    `json:"id"`
	Method         string    `json:"method"`
	Amount         float64   `json:"amount"`
	Status         string    `json:"status"`
	TransactionRef string    `json:"transactionRef,omitempty"`
	PaidAt         time.Time `json:"paidAt"`
}
