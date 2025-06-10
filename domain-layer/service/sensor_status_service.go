package service

import (
	"be-lecsens/asset_management/data-layer/entity"
	"be-lecsens/asset_management/data-layer/repository"
	"be-lecsens/asset_management/helpers/common"
	"be-lecsens/asset_management/helpers/dto"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// SensorStatusService handles business logic for sensor status operations
type SensorStatusService struct {
	repo repository.SensorStatusRepository
}

// NewSensorStatusService creates a new instance of SensorStatusService
func NewSensorStatusService(repo repository.SensorStatusRepository) *SensorStatusService {
	return &SensorStatusService{
		repo: repo,
	}
}

// CreateSensorStatus creates a new sensor status record
func (s *SensorStatusService) CreateSensorStatus(ctx context.Context, req dto.CreateSensorStatusRequest) (*dto.SensorStatusDTO, error) {
	// Get asset sensor context for tenant inheritance
	tenantID, _, err := s.repo.GetAssetSensorContext(ctx, req.AssetSensorID)
	if err != nil {
		return nil, fmt.Errorf("failed to get asset sensor context: %v", err)
	}

	// Convert request to entity
	status := req.ToEntity(tenantID)

	// Create in repository
	err = s.repo.Create(ctx, status)
	if err != nil {
		return nil, fmt.Errorf("failed to create sensor status: %v", err)
	}

	// Convert to DTO and return
	return dto.FromSensorStatusEntity(status), nil
}

// GetSensorStatus retrieves a sensor status by ID
func (s *SensorStatusService) GetSensorStatus(ctx context.Context, id uuid.UUID) (*dto.SensorStatusDTO, error) {
	status, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get sensor status: %v", err)
	}
	if status == nil {
		return nil, fmt.Errorf("sensor status not found")
	}

	return dto.FromSensorStatusEntity(status), nil
}

// GetSensorStatusBySensorID retrieves the current sensor status by asset sensor ID
func (s *SensorStatusService) GetSensorStatusBySensorID(ctx context.Context, assetSensorID uuid.UUID) (*dto.SensorStatusDTO, error) {
	status, err := s.repo.GetBySensorID(ctx, assetSensorID)
	if err != nil {
		return nil, fmt.Errorf("failed to get sensor status by sensor ID: %v", err)
	}
	if status == nil {
		return nil, fmt.Errorf("sensor status not found for sensor ID")
	}

	return dto.FromSensorStatusEntity(status), nil
}

// GetOnlineSensors retrieves all online sensors with pagination
func (s *SensorStatusService) GetOnlineSensors(ctx context.Context, params common.QueryParams) ([]dto.SensorStatusDTO, *common.PaginationResponse, error) {
	statuses, pagination, err := s.repo.GetOnlineSensors(ctx, params)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get online sensors: %v", err)
	}

	dtos := dto.FromSensorStatusEntityList(statuses)
	return dtos, pagination, nil
}

// GetOfflineSensors retrieves all offline sensors with pagination
func (s *SensorStatusService) GetOfflineSensors(ctx context.Context, params common.QueryParams) ([]dto.SensorStatusDTO, *common.PaginationResponse, error) {
	statuses, pagination, err := s.repo.GetOfflineSensors(ctx, params)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get offline sensors: %v", err)
	}

	dtos := dto.FromSensorStatusEntityList(statuses)
	return dtos, pagination, nil
}

// GetLowBatterySensors retrieves sensors with low battery levels
func (s *SensorStatusService) GetLowBatterySensors(ctx context.Context, threshold float64, params common.QueryParams) ([]dto.SensorStatusDTO, *common.PaginationResponse, error) {
	if threshold <= 0 || threshold > 100 {
		threshold = 20.0 // Default threshold: 20%
	}

	statuses, pagination, err := s.repo.GetLowBatterySensors(ctx, threshold, params)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get low battery sensors: %v", err)
	}

	dtos := dto.FromSensorStatusEntityList(statuses)
	return dtos, pagination, nil
}

