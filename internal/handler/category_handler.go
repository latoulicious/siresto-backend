package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/latoulicious/siresto-backend/internal/domain"
	"github.com/latoulicious/siresto-backend/internal/service"
	"github.com/latoulicious/siresto-backend/internal/utils"
	"github.com/latoulicious/siresto-backend/internal/validator"
)

type CategoryHandler struct {
	Service *service.CategoryService
}

// GetAllHandler retrieves all categories
func (h CategoryHandler) ListAllCategories(c *fiber.Ctx) error {
	// Parse query parameter for product inclusion
	includeProducts := c.Query("include_products") == "true"

	categories, err := h.Service.ListAllCategories(includeProducts)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.Error("Failed to retrieve categories", fiber.StatusInternalServerError))
	}
	return c.Status(fiber.StatusOK).JSON(utils.Success("Categories retrieved successfully", categories))
}

// GetByIDHandler retrieves a category by ID
func (h *CategoryHandler) GetCategoryByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("Invalid category ID", fiber.StatusBadRequest))
	}

	includeProducts := c.Query("include_products") == "true"
	category, err := h.Service.GetCategoryByID(id, includeProducts)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(utils.Error("Category not found", fiber.StatusNotFound))
	}

	return c.Status(fiber.StatusOK).JSON(utils.Success("Category found", category))
}

// CreateHandler creates a new category
func (h *CategoryHandler) CreateCategory(c *fiber.Ctx) error {
	var body domain.Category
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("Invalid request body", fiber.StatusBadRequest))
	}

	// Validate category before creation
	if err := validator.ValidateCategory(h.Service.Repo.DB, &body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error(err.Error(), fiber.StatusBadRequest))
	}

	createdCategory, err := h.Service.CreateCategory(&body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.Error("Failed to create category", fiber.StatusInternalServerError))
	}

	return c.Status(fiber.StatusCreated).JSON(utils.Success("Category created successfully", createdCategory))
}

// UpdateHandler updates an existing category
func (h *CategoryHandler) UpdateCategory(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("Invalid category ID", fiber.StatusBadRequest))
	}

	var body domain.Category
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("Invalid request body", fiber.StatusBadRequest))
	}

	// Validate category before update
	if err := validator.ValidateCategory(h.Service.Repo.DB, &body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error(err.Error(), fiber.StatusBadRequest))
	}

	updatedCategory, err := h.Service.UpdateCategory(id, &body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.Error("Failed to update category", fiber.StatusInternalServerError))
	}

	return c.Status(fiber.StatusOK).JSON(utils.Success("Category updated successfully", updatedCategory))
}

// DeleteHandler deletes a category by ID
func (h *CategoryHandler) DeleteCategory(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("Invalid category ID", fiber.StatusBadRequest))
	}

	// Validate category before deletion
	if err := validator.ValidateCategoryDeletable(h.Service.Repo.DB, id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error(err.Error(), fiber.StatusBadRequest))
	}

	if err := h.Service.DeleteCategory(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.Error("Failed to delete category", fiber.StatusInternalServerError))
	}

	return c.Status(fiber.StatusNoContent).JSON(utils.Success("Category deleted successfully", nil))
}
