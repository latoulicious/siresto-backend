package repository

import (
	"github.com/google/uuid"
	"github.com/latoulicious/siresto-backend/internal/domain"
	"gorm.io/gorm"
)

type OrderDetailRepository struct {
	DB *gorm.DB
}

func (r *OrderDetailRepository) Create(orderDetail *domain.OrderDetail) error {
	return r.DB.Create(orderDetail).Error
}

func (r *OrderDetailRepository) GetByOrderID(orderID uuid.UUID) ([]domain.OrderDetail, error) {
	var details []domain.OrderDetail
	err := r.DB.Where("order_id = ?", orderID).Find(&details).Error
	return details, err
}
