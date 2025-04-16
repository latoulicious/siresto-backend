package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/latoulicious/siresto-backend/internal/domain"
	"github.com/latoulicious/siresto-backend/internal/service"
	"github.com/latoulicious/siresto-backend/internal/utils"
	"github.com/latoulicious/siresto-backend/internal/validator"
)

type UserHandler struct {
	Service *service.UserService
}

// GetAllUsersHandler retrieves all users
func (h *UserHandler) ListAllUsers(c *fiber.Ctx) error {
	users, err := h.Service.ListAllUsers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.Error("Failed to retrieve users", fiber.StatusInternalServerError))
	}
	return c.Status(fiber.StatusOK).JSON(utils.Success("Users retrieved successfully", users))
}

// GetUserByIDHandler retrieves a user by ID
func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("Invalid user ID", fiber.StatusBadRequest))
	}

	user, err := h.Service.GetUserByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(utils.Error("User not found", fiber.StatusNotFound))
	}

	return c.Status(fiber.StatusOK).JSON(utils.Success("User found", user))
}

// CreateUserHandler creates a new user
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var req validator.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("Invalid request body", fiber.StatusBadRequest))
	}

	// if err := validator.Validate(req); err != nil {
	// 	return c.Status(fiber.StatusBadRequest).JSON(utils.Error(err.Error(), fiber.StatusBadRequest))
	// }

	createdUser, err := h.Service.CreateUser(&req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.Error("Failed to create user", fiber.StatusInternalServerError))
	}

	return c.Status(fiber.StatusCreated).JSON(utils.Success("User created successfully", createdUser))
}

// UpdateUserHandler updates an existing user
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("Invalid user ID", fiber.StatusBadRequest))
	}

	var body domain.User
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("Invalid request body", fiber.StatusBadRequest))
	}

	updatedUser, err := h.Service.UpdateUser(id, &body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.Error("Failed to update user", fiber.StatusInternalServerError))
	}

	return c.Status(fiber.StatusOK).JSON(utils.Success("User updated successfully", updatedUser))
}

// DeleteUserHandler deletes a user by ID
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("Invalid user ID", fiber.StatusBadRequest))
	}

	err = h.Service.DeleteUser(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.Error("Failed to delete user", fiber.StatusInternalServerError))
	}

	return c.Status(fiber.StatusOK).JSON(utils.Success("User deleted successfully", nil))
}
