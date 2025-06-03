package service

import (
	"be-lecsens/asset_management/data-layer/entity"
	"be-lecsens/asset_management/data-layer/repository"
	"be-lecsens/asset_management/helpers/common"
	"be-lecsens/asset_management/helpers/dto"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

// AssetSensorService handles business logic for asset sensors
type AssetSensorService struct {
	assetSensorRepo repository.AssetSensorRepository
	assetRepo       repository.AssetRepository
}

// NewAssetSensorService creates a new instance of AssetSensorService
func NewAssetSensorService(
	assetSensorRepo repository.AssetSensorRepository,
	assetRepo repository.AssetRepository,
) *AssetSensorService {
	return &AssetSensorService{
		assetSensorRepo: assetSensorRepo,
		assetRepo:       assetRepo,
	}
}

// CreateAssetSensor creates a new asset sensor
func (s *AssetSensorService) CreateAssetSensor(ctx context.Context, req *dto.CreateAssetSensorRequest) (*dto.AssetSensorResponse, error) {
	// Validate request
	if req.Name == "" {
		log.Printf("Validation error: name is required")
		return nil, common.NewValidationError("name is required", nil)
	}

	if req.Status == "" {
		log.Printf("Validation error: status is required")
		return nil, common.NewValidationError("status is required", nil)
	}

	log.Printf("Validating asset with ID: %s", req.AssetID)
	// Validate asset exists and get its tenant_id
	asset, err := s.assetRepo.GetByID(ctx, req.AssetID)
	if err != nil {
		log.Printf("Error validating asset: %v", err)
		return nil, fmt.Errorf("failed to validate asset: %w", err)
	}
	if asset == nil {
		log.Printf("Asset not found with ID: %s", req.AssetID)
		return nil, common.NewNotFoundError("asset", req.AssetID.String())
	}

	log.Printf("Found asset: %+v", asset)

	// Create entity with tenant_id from asset
	sensor := &entity.AssetSensor{
		ID:            uuid.New(),
		TenantID:      asset.TenantID, // Inherit tenant_id from asset
		AssetID:       req.AssetID,
		SensorTypeID:  req.SensorTypeID,
		Name:          req.Name,
		Status:        req.Status,
		Configuration: req.Configuration,
		CreatedAt:     time.Now(),
	}

	log.Printf("Created sensor entity: %+v", sensor)

	// Save to database
	err = s.assetSensorRepo.Create(ctx, sensor)
	if err != nil {
		log.Printf("Error saving sensor to database: %v", err)
		return nil, fmt.Errorf("failed to create asset sensor: %w", err)
	}

	log.Printf("Successfully saved sensor to database with ID: %s", sensor.ID)

	return s.entityToResponse(sensor), nil
}

// GetAssetSensor retrieves an asset sensor by ID
func (s *AssetSensorService) GetAssetSensor(ctx context.Context, id uuid.UUID) (*dto.AssetSensorResponse, error) {
	sensor, err := s.assetSensorRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get asset sensor: %w", err)
	}

	if sensor == nil {
		return nil, common.NewNotFoundError("asset sensor", id.String())
	}

	return s.entityToResponse(sensor.AssetSensor), nil
}

// GetAssetSensors retrieves all sensors for a specific asset
func (s *AssetSensorService) GetAssetSensors(ctx context.Context, assetID uuid.UUID) ([]*dto.AssetSensorResponse, error) {
	// Validate that the asset exists and user has access to it
	asset, err := s.assetRepo.GetByID(ctx, assetID)
	if err != nil {
		return nil, fmt.Errorf("failed to validate asset: %w", err)
	}
	if asset == nil {
		return nil, common.NewNotFoundError("asset", assetID.String())
	}

	sensors, err := s.assetSensorRepo.GetByAssetID(ctx, assetID)
	if err != nil {
		return nil, fmt.Errorf("failed to get asset sensors: %w", err)
	}

	responses := make([]*dto.AssetSensorResponse, len(sensors))
	for i, sensor := range sensors {
		responses[i] = s.entityToResponse(sensor.AssetSensor)
	}
	return responses, nil
}

// GetAssetSensorsDetailed retrieves all sensors for a specific asset with complete details including sensor types and measurement types
func (s *AssetSensorService) GetAssetSensorsDetailed(ctx context.Context, assetID uuid.UUID) ([]*dto.AssetSensorDetailedResponse, error) {
	// Validate that the asset exists and user has access to it
	asset, err := s.assetRepo.GetByID(ctx, assetID)
	if err != nil {
		return nil, fmt.Errorf("failed to validate asset: %w", err)
	}
	if asset == nil {
		return nil, common.NewNotFoundError("asset", assetID.String())
	}

	sensorsWithDetails, err := s.assetSensorRepo.GetByAssetID(ctx, assetID)
	if err != nil {
		return nil, fmt.Errorf("failed to get asset sensors: %w", err)
	}

	responses := make([]*dto.AssetSensorDetailedResponse, len(sensorsWithDetails))
	for i, sensorWithDetails := range sensorsWithDetails {
		responses[i] = s.entityToDetailedResponse(sensorWithDetails)
	}
	return responses, nil
}

