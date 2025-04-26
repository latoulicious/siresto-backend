package middleware

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/latoulicious/siresto-backend/internal/utils"
)

// Role constants for predefined system roles
const (
	RoleSystem  = "System"
	RoleOwner   = "Owner"
	RoleAdmin   = "Admin"
	RoleCashier = "Cashier"
	RoleKitchen = "Kitchen"
	RoleWaiter  = "Waiter"
)

// Permission constants for common operations
const (
	// Resource access levels
	PermissionRead   = "read"
	PermissionCreate = "create"
	PermissionUpdate = "update"
	PermissionDelete = "delete"

	// Prefixes for resource types
	ResourceUser       = "user"
	ResourceRole       = "role"
	ResourcePermission = "permission"
	ResourceMenu       = "menu"
	ResourceOrder      = "order"
	ResourceTable      = "table"
	ResourceInventory  = "inventory"
	ResourceReport     = "report"
	ResourceSetting    = "setting"

	// Special permissions
	PermissionManageUsers       = "manage:users"
	PermissionManageRoles       = "manage:roles"
	PermissionManagePermissions = "manage:permissions"
	PermissionAccessSystem      = "access:system"
	PermissionFullAccess        = "full:access"
)

// Helper function to format permission strings in a consistent way
func FormatPermission(action, resource string) string {
	return fmt.Sprintf("%s:%s", action, resource)
}

// GetDefaultPermissionsForRole returns the default set of permissions for a standard role
func GetDefaultPermissionsForRole(roleName string) []string {
	switch roleName {
	case RoleSystem:
		return []string{PermissionFullAccess}

	case RoleOwner:
		return []string{
			PermissionManageUsers,
			PermissionManageRoles,
			PermissionManagePermissions,
			PermissionAccessSystem,
			FormatPermission(PermissionRead, ResourceReport),
			FormatPermission(PermissionCreate, ResourceReport),
			FormatPermission(PermissionRead, ResourceUser),
			FormatPermission(PermissionCreate, ResourceUser),
			FormatPermission(PermissionUpdate, ResourceUser),
			FormatPermission(PermissionDelete, ResourceUser),
			FormatPermission(PermissionRead, ResourceRole),
			FormatPermission(PermissionCreate, ResourceRole),
			FormatPermission(PermissionUpdate, ResourceRole),
			FormatPermission(PermissionDelete, ResourceRole),
			FormatPermission(PermissionRead, ResourceMenu),
			FormatPermission(PermissionCreate, ResourceMenu),
			FormatPermission(PermissionUpdate, ResourceMenu),
			FormatPermission(PermissionDelete, ResourceMenu),
			FormatPermission(PermissionRead, ResourceOrder),
			FormatPermission(PermissionCreate, ResourceOrder),
			FormatPermission(PermissionUpdate, ResourceOrder),
			FormatPermission(PermissionDelete, ResourceOrder),
			FormatPermission(PermissionRead, ResourceTable),
			FormatPermission(PermissionCreate, ResourceTable),
			FormatPermission(PermissionUpdate, ResourceTable),
			FormatPermission(PermissionDelete, ResourceTable),
			FormatPermission(PermissionRead, ResourceInventory),
			FormatPermission(PermissionCreate, ResourceInventory),
			FormatPermission(PermissionUpdate, ResourceInventory),
			FormatPermission(PermissionDelete, ResourceInventory),
			FormatPermission(PermissionRead, ResourceSetting),
			FormatPermission(PermissionUpdate, ResourceSetting),
		}

	case RoleAdmin:
		return []string{
			PermissionManageUsers,
			FormatPermission(PermissionRead, ResourceReport),
			FormatPermission(PermissionCreate, ResourceReport),
			FormatPermission(PermissionRead, ResourceUser),
			FormatPermission(PermissionCreate, ResourceUser),
			FormatPermission(PermissionUpdate, ResourceUser),
			FormatPermission(PermissionRead, ResourceMenu),
			FormatPermission(PermissionCreate, ResourceMenu),
			FormatPermission(PermissionUpdate, ResourceMenu),
			FormatPermission(PermissionDelete, ResourceMenu),
			FormatPermission(PermissionRead, ResourceOrder),
			FormatPermission(PermissionCreate, ResourceOrder),
			FormatPermission(PermissionUpdate, ResourceOrder),
			FormatPermission(PermissionDelete, ResourceOrder),
			FormatPermission(PermissionRead, ResourceTable),
			FormatPermission(PermissionCreate, ResourceTable),
			FormatPermission(PermissionUpdate, ResourceTable),
			FormatPermission(PermissionDelete, ResourceTable),
			FormatPermission(PermissionRead, ResourceInventory),
			FormatPermission(PermissionCreate, ResourceInventory),
			FormatPermission(PermissionUpdate, ResourceInventory),
			FormatPermission(PermissionDelete, ResourceInventory),
			FormatPermission(PermissionRead, ResourceSetting),
		}

	case RoleCashier:
		return []string{
			FormatPermission(PermissionRead, ResourceMenu),
			FormatPermission(PermissionRead, ResourceOrder),
			FormatPermission(PermissionCreate, ResourceOrder),
			FormatPermission(PermissionUpdate, ResourceOrder),
			FormatPermission(PermissionRead, ResourceTable),
			FormatPermission(PermissionUpdate, ResourceTable),
		}

	case RoleKitchen:
		return []string{
			FormatPermission(PermissionRead, ResourceMenu),
			FormatPermission(PermissionRead, ResourceOrder),
			FormatPermission(PermissionUpdate, ResourceOrder),
			FormatPermission(PermissionRead, ResourceInventory),
			FormatPermission(PermissionUpdate, ResourceInventory),
		}

	case RoleWaiter:
		return []string{
			FormatPermission(PermissionRead, ResourceMenu),
			FormatPermission(PermissionRead, ResourceOrder),
			FormatPermission(PermissionUpdate, ResourceOrder),
			FormatPermission(PermissionRead, ResourceTable),
			FormatPermission(PermissionUpdate, ResourceTable),
		}

	default:
		// For custom roles, return an empty slice
		// Their permissions should be explicitly set
		return []string{}
	}
}

