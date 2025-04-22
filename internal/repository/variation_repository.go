package repository

import (
	"github.com/google/uuid"
	"github.com/latoulicious/siresto-backend/internal/domain"
	"gorm.io/gorm"
)

type VariationRepository struct {
	DB *gorm.DB
}

//! Global Variation Repository

// ListVariations fetches all variations for a product
func (r *VariationRepository) ListAllVariations() ([]domain.Variation, error) {
	var variations []domain.Variation
	err := r.DB.Find(&variations).Error
	if err != nil {
		return nil, err
	}
	return variations, nil
}

// GetVariationByID fetches a variation by its ID
func (r *VariationRepository) GetVariationByID(id uuid.UUID) (*domain.Variation, error) {
	var variation domain.Variation
	err := r.DB.Where("id = ?", id).First(&variation).Error
	if err != nil {
		return nil, err
	}
	return &variation, nil
}

// CreateVariation creates a new variation for a product
func (r *VariationRepository) CreateVariation(variation *domain.Variation) error {
	return r.DB.Create(variation).Error
}

// UpdateVariation updates an existing variation
func (r *VariationRepository) UpdateVariation(variation *domain.Variation) error {
	return r.DB.Save(variation).Error
}

// DeleteVariation deletes a variation by its ID
func (r *VariationRepository) DeleteVariation(id uuid.UUID) error {
	return r.DB.Delete(&domain.Variation{}, id).Error
}

// TODO Implement Function for Product Tied Variations

// Helper Function
func (r *VariationRepository) ProductExists(productID uuid.UUID) (bool, error) {
	var exists bool
	err := r.DB.Raw("SELECT EXISTS(SELECT 1 FROM products WHERE id = ?)", productID).Scan(&exists).Error
	return exists, err
}

// GetVariationsByProductID fetches all variations for a specific product
func (r *VariationRepository) GetVariationsByProductID(productID uuid.UUID) ([]domain.Variation, error) {
	var variations []domain.Variation
	err := r.DB.Where("product_id = ?", productID).Find(&variations).Error
	if err != nil {
		return nil, err
	}
	return variations, nil
}

// CreateProductVariation creates a new variation tied to a specific product
func (r *VariationRepository) CreateProductVariation(variation *domain.Variation) error {
	return r.DB.Create(variation).Error
}
