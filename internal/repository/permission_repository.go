package repository

import (
	"github.com/google/uuid"
	"github.com/latoulicious/siresto-backend/internal/domain"
	"gorm.io/gorm"
)

type PermissionRepository struct {
	DB *gorm.DB
}

// ListAllPermissions now supports pagination
func (r *PermissionRepository) ListAllPermissions(page, limit int) ([]domain.Permission, int64, error) {
	var permissions []domain.Permission
	var totalCount int64

	// Get total count first
	if err := r.DB.Model(&domain.Permission{}).Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Get paginated data with roles preloaded
	err := r.DB.Preload("Roles").
		Limit(limit).
		Offset(offset).
		Find(&permissions).Error

	return permissions, totalCount, err
}

func (r *PermissionRepository) GetPermissionByID(id uuid.UUID) (*domain.Permission, error) {
	var permission domain.Permission
	err := r.DB.Preload("Roles").First(&permission, "id = ?", id).Error
	return &permission, err
}

func (r *PermissionRepository) CreatePermission(permission *domain.Permission) error {
	return r.DB.Create(permission).Error
}

func (r *PermissionRepository) UpdatePermission(permission *domain.Permission) error {
	return r.DB.Save(permission).Error
}

func (r *PermissionRepository) DeletePermission(id uuid.UUID) error {
	return r.DB.Delete(&domain.Permission{}, "id = ?", id).Error
}

// New method to get permissions by IDs for role association
func (r *PermissionRepository) GetPermissionsByIDs(ids []uuid.UUID) ([]domain.Permission, error) {
	var permissions []domain.Permission
	err := r.DB.Where("id IN ?", ids).Find(&permissions).Error
	return permissions, err
}

// GetPermissionByName retrieves a permission by its name
func (r *PermissionRepository) GetPermissionByName(name string) (*domain.Permission, error) {
	var permission domain.Permission
	err := r.DB.Where("name = ?", name).First(&permission).Error
	return &permission, err
}
