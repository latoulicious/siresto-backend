package repository

import (
	"github.com/latoulicious/siresto-backend/internal/domain"
	"gorm.io/gorm"
)

type PaymentRepository struct {
	DB *gorm.DB
}

func (r *PaymentRepository) ListAllPayments() ([]domain.Payment, error) {
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
