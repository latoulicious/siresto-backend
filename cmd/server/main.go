package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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

	// Setup Fiber app
	app := fiber.New()

	// Load allowed origins from env
	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
	if allowedOrigins == "" {
		allowedOrigins = "http://localhost:3000"
	}

	// Setup CORS middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	// Setup routes
	routes.SetupRoutes(app, db, appLogger)

	// Get PORT
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// Start Fiber in goroutine
	go func() {
		appLogger.LogInfo("Server listening on port "+port, nil)
		if err := app.Listen(":" + port); err != nil {
			appLogger.LogError("Failed to start server", map[string]interface{}{
				"error": err.Error(),
			})
			os.Exit(1)
		}
	}()

	// Graceful shutdown on interrupt signals
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	appLogger.LogInfo("Gracefully shutting down server...", nil)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		appLogger.LogError("Shutdown error", map[string]interface{}{
			"error": err.Error(),
		})
	} else {
		appLogger.LogInfo("Server shut down successfully", nil)
	}
}
