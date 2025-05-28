package service

import (
	"be-lecsens/asset_management/data-layer/entity"
	"be-lecsens/asset_management/data-layer/repository"
	"be-lecsens/asset_management/helpers/dto"
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

// IoTSensorReadingService handles business logic for IoT sensor readings
type IoTSensorReadingService struct {
	iotSensorReadingRepo repository.IoTSensorReadingRepository
}

// NewIoTSensorReadingService creates a new instance of IoTSensorReadingService
func NewIoTSensorReadingService(iotSensorReadingRepo repository.IoTSensorReadingRepository) *IoTSensorReadingService {
	return &IoTSensorReadingService{
		iotSensorReadingRepo: iotSensorReadingRepo,
	}
}

// CreateReading creates a new IoT sensor reading
func (s *IoTSensorReadingService) CreateReading(ctx context.Context, req *dto.CreateIoTSensorReadingRequest) (*dto.IoTSensorReadingResponse, error) {
	// Convert DTO to entity
	reading := &entity.IoTSensorReading{
		ID:              uuid.New(),
		SensorTypeID:    req.SensorTypeID,
		MacAddress:      req.MacAddress,
		Location:        req.Location,
		MeasurementData: req.MeasurementData,
		DataX:           req.DataX,
		DataY:           req.DataY,
		PeakX:           req.PeakX,
		PeakY:           req.PeakY,
		PPM:             0,
		Label:           "",
		RawData:         req.RawData,
	}

	// Handle nullable UUID fields - convert from pointer to value
	if req.TenantID != nil {
		reading.TenantID = *req.TenantID
	}

	if req.AssetSensorID != nil {
		reading.AssetSensorID = *req.AssetSensorID
	}

	// Set reading time
	if req.ReadingTime != nil {
		reading.ReadingTime = *req.ReadingTime
	} else {
		reading.ReadingTime = time.Now()
	}

	// Set PPM if provided
	if req.PPM != nil {
		reading.PPM = *req.PPM
	}

	// Set label if provided
	if req.Label != nil {
		reading.Label = *req.Label
	}

	// Extract standard fields from measurement data
	standardFields, err := s.extractStandardFields(req.MeasurementData)
	if err != nil {
		return nil, errors.New("failed to extract standard fields from measurement data")
	}
	reading.StandardFields = standardFields

	// Create the reading
	if err := s.iotSensorReadingRepo.Create(ctx, reading); err != nil {
		return nil, err
	}

	// Convert entity to response DTO
	return s.entityToResponse(reading), nil
}

// GetReading retrieves an IoT sensor reading by ID
func (s *IoTSensorReadingService) GetReading(ctx context.Context, id uuid.UUID) (*dto.IoTSensorReadingResponse, error) {
	reading, err := s.iotSensorReadingRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if reading == nil {
		return nil, errors.New("reading not found")
	}

	return s.entityToResponse(reading), nil
}

// ListReadings retrieves a list of IoT sensor readings with pagination and filtering
func (s *IoTSensorReadingService) ListReadings(ctx context.Context, params *dto.IoTSensorReadingQueryParams) (*dto.IoTSensorReadingListResponse, error) {
	// Set default pagination
	if params.Page == 0 {
		params.Page = 1
	}
	if params.Limit == 0 {
		params.Limit = 10
	}

	readings, err := s.iotSensorReadingRepo.List(ctx, params.TenantID, params.Page, params.Limit)
	if err != nil {
		return nil, err
	}

	// Get total count for pagination
	var total int64
	if params.TenantID != nil {
		total, err = s.iotSensorReadingRepo.GetCountByTenant(ctx, *params.TenantID)
		if err != nil {
			return nil, err
		}
	}

	// Convert entities to response DTOs
	responseReadings := make([]dto.IoTSensorReadingResponse, len(readings))
	for i, reading := range readings {
		responseReadings[i] = *s.entityToResponse(reading)
	}

	// Calculate total pages
	totalPages := int(total) / params.Limit
	if int(total)%params.Limit > 0 {
		totalPages++
	}

	return &dto.IoTSensorReadingListResponse{
		Readings:   responseReadings,
		Page:       params.Page,
		Limit:      params.Limit,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

// ListReadingsByAssetSensor retrieves IoT sensor readings for a specific asset sensor
func (s *IoTSensorReadingService) ListReadingsByAssetSensor(ctx context.Context, assetSensorID uuid.UUID, page, limit int) (*dto.IoTSensorReadingListResponse, error) {
	readings, err := s.iotSensorReadingRepo.ListByAssetSensor(ctx, assetSensorID, page, limit)
	if err != nil {
		return nil, err
	}

	// Convert entities to response DTOs
	responseReadings := make([]dto.IoTSensorReadingResponse, len(readings))
	for i, reading := range readings {
		responseReadings[i] = *s.entityToResponse(reading)
	}

	return &dto.IoTSensorReadingListResponse{
		Readings:   responseReadings,
		Page:       page,
		Limit:      limit,
		Total:      int64(len(readings)),
		TotalPages: 1, // This would need proper count implementation
	}, nil
}

// ListReadingsBySensorType retrieves IoT sensor readings for a specific sensor type
func (s *IoTSensorReadingService) ListReadingsBySensorType(ctx context.Context, sensorTypeID uuid.UUID, page, limit int) (*dto.IoTSensorReadingListResponse, error) {
	readings, err := s.iotSensorReadingRepo.ListBySensorType(ctx, sensorTypeID, page, limit)
	if err != nil {
		return nil, err
	}

	// Convert entities to response DTOs
	responseReadings := make([]dto.IoTSensorReadingResponse, len(readings))
	for i, reading := range readings {
		responseReadings[i] = *s.entityToResponse(reading)
	}

	return &dto.IoTSensorReadingListResponse{
		Readings:   responseReadings,
		Page:       page,
		Limit:      limit,
		Total:      int64(len(readings)),
		TotalPages: 1, // This would need proper count implementation
	}, nil
}

// ListReadingsByMacAddress retrieves IoT sensor readings for a specific MAC address
func (s *IoTSensorReadingService) ListReadingsByMacAddress(ctx context.Context, macAddress string, page, limit int) (*dto.IoTSensorReadingListResponse, error) {
	readings, err := s.iotSensorReadingRepo.ListByMacAddress(ctx, macAddress, page, limit)
	if err != nil {
		return nil, err
	}

	// Convert entities to response DTOs
	responseReadings := make([]dto.IoTSensorReadingResponse, len(readings))
	for i, reading := range readings {
		responseReadings[i] = *s.entityToResponse(reading)
	}

	return &dto.IoTSensorReadingListResponse{
		Readings:   responseReadings,
		Page:       page,
		Limit:      limit,
		Total:      int64(len(readings)),
		TotalPages: 1, // This would need proper count implementation
	}, nil
}

// ListReadingsByTimeRange retrieves IoT sensor readings within a specific time range
func (s *IoTSensorReadingService) ListReadingsByTimeRange(ctx context.Context, startTime, endTime time.Time, page, limit int) (*dto.IoTSensorReadingListResponse, error) {
	readings, err := s.iotSensorReadingRepo.ListByTimeRange(ctx, startTime, endTime, page, limit)
	if err != nil {
		return nil, err
	}

	// Convert entities to response DTOs
	responseReadings := make([]dto.IoTSensorReadingResponse, len(readings))
	for i, reading := range readings {
		responseReadings[i] = *s.entityToResponse(reading)
	}

	return &dto.IoTSensorReadingListResponse{
		Readings:   responseReadings,
		Page:       page,
		Limit:      limit,
		Total:      int64(len(readings)),
		TotalPages: 1, // This would need proper count implementation
	}, nil
}

// UpdateReading updates an existing IoT sensor reading
func (s *IoTSensorReadingService) UpdateReading(ctx context.Context, id uuid.UUID, req *dto.UpdateIoTSensorReadingRequest) (*dto.IoTSensorReadingResponse, error) {
	// Get existing reading
	existingReading, err := s.iotSensorReadingRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if existingReading == nil {
		return nil, errors.New("reading not found")
	}

	// Update only provided fields
	updatedReading := *existingReading

	if req.TenantID != nil {
		updatedReading.TenantID = *req.TenantID
	}

	if req.AssetSensorID != nil {
		updatedReading.AssetSensorID = *req.AssetSensorID
	}

	if req.SensorTypeID != nil {
		updatedReading.SensorTypeID = *req.SensorTypeID
	}

	if req.MacAddress != nil {
		updatedReading.MacAddress = *req.MacAddress
	}

	if req.Location != nil {
		updatedReading.Location = *req.Location
	}

	if req.MeasurementData != nil {
		updatedReading.MeasurementData = req.MeasurementData
		// Re-extract standard fields when measurement data changes
		standardFields, err := s.extractStandardFields(req.MeasurementData)
		if err == nil {
			updatedReading.StandardFields = standardFields
		}
	}

	if req.ReadingTime != nil {
		updatedReading.ReadingTime = *req.ReadingTime
	}

	// Update legacy fields if provided
	if req.DataX != nil {
		updatedReading.DataX = req.DataX
	}

	if req.DataY != nil {
		updatedReading.DataY = req.DataY
	}

	if req.PeakX != nil {
		updatedReading.PeakX = req.PeakX
	}

	if req.PeakY != nil {
		updatedReading.PeakY = req.PeakY
	}

	if req.PPM != nil {
		updatedReading.PPM = *req.PPM
	}

	if req.Label != nil {
		updatedReading.Label = *req.Label
	}

	if req.RawData != nil {
		updatedReading.RawData = req.RawData
	}

	// Update in repository
	if err := s.iotSensorReadingRepo.Update(ctx, &updatedReading); err != nil {
		return nil, err
	}

	return s.entityToResponse(&updatedReading), nil
}

// DeleteReading deletes an IoT sensor reading
func (s *IoTSensorReadingService) DeleteReading(ctx context.Context, id uuid.UUID) error {
	return s.iotSensorReadingRepo.Delete(ctx, id)
}

// GetLatestByMacAddress retrieves the latest IoT sensor reading for a MAC address
func (s *IoTSensorReadingService) GetLatestByMacAddress(ctx context.Context, macAddress string) (*dto.IoTSensorReadingResponse, error) {
	reading, err := s.iotSensorReadingRepo.GetLatestByMacAddress(ctx, macAddress)
	if err != nil {
		return nil, err
	}

	if reading == nil {
		return nil, errors.New("reading not found")
	}

	return s.entityToResponse(reading), nil
}

// Helper method to convert entity to response DTO
func (s *IoTSensorReadingService) entityToResponse(reading *entity.IoTSensorReading) *dto.IoTSensorReadingResponse {
	response := &dto.IoTSensorReadingResponse{
		ID:              reading.ID,
		SensorTypeID:    reading.SensorTypeID,
		MacAddress:      reading.MacAddress,
		Location:        reading.Location,
		MeasurementData: reading.MeasurementData,
		StandardFields:  reading.StandardFields,
		ReadingTime:     reading.ReadingTime.Format(time.RFC3339),
		CreatedAt:       reading.CreatedAt.Format(time.RFC3339),
		DataX:           reading.DataX,
		DataY:           reading.DataY,
		PeakX:           reading.PeakX,
		PeakY:           reading.PeakY,
		PPM:             &reading.PPM,
		Label:           &reading.Label,
		RawData:         reading.RawData,
	}

	// Convert UUID to pointer for response - handle zero UUID as nil
	if reading.TenantID != (uuid.UUID{}) {
		response.TenantID = &reading.TenantID
	}

	if reading.AssetSensorID != (uuid.UUID{}) {
		response.AssetSensorID = &reading.AssetSensorID
	}

	// Format updated_at if it exists
	if reading.UpdatedAt != nil {
		updatedAtStr := reading.UpdatedAt.Format(time.RFC3339)
		response.UpdatedAt = &updatedAtStr
	}

	return response
}

// Helper method to extract standard fields from measurement data
func (s *IoTSensorReadingService) extractStandardFields(measurementData json.RawMessage) (json.RawMessage, error) {
	if measurementData == nil {
		return nil, nil
	}

	var data map[string]interface{}
	if err := json.Unmarshal(measurementData, &data); err != nil {
		return nil, err
	}

	// Extract commonly used standard fields
	standardFields := make(map[string]interface{})

	// Common field mappings
	standardFieldKeys := []string{
		"value", "unit", "timestamp", "quality", "status",
		"temperature", "humidity", "pressure", "voltage",
		"current", "power", "frequency", "ph", "conductivity",
	}

	for _, key := range standardFieldKeys {
		if val, exists := data[key]; exists {
			standardFields[key] = val
		}
	}

	if len(standardFields) == 0 {
		return nil, nil
	}

	standardFieldsBytes, err := json.Marshal(standardFields)
	if err != nil {
		return nil, err
	}

	return json.RawMessage(standardFieldsBytes), nil
}
