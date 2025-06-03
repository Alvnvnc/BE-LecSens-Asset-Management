package service

import (
	"be-lecsens/asset_management/data-layer/entity"
	"be-lecsens/asset_management/data-layer/repository"
	"be-lecsens/asset_management/helpers/common"
	"be-lecsens/asset_management/helpers/dto"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

// IoTSensorReadingService handles business logic for IoT sensor readings
type IoTSensorReadingService struct {
	iotSensorReadingRepo      repository.IoTSensorReadingRepository
	assetSensorRepo           repository.AssetSensorRepository
	sensorTypeRepo            *repository.SensorTypeRepository
	assetRepo                 repository.AssetRepository
	locationRepo              *repository.LocationRepository
	sensorThresholdService    *SensorThresholdService                    // For threshold checking
	sensorMeasurementTypeRepo repository.SensorMeasurementTypeRepository // For getting measurement types
}

// NewIoTSensorReadingService creates a new instance of IoTSensorReadingService
func NewIoTSensorReadingService(
	iotSensorReadingRepo repository.IoTSensorReadingRepository,
	assetSensorRepo repository.AssetSensorRepository,
	sensorTypeRepo *repository.SensorTypeRepository,
	assetRepo repository.AssetRepository,
	locationRepo *repository.LocationRepository,
	sensorThresholdService *SensorThresholdService,
	sensorMeasurementTypeRepo repository.SensorMeasurementTypeRepository,
) *IoTSensorReadingService {
	return &IoTSensorReadingService{
		iotSensorReadingRepo:      iotSensorReadingRepo,
		assetSensorRepo:           assetSensorRepo,
		sensorTypeRepo:            sensorTypeRepo,
		assetRepo:                 assetRepo,
		locationRepo:              locationRepo,
		sensorThresholdService:    sensorThresholdService,
		sensorMeasurementTypeRepo: sensorMeasurementTypeRepo,
	}
}

// CreateIoTSensorReading creates a new IoT sensor reading
func (s *IoTSensorReadingService) CreateIoTSensorReading(ctx context.Context, req *dto.CreateIoTSensorReadingRequest) (*dto.IoTSensorReadingResponse, error) {
	// Validate request
	if err := s.validateCreateRequest(req); err != nil {
		log.Printf("Validation error: %v", err)
		return nil, err
	}

	// Validate asset sensor exists and get its tenant_id
	assetSensor, err := s.assetSensorRepo.GetByID(ctx, req.AssetSensorID)
	if err != nil {
		log.Printf("Error validating asset sensor: %v", err)
		return nil, fmt.Errorf("failed to validate asset sensor: %w", err)
	}
	if assetSensor == nil {
		return nil, common.NewValidationError("asset sensor not found", nil)
	}

	// Validate sensor type exists
	sensorType, err := s.sensorTypeRepo.GetByID(req.SensorTypeID)
	if err != nil {
		log.Printf("Error validating sensor type: %v", err)
		return nil, fmt.Errorf("failed to validate sensor type: %w", err)
	}
	if sensorType == nil {
		return nil, common.NewValidationError("sensor type not found", nil)
	}

	// Get location information from asset
	locationID, locationName, err := s.getLocationFromAssetSensor(ctx, req.AssetSensorID)
	if err != nil {
		log.Printf("Error getting location from asset sensor: %v", err)
		return nil, fmt.Errorf("failed to get location from asset: %w", err)
	}

	// Create entity from request
	var tenantID *uuid.UUID
	if assetSensor.AssetSensor.TenantID != nil {
		tenantID = assetSensor.AssetSensor.TenantID
	}

	macAddress := req.MacAddress
	reading := &entity.IoTSensorReadingFlexible{
		TenantID:      tenantID,
		AssetSensorID: req.AssetSensorID,
		SensorTypeID:  req.SensorTypeID,
		MacAddress:    &macAddress,
		LocationID:    &locationID,
		LocationName:  &locationName,
		ReadingTime:   time.Now(),
	}

	// Use provided reading time if specified
	if req.ReadingTime != nil {
		reading.ReadingTime = *req.ReadingTime
	}

	// Save to database
	if err := s.iotSensorReadingRepo.Create(ctx, reading); err != nil {
		log.Printf("Error creating IoT sensor reading: %v", err)
		return nil, fmt.Errorf("failed to create IoT sensor reading: %w", err)
	}

	log.Printf("Successfully created IoT sensor reading with ID: %s", reading.ID)

	// Check thresholds for the new reading (non-blocking)
	go s.checkThresholdsForReading(ctx, reading)

	// Convert to response DTO
	return s.toResponseDTO(reading), nil
}

// CreateBatchIoTSensorReading creates multiple IoT sensor readings in batch
func (s *IoTSensorReadingService) CreateBatchIoTSensorReading(ctx context.Context, req *dto.CreateBatchIoTSensorReadingRequest) ([]*dto.IoTSensorReadingResponse, error) {
	// Validate batch request
	if len(req.Readings) == 0 {
		return nil, common.NewValidationError("at least one reading is required", nil)
	}
	if len(req.Readings) > 1000 {
		return nil, common.NewValidationError("maximum 1000 readings allowed per batch", nil)
	}

	var readings []*entity.IoTSensorReadingFlexible
	var responses []*dto.IoTSensorReadingResponse

	// Validate and convert each reading
	for i, readingReq := range req.Readings {
		// Validate individual request
		if err := s.validateCreateRequest(&readingReq); err != nil {
			return nil, fmt.Errorf("validation error for reading %d: %w", i, err)
		}

		// Validate asset sensor exists and get its tenant_id
		assetSensor, err := s.assetSensorRepo.GetByID(ctx, readingReq.AssetSensorID)
		if err != nil {
			return nil, fmt.Errorf("failed to validate asset sensor for reading %d: %w", i, err)
		}
		if assetSensor == nil {
			return nil, fmt.Errorf("asset sensor not found for reading %d", i)
		}

		// Get location information from asset
		locationID, locationName, err := s.getLocationFromAssetSensor(ctx, readingReq.AssetSensorID)
		if err != nil {
			return nil, fmt.Errorf("failed to get location from asset for reading %d: %w", i, err)
		}

		// Create entity
		var tenantID *uuid.UUID
		if assetSensor.AssetSensor.TenantID != nil {
			tenantID = assetSensor.AssetSensor.TenantID
		}

		macAddress := readingReq.MacAddress
		reading := &entity.IoTSensorReadingFlexible{
			TenantID:      tenantID,
			AssetSensorID: readingReq.AssetSensorID,
			SensorTypeID:  readingReq.SensorTypeID,
			MacAddress:    &macAddress,
			LocationID:    &locationID,
			LocationName:  &locationName,
			ReadingTime:   time.Now(),
		}

		// Use provided reading time if specified
		if readingReq.ReadingTime != nil {
			reading.ReadingTime = *readingReq.ReadingTime
		}

		readings = append(readings, reading)
	}

	// Save batch to database
	if err := s.iotSensorReadingRepo.CreateBatch(ctx, readings); err != nil {
		log.Printf("Error creating batch IoT sensor readings: %v", err)
		return nil, fmt.Errorf("failed to create batch IoT sensor readings: %w", err)
	}

	// Check thresholds for batch readings (non-blocking)
	go s.checkThresholdsForMultipleReadings(ctx, readings)

	// Convert to response DTOs
	for _, reading := range readings {
		responses = append(responses, s.toResponseDTO(reading))
	}

	log.Printf("Successfully created batch of %d IoT sensor readings", len(readings))
	return responses, nil
}

// GetIoTSensorReadingByID retrieves an IoT sensor reading by its ID
func (s *IoTSensorReadingService) GetIoTSensorReadingByID(ctx context.Context, id uuid.UUID) (*dto.IoTSensorReadingWithDetailsResponse, error) {
	reading, err := s.iotSensorReadingRepo.GetByID(ctx, id)
	if err != nil {
		log.Printf("Error getting IoT sensor reading: %v", err)
		return nil, fmt.Errorf("failed to get IoT sensor reading: %w", err)
	}

	if reading == nil {
		return nil, nil
	}

	// Get detailed reading with all related information
	detailedReading, err := s.iotSensorReadingRepo.GetLatestReading(ctx, reading.AssetSensorID)
	if err != nil {
		return nil, fmt.Errorf("failed to get detailed reading: %w", err)
	}

	return s.toDetailedResponseDTO(detailedReading), nil
}

// UpdateIoTSensorReading updates an existing IoT sensor reading
func (s *IoTSensorReadingService) UpdateIoTSensorReading(ctx context.Context, id uuid.UUID, req *dto.UpdateIoTSensorReadingRequest) (*dto.IoTSensorReadingResponse, error) {
	// Get existing reading
	existingReading, err := s.iotSensorReadingRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing reading: %w", err)
	}
	if existingReading == nil {
		return nil, common.NewValidationError("reading not found", nil)
	}

	// Update fields
	if req.MacAddress != nil {
		existingReading.MacAddress = req.MacAddress
	}
	if req.ReadingTime != nil {
		existingReading.ReadingTime = *req.ReadingTime
	}

	// Save changes
	if err := s.iotSensorReadingRepo.Update(ctx, existingReading); err != nil {
		log.Printf("Error updating IoT sensor reading: %v", err)
		return nil, fmt.Errorf("failed to update IoT sensor reading: %w", err)
	}

	log.Printf("Successfully updated IoT sensor reading: %s", id)
	return s.toResponseDTO(existingReading), nil
}

