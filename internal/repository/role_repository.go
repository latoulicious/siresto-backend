package repository

import (
	"github.com/google/uuid"
	"github.com/latoulicious/siresto-backend/internal/domain"
	"gorm.io/gorm"
)

type RoleRepository struct {
	DB *gorm.DB
}

func (r *RoleRepository) ListAllRoles() ([]domain.Role, error) {
	var roles []domain.Role
	err := r.DB.Find(&roles).Error
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *RoleRepository) GetRoleByID(id uuid.UUID) (*domain.Role, error) {
	var role domain.Role
	err := r.DB.First(&role, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepository) CreateRole(role *domain.Role) error {
	err := r.DB.Create(role).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *RoleRepository) UpdateRole(role *domain.Role) error {
	err := r.DB.Save(role).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *RoleRepository) DeleteRole(id uuid.UUID) error {
	var role domain.Role
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
