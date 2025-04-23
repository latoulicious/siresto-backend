package service

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/latoulicious/siresto-backend/internal/domain"
	"github.com/latoulicious/siresto-backend/internal/repository"
	"github.com/latoulicious/siresto-backend/internal/utils"
	"github.com/latoulicious/siresto-backend/pkg/dto"
)

// Common errors
var (
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type UserService struct {
	Repo *repository.UserRepository
}

// ListAllUsers returns all users as DTOs
func (s *UserService) ListAllUsers() ([]dto.UserResponse, error) {
	users, err := s.Repo.ListAllUsers()
	if err != nil {
		return nil, err
	}

	// Map domain users to DTOs
	userDTOs := make([]dto.UserResponse, len(users))
	for i, user := range users {
		userDTOs[i] = mapToUserResponse(&user)
	}

	return userDTOs, nil
}

// GetUserByID returns a user by ID as DTO
func (s *UserService) GetUserByID(id uuid.UUID) (*dto.UserResponse, error) {
	user, err := s.Repo.GetUserByID(id)
	if err != nil {
		return nil, ErrUserNotFound
	}

	userDTO := mapToUserResponse(user)
	return &userDTO, nil
}

// CreateUser creates a new user from DTO
func (s *UserService) CreateUser(req *dto.CreateUserRequest) (*dto.UserResponse, error) {
	// Check if email already exists
	existingUser, err := s.Repo.FindByEmail(req.Email)
	if err == nil && existingUser != nil {
		return nil, ErrEmailAlreadyExists
	}

	hashedPassword := utils.HashSHA256(req.Password)

	// Map DTO to domain
	user := &domain.User{
		ID:        uuid.New(),
		Name:      req.Name,
		Email:     req.Email,
		Password:  hashedPassword,
		IsStaff:   req.IsStaff,
		RoleID:    req.RoleID,
		CreatedAt: time.Now(),
	}

	createdUser, err := s.Repo.CreateUser(user)
	if err != nil {
		return nil, err
	}

	// Map back to DTO
	userDTO := mapToUserResponse(createdUser)
	return &userDTO, nil
}

// UpdateUser updates a user from DTO
func (s *UserService) UpdateUser(id uuid.UUID, req *dto.UpdateUserRequest) (*dto.UserResponse, error) {
	// Get existing user
	existingUser, err := s.Repo.GetUserByID(id)
	if err != nil {
		return nil, ErrUserNotFound
	}

	// Check if email is being changed and already exists
	if req.Email != nil && *req.Email != existingUser.Email {
		userWithEmail, err := s.Repo.FindByEmail(*req.Email)
		if err == nil && userWithEmail != nil {
			return nil, ErrEmailAlreadyExists
		}
	}

	// Apply updates from DTO to domain model
	if req.Name != nil {
		existingUser.Name = *req.Name
	}

	if req.Email != nil {
		existingUser.Email = *req.Email
	}

	if req.Password != nil {
		existingUser.Password = utils.HashSHA256(*req.Password)
	}

	if req.IsStaff != nil {
		existingUser.IsStaff = *req.IsStaff
	}

	if req.RoleID != nil {
		existingUser.RoleID = *req.RoleID
	}

	// Update user in repository
	updatedUser, err := s.Repo.UpdateUser(existingUser)
	if err != nil {
		return nil, err
	}

	// Map to DTO
	userDTO := mapToUserResponse(updatedUser)
	return &userDTO, nil
}

// DeleteUser deletes a user by ID
func (s *UserService) DeleteUser(id uuid.UUID) error {
	// Check if user exists
	_, err := s.Repo.GetUserByID(id)
	if err != nil {
		return ErrUserNotFound
	}

	return s.Repo.DeleteUser(id)
}

// LoginUser authenticates a user and returns user info with token
func (s *UserService) LoginUser(req *dto.LoginRequest) (*dto.UserLoginResponse, error) {
	user, err := s.Repo.FindByEmail(req.Email)
	if err != nil || user == nil {
		return nil, ErrInvalidCredentials
	}

	hashedInput := utils.HashSHA256(req.Password)
	if user.Password != hashedInput {
		return nil, ErrInvalidCredentials
	}

	// Update last login
	now := time.Now()
	user.LastLoginAt = &now
	if err := s.Repo.UpdateLastLogin(user.ID, now); err != nil {
		// Log the error but continue - non-critical
		return nil, err
	}

	// Get user role if available
	// var roleInfo dto.RoleInfo
	// if user.RoleID != uuid.Nil {
	// 	// Here you would fetch role info from repository
	// 	// For now, assuming a dummy implementation
	// 	roleInfo = dto.RoleInfo{
	// 		ID:   user.RoleID,
	// 		Name: "Dummy Role", // In a real scenario, fetch this from the database
	// 	}
	// }

	// Create user response
	userResponse := mapToUserResponse(user)

	// Generate JWT token (implementation depends on your auth strategy)
	token := generateToken(user.ID)

	loginResponse := &dto.UserLoginResponse{
		User:  userResponse,
		Token: token,
	}

	return loginResponse, nil
}

// Helper functions

// mapToUserResponse maps a domain user to a DTO
func mapToUserResponse(user *domain.User) dto.UserResponse {
	role := dto.RoleInfo{}
	if user.Role != nil {
		role.ID = user.Role.ID
		role.Name = user.Role.Name
	}

	return dto.UserResponse{
		ID:          user.ID,
		Name:        user.Name,
		Email:       user.Email,
		IsStaff:     user.IsStaff,
		CreatedAt:   user.CreatedAt,
		LastLoginAt: user.LastLoginAt,
		Role:        role,
	}
}

// generateToken generates a JWT token for the user
func generateToken(userID uuid.UUID) string {
	// Implement your JWT token generation here
	// This is a placeholder implementation
	return "jwt-token-placeholder-" + userID.String()
}