// DeleteIoTSensorReading deletes an IoT sensor reading
func (s *IoTSensorReadingService) DeleteIoTSensorReading(ctx context.Context, id uuid.UUID) error {
	if err := s.iotSensorReadingRepo.Delete(ctx, id); err != nil {
		log.Printf("Error deleting IoT sensor reading: %v", err)
		return fmt.Errorf("failed to delete IoT sensor reading: %w", err)
	}

	log.Printf("Successfully deleted IoT sensor reading: %s", id)
	return nil
}

// ListIoTSensorReadings retrieves IoT sensor readings with pagination and filtering
func (s *IoTSensorReadingService) ListIoTSensorReadings(ctx context.Context, req *dto.IoTSensorReadingListRequest) (*dto.IoTSensorReadingListResponse, error) {
	// Validate pagination
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}

	// Convert to repository request
	repoReq := repository.IoTSensorReadingListRequest{
		AssetSensorID: req.AssetSensorID,
		SensorTypeID:  req.SensorTypeID,
		MacAddress:    req.MacAddress,
		LocationID:    nil, // Location filtering is not supported in the current DTO
		FromTime:      req.FromTime,
		ToTime:        req.ToTime,
		Page:          req.Page,
		PageSize:      req.PageSize,
	}

	readings, total, err := s.iotSensorReadingRepo.List(ctx, repoReq)
	if err != nil {
		log.Printf("Error listing IoT sensor readings: %v", err)
		return nil, fmt.Errorf("failed to list IoT sensor readings: %w", err)
	}

	// Convert to DTOs
	var responseDTOs []dto.IoTSensorReadingWithDetailsResponse
	for _, reading := range readings {
		responseDTOs = append(responseDTOs, *s.toDetailedResponseDTO(reading))
	}

	// Calculate pagination info
	totalPages := int(math.Ceil(float64(total) / float64(req.PageSize)))

	return &dto.IoTSensorReadingListResponse{
		Readings:   responseDTOs,
		Page:       req.Page,
		Limit:      req.PageSize,
		Total:      int64(total),
		TotalPages: totalPages,
	}, nil
}

// GetReadingsByAssetSensor retrieves readings for a specific asset sensor
func (s *IoTSensorReadingService) GetReadingsByAssetSensor(ctx context.Context, assetSensorID uuid.UUID, limit int) ([]*dto.IoTSensorReadingWithDetailsResponse, error) {
	if limit <= 0 {
		limit = 100
	}
	if limit > 1000 {
		limit = 1000
	}

	readings, err := s.iotSensorReadingRepo.GetByAssetSensorID(ctx, assetSensorID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get readings by asset sensor: %w", err)
	}

	var responses []*dto.IoTSensorReadingWithDetailsResponse
	for _, reading := range readings {
		responses = append(responses, s.toDetailedResponseDTO(reading))
	}

	return responses, nil
}

// GetReadingsBySensorType retrieves readings for a specific sensor type
func (s *IoTSensorReadingService) GetReadingsBySensorType(ctx context.Context, sensorTypeID uuid.UUID, limit int) ([]*dto.IoTSensorReadingWithDetailsResponse, error) {
	if limit <= 0 {
		limit = 100
	}
	if limit > 1000 {
		limit = 1000
	}

	readings, err := s.iotSensorReadingRepo.GetBySensorTypeID(ctx, sensorTypeID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get readings by sensor type: %w", err)
	}

	var responses []*dto.IoTSensorReadingWithDetailsResponse
	for _, reading := range readings {
		responses = append(responses, s.toDetailedResponseDTO(reading))
	}

	return responses, nil
}

