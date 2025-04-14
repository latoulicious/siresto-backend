package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/latoulicious/siresto-backend/internal/service"
	"github.com/latoulicious/siresto-backend/pkg/logger"
)

type LogHandler struct {
	Service *service.LogService
}

// CreateLogHandler will handle the POST request to store a log
func (h *LogHandler) CreateLogHandler(c *fiber.Ctx) error {
	var logRequest service.CreateLogRequest

	// Parse the incoming log request
	if err := c.BodyParser(&logRequest); err != nil {
		logger.Log.WithError(err).Error("Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Create the log via the service
	if err := h.Service.CreateLog(logRequest); err != nil {
		logger.Log.WithError(err).Error("Failed to store log")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to store log"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success"})
}
