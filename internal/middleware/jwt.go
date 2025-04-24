package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/latoulicious/siresto-backend/internal/utils"
	"github.com/latoulicious/siresto-backend/pkg/jwt"
)

// Protected is a middleware that checks for a valid JWT token
func Protected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get token from Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(utils.Error("Missing authorization header", fiber.StatusUnauthorized))
		}

		// Check if the header has the Bearer prefix
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(utils.Error("Invalid authorization header format", fiber.StatusUnauthorized))
		}

		// Extract the token
		token := strings.TrimPrefix(authHeader, "Bearer ")

		// Validate token
		claims, err := jwt.ValidateToken(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(utils.Error("Invalid or expired token", fiber.StatusUnauthorized))
		}

		// Store user ID and role information in context for later use
		c.Locals("userID", claims.UserID)
		c.Locals("roleID", claims.RoleID)
		c.Locals("roleName", claims.RoleName)
		c.Locals("isStaff", claims.IsStaff)

		return c.Next()
	}
}

// GetUserID retrieves the authenticated user's ID from the context
func GetUserID(c *fiber.Ctx) (uuid.UUID, bool) {
	userID, ok := c.Locals("userID").(uuid.UUID)
	return userID, ok
}

// GetRoleID retrieves the authenticated user's role ID from the context
func GetRoleID(c *fiber.Ctx) (uuid.UUID, bool) {
	roleID, ok := c.Locals("roleID").(uuid.UUID)
	return roleID, ok
}

// GetRoleName retrieves the authenticated user's role name from the context
func GetRoleName(c *fiber.Ctx) (string, bool) {
	roleName, ok := c.Locals("roleName").(string)
	return roleName, ok
}

// IsStaff checks if the authenticated user is staff
func IsStaff(c *fiber.Ctx) bool {
	isStaff, ok := c.Locals("isStaff").(bool)
	return ok && isStaff
}
