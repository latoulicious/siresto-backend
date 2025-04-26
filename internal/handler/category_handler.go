package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/latoulicious/siresto-backend/internal/domain"
	"github.com/latoulicious/siresto-backend/internal/service"
	"github.com/latoulicious/siresto-backend/internal/utils"
	"github.com/latoulicious/siresto-backend/pkg/dto"
)

type CategoryHandler struct {
	Service *service.CategoryService
}

// ListAllCategories retrieves all categories, always including products
func (h *CategoryHandler) ListAllCategories(c *fiber.Ctx) error {
	categories, err := h.Service.ListAllCategories(true)
	if err != nil {
		errInfo := utils.NewErrorInfo("CATEGORY_LIST_ERROR", err.Error(), "", nil)
		return c.Status(fiber.StatusInternalServerError).JSON(utils.Error("Failed to retrieve categories", fiber.StatusInternalServerError, errInfo))
	}

	var categoryResponses []dto.CategoryResponse
	for _, category := range categories {
		categoryResponses = append(categoryResponses, *dto.ToCategoryResponse(&category))
	}

	metadata := utils.NewPaginationMetadata(1, len(categoryResponses), len(categoryResponses))
	return c.Status(fiber.StatusOK).JSON(utils.Success("Categories retrieved successfully", categoryResponses, metadata))
}

// GetCategoryByID retrieves a category by ID
func (h *CategoryHandler) GetCategoryByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		errInfo := utils.NewErrorInfo("INVALID_ID", "The provided ID is not a valid UUID", "id", nil)
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("Invalid category ID", fiber.StatusBadRequest, errInfo))
	}

	category, err := h.Service.GetCategoryByID(id, true)
	if err != nil {
		errInfo := utils.NewErrorInfo("CATEGORY_NOT_FOUND", err.Error(), "id", nil)
		return c.Status(fiber.StatusNotFound).JSON(utils.Error("Category not found", fiber.StatusNotFound, errInfo))
	}

	response := dto.ToCategoryResponse(category)
	return c.Status(fiber.StatusOK).JSON(utils.Success("Category retrieved successfully", response))
}

// CreateHandler creates a new category
func (h *CategoryHandler) CreateCategory(c *fiber.Ctx) error {
	var body domain.Category
	if err := c.BodyParser(&body); err != nil {
		errInfo := utils.NewErrorInfo("INVALID_REQUEST", "Failed to parse request body", "", nil)
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("Invalid request body", fiber.StatusBadRequest, errInfo))
	}

	createdCategory, err := h.Service.CreateCategory(&body)
	if err != nil {
		errInfo := utils.NewErrorInfo("CATEGORY_CREATE_ERROR", err.Error(), "", nil)
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("Failed to create category", fiber.StatusBadRequest, errInfo))
	}

	response := dto.ToCategoryResponse(createdCategory)
	return c.Status(fiber.StatusCreated).JSON(utils.Success("Category created successfully", response))
}

// UpdateHandler updates an existing category
func (h *CategoryHandler) UpdateCategory(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		errInfo := utils.NewErrorInfo("INVALID_ID", "The provided ID is not a valid UUID", "id", nil)
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("Invalid category ID", fiber.StatusBadRequest, errInfo))
	}

	var body dto.UpdateCategoryRequest
	if err := c.BodyParser(&body); err != nil {
		errInfo := utils.NewErrorInfo("INVALID_REQUEST", "Failed to parse request body", "", nil)
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("Invalid request body", fiber.StatusBadRequest, errInfo))
	}

	updatedCategory, err := h.Service.UpdateCategory(id, &body)
	if err != nil {
		errInfo := utils.NewErrorInfo("CATEGORY_UPDATE_ERROR", err.Error(), "", nil)
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("Failed to update category", fiber.StatusBadRequest, errInfo))
	}

	return c.Status(fiber.StatusOK).JSON(utils.Success("Category updated successfully", updatedCategory))
}

// DeleteHandler deletes a category by ID
func (h *CategoryHandler) DeleteCategory(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		errInfo := utils.NewErrorInfo("INVALID_ID", "The provided ID is not a valid UUID", "id", nil)
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("Invalid category ID", fiber.StatusBadRequest, errInfo))
	}

	if err := h.Service.DeleteCategory(id); err != nil {
		errInfo := utils.NewErrorInfo("CATEGORY_DELETE_ERROR", err.Error(), "", nil)
		return c.Status(fiber.StatusInternalServerError).JSON(utils.Error("Failed to delete category", fiber.StatusInternalServerError, errInfo))
	}

	return c.Status(fiber.StatusNoContent).JSON(utils.Success("Category deleted successfully", nil))
}
