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
	err := r.DB.Preload("Permissions").Find(&roles).Error
	return roles, err
}

func (r *RoleRepository) GetRoleByID(id uuid.UUID) (*domain.Role, error) {
	var role domain.Role
	err := r.DB.Preload("Permissions").First(&role, "id = ?", id).Error
	return &role, err
}

func (r *RoleRepository) CreateRole(role *domain.Role, tx *gorm.DB) error {
	db := r.DB
	if tx != nil {
		db = tx
	}
	return db.Create(role).Error
}

func (r *RoleRepository) UpdateRole(role *domain.Role) error {
	// Start a transaction
	return r.DB.Transaction(func(tx *gorm.DB) error {
		// Clear existing permission associations
		if err := tx.Model(role).Association("Permissions").Clear(); err != nil {
			return err
		}

		// Save the role with updated permissions
		return tx.Save(role).Error
	})
}

func (r *RoleRepository) DeleteRole(id uuid.UUID) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		// Get the role first
		var role domain.Role
		if err := tx.First(&role, "id = ?", id).Error; err != nil {
			return err
		}

		// Clear permission associations
		if err := tx.Model(&role).Association("Permissions").Clear(); err != nil {
			return err
		}

		// Delete the role
		return tx.Delete(&role).Error
	})
}

// New method to set role permissions
func (r *RoleRepository) SetRolePermissions(roleID uuid.UUID, permissions []domain.Permission, tx *gorm.DB) error {
	db := r.DB
	if tx != nil {
		db = tx
	}

	var role domain.Role
	if err := db.First(&role, "id = ?", roleID).Error; err != nil {
		return err
	}

	return db.Model(&role).Association("Permissions").Replace(permissions)
}
