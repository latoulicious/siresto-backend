package db

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// VariationOption represents one choice inside a variation
type VariationOption struct {
	Label         string   `json:"label"`
	PriceModifier *float64 `json:"price_modifier,omitempty"`
	PriceAbsolute *float64 `json:"price_absolute,omitempty"`
	IsDefault     bool     `json:"is_default,omitempty"`
}

// VariationOptions is a custom JSONB wrapper for an array of VariationOption
type VariationOptions []VariationOption

func (v VariationOptions) Value() (driver.Value, error) {
	bytes, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal VariationOptions: %w", err)
	}
	return bytes, nil
}

func (v *VariationOptions) Scan(src interface{}) error {
	if src == nil {
		*v = nil
		return nil
	}

	bytes, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("VariationOptions scan: type assertion to []byte failed")
	}

	if err := json.Unmarshal(bytes, v); err != nil {
		return fmt.Errorf("VariationOptions scan: failed to unmarshal: %w", err)
	}
	return nil
}
