package handler

import (
	"encoding/base64"
	"fmt"
	"io"
	"mime/multipart"

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
	// Parse form data directly (no DTO)
	var req struct {
		Name            string `json:"name" validate:"required"`
		PrimaryColor    string `json:"primaryColor"`
		SecondaryColor  string `json:"secondaryColor"`
		AccentColor     string `json:"accentColor"`
		BackgroundColor string `json:"backgroundColor"`
		IsDefault       *bool  `json:"isDefault"`
	}

	// Parse body for JSON data (form data will be in multipart)
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("Invalid request payload", fiber.StatusBadRequest))
	}

	// Parse multipart form file for logo and favicon
	logoFile, err := c.FormFile("logo")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("Logo file is required", fiber.StatusBadRequest))
	}

	var faviconBase64 string
	if favicon, err := c.FormFile("favicon"); err == nil {
		// Convert favicon file to Base64
		faviconBase64, err = h.convertFileToBase64(favicon)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(utils.Error("Failed to convert favicon", fiber.StatusInternalServerError))
		}
	}

	// Convert logo to Base64 string
	logoBase64, err := h.convertFileToBase64(logoFile)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.Error("Failed to convert logo", fiber.StatusInternalServerError))
	}

	// Set theme struct directly from parsed data
	theme := &domain.Theme{
		ID:              uuid.New(),
		Name:            req.Name,
		PrimaryColor:    req.PrimaryColor,
		SecondaryColor:  req.SecondaryColor,
		AccentColor:     req.AccentColor,
		BackgroundColor: req.BackgroundColor,
		LogoURL:         logoBase64, // Store the Base64 URL here
		FaviconURL:      faviconBase64,
		IsDefault:       true, // default to true unless explicitly false
	}

	if req.IsDefault != nil {
		theme.IsDefault = *req.IsDefault
	}

	// Call service to create the theme
	if err := h.Service.CreateTheme(theme); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.Error("Failed to create theme", fiber.StatusInternalServerError))
	}

	return c.Status(fiber.StatusCreated).JSON(utils.Success("Theme created successfully", theme))
}

func (h *ThemeHandler) UpdateTheme(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("Invalid theme ID", fiber.StatusBadRequest))
	}

	// Retrieve the existing theme from the database
	theme, err := h.Service.GetThemeByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(utils.Error("Theme not found", fiber.StatusNotFound))
	}

	// Parse the request body into a partial update model
	var body domain.Theme
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("Invalid request body", fiber.StatusBadRequest))
	}

	// Only update the fields that are set in the request body
	if body.Name != "" {
		theme.Name = body.Name
	}
	if body.PrimaryColor != "" {
		theme.PrimaryColor = body.PrimaryColor
	}
	if body.SecondaryColor != "" {
		theme.SecondaryColor = body.SecondaryColor
	}
	if body.AccentColor != "" {
		theme.AccentColor = body.AccentColor
	}
	if body.BackgroundColor != "" {
		theme.BackgroundColor = body.BackgroundColor
	}
	if body.LogoURL != "" {
		theme.LogoURL = body.LogoURL
	}
	if body.FaviconURL != "" {
		theme.FaviconURL = body.FaviconURL
	}
	if body.IsDefault != theme.IsDefault {
		theme.IsDefault = body.IsDefault
	}

	// Save the updated theme back to the database and handle both returned values
	updatedTheme, err := h.Service.UpdateTheme(theme)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.Error("Failed to update theme", fiber.StatusInternalServerError))
	}

	return c.Status(fiber.StatusOK).JSON(utils.Success("Theme updated successfully", updatedTheme))
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

// Helper Function
func (h *ThemeHandler) convertFileToBase64(file *multipart.FileHeader) (string, error) {
	// Open the file
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// Read the file content into a byte slice
	fileBytes, err := io.ReadAll(src)
	if err != nil {
		return "", err
	}

	// Encode the file content into base64
	base64String := base64.StdEncoding.EncodeToString(fileBytes)

	// Return Base64-encoded string with data URL prefix (can be used directly on front-end)
	// For example: "data:image/png;base64,iVBORw0KGgoAAAANS...etc"
	dataURL := fmt.Sprintf("data:%s;base64,%s", file.Header.Get("Content-Type"), base64String)

	return dataURL, nil
}
