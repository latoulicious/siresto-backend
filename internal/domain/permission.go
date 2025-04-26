package domain

import (
	"github.com/google/uuid"
)

type Permission struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name        string    `gorm:"unique;not null"`
	Description string
	Roles       []*Role `gorm:"many2many:role_permissions"`
}
