package dto

import (
	"github.com/google/uuid"
)

// --- Request DTOs ---
type CreateProductRequest struct {
	Name        string                   `json:"name" binding:"required"`
	Description string                   `json:"description" binding:"required"`
	ImageURL    string                   `json:"image_url" binding:"required"` // FIX: Proper JSON binding
	BasePrice   float64                  `json:"base_price" binding:"required"`
	IsAvailable bool                     `json:"is_available" binding:"required"`
	Position    int                      `json:"position" binding:"required"`
	CategoryID  *uuid.UUID               `json:"category_id" binding:"required"`
	Variations  []CreateVariationRequest `json:"variations,omitempty"`
}

type UpdateProductRequest struct {
	Name                  *string                  `json:"name,omitempty"`
	Description           *string                  `json:"description,omitempty"`
	ImageURL              *string                  `json:"image_url,omitempty"`
	BasePrice             *float64                 `json:"base_price,omitempty"`
	IsAvailable           *bool                    `json:"is_available,omitempty"`
	Position              *int                     `json:"position,omitempty"`
	CategoryID            *uuid.UUID               `json:"category_id,omitempty"`
	Variations            []UpdateVariationRequest `json:"variations,omitempty"`
	RemoveOtherVariations *bool                    `json:"remove_other_variations,omitempty"`
}

// --- Response DTOs ---
type ProductResponse struct {
	ID           uuid.UUID          `json:"id"`
	CategoryID   *uuid.UUID         `json:"category_id"`
	CategoryName string             `json:"category_name"`
	Name         string             `json:"name"`
	Description  string             `json:"description"`
	ImageURL     string             `json:"image_url"`
	BasePrice    float64            `json:"base_price"`
	IsAvailable  bool               `json:"is_available"`
	Position     int                `json:"position"`
	Variations   []VariationSummary `json:"variations,omitempty"`
}
