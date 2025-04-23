package service

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/latoulicious/siresto-backend/internal/domain"
	"github.com/latoulicious/siresto-backend/internal/repository"
)

type OrderService struct {
	Repo          *repository.OrderRepository
	ProductRepo   *repository.ProductRepository
	VariationRepo *repository.VariationRepository
}

func (s *OrderService) ListAllOrders() ([]domain.Order, error) {
	// Call the repository method to fetch all orders
	return s.Repo.ListAllOrders()
}

func (s *OrderService) CreateOrder(order *domain.Order, details []domain.OrderDetail) (*domain.Order, error) {
	if order.CustomerName == "" || order.CustomerPhone == "" {
		return nil, errors.New("missing required order fields")
	}

	// Load all related product and variation data
	if err := s.loadProductData(&details); err != nil {
		return nil, err
	}

	// Calculate prices for each order detail
	totalAmount := 0.0
	for i := range details {
		// Calculate unit price based on product price and variation modifiers
		unitPrice := calculateUnitPrice(&details[i])
		details[i].UnitPrice = unitPrice

		// Calculate line total
		details[i].TotalPrice = unitPrice * float64(details[i].Quantity)

		// Add to order total
		totalAmount += details[i].TotalPrice
	}

	// Override any client-provided total with our calculated total
	order.TotalAmount = totalAmount

	createdOrder, err := s.Repo.CreateOrder(order, details)
	if err != nil {
		return nil, err
	}

	// Fetch with all associations
	return s.Repo.GetOrderWithAssociations(createdOrder.ID)
}

// Helper method to load product and variation data
func (s *OrderService) loadProductData(details *[]domain.OrderDetail) error {
	for i := range *details {
		detail := &(*details)[i]

		// Skip if IDs not provided
		if detail.ProductID == nil {
			continue
		}

		// Load product data
		product, err := s.ProductRepo.GetProductByID(*detail.ProductID)
		if err != nil {
			return fmt.Errorf("product not found: %v", err)
		}
		detail.Product = product
		detail.ProductName = product.Name

		// Load variation if specified
		if detail.VariationID != nil {
			variation, err := s.VariationRepo.GetVariationByID(*detail.VariationID)
			if err != nil {
				return fmt.Errorf("variation not found: %v", err)
			}
			detail.Variation = variation

			// Find default variation option for name
			for _, opt := range variation.Options {
				if opt.IsDefault {
					detail.VariationName = opt.Label
					break
				}
			}
		}
	}
	return nil
}

// Calculate unit price considering base price and variation modifiers
func calculateUnitPrice(detail *domain.OrderDetail) float64 {
	// Start with base product price
	price := detail.Product.BasePrice

	// Apply variation price modifiers if applicable
	if detail.Variation != nil && len(detail.Variation.Options) > 0 {
		// Find default option's price modifier
		for _, opt := range detail.Variation.Options {
			if opt.IsDefault {
				// If absolute price is set, use it instead of base price
				if opt.PriceAbsolute != nil {
					price = *opt.PriceAbsolute
				} else if opt.PriceModifier != nil {
					// Otherwise apply modifier to base price
					price += *opt.PriceModifier
				}
				break
			}
		}
	}

	return price
}

func (s *OrderService) GetOrderPaymentStatus(orderID uuid.UUID) (domain.OrderStatus, error) {
	// Get order with payment details
	order, err := s.Repo.GetOrderWithDetails(orderID)
	if err != nil {
		return "", err
	}

	// Return current order status
	return order.Status, nil
}

// Add payment validation method
func (s *OrderService) ValidatePaymentForOrder(orderID uuid.UUID, amount float64) error {
	order, err := s.Repo.GetOrderWithDetails(orderID)
	if err != nil {
		return err
	}

	if order.Status == domain.OrderStatusPaid {
		return errors.New("order already paid")
	}

	if order.Status == domain.OrderStatusCancelled {
		return errors.New("cannot pay for cancelled order")
	}

	if order.TotalAmount != amount {
		return fmt.Errorf("payment amount (%f) does not match order total (%f)",
			amount, order.TotalAmount)
	}

	return nil
}

func (s *OrderService) UpdateDishStatusToCompleted(orderID uuid.UUID) error {
	// Begin transaction
	tx := s.Repo.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Get current order
	var order domain.Order
	if err := tx.First(&order, "id = ?", orderID).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("order not found: %w", err)
	}

	// Validate the current dish status - only allow transition from Diproses to Selesai
	if order.DishStatus != domain.FoodStatusInProcess {
		tx.Rollback()
		return fmt.Errorf("dish status must be 'In Process' to mark as completed, current status: %s", order.DishStatus)
	}

	// Update to Selesai
	if err := tx.Model(&order).Update("dish_status", domain.FoodStatusCompleted).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update dish status: %w", err)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("transaction failed: %w", err)
	}

	return nil
}
