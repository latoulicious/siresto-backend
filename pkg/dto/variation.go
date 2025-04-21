package dto

import "github.com/google/uuid"

// --- Request DTOs ---
type CreateVariationRequest struct {
	ProductID     *uuid.UUID              `json:"product_id" binding:"required"`
	IsDefault     bool                    `json:"is_default" binding:"required"`
	IsAvailable   bool                    `json:"is_available" binding:"required"`
	IsRequired    bool                    `json:"is_required" binding:"required"`
	VariationType string                  `json:"variation_type" binding:"required"`
	Options       []CreateVariationOption `json:"options" binding:"required"`
}

type UpdateVariationRequest struct {
	ID            *string                 `json:"id"`
	IsDefault     *bool                   `json:"isDefault"`
	IsAvailable   *bool                   `json:"isAvailable"`
	IsRequired    *bool                   `json:"isRequired"`
	VariationType *string                 `json:"variationType"`
	Options       []UpdateVariationOption `json:"options"`
}

// --- Response DTOs ---
type VariationResponse struct {
	ID            string                    `json:"id"`
	IsDefault     bool                      `json:"is_default"`
	IsAvailable   bool                      `json:"is_available"`
	IsRequired    bool                      `json:"is_required"`
	VariationType string                    `json:"variation_type"`
	Options       []VariationOptionResponse `json:"options"`
}

type VariationOptionResponse struct {
	Label         string   `json:"label"`
	PriceModifier *float64 `json:"price_modifier,omitempty"`
	PriceAbsolute *float64 `json:"price_absolute,omitempty"`
	IsDefault     bool     `json:"is_default,omitempty"`
}

type CreateVariationResponse struct {
	ProductID     uuid.UUID                       `json:"product_id"`
	IsDefault     bool                            `json:"is_default"`
	IsAvailable   bool                            `json:"is_available"`
	IsRequired    bool                            `json:"is_required"`
	VariationType string                          `json:"variation_type"`
	Options       []CreateVariationOptionResponse `json:"options"`
}

type CreateVariationOptionResponse struct {
	Label         string  `json:"label"`
	PriceModifier float64 `json:"price_modifier"`
	PriceAbsolute float64 `json:"price_absolute"`
	IsDefault     bool    `json:"is_default"`
}

type CreateVariationOption struct {
	Label         string   `json:"label" binding:"required"`
	PriceModifier *float64 `json:"price_modifier,omitempty"`
	PriceAbsolute *float64 `json:"price_absolute,omitempty"`
	IsDefault     bool     `json:"is_default,omitempty"`
}

type UpdateVariationOption struct {
	Label         *string  `json:"label,omitempty"`
	PriceModifier *float64 `json:"price_modifier,omitempty"`
	PriceAbsolute *float64 `json:"price_absolute,omitempty"`
	IsDefault     *bool    `json:"is_default,omitempty"`
}
