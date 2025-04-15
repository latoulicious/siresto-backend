package service

import (
	"github.com/google/uuid"
	"github.com/latoulicious/siresto-backend/internal/domain"
	"github.com/latoulicious/siresto-backend/internal/repository"
)

type CategoryService struct {
	Repo *repository.CategoryRepository
}

// GetAll fetches all categories
func (s *CategoryService) ListAllCategories() ([]domain.Category, error) {
	return s.Repo.ListAllCategories()
}

// GetByID fetches a category by its ID
func (s *CategoryService) GetCategoryByID(id uuid.UUID) (*domain.Category, error) {
	return s.Repo.GetCategoryByID(id)
}

// Create creates a new category
func (s *CategoryService) CreateCategory(category *domain.Category) (*domain.Category, error) {
	category.ID = uuid.New()
	if err := s.Repo.CreateCategory(category); err != nil {
		return nil, err
	}
	return category, nil
}

// Update updates an existing category
func (s *CategoryService) UpdateCategory(id uuid.UUID, update *domain.Category) (*domain.Category, error) {
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
	return s.Repo.DeleteCategory(id)
}
