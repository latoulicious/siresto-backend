package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/latoulicious/siresto-backend/internal/domain"
	"github.com/latoulicious/siresto-backend/internal/service"
)

type PermissionHandler struct {
	Service *service.PermissionService
}

func (h *PermissionHandler) ListAllPermissions(c *fiber.Ctx) error {
	permissions, err := h.Service.ListAllPermissions()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch permissions",
		})
	}
	return c.JSON(permissions)
}

func (h *PermissionHandler) GetPermissionByID(c *fiber.Ctx) error {
	id := c.Params("id")
	roleID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid permission ID",
		})
	}

	permission, err := h.Service.GetPermissionByID(roleID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Permission not found",
		})
	}
	return c.JSON(permission)
}

func (h *PermissionHandler) CreatePermission(c *fiber.Ctx) error {
	var permission domain.Permission
	if err := c.BodyParser(&permission); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	if err := h.Service.CreatePermission(&permission); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create permission",
		})
	}
	return c.Status(fiber.StatusCreated).JSON(permission)
}

func (h *PermissionHandler) UpdatePermission(c *fiber.Ctx) error {
	id := c.Params("id")
	roleID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid permission ID",
		})
	}

	var permission domain.Permission
	if err := c.BodyParser(&permission); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	permission.ID = roleID
	if err := h.Service.UpdatePermission(&permission); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update permission",
		})
	}
	return c.JSON(permission)
}

func (h *PermissionHandler) DeletePermission(c *fiber.Ctx) error {
	id := c.Params("id")
	roleID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid permission ID",
		})
	}

	if err := h.Service.DeletePermission(roleID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete permission",
		})
	}
	return c.SendStatus(fiber.StatusNoContent)
}
