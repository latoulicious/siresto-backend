package handler

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/latoulicious/siresto-backend/internal/service"
	"github.com/latoulicious/siresto-backend/internal/utils"
	"github.com/skip2/go-qrcode"
)

type QRCodeHandler struct {
	Service *service.QRCodeService
}

// ListAllQRCodesHandler lists all QR codes across all stores
func (h *QRCodeHandler) ListAllQRCodesHandler(c *fiber.Ctx) error {
	qrs, err := h.Service.ListAllQRCodes()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.Error("Failed to retrieve QR codes", fiber.StatusInternalServerError))
	}
	return c.Status(fiber.StatusOK).JSON(utils.Success("QR codes retrieved successfully", qrs))
}

// GetQRCodeByIDHandler retrieves a QR by its ID
func (h *QRCodeHandler) GetQRCodeByIDHandler(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("Invalid QR ID format", fiber.StatusBadRequest))
	}

	qr, err := h.Service.GetQRCodeByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(utils.Error("QR ID not found", fiber.StatusNotFound))
	}

	return c.Status(fiber.StatusOK).JSON(utils.Success("QR code retrieved successfully", qr))
}

// ListQRCodesHandler lists all QR codes for a store
func (h *QRCodeHandler) ListQRCodesHandler(c *fiber.Ctx) error {
	storeID, err := uuid.Parse(c.Params("store_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("Invalid store ID", fiber.StatusBadRequest))
	}

	qrs, err := h.Service.ListQRCodes(storeID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.Error(err.Error(), fiber.StatusInternalServerError))
	}

	return c.Status(fiber.StatusOK).JSON(utils.Success("Store QR codes retrieved successfully", qrs))
}

// CreateQRCodeHandler creates a new QR code
func (h *QRCodeHandler) CreateQRCodeHandler(c *fiber.Ctx) error {
	var request struct {
		StoreID     uuid.UUID  `json:"store_id"`
		TableNumber string     `json:"table_number"`
		Type        string     `json:"type"`
		MenuURL     string     `json:"menu_url"`
		ExpiresAt   *time.Time `json:"expires_at"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("Invalid input", fiber.StatusBadRequest))
	}

	qrValue := fmt.Sprintf("%s?store_id=%s&table_number=%s", request.MenuURL, request.StoreID, request.TableNumber)

	qrCode, err := qrcode.Encode(qrValue, qrcode.Medium, 256)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.Error("Failed to generate QR code", fiber.StatusInternalServerError))
	}

	qrCodeBase64 := base64.StdEncoding.EncodeToString(qrCode)

	createdQR, err := h.Service.CreateQRCode(request.StoreID, request.TableNumber, request.Type, request.MenuURL, request.ExpiresAt, qrCodeBase64)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.Error("Failed to save QR code to database", fiber.StatusInternalServerError))
	}

	return c.Status(fiber.StatusCreated).JSON(utils.Success("QR code created successfully", createdQR))
}

// BulkCreateQRCodeHandler creates a new multiple QR code
func (h QRCodeHandler) BulkCreateQRCodeHandler(c *fiber.Ctx) error {
	var request struct {
		StoreID     uuid.UUID  `json:"store_id"`
		TableCount  int        `json:"table_count"`
		StartNumber int        `json:"start_number"`
		Type        string     `json:"type"`
		MenuURL     string     `json:"menu_url"`
		ExpiresAt   *time.Time `json:"expires_at"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("Invalid input", fiber.StatusBadRequest))
	}

	// Validate input
	if request.TableCount <= 0 || request.StartNumber < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("Invalid table count or start number", fiber.StatusBadRequest))
	}

	// Generate QR codes in bulk
	results, err := h.Service.BulkCreateQRCodes(
		request.StoreID,
		request.StartNumber,
		request.TableCount,
		request.Type,
		request.MenuURL,
		request.ExpiresAt,
	)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.Error("Failed to generate bulk QR codes", fiber.StatusInternalServerError))
	}

	return c.Status(fiber.StatusCreated).JSON(utils.Success("QR codes created successfully", results))
}

// UpdateQRCodeHandler updates an existing QR code
func (h *QRCodeHandler) UpdateQRCodeHandler(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("Invalid QR code ID", fiber.StatusBadRequest))
	}

	var request struct {
		TableNumber string     `json:"table_number"`
		Type        string     `json:"type"`
		MenuURL     string     `json:"menu_url"`
		ExpiresAt   *time.Time `json:"expires_at"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("Invalid request", fiber.StatusBadRequest))
	}

	qr, err := h.Service.UpdateQRCode(id, request.TableNumber, request.Type, request.MenuURL, request.ExpiresAt)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.Error(err.Error(), fiber.StatusInternalServerError))
	}

	return c.Status(fiber.StatusOK).JSON(utils.Success("QR code updated successfully", qr))
}

// DeleteQRCodeHandler deletes a QR code by its ID
func (h *QRCodeHandler) DeleteQRCodeHandler(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error("Invalid QR code ID", fiber.StatusBadRequest))
	}

	err = h.Service.DeleteQRCode(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.Error(err.Error(), fiber.StatusInternalServerError))
	}

	return c.Status(fiber.StatusNoContent).JSON(utils.Success("QR code deleted successfully", nil))
}
