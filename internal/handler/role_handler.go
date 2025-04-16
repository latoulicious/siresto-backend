package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/latoulicious/siresto-backend/internal/domain"
	"github.com/latoulicious/siresto-backend/internal/service"
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
	var role domain.Role
	if err := c.BodyParser(&role); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	if err := h.Service.CreateRole(&role); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create role",
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

	var role domain.Role
	if err := c.BodyParser(&role); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	role.ID = roleID
	if err := h.Service.UpdateRole(&role); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update role",
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
			"error": "Failed to delete role",
		})
	}
	return c.SendStatus(fiber.StatusNoContent)
}
