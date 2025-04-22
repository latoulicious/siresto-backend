package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/latoulicious/siresto-backend/internal/domain"
	"github.com/latoulicious/siresto-backend/internal/service"
	"github.com/latoulicious/siresto-backend/internal/utils"
)

type PaymentHandler struct {
	Service *service.PaymentService
}

func (h *PaymentHandler) ListAllPayments(c *fiber.Ctx) error {
	payments, err := h.Service.ListAllPayments()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.Error("Failed to retrieve payments", fiber.StatusInternalServerError))
	}
	return c.Status(fiber.StatusOK).JSON(utils.Success("Payments retrieved successfully", payments))
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
