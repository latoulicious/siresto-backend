package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/latoulicious/siresto-backend/internal/service"
	"github.com/latoulicious/siresto-backend/pkg/logger"
)

type LogHandler struct {
	Service *service.LogService
}

func (h *LogHandler) CreateLogHandler(c *fiber.Ctx) error {
	var logRequest service.CreateLogRequest

	// Parse the incoming log request
	if err := c.BodyParser(&logRequest); err != nil {
		logger.Log.WithError(err).Error("Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request format",
		})
	}

	// Auto-populate context information when available
	ipAddress := c.IP()
	if ipAddress != "" && logRequest.IPAddress == nil {
		logRequest.IPAddress = &ipAddress
	}

	// Extract request ID if available in headers
	requestID := c.Get("X-Request-ID")
	if requestID != "" && logRequest.RequestID == nil {
		logRequest.RequestID = &requestID
	}

	// Extract user ID from context if authenticated
	// (Assuming you store user info in Locals)
	if user := c.Locals("user"); user != nil {
		if userID, ok := user.(string); ok {
			if parsedID, err := uuid.Parse(userID); err == nil {
				logRequest.UserID = &parsedID
			}
		}
	}

	// Create context for the service call
	ctx := c.Context()

	// Create the log via the service
	logEntry, err := h.Service.CreateLog(ctx, logRequest)
	if err != nil {
		logger.Log.WithError(err).Error("Failed to store log")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to store log",
		})
	}

	// Return the created log with its ID
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Log entry created",
		"data": fiber.Map{
			"log_id":    logEntry.ID,
			"timestamp": logEntry.Timestamp,
		},
	})
}
