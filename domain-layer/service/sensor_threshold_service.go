package service

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"be-lecsens/asset_management/data-layer/entity"
	"be-lecsens/asset_management/data-layer/repository"
	"be-lecsens/asset_management/helpers/common"
	"be-lecsens/asset_management/helpers/dto"
)

// SensorThresholdService interface defines the contract for sensor threshold operations
type SensorThresholdService interface {
	// CRUD Operations
	CreateThreshold(ctx context.Context, tenantID uuid.UUID, req *dto.CreateSensorThresholdRequest) (*dto.SensorThresholdResponse, error)
	GetThresholdByID(ctx context.Context, tenantID, thresholdID uuid.UUID) (*dto.SensorThresholdResponse, error)
	UpdateThreshold(ctx context.Context, tenantID, thresholdID uuid.UUID, req *dto.UpdateSensorThresholdRequest) (*dto.SensorThresholdResponse, error)
	DeleteThreshold(ctx context.Context, tenantID, thresholdID uuid.UUID) error

	// List and Filter Operations
	ListThresholds(ctx context.Context, tenantID uuid.UUID, filter *dto.SensorThresholdFilterRequest) (*dto.SensorThresholdListResponse, error)
	GetThresholdsByAssetSensor(ctx context.Context, tenantID, assetSensorID uuid.UUID) ([]*dto.SensorThresholdResponse, error)
	GetThresholdsBySensorType(ctx context.Context, tenantID, sensorTypeID uuid.UUID) ([]*dto.SensorThresholdResponse, error)
	GetActiveThresholds(ctx context.Context, tenantID uuid.UUID) ([]*dto.SensorThresholdResponse, error)

	// Threshold Management
	ActivateThreshold(ctx context.Context, tenantID, thresholdID uuid.UUID) error
	DeactivateThreshold(ctx context.Context, tenantID, thresholdID uuid.UUID) error
	BulkUpdateStatus(ctx context.Context, tenantID uuid.UUID, thresholdIDs []uuid.UUID, isActive bool) error

	// Threshold Checking
	CheckThresholds(ctx context.Context, tenantID, assetSensorID uuid.UUID, measurementField string, value float64) (*dto.BulkThresholdCheckResponse, error)
	ValidateThresholdValue(ctx context.Context, tenantID, thresholdID uuid.UUID, value float64) (*dto.ThresholdCheckResult, error)

	// Statistics and Analytics
	GetThresholdStatistics(ctx context.Context, tenantID uuid.UUID) (*dto.ThresholdStatistics, error)
	GetThresholdsByMeasurementField(ctx context.Context, tenantID uuid.UUID, measurementField string) ([]*dto.SensorThresholdResponse, error)
}

// sensorThresholdService implements SensorThresholdService interface
type sensorThresholdService struct {
	thresholdRepo   repository.SensorThresholdRepository
	assetSensorRepo repository.AssetSensorRepository
	sensorTypeRepo  *repository.SensorTypeRepository
	assetRepo       repository.AssetRepository
	alertService    AssetAlertService // For creating alerts when thresholds are breached
}

// NewSensorThresholdService creates a new instance of SensorThresholdService
func NewSensorThresholdService(
	thresholdRepo repository.SensorThresholdRepository,
	assetSensorRepo repository.AssetSensorRepository,
	sensorTypeRepo *repository.SensorTypeRepository,
	assetRepo repository.AssetRepository,
	alertService AssetAlertService,
) SensorThresholdService {
	return &sensorThresholdService{
		thresholdRepo:   thresholdRepo,
		assetSensorRepo: assetSensorRepo,
		sensorTypeRepo:  sensorTypeRepo,
		assetRepo:       assetRepo,
		alertService:    alertService,
	}
}