// GetReadingsByMacAddress retrieves readings for a specific MAC address
func (s *IoTSensorReadingService) GetReadingsByMacAddress(ctx context.Context, macAddress string, limit int) ([]*dto.IoTSensorReadingWithDetailsResponse, error) {
	if limit <= 0 {
		limit = 100
	}
	if limit > 1000 {
		limit = 1000
	}

	readings, err := s.iotSensorReadingRepo.GetByMacAddress(ctx, macAddress, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get readings by MAC address: %w", err)
	}

	var responses []*dto.IoTSensorReadingWithDetailsResponse
	for _, reading := range readings {
		responses = append(responses, s.toDetailedResponseDTO(reading))
	}

	return responses, nil
}

// GetLatestReading retrieves the most recent reading for an asset sensor
func (s *IoTSensorReadingService) GetLatestReading(ctx context.Context, assetSensorID uuid.UUID) (*dto.IoTSensorReadingWithDetailsResponse, error) {
	reading, err := s.iotSensorReadingRepo.GetLatestReading(ctx, assetSensorID)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest reading: %w", err)
	}

	if reading == nil {
		return nil, nil
	}

	return s.toDetailedResponseDTO(reading), nil
}

// GetReadingsInTimeRange retrieves readings within a time range
func (s *IoTSensorReadingService) GetReadingsInTimeRange(ctx context.Context, req *dto.GetReadingsInTimeRangeRequest) ([]*dto.IoTSensorReadingWithDetailsResponse, error) {
	// Validate time range
	if req.ToTime.Before(req.FromTime) {
		return nil, common.NewValidationError("to_time must be after from_time", nil)
	}

	// Default limit if not specified
	limit := req.Limit
	if limit <= 0 {
		limit = 1000
	}
	if limit > 10000 {
		limit = 10000
	}

	var readings []*repository.IoTSensorReadingWithDetails
	var err error

	if req.AssetSensorID != nil {
		readings, err = s.iotSensorReadingRepo.GetReadingsInTimeRange(ctx, *req.AssetSensorID, req.FromTime, req.ToTime)
	} else {
		// If no specific asset sensor, we'll need to implement a general time range query
		// For now, return error as the repository interface doesn't support this yet
		return nil, common.NewValidationError("asset_sensor_id is required for time range queries", nil)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get readings in time range: %w", err)
	}

	// Apply limit after fetching (ideally this should be done in the query)
	if len(readings) > limit {
		readings = readings[:limit]
	}

	var responses []*dto.IoTSensorReadingWithDetailsResponse
	for _, reading := range readings {
		responses = append(responses, s.toDetailedResponseDTO(reading))
	}

	return responses, nil
}

// GetAggregatedData retrieves aggregated sensor data for analytics
func (s *IoTSensorReadingService) GetAggregatedData(ctx context.Context, req *dto.GetAggregatedDataRequest) (*dto.GetAggregatedDataResponse, error) {
	// Validate time range
	if req.ToTime.Before(req.FromTime) {
		return nil, common.NewValidationError("to_time must be after from_time", nil)
	}

	// Default interval if not specified
	interval := req.Interval
	if interval == "" {
		interval = "hour"
	}

	// Validate interval
	validIntervals := map[string]bool{"hour": true, "day": true, "week": true, "month": true}
	if !validIntervals[interval] {
		return nil, common.NewValidationError("invalid interval, must be: hour, day, week, month", nil)
	}

	var aggregatedData []map[string]interface{}
	var err error

	if req.AssetSensorID != nil {
		aggregatedData, err = s.iotSensorReadingRepo.GetAggregatedData(ctx, *req.AssetSensorID, req.FromTime, req.ToTime, interval)
	} else {
		// If no specific asset sensor, we'll need to implement a general aggregation query
		return nil, common.NewValidationError("asset_sensor_id is required for aggregated data queries", nil)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get aggregated data: %w", err)
	}

	// Convert to response format
	var dataPoints []dto.AggregatedDataPoint
	totalCount := int64(0)

	for _, data := range aggregatedData {
		point := dto.AggregatedDataPoint{
			Data: data,
		}

		// Extract common fields if they exist
		if timeVal, ok := data["time"]; ok {
			if timeStr, ok := timeVal.(string); ok {
				if parsedTime, err := time.Parse(time.RFC3339, timeStr); err == nil {
					point.Time = parsedTime
				}
			}
		}

		if countVal, ok := data["count"]; ok {
			if count, ok := countVal.(int64); ok {
				point.Count = count
				totalCount += count
			}
		}

		// Extract averages, sums, mins, maxs if they exist
		if avgData, ok := data["averages"].(map[string]interface{}); ok {
			point.Averages = make(map[string]float64)
			for k, v := range avgData {
				if val, ok := v.(float64); ok {
					point.Averages[k] = val
				}
			}
		}

		dataPoints = append(dataPoints, point)
	}

	return &dto.GetAggregatedDataResponse{
		DataPoints:  dataPoints,
		TotalCount:  totalCount,
		FromTime:    req.FromTime,
		ToTime:      req.ToTime,
		Interval:    interval,
		AggregateBy: req.AggregateBy,
		RequestedAt: time.Now(),
	}, nil
}

// ValidateAndCreateReading validates measurement data against schemas and creates reading
func (s *IoTSensorReadingService) ValidateAndCreateReading(ctx context.Context, req *dto.ValidateAndCreateRequest) (*dto.ValidateAndCreateResponse, error) {
	// Create entity from request
	reading := &entity.IoTSensorReading{
		AssetSensorID: req.AssetSensorID,
		SensorTypeID:  req.SensorTypeID,
		ReadingTime:   time.Now(),
	}

	if req.ReadingTime != nil {
		reading.ReadingTime = *req.ReadingTime
	}

	// Get asset sensor to set tenant ID
	assetSensor, err := s.assetSensorRepo.GetByID(ctx, req.AssetSensorID)
	if err != nil {
		return nil, fmt.Errorf("failed to validate asset sensor: %w", err)
	}
	if assetSensor == nil {
		return nil, common.NewValidationError("asset sensor not found", nil)
	}

	reading.TenantID = assetSensor.AssetSensor.TenantID

	// Validate and create
	isValid, validationErrors, err := s.iotSensorReadingRepo.ValidateAndCreate(ctx, reading)
	if err != nil {
		return nil, fmt.Errorf("failed to validate and create reading: %w", err)
	}

	response := &dto.ValidateAndCreateResponse{
		IoTSensorReadingResponse: s.toResponseDTO(reading),
		IsValid:                  isValid,
	}

	// Convert validation errors to DTO format
	for _, errMsg := range validationErrors {
		response.ValidationErrors = append(response.ValidationErrors, dto.ValidationError{
			Message: errMsg,
		})
	}

	return response, nil
}

