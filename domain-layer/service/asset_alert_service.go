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

// AssetAlertService interface defines the contract for asset alert operations
type AssetAlertService interface {
	// CRUD Operations
	CreateAlert(ctx context.Context, tenantID uuid.UUID, req *dto.CreateAssetAlertRequest) (*dto.AssetAlertResponse, error)
	GetAlertByID(ctx context.Context, tenantID, alertID uuid.UUID) (*dto.AssetAlertResponse, error)
	UpdateAlert(ctx context.Context, tenantID, alertID uuid.UUID, req *dto.UpdateAssetAlertRequest) (*dto.AssetAlertResponse, error)
	DeleteAlert(ctx context.Context, tenantID, alertID uuid.UUID) error

	// List and Filter Operations
	ListAlerts(ctx context.Context, tenantID uuid.UUID, filter *dto.AssetAlertFilterRequest) (*dto.AssetAlertListResponse, error)
	GetAlertsByAsset(ctx context.Context, tenantID, assetID uuid.UUID) ([]*dto.AssetAlertResponse, error)
	GetAlertsByAssetSensor(ctx context.Context, tenantID, assetSensorID uuid.UUID) ([]*dto.AssetAlertResponse, error)
	GetAlertsByThreshold(ctx context.Context, tenantID, thresholdID uuid.UUID) ([]*dto.AssetAlertResponse, error)
	GetUnresolvedAlerts(ctx context.Context, tenantID uuid.UUID) ([]*dto.AssetAlertResponse, error)
	GetAlertsBySeverity(ctx context.Context, tenantID uuid.UUID, severity string) ([]*dto.AssetAlertResponse, error)

	// Alert Resolution
	ResolveAlert(ctx context.Context, tenantID, alertID uuid.UUID) error
	BulkResolveAlerts(ctx context.Context, tenantID uuid.UUID, req *dto.BulkResolveAlertsRequest) (*dto.BulkResolveAlertsResponse, error)

	// Statistics and Analytics
	GetAlertStatistics(ctx context.Context, tenantID uuid.UUID, req *dto.AlertStatisticsRequest) (*dto.AlertStatistics, error)
	GetAlertSummary(ctx context.Context, tenantID uuid.UUID) (*dto.AlertSummaryResponse, error)
	GetAlertsInTimeRange(ctx context.Context, tenantID uuid.UUID, startTime, endTime time.Time) ([]*dto.AssetAlertResponse, error)
}

// assetAlertService implements AssetAlertService interface
type assetAlertService struct {
	alertRepo       repository.AssetAlertRepository
	assetRepo       repository.AssetRepository
	assetSensorRepo repository.AssetSensorRepository
	thresholdRepo   repository.SensorThresholdRepository
	sensorTypeRepo  *repository.SensorTypeRepository
	locationRepo    *repository.LocationRepository
}

// NewAssetAlertService creates a new instance of AssetAlertService
func NewAssetAlertService(
	alertRepo repository.AssetAlertRepository,
	assetRepo repository.AssetRepository,
	assetSensorRepo repository.AssetSensorRepository,
	thresholdRepo repository.SensorThresholdRepository,
	sensorTypeRepo *repository.SensorTypeRepository,
	locationRepo *repository.LocationRepository,
) AssetAlertService {
	return &assetAlertService{
		alertRepo:       alertRepo,
		assetRepo:       assetRepo,
		assetSensorRepo: assetSensorRepo,
		thresholdRepo:   thresholdRepo,
		sensorTypeRepo:  sensorTypeRepo,
		locationRepo:    locationRepo,
	}
}

