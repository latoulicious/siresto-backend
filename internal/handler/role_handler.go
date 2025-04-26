package handler

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/latoulicious/siresto-backend/internal/domain"
	"github.com/latoulicious/siresto-backend/internal/service"
	"github.com/latoulicious/siresto-backend/internal/utils"
	"github.com/latoulicious/siresto-backend/pkg/dto"
)

type RoleHandler struct {
	Service *service.RoleService
}

func (h *RoleHandler) ListAllRoles(c *fiber.Ctx) error {
	roles, err := h.Service.ListAllRoles()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch roles",
		})
	}
	return c.JSON(roles)
}

func (h *RoleHandler) GetRoleByID(c *fiber.Ctx) error {
	id := c.Params("id")
	roleID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid role ID",
		})
	}

	role, err := h.Service.GetRoleByID(roleID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Role not found",
		})
	}
	return c.JSON(role)
}

func (h *RoleHandler) CreateRole(c *fiber.Ctx) error {
	var request dto.CreateRoleRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	role, err := h.Service.CreateRole(&request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create role: " + err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(role)
}

func (h *RoleHandler) UpdateRole(c *fiber.Ctx) error {
	id := c.Params("id")
	roleID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid role ID",
		})
	}

	var request dto.UpdateRoleRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	role, err := h.Service.UpdateRole(roleID, &request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update role: " + err.Error(),
		})
	}
	return c.JSON(role)
}

// New handler method for adding permissions to a role
func (h *RoleHandler) AddPermissionsToRole(c *fiber.Ctx) error {
	id := c.Params("id")
	roleID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid role ID",
		})
	}

	var request dto.RolePermissionUpdateRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	role, err := h.Service.AddPermissionsToRole(roleID, request.Permissions)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to add permissions to role: " + err.Error(),
		})
	}
	return c.JSON(role)
}

// New handler method for removing permissions from a role
func (h *RoleHandler) RemovePermissionsFromRole(c *fiber.Ctx) error {
	id := c.Params("id")
	roleID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid role ID",
		})
	}

	var request dto.RolePermissionUpdateRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	role, err := h.Service.RemovePermissionsFromRole(roleID, request.Permissions)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to remove permissions from role: " + err.Error(),
		})
	}
	return c.JSON(role)
}

func (h *RoleHandler) DeleteRole(c *fiber.Ctx) error {
	id := c.Params("id")
	roleID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid role ID",
		})
	}

	if err := h.Service.DeleteRole(roleID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete role: " + err.Error(),
		})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// CreateComprehensiveRole creates a role with extensive permissions (but not system-level)
func (h *RoleHandler) CreateComprehensiveRole(c *fiber.Ctx) error {
	// Parse request
	type ComprehensiveRoleRequest struct {
		Name        string   `json:"name" validate:"required"`
		Description string   `json:"description"`
		Resources   []string `json:"resources" validate:"required,min=1"`
	}

	var req ComprehensiveRoleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("Invalid request payload", fiber.StatusBadRequest))
	}

	// Create a new role
	createRoleReq := &dto.CreateRoleRequest{
		Name:        req.Name,
		Description: req.Description,
		Permissions: []uuid.UUID{}, // Will be filled with permissions
	}

	// For each resource, generate all CRUD + manage permissions
	allPermissionIDs := make([]uuid.UUID, 0)

	// First create permissions for all resources
	for _, resource := range req.Resources {
		// Format and validate the resource name
		resourceName, err := utils.ValidateAndFormatResource(resource)
		if err != nil {
			continue // Skip invalid resources
		}

		// Generate CRUD + manage permissions for this resource
		permissions := utils.GeneratePermissionBundle(resourceName)

		// Save each permission
		for i := range permissions {
			// Try to create the permission (it might already exist)
			if err := h.Service.PermissionService.CreatePermission(&permissions[i]); err != nil {
				// Try to fetch it if creation fails (likely due to already existing)
				existingPerm, fetchErr := h.Service.PermissionService.GetPermissionByName(permissions[i].Name)
				if fetchErr != nil {
					continue // Skip if we can't get it
				}
				// Use the existing permission's ID
				allPermissionIDs = append(allPermissionIDs, existingPerm.ID)
			} else {
				// Use the newly created permission's ID
				allPermissionIDs = append(allPermissionIDs, permissions[i].ID)
			}
		}
	}

	// Now add all the standard management permissions
	standardManagementPermissions := []string{
		"manage:users",
		"manage:roles",
		"manage:menu",
		"manage:orders",
		"manage:tables",
		"manage:inventory",
		"manage:reports",
		"manage:settings",
	}

	for _, permName := range standardManagementPermissions {
		// Create a placeholder permission to check if it exists
		placeholderPerm := &domain.Permission{
			ID:   uuid.New(),
			Name: permName,
			Description: fmt.Sprintf("Permission to manage all aspects of %s",
				strings.ReplaceAll(strings.TrimPrefix(permName, "manage:"), "_", " ")),
		}

		// Try to create the permission (it might already exist)
		if err := h.Service.PermissionService.CreatePermission(placeholderPerm); err != nil {
			// Try to fetch it if creation fails
			existingPerm, fetchErr := h.Service.PermissionService.GetPermissionByName(permName)
			if fetchErr != nil {
				continue // Skip if we can't get it
			}
			// Use the existing permission's ID
			allPermissionIDs = append(allPermissionIDs, existingPerm.ID)
		} else {
			// Use the newly created permission's ID
			allPermissionIDs = append(allPermissionIDs, placeholderPerm.ID)
		}
	}

	// Set the permissions for the role
	createRoleReq.Permissions = allPermissionIDs

	// Create the role with all the permissions
	role, err := h.Service.CreateRole(createRoleReq)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.Error(
			"Failed to create comprehensive role",
			fiber.StatusInternalServerError,
			utils.NewErrorInfo("ROLE_CREATION_ERROR", err.Error(), "", nil),
		))
	}

	return c.Status(fiber.StatusCreated).JSON(utils.Success(
		fmt.Sprintf("Created comprehensive role '%s' with %d permissions",
			req.Name, len(allPermissionIDs)),
		role,
	))
}
