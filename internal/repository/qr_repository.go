package repository

import (
	"github.com/google/uuid"
	"github.com/latoulicious/siresto-backend/internal/domain"
	"gorm.io/gorm"
)

type QRCodeRepository struct {
	DB *gorm.DB
}

// ListAllQRCodes fetches all QR codes from the database with pagination
func (r *QRCodeRepository) ListAllQRCodes(page, perPage int) ([]domain.QRCode, int64, error) {
	var qrs []domain.QRCode
	var totalCount int64

	// Get total count
	if err := r.DB.Model(&domain.QRCode{}).Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// Calculate offset
	offset := (page - 1) * perPage

	// Get paginated data
	err := r.DB.Offset(offset).Limit(perPage).Find(&qrs).Error
	if err != nil {
		return nil, 0, err
	}

	return qrs, totalCount, nil
}

// GetQRCodeByID fetches a QR code by its ID
func (r *QRCodeRepository) GetQRCodeByID(id uuid.UUID) (*domain.QRCode, error) {
	var qr domain.QRCode
	err := r.DB.Where("id = ?", id).First(&qr).Error
	if err != nil {
		return nil, err
	}
	return &qr, nil
}

// ListQRCodes fetches all QR codes for a specific store with pagination
func (r *QRCodeRepository) ListQRCodes(storeID uuid.UUID, page, perPage int) ([]domain.QRCode, int64, error) {
	var qrs []domain.QRCode
	var totalCount int64

	// Get total count for the store
	if err := r.DB.Model(&domain.QRCode{}).Where("store_id = ?", storeID).Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// Calculate offset
	offset := (page - 1) * perPage

	// Get paginated data for the store
	err := r.DB.Where("store_id = ?", storeID).Offset(offset).Limit(perPage).Find(&qrs).Error
	if err != nil {
		return nil, 0, err
	}

	return qrs, totalCount, nil
}

// CreateQRCode will create a new QR code record in the database
func (r *QRCodeRepository) CreateQRCode(qr *domain.QRCode) error {
	return r.DB.Create(qr).Error
}

// CreateQRCode will create multiple new QR code record in the database
func (r QRCodeRepository) BulkCreateQRCodes(qrCodes []*domain.QRCode) error {
	// Use a transaction for bulk inserts
	tx := r.DB.Begin()

	for _, qr := range qrCodes {
		if err := tx.Create(qr).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

// DeleteQRCode deletes a QR code from the database
func (r *QRCodeRepository) DeleteQRCode(id uuid.UUID) error {
	return r.DB.Where("id = ?", id).Delete(&domain.QRCode{}).Error
}
