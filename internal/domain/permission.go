package domain

import (
	"github.com/google/uuid"
)

type Permission struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name        string    `gorm:"unique;not null"` // e.g., "read:users"
	Description string
	Roles       []*Role `gorm:"many2many:role_permissions"`
}