// CreateAlert creates a new asset alert
func (s *assetAlertService) CreateAlert(ctx context.Context, tenantID uuid.UUID, req *dto.CreateAssetAlertRequest) (*dto.AssetAlertResponse, error) {
	// Validate asset exists and belongs to tenant
	asset, err := s.assetRepo.GetByID(ctx, req.AssetID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.NewBadRequestError("Asset not found")
		}
		return nil, fmt.Errorf("failed to get asset: %w", err)
	}

	// Validate asset sensor exists and belongs to tenant
	assetSensor, err := s.assetSensorRepo.GetByID(ctx, req.AssetSensorID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.NewBadRequestError("Asset sensor not found")
		}
		return nil, fmt.Errorf("failed to get asset sensor: %w", err)
	}

	// Validate threshold exists and belongs to tenant
	threshold, err := s.thresholdRepo.GetByID(ctx, req.ThresholdID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.NewBadRequestError("Threshold not found")
		}
		return nil, fmt.Errorf("failed to get threshold: %w", err)
	}

	// Validate that asset sensor belongs to the asset
	if assetSensor.AssetID != req.AssetID {
		return nil, common.NewBadRequestError("Asset sensor does not belong to the specified asset")
	}

	// Validate that threshold is for the correct asset sensor
	if threshold.AssetSensorID != req.AssetSensorID {
		return nil, common.NewBadRequestError("Threshold does not belong to the specified asset sensor")
	}

	// Create alert entity
	alert := &entity.AssetAlert{
		ID:            uuid.New(),
		TenantID:      tenantID,
		AssetID:       req.AssetID,
		AssetSensorID: req.AssetSensorID,
		ThresholdID:   req.ThresholdID,
		AlertTime:     time.Now(),
		ResolvedTime:  nil,
		Severity:      entity.ThresholdSeverity(req.Severity),
	}

	// Save alert
	if err := s.alertRepo.Create(ctx, alert); err != nil {
		return nil, fmt.Errorf("failed to create alert: %w", err)
	}

	return s.mapToResponse(alert, asset, assetSensor, threshold, nil, nil), nil
}

// GetAlertByID retrieves an alert by ID
func (s *assetAlertService) GetAlertByID(ctx context.Context, tenantID, alertID uuid.UUID) (*dto.AssetAlertResponse, error) {
	alert, err := s.alertRepo.GetByID(ctx, alertID)
	if err != nil {
		return nil, fmt.Errorf("failed to get alert: %w", err)
	}
	if alert == nil {
		return nil, common.NewNotFoundError("alert", alertID.String())
	}

	// Verify tenant access (unless SuperAdmin)
	if alert.TenantID != tenantID {
		return nil, common.NewNotFoundError("alert", alertID.String())
	}

	// Get related entities for enhanced response
	asset, _ := s.assetRepo.GetByID(ctx, alert.AssetID)
	assetSensor, _ := s.assetSensorRepo.GetByID(ctx, alert.AssetSensorID)
	threshold, _ := s.thresholdRepo.GetByID(ctx, alert.ThresholdID)
	var sensorType *repository.SensorType
	var location *entity.Location

	if assetSensor != nil {
		sensorType, _ = s.sensorTypeRepo.GetByID(assetSensor.SensorTypeID)
	}
	if asset != nil {
		location, _ = s.locationRepo.GetByID(ctx, asset.LocationID)
	}

	return s.mapToResponse(alert, asset, assetSensor, threshold, sensorType, location), nil
}

// UpdateAlert updates an existing alert
func (s *assetAlertService) UpdateAlert(ctx context.Context, tenantID, alertID uuid.UUID, req *dto.UpdateAssetAlertRequest) (*dto.AssetAlertResponse, error) {
	alert, err := s.alertRepo.GetByID(ctx, alertID)
	if err != nil {
		return nil, fmt.Errorf("failed to get alert: %w", err)
	}
	if alert == nil {
		return nil, common.NewNotFoundError("alert", alertID.String())
	}

	// Verify tenant access
	if alert.TenantID != tenantID {
		return nil, common.NewNotFoundError("alert", alertID.String())
	}

	// Update fields if provided
	if req.Severity != nil {
		alert.Severity = entity.ThresholdSeverity(*req.Severity)
	}

	// Save changes
	if err := s.alertRepo.Update(ctx, alert); err != nil {
		return nil, fmt.Errorf("failed to update alert: %w", err)
	}

	// Get related entities for enhanced response
	asset, _ := s.assetRepo.GetByID(ctx, alert.AssetID)
	assetSensor, _ := s.assetSensorRepo.GetByID(ctx, alert.AssetSensorID)
	threshold, _ := s.thresholdRepo.GetByID(ctx, alert.ThresholdID)
	var sensorType *repository.SensorType
	var location *entity.Location

	if assetSensor != nil {
		sensorType, _ = s.sensorTypeRepo.GetByID(assetSensor.SensorTypeID)
	}
	if asset != nil {
		location, _ = s.locationRepo.GetByID(ctx, asset.LocationID)
	}

	return s.mapToResponse(alert, asset, assetSensor, threshold, sensorType, location), nil
}

