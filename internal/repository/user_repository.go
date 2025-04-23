package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/latoulicious/siresto-backend/internal/domain"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func (r *UserRepository) ListAllUsers() ([]domain.User, error) {
	var users []domain.User
	err := r.DB.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) GetUserByID(id uuid.UUID) (*domain.User, error) {
	var user domain.User
	err := r.DB.First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) CreateUser(user *domain.User) (*domain.User, error) {
	if err := r.DB.Create(user).Error; err != nil {
		return nil, err
	}

	// Preload role after creation
	var createdUser domain.User
	if err := r.DB.
		Preload("Role").
		First(&createdUser, "id = ?", user.ID).Error; err != nil {
		return nil, err
	}

	return &createdUser, nil
}

func (r *UserRepository) UpdateUser(user *domain.User) (*domain.User, error) {
	err := r.DB.Save(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) DeleteUser(id uuid.UUID) error {
	err := r.DB.Delete(&domain.User{}, id).Error
	if err != nil {
		return err
	}
	return nil
}

// LoginUser checks if the user exists and returns the user if found
func (r *UserRepository) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	if err := r.DB.Preload("Role").Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) UpdateLastLogin(userID uuid.UUID, time time.Time) error {
	return r.DB.Model(&domain.User{}).
		Where("id = ?", userID).
		Update("last_login_at", time).Error
}
