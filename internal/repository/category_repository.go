package repository

import (
	"github.com/google/uuid"
	"github.com/latoulicious/siresto-backend/internal/domain"
	"gorm.io/gorm"
)

type CategoryRepository struct {
	DB *gorm.DB
}

func (r *CategoryRepository) ListAllCategories() ([]domain.Category, error) {
	var categories []domain.Category
	// Preload products if needed
	err := r.DB.Preload("Products").Find(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (r CategoryRepository) GetCategoryByID(id uuid.UUID) (*domain.Category, error) {
	var category domain.Category
	err := r.DB.Preload("Products").Where("id = ?", id).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// Fetch all categories with products
func (r *CategoryRepository) ListAllCategoriesWithProducts() ([]domain.Category, error) {
	var categories []domain.Category
	// Preload products with their variations and other needed relationships
	err := r.DB.
		Model(&domain.Category{}).
		Preload("Products").
		Preload("Products.Variations"). // If you want variations with products
		Find(&categories).Error
	return categories, err
}

// Fetch single category by ID with products
func (r *CategoryRepository) GetCategoryByIDWithProducts(id uuid.UUID) (*domain.Category, error) {
	var category domain.Category
	err := r.DB.
		Model(&domain.Category{}).
		Preload("Product").
		First(&category, "id = ?", id).Error
	return &category, err
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