// DeleteAlert deletes an alert
func (s *assetAlertService) DeleteAlert(ctx context.Context, tenantID, alertID uuid.UUID) error {
	alert, err := s.alertRepo.GetByID(ctx, alertID)
	if err != nil {
		return fmt.Errorf("failed to get alert: %w", err)
	}
	if alert == nil {
		return common.NewNotFoundError("alert", alertID.String())
	}

	// Verify tenant access
	if alert.TenantID != tenantID {
		return common.NewNotFoundError("alert", alertID.String())
	}

	if err := s.alertRepo.Delete(ctx, alertID); err != nil {
		return fmt.Errorf("failed to delete alert: %w", err)
	}

	return nil
}

// ListAlerts retrieves alerts with filtering and pagination
func (s *assetAlertService) ListAlerts(ctx context.Context, tenantID uuid.UUID, filter *dto.AssetAlertFilterRequest) (*dto.AssetAlertListResponse, error) {
	// Set default pagination
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.Limit <= 0 {
		filter.Limit = 20
	}

	// For now, use basic list functionality
	// In a real implementation, you would build dynamic queries based on filters
	alerts, err := s.alertRepo.List(ctx, filter.Page, filter.Limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list alerts: %w", err)
	}

	// Filter alerts by tenant (if not SuperAdmin)
	var filteredAlerts []*entity.AssetAlert
	for _, alert := range alerts {
		if alert.TenantID == tenantID {
			// Apply additional filters
			if s.matchesFilter(alert, filter) {
				filteredAlerts = append(filteredAlerts, alert)
			}
		}
	}

	// Map to response
	responses := make([]dto.AssetAlertResponse, len(filteredAlerts))
	for i, alert := range filteredAlerts {
		// Get related entities for enhanced response
		asset, _ := s.assetRepo.GetByID(ctx, alert.AssetID)
		assetSensor, _ := s.assetSensorRepo.GetByID(ctx, alert.AssetSensorID)
		threshold, _ := s.thresholdRepo.GetByID(ctx, alert.ThresholdID)
		var sensorType *repository.SensorType
		var location *entity.Location

		if assetSensor != nil {
			sensorType, _ = s.sensorTypeRepo.GetByID(assetSensor.SensorTypeID)
		}
		if asset != nil {
			location, _ = s.locationRepo.GetByID(ctx, asset.LocationID)
		}

		responses[i] = *s.mapToResponse(alert, asset, assetSensor, threshold, sensorType, location)
	}

	totalPages := int(math.Ceil(float64(len(filteredAlerts)) / float64(filter.Limit)))

	return &dto.AssetAlertListResponse{
		Alerts:     responses,
		Page:       filter.Page,
		Limit:      filter.Limit,
		Total:      int64(len(filteredAlerts)),
		TotalPages: totalPages,
	}, nil
}

// GetAlertsByAsset retrieves alerts for a specific asset
func (s *assetAlertService) GetAlertsByAsset(ctx context.Context, tenantID, assetID uuid.UUID) ([]*dto.AssetAlertResponse, error) {
	alerts, err := s.alertRepo.GetByAssetID(ctx, assetID)
	if err != nil {
		return nil, fmt.Errorf("failed to get alerts by asset: %w", err)
	}

	responses := make([]*dto.AssetAlertResponse, 0, len(alerts))
	for _, alert := range alerts {
		// Filter by tenant
		if alert.TenantID == tenantID {
			// Get related entities for enhanced response
			asset, _ := s.assetRepo.GetByID(ctx, alert.AssetID)
			assetSensor, _ := s.assetSensorRepo.GetByID(ctx, alert.AssetSensorID)
			threshold, _ := s.thresholdRepo.GetByID(ctx, alert.ThresholdID)
			var sensorType *repository.SensorType
			var location *entity.Location

			if assetSensor != nil {
				sensorType, _ = s.sensorTypeRepo.GetByID(assetSensor.SensorTypeID)
			}
			if asset != nil {
				location, _ = s.locationRepo.GetByID(ctx, asset.LocationID)
			}

			responses = append(responses, s.mapToResponse(alert, asset, assetSensor, threshold, sensorType, location))
		}
	}

	return responses, nil
}

