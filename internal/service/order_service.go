package service

import (
	"errors"
	"fmt"
	"time"

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

func (s *OrderService) GetOrderByID(orderID uuid.UUID) (*domain.Order, error) {
	return s.Repo.GetOrderWithAssociations(orderID)
}

func (s *OrderService) GetOrderPayments(orderID uuid.UUID) ([]domain.Payment, error) {
	order, err := s.Repo.GetOrderWithAssociations(orderID)
	if err != nil {
		return nil, err
	}

	return order.Payments, nil
}

func (s *OrderService) UpdateOrder(orderID uuid.UUID, orderUpdate *domain.Order, newDetails []domain.OrderDetail, newPayments []domain.Payment, deletedItemIDs []uuid.UUID) (*domain.Order, error) {
	// Start transaction
	tx := s.Repo.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Get existing order with all associations
	existingOrder, err := s.Repo.GetOrderWithAssociations(orderID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Update order fields if provided
	if orderUpdate != nil {
		// Only update allowed fields
		updates := map[string]interface{}{}
		if orderUpdate.CustomerName != "" {
			updates["customer_name"] = orderUpdate.CustomerName
		}
		if orderUpdate.CustomerPhone != "" {
			updates["customer_phone"] = orderUpdate.CustomerPhone
		}
		if orderUpdate.TableNumber != 0 {
			updates["table_number"] = orderUpdate.TableNumber
		}
		if orderUpdate.Notes != "" {
			updates["notes"] = orderUpdate.Notes
		}

		if len(updates) > 0 {
			if err := tx.Model(existingOrder).Updates(updates).Error; err != nil {
				tx.Rollback()
				return nil, fmt.Errorf("failed to update order: %w", err)
			}
		}
	}

	// Delete removed items if any
	if len(deletedItemIDs) > 0 {
		// Delete the items from the database
		if err := tx.Delete(&domain.OrderDetail{}, "order_id = ? AND id IN ?", orderID, deletedItemIDs).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to delete order items: %w", err)
		}

		// Also remove them from the existing order's details
		var remainingDetails []domain.OrderDetail
		for _, detail := range existingOrder.OrderDetails {
			isDeleted := false
			for _, deletedID := range deletedItemIDs {
				if detail.ID == deletedID {
					isDeleted = true
					break
				}
			}
			if !isDeleted {
				remainingDetails = append(remainingDetails, detail)
			}
		}
		existingOrder.OrderDetails = remainingDetails
	}

	// Update or create new order details
	if len(newDetails) > 0 {
		// Load product data for new details
		if err := s.loadProductData(&newDetails); err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to load product data: %w", err)
		}

		// Calculate prices for new details
		for i := range newDetails {
			newDetails[i].OrderID = orderID
			newDetails[i].UnitPrice = calculateUnitPrice(&newDetails[i])
			newDetails[i].TotalPrice = newDetails[i].UnitPrice * float64(newDetails[i].Quantity)
		}

		// Create new order details
		if err := tx.Create(&newDetails).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to create new order details: %w", err)
		}

		// Add new details to existing order's details
		existingOrder.OrderDetails = append(existingOrder.OrderDetails, newDetails...)
	}

	// Process new payments if any
	if len(newPayments) > 0 {
		// If there's an existing payment, update it instead of creating new ones
		if len(existingOrder.Payments) > 0 {
			// Update the first existing payment with new details
			existingPayment := existingOrder.Payments[0]
			updates := map[string]interface{}{
				"method":          newPayments[0].Method,
				"amount":          newPayments[0].Amount,
				"transaction_ref": newPayments[0].TransactionRef,
			}

			if err := tx.Model(&existingPayment).Updates(updates).Error; err != nil {
				tx.Rollback()
				return nil, fmt.Errorf("failed to update payment: %w", err)
			}

			// Update in-memory payment
			existingOrder.Payments[0] = domain.Payment{
				ID:             existingPayment.ID,
				OrderID:        orderID,
				Method:         newPayments[0].Method,
				Amount:         newPayments[0].Amount,
				Status:         domain.PaymentStatusSuccess,
				TransactionRef: newPayments[0].TransactionRef,
				PaidAt:         existingPayment.PaidAt,
			}
		} else {
			// Create new payment if none exists
			for i := range newPayments {
				newPayments[i].OrderID = orderID
				newPayments[i].Status = domain.PaymentStatusSuccess
				if newPayments[i].PaidAt.IsZero() {
					newPayments[i].PaidAt = time.Now()
				}
			}

			if err := tx.Create(&newPayments).Error; err != nil {
				tx.Rollback()
				return nil, fmt.Errorf("failed to create new payments: %w", err)
			}

			existingOrder.Payments = newPayments
		}
	}

	// Recalculate total amount
	var totalAmount float64
	query := tx.Model(&domain.OrderDetail{}).Where("order_id = ?", orderID)
	if len(deletedItemIDs) > 0 {
		query = query.Where("id NOT IN ?", deletedItemIDs)
	}
	if err := query.Select("COALESCE(SUM(total_price), 0)").Scan(&totalAmount).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to calculate total amount: %w", err)
	}

	// Update order total and in-memory value
	if err := tx.Model(existingOrder).Update("total_amount", totalAmount).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update total amount: %w", err)
	}
	existingOrder.TotalAmount = totalAmount

	// Update payment status based on new total
	if len(newPayments) > 0 || len(existingOrder.Payments) > 0 {
		// Calculate total paid amount from all payments
		var totalPaid float64
		for _, payment := range existingOrder.Payments {
			totalPaid += payment.Amount
		}

		// Update order status
		now := time.Now()
		updates := map[string]interface{}{}

		if totalPaid >= totalAmount {
			updates["status"] = domain.OrderStatusPaid
			updates["dish_status"] = domain.FoodStatusInProcess
			updates["paid_at"] = now
		} else {
			updates["status"] = domain.OrderStatusPending
		}

		if len(updates) > 0 {
			if err := tx.Model(existingOrder).Updates(updates).Error; err != nil {
				tx.Rollback()
				return nil, fmt.Errorf("failed to update order status: %w", err)
			}
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Return updated order with all associations
	return s.Repo.GetOrderWithAssociations(orderID)
}

func (s *OrderService) CancelOrder(orderID uuid.UUID) error {
	// Begin transaction
	tx := s.Repo.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Get current order with payments
	var order domain.Order
	if err := tx.Preload("Payments").First(&order, "id = ?", orderID).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("order not found: %w", err)
	}

	// Validate the current status - don't allow cancellation of completed orders
	if order.Status == domain.OrderStatusCancelled {
		tx.Rollback()
		return fmt.Errorf("order is already canceled")
	}

	if order.DishStatus == domain.FoodStatusCompleted {
		tx.Rollback()
		return fmt.Errorf("cannot cancel order that is already completed")
	}

	// Set current time for cancellation timestamp
	now := time.Now()

	// Update to Canceled - just change status, don't delete any records
	updates := map[string]interface{}{
		"status":       domain.OrderStatusCancelled,
		"dish_status":  domain.FoodStatusCancelled,
		"cancelled_at": now,
	}

	if err := tx.Model(&order).Updates(updates).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to cancel order: %w", err)
	}

	// Update payment status to REFUNDED if any payments exist
	if len(order.Payments) > 0 {
		// Update in database
		if err := tx.Model(&domain.Payment{}).
			Where("order_id = ?", orderID).
			Update("status", domain.PaymentStatusRefunded).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to update payment status: %w", err)
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("transaction failed: %w", err)
	}

	return nil
}
