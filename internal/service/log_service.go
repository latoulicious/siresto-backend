package service

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/latoulicious/siresto-backend/internal/domain"
	"github.com/latoulicious/siresto-backend/internal/repository"
	"github.com/latoulicious/siresto-backend/pkg/db"
)

type CreateLogRequest struct {
	Level       string
	Source      string
	UserID      string
	IPAddress   string
	Action      string
	Entity      string
	EntityID    string
	Description string
	Metadata    string
	RequestID   string
	Environment string
	Application string
	Hostname    string
	Type        string
}

type LogService struct {
	Repo repository.LogRepositoryInterface // <-- Change this to use the interface
}

// Make sure the CreateLog method is still the same
func (s *LogService) CreateLog(req CreateLogRequest) error {
	// Parse UserID as *uuid.UUID if it exists
	var userID *uuid.UUID
	if req.UserID != "" {
		parsedUserID, err := uuid.Parse(req.UserID)
		if err != nil {
			return err
		}
		userID = &parsedUserID
	}

	// Parse EntityID as *uuid.UUID if it exists
	var entityID *uuid.UUID
	if req.EntityID != "" {
		parsedEntityID, err := uuid.Parse(req.EntityID)
		if err != nil {
			return err
		}
		entityID = &parsedEntityID
	}

	// Parse IPAddress as *string if it exists
	var ipAddress *string
	if req.IPAddress != "" {
		ipAddress = &req.IPAddress
	}

	// Parse Metadata as db.JSONB if it exists
	var metadata db.JSONB
	if req.Metadata != "" {
		if err := json.Unmarshal([]byte(req.Metadata), &metadata); err != nil {
			return err
		}
	}

	// Parse RequestID as *string if it exists
	var requestID *string
	if req.RequestID != "" {
		requestID = &req.RequestID
	}

	// Parse Hostname as *string if it exists
	var hostname *string
	if req.Hostname != "" {
		hostname = &req.Hostname
	}

	// Create the Log object
	log := domain.Log{
		Level:       req.Level,
		Source:      req.Source,
		UserID:      userID,
		IPAddress:   ipAddress,
		Action:      req.Action,
		Entity:      req.Entity,
		EntityID:    entityID,
		Description: req.Description,
		Metadata:    metadata,
		RequestID:   requestID,
		Environment: req.Environment,
		Application: req.Application,
		Hostname:    hostname,
		Type:        req.Type,
	}

	// Save log into the repository
	if err := s.Repo.SaveLog(log); err != nil {
		return err
	}

	return nil
}
