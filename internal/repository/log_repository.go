package repository

import (
	"github.com/latoulicious/siresto-backend/internal/domain"
	"gorm.io/gorm"
)

type LogRepository struct {
	DB *gorm.DB
}

// LogRepository interface
type LogRepositoryInterface interface {
	SaveLog(log domain.Log) error
}

// SaveLog will store the log in the database
func (r *LogRepository) SaveLog(log domain.Log) error {
	if err := r.DB.Create(&log).Error; err != nil {
		return err
	}
	return nil
}
