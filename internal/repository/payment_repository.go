package repository

import (
	"github.com/google/uuid"
	"github.com/latoulicious/siresto-backend/internal/domain"
	"gorm.io/gorm"
)

type PaymentRepository struct {
	DB *gorm.DB
}

func (r *PaymentRepository) ListAllOrderPayments() ([]domain.Payment, error) {
	var payments []domain.Payment
	err := r.DB.Find(&payments).Error
	if err != nil {
		return nil, err
	}
	return payments, nil
}

func (r *PaymentRepository) Create(payment *domain.Payment) error {
	return r.DB.Create(payment).Error
}

// internal/repository/payment_repository.go - Add these methods
func (r *PaymentRepository) GetByOrderID(orderID uuid.UUID) ([]domain.Payment, error) {
	var payments []domain.Payment
	err := r.DB.Where("order_id = ?", orderID).Find(&payments).Error
	return payments, err
}

func (r *PaymentRepository) GetPaymentsByStatus(status domain.PaymentStatus) ([]domain.Payment, error) {
	var payments []domain.Payment
	err := r.DB.Where("status = ?", status).Find(&payments).Error
	return payments, err
}
