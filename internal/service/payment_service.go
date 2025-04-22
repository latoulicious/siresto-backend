package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/latoulicious/siresto-backend/internal/domain"
	"github.com/latoulicious/siresto-backend/internal/repository"
)

type PaymentService struct {
	Repo *repository.PaymentRepository
}

func (s *PaymentService) ListAllPayments() ([]domain.Payment, error) {
	return s.Repo.ListAllPayments()
}

func (s *PaymentService) CreatePayment(payment *domain.Payment) error {
	return s.Repo.Create(payment)
}

// internal/service/payment_service.go
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

	// 2. Verify payment amount matches order total
	if payment.Amount != order.TotalAmount {
		tx.Rollback()
		return nil, fmt.Errorf("payment amount (%f) does not match order total (%f)",
			payment.Amount, order.TotalAmount)
	}

	// 3. Set payment fields
	payment.OrderID = orderID
	payment.Status = domain.PaymentStatusSuccess

	// Set paid time if not explicitly set
	if payment.PaidAt.IsZero() {
		payment.PaidAt = time.Now()
	}

	// 4. Create payment record
	if err := tx.Create(payment).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create payment: %w", err)
	}

	// 5. Update order status to paid
	now := time.Now()
	if err := tx.Model(&order).Updates(map[string]interface{}{
		"status":  domain.OrderStatusPaid,
		"paid_at": now,
	}).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update order status: %w", err)
	}

	// 6. Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("transaction failed: %w", err)
	}

	return payment, nil
}