// HasPermission checks if a list of permissions includes a specific permission
func HasPermission(userPermissions []string, requiredPermission string) bool {
	// Full access permission grants access to everything
	for _, p := range userPermissions {
		if p == PermissionFullAccess {
			return true
		}
	}

	// Check for exact match
	for _, p := range userPermissions {
		if p == requiredPermission {
			return true
		}
	}

	return false
}

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

// RequireOneOfRoles checks if the user has at least one of the specified roles
func RequireOneOfRoles(roleNames ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRoleName, ok := GetRoleName(c)
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(utils.Error("Insufficient permissions", fiber.StatusForbidden))
		}

		// System role has implicit access to everything
		if userRoleName == RoleSystem {
			return c.Next()
		}

		for _, role := range roleNames {
			if userRoleName == role {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(utils.Error("Insufficient permissions", fiber.StatusForbidden))
	}
}

// RequireAdmin checks if the user is an admin (with role "System", "Owner", or "Admin")
func RequireAdmin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		roleName, ok := GetRoleName(c)
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(utils.Error("Admin access required", fiber.StatusForbidden))
		}

		if roleName == RoleSystem || roleName == RoleOwner || roleName == RoleAdmin {
			return c.Next()
		}

		return c.Status(fiber.StatusForbidden).JSON(utils.Error("Admin access required", fiber.StatusForbidden))
	}
}

// RequireManagement checks if the user has management privileges (System, Owner, Admin)
func RequireManagement() fiber.Handler {
	return RequireOneOfRoles(RoleSystem, RoleOwner, RoleAdmin)
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

// RequireRestaurantStaff checks if user is part of restaurant staff (Cashier, Kitchen, Waiter)
func RequireRestaurantStaff() fiber.Handler {
	return RequireOneOfRoles(RoleCashier, RoleKitchen, RoleWaiter)
}

// HasRolePrefix checks if the user's role starts with a specific prefix
// Useful for dynamic role grouping (e.g., "Kitchen_Level1", "Kitchen_Level2")
func HasRolePrefix(prefix string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRoleName, ok := GetRoleName(c)
		if !ok || !strings.HasPrefix(userRoleName, prefix) {
			return c.Status(fiber.StatusForbidden).JSON(utils.Error("Insufficient permissions", fiber.StatusForbidden))
		}
		return c.Next()
	}
}

