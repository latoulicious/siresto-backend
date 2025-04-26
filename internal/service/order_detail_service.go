package service

import (
	"github.com/google/uuid"
	"github.com/latoulicious/siresto-backend/internal/domain"
	"github.com/latoulicious/siresto-backend/internal/repository"
)

type OrderDetailService struct {
	Repo *repository.OrderDetailRepository
}

func (s *OrderDetailService) Create(details []domain.OrderDetail) error {
	for i := range details {
		if err := s.Repo.Create(&details[i]); err != nil {
			return err
		}
	}
	return nil
}

func (s *OrderDetailService) GetByOrderID(orderID uuid.UUID) ([]domain.OrderDetail, error) {
	var details []domain.OrderDetail
	if err := s.Repo.DB.Where("order_id = ?", orderID).Find(&details).Error; err != nil {
		return nil, err
	}
	return details, nil
}
