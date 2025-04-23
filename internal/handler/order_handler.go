package handler

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/latoulicious/siresto-backend/internal/domain"
	"github.com/latoulicious/siresto-backend/internal/service"
	"github.com/latoulicious/siresto-backend/internal/utils"
	"github.com/latoulicious/siresto-backend/pkg/dto"
	"gorm.io/gorm"
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
	VariationID string `json:"variation_id,omitempty"`
	Quantity    int    `json:"quantity"`
	Note        string `json:"note,omitempty"`
}

type CreateOrderRequest struct {
	Order        OrderRequest         `json:"order"`
	OrderDetails []OrderDetailRequest `json:"order_details"`
	Payments     []PaymentRequest     `json:"payments"`
}

type PaymentRequest struct {
	Method         string  `json:"method" validate:"required,oneof=Tunai Qris Debit Kredit"`
	Amount         float64 `json:"amount" validate:"required,gt=0"`
	TransactionRef string  `json:"transaction_ref,omitempty"`
}

type OrderHandler struct {
	OrderService   *service.OrderService
	PaymentService *service.PaymentService
}

func (h *OrderHandler) ListAllOrders(c *fiber.Ctx) error {
	orders, err := h.OrderService.ListAllOrders()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.Error("Failed to retrieve orders", fiber.StatusInternalServerError))
	}

	// Map all orders to DTOs
	orderDTOs := make([]dto.OrderResponseDTO, len(orders))
	for i, order := range orders {
		orderDTOs[i] = dto.MapToOrderResponseDTO(&order)
	}

	return c.Status(fiber.StatusOK).JSON(utils.Success("Orders retrieved successfully", orderDTOs))
}

func (h *OrderHandler) GetOrderByID(c *fiber.Ctx) error {
	orderID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("Invalid order ID", fiber.StatusBadRequest))
	}

	order, err := h.OrderService.GetOrderByID(orderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(utils.Error("Order not found", fiber.StatusNotFound))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(utils.Error("Failed to retrieve order", fiber.StatusInternalServerError))
	}

	orderDTO := dto.MapToOrderResponseDTO(order)
	return c.Status(fiber.StatusOK).JSON(utils.Success("Order retrieved successfully", orderDTO))
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

	// Map to DTO before returning response
	responseDTO := dto.MapToOrderResponseDTO(createdOrder)

	// Add payment method information to response
	paymentMethods := make([]map[string]string, 0)
	if len(request.Payments) > 0 {
		for _, payment := range request.Payments {
			paymentMethods = append(paymentMethods, map[string]string{"method": payment.Method})
		}
	}

	// Ensure to add paymentMethods to your responseDTO (assuming responseDTO has a field to hold payments)
	if len(paymentMethods) > 0 {
		responseDTO.Methods = paymentMethods
	}

	// Respond with the DTO instead of raw domain model
	return c.Status(fiber.StatusCreated).JSON(responseDTO)
}

// Map request model to domain model
func mapOrderRequestToDomain(req OrderRequest) *domain.Order {

	return &domain.Order{
		CustomerName:  req.CustomerName,
		CustomerPhone: req.CustomerPhone,
		TableNumber:   req.TableNumber,
		Status:        domain.OrderStatus(req.Status),
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

func (h *OrderHandler) MarkOrderAsCompleted(c *fiber.Ctx) error {
	// Parse order ID
	orderID, err := uuid.Parse(c.Params("orderID"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("Invalid order ID", fiber.StatusBadRequest))
	}

	// Update the status
	if err := h.OrderService.UpdateDishStatusToCompleted(orderID); err != nil {
		// Handle different error types
		if strings.Contains(err.Error(), "order not found") {
			return c.Status(fiber.StatusNotFound).JSON(utils.Error(err.Error(), fiber.StatusNotFound))
		}
		if strings.Contains(err.Error(), "must be 'Diproses'") {
			return c.Status(fiber.StatusBadRequest).JSON(utils.Error(err.Error(), fiber.StatusBadRequest))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(utils.Error(err.Error(), fiber.StatusInternalServerError))
	}

	return c.Status(fiber.StatusOK).JSON(utils.Success("Order marked as completed", nil))
}

// ProcessPayment handles payment for an order and updates its status
func (h *OrderHandler) ProcessPayment(c *fiber.Ctx) error {
	// Parse order ID from route parameter
	orderID, err := uuid.Parse(c.Params("orderID"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("Invalid order ID", fiber.StatusBadRequest))
	}

	// Parse payment request
	var paymentReq PaymentRequest
	if err := c.BodyParser(&paymentReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("Invalid payment data", fiber.StatusBadRequest))
	}

	// Validate payment request
	if paymentReq.Method == "" || paymentReq.Amount <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("Invalid payment details", fiber.StatusBadRequest))
	}

	// Map to domain payment
	payment := &domain.Payment{
		Method:         domain.PaymentType(paymentReq.Method),
		Amount:         paymentReq.Amount,
		Status:         domain.PaymentStatusPending,
		TransactionRef: paymentReq.TransactionRef,
	}

	// Process payment
	processedPayment, err := h.PaymentService.ProcessOrderPayment(orderID, payment)
	if err != nil {
		// Handle different error types with appropriate status codes
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(utils.Error("Order not found", fiber.StatusNotFound))
		}

		// Handle validation errors
		if errors.Is(err, service.ErrOrderAlreadyPaid) || errors.Is(err, service.ErrOrderCancelled) {
			return c.Status(fiber.StatusBadRequest).JSON(utils.Error(err.Error(), fiber.StatusBadRequest))
		}

		if errors.Is(err, service.ErrPaymentAmountMismatch) {
			return c.Status(fiber.StatusBadRequest).JSON(utils.Error(err.Error(), fiber.StatusBadRequest))
		}

		return c.Status(fiber.StatusInternalServerError).JSON(utils.Error(err.Error(), fiber.StatusInternalServerError))
	}

	// Get the updated order with payment details
	order, err := h.OrderService.GetOrderByID(orderID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.Error("Payment processed but failed to fetch updated order", fiber.StatusInternalServerError))
	}

	// Map to DTO
	orderDTO := dto.MapToOrderResponseDTO(order)

	return c.Status(fiber.StatusOK).JSON(utils.Success("Payment processed successfully", map[string]interface{}{
		"order":   orderDTO,
		"payment": processedPayment,
	}))
}
