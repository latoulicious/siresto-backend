package handler

import (
	"log"
	"strings"

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
	// Log request info
	log.Printf("Login attempt from IP: %s", c.IP())

	var req dto.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("Login request parsing error: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("Invalid request body", fiber.StatusBadRequest))
	}

	// Log email for debugging (mask part of it for privacy)
	if len(req.Email) > 3 {
		atIndex := strings.Index(req.Email, "@")
		if atIndex > 0 {
			maskedEmail := req.Email[:3] + "***" + req.Email[atIndex:]
			log.Printf("Login attempt for email: %s", maskedEmail)
		} else {
			log.Printf("Login attempt with malformed email (no @ symbol)")
		}
	}

	if err := h.Validate.Struct(req); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		log.Printf("Login validation error: %v", validationErrors)
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error(validationErrors.Error(), fiber.StatusBadRequest))
	}

	// Log that we're calling service layer
	log.Printf("Attempting login via service for email: %s", req.Email)

	loginResponse, err := h.Service.LoginUser(&req)
	if err != nil {
		log.Printf("Login failed: %v", err)
		return c.Status(fiber.StatusUnauthorized).JSON(utils.Error("Invalid email or password", fiber.StatusUnauthorized))
	}

	log.Printf("Login successful for user ID: %s", loginResponse.User.ID)
	return c.Status(fiber.StatusOK).JSON(utils.Success("Login successful", loginResponse))
}
