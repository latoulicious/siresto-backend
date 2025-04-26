package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name        string    `gorm:"type:text;not null"`
	Email       string    `gorm:"type:text;unique;not null"`
	Password    string    `gorm:"type:text" json:"-"`
	IsStaff     bool      `gorm:"default:false"`
	RoleID      uuid.UUID
	Role        *Role     `gorm:"foreignKey:RoleID"`
	CreatedAt   time.Time `gorm:"default:now()"`
	LastLoginAt *time.Time
}
