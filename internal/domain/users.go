package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name         string    `gorm:"type:text;not null"`
	Email        string    `gorm:"type:text;unique;not null"`
	PasswordHash string    `gorm:"type:text"`
	IsStaff      bool      `gorm:"default:false"`
	RoleID       *int      `gorm:"foreignKey:RoleID"`
	Role         *Role     `gorm:"foreignKey:RoleID"`
	CreatedAt    time.Time `gorm:"default:now()"`
	LastLoginAt  *time.Time
}
