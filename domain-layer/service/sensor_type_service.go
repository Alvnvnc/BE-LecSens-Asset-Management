package service

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"be-lecsens/asset_management/data-layer/repository"
	"be-lecsens/asset_management/helpers/dto"

	"github.com/google/uuid"
)

// SensorTypeService handles business logic for sensor types
type SensorTypeService struct {
	sensorTypeRepo *repository.SensorTypeRepository
}

// NewSensorTypeService creates a new instance of SensorTypeService
func NewSensorTypeService(sensorTypeRepo *repository.SensorTypeRepository) *SensorTypeService {
	return &SensorTypeService{
		sensorTypeRepo: sensorTypeRepo,
	}
}

// CreateSensorType creates a new sensor type
func (s *SensorTypeService) CreateSensorType(ctx context.Context, req *dto.CreateSensorTypeRequest) (*dto.SensorTypeDTO, error) {
	// Create new sensor type
	now := time.Now()
	sensorType := &repository.SensorType{
		ID:           uuid.New(),
		Name:         req.Name,
		Description:  req.Description,
		Manufacturer: req.Manufacturer,
		Model:        req.Model,
		Version:      req.Version,
		IsActive:     req.IsActive,
		CreatedAt:    now,
		UpdatedAt:    &now,
	}

	// Save to repository
	err := s.sensorTypeRepo.Create(sensorType)
	if err != nil {
		return nil, err
	}

	// Convert to DTO
	return &dto.SensorTypeDTO{
		ID:           sensorType.ID,
		Name:         sensorType.Name,
		Description:  sensorType.Description,
		Manufacturer: sensorType.Manufacturer,
		Model:        sensorType.Model,
		Version:      sensorType.Version,
		IsActive:     sensorType.IsActive,
		CreatedAt:    sensorType.CreatedAt,
		UpdatedAt:    sensorType.UpdatedAt,
	}, nil
}

// GetSensorType retrieves a sensor type by ID
func (s *SensorTypeService) GetSensorType(ctx context.Context, id uuid.UUID) (*dto.SensorTypeDTO, error) {
	sensorType, err := s.sensorTypeRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if sensorType == nil {
		return nil, errors.New("sensor type not found")
	}

	return &dto.SensorTypeDTO{
		ID:           sensorType.ID,
		Name:         sensorType.Name,
		Description:  sensorType.Description,
		Manufacturer: sensorType.Manufacturer,
		Model:        sensorType.Model,
		Version:      sensorType.Version,
		IsActive:     sensorType.IsActive,
		CreatedAt:    sensorType.CreatedAt,
		UpdatedAt:    sensorType.UpdatedAt,
	}, nil
}

// ListSensorTypes retrieves a list of sensor types with pagination
func (s *SensorTypeService) ListSensorTypes(ctx context.Context, page, pageSize int) ([]*dto.SensorTypeDTO, error) {
	sensorTypes, err := s.sensorTypeRepo.GetAll()
	if err != nil {
		return nil, err
	}

	// Convert to DTOs
	var dtos []*dto.SensorTypeDTO
	for _, st := range sensorTypes {
		dtos = append(dtos, &dto.SensorTypeDTO{
			ID:           st.ID,
			Name:         st.Name,
			Description:  st.Description,
			Manufacturer: st.Manufacturer,
			Model:        st.Model,
			Version:      st.Version,
			IsActive:     st.IsActive,
			CreatedAt:    st.CreatedAt,
			UpdatedAt:    st.UpdatedAt,
		})
	}

	return dtos, nil
}