// GetAlertsByAssetSensor retrieves alerts for a specific asset sensor
func (s *assetAlertService) GetAlertsByAssetSensor(ctx context.Context, tenantID, assetSensorID uuid.UUID) ([]*dto.AssetAlertResponse, error) {
	alerts, err := s.alertRepo.GetByAssetSensorID(ctx, assetSensorID)
	if err != nil {
		return nil, fmt.Errorf("failed to get alerts by asset sensor: %w", err)
	}

	responses := make([]*dto.AssetAlertResponse, 0, len(alerts))
	for _, alert := range alerts {
		// Filter by tenant
		if alert.TenantID == tenantID {
			// Get related entities for enhanced response
			asset, _ := s.assetRepo.GetByID(ctx, alert.AssetID)
			assetSensor, _ := s.assetSensorRepo.GetByID(ctx, alert.AssetSensorID)
			threshold, _ := s.thresholdRepo.GetByID(ctx, alert.ThresholdID)
			var sensorType *repository.SensorType
			var location *entity.Location

			if assetSensor != nil {
				sensorType, _ = s.sensorTypeRepo.GetByID(assetSensor.SensorTypeID)
			}
			if asset != nil {
				location, _ = s.locationRepo.GetByID(ctx, asset.LocationID)
			}

			responses = append(responses, s.mapToResponse(alert, asset, assetSensor, threshold, sensorType, location))
		}
	}

	return responses, nil
}

// GetAlertsByThreshold retrieves alerts for a specific threshold
func (s *assetAlertService) GetAlertsByThreshold(ctx context.Context, tenantID, thresholdID uuid.UUID) ([]*dto.AssetAlertResponse, error) {
	alerts, err := s.alertRepo.GetByThresholdID(ctx, thresholdID)
	if err != nil {
		return nil, fmt.Errorf("failed to get alerts by threshold: %w", err)
	}

	responses := make([]*dto.AssetAlertResponse, 0, len(alerts))
	for _, alert := range alerts {
		// Filter by tenant
		if alert.TenantID == tenantID {
			// Get related entities for enhanced response
			asset, _ := s.assetRepo.GetByID(ctx, alert.AssetID)
			assetSensor, _ := s.assetSensorRepo.GetByID(ctx, alert.AssetSensorID)
			threshold, _ := s.thresholdRepo.GetByID(ctx, alert.ThresholdID)
			var sensorType *repository.SensorType
			var location *entity.Location

			if assetSensor != nil {
				sensorType, _ = s.sensorTypeRepo.GetByID(assetSensor.SensorTypeID)
			}
			if asset != nil {
				location, _ = s.locationRepo.GetByID(ctx, asset.LocationID)
			}

			responses = append(responses, s.mapToResponse(alert, asset, assetSensor, threshold, sensorType, location))
		}
	}

	return responses, nil
}

// GetUnresolvedAlerts retrieves all unresolved alerts
func (s *assetAlertService) GetUnresolvedAlerts(ctx context.Context, tenantID uuid.UUID) ([]*dto.AssetAlertResponse, error) {
	alerts, err := s.alertRepo.GetUnresolvedAlerts(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get unresolved alerts: %w", err)
	}

	responses := make([]*dto.AssetAlertResponse, 0, len(alerts))
	for _, alert := range alerts {
		// Filter by tenant
		if alert.TenantID == tenantID {
			// Get related entities for enhanced response
			asset, _ := s.assetRepo.GetByID(ctx, alert.AssetID)
			assetSensor, _ := s.assetSensorRepo.GetByID(ctx, alert.AssetSensorID)
			threshold, _ := s.thresholdRepo.GetByID(ctx, alert.ThresholdID)
			var sensorType *repository.SensorType
			var location *entity.Location

			if assetSensor != nil {
				sensorType, _ = s.sensorTypeRepo.GetByID(assetSensor.SensorTypeID)
			}
			if asset != nil {
				location, _ = s.locationRepo.GetByID(ctx, asset.LocationID)
			}

			responses = append(responses, s.mapToResponse(alert, asset, assetSensor, threshold, sensorType, location))
		}
	}

	return responses, nil
}

