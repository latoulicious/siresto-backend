package service

import (
	"errors"
	"fmt"

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
