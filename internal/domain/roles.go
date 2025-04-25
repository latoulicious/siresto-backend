package domain

import (
	"time"

	"github.com/google/uuid"
)

type Role struct {
	ID          uuid.UUID     `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name        string        `gorm:"unique;not null"`
	Description string        `gorm:"type:text"`
	Position    int           `gorm:"not null;default:100"`   // Lower numbers = higher privilege
	IsSystem    bool          `gorm:"not null;default:false"` // System roles can't be modified by non-system roles
	IsStaff     bool          `gorm:"not null;default:true"`  // Whether users with this role are considered staff
	Permissions []*Permission `gorm:"many2many:role_permissions"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