// DeleteReadingsByAssetSensor deletes all readings for an asset sensor
func (s *IoTSensorReadingService) DeleteReadingsByAssetSensor(ctx context.Context, assetSensorID uuid.UUID) error {
	if err := s.iotSensorReadingRepo.DeleteByAssetSensorID(ctx, assetSensorID); err != nil {
		log.Printf("Error deleting readings for asset sensor: %v", err)
		return fmt.Errorf("failed to delete readings for asset sensor: %w", err)
	}

	log.Printf("Successfully deleted readings for asset sensor: %s", assetSensorID)
	return nil
}

// Helper methods

// getLocationFromAssetSensor retrieves location information from asset sensor
func (s *IoTSensorReadingService) getLocationFromAssetSensor(ctx context.Context, assetSensorID uuid.UUID) (uuid.UUID, string, error) {
	// Get asset sensor details
	assetSensor, err := s.assetSensorRepo.GetByID(ctx, assetSensorID)
	if err != nil {
		return uuid.Nil, "", fmt.Errorf("failed to get asset sensor: %w", err)
	}
	if assetSensor == nil {
		return uuid.Nil, "", fmt.Errorf("asset sensor not found")
	}

	// Get asset details to retrieve location ID
	asset, err := s.assetRepo.GetByID(ctx, assetSensor.AssetSensor.AssetID)
	if err != nil {
		return uuid.Nil, "", fmt.Errorf("failed to get asset: %w", err)
	}
	if asset == nil {
		return uuid.Nil, "", fmt.Errorf("asset not found")
	}

	// Get location details to retrieve location name
	location, err := s.locationRepo.GetByID(ctx, asset.LocationID)
	if err != nil {
		return uuid.Nil, "", fmt.Errorf("failed to get location: %w", err)
	}
	if location == nil {
		return uuid.Nil, "", fmt.Errorf("location not found")
	}

	return asset.LocationID, location.Name, nil
}

// validateCreateRequest validates the create request
func (s *IoTSensorReadingService) validateCreateRequest(req *dto.CreateIoTSensorReadingRequest) error {
	if req.AssetSensorID == uuid.Nil {
		return common.NewValidationError("asset_sensor_id is required", nil)
	}

	if req.SensorTypeID == uuid.Nil {
		return common.NewValidationError("sensor_type_id is required", nil)
	}

	if req.MacAddress == "" {
		return common.NewValidationError("mac_address is required", nil)
	}

	return nil
}

// toResponseDTO converts entity to response DTO
func (s *IoTSensorReadingService) toResponseDTO(reading *entity.IoTSensorReadingFlexible) *dto.IoTSensorReadingResponse {
	response := &dto.IoTSensorReadingResponse{
		ID:            reading.ID,
		AssetSensorID: reading.AssetSensorID,
		SensorTypeID:  reading.SensorTypeID,
		ReadingTime:   reading.ReadingTime,
		CreatedAt:     reading.CreatedAt,
		UpdatedAt:     reading.UpdatedAt,
	}

	if reading.TenantID != nil {
		response.TenantID = *reading.TenantID
	}
	if reading.MacAddress != nil {
		response.MacAddress = *reading.MacAddress
	}
	if reading.LocationName != nil {
		response.Location = *reading.LocationName
	}

	// Add measurement data if available
	if reading.MeasurementType != "" {
		measurementValue := dto.MeasurementValue{}
		if reading.MeasurementLabel != nil {
			measurementValue.Label = *reading.MeasurementLabel
		}
		if reading.MeasurementUnit != nil {
			measurementValue.Unit = *reading.MeasurementUnit
		}

		// Set value based on type
		if reading.NumericValue != nil {
			measurementValue.Value = *reading.NumericValue
		} else if reading.TextValue != nil {
			measurementValue.Value = *reading.TextValue
		} else if reading.BooleanValue != nil {
			measurementValue.Value = *reading.BooleanValue
		}

		response.MeasurementData = map[string]dto.MeasurementValue{
			reading.MeasurementType: measurementValue,
		}
	}

	return response
}

// toDetailedResponseDTO converts repository model to detailed response DTO
func (s *IoTSensorReadingService) toDetailedResponseDTO(reading *repository.IoTSensorReadingWithDetails) *dto.IoTSensorReadingWithDetailsResponse {
	response := &dto.IoTSensorReadingWithDetailsResponse{
		IoTSensorReadingResponse: s.toResponseDTO(reading.IoTSensorReading),
	}

	// Copy asset sensor details
	response.AssetSensor.ID = reading.AssetSensor.ID
	response.AssetSensor.AssetID = reading.AssetSensor.AssetID
	response.AssetSensor.Name = reading.AssetSensor.Name
	response.AssetSensor.Status = reading.AssetSensor.Status
	response.AssetSensor.Configuration = reading.AssetSensor.Configuration

	// Copy sensor type details
	response.SensorType.ID = reading.SensorType.ID
	response.SensorType.Name = reading.SensorType.Name
	response.SensorType.Description = reading.SensorType.Description
	response.SensorType.Manufacturer = reading.SensorType.Manufacturer
	response.SensorType.Model = reading.SensorType.Model
	response.SensorType.Version = reading.SensorType.Version
	response.SensorType.IsActive = reading.SensorType.IsActive

	// Copy measurement types
	for _, mt := range reading.MeasurementTypes {
		measurementType := struct {
			ID               uuid.UUID       `json:"id"`
			Name             string          `json:"name"`
			Description      string          `json:"description"`
			PropertiesSchema json.RawMessage `json:"properties_schema"`
			UIConfiguration  json.RawMessage `json:"ui_configuration"`
			Version          string          `json:"version"`
			IsActive         bool            `json:"is_active"`
			Fields           []struct {
				ID          uuid.UUID `json:"id"`
				Name        string    `json:"name"`
				Label       string    `json:"label"`
				Description *string   `json:"description"`
				DataType    string    `json:"data_type"`
				Required    bool      `json:"required"`
				Unit        *string   `json:"unit"`
				Min         *float64  `json:"min"`
				Max         *float64  `json:"max"`
			} `json:"fields"`
		}{
			ID:               mt.ID,
			Name:             mt.Name,
			Description:      mt.Description,
			PropertiesSchema: mt.PropertiesSchema,
			UIConfiguration:  mt.UIConfiguration,
			Version:          mt.Version,
			IsActive:         mt.IsActive,
		}

		// Copy fields
		for _, field := range mt.Fields {
			measurementType.Fields = append(measurementType.Fields, struct {
				ID          uuid.UUID `json:"id"`
				Name        string    `json:"name"`
				Label       string    `json:"label"`
				Description *string   `json:"description"`
				DataType    string    `json:"data_type"`
				Required    bool      `json:"required"`
				Unit        *string   `json:"unit"`
				Min         *float64  `json:"min"`
				Max         *float64  `json:"max"`
			}{
				ID:          field.ID,
				Name:        field.Name,
				Label:       field.Label,
				Description: field.Description,
				DataType:    field.DataType,
				Required:    field.Required,
				Unit:        field.Unit,
				Min:         field.Min,
				Max:         field.Max,
			})
		}

		response.MeasurementTypes = append(response.MeasurementTypes, measurementType)
	}

	return response
}

