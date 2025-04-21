package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/latoulicious/siresto-backend/internal/domain"
	"github.com/latoulicious/siresto-backend/internal/service"
	"github.com/latoulicious/siresto-backend/internal/utils"
)

type ThemeHandler struct {
	Service *service.ThemeService
}

func (h *ThemeHandler) ListAllThemes(c *fiber.Ctx) error {
	themes, err := h.Service.ListAllThemes()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.Error("Failed to retrieve Themes", fiber.StatusInternalServerError))
	}
	return c.Status(fiber.StatusOK).JSON(utils.Success("Themes retrieved successfully", themes))
}

func (h *ThemeHandler) GetThemeByID(c *fiber.Ctx) error {
	id := c.Params("id")
	themeID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("Invalid Theme ID format", fiber.StatusBadRequest))
	}

	theme, err := h.Service.GetThemeByID(themeID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(utils.Error("Theme ID not found", fiber.StatusNotFound))
	}
	return c.Status(fiber.StatusOK).JSON(utils.Success("Theme retrieved successfully", theme))
}

func (h *ThemeHandler) CreateTheme(c *fiber.Ctx) error {
	var theme domain.Theme
	if err := c.BodyParser(&theme); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("Invalid request payload", fiber.StatusBadRequest))
	}

	if err := h.Service.CreateTheme(&theme); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.Error("Failed to create theme", fiber.StatusInternalServerError))
	}
	return c.Status(fiber.StatusCreated).JSON(utils.Success("Theme created successfully", theme))
}

func (h *ThemeHandler) UpdateTheme(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("Invalid theme ID format", fiber.StatusBadRequest))
	}

	var theme domain.Theme
	if err := c.BodyParser(&theme); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("Invalid request payload", fiber.StatusBadRequest))
	}

	theme.ID = id
	if err := h.Service.UpdateTheme(&theme); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.Error("Failed to update theme", fiber.StatusInternalServerError))
	}
	return c.Status(fiber.StatusOK).JSON(utils.Success("Theme updated successfully", theme))
}

func (h *ThemeHandler) DeleteTheme(c *fiber.Ctx) error {
	id := c.Params("id")
	themeID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("Invalid Theme ID format", fiber.StatusBadRequest))
	}

	if err := h.Service.DeleteTheme(themeID); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(utils.Error("Theme ID not found", fiber.StatusNotFound))
	}
	return c.Status(fiber.StatusNoContent).JSON(utils.Success("Theme deleted successfully", nil))
}