// GetWeakSignalSensors retrieves sensors with weak signal strength
func (s *SensorStatusService) GetWeakSignalSensors(ctx context.Context, threshold int, params common.QueryParams) ([]dto.SensorStatusDTO, *common.PaginationResponse, error) {
	if threshold >= 0 || threshold < -120 {
		threshold = -70 // Default threshold: -70 dBm
	}

	statuses, pagination, err := s.repo.GetWeakSignalSensors(ctx, threshold, params)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get weak signal sensors: %v", err)
	}

	dtos := dto.FromSensorStatusEntityList(statuses)
	return dtos, pagination, nil
}

// GetUnhealthySensors retrieves sensors that are considered unhealthy
func (s *SensorStatusService) GetUnhealthySensors(ctx context.Context, params common.QueryParams) ([]dto.SensorStatusDTO, *common.PaginationResponse, error) {
	statuses, pagination, err := s.repo.GetUnhealthySensors(ctx, params)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get unhealthy sensors: %v", err)
	}

	dtos := dto.FromSensorStatusEntityList(statuses)
	return dtos, pagination, nil
}

// UpdateSensorStatus updates an existing sensor status record
func (s *SensorStatusService) UpdateSensorStatus(ctx context.Context, id uuid.UUID, req dto.UpdateSensorStatusRequest) (*dto.SensorStatusDTO, error) {
	// Get existing sensor status
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing sensor status: %v", err)
	}
	if existing == nil {
		return nil, fmt.Errorf("sensor status not found")
	}

	// Create partial update entity from request
	updates := req.ToEntity(id)

	// Apply updates to existing entity
	s.applyUpdates(existing, updates)

	// Update in repository
	err = s.repo.Update(ctx, existing)
	if err != nil {
		return nil, fmt.Errorf("failed to update sensor status: %v", err)
	}

	return dto.FromSensorStatusEntity(existing), nil
}

// UpsertSensorStatus creates or updates a sensor status record
func (s *SensorStatusService) UpsertSensorStatus(ctx context.Context, req dto.CreateSensorStatusRequest) (*dto.SensorStatusDTO, error) {
	// Get asset sensor context for tenant inheritance
	tenantID, _, err := s.repo.GetAssetSensorContext(ctx, req.AssetSensorID)
	if err != nil {
		return nil, fmt.Errorf("failed to get asset sensor context: %v", err)
	}

	// Convert request to entity
	status := req.ToEntity(tenantID)

	// Upsert in repository
	err = s.repo.UpsertStatus(ctx, status)
	if err != nil {
		return nil, fmt.Errorf("failed to upsert sensor status: %v", err)
	}

	return dto.FromSensorStatusEntity(status), nil
}

// UpdateBatteryStatus updates only battery-related status fields
func (s *SensorStatusService) UpdateBatteryStatus(ctx context.Context, assetSensorID uuid.UUID, batteryLevel *float64, batteryVoltage *float64, batteryStatus *string) error {
	err := s.repo.UpdateBatteryStatus(ctx, assetSensorID, batteryLevel, batteryVoltage, batteryStatus)
	if err != nil {
		return fmt.Errorf("failed to update battery status: %v", err)
	}
	return nil
}

// UpdateSignalStatus updates only signal-related status fields
func (s *SensorStatusService) UpdateSignalStatus(ctx context.Context, assetSensorID uuid.UUID, rssi *int, snr *float64, quality *int, signalStatus *string) error {
	err := s.repo.UpdateSignalStatus(ctx, assetSensorID, rssi, snr, quality, signalStatus)
	if err != nil {
		return fmt.Errorf("failed to update signal status: %v", err)
	}
	return nil
}

// UpdateConnectionStatus updates only connection-related status fields
func (s *SensorStatusService) UpdateConnectionStatus(ctx context.Context, assetSensorID uuid.UUID, connectionStatus string, connectionType *string, currentIP *string, currentNetwork *string) error {
	err := s.repo.UpdateConnectionStatus(ctx, assetSensorID, connectionStatus, connectionType, currentIP, currentNetwork)
	if err != nil {
		return fmt.Errorf("failed to update connection status: %v", err)
	}
	return nil
}

// UpdateHeartbeat updates the last heartbeat timestamp for a sensor
func (s *SensorStatusService) UpdateHeartbeat(ctx context.Context, assetSensorID uuid.UUID) error {
	err := s.repo.UpdateHeartbeat(ctx, assetSensorID)
	if err != nil {
		return fmt.Errorf("failed to update heartbeat: %v", err)
	}
	return nil
}

