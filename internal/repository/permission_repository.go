package repository

import (
	"github.com/google/uuid"
	"github.com/latoulicious/siresto-backend/internal/domain"
	"gorm.io/gorm"
)

type PermissionRepository struct {
	DB *gorm.DB
}

func (r *PermissionRepository) ListAllRoles() ([]domain.Permission, error) {
	var roles []domain.Permission
	err := r.DB.Find(&roles).Error
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *PermissionRepository) GetRoleByID(id uuid.UUID) (*domain.Permission, error) {
	var role domain.Permission
	err := r.DB.First(&role, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *PermissionRepository) CreateRole(role *domain.Permission) error {
	err := r.DB.Create(role).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *PermissionRepository) UpdateRole(role *domain.Permission) error {
	err := r.DB.Save(role).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *PermissionRepository) DeleteRole(id uuid.UUID) error {
	var role domain.Permission
	err := r.DB.First(&role, "id = ?", id).Error
	if err != nil {
		return err
	}
	err = r.DB.Delete(&role).Error
	if err != nil {
		return err
	}
	return nil
}
