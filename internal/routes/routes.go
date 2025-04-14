package routes

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/latoulicious/siresto-backend/internal/handler"
	"github.com/latoulicious/siresto-backend/internal/repository"
	"github.com/latoulicious/siresto-backend/internal/service"
	"github.com/latoulicious/siresto-backend/pkg/logger"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, db *gorm.DB) {
	// Initialize the repository, service, and handler for QR codes
	qrRepo := &repository.QRCodeRepository{DB: db}
	qrService := &service.QRCodeService{Repo: qrRepo}
	qrHandler := &handler.QRCodeHandler{Service: qrService}

	// Initialize the repository, service, and handler for logs
	logRepo := &repository.LogRepository{DB: db}
	logService := &service.LogService{Repo: logRepo}
	logHandler := &handler.LogHandler{Service: logService}

	// API v1 group
	v1 := app.Group("/api/v1")

	// Generic routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"code":    fiber.StatusOK,
			"status":  "success",
			"message": "SiResto API is running",
		})
	})
	logger.Log.Info("GET / root route registered")

	// Health check route
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":    "healthy",
			"timestamp": time.Now().Format(time.RFC3339),
			"service":   "siresto-api",
			"version":   "1.0.0",
		})
	})
	logger.Log.Info("GET /health route registered")

	// QR Domain routes (v1)
	v1.Get("/qr-codes", qrHandler.ListAllQRCodesHandler)
	logger.Log.Info("GET /api/v1/qr-codes route registered")

	v1.Get("/qr-codes/:id", qrHandler.GetQRCodeByIDHandler)
	logger.Log.Info("GET /api/v1/qr-codes/:id route registered")

	v1.Get("/qr-codes/store/:store_id", qrHandler.ListQRCodesHandler)
	logger.Log.Info("GET /api/v1/qr-codes/store/:store_id route registered")

	v1.Post("/qr-codes", qrHandler.CreateQRCodeHandler)
	logger.Log.Info("POST /api/v1/qr-codes route registered")

	v1.Put("/qr-codes/:id", qrHandler.UpdateQRCodeHandler)
	logger.Log.Info("PUT /api/v1/qr-codes/:id route registered")

	v1.Delete("/qr-codes/:id", qrHandler.DeleteQRCodeHandler)
	logger.Log.Info("DELETE /api/v1/qr-codes/:id route registered")

	// Log Domain routes (v1)
	v1.Post("/logs", logHandler.CreateLogHandler)
	logger.Log.Info("POST /api/v1/logs route registered")
}