// GetAlertsBySeverity retrieves alerts by severity level
func (s *assetAlertService) GetAlertsBySeverity(ctx context.Context, tenantID uuid.UUID, severity string) ([]*dto.AssetAlertResponse, error) {
	alerts, err := s.alertRepo.GetAlertsBySeverity(ctx, entity.ThresholdSeverity(severity))
	if err != nil {
		return nil, fmt.Errorf("failed to get alerts by severity: %w", err)
	}

	responses := make([]*dto.AssetAlertResponse, 0, len(alerts))
	for _, alert := range alerts {
		// Filter by tenant
		if alert.TenantID == tenantID {
			// Get related entities for enhanced response
			asset, _ := s.assetRepo.GetByID(ctx, alert.AssetID)
			assetSensor, _ := s.assetSensorRepo.GetByID(ctx, alert.AssetSensorID)
			threshold, _ := s.thresholdRepo.GetByID(ctx, alert.ThresholdID)
			var sensorType *repository.SensorType
			var location *entity.Location

			if assetSensor != nil {
				sensorType, _ = s.sensorTypeRepo.GetByID(assetSensor.SensorTypeID)
			}
			if asset != nil {
				location, _ = s.locationRepo.GetByID(ctx, asset.LocationID)
			}

			responses = append(responses, s.mapToResponse(alert, asset, assetSensor, threshold, sensorType, location))
		}
	}

	return responses, nil
}

// ResolveAlert marks an alert as resolved
func (s *assetAlertService) ResolveAlert(ctx context.Context, tenantID, alertID uuid.UUID) error {
	alert, err := s.alertRepo.GetByID(ctx, alertID)
	if err != nil {
		return fmt.Errorf("failed to get alert: %w", err)
	}
	if alert == nil {
		return common.NewNotFoundError("alert", alertID.String())
	}

	// Verify tenant access
	if alert.TenantID != tenantID {
		return common.NewNotFoundError("alert", alertID.String())
	}

	// Check if already resolved
	if alert.ResolvedTime != nil {
		return common.NewBadRequestError("Alert is already resolved")
	}

	// Use the repository's ResolveAlert method which handles the update
	if err := s.alertRepo.ResolveAlert(ctx, alertID); err != nil {
		return fmt.Errorf("failed to resolve alert: %w", err)
	}

	return nil
}

// BulkResolveAlerts resolves multiple alerts
func (s *assetAlertService) BulkResolveAlerts(ctx context.Context, tenantID uuid.UUID, req *dto.BulkResolveAlertsRequest) (*dto.BulkResolveAlertsResponse, error) {
	response := &dto.BulkResolveAlertsResponse{
		ResolvedCount: 0,
		FailedCount:   0,
		FailedAlerts:  make([]uuid.UUID, 0),
	}

	for _, alertID := range req.AlertIDs {
		err := s.ResolveAlert(ctx, tenantID, alertID)
		if err != nil {
			response.FailedCount++
			response.FailedAlerts = append(response.FailedAlerts, alertID)
		} else {
			response.ResolvedCount++
		}
	}

	if response.FailedCount == 0 {
		response.Message = fmt.Sprintf("Successfully resolved %d alerts", response.ResolvedCount)
	} else {
		response.Message = fmt.Sprintf("Resolved %d alerts, %d failed", response.ResolvedCount, response.FailedCount)
	}

	return response, nil
}

