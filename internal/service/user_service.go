package service

import (
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
