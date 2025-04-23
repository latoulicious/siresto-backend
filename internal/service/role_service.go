package service

import (
	"fmt"

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
		return nil, fmt.Errorf("failed to get role: %w", err)
	}

	// Update basic fields
	existingRole.Name = req.Name
	existingRole.Description = req.Description

	// Always handle permissions explicitly - empty array will clear permissions
	permissions, err := s.PermissionService.GetPermissionsByIDs(req.Permissions)
	if err != nil && len(req.Permissions) > 0 {
		return nil, fmt.Errorf("failed to get permissions: %w", err)
	}

	// Use transaction to ensure atomicity
	err = s.DB.Transaction(func(tx *gorm.DB) error {
		// Clear existing permissions
		if err := tx.Model(existingRole).Association("Permissions").Clear(); err != nil {
			return fmt.Errorf("failed to clear permissions: %w", err)
		}

		// Add new permissions if any
		if len(permissions) > 0 {
			permPtrs := make([]*domain.Permission, 0, len(permissions))
			for i := range permissions {
				permPtrs = append(permPtrs, &permissions[i]) // Store pointer to array element
			}

			if err := tx.Model(existingRole).Association("Permissions").Append(permPtrs); err != nil {
				return fmt.Errorf("failed to add permissions: %w", err)
			}
		}

		// Update the role itself
		if err := tx.Save(existingRole).Error; err != nil {
			return fmt.Errorf("failed to save role: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Get the updated role with all associations loaded
	updatedRole, err := s.Repo.GetRoleByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated role: %w", err)
	}

	response := mapRoleToResponse(updatedRole)
	return &response, nil
}

// AddPermissionsToRole adds specific permissions to a role
func (s *RoleService) AddPermissionsToRole(roleID uuid.UUID, permissionIDs []uuid.UUID) (*dto.RoleResponse, error) {
	if len(permissionIDs) == 0 {
		return nil, fmt.Errorf("no permissions specified to add")
	}

	// Get the permissions to add
	newPermissions, err := s.PermissionService.GetPermissionsByIDs(permissionIDs)
	if err != nil {
		return nil, err
	}

	// Use transaction to ensure atomicity
	err = s.DB.Transaction(func(tx *gorm.DB) error {
		// Get existing role with its permissions
		existingRole, err := s.Repo.GetRoleByID(roleID)
		if err != nil {
			return err
		}

		// Track existing permission IDs to avoid duplicates
		existingPermIDs := make(map[uuid.UUID]bool)
		for _, perm := range existingRole.Permissions {
			existingPermIDs[perm.ID] = true
		}

		// Add permissions that don't already exist in the role
		for _, newPerm := range newPermissions {
			if !existingPermIDs[newPerm.ID] {
				permCopy := newPerm // Create a copy to avoid pointer issues
				existingRole.Permissions = append(existingRole.Permissions, &permCopy)
				existingPermIDs[newPerm.ID] = true
			}
		}

		return s.Repo.UpdateRole(existingRole)
	})

	if err != nil {
		return nil, err
	}

	// Get the updated role with preloaded permissions
	updatedRole, err := s.Repo.GetRoleByID(roleID)
	if err != nil {
		return nil, err
	}

	response := mapRoleToResponse(updatedRole)
	return &response, nil
}

// RemovePermissionsFromRole removes specific permissions from a role
func (s *RoleService) RemovePermissionsFromRole(roleID uuid.UUID, permissionIDs []uuid.UUID) (*dto.RoleResponse, error) {
	if len(permissionIDs) == 0 {
		return nil, fmt.Errorf("no permissions specified to remove")
	}

	// Get existing role with its permissions
	existingRole, err := s.Repo.GetRoleByID(roleID)
	if err != nil {
		return nil, err
	}

	// Convert to map for faster lookup
	toRemove := make(map[uuid.UUID]bool)
	for _, id := range permissionIDs {
		toRemove[id] = true
	}

	// Filter out permissions that should be removed
	updatedPermissions := make([]*domain.Permission, 0)
	for _, perm := range existingRole.Permissions {
		if !toRemove[perm.ID] {
			updatedPermissions = append(updatedPermissions, perm)
		}
	}

	// Update role with filtered permissions
	existingRole.Permissions = updatedPermissions

	// Save changes
	if err := s.Repo.UpdateRole(existingRole); err != nil {
		return nil, err
	}

	// Get updated role
	updatedRole, err := s.Repo.GetRoleByID(roleID)
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