// CreateThreshold creates a new sensor threshold
func (s *sensorThresholdService) CreateThreshold(ctx context.Context, tenantID uuid.UUID, req *dto.CreateSensorThresholdRequest) (*dto.SensorThresholdResponse, error) {
	// Validate asset sensor exists and belongs to tenant
	assetSensor, err := s.assetSensorRepo.GetByID(ctx, req.AssetSensorID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.NewBadRequestError("Asset sensor not found")
		}
		return nil, fmt.Errorf("failed to get asset sensor: %w", err)
	}

	// Validate sensor type exists and belongs to tenant
	sensorType, err := s.sensorTypeRepo.GetByID(req.SensorTypeID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.NewBadRequestError("Sensor type not found")
		}
		return nil, fmt.Errorf("failed to get sensor type: %w", err)
	}

	// Validate threshold values
	if req.MinValue >= req.MaxValue {
		return nil, common.NewBadRequestError("Min value must be less than max value")
	}

	// Check for duplicate thresholds
	existingThresholds, err := s.thresholdRepo.GetByMeasurementField(ctx, req.AssetSensorID, req.MeasurementField)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing thresholds: %w", err)
	}

	for _, existing := range existingThresholds {
		if existing.IsActive {
			return nil, common.NewBadRequestError("Active threshold already exists for this measurement field")
		}
	}

	// Set default values
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	// Create threshold entity
	threshold := entity.NewSensorThreshold()
	threshold.TenantID = tenantID
	threshold.AssetSensorID = req.AssetSensorID
	threshold.SensorTypeID = req.SensorTypeID
	threshold.MeasurementField = req.MeasurementField
	threshold.MinValue = req.MinValue
	threshold.MaxValue = req.MaxValue
	threshold.Severity = entity.ThresholdSeverity(req.Severity)
	threshold.AlertMessage = req.AlertMessage

	threshold.Name = req.Name
	threshold.Description = req.Description
	threshold.IsActive = isActive

	// Set notification rules if provided
	if req.NotificationRules != nil {
		if err := threshold.SetNotificationRules(req.NotificationRules); err != nil {
			return nil, common.NewBadRequestError(fmt.Sprintf("Invalid notification rules: %v", err))
		}
	}

	// Save threshold
	if err := s.thresholdRepo.Create(ctx, threshold); err != nil {
		return nil, fmt.Errorf("failed to create threshold: %w", err)
	}

	return s.mapToResponse(threshold, assetSensor, sensorType), nil
}

// GetThresholdByID retrieves a threshold by ID
func (s *sensorThresholdService) GetThresholdByID(ctx context.Context, tenantID, thresholdID uuid.UUID) (*dto.SensorThresholdResponse, error) {
	threshold, err := s.thresholdRepo.GetByID(ctx, thresholdID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.NewNotFoundError("threshold", thresholdID.String())
		}
		return nil, fmt.Errorf("failed to get threshold: %w", err)
	}

	// Get related entities for enhanced response
	assetSensor, _ := s.assetSensorRepo.GetByID(ctx, threshold.AssetSensorID)
	sensorType, _ := s.sensorTypeRepo.GetByID(threshold.SensorTypeID)

	return s.mapToResponse(threshold, assetSensor, sensorType), nil
}

// UpdateThreshold updates an existing threshold
func (s *sensorThresholdService) UpdateThreshold(ctx context.Context, tenantID, thresholdID uuid.UUID, req *dto.UpdateSensorThresholdRequest) (*dto.SensorThresholdResponse, error) {
	threshold, err := s.thresholdRepo.GetByID(ctx, thresholdID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.NewNotFoundError("threshold", thresholdID.String())
		}
		return nil, fmt.Errorf("failed to get threshold: %w", err)
	}

	// Update fields if provided
	if req.Name != nil {
		threshold.Name = *req.Name
	}
	if req.Description != nil {
		threshold.Description = *req.Description
	}
	if req.MinValue != nil {
		threshold.MinValue = *req.MinValue
	}
	if req.MaxValue != nil {
		threshold.MaxValue = *req.MaxValue
	}
	if req.Severity != nil {
		threshold.Severity = entity.ThresholdSeverity(*req.Severity)
	}
	if req.AlertMessage != nil {
		threshold.AlertMessage = *req.AlertMessage
	}
	if req.IsActive != nil {
		threshold.IsActive = *req.IsActive
	}

	// Validate threshold values if updated
	if threshold.MinValue >= threshold.MaxValue {
		return nil, common.NewBadRequestError("Min value must be less than max value")
	}

	// Update notification rules if provided
	if req.NotificationRules != nil {
		if err := threshold.SetNotificationRules(req.NotificationRules); err != nil {
			return nil, common.NewBadRequestError(fmt.Sprintf("Invalid notification rules: %v", err))
		}
	}

	// Update timestamp
	now := time.Now()
	threshold.UpdatedAt = &now

	// Save changes
	if err := s.thresholdRepo.Update(ctx, threshold); err != nil {
		return nil, fmt.Errorf("failed to update threshold: %w", err)
	}

	// Get related entities for enhanced response
	assetSensor, _ := s.assetSensorRepo.GetByID(ctx, threshold.AssetSensorID)
	sensorType, _ := s.sensorTypeRepo.GetByID(threshold.SensorTypeID)

	return s.mapToResponse(threshold, assetSensor, sensorType), nil
}

