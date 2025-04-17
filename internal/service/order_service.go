package service

import (
	"github.com/latoulicious/siresto-backend/internal/domain"
	"github.com/latoulicious/siresto-backend/internal/repository"
)

type OrderService struct {
	Repo *repository.OrderRepository
}

func (service *OrderService) CreateOrder(order *domain.Order, orderDetails []domain.OrderDetail) (*domain.Order, error) {
	// Validate or process any logic here, if necessary
	// For example, calculate total amounts or check customer/user

	// Call the repository to create the order and order details
	createdOrder, err := service.Repo.CreateOrder(order, orderDetails)
	if err != nil {
		return nil, err
	}

	// Additional logic can be added here

	return createdOrder, nil
}
