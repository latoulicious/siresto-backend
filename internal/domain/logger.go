package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/latoulicious/siresto-backend/pkg/db"
)

type Log struct {
	ID          uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	Timestamp   time.Time  `gorm:"type:timestamp;not null;default:now()" json:"timestamp"`
	Level       string     `gorm:"type:varchar(20);not null" json:"level"`       // info, warn, error, audit
	Source      string     `gorm:"type:varchar(100);not null" json:"source"`     // service/module/component
	UserID      *uuid.UUID `gorm:"type:uuid" json:"user_id"`                     // optional
	IPAddress   *string    `gorm:"type:inet" json:"ip_address"`                  // optional
	Action      string     `gorm:"type:varchar(100);not null" json:"action"`     // e.g., "menu.created"
	Entity      string     `gorm:"type:varchar(50);not null" json:"entity"`      // e.g., "menu"
	EntityID    *uuid.UUID `gorm:"type:uuid" json:"entity_id"`                   // optional
	Description string     `gorm:"type:text" json:"description"`                 // human readable
	Metadata    db.JSONB   `gorm:"type:jsonb" json:"metadata"`                   // arbitrary context (JSONB)
	RequestID   *string    `gorm:"type:varchar(100)" json:"request_id"`          // optional, for traceability
	Environment string     `gorm:"type:varchar(20);not null" json:"environment"` // e.g., "production"
	Application string     `gorm:"type:varchar(50);not null" json:"application"` // e.g., "restaurant-api"
	Hostname    *string    `gorm:"type:varchar(100)" json:"hostname"`            // optional
	Type        string     `gorm:"type:varchar(20);not null" json:"type"`        // "activity", "audit", etc.
}