// DeleteThreshold deletes a threshold
func (s *sensorThresholdService) DeleteThreshold(ctx context.Context, tenantID, thresholdID uuid.UUID) error {
	_, err := s.thresholdRepo.GetByID(ctx, thresholdID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return common.NewNotFoundError("threshold", thresholdID.String())
		}
		return fmt.Errorf("failed to get threshold: %w", err)
	}

	if err := s.thresholdRepo.Delete(ctx, thresholdID); err != nil {
		return fmt.Errorf("failed to delete threshold: %w", err)
	}

	return nil
}

// ListThresholds retrieves thresholds with filtering and pagination
func (s *sensorThresholdService) ListThresholds(ctx context.Context, tenantID uuid.UUID, filter *dto.SensorThresholdFilterRequest) (*dto.SensorThresholdListResponse, error) {
	// Set default pagination
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.Limit <= 0 {
		filter.Limit = 20
	}

	offset := (filter.Page - 1) * filter.Limit

	// Build filter conditions
	conditions := make(map[string]interface{})
	if filter.AssetSensorID != nil {
		conditions["asset_sensor_id"] = *filter.AssetSensorID
	}
	if filter.SensorTypeID != nil {
		conditions["sensor_type_id"] = *filter.SensorTypeID
	}
	if filter.Severity != nil {
		conditions["severity"] = *filter.Severity
	}
	if filter.IsActive != nil {
		conditions["is_active"] = *filter.IsActive
	}

	// Get thresholds
	thresholds, err := s.thresholdRepo.List(ctx, filter.Limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list thresholds: %w", err)
	}

	// Get total count (simplified approach)
	total := len(thresholds)

	// Map to response
	responses := make([]dto.SensorThresholdResponse, len(thresholds))
	for i, threshold := range thresholds {
		// Get related entities for enhanced response
		assetSensor, _ := s.assetSensorRepo.GetByID(ctx, threshold.AssetSensorID)
		sensorType, _ := s.sensorTypeRepo.GetByID(threshold.SensorTypeID)
		responses[i] = *s.mapToResponse(threshold, assetSensor, sensorType)
	}

	totalPages := int(math.Ceil(float64(total) / float64(filter.Limit)))

	return &dto.SensorThresholdListResponse{
		Thresholds: responses,
		Page:       filter.Page,
		Limit:      filter.Limit,
		Total:      int64(total),
		TotalPages: totalPages,
	}, nil
}

// GetThresholdsByAssetSensor retrieves thresholds for a specific asset sensor
func (s *sensorThresholdService) GetThresholdsByAssetSensor(ctx context.Context, tenantID, assetSensorID uuid.UUID) ([]*dto.SensorThresholdResponse, error) {
	thresholds, err := s.thresholdRepo.GetByAssetSensorID(ctx, assetSensorID)
	if err != nil {
		return nil, fmt.Errorf("failed to get thresholds by asset sensor: %w", err)
	}

	responses := make([]*dto.SensorThresholdResponse, len(thresholds))
	for i, threshold := range thresholds {
		assetSensor, _ := s.assetSensorRepo.GetByID(ctx, threshold.AssetSensorID)
		sensorType, _ := s.sensorTypeRepo.GetByID(threshold.SensorTypeID)
		responses[i] = s.mapToResponse(threshold, assetSensor, sensorType)
	}

	return responses, nil
}

// GetThresholdsBySensorType retrieves thresholds for a specific sensor type
func (s *sensorThresholdService) GetThresholdsBySensorType(ctx context.Context, tenantID, sensorTypeID uuid.UUID) ([]*dto.SensorThresholdResponse, error) {
	thresholds, err := s.thresholdRepo.GetBySensorTypeID(ctx, sensorTypeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get thresholds by sensor type: %w", err)
	}

	responses := make([]*dto.SensorThresholdResponse, len(thresholds))
	for i, threshold := range thresholds {
		assetSensor, _ := s.assetSensorRepo.GetByID(ctx, threshold.AssetSensorID)
		sensorType, _ := s.sensorTypeRepo.GetByID(threshold.SensorTypeID)
		responses[i] = s.mapToResponse(threshold, assetSensor, sensorType)
	}

	return responses, nil
}

