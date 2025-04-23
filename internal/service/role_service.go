package service

import (
	"github.com/google/uuid"
	"github.com/latoulicious/siresto-backend/internal/domain"
	"github.com/latoulicious/siresto-backend/internal/repository"
	"github.com/latoulicious/siresto-backend/pkg/dto"
	"gorm.io/gorm"
)

type RoleService struct {
	Repo              *repository.RoleRepository
	PermissionService *PermissionService
	DB                *gorm.DB
}

func (s *RoleService) ListAllRoles() ([]dto.RoleResponse, error) {
	roles, err := s.Repo.ListAllRoles()
	if err != nil {
		return nil, err
	}

	// Map to response DTOs
	responses := make([]dto.RoleResponse, 0, len(roles))
	for _, role := range roles {
		responses = append(responses, mapRoleToResponse(&role))
	}

	return responses, nil
}

func (s *RoleService) GetRoleByID(id uuid.UUID) (*dto.RoleResponse, error) {
	role, err := s.Repo.GetRoleByID(id)
	if err != nil {
		return nil, err
	}

	response := mapRoleToResponse(role)
	return &response, nil
}

func (s *RoleService) CreateRole(req *dto.CreateRoleRequest) (*dto.RoleResponse, error) {
	// Get permissions if any are specified
	var permissions []domain.Permission
	var err error

	if len(req.Permissions) > 0 {
		permissions, err = s.PermissionService.GetPermissionsByIDs(req.Permissions)
		if err != nil {
			return nil, err
		}
	}

	// Create role with permissions in a transaction
	var role domain.Role
	err = s.DB.Transaction(func(tx *gorm.DB) error {
		role = domain.Role{
			ID:          uuid.New(),
			Name:        req.Name,
			Description: req.Description,
			Permissions: make([]*domain.Permission, len(permissions)),
		}

		// Convert permissions to pointers
		for i, p := range permissions {
			pCopy := p // Create a copy to avoid pointer issues
			role.Permissions[i] = &pCopy
		}

		return s.Repo.CreateRole(&role, tx)
	})

	if err != nil {
		return nil, err
	}

	// Get the complete role with associations
	completeRole, err := s.Repo.GetRoleByID(role.ID)
	if err != nil {
		return nil, err
	}

	response := mapRoleToResponse(completeRole)
	return &response, nil
}

func (s *RoleService) UpdateRole(id uuid.UUID, req *dto.UpdateRoleRequest) (*dto.RoleResponse, error) {
	// Get existing role
	existingRole, err := s.Repo.GetRoleByID(id)
	if err != nil {
		return nil, err
	}

	// Get permissions if specified
	var permissions []domain.Permission
	if len(req.Permissions) > 0 {
		permissions, err = s.PermissionService.GetPermissionsByIDs(req.Permissions)
		if err != nil {
			return nil, err
		}
	}

	// Update fields
	existingRole.Name = req.Name
	existingRole.Description = req.Description

	// Clear and set new permissions
	existingRole.Permissions = make([]*domain.Permission, len(permissions))
	for i, p := range permissions {
		pCopy := p // Create a copy to avoid pointer issues
		existingRole.Permissions[i] = &pCopy
	}

	// Save changes
	if err := s.Repo.UpdateRole(existingRole); err != nil {
		return nil, err
	}

	// Get updated role
	updatedRole, err := s.Repo.GetRoleByID(id)
	if err != nil {
		return nil, err
	}

	response := mapRoleToResponse(updatedRole)
	return &response, nil
}

func (s *RoleService) DeleteRole(id uuid.UUID) error {
	return s.Repo.DeleteRole(id)
}

// Helper function to map domain role to response DTO
func mapRoleToResponse(role *domain.Role) dto.RoleResponse {
	response := dto.RoleResponse{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		Permissions: make([]dto.PermissionResponse, 0, len(role.Permissions)),
	}

	for _, perm := range role.Permissions {
		response.Permissions = append(response.Permissions, dto.PermissionResponse{
			ID:          perm.ID,
			Name:        perm.Name,
			Description: perm.Description,
		})
	}

	return response
}
