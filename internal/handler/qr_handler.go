package handler

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/latoulicious/siresto-backend/internal/service"
	"github.com/skip2/go-qrcode"
)

type QRCodeHandler struct {
	Service *service.QRCodeService
}

// CreateQRCodeHandler creates a new QR code
func (h *QRCodeHandler) CreateQRCodeHandler(c *fiber.Ctx) error {
	// Input structure
	var request struct {
		StoreID     uuid.UUID  `json:"store_id"`
		TableNumber string     `json:"table_number"`
		Type        string     `json:"type"`
		MenuURL     string     `json:"menu_url"`
		ExpiresAt   *time.Time `json:"expires_at"`
	}

	// Parse request body into the request struct
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Generate the QR code data (the value to encode)
	qrValue := fmt.Sprintf("%s?store_id=%s&table_number=%s", request.MenuURL, request.StoreID, request.TableNumber)

	// Generate the QR code image (you can adjust the size as needed)
	qrCode, err := qrcode.Encode(qrValue, qrcode.Medium, 256) // Generate QR code with medium error correction
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate QR code"})
	}

	// Convert QR code image to base64 string
	qrCodeBase64 := base64.StdEncoding.EncodeToString(qrCode)

	// Create a QR code record in the database (saving the base64 image)
	createdQR, err := h.Service.CreateQRCode(request.StoreID, request.TableNumber, request.Type, request.MenuURL, request.ExpiresAt, qrCodeBase64)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save QR code to database"})
	}

	// Return the created QR code information (including the QR image as base64)
	return c.JSON(fiber.Map{
		"id":           createdQR.ID,
		"qr_code":      createdQR.Code,  // Unique code for the QR
		"qr_image":     createdQR.Image, // Base64-encoded image data
		"store_id":     createdQR.StoreID,
		"table_number": createdQR.TableNumber,
		"type":         createdQR.Type,
		"menu_url":     createdQR.MenuURL,
		"expires_at":   createdQR.ExpiresAt,
	})
}

// GetQRCodeByCodeHandler retrieves a QR code by its code
func (h *QRCodeHandler) GetQRCodeByCodeHandler(c *fiber.Ctx) error {
	code := c.Params("code")
	qr, err := h.Service.GetQRCodeByCode(code)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "QR code not found"})
	}

	return c.JSON(qr)
}

// UpdateQRCodeHandler updates an existing QR code
func (h *QRCodeHandler) UpdateQRCodeHandler(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid QR code ID"})
	}

	var request struct {
		TableNumber string     `json:"table_number"`
		Type        string     `json:"type"`
		MenuURL     string     `json:"menu_url"`
		ExpiresAt   *time.Time `json:"expires_at"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	qr, err := h.Service.UpdateQRCode(id, request.TableNumber, request.Type, request.MenuURL, request.ExpiresAt)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(qr)
}

// DeleteQRCodeHandler deletes a QR code by its ID
func (h *QRCodeHandler) DeleteQRCodeHandler(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid QR code ID"})
	}

	err = h.Service.DeleteQRCode(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// ListQRCodesHandler lists all QR codes for a store
func (h *QRCodeHandler) ListQRCodesHandler(c *fiber.Ctx) error {
	storeID, err := uuid.Parse(c.Params("store_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid store ID"})
	}

	qrs, err := h.Service.ListQRCodes(storeID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(qrs)
}
