package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/latoulicious/siresto-backend/internal/domain"
	"github.com/latoulicious/siresto-backend/internal/service"
	"github.com/latoulicious/siresto-backend/internal/utils"
	"github.com/latoulicious/siresto-backend/pkg/db"
)

type VariationHandler struct {
	Service *service.VariationService
}

//! Global Variation Handler

// ListVariationsHandler lists all variations for a product
func (h *VariationHandler) ListAllVariations(c *fiber.Ctx) error {
	variations, err := h.Service.ListAllVariations()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.Error(err.Error(), fiber.StatusInternalServerError))
	}

	return c.Status(fiber.StatusOK).JSON(utils.Success("Product variations retrieved successfully", variations))
}

// GetVariationByIDHandler retrieves a variation by its ID
func (h *VariationHandler) GetVariationByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("Invalid variation ID format", fiber.StatusBadRequest))
	}

	variation, err := h.Service.GetVariationByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(utils.Error("Variation ID not found", fiber.StatusNotFound))
	}

	return c.Status(fiber.StatusOK).JSON(utils.Success("Variation retrieved successfully", variation))
}

// CreateVariationHandler creates a new variation for a product
func (h *VariationHandler) CreateVariation(c *fiber.Ctx) error {
	var request struct {
		ProductID     uuid.UUID            `json:"product_id"`
		VariationType string               `json:"variation_type"`
		IsDefault     bool                 `json:"is_default"`
		IsAvailable   bool                 `json:"is_available"`
		IsRequired    bool                 `json:"is_required"`
		Options       []db.VariationOption `json:"options"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("Invalid request format", fiber.StatusBadRequest))
	}

	// Basic validation
	if request.VariationType == "" {
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("variation_type is required", fiber.StatusBadRequest))
	}
	if len(request.Options) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("options must contain at least one item", fiber.StatusBadRequest))
	}

	variation := &domain.Variation{
		ProductID:     request.ProductID,
		VariationType: request.VariationType,
		IsDefault:     request.IsDefault,
		IsAvailable:   request.IsAvailable,
		IsRequired:    request.IsRequired,
		Options:       request.Options, // Automatically marshaled into jsonb via custom Value()
	}

	created, err := h.Service.CreateVariation(variation)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.Error(err.Error(), fiber.StatusInternalServerError))
	}

	return c.Status(fiber.StatusCreated).JSON(utils.Success("Variation created successfully", created))
}

// UpdateVariationHandler updates an existing variation
func (h *VariationHandler) UpdateVariation(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("Invalid variation ID format", fiber.StatusBadRequest))
	}

	var request struct {
		Description   string   `json:"description"`
		PriceModifier *float64 `json:"price_modifier"`
		PriceAbsolute *float64 `json:"price_absolute"`
		IsDefault     bool     `json:"is_default"`
		IsAvailable   bool     `json:"is_available"`
		IsRequired    bool     `json:"is_required"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("Invalid request format", fiber.StatusBadRequest))
	}

	variation := &domain.Variation{
		ID:          id,
		IsDefault:   request.IsDefault,
		IsAvailable: request.IsAvailable,
		IsRequired:  request.IsRequired,
	}

	updatedVariation, err := h.Service.UpdateVariation(variation)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.Error(err.Error(), fiber.StatusInternalServerError))
	}

	return c.Status(fiber.StatusOK).JSON(utils.Success("Variation updated successfully", updatedVariation))
}

// DeleteVariationHandler deletes a variation by its ID
func (h *VariationHandler) DeleteVariation(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("Invalid variation ID format", fiber.StatusBadRequest))
	}

	err = h.Service.DeleteVariation(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.Error(err.Error(), fiber.StatusInternalServerError))
	}

	return c.Status(fiber.StatusNoContent).JSON(utils.Success("Variation deleted successfully", nil))
}

// TODO Implement Function for Product Tied Variations
