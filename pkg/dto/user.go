package dto

import (
	"time"

	"github.com/google/uuid"
)

// Request DTOs
type CreateUserRequest struct {
	Name     string    `json:"name" validate:"required"`
	Email    string    `json:"email" validate:"required,email"`
	Password string    `json:"password" validate:"required,min=8"`
	IsStaff  bool      `json:"is_staff"`
	RoleID   uuid.UUID `json:"role_id"`
}

type UpdateUserRequest struct {
	Name     *string    `json:"name"`
	Email    *string    `json:"email" validate:"omitempty,email"`
	Password *string    `json:"password" validate:"omitempty,min=8"`
	IsStaff  *bool      `json:"is_staff"`
	RoleID   *uuid.UUID `json:"role_id"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// Response DTOs
type UserResponse struct {
	ID          uuid.UUID  `json:"id"`
	Name        string     `json:"name"`
	Email       string     `json:"email"`
	IsStaff     bool       `json:"is_staff"`
	Role        RoleInfo   `json:"role,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	LastLoginAt *time.Time `json:"last_login_at,omitempty"`
}

type UserLoginResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token,omitempty"`
}

type RoleInfo struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
}
