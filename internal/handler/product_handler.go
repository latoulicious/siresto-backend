package handler

import (
	"encoding/json"
	"path"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/latoulicious/siresto-backend/internal/service"
	"github.com/latoulicious/siresto-backend/internal/utils"
	"github.com/latoulicious/siresto-backend/pkg/dto"
)

type ProductHandler struct {
	Service *service.ProductService
}

// ListAllProductsHandler retrieves all products with pagination
func (h *ProductHandler) ListAllProducts(c *fiber.Ctx) error {
	// Get page from query parameter, default to 1 if not provided
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		errInfo := utils.NewErrorInfo("INVALID_PAGE", "Page number must be a positive integer", "page", nil)
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("Invalid page number", fiber.StatusBadRequest, errInfo))
	}

	const itemsPerPage = 10
	offset := (page - 1) * itemsPerPage

	products, total, err := h.Service.ListAllProducts(offset, itemsPerPage)
	if err != nil {
		errInfo := utils.NewErrorInfo("PRODUCT_LIST_ERROR", err.Error(), "", nil)
		return c.Status(fiber.StatusInternalServerError).JSON(utils.Error("Failed to retrieve products", fiber.StatusInternalServerError, errInfo))
	}

	var responses []dto.ProductResponse
	for _, product := range products {
		responses = append(responses, *dto.ToProductResponse(&product))
	}

	metadata := utils.NewPaginationMetadata(page, itemsPerPage, total)
	return c.Status(fiber.StatusOK).JSON(utils.Success("Products retrieved successfully", responses, metadata))
}

// GetProductByIDHandler retrieves a product by ID
func (h *ProductHandler) GetProductByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		errInfo := utils.NewErrorInfo("INVALID_ID", "The provided ID is not a valid UUID", "id", nil)
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("Invalid product ID", fiber.StatusBadRequest, errInfo))
	}

	product, err := h.Service.GetProductByID(id)
	if err != nil {
		errInfo := utils.NewErrorInfo("PRODUCT_NOT_FOUND", err.Error(), "id", nil)
		return c.Status(fiber.StatusNotFound).JSON(utils.Error("Product not found", fiber.StatusNotFound, errInfo))
	}

	response := dto.ToProductResponse(product)
	return c.Status(fiber.StatusOK).JSON(utils.Success("Product found", response))
}

// CreateProductHandler creates a new product
func (h *ProductHandler) CreateProduct(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		errInfo := utils.NewErrorInfo("INVALID_FORM", "Failed to parse multipart form", "", nil)
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("Invalid form data", fiber.StatusBadRequest, errInfo))
	}

	productData := form.Value["product"]
	if len(productData) == 0 {
		errInfo := utils.NewErrorInfo("MISSING_DATA", "Product data is required in the form", "product", nil)
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("Product data required", fiber.StatusBadRequest, errInfo))
	}

	var body dto.CreateProductRequest
	if err := json.Unmarshal([]byte(productData[0]), &body); err != nil {
		errInfo := utils.NewErrorInfo("INVALID_JSON", err.Error(), "product", nil)
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("Invalid product data", fiber.StatusBadRequest, errInfo))
	}

	product := dto.ToProductDomainFromCreate(&body)

	if err := service.ValidateProduct(h.Service.Repo.DB, product); err != nil {
		errInfo := utils.NewErrorInfo("VALIDATION_ERROR", err.Error(), "", nil)
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("Validation failed", fiber.StatusBadRequest, errInfo))
	}

	if len(form.File["image"]) > 0 {
		if h.Service.Uploader == nil {
			errInfo := utils.NewErrorInfo("UPLOADER_ERROR", "Image uploading service is not configured", "", nil)
			return c.Status(fiber.StatusInternalServerError).JSON(utils.Error("Image uploading not configured properly", fiber.StatusInternalServerError, errInfo))
		}

		fileHeader := form.File["image"][0]
		file, err := fileHeader.Open()
		if err != nil {
			errInfo := utils.NewErrorInfo("FILE_ERROR", "Failed to process uploaded image", "image", nil)
			return c.Status(fiber.StatusInternalServerError).JSON(utils.Error("Failed to open image", fiber.StatusInternalServerError, errInfo))
		}
		defer file.Close()

		filename := uuid.New().String() + path.Ext(fileHeader.Filename)

		if h.Service.Uploader == nil {
			errInfo := utils.NewErrorInfo("UPLOADER_ERROR", "Image uploading service is not available", "", nil)
			return c.Status(fiber.StatusInternalServerError).JSON(utils.Error("Image uploading service is not available", fiber.StatusInternalServerError, errInfo))
		}

		imageURL, err := h.Service.Uploader.Upload(file, "products/"+filename)
		if err != nil {
			errInfo := utils.NewErrorInfo("UPLOAD_ERROR", err.Error(), "image", nil)
			return c.Status(fiber.StatusInternalServerError).JSON(utils.Error("Failed to upload image", fiber.StatusInternalServerError, errInfo))
		}

		product.ImageURL = imageURL
	}

	createdProduct, variations, err := h.Service.CreateProductWithVariations(product, body.Variations)
	if err != nil {
		errInfo := utils.NewErrorInfo("CREATION_ERROR", err.Error(), "", nil)
		return c.Status(fiber.StatusInternalServerError).JSON(utils.Error("Failed to create product", fiber.StatusInternalServerError, errInfo))
	}

	response := dto.ToProductResponse(createdProduct)
	response.Variations = dto.ToVariationResponses(variations)

	metadata := &utils.Metadata{
		CustomData: map[string]interface{}{
			"variation_count": len(variations),
		},
	}

	return c.Status(fiber.StatusCreated).JSON(utils.Success("Product created successfully", response, metadata))
}

