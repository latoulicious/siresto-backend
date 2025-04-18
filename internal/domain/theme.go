package domain

import (
	"github.com/google/uuid"
)

type Theme struct {
	ID              uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name            string    `gorm:"unique;not null"` // e.g., "Dark Mode", "Corporate Blue"
	PrimaryColor    string    `gorm:"type:text"`       // e.g., "#3498db"
	SecondaryColor  string    `gorm:"type:text"`
	AccentColor     string    `gorm:"type:text"`
	BackgroundColor string    `gorm:"type:text"`
	LogoURL         string    `gorm:"type:text"`    // CDN or S3 link to logo
	FaviconURL      string    `gorm:"type:text"`    // optional
	IsDefault       bool      `gorm:"default:true"` // Optional: mark default theme
}
