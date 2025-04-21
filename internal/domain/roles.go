package domain

import (
	"time"

	"github.com/google/uuid"
)

type Role struct {
	ID          uuid.UUID     `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Name        string        `gorm:"unique;not null"`
	Description string        `gorm:"type:text"`
	Permissions []*Permission `gorm:"many2many:role_permissions"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
