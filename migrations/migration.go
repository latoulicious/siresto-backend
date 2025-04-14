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

	log.Println("Running database migrations...")
	// Auto-migrate all domain models
	return db.AutoMigrate(
		// User & authorization models
		&domain.Role{},
		&domain.User{},
		&domain.Customer{},

		// Restaurant models
		&domain.Category{},
		&domain.Product{},
		&domain.Variation{},
		&domain.Menu{},

		// Order processing models
		&domain.Order{},
		&domain.OrderDetail{},
		&domain.Payment{},
		&domain.Invoice{},

		// Utility models
		&domain.QRCode{},
	)
}

// SeedData populates initial reference data required by the application
func SeedData(db *gorm.DB) error {
	// Default roles
	roles := []domain.Role{
		{Name: "admin", Description: "System administrator"},
		{Name: "manager", Description: "Restaurant manager"},
		{Name: "staff", Description: "Staff member"},
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
		}
	}

	return nil
}
