package service

import (
	"be-lecsens/asset_management/data-layer/entity"
	"be-lecsens/asset_management/data-layer/repository"
	"be-lecsens/asset_management/helpers/common"
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
)

// SensorThresholdService handles business logic for sensor thresholds
type SensorThresholdService struct {
	sensorThresholdRepo repository.SensorThresholdRepository
	assetSensorRepo     repository.AssetSensorRepository
	assetAlertRepo      repository.AssetAlertRepository
}

// NewSensorThresholdService creates a new instance of SensorThresholdService
func NewSensorThresholdService(
	sensorThresholdRepo repository.SensorThresholdRepository,
	assetSensorRepo repository.AssetSensorRepository,
	assetAlertRepo repository.AssetAlertRepository,
) *SensorThresholdService {
	return &SensorThresholdService{
		sensorThresholdRepo: sensorThresholdRepo,
		assetSensorRepo:     assetSensorRepo,
		assetAlertRepo:      assetAlertRepo,
	}
}

// CreateSensorThreshold creates a new sensor threshold
func (s *SensorThresholdService) CreateSensorThreshold(ctx context.Context, threshold *entity.SensorThreshold) (*entity.SensorThreshold, error) {
	log.Printf("Creating sensor threshold: %+v", threshold)

	// Validate asset sensor exists
	if s.assetSensorRepo != nil {
		assetSensor, err := s.assetSensorRepo.GetByID(ctx, threshold.AssetSensorID)
		if err != nil {
			return nil, fmt.Errorf("failed to validate asset sensor: %w", err)
		}
		if assetSensor == nil {
			return nil, common.NewValidationError("asset sensor not found", nil)
		}
	}

	// Validate threshold values
	if threshold.MinValue != nil && threshold.MaxValue != nil && *threshold.MinValue >= *threshold.MaxValue {
		return nil, common.NewValidationError("minimum threshold must be less than maximum threshold", nil)
	}

	// Ensure at least one threshold value is set
	if threshold.MinValue == nil && threshold.MaxValue == nil {
		return nil, common.NewValidationError("at least one threshold value (min or max) must be set", nil)
	}

	// Create the threshold
	if err := s.sensorThresholdRepo.Create(ctx, threshold); err != nil {
		log.Printf("Error creating sensor threshold: %v", err)
		return nil, fmt.Errorf("failed to create sensor threshold: %w", err)
	}

	log.Printf("Successfully created sensor threshold with ID: %s", threshold.ID)
	return threshold, nil
}

// GetSensorThresholdByID retrieves a sensor threshold by its ID
func (s *SensorThresholdService) GetSensorThresholdByID(ctx context.Context, id uuid.UUID) (*entity.SensorThreshold, error) {
	threshold, err := s.sensorThresholdRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get sensor threshold: %w", err)
	}

	if threshold == nil {
		return nil, common.NewNotFoundError("sensor threshold", id.String())
	}

	return threshold, nil
}

// ListSensorThresholds retrieves paginated sensor thresholds for a tenant
func (s *SensorThresholdService) ListSensorThresholds(ctx context.Context, tenantID uuid.UUID, limit, offset int) ([]*entity.SensorThreshold, int, error) {
	if limit <= 0 {
		limit = 20 // Default limit
	}
	if limit > 100 {
		limit = 100 // Maximum limit
	}

	thresholds, totalCount, err := s.sensorThresholdRepo.List(ctx, tenantID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list sensor thresholds: %w", err)
	}

	return thresholds, totalCount, nil
}

// UpdateSensorThreshold updates an existing sensor threshold
func (s *SensorThresholdService) UpdateSensorThreshold(ctx context.Context, threshold *entity.SensorThreshold) (*entity.SensorThreshold, error) {
	log.Printf("Updating sensor threshold: %+v", threshold)

	// Validate threshold values
	if threshold.MinValue != nil && threshold.MaxValue != nil && *threshold.MinValue >= *threshold.MaxValue {
		return nil, common.NewValidationError("minimum threshold must be less than maximum threshold", nil)
	}

	// Ensure at least one threshold value is set
	if threshold.MinValue == nil && threshold.MaxValue == nil {
		return nil, common.NewValidationError("at least one threshold value (min or max) must be set", nil)
	}

	// Check if threshold exists
	existing, err := s.sensorThresholdRepo.GetByID(ctx, threshold.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing threshold: %w", err)
	}
	if existing == nil {
		return nil, common.NewNotFoundError("sensor threshold", threshold.ID.String())
	}

	// Update the threshold
	if err := s.sensorThresholdRepo.Update(ctx, threshold); err != nil {
		log.Printf("Error updating sensor threshold: %v", err)
		return nil, fmt.Errorf("failed to update sensor threshold: %w", err)
	}

	// Get updated threshold
	updatedThreshold, err := s.sensorThresholdRepo.GetByID(ctx, threshold.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated threshold: %w", err)
	}

	log.Printf("Successfully updated sensor threshold with ID: %s", threshold.ID)
	return updatedThreshold, nil
}

