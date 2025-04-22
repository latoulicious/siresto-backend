package service

import (
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/latoulicious/siresto-backend/internal/domain"
	"github.com/latoulicious/siresto-backend/internal/repository"
	"github.com/latoulicious/siresto-backend/pkg/dto"
)

type CategoryService struct {
	Repo *repository.CategoryRepository
}

func (s *CategoryService) ListAllCategories(includeProducts bool) ([]domain.Category, error) {
	if includeProducts {
		// Load categories with products included
		return s.Repo.ListAllCategoriesWithProducts()
	}
	// Load categories without products
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
	if strings.TrimSpace(category.Name) == "" {
		return nil, errors.New("category name cannot be empty")
	}
	if category.Position < 0 {
		return nil, errors.New("category position cannot be negative")
	}
	exists, err := s.Repo.ExistsByName(category.Name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("category with the same name already exists")
	}

	category.ID = uuid.New()
	if err := s.Repo.CreateCategory(category); err != nil {
		return nil, err
	}

	category.Products = []domain.Product{}
	return category, nil
}

// Update updates an existing category
func (s *CategoryService) UpdateCategory(id uuid.UUID, update *dto.UpdateCategoryRequest) (*domain.Category, error) {
	existing, err := s.Repo.GetCategoryByID(id)
	if err != nil {
		return nil, err
	}

	if update.Name != nil {
		if strings.TrimSpace(*update.Name) == "" {
			return nil, errors.New("category name cannot be empty")
		}
		exists, err := s.Repo.ExistsByNameExcludingID(*update.Name, existing.ID)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("category with the same name already exists")
		}
		existing.Name = *update.Name
	}

	if update.Position != nil {
		if *update.Position < 0 {
			return nil, errors.New("category position cannot be negative")
		}
		existing.Position = *update.Position
	}

	if update.IsActive != nil {
		existing.IsActive = *update.IsActive
	}

	if err := s.Repo.UpdateCategory(existing); err != nil {
		return nil, err
	}

	return existing, nil
}

// Delete removes a category by ID
func (s *CategoryService) DeleteCategory(id uuid.UUID) error {
	hasProducts, err := s.Repo.HasAssociatedProducts(id)
	if err != nil {
		return err
	}
	if hasProducts {
		return fmt.Errorf("cannot delete category: it has associated products")
	}

	return s.Repo.DeleteCategory(id)
}
