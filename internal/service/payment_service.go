package service

import (
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
