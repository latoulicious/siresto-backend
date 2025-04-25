package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/latoulicious/siresto-backend/internal/domain"
	"github.com/latoulicious/siresto-backend/internal/repository"
	"gorm.io/gorm"
)

type PermissionService struct {
	Repo *repository.PermissionRepository
	DB   *gorm.DB
}

// ListAllPermissions now supports pagination
func (s *PermissionService) ListAllPermissions(page int) ([]domain.Permission, int64, error) {
	const itemsPerPage = 10
	if page < 1 {
		page = 1
	}
	return s.Repo.ListAllPermissions(page, itemsPerPage)
}

func (s *PermissionService) GetPermissionByID(id uuid.UUID) (*domain.Permission, error) {
	return s.Repo.GetPermissionByID(id)
}

func (s *PermissionService) CreatePermission(permission *domain.Permission) error {
	permission.ID = uuid.New()
	return s.Repo.CreatePermission(permission)
}

func (s *PermissionService) UpdatePermission(permission *domain.Permission) error {
	return s.Repo.UpdatePermission(permission)
}

func (s *PermissionService) DeletePermission(id uuid.UUID) error {
	return s.Repo.DeletePermission(id)
}

func (s *PermissionService) GetPermissionsByIDs(ids []uuid.UUID) ([]domain.Permission, error) {
	permissions, err := s.Repo.GetPermissionsByIDs(ids)
	if err != nil {
		return nil, err
	}

	// Verify we got all the permissions requested
	if len(permissions) != len(ids) {
		return nil, errors.New("one or more permissions not found")
	}

	return permissions, nil
}

// GetPermissionByName retrieves a permission by its name
func (s *PermissionService) GetPermissionByName(name string) (*domain.Permission, error) {
	return s.Repo.GetPermissionByName(name)
}
