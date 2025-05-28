package service

import (
	"be-lecsens/asset_management/data-layer/repository"
	"be-lecsens/asset_management/helpers/dto"
	"context"
	"errors"

	"github.com/google/uuid"
)

// SensorMeasurementFieldService handles business logic for sensor measurement fields
type SensorMeasurementFieldService struct {
	repo *repository.SensorMeasurementFieldRepository
}

// NewSensorMeasurementFieldService creates a new instance of SensorMeasurementFieldService
func NewSensorMeasurementFieldService(repo *repository.SensorMeasurementFieldRepository) *SensorMeasurementFieldService {
	return &SensorMeasurementFieldService{
		repo: repo,
	}
}

// GetAll retrieves all sensor measurement fields
func (s *SensorMeasurementFieldService) GetAll(ctx context.Context) ([]dto.SensorMeasurementFieldDTO, error) {
	fields, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	// Convert repository models to DTOs
	var dtos []dto.SensorMeasurementFieldDTO
	for _, field := range fields {
		dto := dto.SensorMeasurementFieldDTO{
			ID:                      field.ID,
			SensorMeasurementTypeID: field.SensorMeasurementTypeID,
			Name:                    field.Name,
			Label:                   field.Label,
			DataType:                field.DataType,
			Required:                field.Required,
			CreatedAt:               field.CreatedAt,
			UpdatedAt:               field.UpdatedAt,
		}

		// Handle nullable fields
		if field.Description.Valid {
			dto.Description = &field.Description.String
		}
		if field.Unit.Valid {
			dto.Unit = &field.Unit.String
		}
		if field.Min.Valid {
			dto.Min = &field.Min.Float64
		}
		if field.Max.Valid {
			dto.Max = &field.Max.Float64
		}

		dtos = append(dtos, dto)
	}

	return dtos, nil
}

// Create creates a new sensor measurement field
func (s *SensorMeasurementFieldService) Create(ctx context.Context, req *dto.CreateSensorMeasurementFieldRequest) (*dto.SensorMeasurementFieldDTO, error) {
	// Validate required fields
	if req.SensorMeasurementTypeID == uuid.Nil {
		return nil, errors.New("sensor measurement type ID is required")
	}
	if req.Name == "" {
		return nil, errors.New("name is required")
	}
	if req.Label == "" {
		return nil, errors.New("label is required")
	}
	if req.DataType == "" {
		return nil, errors.New("data type is required")
	}

	// Create field in repository
	field := &repository.SensorMeasurementField{
		SensorMeasurementTypeID: req.SensorMeasurementTypeID,
		Name:                    req.Name,
		Label:                   req.Label,
		DataType:                req.DataType,
		Required:                req.Required,
	}

	// Handle optional fields
	if req.Description != nil {
		field.Description = repository.NullStringFromPtr(req.Description)
	}
	if req.Unit != nil {
		field.Unit = repository.NullStringFromPtr(req.Unit)
	}
	if req.Min != nil {
		field.Min = repository.NullFloat64FromPtr(req.Min)
	}
	if req.Max != nil {
		field.Max = repository.NullFloat64FromPtr(req.Max)
	}

	createdField, err := s.repo.Create(ctx, field)
	if err != nil {
		return nil, err
	}

	// Convert to DTO
	dto := dto.SensorMeasurementFieldDTO{
		ID:                      createdField.ID,
		SensorMeasurementTypeID: createdField.SensorMeasurementTypeID,
		Name:                    createdField.Name,
		Label:                   createdField.Label,
		DataType:                createdField.DataType,
		Required:                createdField.Required,
		CreatedAt:               createdField.CreatedAt,
		UpdatedAt:               createdField.UpdatedAt,
	}

	// Handle nullable fields
	if createdField.Description.Valid {
		dto.Description = &createdField.Description.String
	}
	if createdField.Unit.Valid {
		dto.Unit = &createdField.Unit.String
	}
	if createdField.Min.Valid {
		dto.Min = &createdField.Min.Float64
	}
	if createdField.Max.Valid {
		dto.Max = &createdField.Max.Float64
	}

	return &dto, nil
}

// GetByID retrieves a sensor measurement field by its ID
func (s *SensorMeasurementFieldService) GetByID(ctx context.Context, id uuid.UUID) (*dto.SensorMeasurementFieldDTO, error) {
	field, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Check if field was found
	if field == nil {
		return nil, errors.New("sensor measurement field not found")
	}

	// Convert to DTO
	dto := dto.SensorMeasurementFieldDTO{
		ID:                      field.ID,
		SensorMeasurementTypeID: field.SensorMeasurementTypeID,
		Name:                    field.Name,
		Label:                   field.Label,
		DataType:                field.DataType,
		Required:                field.Required,
		CreatedAt:               field.CreatedAt,
		UpdatedAt:               field.UpdatedAt,
	}

	// Handle nullable fields
	if field.Description.Valid {
		dto.Description = &field.Description.String
	}
	if field.Unit.Valid {
		dto.Unit = &field.Unit.String
	}
	if field.Min.Valid {
		dto.Min = &field.Min.Float64
	}
	if field.Max.Valid {
		dto.Max = &field.Max.Float64
	}

	return &dto, nil
}

