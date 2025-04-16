package service

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/latoulicious/siresto-backend/internal/domain"
	"github.com/latoulicious/siresto-backend/internal/repository"
	"github.com/latoulicious/siresto-backend/internal/utils"
	"github.com/latoulicious/siresto-backend/internal/validator"
)

type UserService struct {
	Repo *repository.UserRepository
}

func (s *UserService) ListAllUsers() ([]domain.User, error) {
	users, err := s.Repo.ListAllUsers()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *UserService) GetUserByID(id uuid.UUID) (*domain.User, error) {
	user, err := s.Repo.GetUserByID(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) CreateUser(req *validator.CreateUserRequest) (*domain.User, error) {
	hashedPassword := utils.HashSHA256(req.Password)

	user := &domain.User{
		Name:      req.Name,
		Email:     req.Email,
		Password:  hashedPassword,
		IsStaff:   false, // Default unless you need otherwise
		CreatedAt: time.Now(),
	}

	createdUser, err := s.Repo.CreateUser(user)
	if err != nil {
		return nil, err
	}
	return createdUser, nil
}

func (s *UserService) UpdateUser(id uuid.UUID, update *domain.User) (*domain.User, error) {
	updatedUser, err := s.Repo.UpdateUser(update)
	if err != nil {
		return nil, err
	}
	return updatedUser, nil
}

func (s *UserService) DeleteUser(id uuid.UUID) error {
	err := s.Repo.DeleteUser(id)
	if err != nil {
		return err
	}
	return nil
}

// LoginUser checks if the user exists and returns the user if found
func (s *UserService) LoginUser(req *validator.LoginRequest) (*domain.User, error) {
	user, err := s.Repo.FindByEmail(req.Email)
	if err != nil || user == nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	hashedInput := utils.HashSHA256(req.Password)
	if user.Password != hashedInput {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Update last login
	now := time.Now()
	user.LastLoginAt = &now
	_ = s.Repo.UpdateLastLogin(user.ID, now) // ignore error for now

	return user, nil
}
