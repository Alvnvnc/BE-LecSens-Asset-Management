package service

import (
	"be-lecsens/asset_management/data-layer/entity"
	"be-lecsens/asset_management/data-layer/repository"
	"be-lecsens/asset_management/helpers/common"
	"be-lecsens/asset_management/helpers/dto"
	"context"
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
)

// SensorLogsService handles business logic for sensor logs operations
type SensorLogsService struct {
	repo repository.SensorLogsRepository
}

// NewSensorLogsService creates a new instance of SensorLogsService
func NewSensorLogsService(repo repository.SensorLogsRepository) *SensorLogsService {
	return &SensorLogsService{
		repo: repo,
	}
}

// CreateSensorLog creates a new sensor log entry
func (s *SensorLogsService) CreateSensorLog(ctx context.Context, req dto.CreateSensorLogsRequest) (*dto.SensorLogsDTO, error) {
	// Get asset sensor context for tenant inheritance
	tenantID, _, err := s.repo.GetAssetSensorContext(ctx, req.AssetSensorID)
	if err != nil {
		return nil, fmt.Errorf("failed to get asset sensor context: %v", err)
	}

	// Convert request to entity
	log := req.ToEntity(tenantID)

	// Create in repository
	err = s.repo.Create(ctx, log)
	if err != nil {
		return nil, fmt.Errorf("failed to create sensor log: %v", err)
	}

	// Convert to DTO and return
	return dto.FromSensorLogsEntity(log), nil
}

// GetSensorLog retrieves a sensor log by ID
func (s *SensorLogsService) GetSensorLog(ctx context.Context, id uuid.UUID) (*dto.SensorLogsDTO, error) {
	log, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get sensor log: %v", err)
	}
	if log == nil {
		return nil, fmt.Errorf("sensor log not found")
	}

	return dto.FromSensorLogsEntity(log), nil
}

// GetSensorLogsBySensorID retrieves logs for a specific sensor with pagination
func (s *SensorLogsService) GetSensorLogsBySensorID(ctx context.Context, assetSensorID uuid.UUID, params common.QueryParams) ([]dto.SensorLogsDTO, *common.PaginationResponse, error) {
	logs, pagination, err := s.repo.GetBySensorID(ctx, assetSensorID, params)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get sensor logs by sensor ID: %v", err)
	}

	dtos := dto.FromSensorLogsEntityList(logs)
	return dtos, pagination, nil
}

// GetSensorLogsByType retrieves logs filtered by log type with pagination
func (s *SensorLogsService) GetSensorLogsByType(ctx context.Context, logType string, params common.QueryParams) ([]dto.SensorLogsDTO, *common.PaginationResponse, error) {
	logs, pagination, err := s.repo.GetByLogType(ctx, logType, params)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get sensor logs by type: %v", err)
	}

	dtos := dto.FromSensorLogsEntityList(logs)
	return dtos, pagination, nil
}

// GetSensorLogsByLevel retrieves logs filtered by log level with pagination
func (s *SensorLogsService) GetSensorLogsByLevel(ctx context.Context, logLevel string, params common.QueryParams) ([]dto.SensorLogsDTO, *common.PaginationResponse, error) {
	logs, pagination, err := s.repo.GetByLogLevel(ctx, logLevel, params)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get sensor logs by level: %v", err)
	}

	dtos := dto.FromSensorLogsEntityList(logs)
	return dtos, pagination, nil
}

// GetConnectionHistory retrieves connection history logs for a sensor
func (s *SensorLogsService) GetConnectionHistory(ctx context.Context, assetSensorID uuid.UUID, params common.QueryParams) ([]dto.SensorLogsDTO, *common.PaginationResponse, error) {
	logs, pagination, err := s.repo.GetConnectionHistory(ctx, assetSensorID, params)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get connection history: %v", err)
	}

	dtos := dto.FromSensorLogsEntityList(logs)
	return dtos, pagination, nil
}

// SearchSensorLogs searches logs by message content with pagination
func (s *SensorLogsService) SearchSensorLogs(ctx context.Context, searchQuery string, params common.QueryParams) ([]dto.SensorLogsDTO, *common.PaginationResponse, error) {
	logs, pagination, err := s.repo.SearchLogs(ctx, searchQuery, params)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to search sensor logs: %v", err)
	}

	dtos := dto.FromSensorLogsEntityList(logs)
	return dtos, pagination, nil
}

