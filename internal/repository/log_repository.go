package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/latoulicious/siresto-backend/internal/domain"
	"gorm.io/gorm"
)

type LogRepositoryInterface interface {
	SaveLog(ctx context.Context, log domain.Log) error
	GetLogByID(ctx context.Context, id string) (*domain.Log, error)
}

type LogRepository struct {
	DB *gorm.DB
}

func (r *LogRepository) SaveLog(ctx context.Context, log domain.Log) error {
	// Auto-populate UUID if not set
	if log.ID == uuid.Nil {
		log.ID = uuid.New()
	}

	// Auto-populate timestamp if zero
	if log.Timestamp.IsZero() {
		log.Timestamp = time.Now()
	}

	result := r.DB.Create(&log)
	return result.Error
}

func (r *LogRepository) GetLogByID(ctx context.Context, id string) (*domain.Log, error) {
	var log domain.Log

	parsedID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	if err := r.DB.Where("id = ?", parsedID).First(&log).Error; err != nil {
		return nil, err
	}

	return &log, nil
}
