package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/latoulicious/siresto-backend/internal/domain"
	"github.com/latoulicious/siresto-backend/internal/repository"
)

var (
	ErrOrderAlreadyPaid      = errors.New("order is already paid")
	ErrOrderCancelled        = errors.New("cannot process payment for cancelled order")
	ErrPaymentAmountMismatch = errors.New("payment amount does not match order total")
)

type PaymentService struct {
	Repo *repository.PaymentRepository
}

func (s *PaymentService) ListAllOrderPayments() ([]domain.Payment, error) {
	return s.Repo.ListAllOrderPayments()
}

func (s *PaymentService) CreatePayment(payment *domain.Payment) error {
	return s.Repo.Create(payment)
}

func (s *PaymentService) ProcessOrderPayment(orderID uuid.UUID, payment *domain.Payment) (*domain.Payment, error) {
	// Validate payment data
	if payment.Method == "" || payment.Amount <= 0 {
		return nil, errors.New("invalid payment data")
	}

	// Begin transaction
	tx := s.Repo.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. Fetch order with details to verify payment amount
	var order domain.Order
	if err := tx.Preload("OrderDetails").First(&order, "id = ?", orderID).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("order not found: %w", err)
	}

	// 2. Validate current order state
	if order.Status == domain.OrderStatusPaid {
		tx.Rollback()
		return nil, errors.New("order is already paid")
	}

	if order.Status == domain.OrderStatusCancelled {
		tx.Rollback()
		return nil, errors.New("cannot process payment for cancelled order")
	}

	// 3. Verify payment amount matches order total
	if payment.Amount != order.TotalAmount {
		tx.Rollback()
		return nil, fmt.Errorf("payment amount (%f) does not match order total (%f)",
			payment.Amount, order.TotalAmount)
	}

	// 4. Set payment fields
	payment.OrderID = orderID
	payment.Status = domain.PaymentStatusSuccess

	// Set paid time if not explicitly set
	if payment.PaidAt.IsZero() {
		payment.PaidAt = time.Now()
	}

	// 5. Create payment record
	if err := tx.Create(payment).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create payment: %w", err)
	}

	// 6. Update order status to paid AND dish status to Diproses
	now := time.Now()
	if err := tx.Model(&order).Updates(map[string]interface{}{
		"status":      domain.OrderStatusPaid,
		"dish_status": domain.FoodStatusInProcess, // Auto-transition to Diproses on payment
		"paid_at":     now,
	}).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update order status: %w", err)
	}

	// 7. Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("transaction failed: %w", err)
	}

	return payment, nil
}
