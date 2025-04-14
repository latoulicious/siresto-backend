package domain

import (
	"time"

	"github.com/google/uuid"
)

type QRCode struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Code        string    `gorm:"type:text;unique;not null"`
	StoreID     uuid.UUID `gorm:"type:uuid"`
	TableNumber string    `gorm:"type:text"`
	Type        string    `gorm:"type:text;default:menu"`
	MenuURL     string    `gorm:"type:text"`
	ExpiresAt   *time.Time
	Orders      []Order `gorm:"foreignKey:QRID"`
}
