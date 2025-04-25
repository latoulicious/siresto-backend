package utils

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/latoulicious/siresto-backend/internal/domain"
)

// Permission constants for common operations
const (
	// Resource access levels
	PermissionRead   = "read"
	PermissionCreate = "create"
	PermissionUpdate = "update"
	PermissionDelete = "delete"
)

// GenerateCRUDPermissions creates a set of standard CRUD permissions for a resource
func GenerateCRUDPermissions(resourceName string) []domain.Permission {
	// Ensure resource name is valid (lowercase, no spaces)
	resourceName = strings.ToLower(strings.ReplaceAll(resourceName, " ", "_"))

	// Create the four standard CRUD permissions
	permissions := []domain.Permission{
		{
			ID:          uuid.New(),
			Name:        fmt.Sprintf("%s:%s", PermissionRead, resourceName),
			Description: fmt.Sprintf("Permission to read %s", strings.ReplaceAll(resourceName, "_", " ")),
		},
		{
			ID:          uuid.New(),
			Name:        fmt.Sprintf("%s:%s", PermissionCreate, resourceName),
			Description: fmt.Sprintf("Permission to create %s", strings.ReplaceAll(resourceName, "_", " ")),
		},
		{
			ID:          uuid.New(),
			Name:        fmt.Sprintf("%s:%s", PermissionUpdate, resourceName),
			Description: fmt.Sprintf("Permission to update %s", strings.ReplaceAll(resourceName, "_", " ")),
		},
		{
			ID:          uuid.New(),
			Name:        fmt.Sprintf("%s:%s", PermissionDelete, resourceName),
			Description: fmt.Sprintf("Permission to delete %s", strings.ReplaceAll(resourceName, "_", " ")),
		},
	}

	return permissions
}

// GenerateManagementPermission creates a management permission for a resource
func GenerateManagementPermission(resourceName string) domain.Permission {
	// Ensure resource name is valid (lowercase, no spaces)
	resourceName = strings.ToLower(strings.ReplaceAll(resourceName, " ", "_"))

	return domain.Permission{
		ID:          uuid.New(),
		Name:        fmt.Sprintf("manage:%s", resourceName),
		Description: fmt.Sprintf("Permission to manage all aspects of %s", strings.ReplaceAll(resourceName, "_", " ")),
	}
}

// GeneratePermissionBundle creates a complete set of permissions for a resource (CRUD + management)
func GeneratePermissionBundle(resourceName string) []domain.Permission {
	permissions := GenerateCRUDPermissions(resourceName)
	permissions = append(permissions, GenerateManagementPermission(resourceName))

	return permissions
}

// ValidateAndFormatResource ensures a resource name is valid for permission creation
func ValidateAndFormatResource(resourceName string) (string, error) {
	// Trim whitespace
	resourceName = strings.TrimSpace(resourceName)

	// Check for empty string
	if resourceName == "" {
		return "", fmt.Errorf("resource name cannot be empty")
	}

	// Convert to lowercase and replace spaces with underscores
	formatted := strings.ToLower(strings.ReplaceAll(resourceName, " ", "_"))

	// Validate against a regex pattern - only allow alphanumeric and underscores
	validResource := regexp.MustCompile(`^[a-z0-9_]+$`)
	if !validResource.MatchString(formatted) {
		return "", fmt.Errorf("resource name must contain only letters, numbers, and underscores")
	}

	return formatted, nil
}
