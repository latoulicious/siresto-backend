package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/latoulicious/siresto-backend/internal/config"
	"github.com/latoulicious/siresto-backend/internal/domain"
	"github.com/latoulicious/siresto-backend/pkg/crypto"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Error loading .env file: %v\n", err)
	}

	// Connect to database
	db, err := config.NewGormDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Get all users
	var users []domain.User
	if err := db.Find(&users).Error; err != nil {
		log.Fatal("Failed to fetch users:", err)
	}

	fmt.Printf("Found %d users to update\n", len(users))

	// Ask for confirmation
	fmt.Print("This will update all user passwords to use bcrypt hashing. Continue? (y/n): ")
	var confirm string
	fmt.Scanln(&confirm)
	if confirm != "y" {
		fmt.Println("Operation cancelled")
		os.Exit(0)
	}

	// Update each user's password
	for i, user := range users {
		// Since we can't decrypt the old hash, we'll set a temporary password
		// that users will need to reset
		tempPass := "siresto@123" // You can change this temporary password

		// Hash the temporary password with bcrypt
		hashedPass, err := crypto.HashPassword(tempPass)
		if err != nil {
			log.Printf("Failed to hash password for user %s: %v\n", user.Email, err)
			continue
		}

		// Update the user's password
		if err := db.Model(&user).Update("password", hashedPass).Error; err != nil {
			log.Printf("Failed to update password for user %s: %v\n", user.Email, err)
			continue
		}

		fmt.Printf("Updated password for user %d/%d: %s\n", i+1, len(users), user.Email)
	}

	fmt.Println("\nMigration completed successfully!")
	fmt.Println("All users can now login with the temporary password: siresto@123")
	fmt.Println("Please inform users to change their passwords after logging in.")
}
