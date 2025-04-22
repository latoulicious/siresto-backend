package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/latoulicious/siresto-backend/internal/domain"
	"github.com/latoulicious/siresto-backend/internal/service"
	"github.com/latoulicious/siresto-backend/internal/utils"
)

type OrderRequest struct {
	CustomerName  string  `json:"customer_name"`
	CustomerPhone string  `json:"customer_phone"`
	TableNumber   int     `json:"table_number"`
	Status        string  `json:"status"`
	TotalAmount   float64 `json:"total_amount"`
	Notes         string  `json:"notes"`
}

type OrderDetailRequest struct {
	ProductID   string `json:"product_id"`
	VariationID string `json:"variation_id,omitempty"` // Optional variation
	Quantity    int    `json:"quantity"`
	Note        string `json:"note,omitempty"` // Optional note
}

type CreateOrderRequest struct {
	Order        OrderRequest         `json:"order"`
	OrderDetails []OrderDetailRequest `json:"order_details"`
}

type OrderHandler struct {
	OrderService *service.OrderService
}

func (h *OrderHandler) ListAllOrders(c *fiber.Ctx) error {
	orders, err := h.OrderService.ListAllOrders()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.Error("Failed to retrieve orders", fiber.StatusInternalServerError))
	}
	return c.Status(fiber.StatusOK).JSON(utils.Success("Orders retrieved successfully", orders))
}

func (handler *OrderHandler) CreateOrder(c *fiber.Ctx) error {
	// Parse the request using our custom request structs
	var request CreateOrderRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Map request to domain models
	order := mapOrderRequestToDomain(request.Order)
	orderDetails := mapOrderDetailsRequestToDomain(request.OrderDetails)

	// Call the service to create the order with details
	createdOrder, err := handler.OrderService.CreateOrder(order, orderDetails)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Respond with the created order, including details
	return c.Status(fiber.StatusCreated).JSON(createdOrder)
}

// Map request model to domain model
func mapOrderRequestToDomain(req OrderRequest) *domain.Order {

	return &domain.Order{
		CustomerName:  req.CustomerName,
		CustomerPhone: req.CustomerPhone,
		TableNumber:   req.TableNumber,
		Status:        domain.OrderStatus(req.Status),
		TotalAmount:   req.TotalAmount,
		Notes:         req.Notes,
	}
}

// Map request models to domain models
func mapOrderDetailsRequestToDomain(reqs []OrderDetailRequest) []domain.OrderDetail {
	details := make([]domain.OrderDetail, len(reqs))

	for i, req := range reqs {
		productID, _ := uuid.Parse(req.ProductID)     // Handle error in production
		variationID, _ := uuid.Parse(req.VariationID) // Handle error in production

		details[i] = domain.OrderDetail{
			ProductID:   &productID,
			VariationID: &variationID,
			Quantity:    req.Quantity,
			Note:        req.Note,
		}
	}

	return details
}
