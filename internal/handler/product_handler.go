package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/latoulicious/siresto-backend/internal/domain"
	"github.com/latoulicious/siresto-backend/internal/service"
	"github.com/latoulicious/siresto-backend/internal/utils"
	"github.com/latoulicious/siresto-backend/internal/validator"
)

type ProductHandler struct {
	Service *service.ProductService
}

// ListAllProductsHandler retrieves all products
func (h *ProductHandler) ListAllProducts(c *fiber.Ctx) error {
	products, err := h.Service.ListAllProducts()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.Error("Failed to retrieve products", fiber.StatusInternalServerError))
	}
	return c.Status(fiber.StatusOK).JSON(utils.Success("Products retrieved successfully", products))
}

// GetProductByIDHandler retrieves a product by ID
func (h *ProductHandler) GetProductByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("Invalid product ID", fiber.StatusBadRequest))
	}
	product, err := h.Service.GetProductByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(utils.Error("Product not found", fiber.StatusNotFound))
	}
	return c.Status(fiber.StatusOK).JSON(utils.Success("Product found", product))
}

// CreateProductHandler creates a new product
func (h *ProductHandler) CreateProduct(c *fiber.Ctx) error {
	var body domain.Product
	if err := validator.ValidateProduct(h.Service.Repo.DB, &body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("Validation failed: "+err.Error(), fiber.StatusBadRequest))
	}

	createdProduct, err := h.Service.CreateProduct(&body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.Error("Failed to create product: "+err.Error(), fiber.StatusInternalServerError))
	}

	return c.Status(fiber.StatusCreated).JSON(utils.Success("Product created successfully", createdProduct))
}

// UpdateProductHandler updates an existing product
func (h *ProductHandler) UpdateProduct(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("Invalid product ID", fiber.StatusBadRequest))
	}

	var body domain.Product
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("Invalid request body", fiber.StatusBadRequest))
	}

	if err := validator.ValidateProductForUpdate(h.Service.Repo.DB, &body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("Validation failed: "+err.Error(), fiber.StatusBadRequest))
	}

	updatedProduct, err := h.Service.UpdateProduct(id, &body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.Error("Failed to update product", fiber.StatusInternalServerError))
	}

	return c.Status(fiber.StatusOK).JSON(utils.Success("Product updated successfully", updatedProduct))
}

// DeleteProductHandler removes a product by ID
func (h *ProductHandler) DeleteProduct(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("Invalid product ID", fiber.StatusBadRequest))
	}

	// Validate if deletable
	if err := validator.ValidateProductDeletable(h.Service.Repo.DB, id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("Validation failed: "+err.Error(), fiber.StatusBadRequest))
	}

	if err := h.Service.DeleteProduct(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.Error("Failed to delete product", fiber.StatusInternalServerError))
	}

	return c.Status(fiber.StatusOK).JSON(utils.Success("Product deleted successfully", nil))
}
