package service

import (
	"time"

	"github.com/google/uuid"
	"github.com/latoulicious/siresto-backend/internal/domain"
	"github.com/latoulicious/siresto-backend/internal/repository"
)

type QRCodeService struct {
	Repo *repository.QRCodeRepository
}

// ListAllQRCodes fetches all QR codes across all stores
func (s *QRCodeService) ListAllQRCodes() ([]domain.QRCode, error) {
	return s.Repo.ListAllQRCodes()
}

// GetQRCodeByID fetches a QR code by its ID
func (s *QRCodeService) GetQRCodeByID(id uuid.UUID) (*domain.QRCode, error) {
	return s.Repo.GetQRCodeByID(id)
}

// ListQRCodes fetches all QR codes for a store
func (s *QRCodeService) ListQRCodes(storeID uuid.UUID) ([]domain.QRCode, error) {
	return s.Repo.ListQRCodes(storeID)
}

// CreateQRCode creates a new QR code
func (s *QRCodeService) CreateQRCode(storeID uuid.UUID, tableNumber string, qrType string, menuURL string, expiresAt *time.Time, qrImage string) (*domain.QRCode, error) {
	qr := &domain.QRCode{
		Code:        uuid.New().String(), // Generate a unique code for the QR
		StoreID:     storeID,
		TableNumber: tableNumber,
		Type:        qrType,
		MenuURL:     menuURL,
		ExpiresAt:   expiresAt,
		Image:       qrImage, // Save the base64-encoded image
	}

	// Save the QR code to the database
	err := s.Repo.CreateQRCode(qr)
	if err != nil {
		return nil, err
	}
	return qr, nil
}

// UpdateQRCode updates an existing QR code
func (s *QRCodeService) UpdateQRCode(id uuid.UUID, tableNumber string, qrType string, menuURL string, expiresAt *time.Time) (*domain.QRCode, error) {
	qr, err := s.Repo.GetQRCodeByID(id)
	if err != nil {
		return nil, err
	}
	qr.TableNumber = tableNumber
	qr.Type = qrType
	qr.MenuURL = menuURL
	qr.ExpiresAt = expiresAt
	err = s.Repo.UpdateQRCode(qr)
	if err != nil {
		return nil, err
	}
	return qr, nil
}

// DeleteQRCode deletes a QR code by its ID
func (s *QRCodeService) DeleteQRCode(id uuid.UUID) error {
	return s.Repo.DeleteQRCode(id)
}
