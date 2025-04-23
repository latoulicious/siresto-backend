package repository

import (
	"github.com/google/uuid"
	"github.com/latoulicious/siresto-backend/internal/domain"
	"gorm.io/gorm"
)

type PermissionRepository struct {
	DB *gorm.DB
}

func (r *PermissionRepository) ListAllPermissions() ([]domain.Permission, error) {
	var permissions []domain.Permission
	err := r.DB.Find(&permissions).Error
	return permissions, err
}

func (r *PermissionRepository) GetPermissionByID(id uuid.UUID) (*domain.Permission, error) {
	var permission domain.Permission
	err := r.DB.First(&permission, "id = ?", id).Error
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