// Manual Data Input Methods

// CreateFromJSONString creates a sensor reading from raw JSON string
func (s *IoTSensorReadingService) CreateFromJSONString(ctx context.Context, jsonData string) (*dto.IoTSensorReadingResponse, error) {
	var req dto.CreateIoTSensorReadingRequest
	if err := json.Unmarshal([]byte(jsonData), &req); err != nil {
		log.Printf("Error parsing JSON data: %v", err)
		return nil, fmt.Errorf("invalid JSON format: %w", err)
	}

	return s.CreateIoTSensorReading(ctx, &req)
}

// CreateFromJSONFile creates a sensor reading from JSON file content
func (s *IoTSensorReadingService) CreateFromJSONFile(ctx context.Context, fileContent []byte) (*dto.IoTSensorReadingResponse, error) {
	var req dto.CreateIoTSensorReadingRequest
	if err := json.Unmarshal(fileContent, &req); err != nil {
		log.Printf("Error parsing JSON file: %v", err)
		return nil, fmt.Errorf("invalid JSON file format: %w", err)
	}

	return s.CreateIoTSensorReading(ctx, &req)
}

// CreateBatchFromJSONString creates multiple sensor readings from JSON string
func (s *IoTSensorReadingService) CreateBatchFromJSONString(ctx context.Context, jsonData string) ([]*dto.IoTSensorReadingResponse, error) {
	var req dto.CreateBatchIoTSensorReadingRequest
	if err := json.Unmarshal([]byte(jsonData), &req); err != nil {
		log.Printf("Error parsing batch JSON data: %v", err)
		return nil, fmt.Errorf("invalid batch JSON format: %w", err)
	}

	return s.CreateBatchIoTSensorReading(ctx, &req)
}

// CreateSimpleReading creates a sensor reading with minimal required fields
func (s *IoTSensorReadingService) CreateSimpleReading(ctx context.Context, assetSensorID, sensorTypeID uuid.UUID, macAddress string, measurementData map[string]interface{}) (*dto.IoTSensorReadingResponse, error) {
	req := &dto.CreateIoTSensorReadingRequest{
		AssetSensorID: assetSensorID,
		SensorTypeID:  sensorTypeID,
		MacAddress:    macAddress,
		ReadingTime:   nil, // Will use current time
	}

	return s.CreateIoTSensorReading(ctx, req)
}

// CreateDummyReading creates a dummy sensor reading for testing purposes
func (s *IoTSensorReadingService) CreateDummyReading(ctx context.Context, sensorType string) (*dto.IoTSensorReadingResponse, error) {
	// Generate dummy data based on sensor type
	var macAddress string

	switch sensorType {
	case "temperature":
		macAddress = fmt.Sprintf("temp_%d", time.Now().Unix()%1000)
	case "vibration":
		macAddress = fmt.Sprintf("vib_%d", time.Now().Unix()%1000)
	case "pressure":
		macAddress = fmt.Sprintf("press_%d", time.Now().Unix()%1000)
	case "gas":
		macAddress = fmt.Sprintf("gas_%d", time.Now().Unix()%1000)
	default:
		macAddress = fmt.Sprintf("sensor_%d", time.Now().Unix()%1000)
	}

	// You'll need to provide actual UUIDs from your database
	// For now, this will generate random UUIDs (you should replace with actual values)
	assetSensorID := uuid.New() // Replace with actual asset sensor ID
	sensorTypeID := uuid.New()  // Replace with actual sensor type ID

	return s.CreateSimpleReading(ctx, assetSensorID, sensorTypeID, macAddress, nil)
}

// CreateMultipleDummyReadings creates multiple dummy readings for different sensor types
func (s *IoTSensorReadingService) CreateMultipleDummyReadings(ctx context.Context, count int, sensorTypes []string) ([]*dto.IoTSensorReadingResponse, error) {
	var responses []*dto.IoTSensorReadingResponse

	for i := 0; i < count; i++ {
		sensorType := sensorTypes[i%len(sensorTypes)] // Cycle through sensor types

		// Add delay to create different timestamps
		time.Sleep(100 * time.Millisecond)

		response, err := s.CreateDummyReading(ctx, sensorType)
		if err != nil {
			log.Printf("Error creating dummy reading %d: %v", i+1, err)
			continue
		}

		responses = append(responses, response)
	}

	return responses, nil
}

// GetJSONTemplate creates a template JSON for different sensor types
func (s *IoTSensorReadingService) GetJSONTemplate(sensorType string) (string, error) {
	var template map[string]interface{}

	switch sensorType {
	case "temperature":
		template = map[string]interface{}{
			"asset_sensor_id": "00000000-0000-0000-0000-000000000000", // Replace with actual UUID
			"sensor_type_id":  "00000000-0000-0000-0000-000000000000", // Replace with actual UUID
			"mac_address":     "temp_sensor_001",
			"reading_time":    time.Now().Format(time.RFC3339),
		}

	case "vibration":
		template = map[string]interface{}{
			"asset_sensor_id": "00000000-0000-0000-0000-000000000000", // Replace with actual UUID
			"sensor_type_id":  "00000000-0000-0000-0000-000000000000", // Replace with actual UUID
			"mac_address":     "vib_sensor_001",
			"reading_time":    time.Now().Format(time.RFC3339),
		}

	case "pressure":
		template = map[string]interface{}{
			"asset_sensor_id": "00000000-0000-0000-0000-000000000000", // Replace with actual UUID
			"sensor_type_id":  "00000000-0000-0000-0000-000000000000", // Replace with actual UUID
			"mac_address":     "press_sensor_001",
			"reading_time":    time.Now().Format(time.RFC3339),
		}

	case "gas":
		template = map[string]interface{}{
			"asset_sensor_id": "00000000-0000-0000-0000-000000000000", // Replace with actual UUID
			"sensor_type_id":  "00000000-0000-0000-0000-000000000000", // Replace with actual UUID
			"mac_address":     "gas_sensor_001",
			"reading_time":    time.Now().Format(time.RFC3339),
		}

	default:
		return "", fmt.Errorf("unsupported sensor type: %s", sensorType)
	}

	jsonBytes, err := json.MarshalIndent(template, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal template: %w", err)
	}

	return string(jsonBytes), nil
}