// ListAssetSensors retrieves asset sensors with pagination
func (s *AssetSensorService) ListAssetSensors(ctx context.Context, page, pageSize int) (*dto.AssetSensorListResponse, error) {
	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10 // Default page size
	}

	sensors, err := s.assetSensorRepo.List(ctx, page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("failed to list asset sensors: %w", err)
	}

	// TODO: Get total count from repository
	total := int64(len(sensors))
	totalPages := (int(total) + pageSize - 1) / pageSize

	// Convert sensors to response type
	sensorResponses := make([]dto.AssetSensorResponse, len(sensors))
	for i, sensor := range sensors {
		response := s.entityToResponse(sensor.AssetSensor)
		sensorResponses[i] = *response
	}

	return &dto.AssetSensorListResponse{
		Sensors:    sensorResponses,
		Page:       page,
		Limit:      pageSize,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

// UpdateAssetSensor updates an existing asset sensor
func (s *AssetSensorService) UpdateAssetSensor(ctx context.Context, id uuid.UUID, req *dto.UpdateAssetSensorRequest) (*dto.AssetSensorResponse, error) {
	// Get existing sensor to validate ownership and existence
	existingSensor, err := s.assetSensorRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing sensor: %w", err)
	}
	if existingSensor == nil {
		return nil, common.NewNotFoundError("asset sensor", id.String())
	}

	// Start with existing sensor
	updatedSensor := *existingSensor.AssetSensor

	// Update fields if provided
	if req.Name != nil {
		if *req.Name == "" {
			return nil, common.NewValidationError("name cannot be empty", nil)
		}
		updatedSensor.Name = *req.Name
	}

	if req.Status != nil {
		if *req.Status == "" {
			return nil, common.NewValidationError("status cannot be empty", nil)
		}
		updatedSensor.Status = *req.Status
	}

	if req.Configuration != nil {
		updatedSensor.Configuration = *req.Configuration
	}

	// Update in repository
	err = s.assetSensorRepo.Update(ctx, &updatedSensor)
	if err != nil {
		return nil, fmt.Errorf("failed to update asset sensor: %w", err)
	}

	return s.entityToResponse(&updatedSensor), nil
}

// DeleteAssetSensor deletes an asset sensor
func (s *AssetSensorService) DeleteAssetSensor(ctx context.Context, id uuid.UUID) error {
	// Validate sensor exists and user has access
	sensor, err := s.assetSensorRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get asset sensor: %w", err)
	}
	if sensor == nil {
		return common.NewNotFoundError("asset sensor", id.String())
	}

	// Delete from database
	err = s.assetSensorRepo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete asset sensor: %w", err)
	}

	return nil
}

// DeleteAssetSensors deletes all sensors for a specific asset
func (s *AssetSensorService) DeleteAssetSensors(ctx context.Context, assetID uuid.UUID) error {
	// Validate that the asset exists and user has access to it
	asset, err := s.assetRepo.GetByID(ctx, assetID)
	if err != nil {
		return fmt.Errorf("failed to validate asset: %w", err)
	}
	if asset == nil {
		return common.NewNotFoundError("asset", assetID.String())
	}

	// Delete from database
	err = s.assetSensorRepo.DeleteByAssetID(ctx, assetID)
	if err != nil {
		return fmt.Errorf("failed to delete asset sensors: %w", err)
	}

	return nil
}

// UpdateSensorReading updates the sensor's reading values
func (s *AssetSensorService) UpdateSensorReading(ctx context.Context, id uuid.UUID, req *dto.UpdateSensorReadingRequest) error {
	// Validate sensor exists and user has access
	sensor, err := s.assetSensorRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get asset sensor: %w", err)
	}
	if sensor == nil {
		return common.NewNotFoundError("asset sensor", id.String())
	}

	// Update reading in repository
	err = s.assetSensorRepo.UpdateLastReading(ctx, id, req.Value, req.Readings)
	if err != nil {
		return fmt.Errorf("failed to update sensor reading: %w", err)
	}

	return nil
}

// GetActiveSensors retrieves all active sensors
func (s *AssetSensorService) GetActiveSensors(ctx context.Context) ([]*dto.AssetSensorResponse, error) {
	sensors, err := s.assetSensorRepo.GetActiveSensors(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get active sensors: %w", err)
	}

	responses := make([]*dto.AssetSensorResponse, len(sensors))
	for i, sensor := range sensors {
		responses[i] = s.entityToResponse(sensor.AssetSensor)
	}
	return responses, nil
}

// GetSensorsByStatus retrieves all sensors with a specific status
func (s *AssetSensorService) GetSensorsByStatus(ctx context.Context, status string) ([]*dto.AssetSensorResponse, error) {
	if status == "" {
		return nil, common.NewValidationError("status is required", nil)
	}

	sensors, err := s.assetSensorRepo.GetSensorsByStatus(ctx, status)
	if err != nil {
		return nil, fmt.Errorf("failed to get sensors by status: %w", err)
	}

	responses := make([]*dto.AssetSensorResponse, len(sensors))
	for i, sensor := range sensors {
		responses[i] = s.entityToResponse(sensor.AssetSensor)
	}
	return responses, nil
}