// DeleteSensorStatus deletes a sensor status record
func (s *SensorStatusService) DeleteSensorStatus(ctx context.Context, id uuid.UUID) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete sensor status: %v", err)
	}
	return nil
}

// DeleteSensorStatusBySensorID deletes all status records for a specific sensor
func (s *SensorStatusService) DeleteSensorStatusBySensorID(ctx context.Context, assetSensorID uuid.UUID) error {
	err := s.repo.DeleteBySensorID(ctx, assetSensorID)
	if err != nil {
		return fmt.Errorf("failed to delete sensor status by sensor ID: %v", err)
	}
	return nil
}

// GetSensorHealthSummary provides aggregated health statistics
func (s *SensorStatusService) GetSensorHealthSummary(ctx context.Context, params common.QueryParams) (*dto.SensorHealthSummaryResponse, error) {
	// Get online sensors
	onlineSensors, _, err := s.repo.GetOnlineSensors(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get online sensors for summary: %v", err)
	}

	// Get offline sensors
	offlineSensors, _, err := s.repo.GetOfflineSensors(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get offline sensors for summary: %v", err)
	}

	// Get low battery sensors (threshold 20%)
	lowBatterySensors, _, err := s.repo.GetLowBatterySensors(ctx, 20.0, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get low battery sensors for summary: %v", err)
	}

	// Get critical battery sensors (threshold 10%)
	criticalBatterySensors, _, err := s.repo.GetLowBatterySensors(ctx, 10.0, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get critical battery sensors for summary: %v", err)
	}

	// Get weak signal sensors (threshold -70 dBm)
	weakSignalSensors, _, err := s.repo.GetWeakSignalSensors(ctx, -70, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get weak signal sensors for summary: %v", err)
	}

	// Get unhealthy sensors
	unhealthySensors, _, err := s.repo.GetUnhealthySensors(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get unhealthy sensors for summary: %v", err)
	}

	// Calculate totals
	totalOnline := len(onlineSensors)
	totalOffline := len(offlineSensors)
	totalSensors := totalOnline + totalOffline
	lowBattery := len(lowBatterySensors)
	criticalBattery := len(criticalBatterySensors)
	weakSignal := len(weakSignalSensors)
	errorSensors := len(unhealthySensors)

	// Calculate health percentage
	var healthyPercentage float64
	if totalSensors > 0 {
		healthySensors := totalSensors - errorSensors - criticalBattery
		if healthySensors < 0 {
			healthySensors = 0
		}
		healthyPercentage = (float64(healthySensors) / float64(totalSensors)) * 100
	}

	return &dto.SensorHealthSummaryResponse{
		TotalSensors:      totalSensors,
		OnlineSensors:     totalOnline,
		OfflineSensors:    totalOffline,
		LowBattery:        lowBattery,
		CriticalBattery:   criticalBattery,
		WeakSignal:        weakSignal,
		ErrorSensors:      errorSensors,
		HealthyPercentage: healthyPercentage,
	}, nil
}

