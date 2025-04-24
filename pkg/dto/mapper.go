package dto

import (
	"github.com/latoulicious/siresto-backend/internal/domain"
	"github.com/latoulicious/siresto-backend/pkg/db"
)

// Category DTO
func ToCategoryResponse(c *domain.Category) *CategoryResponse {
	category := &CategoryResponse{
		ID:       c.ID,
		Name:     c.Name,
		IsActive: c.IsActive,
		Position: c.Position,
	}

	for _, product := range c.Products {
		category.Products = append(category.Products, ToProductSummary(&product))
	}

	return category
}

func ToProductSummary(p *domain.Product) ProductSummary {
	product := ProductSummary{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		ImageURL:    p.ImageURL,
		BasePrice:   p.BasePrice,
		IsAvailable: p.IsAvailable,
		Position:    p.Position,
	}

	for _, v := range p.Variations {
		product.Variations = append(product.Variations, ToVariationSummary(&v))
	}

	return product
}

func ToVariationSummary(v *domain.Variation) VariationSummary {
	return VariationSummary{
		ID:            v.ID,
		IsDefault:     v.IsDefault,
		IsAvailable:   v.IsAvailable,
		IsRequired:    v.IsRequired,
		VariationType: v.VariationType,
		Options:       toVariationOptions(v.Options),
	}
}

func toVariationOptions(options db.VariationOptions) []VariationOption {
	var dtoOptions []VariationOption
	for _, opt := range options {
		dtoOptions = append(dtoOptions, VariationOption{
			Label:         opt.Label,
			PriceModifier: opt.PriceModifier,
			PriceAbsolute: opt.PriceAbsolute,
			IsDefault:     opt.IsDefault,
		})
	}
	return dtoOptions
}

// Product DTO
func ToProductResponse(p *domain.Product) *ProductResponse {
	product := &ProductResponse{
		ID:           p.ID,
		Name:         p.Name,
		Description:  p.Description,
		ImageURL:     p.ImageURL,
		BasePrice:    p.BasePrice,
		IsAvailable:  p.IsAvailable,
		Position:     p.Position,
		CategoryID:   p.CategoryID,
		CategoryName: p.CategoryName,
	}

	for _, v := range p.Variations {
		product.Variations = append(product.Variations, ToVariationSummary(&v))
	}
	return product
}

func ToCreateProductRequest(p *domain.Product) *CreateProductRequest {
	return &CreateProductRequest{
		Name:        p.Name,
		Description: p.Description,
		ImageURL:    p.ImageURL,
		BasePrice:   p.BasePrice,
		IsAvailable: p.IsAvailable,
		Position:    p.Position,
		CategoryID:  p.CategoryID,
	}
}

func ToUpdateProductRequest(p *domain.Product) *UpdateProductRequest {
	return &UpdateProductRequest{
		Name:        &p.Name,
		Description: &p.Description,
		ImageURL:    &p.ImageURL,
		BasePrice:   &p.BasePrice,
		IsAvailable: &p.IsAvailable,
		Position:    &p.Position,
		CategoryID:  p.CategoryID,
	}
}

// Mapping DTO back to &Domain
func ToProductDomainFromCreate(request *CreateProductRequest) *domain.Product {
	return &domain.Product{
		Name:        request.Name,
		Description: request.Description,
		ImageURL:    request.ImageURL,
		BasePrice:   request.BasePrice,
		IsAvailable: request.IsAvailable,
		Position:    request.Position,
		CategoryID:  request.CategoryID,
	}
}

// Mapping DTO back to &Domain
func ToProductDomainFromUpdate(request *UpdateProductRequest, existingProduct *domain.Product) *domain.Product {
	// Create a copy of the existing product to preserve unchanged values
	updatedProduct := *existingProduct

	// Only update non-nil values
	if request.Name != nil {
		updatedProduct.Name = *request.Name
	}
	if request.Description != nil {
		updatedProduct.Description = *request.Description
	}
	if request.ImageURL != nil {
		updatedProduct.ImageURL = *request.ImageURL
	}
	if request.BasePrice != nil {
		updatedProduct.BasePrice = *request.BasePrice
	}
	if request.IsAvailable != nil {
		updatedProduct.IsAvailable = *request.IsAvailable
	}
	if request.Position != nil {
		updatedProduct.Position = *request.Position
	}
	if request.CategoryID != nil {
		updatedProduct.CategoryID = request.CategoryID
	}

	return &updatedProduct
}