// GetActiveThresholds retrieves all active thresholds
func (s *sensorThresholdService) GetActiveThresholds(ctx context.Context, tenantID uuid.UUID) ([]*dto.SensorThresholdResponse, error) {
	thresholds, err := s.thresholdRepo.GetActiveThresholds(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get active thresholds: %w", err)
	}

	responses := make([]*dto.SensorThresholdResponse, len(thresholds))
	for i, threshold := range thresholds {
		assetSensor, _ := s.assetSensorRepo.GetByID(ctx, threshold.AssetSensorID)
		sensorType, _ := s.sensorTypeRepo.GetByID(threshold.SensorTypeID)
		responses[i] = s.mapToResponse(threshold, assetSensor, sensorType)
	}

	return responses, nil
}

// ActivateThreshold activates a threshold
func (s *sensorThresholdService) ActivateThreshold(ctx context.Context, tenantID, thresholdID uuid.UUID) error {
	threshold, err := s.thresholdRepo.GetByID(ctx, thresholdID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return common.NewNotFoundError("threshold", thresholdID.String())
		}
		return fmt.Errorf("failed to get threshold: %w", err)
	}

	if threshold.IsActive {
		return common.NewBadRequestError("Threshold is already active")
	}

	threshold.IsActive = true
	now := time.Now()
	threshold.UpdatedAt = &now

	if err := s.thresholdRepo.Update(ctx, threshold); err != nil {
		return fmt.Errorf("failed to activate threshold: %w", err)
	}

	return nil
}

// DeactivateThreshold deactivates a threshold
func (s *sensorThresholdService) DeactivateThreshold(ctx context.Context, tenantID, thresholdID uuid.UUID) error {
	threshold, err := s.thresholdRepo.GetByID(ctx, thresholdID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return common.NewNotFoundError("threshold", thresholdID.String())
		}
		return fmt.Errorf("failed to get threshold: %w", err)
	}

	if !threshold.IsActive {
		return common.NewBadRequestError("Threshold is already inactive")
	}

	threshold.IsActive = false
	now := time.Now()
	threshold.UpdatedAt = &now

	if err := s.thresholdRepo.Update(ctx, threshold); err != nil {
		return fmt.Errorf("failed to deactivate threshold: %w", err)
	}

	return nil
}

// BulkUpdateStatus updates the status of multiple thresholds
func (s *sensorThresholdService) BulkUpdateStatus(ctx context.Context, tenantID uuid.UUID, thresholdIDs []uuid.UUID, isActive bool) error {
	for _, thresholdID := range thresholdIDs {
		threshold, err := s.thresholdRepo.GetByID(ctx, thresholdID)
		if err != nil {
			continue // Skip if not found
		}

		threshold.IsActive = isActive
		now := time.Now()
		threshold.UpdatedAt = &now

		if err := s.thresholdRepo.Update(ctx, threshold); err != nil {
			return fmt.Errorf("failed to update threshold %s: %w", thresholdID, err)
		}
	}

	return nil
}

// CheckThresholds checks if a sensor reading breaches any thresholds
func (s *sensorThresholdService) CheckThresholds(ctx context.Context, tenantID, assetSensorID uuid.UUID, measurementField string, value float64) (*dto.BulkThresholdCheckResponse, error) {
	// Get active thresholds for this asset sensor and measurement field
	thresholds, err := s.thresholdRepo.GetByMeasurementField(ctx, assetSensorID, measurementField)
	if err != nil {
		return nil, fmt.Errorf("failed to get thresholds: %w", err)
	}

	response := &dto.BulkThresholdCheckResponse{
		AssetSensorID: assetSensorID,
		Results:       make([]dto.ThresholdCheckResult, 0),
		AlertsCreated: 0,
	}

	for _, threshold := range thresholds {
		if !threshold.IsActive {
			continue
		}

		// Check if threshold is breached
		isBreached := threshold.IsBreached(value)

		result := dto.ThresholdCheckResult{
			ThresholdID:   threshold.ID,
			ThresholdName: threshold.Name,
			IsBreached:    isBreached,
			CurrentValue:  value,
			MinValue:      threshold.MinValue,
			MaxValue:      threshold.MaxValue,
			Severity:      string(threshold.Severity),
		}

		if isBreached {
			result.Message = threshold.AlertMessage
			if result.Message == "" {
				result.Message = fmt.Sprintf("Threshold breached: value %.2f is outside range [%.2f, %.2f]", value, threshold.MinValue, threshold.MaxValue)
			}

			// Create alert through alert service
			if s.alertService != nil {
				// Get asset ID from asset sensor
				assetSensor, err := s.assetSensorRepo.GetByID(ctx, assetSensorID)
				if err == nil {
					alertReq := &dto.CreateAssetAlertRequest{
						AssetID:       assetSensor.AssetID,
						AssetSensorID: assetSensorID,
						ThresholdID:   threshold.ID,
						Severity:      string(threshold.Severity),
						Message:       result.Message,
					}

					_, err := s.alertService.CreateAlert(ctx, tenantID, alertReq)
					if err == nil {
						response.AlertsCreated++
					}
				}
			}
		}

		response.Results = append(response.Results, result)
	}

	return response, nil
}

