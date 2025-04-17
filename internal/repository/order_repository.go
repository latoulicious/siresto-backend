package repository

import (
	"errors"

	"github.com/google/uuid"
	"github.com/latoulicious/siresto-backend/internal/domain"
	"gorm.io/gorm"
)

type OrderRepository struct {
	DB *gorm.DB
}

func (repo *OrderRepository) CreateOrder(order *domain.Order, orderDetails []domain.OrderDetail) (*domain.Order, error) {
	// Start a transaction
	tx := repo.DB.Begin()

	// Create the Order
	if err := tx.Create(order).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Ensure the order ID is populated
	if order.ID == uuid.Nil {
		tx.Rollback()
		return nil, errors.New("failed to create order, missing ID")
	}

	// Create the OrderDetails
	for i := range orderDetails {
		orderDetails[i].OrderID = order.ID // Set the OrderID for each order detail
		if err := tx.Create(&orderDetails[i]).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return order, nil
}
