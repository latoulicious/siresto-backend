package repository

import (
	"github.com/google/uuid"
	"github.com/latoulicious/siresto-backend/internal/domain"
	"gorm.io/gorm"
)

type CategoryRepository struct {
	DB *gorm.DB
}

// FindAll fetches all categories from the database
func (r *CategoryRepository) ListAllCategories() ([]domain.Category, error) {
	var categories []domain.Category
	err := r.DB.Find(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}

// FindByID fetches a category by its ID
func (r *CategoryRepository) GetCategoryByID(id uuid.UUID) (*domain.Category, error) {
	var category domain.Category
	err := r.DB.Where("id = ?", id).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// Create inserts a new category into the database
func (r *CategoryRepository) CreateCategory(category *domain.Category) error {
	return r.DB.Create(category).Error
}

// Update saves changes to an existing category
func (r *CategoryRepository) UpdateCategory(category *domain.Category) error {
	return r.DB.Save(category).Error
}

// Delete removes a category by its ID
func (r *CategoryRepository) DeleteCategory(id uuid.UUID) error {
	return r.DB.Delete(&domain.Category{}, "id = ?", id).Error
}