// Variation DTO
func ToVariationResponses(variations []*domain.Variation) []VariationSummary {
	var variationResponses []VariationSummary
	for _, variation := range variations {
		if variation == nil {
			continue
		}

		variationResponse := VariationSummary{
			ID:            variation.ID,
			IsDefault:     variation.IsDefault,
			IsAvailable:   variation.IsAvailable,
			IsRequired:    variation.IsRequired,
			VariationType: variation.VariationType,
			Options:       toVariationOptions(variation.Options),
		}
		variationResponses = append(variationResponses, variationResponse)
	}
	return variationResponses
}

func ToCreateVariationRequest(v *domain.Variation) *CreateVariationRequest {
	if v == nil {
		return nil
	}

	return &CreateVariationRequest{
		IsDefault:     v.IsDefault,
		IsAvailable:   v.IsAvailable,
		IsRequired:    v.IsRequired,
		VariationType: v.VariationType,
		Options:       toCreateVariationOptions(v.Options),
	}
}

func ToUpdateVariationRequest(v *domain.Variation) *UpdateVariationRequest {
	if v == nil {
		return nil
	}

	return &UpdateVariationRequest{
		IsDefault:     &v.IsDefault,
		IsAvailable:   &v.IsAvailable,
		IsRequired:    &v.IsRequired,
		VariationType: &v.VariationType,
		Options:       toVariationOptionsUpdate(v.Options),
	}
}

func toCreateVariationOptions(options db.VariationOptions) []CreateVariationOption {
	var dtoOptions []CreateVariationOption
	for _, opt := range options {
		dtoOptions = append(dtoOptions, CreateVariationOption{
			Label:         opt.Label,
			PriceModifier: opt.PriceModifier,
			PriceAbsolute: opt.PriceAbsolute,
			IsDefault:     opt.IsDefault,
		})
	}
	return dtoOptions
}

func toVariationOptionsUpdate(options db.VariationOptions) []UpdateVariationOption {
	var updateOpts []UpdateVariationOption
	for _, opt := range options {
		label := opt.Label
		isDefault := opt.IsDefault
		updateOpts = append(updateOpts, UpdateVariationOption{
			Label:         &label,
			PriceModifier: opt.PriceModifier,
			PriceAbsolute: opt.PriceAbsolute,
			IsDefault:     &isDefault,
		})
	}
	return updateOpts
}

// Mapping DTO back to &Domain
func toVariationOptionsDomain(options []UpdateVariationOption) db.VariationOptions {
	var dbOptions db.VariationOptions
	for _, opt := range options {
		// Convert UpdateVariationOption to CreateVariationOption
		dbOptions = append(dbOptions, db.VariationOption{
			Label:         *opt.Label,
			PriceModifier: opt.PriceModifier,
			PriceAbsolute: opt.PriceAbsolute,
			IsDefault:     *opt.IsDefault,
		})
	}
	return dbOptions
}

func convertCreateToUpdateOptions(createOpts []CreateVariationOption) []UpdateVariationOption {
	var updateOpts []UpdateVariationOption
	for _, opt := range createOpts {
		updateOpts = append(updateOpts, UpdateVariationOption{
			Label:         &opt.Label,
			PriceModifier: opt.PriceModifier,
			PriceAbsolute: opt.PriceAbsolute,
			IsDefault:     &opt.IsDefault,
		})
	}
	return updateOpts
}

// Mapping DTO back to &Domain
func ToVariationDomain(request *CreateVariationRequest) *domain.Variation {
	return &domain.Variation{
		ProductID:     *request.ProductID,
		IsDefault:     request.IsDefault,
		IsAvailable:   request.IsAvailable,
		IsRequired:    request.IsRequired,
		VariationType: request.VariationType,
		Options:       toVariationOptionsDomain(convertCreateToUpdateOptions(request.Options)),
	}
}

