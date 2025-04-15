package routes

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/latoulicious/siresto-backend/internal/domain"
	"github.com/latoulicious/siresto-backend/internal/handler"
	"github.com/latoulicious/siresto-backend/internal/repository"
	"github.com/latoulicious/siresto-backend/internal/service"
	"github.com/latoulicious/siresto-backend/pkg/core/logging"
	"github.com/latoulicious/siresto-backend/pkg/logutil"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, db *gorm.DB, logger logging.Logger) {
	// QR Code domain
	qrRepo := &repository.QRCodeRepository{DB: db}
	qrService := &service.QRCodeService{Repo: qrRepo}
	qrHandler := &handler.QRCodeHandler{Service: qrService}

	// API v1
	v1 := app.Group("/api/v1")

	// General routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"code":    fiber.StatusOK,
			"status":  "success",
			"message": "SiResto API is running",
		})
	})
	logger.LogInfo("GET / root route registered", logutil.Route("GET", "/"))

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":    "healthy",
			"timestamp": time.Now().Format(time.RFC3339),
			"service":   "siresto-api",
			"version":   "1.0.0",
		})
	})
	logger.LogInfo("GET /health route registered", logutil.Route("GET", "/health"))

	// QR Code routes
	v1.Get("/qr-codes", qrHandler.ListAllQRCodesHandler)
	logger.LogInfo("GET /api/v1/qr-codes route registered", logutil.Route("GET", "/api/v1/qr-codes"))

	v1.Get("/qr-codes/:id", qrHandler.GetQRCodeByIDHandler)
	logger.LogInfo("GET /api/v1/qr-codes/:id route registered", logutil.Route("GET", "/api/v1/qr-codes/:id"))

	v1.Get("/qr-codes/store/:store_id", qrHandler.ListQRCodesHandler)
	logger.LogInfo("GET /api/v1/qr-codes/store/:store_id route registered", logutil.Route("GET", "/api/v1/qr-codes/store/:store_id"))

	v1.Post("/qr-codes", qrHandler.CreateQRCodeHandler)
	logger.LogInfo("POST /api/v1/qr-codes route registered", logutil.Route("POST", "/api/v1/qr-codes"))

	v1.Put("/qr-codes/:id", qrHandler.UpdateQRCodeHandler)
	logger.LogInfo("PUT /api/v1/qr-codes/:id route registered", logutil.Route("PUT", "/api/v1/qr-codes/:id"))

	v1.Delete("/qr-codes/:id", qrHandler.DeleteQRCodeHandler)
	logger.LogInfo("DELETE /api/v1/qr-codes/:id route registered", logutil.Route("DELETE", "/api/v1/qr-codes/:id"))

	// Log routes
	v1.Get("/logs", func(c *fiber.Ctx) error {
		var logs []domain.Log // Assuming your Log model is in the db package
		if err := db.Find(&logs).Error; err != nil {
			logger.LogError("Failed to fetch logs", logutil.Route("GET", "/api/v1/logs"))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to fetch logs",
			})
		}

		// Return the fetched logs
		return c.JSON(fiber.Map{
			"status": "success",
			"logs":   logs,
		})
	})
	logger.LogInfo("GET /api/v1/logs route registered", logutil.Route("GET", "/api/v1/logs"))
}
