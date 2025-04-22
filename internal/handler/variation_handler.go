package handler

import (
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/latoulicious/siresto-backend/internal/domain"
	"github.com/latoulicious/siresto-backend/internal/service"
	"github.com/latoulicious/siresto-backend/internal/utils"
	"github.com/latoulicious/siresto-backend/pkg/db"
)

type VariationHandler struct {
	Service        *service.VariationService
	ProductService *service.ProductService
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

// GetProductVariations lists all variations for a specific product
func (h *VariationHandler) GetProductVariations(c *fiber.Ctx) error {
	productID, err := uuid.Parse(c.Params("product_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("Invalid product ID format", fiber.StatusBadRequest))
	}

	variations, err := h.Service.GetVariationsByProductID(productID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.Error(err.Error(), fiber.StatusInternalServerError))
	}

	return c.Status(fiber.StatusOK).JSON(utils.Success("Product variations retrieved successfully", variations))
}

// CreateProductVariation creates a new variation for a specific product
func (h *VariationHandler) CreateProductVariation(c *fiber.Ctx) error {
	// Parse product ID from URL parameter
	productID, err := uuid.Parse(c.Params("product_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("Invalid product ID format", fiber.StatusBadRequest))
	}

	// Check if product exists before creating a variation
	product, err := h.ProductService.GetProductByID(productID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(utils.Error("Product not found", fiber.StatusNotFound))
	}

	// Optional - verify product is in a valid state to accept variations
	if !product.IsAvailable {
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("Cannot add variations to unavailable product", fiber.StatusBadRequest))
	}
	// This would require a productService dependency in your VariationHandler

	var request struct {
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
		ProductID:     productID, // Use the product ID from the URL
		VariationType: request.VariationType,
		IsDefault:     request.IsDefault,
		IsAvailable:   request.IsAvailable,
		IsRequired:    request.IsRequired,
		Options:       request.Options, // Automatically marshaled into jsonb via custom Value()
	}

	if variation.ProductID == uuid.Nil {
		log.Printf("Warning: ProductID was set to nil UUID after object creation")
		variation.ProductID = productID // Force it again
	}

	created, err := h.Service.CreateProductVariation(productID, variation)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.Status(fiber.StatusNotFound).JSON(utils.Error(err.Error(), fiber.StatusNotFound))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(utils.Error(err.Error(), fiber.StatusInternalServerError))
	}
	return c.Status(fiber.StatusCreated).JSON(utils.Success("Variation created successfully for product", created))
}
