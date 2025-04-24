package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/latoulicious/siresto-backend/internal/utils"
)

// RequireRole checks if the user has a specific role
func RequireRole(roleName string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRoleName, ok := GetRoleName(c)
		if !ok || userRoleName != roleName {
			return c.Status(fiber.StatusForbidden).JSON(utils.Error("Insufficient permissions", fiber.StatusForbidden))
		}
		return c.Next()
	}
}

// RequireAdmin checks if the user is an admin (with role "Admin" or is staff)
func RequireAdmin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		roleName, _ := GetRoleName(c)
		if roleName == "System" || IsStaff(c) {
			return c.Next()
		}
		return c.Status(fiber.StatusForbidden).JSON(utils.Error("Admin access required", fiber.StatusForbidden))
	}
}

// RequireStaff checks if the user is staff
func RequireStaff() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if !IsStaff(c) {
			return c.Status(fiber.StatusForbidden).JSON(utils.Error("Staff access required", fiber.StatusForbidden))
		}
		return c.Next()
	}
}
