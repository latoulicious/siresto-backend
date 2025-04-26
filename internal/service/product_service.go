package service

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/latoulicious/siresto-backend/internal/domain"
	"github.com/latoulicious/siresto-backend/internal/repository"
	"github.com/latoulicious/siresto-backend/internal/utils"
	"github.com/latoulicious/siresto-backend/pkg/dto"
	"gorm.io/gorm"
)

type ProductService struct {
	Repo     *repository.ProductRepository
	Uploader utils.Uploader
}

// ListAllProducts retrieves all products from the repository
// If limit > 0, pagination is enabled and offset must be provided
func (s *ProductService) ListAllProducts(offset, limit int) ([]domain.Product, int, error) {
	return s.Repo.ListProductsPaginated(offset, limit)
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
	if err := ValidateProduct(s.Repo.DB, product); err != nil {
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

func (s *ProductService) CreateProductWithVariations(product *domain.Product, variations []dto.CreateVariationRequest) (*domain.Product, []*domain.Variation, error) {
	// Start a database transaction to ensure atomicity
	tx := s.Repo.DB.Begin()

	// Create the product
	if err := tx.Create(&product).Error; err != nil {
		tx.Rollback()
		return nil, nil, fmt.Errorf("failed to create product: %w", err)
	}

	// Prepare a slice to hold the created variations (pointer slice)
	var createdVariations []*domain.Variation

	// Handle variations, if provided
	for _, variationDTO := range variations {
		// Map the DTO to domain variation
		variation := dto.ToVariationDomain(&variationDTO)
		variation.ProductID = product.ID

		// Save the variation
		if err := tx.Create(&variation).Error; err != nil {
			tx.Rollback()
			return nil, nil, fmt.Errorf("failed to create variation: %w", err)
		}

		// Append the created variation to the slice (pointers)
		createdVariations = append(createdVariations, variation)
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return nil, nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Return the created product and its variations
	return product, createdVariations, nil
}

// UpdateProduct updates an existing product in the repository
func (s *ProductService) UpdateProduct(id uuid.UUID, update *domain.Product) (*domain.Product, error) {
	// Fetch the existing product
	existing, err := s.Repo.GetProductByID(id)
	if err != nil {
		return nil, err
	}

	// Perform validation
	if err := ValidateProductForUpdate(s.Repo.DB, update); err != nil {
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

func (s *ProductService) UpdateProductWithVariations(id uuid.UUID, product *domain.Product, variations []dto.UpdateVariationRequest, removeOthers bool) (*domain.Product, []*domain.Variation, error) {
	// Start transaction
	tx := s.Repo.DB.Begin()
	if tx.Error != nil {
		return nil, nil, fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	// Fetch existing product with variations
	existingProduct := &domain.Product{}
	if err := tx.Preload("Variations").First(existingProduct, "id = ?", id).Error; err != nil {
		tx.Rollback()
		return nil, nil, fmt.Errorf("product not found: %w", err)
	}

	// Update product fields
	existingProduct.Name = product.Name
	existingProduct.Description = product.Description
	existingProduct.ImageURL = product.ImageURL
	existingProduct.BasePrice = product.BasePrice
	existingProduct.IsAvailable = product.IsAvailable
	existingProduct.Position = product.Position
	if product.CategoryID != nil {
		existingProduct.CategoryID = product.CategoryID
	}

	// Save product updates
	if err := tx.Save(existingProduct).Error; err != nil {
		tx.Rollback()
		return nil, nil, fmt.Errorf("failed to update product: %w", err)
	}

	// Handle variations only if they're provided
	var updatedVariations []*domain.Variation

	// Map existing variations by ID for efficient lookup
	existingVariationsMap := make(map[uuid.UUID]*domain.Variation)
	for i := range existingProduct.Variations {
		existingVariationsMap[existingProduct.Variations[i].ID] = &existingProduct.Variations[i]
	}

	// Track which variation IDs are being updated
	processedIDs := make(map[uuid.UUID]bool)

	// Process each variation in the request
	for _, varDTO := range variations {
		if varDTO.ID != nil {
			// Update existing variation
			if existingVar, found := existingVariationsMap[*varDTO.ID]; found {
				// Update fields that are present in the request
				if varDTO.IsDefault != nil {
					existingVar.IsDefault = *varDTO.IsDefault
				}
				if varDTO.IsAvailable != nil {
					existingVar.IsAvailable = *varDTO.IsAvailable
				}
				if varDTO.IsRequired != nil {
					existingVar.IsRequired = *varDTO.IsRequired
				}
				if varDTO.VariationType != nil {
					existingVar.VariationType = *varDTO.VariationType
				}

				// Handle options update if provided
				if len(varDTO.Options) > 0 {
					existingVar.Options = dto.ToVariationOptionsDomainFromUpdate(varDTO.Options)
				}

				// Save the updated variation
				if err := tx.Save(existingVar).Error; err != nil {
					tx.Rollback()
					return nil, nil, fmt.Errorf("failed to update variation: %w", err)
				}

				updatedVariations = append(updatedVariations, existingVar)
				processedIDs[*varDTO.ID] = true
			}
		} else {
			// Create new variation
			newVariation := &domain.Variation{
				ID:        uuid.New(),
				ProductID: existingProduct.ID,
			}

			// Set optional fields
			if varDTO.IsDefault != nil {
				newVariation.IsDefault = *varDTO.IsDefault
			}
			if varDTO.IsAvailable != nil {
				newVariation.IsAvailable = *varDTO.IsAvailable
			}
			if varDTO.IsRequired != nil {
				newVariation.IsRequired = *varDTO.IsRequired
			}
			if varDTO.VariationType != nil {
				newVariation.VariationType = *varDTO.VariationType
			}

			// Handle options
			if len(varDTO.Options) > 0 {
				newVariation.Options = dto.ToVariationOptionsDomainFromUpdate(varDTO.Options)
			}

			if err := tx.Create(newVariation).Error; err != nil {
				tx.Rollback()
				return nil, nil, fmt.Errorf("failed to create new variation: %w", err)
			}

			updatedVariations = append(updatedVariations, newVariation)
		}
	}

	// Delete variations that aren't in the update request if removeOthers flag is set
	if removeOthers && len(processedIDs) > 0 {
		for id := range existingVariationsMap {
			if !processedIDs[id] {
				if err := tx.Delete(&domain.Variation{}, "id = ?", id).Error; err != nil {
					tx.Rollback()
					return nil, nil, fmt.Errorf("failed to delete variation: %w", err)
				}
			}
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Return updated entities
	return existingProduct, updatedVariations, nil
}

// DeleteProduct removes a product by its ID from the repository
func (s *ProductService) DeleteProduct(id uuid.UUID) error {
	// Perform deletability validation
	if err := ValidateProductDeletable(s.Repo.DB, id); err != nil {
		return fmt.Errorf("product cannot be deleted: %w", err)
	}

	// Proceed with deletion
	return s.Repo.DeleteProduct(id)
}

// Helper Function

var (
	ErrMissingCategoryID = errors.New("categoryID is required")
	ErrInvalidBasePrice  = errors.New("basePrice must be greater than 0")
	ErrEmptyName         = errors.New("product name cannot be empty")
	ErrCategoryNotFound  = errors.New("category does not exist")
)

// ValidateProduct performs all validation on a product for creation or update
func ValidateProduct(db *gorm.DB, p *domain.Product) error {
	if p.CategoryID == nil || *p.CategoryID == uuid.Nil {
		return ErrMissingCategoryID
	}
	if p.Name == "" {
		return ErrEmptyName
	}
	if p.BasePrice <= 0 {
		return ErrInvalidBasePrice
	}
	exists, err := repository.CategoryExists(db, *p.CategoryID)
	if err != nil {
		return fmt.Errorf("error checking category: %w", err)
	}
	if !exists {
		return ErrCategoryNotFound
	}
	return nil
}

// ValidateProductForUpdate performs validation on a product for update
func ValidateProductForUpdate(db *gorm.DB, p *domain.Product) error {
	if p.Name == "" {
		return ErrEmptyName
	}
	if p.BasePrice <= 0 {
		return ErrInvalidBasePrice
	}
	// only check category if it's explicitly allowed to change
	if p.CategoryID != nil && *p.CategoryID != uuid.Nil {
		exists, err := repository.CategoryExists(db, *p.CategoryID)
		if err != nil {
			return fmt.Errorf("error checking category: %w", err)
		}
		if !exists {
			return ErrCategoryNotFound
		}
	}
	return nil
}

// ValidateProductDeletable ensures the product can be safely deleted
func ValidateProductDeletable(db *gorm.DB, productID uuid.UUID) error {
	var count int64
	if err := db.Model(&domain.Variation{}).
		Where("product_id = ?", productID).
		Count(&count).Error; err != nil {
		return fmt.Errorf("error checking related variations: %w", err)
	}

	if count > 0 {
		return fmt.Errorf("cannot delete product: %d associated variations found", count)
	}

	return nil
}