// GetErrorLogs retrieves error and critical level logs with pagination
func (s *SensorLogsService) GetErrorLogs(ctx context.Context, assetSensorID *uuid.UUID, params common.QueryParams) ([]dto.SensorLogsDTO, *common.PaginationResponse, error) {
	logs, pagination, err := s.repo.GetErrorLogs(ctx, assetSensorID, params)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get error logs: %v", err)
	}

	dtos := dto.FromSensorLogsEntityList(logs)
	return dtos, pagination, nil
}

// UpdateSensorLog updates an existing sensor log entry
func (s *SensorLogsService) UpdateSensorLog(ctx context.Context, id uuid.UUID, req dto.UpdateSensorLogsRequest) (*dto.SensorLogsDTO, error) {
	// Get existing sensor log
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing sensor log: %v", err)
	}
	if existing == nil {
		return nil, fmt.Errorf("sensor log not found")
	}

	// Create partial update entity from request
	updates := req.ToEntity(id)

	// Apply updates to existing entity
	s.applyUpdates(existing, updates)

	// Update in repository
	err = s.repo.Update(ctx, existing)
	if err != nil {
		return nil, fmt.Errorf("failed to update sensor log: %v", err)
	}

	return dto.FromSensorLogsEntity(existing), nil
}

// DeleteSensorLog deletes a sensor log entry
func (s *SensorLogsService) DeleteSensorLog(ctx context.Context, id uuid.UUID) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete sensor log: %v", err)
	}
	return nil
}

// DeleteOldLogs deletes logs older than the specified date
func (s *SensorLogsService) DeleteOldLogs(ctx context.Context, olderThan time.Time) (int64, error) {
	deletedCount, err := s.repo.DeleteOldLogs(ctx, olderThan)
	if err != nil {
		return 0, fmt.Errorf("failed to delete old logs: %v", err)
	}
	return deletedCount, nil
}

// ListSensorLogs retrieves sensor logs with filtering and pagination
func (s *SensorLogsService) ListSensorLogs(ctx context.Context, filter dto.SensorLogsFilter) ([]dto.SensorLogsDTO, *common.PaginationResponse, error) {
	// Handle different filter conditions and delegate to appropriate repository methods
	if filter.AssetSensorID != nil {
		// Get logs for specific sensor
		logs, pagination, err := s.repo.GetBySensorID(ctx, *filter.AssetSensorID, filter.QueryParams)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get sensor logs by sensor ID: %v", err)
		}
		dtos := dto.FromSensorLogsEntityList(logs)
		return dtos, pagination, nil
	}

	if filter.LogType != nil {
		// Get logs by type
		logs, pagination, err := s.repo.GetByLogType(ctx, *filter.LogType, filter.QueryParams)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get sensor logs by type: %v", err)
		}
		dtos := dto.FromSensorLogsEntityList(logs)
		return dtos, pagination, nil
	}

	if filter.LogLevel != nil {
		// Get logs by level
		logs, pagination, err := s.repo.GetByLogLevel(ctx, *filter.LogLevel, filter.QueryParams)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get sensor logs by level: %v", err)
		}
		dtos := dto.FromSensorLogsEntityList(logs)
		return dtos, pagination, nil
	}

	if filter.SearchMessage != nil {
		// Search logs
		logs, pagination, err := s.repo.SearchLogs(ctx, *filter.SearchMessage, filter.QueryParams)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to search sensor logs: %v", err)
		}
		dtos := dto.FromSensorLogsEntityList(logs)
		return dtos, pagination, nil
	}

	// Default: get error logs (as a fallback since we don't have GetAll)
	logs, pagination, err := s.repo.GetErrorLogs(ctx, nil, filter.QueryParams)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get sensor logs: %v", err)
	}
	dtos := dto.FromSensorLogsEntityList(logs)
	return dtos, pagination, nil
}

