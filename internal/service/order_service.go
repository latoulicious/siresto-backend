package service

import (
	"errors"

	"github.com/latoulicious/siresto-backend/internal/domain"
	"github.com/latoulicious/siresto-backend/internal/repository"
)

type OrderService struct {
	Repo *repository.OrderRepository
}

func (service *OrderService) CreateOrder(order *domain.Order, orderDetails []domain.OrderDetail) (*domain.Order, error) {
	// Perform business validations
	if order.UserID == nil || order.CustomerName == "" || order.CustomerPhone == "" {
		return nil, errors.New("missing required order fields")
	}

	// Calculate total amount from order details if not provided
	if order.TotalAmount == 0 {
		var total float64
		for _, detail := range orderDetails {
			total += detail.TotalPrice
		}
		order.TotalAmount = total
	}

	// Create the order and details in the repository
	createdOrder, err := service.Repo.CreateOrder(order, orderDetails)
	if err != nil {
		return nil, err
	}

	// Add logic to retrieve complete order with details if needed
	// This would require enhancing OrderRepository with a getWithDetails method

	return createdOrder, nil
}