// GetBatchJSONTemplate creates a batch template with multiple sensor readings
func (s *IoTSensorReadingService) GetBatchJSONTemplate(sensorTypes []string, count int) (string, error) {
	var readings []map[string]interface{}

	for i := 0; i < count; i++ {
		sensorType := sensorTypes[i%len(sensorTypes)]

		templateStr, err := s.GetJSONTemplate(sensorType)
		if err != nil {
			continue
		}

		var reading map[string]interface{}
		if err := json.Unmarshal([]byte(templateStr), &reading); err != nil {
			continue
		}

		// Modify MAC address to be unique
		reading["mac_address"] = fmt.Sprintf("%s_%03d", reading["mac_address"], i+1)

		readings = append(readings, reading)
	}

	batch := map[string]interface{}{
		"readings": readings,
	}

	jsonBytes, err := json.MarshalIndent(batch, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal batch template: %w", err)
	}

	return string(jsonBytes), nil
}

// CreateFlexibleIoTSensorReading creates a flexible IoT sensor reading
func (s *IoTSensorReadingService) CreateFlexibleIoTSensorReading(ctx context.Context, req *dto.FlexibleIoTSensorReadingRequest) (*dto.IoTSensorReadingResponse, error) {
	// Validate basic required fields
	if req.AssetSensorID == uuid.Nil {
		return nil, common.NewValidationError("asset_sensor_id is required", nil)
	}

	// Validate asset sensor exists and get its tenant_id
	assetSensor, err := s.assetSensorRepo.GetByID(ctx, req.AssetSensorID)
	if err != nil {
		log.Printf("Error validating asset sensor: %v", err)
		return nil, fmt.Errorf("failed to validate asset sensor: %w", err)
	}
	if assetSensor == nil {
		return nil, common.NewValidationError("asset sensor not found", nil)
	}

	// Get measurement types and fields for validation
	measurementTypesQuery := `
		SELECT 
			smt.id, smt.name, smt.description, smt.properties_schema,
			COALESCE(
				(SELECT json_agg(
					json_build_object(
						'id', smf.id,
						'name', smf.name,
						'label', smf.label,
						'description', smf.description,
						'data_type', smf.data_type,
						'required', smf.required,
						'unit', smf.unit,
						'min', smf.min,
						'max', smf.max
					)
				)
				FROM sensor_measurement_fields smf
				WHERE smf.sensor_measurement_type_id = smt.id), '[]'::json
			) as fields
		FROM sensor_measurement_types smt
		WHERE smt.sensor_type_id = $1 AND smt.is_active = true`

	rows, err := s.iotSensorReadingRepo.GetDB().QueryContext(ctx, measurementTypesQuery, req.SensorTypeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get measurement types and fields: %w", err)
	}
	defer rows.Close()

	// Map to store valid fields by name
	validFields := make(map[string]struct {
		Unit string
		Min  *float64
		Max  *float64
	})

	// Process measurement types and fields
	for rows.Next() {
		var mt struct {
			ID               uuid.UUID
			Name             string
			Description      string
			PropertiesSchema json.RawMessage
			Fields           []struct {
				ID          uuid.UUID `json:"id"`
				Name        string    `json:"name"`
				Label       string    `json:"label"`
				Description *string   `json:"description"`
				DataType    string    `json:"data_type"`
				Required    bool      `json:"required"`
				Unit        *string   `json:"unit"`
				Min         *float64  `json:"min"`
				Max         *float64  `json:"max"`
			} `json:"fields"`
		}

		var fieldsJSON []byte
		err := rows.Scan(&mt.ID, &mt.Name, &mt.Description, &mt.PropertiesSchema, &fieldsJSON)
		if err != nil {
			return nil, fmt.Errorf("failed to scan measurement type: %w", err)
		}

		// Parse fields JSON
		if len(fieldsJSON) > 0 && string(fieldsJSON) != "[]" {
			if err := json.Unmarshal(fieldsJSON, &mt.Fields); err != nil {
				return nil, fmt.Errorf("failed to parse measurement fields: %w", err)
			}

			// Add fields to valid fields map
			for _, field := range mt.Fields {
				validFields[field.Name] = struct {
					Unit string
					Min  *float64
					Max  *float64
				}{
					Unit: func() string {
						if field.Unit != nil {
							return *field.Unit
						}
						// Auto-fill unit based on label
						switch strings.ToLower(field.Label) {
						case "temperature":
							return "°C"
						case "humidity":
							return "%"
						case "pressure":
							return "Pa"
						case "voltage":
							return "V"
						case "current":
							return "A"
						case "power":
							return "W"
						case "energy":
							return "kWh"
						case "speed":
							return "m/s"
						case "distance":
							return "m"
						case "weight":
							return "kg"
						case "volume":
							return "m³"
						case "flow":
							return "m³/s"
						case "concentration":
							return "ppm"
						case "ph":
							return "pH"
						case "turbidity":
							return "NTU"
						case "conductivity":
							return "μS/cm"
						default:
							return ""
						}
					}(),
					Min: field.Min,
					Max: field.Max,
				}
			}
		}
	}

	// Validate measurement data against valid fields
	var warnings []string
	validMeasurementData := make(map[string]dto.MeasurementValue)
	for key, measurement := range req.MeasurementData {
		field, exists := validFields[key]
		if !exists {
			warnings = append(warnings, fmt.Sprintf("invalid measurement field: %s", key))
			continue
		}

		// Validate numeric values against min/max if applicable
		if numValue, ok := measurement.Value.(float64); ok {
			if field.Min != nil && numValue < *field.Min {
				warnings = append(warnings, fmt.Sprintf("value for %s is below minimum allowed value", key))
				continue
			}
			if field.Max != nil && numValue > *field.Max {
				warnings = append(warnings, fmt.Sprintf("value for %s is above maximum allowed value", key))
				continue
			}
		}

		// Auto-fill unit if not provided
		if measurement.Unit == "" && field.Unit != "" {
			measurement.Unit = field.Unit
		}

		validMeasurementData[key] = measurement
	}

	if len(validMeasurementData) == 0 {
		return nil, common.NewValidationError(strings.Join(warnings, "; "), nil)
	}

	// Create flexible readings for each valid measurement
	var flexibleReadings []*entity.IoTSensorReadingFlexible
	for key, measurement := range validMeasurementData {
		flexibleReading := &entity.IoTSensorReadingFlexible{
			ID:            uuid.New(),
			TenantID:      assetSensor.AssetSensor.TenantID,
			AssetSensorID: req.AssetSensorID,
			SensorTypeID:  req.SensorTypeID,
			ReadingTime:   time.Now(),
		}

		// Set reading time if provided
		if req.ReadingTime != nil {
			flexibleReading.ReadingTime = *req.ReadingTime
		}

		// Set MAC address if provided
		if req.MacAddress != "" {
			macAddressCopy := req.MacAddress
			flexibleReading.MacAddress = &macAddressCopy
		}

		// Set measurement type and label
		flexibleReading.MeasurementType = key
		labelCopy := measurement.Label
		flexibleReading.MeasurementLabel = &labelCopy

		// Set unit if provided
		if measurement.Unit != "" {
			unitCopy := measurement.Unit
			flexibleReading.MeasurementUnit = &unitCopy
		}

		// Set value based on type
		switch v := measurement.Value.(type) {
		case float64:
			valueCopy := v
			flexibleReading.NumericValue = &valueCopy
			log.Printf("Set numeric value: %f for measurement type: %s", v, key)
		case string:
			valueCopy := v
			flexibleReading.TextValue = &valueCopy
			log.Printf("Set text value: %s for measurement type: %s", v, key)
		case bool:
			valueCopy := v
			flexibleReading.BooleanValue = &valueCopy
			log.Printf("Set boolean value: %t for measurement type: %s", v, key)
		default:
			log.Printf("Warning: Unknown value type %T for measurement type: %s", v, key)
		}

		flexibleReadings = append(flexibleReadings, flexibleReading)
	}

	log.Printf("Created %d flexible reading entities", len(flexibleReadings))

	// Store all flexible readings using batch create
	if err := s.iotSensorReadingRepo.CreateFlexibleBatch(ctx, flexibleReadings); err != nil {
		log.Printf("Error creating flexible IoT sensor readings: %v", err)
		return nil, fmt.Errorf("failed to create flexible IoT sensor readings: %w", err)
	}

	log.Printf("Successfully created %d flexible IoT sensor readings", len(flexibleReadings))

	// Check thresholds for flexible readings (non-blocking)
	go s.checkThresholdsForMultipleReadings(ctx, flexibleReadings)

	// Convert to response using the first reading as base (all have same basic info)
	if len(flexibleReadings) > 0 {
		resp := s.toResponseDTO(flexibleReadings[0])
		if len(warnings) > 0 {
			resp.Message = "Flexible IoT sensor reading created successfully (some fields skipped)"
			resp.Warnings = warnings
		}
		return resp, nil
	}

	return nil, fmt.Errorf("no measurements were processed")
}

