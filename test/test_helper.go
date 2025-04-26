package test

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/latoulicious/siresto-backend/internal/domain"
	"github.com/latoulicious/siresto-backend/internal/middleware"
	"github.com/latoulicious/siresto-backend/internal/routes"
	"github.com/latoulicious/siresto-backend/pkg/logger"
	"gorm.io/gorm"
)

// SetupTestApp creates a new Fiber app instance with test configuration
func SetupTestApp(db *gorm.DB) *fiber.App {
	app := fiber.New()

	// Setup CORS for testing - fixed configuration for security
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000", // Specific origin instead of wildcard
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	// Setup routes with test configuration
	routes.SetupRoutes(app, db, logger.NewLogger(nil))

	return app
}

// SetupProtectedRoute creates a test route with JWT protection
func SetupProtectedRoute(app *fiber.App, method, path string, handler fiber.Handler) {
	app.Add(method, path, middleware.Protected(), handler)
}

// CreateTestUser creates a test user in the database using direct SQL
// This has been replaced with direct SQL in the auth_test.go file
func CreateTestUser(db *gorm.DB, user *domain.User) error {
	result := db.Exec(
		"INSERT INTO users (id, name, email, password, is_staff, created_at) VALUES (?, ?, ?, ?, ?, ?)",
		user.ID.String(),
		user.Name,
		user.Email,
		user.Password,
		user.IsStaff,
		user.CreatedAt,
	)
	return result.Error
}

// CleanupTestUser removes a test user from the database
func CleanupTestUser(db *gorm.DB, userID interface{}) error {
	result := db.Exec("DELETE FROM users WHERE id = ?", userID)
	return result.Error
}

// CleanupTestData removes all test data from the database
func CleanupTestData(db *gorm.DB) error {
	// Clean up tables in the correct order to avoid foreign key constraints
	tables := []string{
		"permissions",
		"roles",
		"users",
	}

	for _, table := range tables {
		if err := db.Exec("DELETE FROM " + table).Error; err != nil {
			return err
		}
	}

	return nil
}
