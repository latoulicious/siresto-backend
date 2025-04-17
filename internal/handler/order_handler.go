package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/latoulicious/siresto-backend/internal/domain"
	"github.com/latoulicious/siresto-backend/internal/service"
)

type OrderHandler struct {
	OrderService *service.OrderService
}

func (handler *OrderHandler) CreateOrder(c *fiber.Ctx) error {
	// Parse the request body into the Order and OrderDetails domain
	var request struct {
		Order        domain.Order         `json:"order"`
		OrderDetails []domain.OrderDetail `json:"order_details"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Call the service to create the order
	createdOrder, err := handler.OrderService.CreateOrder(&request.Order, request.OrderDetails)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Respond with the created order
	return c.Status(fiber.StatusCreated).JSON(createdOrder)
}