// GetFlexibleIoTSensorReading gets a flexible IoT sensor reading by ID
func (s *IoTSensorReadingService) GetFlexibleIoTSensorReading(ctx context.Context, id uuid.UUID) (*dto.FlexibleIoTSensorReadingResponse, error) {
	flexibleReading, err := s.iotSensorReadingRepo.GetByID(ctx, id)
	if err != nil {
		log.Printf("Error getting flexible IoT sensor reading: %v", err)
		return nil, fmt.Errorf("failed to get flexible IoT sensor reading: %w", err)
	}

	if flexibleReading == nil {
		return nil, common.NewValidationError("flexible IoT sensor reading not found", nil)
	}

	// Convert to response
	baseResponse := s.toResponseDTO(flexibleReading)

	// Convert measurement data to proper format
	measurementData := make(map[string]dto.MeasurementValue)

	// Add measurement data from the reading
	measurementValue := dto.MeasurementValue{}
	if flexibleReading.MeasurementLabel != nil {
		measurementValue.Label = *flexibleReading.MeasurementLabel
	}
	if flexibleReading.MeasurementUnit != nil {
		measurementValue.Unit = *flexibleReading.MeasurementUnit
	}

	// Set value based on type
	if flexibleReading.NumericValue != nil {
		measurementValue.Value = *flexibleReading.NumericValue
	} else if flexibleReading.TextValue != nil {
		measurementValue.Value = *flexibleReading.TextValue
	} else if flexibleReading.BooleanValue != nil {
		measurementValue.Value = *flexibleReading.BooleanValue
	}

	measurementData[flexibleReading.MeasurementType] = measurementValue

	response := &dto.FlexibleIoTSensorReadingResponse{
		IoTSensorReadingResponse: baseResponse,
		MeasurementData:          measurementData,
	}

	return response, nil
}

// CreateBulkFlexibleIoTSensorReadings creates multiple flexible IoT sensor readings in batch
func (s *IoTSensorReadingService) CreateBulkFlexibleIoTSensorReadings(ctx context.Context, requests []*dto.FlexibleIoTSensorReadingRequest) ([]*dto.IoTSensorReadingResponse, error) {
	if len(requests) == 0 {
		return nil, common.NewValidationError("at least one reading is required", nil)
	}
	if len(requests) > 1000 {
		return nil, common.NewValidationError("maximum 1000 readings allowed per batch", nil)
	}

	var readings []*entity.IoTSensorReadingFlexible
	var responses []*dto.IoTSensorReadingResponse

	now := time.Now()

	for i, req := range requests {
		// Validate basic required fields
		if req.AssetSensorID == uuid.Nil {
			return nil, fmt.Errorf("validation error for reading %d: asset_sensor_id is required", i)
		}

		// Validate asset sensor exists and get its tenant_id
		assetSensor, err := s.assetSensorRepo.GetByID(ctx, req.AssetSensorID)
		if err != nil {
			return nil, fmt.Errorf("failed to validate asset sensor for reading %d: %w", i, err)
		}
		if assetSensor == nil {
			return nil, fmt.Errorf("asset sensor not found for reading %d", i)
		}

		// Get location information from asset
		locationID, locationName, err := s.getLocationFromAssetSensor(ctx, req.AssetSensorID)
		if err != nil {
			log.Printf("Warning: Failed to get location for reading %d: %v", i, err)
			// Continue without location
		}

		// Create flexible reading entity
		flexibleReading := &entity.IoTSensorReadingFlexible{
			ID:            uuid.New(),
			AssetSensorID: req.AssetSensorID,
			SensorTypeID:  req.SensorTypeID,
			ReadingTime:   now,
			CreatedAt:     now,
		}

		// Set location fields if available
		if locationID != uuid.Nil {
			locationIDCopy := locationID
			flexibleReading.LocationID = &locationIDCopy
		}
		if locationName != "" {
			locationNameCopy := locationName
			flexibleReading.LocationName = &locationNameCopy
		}

		// Set tenant ID
		if assetSensor.AssetSensor.TenantID != nil {
			flexibleReading.TenantID = assetSensor.AssetSensor.TenantID
		}

		// Set MAC address if provided
		if req.MacAddress != "" {
			macAddressCopy := req.MacAddress
			flexibleReading.MacAddress = &macAddressCopy
		}

		// Convert measurement data to proper format
		for key, measurement := range req.MeasurementData {
			// Set measurement type
			flexibleReading.MeasurementType = key

			// Set label if provided
			if measurement.Label != "" {
				labelCopy := measurement.Label
				flexibleReading.MeasurementLabel = &labelCopy
			}

			// Set unit if provided
			if measurement.Unit != "" {
				unitCopy := measurement.Unit
				flexibleReading.MeasurementUnit = &unitCopy
			}

			// Set value based on type
			switch v := measurement.Value.(type) {
			case float64:
				valueCopy := v
				flexibleReading.NumericValue = &valueCopy
			case string:
				valueCopy := v
				flexibleReading.TextValue = &valueCopy
			case bool:
				valueCopy := v
				flexibleReading.BooleanValue = &valueCopy
			}
		}

		readings = append(readings, flexibleReading)
	}

	// Store the flexible readings in batch
	if err := s.iotSensorReadingRepo.CreateFlexibleBatch(ctx, readings); err != nil {
		log.Printf("Error creating flexible IoT sensor readings in batch: %v", err)
		return nil, fmt.Errorf("failed to create flexible IoT sensor readings in batch: %w", err)
	}

	// Convert to responses
	for _, reading := range readings {
		responses = append(responses, s.toResponseDTO(reading))
	}

	return responses, nil
}