// GetCompleteSensorInfo retrieves complete sensor information including sensor type and measurement types
func (s *AssetSensorService) GetCompleteSensorInfo(ctx context.Context, id uuid.UUID) (*repository.AssetSensorWithDetails, error) {
	log.Printf("DEBUG: GetCompleteSensorInfo service called with sensor ID: %s", id)

	sensor, err := s.assetSensorRepo.GetByID(ctx, id)
	if err != nil {
		log.Printf("DEBUG: Error from repository GetByID: %v", err)
		return nil, fmt.Errorf("failed to get complete sensor info: %w", err)
	}
	if sensor == nil {
		log.Printf("DEBUG: Sensor not found in repository, returning NotFoundError")
		return nil, common.NewNotFoundError("asset sensor", id.String())
	}

	log.Printf("DEBUG: Successfully retrieved sensor from repository")
	log.Printf("DEBUG: Sensor details - ID: %s, Name: %s, SensorType.Version: %s",
		sensor.AssetSensor.ID, sensor.AssetSensor.Name, sensor.SensorType.Version)
	log.Printf("DEBUG: Number of measurement types: %d", len(sensor.MeasurementTypes))

	return sensor, nil
}

// entityToDetailedResponse converts an AssetSensorWithDetails to detailed response DTO
func (s *AssetSensorService) entityToDetailedResponse(sensorWithDetails *repository.AssetSensorWithDetails) *dto.AssetSensorDetailedResponse {
	sensor := sensorWithDetails.AssetSensor

	var tenantID uuid.UUID
	if sensor.TenantID != nil {
		tenantID = *sensor.TenantID
	}

	// Convert sensor type
	sensorType := dto.SensorTypeInfo{
		ID:           sensorWithDetails.SensorType.ID,
		Name:         sensorWithDetails.SensorType.Name,
		Description:  sensorWithDetails.SensorType.Description,
		Manufacturer: sensorWithDetails.SensorType.Manufacturer,
		Model:        sensorWithDetails.SensorType.Model,
		Version:      sensorWithDetails.SensorType.Version,
		IsActive:     sensorWithDetails.SensorType.IsActive,
	}

	// Convert measurement types
	measurementTypes := make([]dto.MeasurementTypeInfo, len(sensorWithDetails.MeasurementTypes))
	for i, mt := range sensorWithDetails.MeasurementTypes {
		// Convert fields
		fields := make([]dto.MeasurementFieldInfo, len(mt.Fields))
		for j, field := range mt.Fields {
			fields[j] = dto.MeasurementFieldInfo{
				ID:          field.ID,
				Name:        field.Name,
				Label:       field.Label,
				Description: field.Description,
				DataType:    field.DataType,
				Required:    field.Required,
				Unit:        field.Unit,
				Min:         field.Min,
				Max:         field.Max,
			}
		}

		measurementTypes[i] = dto.MeasurementTypeInfo{
			ID:               mt.ID,
			Name:             mt.Name,
			Description:      mt.Description,
			PropertiesSchema: mt.PropertiesSchema,
			UIConfiguration:  mt.UIConfiguration,
			Version:          mt.Version,
			IsActive:         mt.IsActive,
			Fields:           fields,
		}
	}

	return &dto.AssetSensorDetailedResponse{
		ID:                sensor.ID,
		TenantID:          tenantID,
		AssetID:           sensor.AssetID,
		SensorTypeID:      sensor.SensorTypeID,
		Name:              sensor.Name,
		Status:            sensor.Status,
		Configuration:     sensor.Configuration,
		LastReadingValue:  sensor.LastReadingValue,
		LastReadingTime:   sensor.LastReadingTime,
		LastReadingValues: sensor.LastReadingValues,
		CreatedAt:         sensor.CreatedAt,
		UpdatedAt:         sensor.UpdatedAt,
		SensorType:        sensorType,
		MeasurementTypes:  measurementTypes,
	}
}

// entityToResponse converts an entity to response DTO
func (s *AssetSensorService) entityToResponse(sensor *entity.AssetSensor) *dto.AssetSensorResponse {
	var tenantID uuid.UUID
	if sensor.TenantID != nil {
		tenantID = *sensor.TenantID
	}

	return &dto.AssetSensorResponse{
		ID:                sensor.ID,
		TenantID:          tenantID,
		AssetID:           sensor.AssetID,
		SensorTypeID:      sensor.SensorTypeID,
		Name:              sensor.Name,
		Status:            sensor.Status,
		Configuration:     sensor.Configuration,
		LastReadingValue:  sensor.LastReadingValue,
		LastReadingTime:   sensor.LastReadingTime,
		LastReadingValues: sensor.LastReadingValues,
		CreatedAt:         sensor.CreatedAt,
		UpdatedAt:         sensor.UpdatedAt,
	}
}