// RequirePermission checks if the user's role has a specific permission
func RequirePermission(permissionName string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// System, Owner and Admin roles implicitly have all permissions
		roleName, ok := GetRoleName(c)
		if ok && (roleName == RoleSystem || roleName == RoleOwner || roleName == RoleAdmin) {
			return c.Next()
		}

		// For other roles, check if user has specific permission
		userPermissions, ok := GetUserPermissions(c)
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(utils.Error("Insufficient permissions", fiber.StatusForbidden))
		}

		// Check if the permission exists in the user's permissions
		for _, p := range userPermissions {
			if p == permissionName {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(utils.Error("Insufficient permissions", fiber.StatusForbidden))
	}
}

// RequireAnyPermission checks if the user's role has at least one of the specified permissions
func RequireAnyPermission(permissionNames ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// System, Owner and Admin roles implicitly have all permissions
		roleName, ok := GetRoleName(c)
		if ok && (roleName == RoleSystem || roleName == RoleOwner || roleName == RoleAdmin) {
			return c.Next()
		}

		// For other roles, check if user has any of the specified permissions
		userPermissions, ok := GetUserPermissions(c)
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(utils.Error("Insufficient permissions", fiber.StatusForbidden))
		}

		// Convert user permissions to a map for faster lookup
		permMap := make(map[string]bool)
		for _, p := range userPermissions {
			permMap[p] = true
		}

		// Check if any required permission exists
		for _, required := range permissionNames {
			if permMap[required] {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(utils.Error("Insufficient permissions", fiber.StatusForbidden))
	}
}

// RequireAllPermissions checks if the user's role has all specified permissions
func RequireAllPermissions(permissionNames ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// System, Owner and Admin roles implicitly have all permissions
		roleName, ok := GetRoleName(c)
		if ok && (roleName == RoleSystem || roleName == RoleOwner || roleName == RoleAdmin) {
			return c.Next()
		}

		// For other roles, check if user has all the specified permissions
		userPermissions, ok := GetUserPermissions(c)
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(utils.Error("Insufficient permissions", fiber.StatusForbidden))
		}

		// Convert user permissions to a map for faster lookup
		permMap := make(map[string]bool)
		for _, p := range userPermissions {
			permMap[p] = true
		}

		// Check if all required permissions exist
		for _, required := range permissionNames {
			if !permMap[required] {
				return c.Status(fiber.StatusForbidden).JSON(utils.Error("Insufficient permissions", fiber.StatusForbidden))
			}
		}

		return c.Next()
	}
}

// RequireResourcePermission is a helper for common CRUD operations on resources
func RequireResourcePermission(action, resource string) fiber.Handler {
	permission := FormatPermission(action, resource)
	return RequirePermission(permission)
}

// GetUserPermissions retrieves the authenticated user's permissions from context
func GetUserPermissions(c *fiber.Ctx) ([]string, bool) {
	perms, ok := c.Locals("permissions").([]string)
	return perms, ok
}

// IsSameUserOrHigherRole ensures the authenticated user is either:
// 1. The same user being accessed (based on ID)
// 2. Has a higher privilege role (System, Owner, Admin)
func IsSameUserOrHigherRole() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get target user ID from params
		targetUserID := c.Params("id")
		if targetUserID == "" {
			targetUserID = c.Query("id")
		}

		// Get authenticated user
		userID, ok := GetUserID(c)
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(utils.Error("Authentication required", fiber.StatusForbidden))
		}

		// If same user, allow access
		if userID.String() == targetUserID {
			return c.Next()
		}

		// Otherwise, check if user has admin privileges
		roleName, ok := GetRoleName(c)
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(utils.Error("Insufficient permissions", fiber.StatusForbidden))
		}

		// Allow if user has higher-level role
		if roleName == RoleSystem || roleName == RoleOwner || roleName == RoleAdmin {
			return c.Next()
		}

		return c.Status(fiber.StatusForbidden).JSON(utils.Error("Insufficient permissions", fiber.StatusForbidden))
	}
}
