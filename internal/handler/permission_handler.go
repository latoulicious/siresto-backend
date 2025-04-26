package handler

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/latoulicious/siresto-backend/internal/domain"
	"github.com/latoulicious/siresto-backend/internal/service"
	"github.com/latoulicious/siresto-backend/internal/utils"
	"github.com/latoulicious/siresto-backend/internal/validator"
)

type PermissionHandler struct {
	Service *service.PermissionService
}

func (h *PermissionHandler) ListAllPermissions(c *fiber.Ctx) error {
	// Get page from query parameter, default to 1 if not provided
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	// Get permissions with pagination
	permissions, totalCount, err := h.Service.ListAllPermissions(page)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.Error(
			"Failed to fetch permissions",
			fiber.StatusInternalServerError,
			utils.NewErrorInfo("FETCH_ERROR", err.Error(), "", nil),
		))
	}

	// Create pagination metadata
	const itemsPerPage = 10
	metadata := utils.NewPaginationMetadata(page, itemsPerPage, int(totalCount))

	return c.Status(fiber.StatusOK).JSON(utils.Success(
		"Permissions retrieved successfully",
		permissions,
		metadata,
	))
}

func (h *PermissionHandler) GetPermissionByID(c *fiber.Ctx) error {
	id := c.Params("id")
	roleID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("Invalid permission ID format", fiber.StatusBadRequest))
	}

	permission, err := h.Service.GetPermissionByID(roleID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(utils.Error("Permission not found", fiber.StatusNotFound))
	}
	return c.Status(fiber.StatusOK).JSON(utils.Success("Permission retrieved successfully", permission))
}

func (h *PermissionHandler) CreatePermission(c *fiber.Ctx) error {
	var permission domain.Permission
	if err := c.BodyParser(&permission); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("Invalid request payload", fiber.StatusBadRequest))
	}

	// Validate permission format
	if err := validator.ValidatePermission(&permission); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error(
			err.Error(),
			fiber.StatusBadRequest,
			utils.NewErrorInfo("VALIDATION_ERROR", "Permission format is invalid", "", nil),
		))
	}

	// Auto-generate description if not provided
	if permission.Description == "" {
		permission.Description = validator.GetPermissionDescription(permission.Name)
	}

	if err := h.Service.CreatePermission(&permission); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.Error("Failed to create permission", fiber.StatusInternalServerError))
	}
	return c.Status(fiber.StatusCreated).JSON(utils.Success("Permission created successfully", permission))
}

func (h *PermissionHandler) UpdatePermission(c *fiber.Ctx) error {
	id := c.Params("id")
	roleID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("Invalid permission ID format", fiber.StatusBadRequest))
	}

	var permission domain.Permission
	if err := c.BodyParser(&permission); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("Invalid request payload", fiber.StatusBadRequest))
	}

	// Validate permission format if name is being updated
	if permission.Name != "" {
		if err := validator.ValidatePermission(&permission); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(utils.Error(
				err.Error(),
				fiber.StatusBadRequest,
				utils.NewErrorInfo("VALIDATION_ERROR", "Permission format is invalid", "", nil),
			))
		}
	}

	// Auto-generate description if updating name but not description
	if permission.Name != "" && permission.Description == "" {
		permission.Description = validator.GetPermissionDescription(permission.Name)
	}

	permission.ID = roleID
	if err := h.Service.UpdatePermission(&permission); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.Error("Failed to update permission", fiber.StatusInternalServerError))
	}
	return c.Status(fiber.StatusOK).JSON(utils.Success("Permission updated successfully", permission))
}

func (h *PermissionHandler) DeletePermission(c *fiber.Ctx) error {
	id := c.Params("id")
	roleID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("Invalid permission ID format", fiber.StatusBadRequest))
	}

	if err := h.Service.DeletePermission(roleID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.Error("Failed to delete permission", fiber.StatusInternalServerError))
	}
	return c.Status(fiber.StatusNoContent).JSON(utils.Success("Permission deleted successfully", nil))
}

// GenerateResourcePermissions handles bulk permission generation for a new resource
func (h *PermissionHandler) GenerateResourcePermissions(c *fiber.Ctx) error {
	// Parse request
	type ResourcePermissionRequest struct {
		ResourceName  string `json:"resource_name" validate:"required"`
		IncludeManage bool   `json:"include_manage"`
	}

	var req ResourcePermissionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("Invalid request payload", fiber.StatusBadRequest))
	}

	// Validate and format the resource name
	resourceName, err := utils.ValidateAndFormatResource(req.ResourceName)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error(
			err.Error(),
			fiber.StatusBadRequest,
			utils.NewErrorInfo("VALIDATION_ERROR", "Resource name is invalid", "", nil),
		))
	}

	// Generate the permissions
	var permissions []domain.Permission
	if req.IncludeManage {
		permissions = utils.GeneratePermissionBundle(resourceName)
	} else {
		permissions = utils.GenerateCRUDPermissions(resourceName)
	}

	// Save all permissions in a transaction
	createdPermissions := make([]domain.Permission, 0, len(permissions))
	for _, perm := range permissions {
		if err := h.Service.CreatePermission(&perm); err != nil {
			// If permission already exists, continue to the next one
			continue
		}
		createdPermissions = append(createdPermissions, perm)
	}

	return c.Status(fiber.StatusCreated).JSON(utils.Success(
		fmt.Sprintf("Created %d permissions for resource '%s'", len(createdPermissions), req.ResourceName),
		createdPermissions,
	))
}
