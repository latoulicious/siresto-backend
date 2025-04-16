package validator

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/latoulicious/siresto-backend/internal/domain"
	"gorm.io/gorm"
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
