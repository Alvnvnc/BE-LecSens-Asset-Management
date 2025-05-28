package service

import (
	"be-lecsens/asset_management/data-layer/repository"
	"be-lecsens/asset_management/helpers/dto"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

// SensorMeasurementTypeService handles business logic for sensor measurement types
type SensorMeasurementTypeService struct {
	repo *repository.SensorMeasurementTypeRepository
}

// NewSensorMeasurementTypeService creates a new instance of SensorMeasurementTypeService
func NewSensorMeasurementTypeService(repo *repository.SensorMeasurementTypeRepository) *SensorMeasurementTypeService {
	return &SensorMeasurementTypeService{
		repo: repo,
	}
}

// CreateSensorMeasurementType creates a new sensor measurement type
func (s *SensorMeasurementTypeService) CreateSensorMeasurementType(ctx context.Context, req dto.CreateSensorMeasurementTypeRequest) (*dto.SensorMeasurementTypeDTO, error) {
	// Validate min/max values if provided
	if req.MinAcceptedValue != nil && req.MaxAcceptedValue != nil {
		if *req.MinAcceptedValue >= *req.MaxAcceptedValue {
			return nil, errors.New("min_accepted_value must be less than max_accepted_value")
		}
	}

	now := time.Now()
	sensorMeasurementType := &dto.SensorMeasurementTypeDTO{
		ID:               uuid.New(),
		SensorTypeID:     req.SensorTypeID,
		Name:             req.Name,
		Description:      req.Description,
		UnitOfMeasure:    req.UnitOfMeasure,
		MinAcceptedValue: req.MinAcceptedValue,
		MaxAcceptedValue: req.MaxAcceptedValue,
		PropertiesSchema: req.PropertiesSchema,
		UIConfiguration:  req.UIConfiguration,
		Version:          req.Version,
		IsActive:         req.IsActive,
		CreatedAt:        now,
		UpdatedAt:        &now,
	}

	err := s.repo.Create(ctx, sensorMeasurementType)
	if err != nil {
		return nil, err
	}

	return sensorMeasurementType, nil
}

// GetSensorMeasurementType retrieves a sensor measurement type by ID
func (s *SensorMeasurementTypeService) GetSensorMeasurementType(ctx context.Context, id uuid.UUID) (*dto.SensorMeasurementTypeDTO, error) {
	return s.repo.GetByID(ctx, id)
}

// ListSensorMeasurementTypes retrieves a paginated list of sensor measurement types
func (s *SensorMeasurementTypeService) ListSensorMeasurementTypes(ctx context.Context, page, limit int) ([]dto.SensorMeasurementTypeDTO, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit
	return s.repo.List(ctx, offset, limit)
}

// UpdateSensorMeasurementType updates an existing sensor measurement type
func (s *SensorMeasurementTypeService) UpdateSensorMeasurementType(ctx context.Context, id uuid.UUID, req dto.UpdateSensorMeasurementTypeRequest) (*dto.SensorMeasurementTypeDTO, error) {
	// Get existing sensor measurement type
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Name != nil {
		existing.Name = *req.Name
	}
	if req.Description != nil {
		existing.Description = req.Description
	}
	if req.UnitOfMeasure != nil {
		existing.UnitOfMeasure = req.UnitOfMeasure
	}
	if req.MinAcceptedValue != nil {
		existing.MinAcceptedValue = req.MinAcceptedValue
	}
	if req.MaxAcceptedValue != nil {
		existing.MaxAcceptedValue = req.MaxAcceptedValue
	}
	if req.PropertiesSchema != nil {
		existing.PropertiesSchema = req.PropertiesSchema
	}
	if req.UIConfiguration != nil {
		existing.UIConfiguration = req.UIConfiguration
	}
	if req.Version != nil {
		existing.Version = *req.Version
	}
	if req.IsActive != nil {
		existing.IsActive = *req.IsActive
	}

	// Validate min/max values if both are provided
	if existing.MinAcceptedValue != nil && existing.MaxAcceptedValue != nil {
		if *existing.MinAcceptedValue >= *existing.MaxAcceptedValue {
			return nil, errors.New("min_accepted_value must be less than max_accepted_value")
		}
	}

	now := time.Now()
	existing.UpdatedAt = &now

	err = s.repo.Update(ctx, existing)
	if err != nil {
		return nil, err
	}

	return existing, nil
}

// DeleteSensorMeasurementType deletes a sensor measurement type
func (s *SensorMeasurementTypeService) DeleteSensorMeasurementType(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

// GetActiveSensorMeasurementTypes retrieves all active sensor measurement types
func (s *SensorMeasurementTypeService) GetActiveSensorMeasurementTypes(ctx context.Context) ([]dto.SensorMeasurementTypeDTO, error) {
	return s.repo.GetActive(ctx)
}

// UpdateSensorMeasurementTypePartial updates specific fields of a sensor measurement type
func (s *SensorMeasurementTypeService) UpdateSensorMeasurementTypePartial(ctx context.Context, id uuid.UUID, req dto.UpdateSensorMeasurementTypeRequest) (*dto.SensorMeasurementTypeDTO, error) {
	return s.UpdateSensorMeasurementType(ctx, id, req)
}

// GetSensorMeasurementTypesBySensorTypeID retrieves all measurement types for a specific sensor type
func (s *SensorMeasurementTypeService) GetSensorMeasurementTypesBySensorTypeID(ctx context.Context, sensorTypeID uuid.UUID) ([]dto.SensorMeasurementTypeDTO, error) {
	return s.repo.GetBySensorTypeID(ctx, sensorTypeID)
}
