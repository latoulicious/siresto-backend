package service

import (
	"errors"

	"github.com/latoulicious/siresto-backend/internal/domain"
	"github.com/latoulicious/siresto-backend/internal/repository"
)

type OrderService struct {
	Repo *repository.OrderRepository
}

func (s *OrderService) ListAllOrders() ([]domain.Order, error) {
	// Call the repository method to fetch all orders
	return s.Repo.ListAllOrders()
}

func (s *OrderService) CreateOrder(order *domain.Order, details []domain.OrderDetail) (*domain.Order, error) {
	if order.CustomerName == "" || order.CustomerPhone == "" {
		return nil, errors.New("missing required order fields")
	}

	if order.TotalAmount == 0 {
		var total float64
		for _, detail := range details {
			total += detail.TotalPrice
		}
		order.TotalAmount = total
	}

	createdOrder, err := s.Repo.CreateOrder(order, details)
	if err != nil {
		return nil, err
	}

	// Fetch with all associations
	return s.Repo.GetOrderWithAssociations(createdOrder.ID)
}
