package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/latoulicious/siresto-backend/internal/handler"
	"github.com/latoulicious/siresto-backend/internal/repository"
	"github.com/latoulicious/siresto-backend/internal/service"
	"github.com/latoulicious/siresto-backend/pkg/logger"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, db *gorm.DB) {
	// QR Code domain setup
	qrRepo := &repository.QRCodeRepository{DB: db}
	qrService := &service.QRCodeService{Repo: qrRepo}
	qrHandler := &handler.QRCodeHandler{Service: qrService}

	// Generic routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "SiResto API is running",
		})
	})
	logger.Log.Info("GET / root route registered")

	// Health check route
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})
	logger.Log.Info("GET /health route registered")

	// Domain routes
	app.Post("/qr-codes", qrHandler.CreateQRCodeHandler)
	logger.Log.Info("POST /qr-codes route registered")

	app.Get("/qr-codes/:code", qrHandler.GetQRCodeByCodeHandler)
	logger.Log.Info("GET /qr-codes/:code route registered")

	app.Put("/qr-codes/:id", qrHandler.UpdateQRCodeHandler)
	logger.Log.Info("PUT /qr-codes/:id route registered")

	app.Delete("/qr-codes/:id", qrHandler.DeleteQRCodeHandler)
	logger.Log.Info("DELETE /qr-codes/:id route registered")

	app.Get("/qr-codes/store/:store_id", qrHandler.ListQRCodesHandler)
	logger.Log.Info("GET /qr-codes/store/:store_id route registered")

}