// ParseTextToFlexibleReading parses text data into a flexible IoT sensor reading
func (s *IoTSensorReadingService) ParseTextToFlexibleReading(ctx context.Context, req *dto.TextToJSONRequest) (*dto.FlexibleIoTSensorReadingResponse, error) {
	// Basic validation
	if req.TextData == "" {
		return nil, common.NewValidationError("text_data is required", nil)
	}

	// Parse text data into key-value pairs
	// This is a simple implementation - you might want to enhance this based on your needs
	lines := strings.Split(req.TextData, "\n")
	measurementData := make(map[string]dto.MeasurementValue)

	for _, line := range lines {
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Try to parse as number first
		if num, err := strconv.ParseFloat(value, 64); err == nil {
			measurementData[key] = dto.MeasurementValue{
				Value: num,
			}
		} else if boolVal, err := strconv.ParseBool(value); err == nil {
			measurementData[key] = dto.MeasurementValue{
				Value: boolVal,
			}
		} else {
			measurementData[key] = dto.MeasurementValue{
				Value: value,
			}
		}
	}

	// Parse UUIDs
	assetSensorID, err := uuid.Parse(req.AssetSensorID)
	if err != nil {
		return nil, common.NewValidationError("invalid asset_sensor_id format", nil)
	}

	sensorTypeID, err := uuid.Parse(req.SensorTypeID)
	if err != nil {
		return nil, common.NewValidationError("invalid sensor_type_id format", nil)
	}

	// Create flexible reading request
	flexibleReq := &dto.FlexibleIoTSensorReadingRequest{
		AssetSensorID:   assetSensorID,
		SensorTypeID:    sensorTypeID,
		MacAddress:      req.MacAddress,
		MeasurementData: measurementData,
	}

	// Create the reading
	reading, err := s.CreateFlexibleIoTSensorReading(ctx, flexibleReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create reading from text: %w", err)
	}

	// Convert to flexible response
	return &dto.FlexibleIoTSensorReadingResponse{
		IoTSensorReadingResponse: reading,
		MeasurementData:          measurementData,
	}, nil
}

// checkThresholdsForReading checks if a sensor reading breaches any thresholds and creates alerts
func (s *IoTSensorReadingService) checkThresholdsForReading(
	ctx context.Context,
	reading *entity.IoTSensorReadingFlexible,
) {
	// Skip threshold checking if service is not available
	if s.sensorThresholdService == nil {
		return
	}

	// Only check numeric values for threshold breaches
	if reading.NumericValue == nil {
		return
	}

	// Get measurement types for this sensor type
	measurementTypes, err := s.sensorMeasurementTypeRepo.GetBySensorTypeID(ctx, reading.SensorTypeID)
	if err != nil {
		log.Printf("Failed to get measurement types for threshold checking: %v", err)
		return
	}

	// Find the measurement type that matches the reading's measurement type
	var targetMeasurementTypeID uuid.UUID
	for _, mt := range measurementTypes {
		if mt.Name == reading.MeasurementType {
			targetMeasurementTypeID = mt.ID
			break
		}
	}

	// If no matching measurement type found, skip threshold checking
	if targetMeasurementTypeID == uuid.Nil {
		log.Printf("No measurement type found for '%s', skipping threshold check", reading.MeasurementType)
		return
	}

	// Get asset information for alert context
	assetSensor, err := s.assetSensorRepo.GetByID(ctx, reading.AssetSensorID)
	if err != nil {
		log.Printf("Failed to get asset sensor for threshold checking: %v", err)
		return
	}

	if assetSensor == nil {
		log.Printf("Asset sensor not found for threshold checking")
		return
	}

	// Convert flexible reading to IoTSensorReading
	iotReading := &entity.IoTSensorReading{
		ID:            reading.ID,
		TenantID:      reading.TenantID,
		AssetSensorID: reading.AssetSensorID,
		SensorTypeID:  reading.SensorTypeID,
		ReadingTime:   reading.ReadingTime,
		CreatedAt:     reading.CreatedAt,
		UpdatedAt:     reading.UpdatedAt,
	}

	// Check thresholds for this measurement value
	err = s.sensorThresholdService.CheckThresholdsForValue(
		ctx,
		iotReading,
		reading.MeasurementType,
		*reading.NumericValue,
	)

	if err != nil {
		log.Printf("Error checking thresholds for reading %s: %v", reading.ID, err)
		return
	}

	log.Printf("Successfully checked thresholds for sensor reading %s (value: %.2f, type: %s)",
		reading.ID, *reading.NumericValue, reading.MeasurementType)
}

// checkThresholdsForMultipleReadings checks thresholds for multiple readings in batch
func (s *IoTSensorReadingService) checkThresholdsForMultipleReadings(
	ctx context.Context,
	readings []*entity.IoTSensorReadingFlexible,
) {
	// Skip if threshold service is not available
	if s.sensorThresholdService == nil {
		return
	}

	// Process each reading for threshold checking
	for _, reading := range readings {
		// Use a separate goroutine for non-blocking threshold checking
		go func(r *entity.IoTSensorReadingFlexible) {
			s.checkThresholdsForReading(ctx, r)
		}(reading)
	}
}
