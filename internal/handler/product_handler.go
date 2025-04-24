package handler

import (
	"encoding/json"
	"path"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/latoulicious/siresto-backend/internal/service"
	"github.com/latoulicious/siresto-backend/internal/utils"
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
	if err := service.ValidateProduct(h.Service.Repo.DB, product); err != nil {
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

	// Retrieve the existing product from the database
	existingProduct, err := h.Service.GetProductByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(utils.Error("Product not found", fiber.StatusNotFound))
	}

	// Map the request body to the domain model
	updatedProduct := dto.ToProductDomainFromUpdate(&body, existingProduct)

	// Validate the updated product
	if err := service.ValidateProductForUpdate(h.Service.Repo.DB, updatedProduct); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("Validation failed: "+err.Error(), fiber.StatusBadRequest))
	}

	// Default removeOthers to false if not specified
	removeOthers := false
	if body.RemoveOtherVariations != nil {
		removeOthers = *body.RemoveOtherVariations
	}

	// Perform the update operation with variations in the service layer
	updatedProduct, variations, err := h.Service.UpdateProductWithVariations(id, updatedProduct, body.Variations, removeOthers)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.Error("Failed to update product: "+err.Error(), fiber.StatusInternalServerError))
	}

	// Map the updated product to a response DTO
	response := dto.ToProductResponse(updatedProduct)
	// Include variations in the response
	response.Variations = dto.ToVariationResponses(variations)

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
	if err := service.ValidateProductDeletable(h.Service.Repo.DB, id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("Validation failed: "+err.Error(), fiber.StatusBadRequest))
	}

	if err := h.Service.DeleteProduct(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.Error("Failed to delete product", fiber.StatusInternalServerError))
	}

	return c.Status(fiber.StatusOK).JSON(utils.Success("Product deleted successfully", nil))
}

// Helper function to upload a product image
func (h *ProductHandler) UploadProductImage(c *fiber.Ctx) error {
	// Check if uploader is available
	if h.Service.Uploader == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.Error("image uploading not configured", fiber.StatusInternalServerError))
	}

	fileHeader, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("image file required", fiber.StatusBadRequest))
	}

	file, err := fileHeader.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.Error("failed to open image", fiber.StatusInternalServerError))
	}
	defer file.Close()

	// Generate a unique file name, e.g., UUID + extension
	filename := uuid.New().String() + path.Ext(fileHeader.Filename)

	url, err := h.Service.Uploader.Upload(file, "products/"+filename)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.Error("failed to upload image: "+err.Error(), fiber.StatusInternalServerError))
	}

	return c.Status(fiber.StatusOK).JSON(utils.Success("upload success", map[string]string{"image_url": url}))
}

// CreateProductWithImage creates a new product with image upload in a single request
func (h *ProductHandler) CreateProductWithImage(c *fiber.Ctx) error {
	// Get the multipart form
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("Invalid form data", fiber.StatusBadRequest))
	}

	// Get the product data from form field
	productData := form.Value["product"]
	if len(productData) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("Product data required", fiber.StatusBadRequest))
	}

	// Parse the product data JSON
	var body dto.CreateProductRequest
	if err := json.Unmarshal([]byte(productData[0]), &body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("Invalid product data: "+err.Error(), fiber.StatusBadRequest))
	}

	// Map DTO to domain model
	product := dto.ToProductDomainFromCreate(&body)

	// Validate the domain model
	if err := service.ValidateProduct(h.Service.Repo.DB, product); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("Validation failed: "+err.Error(), fiber.StatusBadRequest))
	}

	// Check if there's an image to upload
	if len(form.File["image"]) > 0 {
		// Check if uploader is configured - EARLY CHECK to prevent nil pointer
		if h.Service.Uploader == nil {
			return c.Status(fiber.StatusInternalServerError).JSON(utils.Error("Image uploading not configured properly. R2 configuration is missing.", fiber.StatusInternalServerError))
		}

		fileHeader := form.File["image"][0]
		file, err := fileHeader.Open()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(utils.Error("Failed to open image", fiber.StatusInternalServerError))
		}
		defer file.Close()

		// Generate a unique filename
		filename := uuid.New().String() + path.Ext(fileHeader.Filename)

		// Extra check for nil uploader right before using it
		if h.Service.Uploader == nil {
			return c.Status(fiber.StatusInternalServerError).JSON(utils.Error("Image uploading service is not available", fiber.StatusInternalServerError))
		}

		// Upload the image
		imageURL, err := h.Service.Uploader.Upload(file, "products/"+filename)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(utils.Error("Failed to upload image: "+err.Error(), fiber.StatusInternalServerError))
		}

		// Set the image URL on the product
		product.ImageURL = imageURL
	}

	// Create the product with variations
	createdProduct, variations, err := h.Service.CreateProductWithVariations(product, body.Variations)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.Error("Failed to create product: "+err.Error(), fiber.StatusInternalServerError))
	}

	// Map the domain model back to a DTO for the response
	response := dto.ToProductResponse(createdProduct)
	// Include variations in the response
	response.Variations = dto.ToVariationResponses(variations)

	// Return the successful response
	return c.Status(fiber.StatusCreated).JSON(utils.Success("Product created successfully", response))
}
