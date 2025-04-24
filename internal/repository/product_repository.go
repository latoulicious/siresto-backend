package repository

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/latoulicious/siresto-backend/internal/domain"
	"gorm.io/gorm"
)

type ProductRepository struct {
	DB *gorm.DB
}

// ListAllProducts fetches all products from the database
func (r *ProductRepository) ListAllProducts() ([]domain.Product, error) {
	var products []domain.Product
	err := r.DB.
		Preload("Category").
		Preload("Variations").
		Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}

// ListProductsPaginated fetches products with pagination
func (r *ProductRepository) ListProductsPaginated(offset, limit int) ([]domain.Product, int, error) {
	var products []domain.Product
	var total int64

	// Get total count
	if err := r.DB.Model(&domain.Product{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count products: %w", err)
	}

	// Get paginated data
	err := r.DB.
		Preload("Category").
		Preload("Variations").
		Offset(offset).
		Limit(limit).
		Find(&products).Error

	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch products: %w", err)
	}

	return products, int(total), nil
}

// GetProductByID fetches a product by its ID
func (r *ProductRepository) GetProductByID(id uuid.UUID) (*domain.Product, error) {
	var product domain.Product
	err := r.DB.Preload("Category").Preload("Variations").Where("id = ?", id).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepository) LoadProductWithRelations(id uuid.UUID, dest *domain.Product) error {
	return r.DB.
		Model(&domain.Product{}).
		Preload("Category").
		Preload("Variations").
		First(dest, "id = ?", id).Error
}

// CreateProduct inserts a new product into the database
func (r *ProductRepository) CreateProduct(product *domain.Product) error {
	return r.DB.Create(product).Error
}

// UpdateProduct saves changes to an existing product
func (r *ProductRepository) UpdateProduct(product *domain.Product) error {
	return r.DB.Save(product).Error
}

// DeleteProduct removes a product by its ID
func (r *ProductRepository) DeleteProduct(id uuid.UUID) error {
	return r.DB.Delete(&domain.Product{}, "id = ?", id).Error
}

// Helper Function

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