// ListSensorStatuses retrieves a paginated list of sensor statuses with filtering
func (s *SensorStatusService) ListSensorStatuses(ctx context.Context, page, pageSize int, sensorUUID *uuid.UUID, status string) (*dto.SensorStatusListResponse, error) {
	// Set default pagination
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	// Build query parameters
	params := common.QueryParams{
		Page:     page,
		PageSize: pageSize,
	}

	// Create filter
	filter := dto.SensorStatusFilter{
		QueryParams: params,
	}

	if sensorUUID != nil {
		filter.AssetSensorID = sensorUUID
	}

	if status != "" {
		switch status {
		case "online":
			online := true
			filter.IsOnline = &online
		case "offline":
			offline := false
			filter.IsOnline = &offline
		case "low_battery":
			lowBattery := 20.0
			filter.BatteryLevelMax = &lowBattery
		case "weak_signal":
			weakSignal := "poor"
			filter.SignalStatus = &weakSignal
		case "unhealthy":
			hasErrors := true
			filter.HasErrors = &hasErrors
		}
	}

	// Get data from repository (this would need to be implemented in repository)
	// For now, using existing methods as fallback
	var statuses []*entity.SensorStatus
	var pagination *common.PaginationResponse
	var err error

	switch status {
	case "online":
		statuses, pagination, err = s.repo.GetOnlineSensors(ctx, params)
	case "offline":
		statuses, pagination, err = s.repo.GetOfflineSensors(ctx, params)
	case "low_battery":
		statuses, pagination, err = s.repo.GetLowBatterySensors(ctx, 20.0, params)
	case "weak_signal":
		statuses, pagination, err = s.repo.GetWeakSignalSensors(ctx, -70, params)
	case "unhealthy":
		statuses, pagination, err = s.repo.GetUnhealthySensors(ctx, params)
	default:
		// Get all statuses
		statuses, pagination, err = s.repo.GetAll(ctx, params)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to list sensor statuses: %v", err)
	}

	dtos := dto.FromSensorStatusEntityList(statuses)
	return &dto.SensorStatusListResponse{
		Data:       dtos,
		Pagination: *pagination,
		Message:    "Sensor statuses retrieved successfully",
	}, nil
}

// UpdateSensorStatusBySensorID updates a sensor status by asset sensor ID
func (s *SensorStatusService) UpdateSensorStatusBySensorID(ctx context.Context, assetSensorID uuid.UUID, req *dto.UpdateSensorStatusRequest) (*dto.SensorStatusDTO, error) {
	// Get existing sensor status by sensor ID
	existing, err := s.repo.GetBySensorID(ctx, assetSensorID)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing sensor status: %v", err)
	}
	if existing == nil {
		return nil, fmt.Errorf("sensor status not found for sensor ID")
	}

	// Create partial update entity from request
	updates := req.ToEntity(existing.ID)

	// Apply updates to existing entity
	s.applyUpdates(existing, updates)

	// Update in repository
	err = s.repo.Update(ctx, existing)
	if err != nil {
		return nil, fmt.Errorf("failed to update sensor status: %v", err)
	}

	return dto.FromSensorStatusEntity(existing), nil
}

// RecordHeartbeat updates the heartbeat timestamp for a sensor
func (s *SensorStatusService) RecordHeartbeat(ctx context.Context, assetSensorID uuid.UUID) (*dto.SensorStatusDTO, error) {
	// Try to get existing status first
	existing, err := s.repo.GetBySensorID(ctx, assetSensorID)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing sensor status: %v", err)
	}

	now := time.Now()

	if existing == nil {
		// Create a minimal status record if none exists
		status := &entity.SensorStatus{
			ID:               uuid.New(),
			AssetSensorID:    assetSensorID,
			ConnectionStatus: "online",
			IsOnline:         true,
			LastHeartbeat:    &now,
			RecordedAt:       now,
			CreatedAt:        now,
		}

		// Get tenant context
		tenantID, _, err := s.repo.GetAssetSensorContext(ctx, assetSensorID)
		if err != nil {
			return nil, fmt.Errorf("failed to get asset sensor context: %v", err)
		}
		status.TenantID = tenantID

		// Create in repository
		err = s.repo.Create(ctx, status)
		if err != nil {
			return nil, fmt.Errorf("failed to create sensor status for heartbeat: %v", err)
		}

		return dto.FromSensorStatusEntity(status), nil
	}

	// Update existing record
	existing.LastHeartbeat = &now
	existing.IsOnline = true
	if existing.ConnectionStatus == "offline" {
		existing.ConnectionStatus = "online"
		existing.LastConnectedAt = &now
	}
	existing.UpdatedAt = &now

	err = s.repo.Update(ctx, existing)
	if err != nil {
		return nil, fmt.Errorf("failed to update heartbeat: %v", err)
	}

	return dto.FromSensorStatusEntity(existing), nil
}

// GetHealthSummary retrieves aggregated health summary
func (s *SensorStatusService) GetHealthSummary(ctx context.Context) (*dto.SensorHealthSummaryResponse, error) {
	params := common.QueryParams{
		Page:     1,
		PageSize: 1000, // Get a large set for summary
	}

	return s.GetSensorHealthSummary(ctx, params)
}

