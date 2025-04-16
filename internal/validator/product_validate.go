package validator

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/latoulicious/siresto-backend/internal/domain"
	"gorm.io/gorm"
)

var (
	ErrMissingCategoryID = errors.New("CategoryID is required")
	ErrInvalidBasePrice  = errors.New("BasePrice must be greater than 0")
	ErrEmptyName         = errors.New("Product name cannot be empty")
	ErrCategoryNotFound  = errors.New("Category does not exist")
)

// CategoryExists checks if a category exists in the database by its ID.
func CategoryExists(db *gorm.DB, categoryID uuid.UUID) (bool, error) {
	var count int64
	// Check if the category exists in the database
	if err := db.Model(&domain.Category{}).Where("id = ?", categoryID).Count(&count).Error; err != nil {
		return false, fmt.Errorf("error checking category existence: %w", err)
	}

	// Return true if category exists, otherwise false
	return count > 0, nil
}

// ValidateProduct performs all validation on a product for creation or update
func ValidateProduct(db *gorm.DB, p *domain.Product) error {
	if p.CategoryID == nil || *p.CategoryID == uuid.Nil {
		return ErrMissingCategoryID
	}
	if p.Name == "" {
		return ErrEmptyName
	}
	if p.BasePrice <= 0 {
		return ErrInvalidBasePrice
	}
	exists, err := CategoryExists(db, *p.CategoryID)
	if err != nil {
		return fmt.Errorf("error checking category: %w", err)
	}
	if !exists {
		return ErrCategoryNotFound
	}
	return nil
}

// ValidateProductForUpdate performs validation on a product for update
func ValidateProductForUpdate(db *gorm.DB, p *domain.Product) error {
	if p.Name == "" {
		return ErrEmptyName
	}
	if p.BasePrice <= 0 {
		return ErrInvalidBasePrice
	}
	// only check category if it's explicitly allowed to change
	if p.CategoryID != nil && *p.CategoryID != uuid.Nil {
		exists, err := CategoryExists(db, *p.CategoryID)
		if err != nil {
			return fmt.Errorf("error checking category: %w", err)
		}
		if !exists {
			return ErrCategoryNotFound
		}
	}
	return nil
}

// ValidateProductDeletable ensures the product can be safely deleted
func ValidateProductDeletable(db *gorm.DB, productID uuid.UUID) error {
	var count int64
	if err := db.Model(&domain.Variation{}).
		Where("product_id = ?", productID).
		Count(&count).Error; err != nil {
		return fmt.Errorf("error checking related variations: %w", err)
	}

	if count > 0 {
		return fmt.Errorf("cannot delete product: %d associated variations found", count)
	}

	return nil
}