// UpdateSensorType updates an existing sensor type
func (s *SensorTypeService) UpdateSensorType(ctx context.Context, id uuid.UUID, req *dto.UpdateSensorTypeRequest) (*dto.SensorTypeDTO, error) {
	// Get existing sensor type
	existingSensorType, err := s.sensorTypeRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if existingSensorType == nil {
		return nil, errors.New("sensor type not found")
	}

	// Update fields if provided
	if req.Name != "" {
		existingSensorType.Name = req.Name
	}
	if req.Description != "" {
		existingSensorType.Description = req.Description
	}
	if req.Manufacturer != "" {
		existingSensorType.Manufacturer = req.Manufacturer
	}
	if req.Model != "" {
		existingSensorType.Model = req.Model
	}
	if req.Version != "" {
		existingSensorType.Version = req.Version
	}
	existingSensorType.IsActive = req.IsActive
	now := time.Now()
	existingSensorType.UpdatedAt = &now

	// Save to repository
	err = s.sensorTypeRepo.Update(existingSensorType)
	if err != nil {
		return nil, err
	}

	// Convert to DTO
	return &dto.SensorTypeDTO{
		ID:           existingSensorType.ID,
		Name:         existingSensorType.Name,
		Description:  existingSensorType.Description,
		Manufacturer: existingSensorType.Manufacturer,
		Model:        existingSensorType.Model,
		Version:      existingSensorType.Version,
		IsActive:     existingSensorType.IsActive,
		CreatedAt:    existingSensorType.CreatedAt,
		UpdatedAt:    existingSensorType.UpdatedAt,
	}, nil
}

// DeleteSensorType deletes a sensor type
func (s *SensorTypeService) DeleteSensorType(ctx context.Context, id uuid.UUID) error {
	return s.sensorTypeRepo.Delete(id)
}

// GetActiveSensorTypes retrieves all active sensor types
func (s *SensorTypeService) GetActiveSensorTypes(ctx context.Context) ([]*dto.SensorTypeDTO, error) {
	sensorTypes, err := s.sensorTypeRepo.GetActive()
	if err != nil {
		return nil, err
	}

	// Convert to DTOs
	var dtos []*dto.SensorTypeDTO
	for _, st := range sensorTypes {
		dtos = append(dtos, &dto.SensorTypeDTO{
			ID:           st.ID,
			Name:         st.Name,
			Description:  st.Description,
			Manufacturer: st.Manufacturer,
			Model:        st.Model,
			Version:      st.Version,
			IsActive:     st.IsActive,
			CreatedAt:    st.CreatedAt,
			UpdatedAt:    st.UpdatedAt,
		})
	}

	return dtos, nil
}

// UpdateSensorTypePartial updates specific fields of an existing sensor type
func (s *SensorTypeService) UpdateSensorTypePartial(ctx context.Context, id uuid.UUID, updateReq interface{}) (*dto.SensorTypeDTO, error) {
	// Get existing sensor type
	existingSensorType, err := s.sensorTypeRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if existingSensorType == nil {
		return nil, errors.New("sensor type not found")
	}

	// Convert updateReq to map for easier handling
	updateBytes, err := json.Marshal(updateReq)
	if err != nil {
		return nil, err
	}

	var updateRequest map[string]interface{}
	err = json.Unmarshal(updateBytes, &updateRequest)
	if err != nil {
		return nil, err
	}

	// Update only provided fields
	if name, exists := updateRequest["name"]; exists && name != nil {
		existingSensorType.Name = name.(string)
	}
	if description, exists := updateRequest["description"]; exists && description != nil {
		existingSensorType.Description = description.(string)
	}
	if manufacturer, exists := updateRequest["manufacturer"]; exists && manufacturer != nil {
		existingSensorType.Manufacturer = manufacturer.(string)
	}
	if model, exists := updateRequest["model"]; exists && model != nil {
		existingSensorType.Model = model.(string)
	}
	if version, exists := updateRequest["version"]; exists && version != nil {
		existingSensorType.Version = version.(string)
	}
	if isActive, exists := updateRequest["is_active"]; exists && isActive != nil {
		existingSensorType.IsActive = isActive.(bool)
	}

	// Update timestamp
	now := time.Now()
	existingSensorType.UpdatedAt = &now

	// Update in repository
	err = s.sensorTypeRepo.Update(existingSensorType)
	if err != nil {
		return nil, err
	}

	// Convert to DTO
	return &dto.SensorTypeDTO{
		ID:           existingSensorType.ID,
		Name:         existingSensorType.Name,
		Description:  existingSensorType.Description,
		Manufacturer: existingSensorType.Manufacturer,
		Model:        existingSensorType.Model,
		Version:      existingSensorType.Version,
		IsActive:     existingSensorType.IsActive,
		CreatedAt:    existingSensorType.CreatedAt,
		UpdatedAt:    existingSensorType.UpdatedAt,
	}, nil
}