// GetAlertStatistics returns statistics about alerts
func (s *assetAlertService) GetAlertStatistics(ctx context.Context, tenantID uuid.UUID, req *dto.AlertStatisticsRequest) (*dto.AlertStatistics, error) {
	// Get alerts in time range or all alerts
	var alerts []*entity.AssetAlert
	var err error

	if req.StartTime != nil && req.EndTime != nil {
		alerts, err = s.alertRepo.GetAlertsInTimeRange(ctx, *req.StartTime, *req.EndTime)
	} else {
		// Get all alerts for the tenant (this is a simplified approach)
		alerts, err = s.alertRepo.List(ctx, 1, 10000) // Large page size to get all
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get alerts: %w", err)
	}

	// Filter by tenant
	var tenantAlerts []*entity.AssetAlert
	for _, alert := range alerts {
		if alert.TenantID == tenantID {
			// Apply asset filter if provided
			if req.AssetID == nil || alert.AssetID == *req.AssetID {
				tenantAlerts = append(tenantAlerts, alert)
			}
		}
	}

	stats := &dto.AlertStatistics{
		TotalAlerts:      len(tenantAlerts),
		AlertsBySeverity: make(map[string]int),
		AlertsByAsset:    make(map[string]int),
	}

	for _, alert := range tenantAlerts {
		// Count by resolution status
		if alert.ResolvedTime == nil {
			stats.UnresolvedAlerts++
		} else {
			stats.ResolvedAlerts++
		}

		// Count by severity
		severityStr := string(alert.Severity)
		stats.AlertsBySeverity[severityStr]++

		switch alert.Severity {
		case entity.ThresholdSeverityWarning:
			stats.WarningAlerts++
		case entity.ThresholdSeverityCritical:
			stats.CriticalAlerts++
		}

		// Count by asset (simplified - using asset ID as key)
		assetKey := alert.AssetID.String()
		stats.AlertsByAsset[assetKey]++
	}

	// Get recent alerts (last 10)
	recentCount := 10
	if len(tenantAlerts) < recentCount {
		recentCount = len(tenantAlerts)
	}

	stats.RecentAlerts = make([]dto.AssetAlertResponse, recentCount)
	for i := 0; i < recentCount; i++ {
		alert := tenantAlerts[i]
		// Get related entities for enhanced response
		asset, _ := s.assetRepo.GetByID(ctx, alert.AssetID)
		assetSensor, _ := s.assetSensorRepo.GetByID(ctx, alert.AssetSensorID)
		threshold, _ := s.thresholdRepo.GetByID(ctx, alert.ThresholdID)
		var sensorType *repository.SensorType
		var location *entity.Location

		if assetSensor != nil {
			sensorType, _ = s.sensorTypeRepo.GetByID(assetSensor.SensorTypeID)
		}
		if asset != nil {
			location, _ = s.locationRepo.GetByID(ctx, asset.LocationID)
		}

		stats.RecentAlerts[i] = *s.mapToResponse(alert, asset, assetSensor, threshold, sensorType, location)
	}

	return stats, nil
}

// GetAlertSummary returns a summary of alerts for dashboard
func (s *assetAlertService) GetAlertSummary(ctx context.Context, tenantID uuid.UUID) (*dto.AlertSummaryResponse, error) {
	// Get all alerts for the tenant (simplified approach)
	alerts, err := s.alertRepo.List(ctx, 1, 10000) // Large page size to get all
	if err != nil {
		return nil, fmt.Errorf("failed to get alerts: %w", err)
	}

	// Filter by tenant
	var tenantAlerts []*entity.AssetAlert
	for _, alert := range alerts {
		if alert.TenantID == tenantID {
			tenantAlerts = append(tenantAlerts, alert)
		}
	}

	summary := &dto.AlertSummaryResponse{
		TotalAlerts:       len(tenantAlerts),
		TopAffectedAssets: make([]dto.AlertAssetSummary, 0),
	}

	assetAlertCount := make(map[uuid.UUID]int)
	assetLastAlert := make(map[uuid.UUID]time.Time)

	for _, alert := range tenantAlerts {
		// Count by resolution status
		if alert.ResolvedTime == nil {
			summary.UnresolvedAlerts++
		}

		// Count by severity
		switch alert.Severity {
		case entity.ThresholdSeverityWarning:
			summary.WarningAlerts++
		case entity.ThresholdSeverityCritical:
			summary.CriticalAlerts++
		}

		// Track by asset
		assetAlertCount[alert.AssetID]++
		if assetLastAlert[alert.AssetID].Before(alert.AlertTime) {
			assetLastAlert[alert.AssetID] = alert.AlertTime
		}
	}

	// Create top affected assets (simplified - just take first few)
	count := 0
	maxAssets := 5
	for assetID, alertCount := range assetAlertCount {
		if count >= maxAssets {
			break
		}

		asset, _ := s.assetRepo.GetByID(ctx, assetID)
		assetName := assetID.String() // fallback
		if asset != nil {
			assetName = asset.Name
		}

		summary.TopAffectedAssets = append(summary.TopAffectedAssets, dto.AlertAssetSummary{
			AssetID:    assetID,
			AssetName:  assetName,
			AlertCount: alertCount,
			LastAlert:  assetLastAlert[assetID],
		})
		count++
	}

	return summary, nil
}