// UpdateProductHandler updates an existing product
func (h *ProductHandler) UpdateProduct(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		errInfo := utils.NewErrorInfo("INVALID_ID", "The provided ID is not a valid UUID", "id", nil)
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("Invalid product ID", fiber.StatusBadRequest, errInfo))
	}

	var body dto.UpdateProductRequest
	if err := c.BodyParser(&body); err != nil {
		errInfo := utils.NewErrorInfo("INVALID_BODY", "Failed to parse request body", "", nil)
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("Invalid request body", fiber.StatusBadRequest, errInfo))
	}

	existingProduct, err := h.Service.GetProductByID(id)
	if err != nil {
		errInfo := utils.NewErrorInfo("PRODUCT_NOT_FOUND", err.Error(), "id", nil)
		return c.Status(fiber.StatusNotFound).JSON(utils.Error("Product not found", fiber.StatusNotFound, errInfo))
	}

	updatedProduct := dto.ToProductDomainFromUpdate(&body, existingProduct)

	if err := service.ValidateProductForUpdate(h.Service.Repo.DB, updatedProduct); err != nil {
		errInfo := utils.NewErrorInfo("VALIDATION_ERROR", err.Error(), "", nil)
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("Validation failed", fiber.StatusBadRequest, errInfo))
	}

	removeOthers := false
	if body.RemoveOtherVariations != nil {
		removeOthers = *body.RemoveOtherVariations
	}

	updatedProduct, variations, err := h.Service.UpdateProductWithVariations(id, updatedProduct, body.Variations, removeOthers)
	if err != nil {
		errInfo := utils.NewErrorInfo("UPDATE_ERROR", err.Error(), "", nil)
		return c.Status(fiber.StatusInternalServerError).JSON(utils.Error("Failed to update product", fiber.StatusInternalServerError, errInfo))
	}

	response := dto.ToProductResponse(updatedProduct)
	response.Variations = dto.ToVariationResponses(variations)

	metadata := &utils.Metadata{
		CustomData: map[string]interface{}{
			"variation_count": len(variations),
			"removed_others":  removeOthers,
		},
	}

	return c.Status(fiber.StatusOK).JSON(utils.Success("Product updated successfully", response, metadata))
}

// DeleteProductHandler removes a product by ID
func (h *ProductHandler) DeleteProduct(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		errInfo := utils.NewErrorInfo("INVALID_ID", "The provided ID is not a valid UUID", "id", nil)
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("Invalid product ID", fiber.StatusBadRequest, errInfo))
	}

	if err := service.ValidateProductDeletable(h.Service.Repo.DB, id); err != nil {
		errInfo := utils.NewErrorInfo("VALIDATION_ERROR", err.Error(), "", nil)
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("Product cannot be deleted", fiber.StatusBadRequest, errInfo))
	}

	if err := h.Service.DeleteProduct(id); err != nil {
		errInfo := utils.NewErrorInfo("DELETE_ERROR", err.Error(), "", nil)
		return c.Status(fiber.StatusInternalServerError).JSON(utils.Error("Failed to delete product", fiber.StatusInternalServerError, errInfo))
	}

	return c.Status(fiber.StatusOK).JSON(utils.Success("Product deleted successfully", nil))
}