// Mapping DTO back to &Domain
func ToVariationDomainFromUpdate(request *UpdateVariationRequest) *domain.Variation {
	variation := &domain.Variation{}

	// Set ID if provided
	if request.ID != nil {
		variation.ID = *request.ID
	}

	// Handle all pointer fields properly
	if request.IsDefault != nil {
		variation.IsDefault = *request.IsDefault
	}
	if request.IsAvailable != nil {
		variation.IsAvailable = *request.IsAvailable
	}
	if request.IsRequired != nil {
		variation.IsRequired = *request.IsRequired
	}
	if request.VariationType != nil {
		variation.VariationType = *request.VariationType
	}

	// Handle options
	if len(request.Options) > 0 {
		variation.Options = toVariationOptionsDomain(request.Options)
	}

	return variation
}

// Mapping DTO back to &Domain
func ToVariationOptionsDomainFromUpdate(options []UpdateVariationOption) db.VariationOptions {
	var dbOptions db.VariationOptions
	for _, opt := range options {
		// Skip if required fields are nil
		if opt.Label == nil {
			continue
		}

		option := db.VariationOption{
			Label: *opt.Label,
		}

		// Handle optional fields
		if opt.PriceModifier != nil {
			option.PriceModifier = opt.PriceModifier
		}

		if opt.PriceAbsolute != nil {
			option.PriceAbsolute = opt.PriceAbsolute
		}

		if opt.IsDefault != nil {
			option.IsDefault = *opt.IsDefault
		}

		dbOptions = append(dbOptions, option)
	}
	return dbOptions
}

// Order DTO
func MapToOrderResponseDTO(order *domain.Order) OrderResponseDTO {
	// Map order details to order items
	items := make([]OrderItemDTO, 0, len(order.OrderDetails))

	for _, detail := range order.OrderDetails {
		// Extract variation information
		variationName := ""
		if detail.Variation != nil && len(detail.Variation.Options) > 0 {
			// Find default option for variation name
			for _, opt := range detail.Variation.Options {
				if opt.IsDefault {
					variationName = opt.Label
					break
				}
			}
		}

		// Calculate actual price factoring in variations
		unitPrice := detail.UnitPrice
		if unitPrice == 0 && detail.Product != nil {
			unitPrice = detail.Product.BasePrice
		}

		totalPrice := detail.TotalPrice
		if totalPrice == 0 {
			totalPrice = float64(detail.Quantity) * unitPrice
		}

		// Get product name either from detail or related product
		productName := detail.ProductName
		if productName == "" && detail.Product != nil {
			productName = detail.Product.Name
		}

		productID := ""
		if detail.ProductID != nil {
			productID = detail.ProductID.String()
		}

		// Get image URL from the product
		imageURL := ""
		if detail.Product != nil {
			imageURL = detail.Product.ImageURL
		}

		items = append(items, OrderItemDTO{
			ID:          detail.ID.String(),
			ProductID:   productID,
			ProductName: productName,
			Variation:   variationName,
			Quantity:    detail.Quantity,
			UnitPrice:   unitPrice,
			TotalPrice:  totalPrice,
			Note:        detail.Note,
			ImageURL:    imageURL,
		})
	}

	// Map payment methods - just include methods without other details
	paymentMethods := make([]map[string]string, 0, len(order.Payments))
	for _, payment := range order.Payments {
		paymentMethods = append(paymentMethods, map[string]string{
			"method": string(payment.Method),
		})
	}

	return OrderResponseDTO{
		ID:            order.ID.String(),
		CustomerName:  order.CustomerName,
		CustomerPhone: order.CustomerPhone,
		TableNumber:   order.TableNumber,
		Status:        string(order.Status),
		DishStatus:    string(order.DishStatus),
		TotalAmount:   order.TotalAmount,
		Notes:         order.Notes,
		CreatedAt:     order.CreatedAt,
		PaidAt:        order.PaidAt,
		Items:         items,
		Methods:       paymentMethods, // Add the methods to the response
	}
}
