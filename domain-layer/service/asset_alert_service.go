package service

import (
	"be-lecsens/asset_management/data-layer/entity"
	"be-lecsens/asset_management/data-layer/repository"
	"be-lecsens/asset_management/helpers/common"
	"be-lecsens/asset_management/helpers/dto"
	"context"
	"fmt"
	"log"
	"math"

	"github.com/google/uuid"
)

// AssetAlertService handles business logic for asset alerts
type AssetAlertService struct {
	assetAlertRepo  repository.AssetAlertRepository
	assetRepo       repository.AssetRepository
	assetSensorRepo repository.AssetSensorRepository
}

// NewAssetAlertService creates a new instance of AssetAlertService
func NewAssetAlertService(
	assetAlertRepo repository.AssetAlertRepository,
	assetRepo repository.AssetRepository,
	assetSensorRepo repository.AssetSensorRepository,
) *AssetAlertService {
	return &AssetAlertService{
		assetAlertRepo:  assetAlertRepo,
		assetRepo:       assetRepo,
		assetSensorRepo: assetSensorRepo,
	}
}

// CreateAssetAlert creates a new asset alert
func (s *AssetAlertService) CreateAssetAlert(ctx context.Context, alert *entity.AssetAlert) (*dto.AssetAlertResponse, error) {
	log.Printf("Creating asset alert: %+v", alert)

	// Validate asset exists
	if s.assetRepo != nil {
		asset, err := s.assetRepo.GetByID(ctx, alert.AssetID)
		if err != nil {
			return nil, fmt.Errorf("failed to validate asset: %w", err)
		}
		if asset == nil {
			return nil, common.NewValidationError("asset not found", nil)
		}
	}

	// Validate asset sensor exists
	if s.assetSensorRepo != nil {
		assetSensor, err := s.assetSensorRepo.GetByID(ctx, alert.AssetSensorID)
		if err != nil {
			return nil, fmt.Errorf("failed to validate asset sensor: %w", err)
		}
		if assetSensor == nil {
			return nil, common.NewValidationError("asset sensor not found", nil)
		}
	}

	// Create the alert
	if err := s.assetAlertRepo.Create(ctx, alert); err != nil {
		log.Printf("Error creating asset alert: %v", err)
		return nil, fmt.Errorf("failed to create asset alert: %w", err)
	}

	log.Printf("Successfully created asset alert with ID: %s", alert.ID)
	return s.toResponseDTO(alert), nil
}

// GetAssetAlertByID retrieves an asset alert by its ID
func (s *AssetAlertService) GetAssetAlertByID(ctx context.Context, id uuid.UUID) (*dto.AssetAlertResponse, error) {
	alert, err := s.assetAlertRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get asset alert: %w", err)
	}

	if alert == nil {
		return nil, common.NewNotFoundError("asset alert", id.String())
	}

	return s.toResponseDTO(alert), nil
}

// GetAssetAlertsByTenant retrieves all asset alerts for a tenant
func (s *AssetAlertService) GetAssetAlertsByTenant(ctx context.Context, tenantID uuid.UUID) ([]*entity.AssetAlert, error) {
	alerts, err := s.assetAlertRepo.GetByTenantID(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get asset alerts: %w", err)
	}

	return alerts, nil
}

// GetAssetAlertsByAsset retrieves alerts for a specific asset
func (s *AssetAlertService) GetAssetAlertsByAsset(ctx context.Context, assetID uuid.UUID) ([]*entity.AssetAlert, error) {
	alerts, err := s.assetAlertRepo.GetByAssetID(ctx, assetID)
	if err != nil {
		return nil, fmt.Errorf("failed to get asset alerts: %w", err)
	}

	return alerts, nil
}