// GetHealthAnalytics retrieves detailed health analytics over a timeframe
func (s *SensorStatusService) GetHealthAnalytics(ctx context.Context, timeframe string) (*dto.SensorHealthAnalyticsResponse, error) {
	// Validate timeframe
	var duration time.Duration
	var interval time.Duration

	switch timeframe {
	case "24h":
		duration = 24 * time.Hour
		interval = time.Hour
	case "7d":
		duration = 7 * 24 * time.Hour
		interval = 6 * time.Hour
	case "30d":
		duration = 30 * 24 * time.Hour
		interval = 24 * time.Hour
	default:
		return nil, fmt.Errorf("invalid timeframe: %s. Supported: 24h, 7d, 30d", timeframe)
	}

	now := time.Now()
	startTime := now.Add(-duration)

	// Generate data points for the timeframe
	// This is a simplified implementation - in a real system you'd query historical data
	var dataPoints []dto.SensorHealthDataPoint

	for t := startTime; t.Before(now); t = t.Add(interval) {
		// For demo purposes, we'll generate mock data based on current status
		// In a real implementation, you'd query historical sensor status data

		params := common.QueryParams{Page: 1, PageSize: 1000}
		currentSummary, err := s.GetSensorHealthSummary(ctx, params)
		if err != nil {
			return nil, fmt.Errorf("failed to get current summary for analytics: %v", err)
		}

		// Add some variation for demo purposes
		variation := float64(t.Unix()%10) / 10.0

		dataPoint := dto.SensorHealthDataPoint{
			Timestamp:       t,
			OnlineSensors:   int(float64(currentSummary.OnlineSensors) * (0.9 + variation*0.2)),
			OfflineSensors:  int(float64(currentSummary.OfflineSensors) * (0.9 + variation*0.2)),
			LowBattery:      int(float64(currentSummary.LowBattery) * (0.9 + variation*0.2)),
			CriticalBattery: int(float64(currentSummary.CriticalBattery) * (0.9 + variation*0.2)),
			WeakSignal:      int(float64(currentSummary.WeakSignal) * (0.9 + variation*0.2)),
			ErrorSensors:    int(float64(currentSummary.ErrorSensors) * (0.9 + variation*0.2)),
		}
		dataPoints = append(dataPoints, dataPoint)
	}

	// Calculate summary and trends
	totalSensors := 0
	if len(dataPoints) > 0 {
		latest := dataPoints[len(dataPoints)-1]
		totalSensors = latest.OnlineSensors + latest.OfflineSensors
	}

	// Calculate averages
	var totalOnlinePercentage, totalLowBatteryPercentage, totalWeakSignalPercentage float64
	for _, dp := range dataPoints {
		total := dp.OnlineSensors + dp.OfflineSensors
		if total > 0 {
			totalOnlinePercentage += float64(dp.OnlineSensors) / float64(total) * 100
			totalLowBatteryPercentage += float64(dp.LowBattery) / float64(total) * 100
			totalWeakSignalPercentage += float64(dp.WeakSignal) / float64(total) * 100
		}
	}

	pointCount := float64(len(dataPoints))
	summary := dto.SensorHealthSummary{
		AverageOnlinePercentage:     totalOnlinePercentage / pointCount,
		AverageLowBatteryPercentage: totalLowBatteryPercentage / pointCount,
		AverageWeakSignalPercentage: totalWeakSignalPercentage / pointCount,
		TotalDowntimeHours:          duration.Hours() * (1 - (totalOnlinePercentage/pointCount)/100),
		MostCommonIssue:             "low_battery", // Simplified
	}

	// Calculate trends (compare first vs last)
	trends := dto.SensorHealthTrends{}
	if len(dataPoints) >= 2 {
		first := dataPoints[0]
		last := dataPoints[len(dataPoints)-1]

		firstTotal := first.OnlineSensors + first.OfflineSensors
		lastTotal := last.OnlineSensors + last.OfflineSensors

		if firstTotal > 0 && lastTotal > 0 {
			firstOnlinePercentage := float64(first.OnlineSensors) / float64(firstTotal) * 100
			lastOnlinePercentage := float64(last.OnlineSensors) / float64(lastTotal) * 100
			trends.OnlinePercentageChange = lastOnlinePercentage - firstOnlinePercentage

			firstLowBatteryPercentage := float64(first.LowBattery) / float64(firstTotal) * 100
			lastLowBatteryPercentage := float64(last.LowBattery) / float64(lastTotal) * 100
			trends.LowBatteryPercentageChange = lastLowBatteryPercentage - firstLowBatteryPercentage

			firstWeakSignalPercentage := float64(first.WeakSignal) / float64(firstTotal) * 100
			lastWeakSignalPercentage := float64(last.WeakSignal) / float64(lastTotal) * 100
			trends.WeakSignalPercentageChange = lastWeakSignalPercentage - firstWeakSignalPercentage
		}

		trends.ErrorCountChange = last.ErrorSensors - first.ErrorSensors
	}

	return &dto.SensorHealthAnalyticsResponse{
		Timeframe:    timeframe,
		TotalSensors: totalSensors,
		Analytics:    dataPoints,
		Summary:      summary,
		Trends:       trends,
	}, nil
}

