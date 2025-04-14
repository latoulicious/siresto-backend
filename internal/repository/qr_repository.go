package repository

import (
	"github.com/google/uuid"
	"github.com/latoulicious/siresto-backend/internal/domain"
	"gorm.io/gorm"
)

type QRCodeRepository struct {
	DB *gorm.DB
}

// ListAllQRCodes fetches all QR codes from the database
func (r *QRCodeRepository) ListAllQRCodes() ([]domain.QRCode, error) {
	var qrs []domain.QRCode
	err := r.DB.Find(&qrs).Error
	if err != nil {
		return nil, err
	}
	return qrs, nil
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

// GetQRCodeByCode fetches a QR code by its unique code
// func (r *QRCodeRepository) GetQRCodeByCode(code string) (*domain.QRCode, error) {
// 	var qr domain.QRCode
// 	err := r.DB.Where("code = ?", code).First(&qr).Error
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &qr, nil
// }

// ListQRCodes fetches all QR codes for a specific store
func (r *QRCodeRepository) ListQRCodes(storeID uuid.UUID) ([]domain.QRCode, error) {
	var qrs []domain.QRCode
	err := r.DB.Where("store_id = ?", storeID).Find(&qrs).Error
	if err != nil {
		return nil, err
	}
	return qrs, nil
}

// CreateQRCode will create a new QR code record in the database
func (r *QRCodeRepository) CreateQRCode(qr *domain.QRCode) error {
	return r.DB.Create(qr).Error
}

// UpdateQRCode updates an existing QR code record
func (r *QRCodeRepository) UpdateQRCode(qr *domain.QRCode) error {
	return r.DB.Save(qr).Error
}

// DeleteQRCode deletes a QR code from the database
func (r *QRCodeRepository) DeleteQRCode(id uuid.UUID) error {
	return r.DB.Where("id = ?", id).Delete(&domain.QRCode{}).Error
}
