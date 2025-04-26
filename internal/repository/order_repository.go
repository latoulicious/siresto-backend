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

func (repo *OrderRepository) ListAllOrders() ([]domain.Order, error) {
	var orders []domain.Order
	if err := repo.DB.
		Preload("OrderDetails").
		Preload("OrderDetails.Product").
		Preload("OrderDetails.Variation").
		Preload("Payments").
		Preload("Invoice").
		Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (repo *OrderRepository) GetOrderByID(orderID uuid.UUID) (*domain.Order, error) {
	var order domain.Order
	if err := repo.DB.First(&order, "id = ?", orderID).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (repo *OrderRepository) CreateOrder(order *domain.Order, orderDetails []domain.OrderDetail) (*domain.Order, error) {
	// Start a transaction
	tx := repo.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

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

	// Assign order details to the order for returning
	order.OrderDetails = orderDetails

	return order, nil
}

// GetOrderWithDetails retrieves an order with its related details
func (repo *OrderRepository) GetOrderWithDetails(orderID uuid.UUID) (*domain.Order, error) {
	var order domain.Order

	// Use preload to retrieve the order with all its relationships
	if err := repo.DB.
		Preload("OrderDetails").
		Preload("User").
		Preload("Payments").
		Preload("Invoice").
		First(&order, "id = ?", orderID).Error; err != nil {
		return nil, err
	}

	return &order, nil
}

func (r *OrderRepository) GetOrderWithAssociations(orderID uuid.UUID) (*domain.Order, error) {
	var order domain.Order
	if err := r.DB.
		Preload("OrderDetails").
		Preload("OrderDetails.Product").
		Preload("OrderDetails.Variation").
		Preload("Payments").
		Preload("Invoice").
		First(&order, "id = ?", orderID).Error; err != nil {
		return nil, err
	}
	return &order, nil
}
