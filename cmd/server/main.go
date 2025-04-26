package main

import (
	"context"
	"fmt"
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
	"github.com/latoulicious/siresto-backend/pkg/logutil"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Warning: Error loading .env file: %v\n", err)
		fmt.Println("Continuing with existing environment variables")
	} else {
		fmt.Println("Successfully loaded .env file")
	}

	// Connect to DB
	db, err := config.NewGormDB()
	if err != nil {
		logger.NewLogger(nil).LogError("Failed to connect to DB", logutil.MainCall("connect", "database", map[string]interface{}{
			"error": err.Error(),
		}))
		panic("Failed to connect to DB: " + err.Error())
	}
	sqlDB, err := db.DB()
	if err != nil {
		logger.NewLogger(nil).LogError("Failed to get sql.DB from GORM", logutil.MainCall("init", "database", map[string]interface{}{
			"error": err.Error(),
		}))
		panic("Failed to get sql.DB from GORM: " + err.Error())
	}
	defer sqlDB.Close()

	// Init logger persister and logger
	persister := logger.NewLogServicePersister(db)
	appLogger := logger.NewLogger(persister)

	appLogger.LogInfo("Connected to DB successfully", logutil.MainCall("connect", "database", nil))

	// Setup Fiber app
	app := fiber.New()

	// Load allowed origins from env
	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
	if allowedOrigins == "" {
		allowedOrigins = "http://localhost:3000"
	}

	appLogger.LogInfo("CORS initialized", logutil.MainCall("init", "cors", map[string]interface{}{
		"allowed_origins": allowedOrigins,
	}))

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
		appLogger.LogInfo("Server starting...", logutil.MainCall("start", "server", map[string]interface{}{
			"port": port,
		}))
		if err := app.Listen(":" + port); err != nil {
			appLogger.LogError("Failed to start server", logutil.MainCall("start", "server", map[string]interface{}{
				"error": err.Error(),
			}))
			os.Exit(1)
		}
	}()

	// Graceful shutdown on interrupt signals
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	appLogger.LogInfo("Gracefully shutting down server...", logutil.MainCall("shutdown", "server", nil))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		appLogger.LogError("Shutdown error", logutil.MainCall("shutdown", "server", map[string]interface{}{
			"error": err.Error(),
		}))
	} else {
		appLogger.LogInfo("Server shut down successfully", logutil.MainCall("shutdown", "server", nil))
	}
}