// GetByMeasurementTypeID retrieves all fields for a measurement type
func (s *SensorMeasurementFieldService) GetByMeasurementTypeID(ctx context.Context, measurementTypeID uuid.UUID) ([]dto.SensorMeasurementFieldDTO, error) {
	fields, err := s.repo.GetByMeasurementTypeID(ctx, measurementTypeID)
	if err != nil {
		return nil, err
	}

	// Convert to DTOs
	var dtos []dto.SensorMeasurementFieldDTO
	for _, field := range fields {
		dto := dto.SensorMeasurementFieldDTO{
			ID:                      field.ID,
			SensorMeasurementTypeID: field.SensorMeasurementTypeID,
			Name:                    field.Name,
			Label:                   field.Label,
			DataType:                field.DataType,
			Required:                field.Required,
			CreatedAt:               field.CreatedAt,
			UpdatedAt:               field.UpdatedAt,
		}

		// Handle nullable fields
		if field.Description.Valid {
			dto.Description = &field.Description.String
		}
		if field.Unit.Valid {
			dto.Unit = &field.Unit.String
		}
		if field.Min.Valid {
			dto.Min = &field.Min.Float64
		}
		if field.Max.Valid {
			dto.Max = &field.Max.Float64
		}

		dtos = append(dtos, dto)
	}

	return dtos, nil
}

// Update updates an existing sensor measurement field
func (s *SensorMeasurementFieldService) Update(ctx context.Context, id uuid.UUID, req *dto.UpdateSensorMeasurementFieldRequest) (*dto.SensorMeasurementFieldDTO, error) {
	// Get existing field
	field, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Check if field was found
	if field == nil {
		return nil, errors.New("sensor measurement field not found")
	}

	// Update fields if provided
	if req.Name != nil {
		field.Name = *req.Name
	}
	if req.Label != nil {
		field.Label = *req.Label
	}
	if req.Description != nil {
		field.Description = repository.NullStringFromPtr(req.Description)
	}
	if req.DataType != nil {
		field.DataType = *req.DataType
	}
	if req.Required != nil {
		field.Required = *req.Required
	}
	if req.Unit != nil {
		field.Unit = repository.NullStringFromPtr(req.Unit)
	}
	if req.Min != nil {
		field.Min = repository.NullFloat64FromPtr(req.Min)
	}
	if req.Max != nil {
		field.Max = repository.NullFloat64FromPtr(req.Max)
	}

	// Update in repository
	updatedField, err := s.repo.Update(ctx, field)
	if err != nil {
		return nil, err
	}

	// Convert to DTO
	dto := dto.SensorMeasurementFieldDTO{
		ID:                      updatedField.ID,
		SensorMeasurementTypeID: updatedField.SensorMeasurementTypeID,
		Name:                    updatedField.Name,
		Label:                   updatedField.Label,
		DataType:                updatedField.DataType,
		Required:                updatedField.Required,
		CreatedAt:               updatedField.CreatedAt,
		UpdatedAt:               updatedField.UpdatedAt,
	}

	// Handle nullable fields
	if updatedField.Description.Valid {
		dto.Description = &updatedField.Description.String
	}
	if updatedField.Unit.Valid {
		dto.Unit = &updatedField.Unit.String
	}
	if updatedField.Min.Valid {
		dto.Min = &updatedField.Min.Float64
	}
	if updatedField.Max.Valid {
		dto.Max = &updatedField.Max.Float64
	}

	return &dto, nil
}

// Delete deletes a sensor measurement field
func (s *SensorMeasurementFieldService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

// GetRequiredFields retrieves all required fields for a measurement type
func (s *SensorMeasurementFieldService) GetRequiredFields(ctx context.Context, measurementTypeID uuid.UUID) ([]dto.SensorMeasurementFieldDTO, error) {
	fields, err := s.repo.GetRequiredFields(ctx, measurementTypeID)
	if err != nil {
		return nil, err
	}

	// Convert to DTOs
	var dtos []dto.SensorMeasurementFieldDTO
	for _, field := range fields {
		dto := dto.SensorMeasurementFieldDTO{
			ID:                      field.ID,
			SensorMeasurementTypeID: field.SensorMeasurementTypeID,
			Name:                    field.Name,
			Label:                   field.Label,
			DataType:                field.DataType,
			Required:                field.Required,
			CreatedAt:               field.CreatedAt,
			UpdatedAt:               field.UpdatedAt,
		}

		// Handle nullable fields
		if field.Description.Valid {
			dto.Description = &field.Description.String
		}
		if field.Unit.Valid {
			dto.Unit = &field.Unit.String
		}
		if field.Min.Valid {
			dto.Min = &field.Min.Float64
		}
		if field.Max.Valid {
			dto.Max = &field.Max.Float64
		}

		dtos = append(dtos, dto)
	}

	return dtos, nil
}
