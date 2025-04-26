package migrations

import (
	"log"

	"github.com/google/uuid"
	"github.com/latoulicious/siresto-backend/internal/domain"
	"github.com/latoulicious/siresto-backend/pkg/crypto"
	"gorm.io/gorm"
)

// SeedUsers creates default users for each role in the system
func SeedUsers(db *gorm.DB) error {
	log.Println("Seeding users for each role...")

	// First, ensure roles exist
	if err := SeedData(db); err != nil {
		return err
	}

	// Get role IDs
	roles := make(map[string]uuid.UUID)
	var rolesList []domain.Role
	if err := db.Find(&rolesList).Error; err != nil {
		return err
	}

	for _, role := range rolesList {
		roles[role.Name] = role.ID
	}

	// Default password for test accounts
	defaultPassword, err := crypto.HashPassword("password123")
	if err != nil {
		return err
	}

	// Define users for each role
	users := []domain.User{
		{
			Name:     "Victoria",
			Email:    "system@siresto.com",
			Password: defaultPassword,
			RoleID:   roles["System"],
			IsStaff:  true,
		},
		{
			Name:     "Owner User",
			Email:    "owner@siresto.com",
			Password: defaultPassword,
			RoleID:   roles["Owner"],
			IsStaff:  true,
		},
		{
			Name:     "Admin User",
			Email:    "admin@siresto.com",
			Password: defaultPassword,
			RoleID:   roles["Admin"],
			IsStaff:  true,
		},
		{
			Name:     "Cashier User",
			Email:    "cashier@siresto.com",
			Password: defaultPassword,
			RoleID:   roles["Cashier"],
			IsStaff:  true,
		},
		{
			Name:     "Kitchen User",
			Email:    "kitchen@siresto.com",
			Password: defaultPassword,
			RoleID:   roles["Kitchen"],
			IsStaff:  true,
		},
		{
			Name:     "Waiter User",
			Email:    "waiter@siresto.com",
			Password: defaultPassword,
			RoleID:   roles["Waiter"],
			IsStaff:  true,
		},
	}

	// Create users if they don't exist
	for _, user := range users {
		var existingUser domain.User
		result := db.Where("email = ?", user.Email).First(&existingUser)

		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				// Use a local variable to get the role name for logging
				var roleName string
				for name, id := range roles {
					if id == user.RoleID {
						roleName = name
						break
					}
				}

				if err := db.Create(&user).Error; err != nil {
					return err
				}
				log.Printf("Created user: %s with role: %s", user.Email, roleName)
			} else {
				return result.Error
			}
		} else {
			log.Printf("User already exists: %s", user.Email)
		}
	}

	log.Println("User seeding completed successfully")
	return nil
}