// applyUpdates applies partial updates to an existing entity
func (s *SensorStatusService) applyUpdates(existing *entity.SensorStatus, updates *entity.SensorStatus) {
	now := time.Now()
	existing.UpdatedAt = &now
	existing.RecordedAt = updates.RecordedAt

	// Apply battery updates
	if updates.BatteryLevel != nil {
		existing.BatteryLevel = updates.BatteryLevel
	}
	if updates.BatteryVoltage != nil {
		existing.BatteryVoltage = updates.BatteryVoltage
	}
	if updates.BatteryStatus != nil {
		existing.BatteryStatus = updates.BatteryStatus
	}
	if updates.BatteryLastCharged != nil {
		existing.BatteryLastCharged = updates.BatteryLastCharged
	}
	if updates.BatteryEstimatedLife != nil {
		existing.BatteryEstimatedLife = updates.BatteryEstimatedLife
	}
	if updates.BatteryType != nil {
		existing.BatteryType = updates.BatteryType
	}

	// Apply signal updates
	if updates.SignalType != nil {
		existing.SignalType = updates.SignalType
	}
	if updates.SignalRSSI != nil {
		existing.SignalRSSI = updates.SignalRSSI
	}
	if updates.SignalSNR != nil {
		existing.SignalSNR = updates.SignalSNR
	}
	if updates.SignalQuality != nil {
		existing.SignalQuality = updates.SignalQuality
	}
	if updates.SignalFrequency != nil {
		existing.SignalFrequency = updates.SignalFrequency
	}
	if updates.SignalChannel != nil {
		existing.SignalChannel = updates.SignalChannel
	}
	if updates.SignalStatus != nil {
		existing.SignalStatus = updates.SignalStatus
	}

	// Apply connection updates
	if updates.ConnectionType != nil {
		existing.ConnectionType = updates.ConnectionType
	}
	if updates.ConnectionStatus != "" {
		existing.ConnectionStatus = updates.ConnectionStatus
	}
	if updates.LastConnectedAt != nil {
		existing.LastConnectedAt = updates.LastConnectedAt
	}
	if updates.LastDisconnectedAt != nil {
		existing.LastDisconnectedAt = updates.LastDisconnectedAt
	}
	if updates.CurrentIP != nil {
		existing.CurrentIP = updates.CurrentIP
	}
	if updates.CurrentNetwork != nil {
		existing.CurrentNetwork = updates.CurrentNetwork
	}

	// Apply additional status updates
	if updates.Temperature != nil {
		existing.Temperature = updates.Temperature
	}
	if updates.Humidity != nil {
		existing.Humidity = updates.Humidity
	}
	// Handle boolean field carefully - default value is false
	existing.IsOnline = updates.IsOnline
	if updates.LastHeartbeat != nil {
		existing.LastHeartbeat = updates.LastHeartbeat
	}
	if updates.FirmwareVersion != nil {
		existing.FirmwareVersion = updates.FirmwareVersion
	}
	if updates.ErrorCount != nil {
		existing.ErrorCount = updates.ErrorCount
	}
	if updates.LastErrorAt != nil {
		existing.LastErrorAt = updates.LastErrorAt
	}
}