// GetAssetAlertsByAssetSensor retrieves alerts for a specific asset sensor
func (s *AssetAlertService) GetAssetAlertsByAssetSensor(ctx context.Context, assetSensorID uuid.UUID) ([]*entity.AssetAlert, error) {
	alerts, err := s.assetAlertRepo.GetByAssetSensorID(ctx, assetSensorID)
	if err != nil {
		return nil, fmt.Errorf("failed to get asset alerts: %w", err)
	}

	return alerts, nil
}

// GetActiveAlerts retrieves all active (unresolved) alerts for a tenant
func (s *AssetAlertService) GetActiveAlerts(ctx context.Context, tenantID uuid.UUID) ([]*entity.AssetAlert, error) {
	alerts, err := s.assetAlertRepo.GetActiveAlerts(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get active alerts: %w", err)
	}

	return alerts, nil
}

// GetActiveAlertsByAssetSensor retrieves active alerts for a specific asset sensor
func (s *AssetAlertService) GetActiveAlertsByAssetSensor(ctx context.Context, assetSensorID uuid.UUID) ([]*entity.AssetAlert, error) {
	alerts, err := s.assetAlertRepo.GetActiveAlertsByAssetSensor(ctx, assetSensorID)
	if err != nil {
		return nil, fmt.Errorf("failed to get active alerts for asset sensor: %w", err)
	}

	return alerts, nil
}

// ResolveAlert marks an alert as resolved
func (s *AssetAlertService) ResolveAlert(ctx context.Context, id uuid.UUID) (*entity.AssetAlert, error) {
	log.Printf("Resolving asset alert with ID: %s", id)

	// Check if alert exists and is not already resolved
	alert, err := s.assetAlertRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get alert: %w", err)
	}
	if alert == nil {
		return nil, common.NewNotFoundError("asset alert", id.String())
	}
	if alert.IsResolved {
		return nil, common.NewValidationError("alert is already resolved", nil)
	}

	// Resolve the alert
	if err := s.assetAlertRepo.ResolveAlert(ctx, id); err != nil {
		log.Printf("Error resolving asset alert: %v", err)
		return nil, fmt.Errorf("failed to resolve asset alert: %w", err)
	}

	// Get updated alert
	resolvedAlert, err := s.assetAlertRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get resolved alert: %w", err)
	}

	log.Printf("Successfully resolved asset alert with ID: %s", id)
	return resolvedAlert, nil
}

// ResolveAlertsByAssetSensor resolves all active alerts for a specific asset sensor
func (s *AssetAlertService) ResolveAlertsByAssetSensor(ctx context.Context, assetSensorID uuid.UUID) (int, error) {
	log.Printf("Resolving all active alerts for asset sensor: %s", assetSensorID)

	// Get active alerts for the asset sensor
	activeAlerts, err := s.assetAlertRepo.GetActiveAlertsByAssetSensor(ctx, assetSensorID)
	if err != nil {
		return 0, fmt.Errorf("failed to get active alerts: %w", err)
	}

	resolvedCount := 0
	for _, alert := range activeAlerts {
		if err := s.assetAlertRepo.ResolveAlert(ctx, alert.ID); err != nil {
			log.Printf("Failed to resolve alert %s: %v", alert.ID, err)
			continue
		}
		resolvedCount++
	}

	log.Printf("Successfully resolved %d alerts for asset sensor %s", resolvedCount, assetSensorID)
	return resolvedCount, nil
}

// UpdateAlert updates an existing asset alert
func (s *AssetAlertService) UpdateAlert(ctx context.Context, alert *entity.AssetAlert) (*entity.AssetAlert, error) {
	log.Printf("Updating asset alert: %+v", alert)

	// Check if alert exists
	existing, err := s.assetAlertRepo.GetByID(ctx, alert.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing alert: %w", err)
	}
	if existing == nil {
		return nil, common.NewNotFoundError("asset alert", alert.ID.String())
	}

	// Update the alert
	if err := s.assetAlertRepo.Update(ctx, alert); err != nil {
		log.Printf("Error updating asset alert: %v", err)
		return nil, fmt.Errorf("failed to update asset alert: %w", err)
	}

	// Get updated alert
	updatedAlert, err := s.assetAlertRepo.GetByID(ctx, alert.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated alert: %w", err)
	}

	log.Printf("Successfully updated asset alert with ID: %s", alert.ID)
	return updatedAlert, nil
}

