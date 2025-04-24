package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/latoulicious/siresto-backend/internal/domain"
	"github.com/latoulicious/siresto-backend/internal/handler"
	"github.com/latoulicious/siresto-backend/internal/repository"
	"github.com/latoulicious/siresto-backend/internal/service"
	"github.com/latoulicious/siresto-backend/pkg/crypto"
	"github.com/latoulicious/siresto-backend/pkg/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type AuthTestSuite struct {
	suite.Suite
	db          *gorm.DB
	app         *fiber.App
	userHandler *handler.UserHandler
	testUser    *domain.User
	testToken   string
}

func (s *AuthTestSuite) SetupSuite() {
	// Setup in-memory SQLite database for testing
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		s.T().Fatal("Failed to connect to test database:", err)
	}
	s.db = db

	// Modify SQLite to handle UUIDs better (as strings)
	s.db.Exec("PRAGMA foreign_keys = ON")

	// Create simplified test schema for users table with TEXT fields for UUIDs
	s.db.Exec(`CREATE TABLE IF NOT EXISTS "users" (
		"id" TEXT PRIMARY KEY,
		"name" TEXT NOT NULL,
		"email" TEXT NOT NULL UNIQUE,
		"password" TEXT NOT NULL,
		"is_staff" BOOLEAN DEFAULT false,
		"created_at" DATETIME,
		"updated_at" DATETIME,
		"role_id" TEXT
	)`)

	// Create simplified roles table
	s.db.Exec(`CREATE TABLE IF NOT EXISTS "roles" (
		"id" TEXT PRIMARY KEY,
		"name" TEXT NOT NULL UNIQUE,
		"description" TEXT,
		"created_at" DATETIME,
		"updated_at" DATETIME
	)`)

	// Create simplified permissions table
	s.db.Exec(`CREATE TABLE IF NOT EXISTS "permissions" (
		"id" TEXT PRIMARY KEY,
		"name" TEXT NOT NULL UNIQUE,
		"description" TEXT,
		"created_at" DATETIME,
		"updated_at" DATETIME
	)`)

	// Setup repositories and services
	userRepo := &repository.UserRepository{DB: s.db}
	userService := &service.UserService{Repo: userRepo}

	// Setup handlers
	s.userHandler = handler.NewUserHandler(userService, nil)

	// Setup Fiber app with routes
	s.app = SetupTestApp(s.db)

	// Setup login route
	s.app.Post("/api/v1/auth/login", s.userHandler.LoginUser)

	// Setup protected route for testing
	SetupProtectedRoute(s.app, "GET", "/api/v1/users", s.userHandler.ListAllUsers)

	// Create test user
	hashedPassword, err := crypto.HashPassword("testpassword")
	if err != nil {
		s.T().Fatal("Failed to hash password:", err)
	}

	s.testUser = &domain.User{
		ID:        uuid.New(),
		Name:      "Test User",
		Email:     "test@example.com",
		Password:  hashedPassword,
		IsStaff:   true,
		CreatedAt: time.Now(),
	}

	// Insert test user directly using SQL to avoid GORM issues with SQLite and UUIDs
	result := s.db.Exec(
		"INSERT INTO users (id, name, email, password, is_staff, created_at) VALUES (?, ?, ?, ?, ?, ?)",
		s.testUser.ID.String(),
		s.testUser.Name,
		s.testUser.Email,
		s.testUser.Password,
		s.testUser.IsStaff,
		s.testUser.CreatedAt,
	)

	if result.Error != nil {
		s.T().Fatal("Failed to create test user:", result.Error)
	}
}

func (s *AuthTestSuite) TearDownSuite() {
	// Clean up test data
	s.db.Exec("DELETE FROM users")
	s.db.Exec("DELETE FROM roles")
	s.db.Exec("DELETE FROM permissions")

	sqlDB, err := s.db.DB()
	if err != nil {
		s.T().Error("Failed to get database instance:", err)
		return
	}
	if err := sqlDB.Close(); err != nil {
		s.T().Error("Failed to close database:", err)
	}
}

func (s *AuthTestSuite) TestLogin() {
	s.Run("Successful Login", func() {
		loginReq := dto.LoginRequest{
			Email:    "test@example.com",
			Password: "testpassword",
		}

		body, err := json.Marshal(loginReq)
		assert.NoError(s.T(), err)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := s.app.Test(req)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), http.StatusOK, resp.StatusCode)

		var loginResp map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&loginResp)
		assert.NoError(s.T(), err)

		// Store token for protected route tests
		data, ok := loginResp["data"].(map[string]interface{})
		assert.True(s.T(), ok, "Response data should be a map")
		token, ok := data["token"].(string)
		assert.True(s.T(), ok, "Token should be a string")
		s.testToken = token
	})

	s.Run("Invalid Credentials", func() {
		loginReq := dto.LoginRequest{
			Email:    "test@example.com",
			Password: "wrongpassword",
		}

		body, err := json.Marshal(loginReq)
		assert.NoError(s.T(), err)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := s.app.Test(req)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), http.StatusUnauthorized, resp.StatusCode)
	})
}

func (s *AuthTestSuite) TestProtectedRoutes() {
	s.Run("Access Protected Route With Valid Token", func() {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.testToken))

		resp, err := s.app.Test(req)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), http.StatusOK, resp.StatusCode)
	})

	s.Run("Access Protected Route Without Token", func() {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)

		resp, err := s.app.Test(req)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), http.StatusUnauthorized, resp.StatusCode)
	})

	s.Run("Access Protected Route With Invalid Token", func() {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)
		req.Header.Set("Authorization", "Bearer invalid.token.here")

		resp, err := s.app.Test(req)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), http.StatusUnauthorized, resp.StatusCode)
	})
}

func TestAuthSuite(t *testing.T) {
	suite.Run(t, new(AuthTestSuite))
}