// DeleteSensorLogsBySensorID deletes all logs for a specific sensor
func (s *SensorLogsService) DeleteSensorLogsBySensorID(ctx context.Context, assetSensorID uuid.UUID) (int64, error) {
	// Since there's no direct DeleteBySensorID method, we need to implement it differently
	// For now, we'll return an error indicating this functionality needs to be implemented
	return 0, fmt.Errorf("DeleteSensorLogsBySensorID not yet implemented in repository")
}

// CleanupOldLogs removes logs older than the specified duration
func (s *SensorLogsService) CleanupOldLogs(ctx context.Context, olderThanDays int) (int64, error) {
	if olderThanDays <= 0 {
		return 0, fmt.Errorf("olderThanDays must be positive")
	}

	cutoffTime := time.Now().AddDate(0, 0, -olderThanDays)
	deletedCount, err := s.repo.DeleteOldLogs(ctx, cutoffTime)
	if err != nil {
		return 0, fmt.Errorf("failed to cleanup old logs: %v", err)
	}
	return deletedCount, nil
}

// GetLogStatistics provides comprehensive statistics about sensor logs
func (s *SensorLogsService) GetLogStatistics(ctx context.Context, assetSensorID *uuid.UUID, startDate, endDate time.Time) (*dto.LogStatisticsResponse, error) {
	// Set up query parameters for the time range
	params := common.QueryParams{
		Page:     1,
		PageSize: 1000, // Large limit to get comprehensive statistics
	}

	// Get all logs in the time range - we'll filter in memory for now
	// In a production system, you'd want to add time range filtering to the repository methods
	var allLogs []*entity.SensorLogs
	var err error

	if assetSensorID != nil {
		logs, _, err := s.repo.GetBySensorID(ctx, *assetSensorID, params)
		if err != nil {
			return nil, fmt.Errorf("failed to get logs for statistics: %v", err)
		}
		allLogs = logs
	} else {
		// For system-wide statistics, we'd need a GetAll method
		// For now, return error for system-wide stats
		return nil, fmt.Errorf("system-wide statistics not implemented yet")
	}

	// Filter logs by time range and calculate statistics
	filteredLogs := s.filterLogsByTimeRange(allLogs, startDate, endDate)

	// Calculate statistics
	stats := s.calculateLogStatistics(filteredLogs)

	// Get recent errors (last 10)
	errorParams := common.QueryParams{Page: 1, PageSize: 10}
	recentErrors, _, err := s.repo.GetErrorLogs(ctx, assetSensorID, errorParams)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent errors: %v", err)
	}

	stats.RecentErrors = dto.FromSensorLogsEntityList(recentErrors)
	stats.TimeRange = dto.LogTimeRangeStats{
		StartDate:    startDate,
		EndDate:      endDate,
		DurationDays: int(endDate.Sub(startDate).Hours() / 24),
		LogsPerDay:   s.calculateLogsPerDay(filteredLogs, startDate, endDate),
	}

	return stats, nil
}

// GetConnectionHistoryAnalysis provides detailed connection history analysis
func (s *SensorLogsService) GetConnectionHistoryAnalysis(ctx context.Context, assetSensorID uuid.UUID, startDate, endDate time.Time) (*dto.ConnectionHistoryResponse, error) {
	params := common.QueryParams{
		Page:     1,
		PageSize: 1000, // Large limit to get comprehensive data
	}

	// Get connection history
	logs, _, err := s.repo.GetConnectionHistory(ctx, assetSensorID, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get connection history: %v", err)
	}

	// Filter by time range
	filteredLogs := s.filterLogsByTimeRange(logs, startDate, endDate)

	// Analyze connection patterns
	analysis := s.analyzeConnectionHistory(assetSensorID, filteredLogs, startDate, endDate)

	return analysis, nil
}

