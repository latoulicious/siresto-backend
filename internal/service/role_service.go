package service

import (
	"github.com/google/uuid"
	"github.com/latoulicious/siresto-backend/internal/domain"
	"github.com/latoulicious/siresto-backend/internal/repository"
)

type RoleService struct {
	Repo *repository.RoleRepository
}

func (s *RoleService) ListAllRoles() ([]domain.Role, error) {
	roles, err := s.Repo.ListAllRoles()
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func (s *RoleService) GetRoleByID(id uuid.UUID) (*domain.Role, error) {
	role, err := s.Repo.GetRoleByID(id)
	if err != nil {
		return nil, err
	}
	return role, nil
}

func (s *RoleService) CreateRole(role *domain.Role) error {
	role.ID = uuid.New()
	err := s.Repo.CreateRole(role)
	if err != nil {
		return err
	}
	return nil
}

func (s *RoleService) UpdateRole(role *domain.Role) error {
	err := s.Repo.UpdateRole(role)
	if err != nil {
		return err
	}
	return nil
}

func (s *RoleService) DeleteRole(id uuid.UUID) error {
	err := s.Repo.DeleteRole(id)
	if err != nil {
		return err
	}
	return nil
}