// DeleteSensorThreshold deletes a sensor threshold
func (s *SensorThresholdService) DeleteSensorThreshold(ctx context.Context, id uuid.UUID) error {
	log.Printf("Deleting sensor threshold with ID: %s", id)

	// Check if threshold exists
	existing, err := s.sensorThresholdRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to check existing threshold: %w", err)
	}
	if existing == nil {
		return common.NewNotFoundError("sensor threshold", id.String())
	}

	// Delete the threshold
	if err := s.sensorThresholdRepo.Delete(ctx, id); err != nil {
		log.Printf("Error deleting sensor threshold: %v", err)
		return fmt.Errorf("failed to delete sensor threshold: %w", err)
	}

	log.Printf("Successfully deleted sensor threshold with ID: %s", id)
	return nil
}

// GetThresholdsByAssetSensor retrieves thresholds for a specific asset sensor
func (s *SensorThresholdService) GetThresholdsByAssetSensor(ctx context.Context, assetSensorID uuid.UUID) ([]*entity.SensorThreshold, error) {
	// Validate asset sensor exists
	if s.assetSensorRepo != nil {
		assetSensor, err := s.assetSensorRepo.GetByID(ctx, assetSensorID)
		if err != nil {
			return nil, fmt.Errorf("failed to validate asset sensor: %w", err)
		}
		if assetSensor == nil {
			return nil, common.NewNotFoundError("asset sensor", assetSensorID.String())
		}
	}

	thresholds, err := s.sensorThresholdRepo.GetByAssetSensorID(ctx, assetSensorID)
	if err != nil {
		return nil, fmt.Errorf("failed to get thresholds by asset sensor: %w", err)
	}

	return thresholds, nil
}

// GetThresholdsByMeasurementType retrieves thresholds for a specific measurement type
func (s *SensorThresholdService) GetThresholdsByMeasurementType(ctx context.Context, measurementTypeID uuid.UUID) ([]*entity.SensorThreshold, error) {
	thresholds, err := s.sensorThresholdRepo.GetByMeasurementTypeID(ctx, measurementTypeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get thresholds by measurement type: %w", err)
	}

	return thresholds, nil
}

// ActivateSensorThreshold activates a sensor threshold
func (s *SensorThresholdService) ActivateSensorThreshold(ctx context.Context, id uuid.UUID) (*entity.SensorThreshold, error) {
	// Check if threshold exists
	threshold, err := s.sensorThresholdRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get threshold: %w", err)
	}
	if threshold == nil {
		return nil, common.NewNotFoundError("sensor threshold", id.String())
	}

	// Activate the threshold
	threshold.IsActive = true
	if err := s.sensorThresholdRepo.Update(ctx, threshold); err != nil {
		return nil, fmt.Errorf("failed to activate threshold: %w", err)
	}

	return threshold, nil
}

// DeactivateSensorThreshold deactivates a sensor threshold
func (s *SensorThresholdService) DeactivateSensorThreshold(ctx context.Context, id uuid.UUID) (*entity.SensorThreshold, error) {
	// Check if threshold exists
	threshold, err := s.sensorThresholdRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get threshold: %w", err)
	}
	if threshold == nil {
		return nil, common.NewNotFoundError("sensor threshold", id.String())
	}

	// Deactivate the threshold
	threshold.IsActive = false
	if err := s.sensorThresholdRepo.Update(ctx, threshold); err != nil {
		return nil, fmt.Errorf("failed to deactivate threshold: %w", err)
	}

	return threshold, nil
}

// ListAllSensorThresholds retrieves all sensor thresholds across all tenants
func (s *SensorThresholdService) ListAllSensorThresholds(ctx context.Context, limit, offset int) ([]*entity.SensorThreshold, int, error) {
	if limit <= 0 {
		limit = 20 // Default limit
	}
	if limit > 100 {
		limit = 100 // Maximum limit
	}

	thresholds, totalCount, err := s.sensorThresholdRepo.ListAll(ctx, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list all sensor thresholds: %w", err)
	}

	return thresholds, totalCount, nil
}

// CheckThresholdsForValue checks if a value breaches any thresholds and creates alerts
func (s *SensorThresholdService) CheckThresholdsForValue(
	ctx context.Context,
	reading *entity.IoTSensorReading,
	fieldName string,
	value float64,
) error {
	// Get all thresholds for this measurement type
	thresholds, err := s.sensorThresholdRepo.GetByMeasurementTypeID(ctx, reading.SensorTypeID)
	if err != nil {
		return fmt.Errorf("failed to get thresholds: %w", err)
	}

	// Check each threshold
	for _, threshold := range thresholds {
		// Skip if threshold is not active
		if !threshold.IsActive {
			continue
		}

		// Skip if field name doesn't match
		if threshold.MeasurementFieldName != fieldName {
			continue
		}

		// Check if value breaches threshold
		status := threshold.CheckValue(value)

		if status != entity.ThresholdStatusNormal {
			// Create or update alert
			if err := s.assetAlertRepo.CreateAlertFromReading(ctx, reading, threshold, value); err != nil {
				log.Printf("Error creating alert for threshold breach: %v", err)
				continue
			}
		} else {
			// Resolve any active alerts for this threshold
			if err := s.assetAlertRepo.ResolveAlertsForReading(ctx, reading, threshold, value); err != nil {
				log.Printf("Error resolving alerts for normal value: %v", err)
				continue
			}
		}
	}

	return nil
}