// DeleteAlert deletes an asset alert
func (s *AssetAlertService) DeleteAlert(ctx context.Context, id uuid.UUID) error {
	log.Printf("Deleting asset alert with ID: %s", id)

	// Check if alert exists
	existing, err := s.assetAlertRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to check existing alert: %w", err)
	}
	if existing == nil {
		return common.NewNotFoundError("asset alert", id.String())
	}

	// Delete the alert
	if err := s.assetAlertRepo.Delete(ctx, id); err != nil {
		log.Printf("Error deleting asset alert: %v", err)
		return fmt.Errorf("failed to delete asset alert: %w", err)
	}

	log.Printf("Successfully deleted asset alert with ID: %s", id)
	return nil
}

// ListAlerts retrieves paginated asset alerts for a tenant
func (s *AssetAlertService) ListAlerts(ctx context.Context, tenantID uuid.UUID, limit, offset int) ([]*entity.AssetAlert, int, error) {
	if limit <= 0 {
		limit = 20 // Default limit
	}
	if limit > 100 {
		limit = 100 // Maximum limit
	}

	alerts, totalCount, err := s.assetAlertRepo.List(ctx, tenantID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list asset alerts: %w", err)
	}

	return alerts, totalCount, nil
}

// GetAlertStatistics retrieves alert statistics for a tenant
func (s *AssetAlertService) GetAlertStatistics(ctx context.Context, tenantID uuid.UUID, filter dto.AssetAlertFilter) (*dto.AssetAlertStatisticsResponse, error) {
	stats, err := s.assetAlertRepo.GetAlertStatistics(ctx, tenantID, filter.AssetID, filter.FromTime, filter.ToTime)
	if err != nil {
		return nil, fmt.Errorf("failed to get alert statistics: %w", err)
	}

	return &dto.AssetAlertStatisticsResponse{
		TotalAlerts:    stats["total_alerts"].(int),
		ActiveAlerts:   stats["active_alerts"].(int),
		ResolvedAlerts: stats["resolved_alerts"].(int),
		CriticalAlerts: stats["critical_alerts"].(int),
		WarningAlerts:  stats["warning_alerts"].(int),
		Alerts24h:      stats["alerts_24h"].(int),
		Alerts7d:       stats["alerts_7d"].(int),
	}, nil
}

// GetCriticalAlerts retrieves all critical unresolved alerts for a tenant
func (s *AssetAlertService) GetCriticalAlerts(ctx context.Context, tenantID uuid.UUID) ([]*entity.AssetAlert, error) {
	allActiveAlerts, err := s.assetAlertRepo.GetActiveAlerts(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get active alerts: %w", err)
	}

	var criticalAlerts []*entity.AssetAlert
	for _, alert := range allActiveAlerts {
		if alert.Severity == entity.ThresholdSeverityCritical {
			criticalAlerts = append(criticalAlerts, alert)
		}
	}

	return criticalAlerts, nil
}

// GetRecentAlerts retrieves recent alerts (last 24 hours) for a tenant
func (s *AssetAlertService) GetRecentAlerts(ctx context.Context, tenantID uuid.UUID) ([]*entity.AssetAlert, error) {
	allAlerts, err := s.assetAlertRepo.GetByTenantID(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get alerts: %w", err)
	}

	var recentAlerts []*entity.AssetAlert
	for _, alert := range allAlerts {
		if alert.GetDuration().Hours() <= 24 {
			recentAlerts = append(recentAlerts, alert)
		}
	}

	return recentAlerts, nil
}

