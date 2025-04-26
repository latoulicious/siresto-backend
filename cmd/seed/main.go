package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/latoulicious/siresto-backend/internal/config"
	"github.com/latoulicious/siresto-backend/migrations"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Error loading .env file: %v\n", err)
	}

	// Initialize database connection using the project's config
	db, err := config.NewGormDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
		os.Exit(1)
	}

	// Run seeds
	if err := migrations.RunSeeds(db); err != nil {
		log.Fatalf("Failed to seed database: %v", err)
		os.Exit(1)
	}

	log.Println("Database seeding completed successfully!")
}
