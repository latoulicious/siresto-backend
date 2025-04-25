package validator

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/latoulicious/siresto-backend/internal/domain"
	"github.com/latoulicious/siresto-backend/internal/middleware"
)

var (
	// Standard action prefixes
	standardActions = []string{
		middleware.PermissionRead,
		middleware.PermissionCreate,
		middleware.PermissionUpdate,
		middleware.PermissionDelete,
		"manage",  // For management operations
		"approve", // For approval workflows
		"reject",  // For rejection workflows
		"export",  // For export operations
		"import",  // For import operations
		"assign",  // For assignment operations
		"view",    // Alternative to read
		"list",    // For listing operations
		"access",  // For general access
		"full",    // For full access
	}

	// Special prefixes for more complex permissions
	specialPermissionPrefixes = []string{
		"manage:",
		"access:",
		"full:",
	}

	// Valid permission format: action:resource or special:permission
	permissionFormatRegex = regexp.MustCompile(`^([a-z]+):([a-z_]+)$`)
)

// ValidatePermission ensures that a permission follows the established patterns
func ValidatePermission(permission *domain.Permission) error {
	// Check for empty name
	if permission.Name == "" {
		return fmt.Errorf("permission name cannot be empty")
	}

	// Check for special permissions
	for _, prefix := range specialPermissionPrefixes {
		if strings.HasPrefix(permission.Name, prefix) {
			return nil // Special permissions are allowed
		}
	}

	// Validate against permission format regex
	if !permissionFormatRegex.MatchString(permission.Name) {
		return fmt.Errorf("permission must follow format 'action:resource' (e.g., 'read:user')")
	}

	// Extract action and resource parts
	parts := strings.Split(permission.Name, ":")
	if len(parts) != 2 {
		return fmt.Errorf("permission must have exactly one colon separator")
	}

	action := parts[0]

	// Validate action against standard list
	validAction := false
	for _, standardAction := range standardActions {
		if action == standardAction {
			validAction = true
			break
		}
	}

	if !validAction {
		return fmt.Errorf("invalid action '%s', must be one of: %s", action, strings.Join(standardActions, ", "))
	}

	return nil
}

// GetPermissionDescription generates a standardized description for a permission
func GetPermissionDescription(permissionName string) string {
	// Check for special permissions
	for _, prefix := range specialPermissionPrefixes {
		if strings.HasPrefix(permissionName, prefix) {
			resource := strings.TrimPrefix(permissionName, prefix)
			return fmt.Sprintf("Special permission to %s %s", prefix[:len(prefix)-1], resource)
		}
	}

	// Handle standard permissions
	parts := strings.Split(permissionName, ":")
	if len(parts) != 2 {
		return permissionName // Return as is if not in expected format
	}

	action := parts[0]
	resource := parts[1]

	// Create a user-friendly description
	return fmt.Sprintf("Permission to %s %s", action, strings.ReplaceAll(resource, "_", " "))
}
