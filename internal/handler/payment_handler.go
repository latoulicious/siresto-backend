package handler

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/latoulicious/siresto-backend/internal/domain"
	"github.com/latoulicious/siresto-backend/internal/service"
	"github.com/latoulicious/siresto-backend/internal/utils"
	"gorm.io/gorm"
)

type PaymentHandler struct {
	Service *service.PaymentService
}

func (h *PaymentHandler) ListAllOrderPayments(c *fiber.Ctx) error {
	payments, err := h.Service.ListAllOrderPayments()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.Error("Failed to retrieve order payments", fiber.StatusInternalServerError))
	}

	// If no payments found, return empty array instead of error
	if len(payments) == 0 {
		return c.Status(fiber.StatusOK).JSON(utils.Success("No order payments found", []domain.Payment{}))
	}

	return c.Status(fiber.StatusOK).JSON(utils.Success("Order payments retrieved successfully", payments))
}

func (h *PaymentHandler) CreatePayment(c *fiber.Ctx) error {
	var payment domain.Payment
	if err := c.BodyParser(&payment); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("Invalid request", fiber.StatusBadRequest))
	}

	if err := h.Service.CreatePayment(&payment); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.Error("Failed to create payment", fiber.StatusInternalServerError))
	}

	return c.Status(fiber.StatusCreated).JSON(utils.Success("Payment created", payment))
}

func (h *PaymentHandler) ProcessOrderPayment(c *fiber.Ctx) error {
	// Parse order ID from route parameter
	orderID, err := uuid.Parse(c.Params("orderID"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("Invalid order ID", fiber.StatusBadRequest))
	}

	// Parse payment request
	var payment domain.Payment
	if err := c.BodyParser(&payment); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("Invalid payment data", fiber.StatusBadRequest))
	}

	// Process payment
	processedPayment, err := h.Service.ProcessOrderPayment(orderID, &payment)
	if err != nil {
		// Handle different error types with appropriate status codes
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(utils.Error("Order not found", fiber.StatusNotFound))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(utils.Error(err.Error(), fiber.StatusInternalServerError))
	}

	return c.Status(fiber.StatusOK).JSON(utils.Success("Payment processed successfully", processedPayment))
}

// GetOrderPayments returns all payments for a specific order
func (h *PaymentHandler) GetOrderPayments(c *fiber.Ctx) error {
	// Parse order ID from route parameter
	orderID, err := uuid.Parse(c.Params("orderID"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("Invalid order ID", fiber.StatusBadRequest))
	}

	// Use the payment repository's method to get payments by order ID
	payments, err := h.Service.Repo.GetByOrderID(orderID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.Error("Failed to retrieve payments", fiber.StatusInternalServerError))
	}

	// If no payments found, return empty array instead of error
	if len(payments) == 0 {
		return c.Status(fiber.StatusOK).JSON(utils.Success("No payments found for this order", []domain.Payment{}))
	}

	return c.Status(fiber.StatusOK).JSON(utils.Success("Payments retrieved successfully", payments))
}
