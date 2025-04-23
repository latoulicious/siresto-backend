package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/latoulicious/siresto-backend/internal/service"
	"github.com/latoulicious/siresto-backend/internal/utils"
	"github.com/latoulicious/siresto-backend/pkg/dto"
)

type UserHandler struct {
	Service  *service.UserService
	Validate *validator.Validate
}

func NewUserHandler(service *service.UserService, validate *validator.Validate) *UserHandler {
	return &UserHandler{
		Service:  service,
		Validate: validate,
	}
}

// ListAllUsers retrieves all users
func (h *UserHandler) ListAllUsers(c *fiber.Ctx) error {
	users, err := h.Service.ListAllUsers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.Error("Failed to retrieve users", fiber.StatusInternalServerError))
	}
	return c.Status(fiber.StatusOK).JSON(utils.Success("Users retrieved successfully", users))
}

// GetUserByID retrieves a user by ID
func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("Invalid user ID", fiber.StatusBadRequest))
	}

	user, err := h.Service.GetUserByID(id)
	if err != nil {
		switch err {
		case service.ErrUserNotFound:
			return c.Status(fiber.StatusNotFound).JSON(utils.Error("User not found", fiber.StatusNotFound))
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(utils.Error("Failed to retrieve user", fiber.StatusInternalServerError))
		}
	}

	return c.Status(fiber.StatusOK).JSON(utils.Success("User found", user))
}

// CreateUser creates a new user
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var req dto.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("Invalid request body", fiber.StatusBadRequest))
	}

	if err := h.Validate.Struct(req); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error(validationErrors.Error(), fiber.StatusBadRequest))
	}

	createdUser, err := h.Service.CreateUser(&req)
	if err != nil {
		switch err {
		case service.ErrEmailAlreadyExists:
			return c.Status(fiber.StatusConflict).JSON(utils.Error("Email already in use", fiber.StatusConflict))
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(utils.Error("Failed to create user", fiber.StatusInternalServerError))
		}
	}

	return c.Status(fiber.StatusCreated).JSON(utils.Success("User created successfully", createdUser))
}

// UpdateUser updates an existing user
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("Invalid user ID", fiber.StatusBadRequest))
	}

	var req dto.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("Invalid request body", fiber.StatusBadRequest))
	}

	if err := h.Validate.Struct(req); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error(validationErrors.Error(), fiber.StatusBadRequest))
	}

	updatedUser, err := h.Service.UpdateUser(id, &req)
	if err != nil {
		switch err {
		case service.ErrUserNotFound:
			return c.Status(fiber.StatusNotFound).JSON(utils.Error("User not found", fiber.StatusNotFound))
		case service.ErrEmailAlreadyExists:
			return c.Status(fiber.StatusConflict).JSON(utils.Error("Email already in use", fiber.StatusConflict))
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(utils.Error("Failed to update user", fiber.StatusInternalServerError))
		}
	}

	return c.Status(fiber.StatusOK).JSON(utils.Success("User updated successfully", updatedUser))
}

// DeleteUser deletes a user by ID
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("Invalid user ID", fiber.StatusBadRequest))
	}

	err = h.Service.DeleteUser(id)
	if err != nil {
		switch err {
		case service.ErrUserNotFound:
			return c.Status(fiber.StatusNotFound).JSON(utils.Error("User not found", fiber.StatusNotFound))
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(utils.Error("Failed to delete user", fiber.StatusInternalServerError))
		}
	}

	return c.Status(fiber.StatusOK).JSON(utils.Success("User deleted successfully", nil))
}

// LoginUser handles user login
func (h *UserHandler) LoginUser(c *fiber.Ctx) error {
	var req dto.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("Invalid request body", fiber.StatusBadRequest))
	}

	if err := h.Validate.Struct(req); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error(validationErrors.Error(), fiber.StatusBadRequest))
	}

	loginResponse, err := h.Service.LoginUser(&req)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(utils.Error("Invalid email or password", fiber.StatusUnauthorized))
	}

	return c.Status(fiber.StatusOK).JSON(utils.Success("Login successful", loginResponse))
}
