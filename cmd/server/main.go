package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/latoulicious/siresto-backend/internal/config"
)

func main() {
	// Auto-load .env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or failed to load, continuing...")
	}

	// Connect to database
	db, err := config.NewGormDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get sql.DB from GORM: %v", err)
	}
	defer sqlDB.Close()
	log.Println("Connected to Vercel Postgres using GORM!")

	// Initialize Fiber
	app := fiber.New()

	// Root endpoint
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "SiResto API is running",
		})
	})

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Server starting on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
