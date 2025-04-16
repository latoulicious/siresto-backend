package service

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/latoulicious/siresto-backend/internal/domain"
	"github.com/latoulicious/siresto-backend/internal/repository"
	"github.com/latoulicious/siresto-backend/internal/validator"
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
func (s ProductService) CreateProduct(product *domain.Product) (*domain.Product, error) {
	// Set the product ID if it's not set
	if product.ID == uuid.Nil {
		product.ID = uuid.New()
	}

	// Perform validation
	if err := validator.ValidateProduct(s.Repo.DB, product); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Create the product in the repository
	if err := s.Repo.CreateProduct(product); err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	// Load the product with relations after creation
	created := &domain.Product{}
	if err := s.Repo.LoadProductWithRelations(product.ID, created); err != nil {
		return nil, fmt.Errorf("failed to reload product after creation: %w", err)
	}

	return created, nil
}

// UpdateProduct updates an existing product in the repository
func (s *ProductService) UpdateProduct(id uuid.UUID, update *domain.Product) (*domain.Product, error) {
	// Fetch the existing product
	existing, err := s.Repo.GetProductByID(id)
	if err != nil {
		return nil, err
	}

	// Perform validation
	if err := validator.ValidateProductForUpdate(s.Repo.DB, update); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Update the product fields
	existing.Name = update.Name
	existing.Description = update.Description
	existing.ImageURL = update.ImageURL
	existing.BasePrice = update.BasePrice
	existing.IsAvailable = update.IsAvailable
	existing.Position = update.Position

	// Save the updated product
	if err := s.Repo.UpdateProduct(existing); err != nil {
		return nil, err
	}

	return existing, nil
}

// DeleteProduct removes a product by its ID from the repository
func (s *ProductService) DeleteProduct(id uuid.UUID) error {
	// Perform deletability validation
	if err := validator.ValidateProductDeletable(s.Repo.DB, id); err != nil {
		return fmt.Errorf("product cannot be deleted: %w", err)
	}

	// Proceed with deletion
	return s.Repo.DeleteProduct(id)
}