// BulkResolveAlerts resolves multiple alerts by their IDs
func (s *AssetAlertService) BulkResolveAlerts(ctx context.Context, alertIDs []uuid.UUID) (int, error) {
	resolvedCount := 0
	for _, id := range alertIDs {
		if err := s.assetAlertRepo.ResolveAlert(ctx, id); err != nil {
			log.Printf("Failed to resolve alert %s: %v", id, err)
			continue
		}
		resolvedCount++
	}

	return resolvedCount, nil
}

// ListAssetAlerts retrieves paginated asset alerts with filters
func (s *AssetAlertService) ListAssetAlerts(ctx context.Context, tenantID uuid.UUID, filter dto.AssetAlertFilter) (*dto.AssetAlertListResponse, error) {
	// Set default pagination values
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.Limit <= 0 {
		filter.Limit = 20
	}
	if filter.Limit > 100 {
		filter.Limit = 100
	}

	offset := (filter.Page - 1) * filter.Limit

	alerts, totalCount, err := s.assetAlertRepo.ListWithFilters(
		ctx,
		tenantID,
		filter.Limit,
		offset,
		filter.AssetID,
		filter.AssetSensorID,
		filter.Severity,
		filter.IsResolved,
		filter.FromTime,
		filter.ToTime,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list asset alerts: %w", err)
	}

	// Convert to DTOs
	var responseDTOs []dto.AssetAlertResponse
	for _, alert := range alerts {
		responseDTOs = append(responseDTOs, *s.toResponseDTO(alert))
	}

	// Calculate pagination info
	totalPages := int(math.Ceil(float64(totalCount) / float64(filter.Limit)))

	return &dto.AssetAlertListResponse{
		Data: responseDTOs,
		Pagination: dto.PaginationInfo{
			Page:        filter.Page,
			Limit:       filter.Limit,
			TotalItems:  int64(totalCount),
			TotalPages:  totalPages,
			HasNext:     filter.Page < totalPages,
			HasPrevious: filter.Page > 1,
		},
	}, nil
}

// ResolveAssetAlert resolves an asset alert
func (s *AssetAlertService) ResolveAssetAlert(ctx context.Context, id uuid.UUID) (*dto.AssetAlertResponse, error) {
	log.Printf("Resolving asset alert with ID: %s", id)

	// Check if alert exists and is not already resolved
	alert, err := s.assetAlertRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get alert: %w", err)
	}
	if alert == nil {
		return nil, common.NewNotFoundError("asset alert", id.String())
	}
	if alert.IsResolved {
		return nil, common.NewValidationError("alert is already resolved", nil)
	}

	// Resolve the alert
	if err := s.assetAlertRepo.ResolveAlert(ctx, id); err != nil {
		log.Printf("Error resolving asset alert: %v", err)
		return nil, fmt.Errorf("failed to resolve asset alert: %w", err)
	}

	// Get updated alert
	resolvedAlert, err := s.assetAlertRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get resolved alert: %w", err)
	}

	log.Printf("Successfully resolved asset alert with ID: %s", id)
	return s.toResponseDTO(resolvedAlert), nil
}

// ResolveMultipleAssetAlerts resolves multiple asset alerts
func (s *AssetAlertService) ResolveMultipleAssetAlerts(ctx context.Context, request dto.ResolveMultipleAlertsRequest) (*dto.ResolveMultipleAlertsResponse, error) {
	successCount, failureCount, err := s.assetAlertRepo.ResolveMultipleAlerts(ctx, request.AlertIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve multiple alerts: %w", err)
	}

	return &dto.ResolveMultipleAlertsResponse{
		ResolvedCount:  successCount,
		FailedCount:    failureCount,
		TotalRequested: len(request.AlertIDs),
	}, nil
}

// DeleteAssetAlert deletes an asset alert
func (s *AssetAlertService) DeleteAssetAlert(ctx context.Context, id uuid.UUID) error {
	log.Printf("Deleting asset alert with ID: %s", id)

	// Check if alert exists
	existing, err := s.assetAlertRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to check existing alert: %w", err)
	}
	if existing == nil {
		return common.NewNotFoundError("asset alert", id.String())
	}

	// Delete the alert
	if err := s.assetAlertRepo.Delete(ctx, id); err != nil {
		log.Printf("Error deleting asset alert: %v", err)
		return fmt.Errorf("failed to delete asset alert: %w", err)
	}

	log.Printf("Successfully deleted asset alert with ID: %s", id)
	return nil
}

