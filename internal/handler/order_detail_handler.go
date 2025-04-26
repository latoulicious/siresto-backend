package handler

import (
	"github.com/latoulicious/siresto-backend/internal/domain"
	"github.com/latoulicious/siresto-backend/internal/service"

	"github.com/gofiber/fiber/v2"
)

type OrderDetailHandler struct {
	OrderDetailService *service.OrderDetailService
}

func (h *OrderDetailHandler) CreateOrderDetails(c *fiber.Ctx) error {
	var details []domain.OrderDetail
	if err := c.BodyParser(&details); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	err := h.OrderDetailService.Create(details)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Order details created successfully"})
}
