package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/latoulicious/siresto-backend/internal/domain"
	"github.com/latoulicious/siresto-backend/internal/repository"
	"github.com/latoulicious/siresto-backend/pkg/db"
)

type CreateLogRequest struct {
	Level       string                 `json:"level"`
	Source      string                 `json:"source"`
	Action      string                 `json:"action"`
	Entity      string                 `json:"entity"`
	EntityID    *uuid.UUID             `json:"entity_id"`
	UserID      *uuid.UUID             `json:"user_id"`
	IPAddress   *string                `json:"ip_address"`
	Description string                 `json:"description"`
	Metadata    map[string]interface{} `json:"metadata"`
	RequestID   *string                `json:"request_id"`
	Environment string                 `json:"environment"`
	Application string                 `json:"application"`
	Hostname    *string                `json:"hostname"`
	Type        string                 `json:"type"`
}

type LogService struct {
	Repo repository.LogRepositoryInterface
}

func (s *LogService) CreateLog(ctx context.Context, req CreateLogRequest) (*domain.Log, error) {
	// Validate required fields
	if req.Level == "" || req.Source == "" || req.Action == "" ||
		req.Entity == "" || req.Environment == "" || req.Application == "" || req.Type == "" {
		return nil, errors.New("missing required log fields")
	}

	// Convert Metadata to JSONB
	metadataJSON := db.JSONB{}
	if req.Metadata != nil {
		metadataJSON = db.JSONB(req.Metadata)
	}

	// Create the log entity with all fields
	log := domain.Log{
		ID:          uuid.New(),
		Timestamp:   time.Now(),
		Level:       req.Level,
		Source:      req.Source,
		UserID:      req.UserID,
		IPAddress:   req.IPAddress,
		Action:      req.Action,
		Entity:      req.Entity,
		EntityID:    req.EntityID,
		Description: req.Description,
		Metadata:    metadataJSON,
		RequestID:   req.RequestID,
		Environment: req.Environment,
		Application: req.Application,
		Hostname:    req.Hostname,
		Type:        req.Type,
	}

	// Save the log
	if err := s.Repo.SaveLog(ctx, log); err != nil {
		return nil, err
	}

	return &log, nil
}
