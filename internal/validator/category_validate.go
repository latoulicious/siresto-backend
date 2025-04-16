package validator

import (
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/latoulicious/siresto-backend/internal/domain"
	"gorm.io/gorm"
)

var (
	ErrCategoryNameEmpty    = errors.New("category name cannot be empty")
	ErrInvalidCategoryPos   = errors.New("category position cannot be negative")
	ErrCategoryAlreadyExist = errors.New("category with the same name already exists")
)

// CategoryExistsByID checks if a category exists in the DB by ID
func CategoryExistsByID(db *gorm.DB, id uuid.UUID) (bool, error) {
	var count int64
	err := db.Model(&domain.Category{}).Where("id = ?", id).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("error checking category existence by ID: %w", err)
	}
	return count > 0, nil
}

// CategoryExistsByName checks if a category exists in the DB by name (case insensitive)
func CategoryExistsByName(db *gorm.DB, name string) (bool, error) {
	var count int64
	err := db.Model(&domain.Category{}).Where("LOWER(name) = LOWER(?)", name).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("error checking category existence by name: %w", err)
	}
	return count > 0, nil
}

// ValidateCategory validates category fields before insert
func ValidateCategory(db *gorm.DB, c *domain.Category) error {
	if strings.TrimSpace(c.Name) == "" {
		return ErrCategoryNameEmpty
	}
	if c.Position < 0 {
		return ErrInvalidCategoryPos
	}
	exists, err := CategoryExistsByName(db, c.Name)
	if err != nil {
		return err
	}
	if exists {
		return ErrCategoryAlreadyExist
	}
	return nil
}

// ValidateCategoryForUpdate allows name/position update with uniqueness check
func ValidateCategoryForUpdate(db *gorm.DB, c *domain.Category) error {
	if strings.TrimSpace(c.Name) == "" {
		return ErrCategoryNameEmpty
	}
	if c.Position < 0 {
		return ErrInvalidCategoryPos
	}

	// Ensure name is unique (excluding current category)
	var count int64
	err := db.Model(&domain.Category{}).
		Where("LOWER(name) = LOWER(?) AND id != ?", c.Name, c.ID).
		Count(&count).Error
	if err != nil {
		return fmt.Errorf("error validating category name uniqueness: %w", err)
	}
	if count > 0 {
		return ErrCategoryAlreadyExist
	}
	return nil
}

// ValidateCategoryDeletable ensures the category can be safely deleted
func ValidateCategoryDeletable(db *gorm.DB, categoryID uuid.UUID) error {
	var count int64
	if err := db.Model(&domain.Product{}).
		Where("category_id = ?", categoryID).
		Count(&count).Error; err != nil {
		return fmt.Errorf("error checking related products: %w", err)
	}

	if count > 0 {
		return fmt.Errorf("cannot delete category: %d associated products found", count)
	}

	return nil
}
