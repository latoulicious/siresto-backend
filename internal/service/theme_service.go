package service

import (
	"github.com/google/uuid"
	"github.com/latoulicious/siresto-backend/internal/domain"
	"github.com/latoulicious/siresto-backend/internal/repository"
)

type ThemeService struct {
	Repo *repository.ThemeRepository
}

func (s *ThemeService) ListAllThemes() ([]domain.Theme, error) {
	themes, err := s.Repo.ListAllThemes()
	if err != nil {
		return nil, err
	}
	return themes, nil
}

func (s *ThemeService) GetThemeByID(id uuid.UUID) (*domain.Theme, error) {
	theme, err := s.Repo.GetThemeByID(id)
	if err != nil {
		return nil, err
	}
	return theme, nil
}

func (s *ThemeService) CreateTheme(theme *domain.Theme) error {
	err := s.Repo.CreateTheme(theme)
	if err != nil {
		return err
	}
	return nil
}

func (s *ThemeService) UpdateTheme(theme *domain.Theme) error {
	err := s.Repo.UpdateTheme(theme)
	if err != nil {
		return err
	}
	return nil
}

func (s *ThemeService) DeleteTheme(id uuid.UUID) error {
	err := s.Repo.DeleteTheme(id)
	if err != nil {
		return err
	}
	return nil
}
