package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/latoulicious/siresto-backend/internal/domain"
	"github.com/latoulicious/siresto-backend/internal/service"
	"github.com/latoulicious/siresto-backend/internal/utils"
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
