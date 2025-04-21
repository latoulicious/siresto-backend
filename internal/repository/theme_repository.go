package repository

import (
	"github.com/google/uuid"
	"github.com/latoulicious/siresto-backend/internal/domain"
	"gorm.io/gorm"
)

type ThemeRepository struct {
	DB *gorm.DB
}

func (r *ThemeRepository) ListAllThemes() ([]domain.Theme, error) {
	var themes []domain.Theme
	err := r.DB.Find(&themes).Error
	if err != nil {
		return nil, err
	}
	return themes, nil
}

func (r *ThemeRepository) GetThemeByID(id uuid.UUID) (*domain.Theme, error) {
	var theme domain.Theme
	err := r.DB.First(&theme, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &theme, nil
}

func (r *ThemeRepository) CreateTheme(theme *domain.Theme) error {
	err := r.DB.Create(theme).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *ThemeRepository) UpdateTheme(theme *domain.Theme) error {
	err := r.DB.Save(theme).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *ThemeRepository) DeleteTheme(id uuid.UUID) error {
	var theme domain.Theme
	err := r.DB.First(&theme, "id = ?", id).Error
	if err != nil {
		return err
	}
	err = r.DB.Delete(&theme).Error
	if err != nil {
		return err
	}
	return nil
}