// GetAlertsInTimeRange retrieves alerts within a specific time range
func (s *assetAlertService) GetAlertsInTimeRange(ctx context.Context, tenantID uuid.UUID, startTime, endTime time.Time) ([]*dto.AssetAlertResponse, error) {
	alerts, err := s.alertRepo.GetAlertsInTimeRange(ctx, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to get alerts in time range: %w", err)
	}

	responses := make([]*dto.AssetAlertResponse, 0, len(alerts))
	for _, alert := range alerts {
		// Filter by tenant
		if alert.TenantID == tenantID {
			// Get related entities for enhanced response
			asset, _ := s.assetRepo.GetByID(ctx, alert.AssetID)
			assetSensor, _ := s.assetSensorRepo.GetByID(ctx, alert.AssetSensorID)
			threshold, _ := s.thresholdRepo.GetByID(ctx, alert.ThresholdID)
			var sensorType *repository.SensorType
			var location *entity.Location

			if assetSensor != nil {
				sensorType, _ = s.sensorTypeRepo.GetByID(assetSensor.SensorTypeID)
			}
			if asset != nil {
				location, _ = s.locationRepo.GetByID(ctx, asset.LocationID)
			}

			responses = append(responses, s.mapToResponse(alert, asset, assetSensor, threshold, sensorType, location))
		}
	}

	return responses, nil
}

// Helper methods

// matchesFilter checks if an alert matches the given filter criteria
func (s *assetAlertService) matchesFilter(alert *entity.AssetAlert, filter *dto.AssetAlertFilterRequest) bool {
	if filter.AssetID != nil && alert.AssetID != *filter.AssetID {
		return false
	}
	if filter.AssetSensorID != nil && alert.AssetSensorID != *filter.AssetSensorID {
		return false
	}
	if filter.ThresholdID != nil && alert.ThresholdID != *filter.ThresholdID {
		return false
	}
	if filter.Severity != nil && string(alert.Severity) != *filter.Severity {
		return false
	}
	if filter.IsResolved != nil {
		isResolved := alert.ResolvedTime != nil
		if isResolved != *filter.IsResolved {
			return false
		}
	}
	if filter.StartTime != nil && alert.AlertTime.Before(*filter.StartTime) {
		return false
	}
	if filter.EndTime != nil && alert.AlertTime.After(*filter.EndTime) {
		return false
	}
	return true
}

// mapToResponse maps entity to response DTO
func (s *assetAlertService) mapToResponse(
	alert *entity.AssetAlert,
	asset *entity.Asset,
	assetSensor *repository.AssetSensorWithDetails,
	threshold *entity.SensorThreshold,
	sensorType *repository.SensorType,
	location *entity.Location,
) *dto.AssetAlertResponse {
	response := &dto.AssetAlertResponse{
		ID:            alert.ID,
		TenantID:      alert.TenantID,
		AssetID:       alert.AssetID,
		AssetSensorID: alert.AssetSensorID,
		ThresholdID:   alert.ThresholdID,
		AlertTime:     alert.AlertTime,
		ResolvedTime:  alert.ResolvedTime,
		Severity:      string(alert.Severity),
		IsResolved:    alert.ResolvedTime != nil,
	}

	// Add related entity names if available
	if asset != nil {
		response.AssetName = asset.Name
	}
	if assetSensor != nil {
		response.AssetSensorName = assetSensor.Name
	}
	if threshold != nil {
		response.ThresholdName = threshold.Name
	}
	if sensorType != nil {
		response.SensorTypeName = sensorType.Name
	}
	if location != nil {
		response.LocationName = location.Name
	}

	return response
}
