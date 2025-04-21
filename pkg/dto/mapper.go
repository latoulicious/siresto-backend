package dto

import (
	"github.com/latoulicious/siresto-backend/internal/domain"
	"github.com/latoulicious/siresto-backend/pkg/db"
)

// Convert domain.Category to dto.CategoryResponse
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
func ToCreateVariationRequest(v *domain.Variation) *CreateVariationRequest {
	return &CreateVariationRequest{
		IsDefault:     v.IsDefault,
		IsAvailable:   v.IsAvailable,
		IsRequired:    v.IsRequired,
		VariationType: v.VariationType,
		Options:       toCreateVariationOptions(v.Options),
	}
}

func ToUpdateVariationRequest(v *domain.Variation) *UpdateVariationRequest {
	return &UpdateVariationRequest{
		IsDefault:     &v.IsDefault,
		IsAvailable:   &v.IsAvailable,
		IsRequired:    &v.IsRequired,
		VariationType: &v.VariationType,
		Options:       toUpdateVariationOptions(v.Options),
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

func toUpdateVariationOptions(options db.VariationOptions) []UpdateVariationOption {
	var dtoOptions []UpdateVariationOption
	for _, opt := range options {
		dtoOptions = append(dtoOptions, UpdateVariationOption{
			Label:         &opt.Label,
			PriceModifier: opt.PriceModifier,
			PriceAbsolute: opt.PriceAbsolute,
			IsDefault:     &opt.IsDefault,
		})
	}
	return dtoOptions
}
