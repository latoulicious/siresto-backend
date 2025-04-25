package dto

import (
	"github.com/google/uuid"
)

// Request DTOs
type CreateRoleRequest struct {
	Name        string      `json:"name" validate:"required"`
	Description string      `json:"description"`
	Position    int         `json:"position"`
	Permissions []uuid.UUID `json:"permissions"`
}

type UpdateRoleRequest struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Position    *int        `json:"position,omitempty"`
	Permissions []uuid.UUID `json:"permissions"`
}

// New DTOs for permission operations
type RolePermissionUpdateRequest struct {
	Permissions []uuid.UUID `json:"permissions" validate:"required,min=1"`
}

// Response DTOs
type PermissionResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

type RoleResponse struct {
	ID          uuid.UUID            `json:"id"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Position    int                  `json:"position"`
	IsSystem    bool                 `json:"is_system"`
	Permissions []PermissionResponse `json:"permissions"`
}
