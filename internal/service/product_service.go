package service

import (
	"github.com/google/uuid"
	"github.com/latoulicious/siresto-backend/internal/domain"
	"github.com/latoulicious/siresto-backend/internal/repository"
)

type ProductService struct {
	Repo *repository.ProductRepository
}

// ListAllProducts retrieves all products from the repository
func (s *ProductService) ListAllProducts() ([]domain.Product, error) {
	return s.Repo.ListAllProducts()
}

// GetProductByID retrieves a product by its ID from the repository
func (s *ProductService) GetProductByID(id uuid.UUID) (*domain.Product, error) {
	return s.Repo.GetProductByID(id)
}

// CreateProduct creates a new product in the repository
func (s *ProductService) CreateProduct(product *domain.Product) (*domain.Product, error) {
	product.ID = uuid.New()
	if err := s.Repo.CreateProduct(product); err != nil {
		return nil, err
	}
	return product, nil
}

// UpdateProduct updates an existing product in the repository
func (s *ProductService) UpdateProduct(id uuid.UUID, update *domain.Product) (*domain.Product, error) {
	existing, err := s.Repo.GetProductByID(id)
	if err != nil {
		return nil, err
	}
	existing.Name = update.Name
	existing.Description = update.Description
	existing.ImageURL = update.ImageURL
	existing.BasePrice = update.BasePrice
	existing.IsAvailable = update.IsAvailable
	existing.Position = update.Position

	if err := s.Repo.UpdateProduct(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

// DeleteProduct removes a product by its ID from the repository
func (s *ProductService) DeleteProduct(id uuid.UUID) error {
	return s.Repo.DeleteProduct(id)
}