// DeleteMultipleAssetAlerts deletes multiple asset alerts
func (s *AssetAlertService) DeleteMultipleAssetAlerts(ctx context.Context, request dto.DeleteMultipleAlertsRequest) (*dto.DeleteMultipleAlertsResponse, error) {
	successCount, failureCount, err := s.assetAlertRepo.DeleteMultipleAlerts(ctx, request.AlertIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to delete multiple alerts: %w", err)
	}

	return &dto.DeleteMultipleAlertsResponse{
		DeletedCount:   successCount,
		FailedCount:    failureCount,
		TotalRequested: len(request.AlertIDs),
	}, nil
}

// GetAlertsByMeasurementType retrieves alerts for a specific measurement type
func (s *AssetAlertService) GetAlertsByMeasurementType(ctx context.Context, measurementTypeID uuid.UUID) ([]*entity.AssetAlert, error) {
	alerts, err := s.assetAlertRepo.GetByMeasurementTypeID(ctx, measurementTypeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get alerts by measurement type: %w", err)
	}

	return alerts, nil
}

// ListAllAssetAlerts retrieves all asset alerts across all tenants
func (s *AssetAlertService) ListAllAssetAlerts(ctx context.Context, limit, offset int) ([]*entity.AssetAlert, int, error) {
	if limit <= 0 {
		limit = 20 // Default limit
	}
	if limit > 100 {
		limit = 100 // Maximum limit
	}

	alerts, totalCount, err := s.assetAlertRepo.ListAll(ctx, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list all asset alerts: %w", err)
	}

	return alerts, totalCount, nil
}

// GetGlobalAlertStatistics retrieves alert statistics across all tenants
func (s *AssetAlertService) GetGlobalAlertStatistics(ctx context.Context) (*dto.AssetAlertStatisticsResponse, error) {
	stats, err := s.assetAlertRepo.GetGlobalAlertStatistics(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get global alert statistics: %w", err)
	}

	return &dto.AssetAlertStatisticsResponse{
		TotalAlerts:    stats["total_alerts"].(int),
		ActiveAlerts:   stats["active_alerts"].(int),
		ResolvedAlerts: stats["resolved_alerts"].(int),
		CriticalAlerts: stats["critical_alerts"].(int),
		WarningAlerts:  stats["warning_alerts"].(int),
		Alerts24h:      stats["alerts_24h"].(int),
		Alerts7d:       stats["alerts_7d"].(int),
		TotalTenants:   stats["total_tenants"].(int),
	}, nil
}

// Helper methods

// toResponseDTO converts entity to response DTO
func (s *AssetAlertService) toResponseDTO(alert *entity.AssetAlert) *dto.AssetAlertResponse {
	return &dto.AssetAlertResponse{
		ID:                   alert.ID,
		TenantID:             alert.TenantID,
		AssetID:              alert.AssetID,
		AssetSensorID:        alert.AssetSensorID,
		ThresholdID:          alert.ThresholdID,
		MeasurementFieldName: alert.MeasurementFieldName,
		AlertTime:            alert.AlertTime,
		ResolvedTime:         alert.ResolvedTime,
		Severity:             alert.Severity,
		TriggerValue:         alert.TriggerValue,
		ThresholdMinValue:    alert.ThresholdMinValue,
		ThresholdMaxValue:    alert.ThresholdMaxValue,
		AlertMessage:         alert.AlertMessage,
		AlertType:            alert.AlertType,
		IsResolved:           alert.IsResolved,
		CreatedAt:            alert.CreatedAt,
		UpdatedAt:            alert.UpdatedAt,
	}
}
