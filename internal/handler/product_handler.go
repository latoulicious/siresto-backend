package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/latoulicious/siresto-backend/internal/service"
	"github.com/latoulicious/siresto-backend/internal/utils"
	"github.com/latoulicious/siresto-backend/internal/validator"
	"github.com/latoulicious/siresto-backend/pkg/dto"
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

	var responses []dto.ProductResponse
	for _, product := range products {
		responses = append(responses, *dto.ToProductResponse(&product))
	}

	return c.Status(fiber.StatusOK).JSON(utils.Success("Products retrieved successfully", responses))
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

	response := dto.ToProductResponse(product)
	return c.Status(fiber.StatusOK).JSON(utils.Success("Product found", response))
}

// CreateProductHandler creates a new product
func (h *ProductHandler) CreateProduct(c *fiber.Ctx) error {
	// Use the CreateProductRequest DTO to parse the request body
	var body dto.CreateProductRequest
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("Invalid request body", fiber.StatusBadRequest))
	}

	// Map DTO to domain model
	product := dto.ToProductDomainFromCreate(&body)

	// Validate the domain model (since we expect CategoryID and other fields to be in PascalCase inside the domain)
	if err := validator.ValidateProduct(h.Service.Repo.DB, product); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("Validation failed: "+err.Error(), fiber.StatusBadRequest))
	}

	// Create the product (interact with the service layer)
	createdProduct, variations, err := h.Service.CreateProductWithVariations(product, body.Variations)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.Error("Failed to create product: "+err.Error(), fiber.StatusInternalServerError))
	}

	// Map the domain model back to a DTO for the response (snake_case for frontend)
	response := dto.ToProductResponse(createdProduct)
	// Include variations in the response
	response.Variations = dto.ToVariationResponses(variations)

	// Return the successful response with the product details and variations
	return c.Status(fiber.StatusCreated).JSON(utils.Success("Product created successfully", response))
}

// UpdateProductHandler updates an existing product
func (h *ProductHandler) UpdateProduct(c *fiber.Ctx) error {
	// Parse the product ID from the URL parameter
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("Invalid product ID", fiber.StatusBadRequest))
	}

	// Parse the request body into the UpdateProductRequest DTO
	var body dto.UpdateProductRequest
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("Invalid request body", fiber.StatusBadRequest))
	}

	// Retrieve the existing product from the database to preserve unchanged fields
	existingProduct, err := h.Service.GetProductByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(utils.Error("Product not found", fiber.StatusNotFound))
	}

	// Map the request body to the domain model, passing the existing product to preserve unchanged fields
	updatedProduct := dto.ToProductDomainFromUpdate(&body, existingProduct)

	// Validate the updated product (we can skip validation for `nil` values to support partial updates)
	if err := validator.ValidateProductForUpdate(h.Service.Repo.DB, updatedProduct); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("Validation failed: "+err.Error(), fiber.StatusBadRequest))
	}

	// Perform the update operation in the service layer
	updatedProduct, err = h.Service.UpdateProduct(id, updatedProduct)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.Error("Failed to update product", fiber.StatusInternalServerError))
	}

	// Map the updated product to a response DTO for frontend consumption (e.g., snake_case)
	response := dto.ToProductResponse(updatedProduct)

	// Return the updated product in the response
	return c.Status(fiber.StatusOK).JSON(utils.Success("Product updated successfully", response))
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
