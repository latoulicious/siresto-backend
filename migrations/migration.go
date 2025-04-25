package migrations

import (
	"log"

	"github.com/latoulicious/siresto-backend/internal/domain"
	"gorm.io/gorm"
)

// RunMigrations executes all database migrations in the correct order
func RunMigrations(db *gorm.DB) error {
	log.Println("Creating PostgreSQL extensions...")
	// Create extensions (must be done before migrations)
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";").Error; err != nil {
		return err
	}

	// log.Println("Dropping menus table...")
	// Drop the unused table if it exists
	// if err := db.Migrator().DropTable(&domain.Menu{}); err != nil {
	// 	return err
	// }

	log.Println("Running database migrations...")
	// Auto-migrate all domain models
	if err := db.AutoMigrate(
		// User & authorization models
		&domain.Role{},
		&domain.User{},

		// Restaurant models
		&domain.Category{},
		&domain.Product{},
		&domain.Variation{},

		// Order processing models
		&domain.Order{},
		&domain.OrderDetail{},
		&domain.Payment{},
		&domain.Invoice{},

		// Utility models
		&domain.QRCode{},
		&domain.Log{},
		&domain.Theme{},
	); err != nil {
		return err
	}

	// Apply manual migration for setting role positions
	log.Println("Setting up role hierarchy...")
	if err := setupRoleHierarchy(db); err != nil {
		return err
	}

	return nil
}

// setupRoleHierarchy sets the position and isSystem fields for predefined roles
func setupRoleHierarchy(db *gorm.DB) error {
	// System role (highest privilege)
	if err := db.Exec("UPDATE roles SET position = 1, is_system = true WHERE name = 'System'").Error; err != nil {
		return err
	}

	// Owner role (second highest)
	if err := db.Exec("UPDATE roles SET position = 2, is_system = true WHERE name = 'Owner'").Error; err != nil {
		return err
	}

	// Admin role (third highest)
	if err := db.Exec("UPDATE roles SET position = 3, is_system = true WHERE name = 'Admin'").Error; err != nil {
		return err
	}

	// Other standard roles
	if err := db.Exec("UPDATE roles SET position = 10 WHERE name = 'Cashier'").Error; err != nil {
		return err
	}

	if err := db.Exec("UPDATE roles SET position = 11 WHERE name = 'Kitchen'").Error; err != nil {
		return err
	}

	if err := db.Exec("UPDATE roles SET position = 12 WHERE name = 'Waiter'").Error; err != nil {
		return err
	}

	return nil
}

// SeedData populates initial reference data required by the application
func SeedData(db *gorm.DB) error {
	// Default roles
	roles := []domain.Role{
		{Name: "System", Description: "System administrator with full access", Position: 1, IsSystem: true},
		{Name: "Owner", Description: "Business owner with extensive access", Position: 2, IsSystem: true},
		{Name: "Admin", Description: "Administrator with management access", Position: 3, IsSystem: true},
		{Name: "Cashier", Description: "Handles payments and orders", Position: 10, IsSystem: false},
		{Name: "Kitchen", Description: "Kitchen staff", Position: 11, IsSystem: false},
		{Name: "Waiter", Description: "Serving staff", Position: 12, IsSystem: false},
	}

	// Generate UUIDs if not provided and upsert by name
	for _, role := range roles {
		// Use name as the conflict check
		var existingRole domain.Role
		result := db.Where("name = ?", role.Name).First(&existingRole)

		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				// Create new with UUID assigned by Postgres
				if err := db.Create(&role).Error; err != nil {
					return err
				}
			} else {
				return result.Error
			}
		} else {
			// Update position and isSystem if role exists
			existingRole.Position = role.Position
			existingRole.IsSystem = role.IsSystem
			existingRole.Description = role.Description
			if err := db.Save(&existingRole).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

// RunSeeds executes all seed functions to populate the database with initial data
func RunSeeds(db *gorm.DB) error {
	log.Println("Running database seeds...")

	// Seed roles (already called by SeedUsers, but keeping for consistency)
	if err := SeedData(db); err != nil {
		return err
	}

	// Seed users for each role
	if err := SeedUsers(db); err != nil {
		return err
	}

	log.Println("All seeds completed successfully")
	return nil
}
