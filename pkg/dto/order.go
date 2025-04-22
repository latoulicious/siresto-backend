package dto

import "time"

// OrderResponseDTO - Clean representation for client consumption
type OrderResponseDTO struct {
	ID            string         `json:"id"`
	CustomerName  string         `json:"customerName"`
	CustomerPhone string         `json:"customerPhone,omitempty"`
	TableNumber   int            `json:"tableNumber"`
	Status        string         `json:"status"`
	DishStatus    string         `json:"dishStatus,omitempty"`
	TotalAmount   float64        `json:"totalAmount"`
	Notes         string         `json:"notes,omitempty"`
	CreatedAt     time.Time      `json:"createdAt"`
	PaidAt        *time.Time     `json:"paidAt,omitempty"`
	Items         []OrderItemDTO `json:"items"`
}

// OrderItemDTO - Simplified order item representation
type OrderItemDTO struct {
	ID          string  `json:"id"`
	ProductID   string  `json:"productId"`
	ProductName string  `json:"productName"`
	Variation   string  `json:"variation,omitempty"`
	Quantity    int     `json:"quantity"`
	UnitPrice   float64 `json:"unitPrice"`
	TotalPrice  float64 `json:"totalPrice"`
	Note        string  `json:"note,omitempty"`
}
