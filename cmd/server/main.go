package main

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/latoulicious/siresto-backend/internal/config"
	"github.com/latoulicious/siresto-backend/internal/routes"
	"github.com/latoulicious/siresto-backend/pkg/logger"
)

func main() {
	_ = godotenv.Load()

	// Connect to DB
	db, err := config.NewGormDB()
	if err != nil {
		panic("Failed to connect to DB: " + err.Error())
	}
	sqlDB, err := db.DB()
	if err != nil {
		panic("Failed to get sql.DB from GORM: " + err.Error())
	}
	defer sqlDB.Close()

	// Init logger persister and logger
	persister := logger.NewLogServicePersister(db)
	appLogger := logger.NewLogger(persister)

	// Log example
	appLogger.LogInfo("Connected to DB successfully", nil)
	appLogger.LogInfo("Starting the server...", nil)

	// Setup Fiber app
	app := fiber.New()

	// Setup routes (pass logger)
	routes.SetupRoutes(app, db, appLogger)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	appLogger.LogInfo("Server listening on port "+port, nil)
	if err := app.Listen(":" + port); err != nil {
		appLogger.LogError("Failed to start server", map[string]interface{}{
			"error": err.Error(),
		})
	}
}