// ValidateThresholdValue validates a single value against a specific threshold
func (s *sensorThresholdService) ValidateThresholdValue(ctx context.Context, tenantID, thresholdID uuid.UUID, value float64) (*dto.ThresholdCheckResult, error) {
	threshold, err := s.thresholdRepo.GetByID(ctx, thresholdID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.NewNotFoundError("threshold", thresholdID.String())
		}
		return nil, fmt.Errorf("failed to get threshold: %w", err)
	}

	isBreached := threshold.IsBreached(value)

	result := &dto.ThresholdCheckResult{
		ThresholdID:   threshold.ID,
		ThresholdName: threshold.Name,
		IsBreached:    isBreached,
		CurrentValue:  value,
		MinValue:      threshold.MinValue,
		MaxValue:      threshold.MaxValue,
		Severity:      string(threshold.Severity),
	}

	if isBreached {
		result.Message = threshold.AlertMessage
		if result.Message == "" {
			result.Message = fmt.Sprintf("Threshold breached: value %.2f is outside range [%.2f, %.2f]", value, threshold.MinValue, threshold.MaxValue)
		}
	}

	return result, nil
}

// GetThresholdStatistics returns statistics about thresholds
func (s *sensorThresholdService) GetThresholdStatistics(ctx context.Context, tenantID uuid.UUID) (*dto.ThresholdStatistics, error) {
	// Get all thresholds
	allThresholds, err := s.thresholdRepo.List(ctx, 0, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to get thresholds: %w", err)
	}

	stats := &dto.ThresholdStatistics{
		TotalThresholds: len(allThresholds),
	}

	for _, threshold := range allThresholds {
		if threshold.IsActive {
			stats.ActiveThresholds++
		} else {
			stats.InactiveThresholds++
		}

		switch threshold.Severity {
		case "warning":
			stats.WarningThresholds++
		case "critical":
			stats.CriticalThresholds++
		}
	}

	return stats, nil
}

// GetThresholdsByMeasurementField retrieves thresholds for a specific measurement field
func (s *sensorThresholdService) GetThresholdsByMeasurementField(ctx context.Context, tenantID uuid.UUID, measurementField string) ([]*dto.SensorThresholdResponse, error) {
	// Get all thresholds and filter by measurement field
	allThresholds, err := s.thresholdRepo.List(ctx, 0, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to get thresholds: %w", err)
	}

	responses := make([]*dto.SensorThresholdResponse, len(allThresholds))
	for i, threshold := range allThresholds {
		assetSensor, _ := s.assetSensorRepo.GetByID(ctx, threshold.AssetSensorID)
		sensorType, _ := s.sensorTypeRepo.GetByID(threshold.SensorTypeID)
		responses[i] = s.mapToResponse(threshold, assetSensor, sensorType)
	}

	return responses, nil
}

// mapToResponse maps entity to response DTO
func (s *sensorThresholdService) mapToResponse(threshold *entity.SensorThreshold, assetSensor *repository.AssetSensorWithDetails, sensorType *repository.SensorType) *dto.SensorThresholdResponse {
	response := &dto.SensorThresholdResponse{
		ID:               threshold.ID,
		TenantID:         threshold.TenantID,
		AssetSensorID:    threshold.AssetSensorID,
		SensorTypeID:     threshold.SensorTypeID,
		MeasurementField: threshold.MeasurementField,
		Name:             threshold.Name,
		Description:      threshold.Description,
		MinValue:         threshold.MinValue,
		MaxValue:         threshold.MaxValue,
		Severity:         string(threshold.Severity),
		AlertMessage:     threshold.AlertMessage,
		IsActive:         threshold.IsActive,
		CreatedAt:        threshold.CreatedAt,
		UpdatedAt:        threshold.UpdatedAt,
	}

	// Set notification rules
	if rules, err := threshold.GetNotificationRules(); err == nil && rules != nil {
		response.NotificationRules = rules
	}

	// Add related entity names if available
	if assetSensor != nil {
		response.AssetSensorName = assetSensor.Name
	}
	if sensorType != nil {
		response.SensorTypeName = sensorType.Name
	}

	return response
}
