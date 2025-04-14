package routes

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/latoulicious/siresto-backend/internal/handler"
	"github.com/latoulicious/siresto-backend/internal/repository"
	"github.com/latoulicious/siresto-backend/internal/service"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, db *gorm.DB) {
	// Initialize the repository, service, and handler
	qrRepo := &repository.QRCodeRepository{DB: db}
	qrService := &service.QRCodeService{Repo: qrRepo}
	qrHandler := &handler.QRCodeHandler{Service: qrService}

	// Register routes for QR code management
	app.Post("/qr-codes", qrHandler.CreateQRCodeHandler)
	log.Println("POST /qr-codes route registered") // <-- Add logging

	app.Get("/qr-codes/:code", qrHandler.GetQRCodeByCodeHandler)
	log.Println("GET /qr-codes/:code route registered") // <-- Add logging

	app.Put("/qr-codes/:id", qrHandler.UpdateQRCodeHandler)
	log.Println("PUT /qr-codes/:id route registered") // <-- Add logging

	app.Delete("/qr-codes/:id", qrHandler.DeleteQRCodeHandler)
	log.Println("DELETE /qr-codes/:id route registered") // <-- Add logging

	app.Get("/qr-codes/store/:store_id", qrHandler.ListQRCodesHandler)
	log.Println("GET /qr-codes/store/:store_id route registered") // <-- Add logging

	// Add other routes as needed
}
