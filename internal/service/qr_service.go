package service

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/latoulicious/siresto-backend/internal/domain"
	"github.com/latoulicious/siresto-backend/internal/repository"
	"github.com/skip2/go-qrcode"
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

// BulkCreateQRCode creates a new multiple QR code
func (s QRCodeService) BulkCreateQRCodes(
	storeID uuid.UUID,
	startNumber int,
	count int,
	qrType string,
	menuURL string,
	expiresAt *time.Time,
) ([]*domain.QRCode, error) {
	results := make([]*domain.QRCode, 0, count)

	for i := 0; i < count; i++ {
		tableNumber := fmt.Sprintf("%d", startNumber+i)

		// Generate QR code base64 image
		qrValue := fmt.Sprintf("%s?store_id=%s&table_number=%s", menuURL, storeID, tableNumber)
		qrCode, err := qrcode.Encode(qrValue, qrcode.Medium, 256)
		if err != nil {
			return nil, fmt.Errorf("failed to generate QR code for table %s: %w", tableNumber, err)
		}
		qrCodeBase64 := base64.StdEncoding.EncodeToString(qrCode)

		// Create QR code entry
		qr := &domain.QRCode{
			Code:        uuid.New().String(),
			StoreID:     storeID,
			TableNumber: tableNumber,
			Type:        qrType,
			MenuURL:     menuURL,
			ExpiresAt:   expiresAt,
			Image:       qrCodeBase64,
		}

		// Save the QR code to the database
		err = s.Repo.CreateQRCode(qr)
		if err != nil {
			return nil, fmt.Errorf("failed to save QR code for table %s: %w", tableNumber, err)
		}

		results = append(results, qr)
	}

	return results, nil
}

// DeleteQRCode deletes a QR code by its ID
func (s *QRCodeService) DeleteQRCode(id uuid.UUID) error {
	return s.Repo.DeleteQRCode(id)
}
