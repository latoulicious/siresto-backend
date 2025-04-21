package dto

import (
	"github.com/google/uuid"
)

// --- Request DTOs ---
type CreateCategoryRequest struct {
	Name     string `json:"name" binding:"required"`
	IsActive *bool  `json:"is_active,omitempty"`
	Position *int   `json:"position,omitempty"`
}

type UpdateCategoryRequest struct {
	Name     *string `json:"name,omitempty"`
	IsActive *bool   `json:"is_active,omitempty"`
	Position *int    `json:"position,omitempty"`
}

// --- Response DTOs ---
type CategoryResponse struct {
	ID       uuid.UUID        `json:"id"`
	Name     string           `json:"name"`
	IsActive bool             `json:"is_active"`
	Position int              `json:"position"`
	Products []ProductSummary `json:"products,omitempty"`
}

type ProductSummary struct {
	ID          uuid.UUID          `json:"id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	ImageURL    string             `json:"image_url"`
	BasePrice   float64            `json:"base_price"`
	IsAvailable bool               `json:"is_available"`
	Position    int                `json:"position"`
	Variations  []VariationSummary `json:"variations,omitempty"`
}

type VariationSummary struct {
	ID            uuid.UUID         `json:"id"`
	IsDefault     bool              `json:"is_default"`
	IsAvailable   bool              `json:"is_available"`
	IsRequired    bool              `json:"is_required"`
	VariationType string            `json:"variation_type"`
	Options       []VariationOption `json:"options"`
}

type VariationOption struct {
	Label         string   `json:"label"`
	PriceModifier *float64 `json:"price_modifier,omitempty"`
	PriceAbsolute *float64 `json:"price_absolute,omitempty"`
	IsDefault     bool     `json:"is_default,omitempty"`
}
