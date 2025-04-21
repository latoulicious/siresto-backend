package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/latoulicious/siresto-backend/internal/domain"
	"github.com/latoulicious/siresto-backend/internal/service"
)

type ThemeHandler struct {
	Service *service.ThemeService
}

func (h *ThemeHandler) ListAllThemes(c *fiber.Ctx) error {
	themes, err := h.Service.ListAllThemes()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch themes",
		})
	}
	return c.JSON(themes)
}

func (h *ThemeHandler) GetThemeByID(c *fiber.Ctx) error {
	id := c.Params("id")
	themeID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid theme ID",
		})
	}

	theme, err := h.Service.GetThemeByID(themeID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Theme not found",
		})
	}
	return c.JSON(theme)
}

func (h *ThemeHandler) CreateTheme(c *fiber.Ctx) error {
	var theme domain.Theme
	if err := c.BodyParser(&theme); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	if err := h.Service.CreateTheme(&theme); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create theme",
		})
	}
	return c.Status(fiber.StatusCreated).JSON(theme)
}

func (h *ThemeHandler) UpdateTheme(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid theme ID",
		})
	}

	var theme domain.Theme
	if err := c.BodyParser(&theme); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	theme.ID = id
	if err := h.Service.UpdateTheme(&theme); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update theme",
		})
	}
	return c.JSON(theme)
}

func (h *ThemeHandler) DeleteTheme(c *fiber.Ctx) error {
	id := c.Params("id")
	themeID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid theme ID",
		})
	}

	if err := h.Service.DeleteTheme(themeID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete theme",
		})
	}
	return c.SendStatus(fiber.StatusNoContent)
}
