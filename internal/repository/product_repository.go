package repository

import (
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

// GetProductByID fetches a product by its ID
func (r *ProductRepository) GetProductByID(id uuid.UUID) (*domain.Product, error) {
	var product domain.Product
	err := r.DB.Preload("Category").Where("id = ?", id).First(&product).Error
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
