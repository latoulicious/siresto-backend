package domain

import (
	"github.com/google/uuid"
)

type Theme struct {
	ID              uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	Name            string    `gorm:"unique;not null" json:"name"`    // e.g., "Dark Mode", "Corporate Blue"
	PrimaryColor    string    `gorm:"type:text" json:"primary_color"` // e.g., "#3498db"
	SecondaryColor  string    `gorm:"type:text" json:"secondary_color"`
	AccentColor     string    `gorm:"type:text" json:"accent_color"`
	BackgroundColor string    `gorm:"type:text" json:"background_color"`
	LogoURL         string    `gorm:"type:text" json:"logo_url"`      // CDN or S3 link to logo
	FaviconURL      string    `gorm:"type:text" json:"favicon_url"`   // optional
	IsDefault       bool      `gorm:"default:true" json:"is_default"` // Optional: mark default theme
}
