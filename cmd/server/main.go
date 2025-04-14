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
	// Auto-load .env
	if err := godotenv.Load(); err != nil {
		logger.Log.Warn("No .env file found or failed to load, continuing...")
	}

	// Initialize the logger (only if you made InitLogger separate)
	logger.InitLogger()

	// Log info on start
	logger.Log.Info("Starting the server...")

	// Connect to database
	db, err := config.NewGormDB()
	if err != nil {
		logger.Log.WithError(err).Fatal("Failed to connect to database")
	}

	sqlDB, err := db.DB()
	if err != nil {
		logger.Log.WithError(err).Fatal("Failed to get sql.DB from GORM")
	}
	defer func() {
		if err := sqlDB.Close(); err != nil {
			logger.Log.WithError(err).Warn("Failed to close database connection cleanly")
		}
	}()
	logger.Log.Info("Connected to Vercel Postgres using GORM")

	// Initialize Fiber
	app := fiber.New()

	// Setup routes
	routes.SetupRoutes(app, db)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	logger.Log.Info("Server starting on port ", port)
	if err := app.Listen(":" + port); err != nil {
		logger.Log.WithError(err).Fatal("Failed to start server")
	}
}
