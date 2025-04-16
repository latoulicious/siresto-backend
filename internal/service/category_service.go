package service

import (
	"github.com/google/uuid"
	"github.com/latoulicious/siresto-backend/internal/domain"
	"github.com/latoulicious/siresto-backend/internal/repository"
	"github.com/latoulicious/siresto-backend/internal/validator"
)

type CategoryService struct {
	Repo *repository.CategoryRepository
}

func (s *CategoryService) ListAllCategories(includeProducts bool) ([]domain.Category, error) {
	if includeProducts {
		return s.Repo.ListAllCategoriesWithProducts()
	}
	return s.Repo.ListAllCategories()
}

func (s *CategoryService) GetCategoryByID(id uuid.UUID, includeProducts bool) (*domain.Category, error) {
	if includeProducts {
		return s.Repo.GetCategoryByIDWithProducts(id)
	}
	return s.Repo.GetCategoryByID(id)
}

// Create creates a new category
func (s *CategoryService) CreateCategory(category *domain.Category) (*domain.Category, error) {
	// Validate category before creation
	if err := validator.ValidateCategory(s.Repo.DB, category); err != nil {
		return nil, err
	}

	category.ID = uuid.New()

	if err := s.Repo.CreateCategory(category); err != nil {
		return nil, err
	}

	// Ensure Products is an empty slice instead of nil
	category.Products = []domain.Product{}

	return category, nil
}

// Update updates an existing category
func (s *CategoryService) UpdateCategory(id uuid.UUID, update *domain.Category) (*domain.Category, error) {
	// Validate category before update
	if err := validator.ValidateCategory(s.Repo.DB, update); err != nil {
		return nil, err
	}

	existing, err := s.Repo.GetCategoryByID(id)
	if err != nil {
		return nil, err
	}
	existing.Name = update.Name
	existing.IsActive = update.IsActive
	existing.Position = update.Position

	if err := s.Repo.UpdateCategory(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

// Delete removes a category by ID
func (s *CategoryService) DeleteCategory(id uuid.UUID) error {
	// Validate category before deletion
	if err := validator.ValidateCategoryDeletable(s.Repo.DB, id); err != nil {
		return err
	}

	return s.Repo.DeleteCategory(id)
}
