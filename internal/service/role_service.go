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

	// Set a default position if not specified (100 is a relatively low privilege level)
	if req.Position == 0 {
		req.Position = 100
	}

	// Check if the position is valid (must be > caller's position)
	callerPosition, err := s.GetCallerPosition()
	if err != nil {
		return nil, fmt.Errorf("failed to validate role position: %w", err)
	}

	// Cannot create a role with higher or equal privilege than the caller
	if req.Position <= callerPosition {
		return nil, fmt.Errorf("cannot create a role with higher privilege than your own role")
	}

	// Create role with permissions in a transaction
	var role domain.Role
	err = s.DB.Transaction(func(tx *gorm.DB) error {
		role = domain.Role{
			ID:          uuid.New(),
			Name:        req.Name,
			Description: req.Description,
			Position:    req.Position,
			IsSystem:    false, // New roles can't be system roles
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

	// Get caller's position to check permissions
	callerPosition, err := s.GetCallerPosition()
	if err != nil {
		return nil, fmt.Errorf("failed to validate permissions: %w", err)
	}

	// Cannot modify a role with higher or equal privilege
	if existingRole.Position <= callerPosition {
		return nil, fmt.Errorf("insufficient privileges to modify this role")
	}

	// Update basic fields if provided
	if req.Name != "" {
		existingRole.Name = req.Name
	}

	if req.Description != "" {
		existingRole.Description = req.Description
	}

	// Update position if provided
	if req.Position != nil {
		// Cannot set position to be higher or equal to caller's privilege
		if *req.Position <= callerPosition {
			return nil, fmt.Errorf("cannot set position to higher privilege than your role")
		}
		existingRole.Position = *req.Position
	}

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

		// Check if caller has sufficient privileges to modify this role
		callerPosition, err := s.GetCallerPosition()
		if err != nil {
			return fmt.Errorf("failed to validate permissions: %w", err)
		}

		// Cannot modify a role with higher or equal privilege
		if existingRole.Position <= callerPosition {
			return fmt.Errorf("insufficient privileges to modify this role")
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
	// Get the role to be deleted first
	role, err := s.Repo.GetRoleByID(id)
	if err != nil {
		return fmt.Errorf("failed to get role: %w", err)
	}

	// Check if caller has sufficient privileges to delete this role
	callerPosition, err := s.GetCallerPosition()
	if err != nil {
		return fmt.Errorf("failed to validate permissions: %w", err)
	}

	// Cannot delete a role with higher or equal privilege
	if role.Position <= callerPosition {
		return fmt.Errorf("insufficient privileges to delete this role")
	}

	// Cannot delete system roles
	if role.IsSystem {
		return fmt.Errorf("system roles cannot be deleted")
	}

	return s.Repo.DeleteRole(id)
}

// Helper function to map domain role to response DTO
func mapRoleToResponse(role *domain.Role) dto.RoleResponse {
	response := dto.RoleResponse{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		Position:    role.Position,
		IsSystem:    role.IsSystem,
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

// GetCallerPosition retrieves the position of the caller's role
// In a real implementation, this would extract the caller's role from a context
// For now, we use a middleware that sets role info in fiber.Ctx.Locals
func (s *RoleService) GetCallerPosition() (int, error) {
	// TODO: In a real implementation, get the role position from the request context
	// This would be implemented with context propagation through the service layer

	// For demonstration purposes, we return a default position of 2 (Owner)
	// In production, this would extract the actual role position from a user context
	return 2, nil
}
