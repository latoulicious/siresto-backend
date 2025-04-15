package service

import (
	"github.com/google/uuid"
	"github.com/latoulicious/siresto-backend/internal/domain"
	"github.com/latoulicious/siresto-backend/internal/repository"
)

type VariationService struct {
	Repo *repository.VariationRepository
}

// ListVariations fetches all variations for a product
func (s *VariationService) ListAllVariations() ([]domain.Variation, error) {
	return s.Repo.ListAllVariations()
}

// GetVariationByID fetches a variation by its ID
func (s *VariationService) GetVariationByID(id uuid.UUID) (*domain.Variation, error) {
	return s.Repo.GetVariationByID(id)
}

// CreateVariation creates a new variation for a product
func (s *VariationService) CreateVariation(variation *domain.Variation) (*domain.Variation, error) {
	err := s.Repo.CreateVariation(variation)
	if err != nil {
		return nil, err
	}
	return variation, nil
}

// UpdateVariation updates an existing variation
func (s *VariationService) UpdateVariation(variation *domain.Variation) (*domain.Variation, error) {
	err := s.Repo.UpdateVariation(variation)
	if err != nil {
		return nil, err
	}
	return variation, nil
}

// DeleteVariation deletes a variation by its ID
func (s *VariationService) DeleteVariation(id uuid.UUID) error {
	return s.Repo.DeleteVariation(id)
}
