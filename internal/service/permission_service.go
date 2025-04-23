package service

import (
	"github.com/google/uuid"
	"github.com/latoulicious/siresto-backend/internal/domain"
	"github.com/latoulicious/siresto-backend/internal/repository"
)

type PermissionService struct {
	Repo *repository.PermissionRepository
}

func (s *PermissionService) ListAllRoles() ([]domain.Permission, error) {
	permissions, err := s.Repo.ListAllRoles()
	if err != nil {
		return nil, err
	}
	return permissions, nil
}

func (s *PermissionService) GetRoleByID(id uuid.UUID) (*domain.Permission, error) {
	permission, err := s.Repo.GetRoleByID(id)
	if err != nil {
		return nil, err
	}
	return permission, nil
}

func (s *PermissionService) CreateRole(permission *domain.Permission) error {
	permission.ID = uuid.New()
	err := s.Repo.CreateRole(permission)
	if err != nil {
		return err
	}
	return nil
}

func (s *PermissionService) UpdateRole(permission *domain.Permission) error {
	err := s.Repo.UpdateRole(permission)
	if err != nil {
		return err
	}
	return nil
}

func (s *PermissionService) DeleteRole(id uuid.UUID) error {
	err := s.Repo.DeleteRole(id)
	if err != nil {
		return err
	}
	return nil
}
