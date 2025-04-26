package handler

import (
	"encoding/base64"
	"fmt"
	"strconv"
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
	// Get pagination parameters from query
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "10"))

	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 10
	}

	qrs, totalCount, err := h.Service.ListAllQRCodes(page, perPage)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.Error(
			"Failed to retrieve QR codes",
			fiber.StatusInternalServerError,
			utils.NewErrorInfo("INTERNAL_ERROR", err.Error(), "", nil),
		))
	}

	metadata := utils.NewPaginationMetadata(page, perPage, int(totalCount))
	return c.Status(fiber.StatusOK).JSON(utils.Success("QR codes retrieved successfully", qrs, metadata))
}

// GetQRCodeByIDHandler retrieves a QR by its ID
func (h *QRCodeHandler) GetQRCodeByIDHandler(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error(
			"Invalid QR ID format",
			fiber.StatusBadRequest,
			utils.NewErrorInfo("INVALID_INPUT", "QR ID must be a valid UUID", "id", nil),
		))
	}

	qr, err := h.Service.GetQRCodeByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(utils.Error(
			"QR ID not found",
			fiber.StatusNotFound,
			utils.NewErrorInfo("NOT_FOUND", fmt.Sprintf("QR code with ID %s not found", id), "id", nil),
		))
	}

	return c.Status(fiber.StatusOK).JSON(utils.Success("QR code retrieved successfully", qr, nil))
}

// ListQRCodesHandler lists all QR codes for a store
func (h *QRCodeHandler) ListQRCodesHandler(c *fiber.Ctx) error {
	storeID, err := uuid.Parse(c.Params("store_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error(
			"Invalid store ID",
			fiber.StatusBadRequest,
			utils.NewErrorInfo("INVALID_INPUT", "Store ID must be a valid UUID", "store_id", nil),
		))
	}

	// Get pagination parameters from query
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "10"))

	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 10
	}

	qrs, totalCount, err := h.Service.ListQRCodes(storeID, page, perPage)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.Error(
			"Failed to retrieve store QR codes",
			fiber.StatusInternalServerError,
			utils.NewErrorInfo("INTERNAL_ERROR", err.Error(), "", nil),
		))
	}

	metadata := utils.NewPaginationMetadata(page, perPage, int(totalCount))
	return c.Status(fiber.StatusOK).JSON(utils.Success("Store QR codes retrieved successfully", qrs, metadata))
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
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error(
			"Invalid input",
			fiber.StatusBadRequest,
			utils.NewErrorInfo("INVALID_INPUT", "Failed to parse request body", "", nil),
		))
	}

	qrValue := fmt.Sprintf("%s?store_id=%s&table_number=%s", request.MenuURL, request.StoreID, request.TableNumber)

	qrCode, err := qrcode.Encode(qrValue, qrcode.Medium, 256)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.Error(
			"Failed to generate QR code",
			fiber.StatusInternalServerError,
			utils.NewErrorInfo("QR_GENERATION_ERROR", err.Error(), "", nil),
		))
	}

	qrCodeBase64 := base64.StdEncoding.EncodeToString(qrCode)

	createdQR, err := h.Service.CreateQRCode(request.StoreID, request.TableNumber, request.Type, request.MenuURL, request.ExpiresAt, qrCodeBase64)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.Error(
			"Failed to save QR code to database",
			fiber.StatusInternalServerError,
			utils.NewErrorInfo("DB_ERROR", err.Error(), "", nil),
		))
	}

	return c.Status(fiber.StatusCreated).JSON(utils.Success("QR code created successfully", createdQR, nil))
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
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error(
			"Invalid input",
			fiber.StatusBadRequest,
			utils.NewErrorInfo("INVALID_INPUT", "Failed to parse request body", "", nil),
		))
	}

	if request.TableCount <= 0 || request.StartNumber < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error(
			"Invalid table count or start number",
			fiber.StatusBadRequest,
			utils.NewErrorInfo("VALIDATION_ERROR", "Table count must be positive and start number must be non-negative", "", []string{
				"table_count must be greater than 0",
				"start_number must be greater than or equal to 0",
			}),
		))
	}

	results, err := h.Service.BulkCreateQRCodes(
		request.StoreID,
		request.StartNumber,
		request.TableCount,
		request.Type,
		request.MenuURL,
		request.ExpiresAt,
	)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.Error(
			"Failed to generate bulk QR codes",
			fiber.StatusInternalServerError,
			utils.NewErrorInfo("BULK_GENERATION_ERROR", err.Error(), "", nil),
		))
	}

	metadata := utils.NewPaginationMetadata(1, request.TableCount, len(results))
	return c.Status(fiber.StatusCreated).JSON(utils.Success("QR codes created successfully", results, metadata))
}

// DeleteQRCodeHandler deletes a QR code by its ID
func (h *QRCodeHandler) DeleteQRCodeHandler(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.Error(
			"Invalid QR code ID",
			fiber.StatusBadRequest,
			utils.NewErrorInfo("INVALID_INPUT", "QR code ID must be a valid UUID", "id", nil),
		))
	}

	err = h.Service.DeleteQRCode(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.Error(
			"Failed to delete QR code",
			fiber.StatusInternalServerError,
			utils.NewErrorInfo("DB_ERROR", err.Error(), "", nil),
		))
	}

	return c.Status(fiber.StatusNoContent).JSON(utils.Success("QR code deleted successfully", nil, nil))
}