// GetLogAnalytics provides comprehensive log analytics
func (s *SensorLogsService) GetLogAnalytics(ctx context.Context, req dto.LogAnalyticsRequest) (*dto.LogAnalyticsResponse, error) {
	// Get statistics
	stats, err := s.GetLogStatistics(ctx, req.AssetSensorID, req.StartDate, req.EndDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get log statistics: %v", err)
	}

	// Get connection history if sensor ID is provided
	var connectionHistory dto.ConnectionHistoryResponse
	if req.AssetSensorID != nil {
		history, err := s.GetConnectionHistoryAnalysis(ctx, *req.AssetSensorID, req.StartDate, req.EndDate)
		if err != nil {
			return nil, fmt.Errorf("failed to get connection history: %v", err)
		}
		connectionHistory = *history
	}

	// Generate trend analysis based on groupBy parameter
	trendAnalysis := s.generateTrendAnalysis(stats.TimeRange.LogsPerDay)

	// Detect anomalies
	anomalies := s.detectLogAnomalies(stats.TimeRange.LogsPerDay)

	return &dto.LogAnalyticsResponse{
		Statistics:        *stats,
		ConnectionHistory: connectionHistory,
		TrendAnalysis:     trendAnalysis,
		Anomalies:         anomalies,
	}, nil
}

// Utility methods for internal processing

func (s *SensorLogsService) applyUpdates(existing *entity.SensorLogs, updates *entity.SensorLogs) {
	now := time.Now()
	existing.UpdatedAt = &now
	existing.RecordedAt = updates.RecordedAt

	// Apply updates only for non-empty fields
	if updates.LogType != "" {
		existing.LogType = updates.LogType
	}
	if updates.LogLevel != "" {
		existing.LogLevel = updates.LogLevel
	}
	if updates.Message != "" {
		existing.Message = updates.Message
	}
	if updates.Component != nil {
		existing.Component = updates.Component
	}
	if updates.EventType != nil {
		existing.EventType = updates.EventType
	}
	if updates.ErrorCode != nil {
		existing.ErrorCode = updates.ErrorCode
	}
	if updates.ConnectionType != nil {
		existing.ConnectionType = updates.ConnectionType
	}
	if updates.ConnectionStatus != nil {
		existing.ConnectionStatus = updates.ConnectionStatus
	}
	if updates.IPAddress != nil {
		existing.IPAddress = updates.IPAddress
	}
	if updates.MACAddress != nil {
		existing.MACAddress = updates.MACAddress
	}
	if updates.NetworkName != nil {
		existing.NetworkName = updates.NetworkName
	}
	if updates.ConnectionDuration != nil {
		existing.ConnectionDuration = updates.ConnectionDuration
	}
	if updates.Metadata != nil {
		existing.Metadata = updates.Metadata
	}
	if updates.SourceIP != nil {
		existing.SourceIP = updates.SourceIP
	}
	if updates.UserAgent != nil {
		existing.UserAgent = updates.UserAgent
	}
	if updates.SessionID != nil {
		existing.SessionID = updates.SessionID
	}
}

func (s *SensorLogsService) filterLogsByTimeRange(logs []*entity.SensorLogs, startDate, endDate time.Time) []*entity.SensorLogs {
	var filtered []*entity.SensorLogs
	for _, log := range logs {
		if log.RecordedAt.After(startDate) && log.RecordedAt.Before(endDate) {
			filtered = append(filtered, log)
		}
	}
	return filtered
}

func (s *SensorLogsService) calculateLogStatistics(logs []*entity.SensorLogs) *dto.LogStatisticsResponse {
	stats := &dto.LogStatisticsResponse{
		TotalLogs:       int64(len(logs)),
		LogsByType:      make(map[string]int64),
		LogsByLevel:     make(map[string]int64),
		LogsByComponent: make(map[string]int64),
	}

	errorCount := int64(0)
	for _, log := range logs {
		// Count by type
		stats.LogsByType[log.LogType]++

		// Count by level
		stats.LogsByLevel[log.LogLevel]++

		// Count by component
		if log.Component != nil {
			stats.LogsByComponent[*log.Component]++
		}

		// Count errors
		if log.LogLevel == "error" || log.LogLevel == "critical" {
			errorCount++
		}
	}

	// Calculate error rate
	if stats.TotalLogs > 0 {
		stats.ErrorRate = (float64(errorCount) / float64(stats.TotalLogs)) * 100
	}

	return stats
}

func (s *SensorLogsService) calculateLogsPerDay(logs []*entity.SensorLogs, startDate, endDate time.Time) []dto.TimeSeriesData {
	// Create a map to count logs per day
	dailyCounts := make(map[string]int64)

	// Initialize all days in range with zero counts
	for d := startDate; d.Before(endDate); d = d.AddDate(0, 0, 1) {
		dateKey := d.Format("2006-01-02")
		dailyCounts[dateKey] = 0
	}

	// Count actual logs
	for _, log := range logs {
		dateKey := log.RecordedAt.Format("2006-01-02")
		dailyCounts[dateKey]++
	}

	// Convert to time series data
	var timeSeries []dto.TimeSeriesData
	for d := startDate; d.Before(endDate); d = d.AddDate(0, 0, 1) {
		dateKey := d.Format("2006-01-02")
		timeSeries = append(timeSeries, dto.TimeSeriesData{
			Date:  d,
			Count: dailyCounts[dateKey],
			Value: float64(dailyCounts[dateKey]),
		})
	}

	return timeSeries
}

func (s *SensorLogsService) analyzeConnectionHistory(assetSensorID uuid.UUID, logs []*entity.SensorLogs, startDate, endDate time.Time) *dto.ConnectionHistoryResponse {
	analysis := &dto.ConnectionHistoryResponse{
		AssetSensorID:     assetSensorID,
		ConnectionsByType: make(map[string]int64),
		TimeRange: dto.LogTimeRangeStats{
			StartDate:    startDate,
			EndDate:      endDate,
			DurationDays: int(endDate.Sub(startDate).Hours() / 24),
		},
	}

	var connections, disconnections int64
	var totalUptime float64

	for _, log := range logs {
		if log.IsConnectionLog() {
			if log.ConnectionType != nil {
				analysis.ConnectionsByType[*log.ConnectionType]++
			}

			if log.ConnectionStatus != nil {
				switch *log.ConnectionStatus {
				case "connected":
					connections++
				case "disconnected":
					disconnections++
				}
			}

			// Calculate uptime if duration is available
			if log.ConnectionDuration != nil {
				totalUptime += float64(*log.ConnectionDuration) / 3600.0 // Convert to hours
			}
		}
	}

	analysis.TotalConnections = connections
	analysis.TotalDisconnections = disconnections

	if connections > 0 {
		analysis.AverageUptime = totalUptime / float64(connections)
	}

	// Calculate uptime percentage
	totalHours := endDate.Sub(startDate).Hours()
	if totalHours > 0 {
		analysis.UptimePercentage = (totalUptime / totalHours) * 100
		if analysis.UptimePercentage > 100 {
			analysis.UptimePercentage = 100
		}
	}

	// Get recent activity (last 10 connection-related logs)
	recentLogs := logs
	if len(logs) > 10 {
		recentLogs = logs[len(logs)-10:]
	}
	analysis.RecentActivity = dto.FromSensorLogsEntityList(recentLogs)

	return analysis
}

func (s *SensorLogsService) generateTrendAnalysis(timeSeries []dto.TimeSeriesData) []dto.TimeSeriesData {
	// For now, return the daily data as-is
	// In a more sophisticated implementation, you would group by week/month as requested
	return timeSeries
}

func (s *SensorLogsService) detectLogAnomalies(timeSeries []dto.TimeSeriesData) []dto.LogAnomalyData {
	var anomalies []dto.LogAnomalyData

	if len(timeSeries) < 7 {
		return anomalies // Need at least a week of data
	}

	// Calculate average and standard deviation
	var sum, sumSquares float64
	for _, point := range timeSeries {
		sum += point.Value
		sumSquares += point.Value * point.Value
	}

	mean := sum / float64(len(timeSeries))
	variance := (sumSquares / float64(len(timeSeries))) - (mean * mean)
	stdDev := math.Sqrt(variance)

	// Detect anomalies (values more than 2 standard deviations from mean)
	threshold := 2.0
	for _, point := range timeSeries {
		deviation := math.Abs(point.Value - mean)
		if deviation > threshold*stdDev && point.Count > 0 {
			severity := "warning"
			if deviation > 3*stdDev {
				severity = "critical"
			}

			anomalies = append(anomalies, dto.LogAnomalyData{
				Date:        point.Date,
				Type:        "volume_spike",
				Description: fmt.Sprintf("Unusual log volume: %.0f logs (%.1f standard deviations from normal)", point.Value, deviation/stdDev),
				Severity:    severity,
				Count:       point.Count,
			})
		}
	}

	return anomalies
}
