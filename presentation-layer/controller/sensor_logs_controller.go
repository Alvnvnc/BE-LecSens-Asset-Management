package controller

import (
	"be-lecsens/asset_management/domain-layer/service"
	"be-lecsens/asset_management/helpers/common"
	"be-lecsens/asset_management/helpers/dto"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// SensorLogsController handles HTTP requests for sensor logs operations
type SensorLogsController struct {
	sensorLogsService *service.SensorLogsService
}

// NewSensorLogsController creates a new SensorLogsController
func NewSensorLogsController(sensorLogsService *service.SensorLogsService) *SensorLogsController {
	return &SensorLogsController{
		sensorLogsService: sensorLogsService,
	}
}

// CreateSensorLog handles POST /api/v1/sensor-logs
func (c *SensorLogsController) CreateSensorLog(ctx *gin.Context) {
	var req dto.CreateSensorLogsRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error binding request: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}

	log.Printf("Creating sensor log with request: %+v", req)

	sensorLog, err := c.sensorLogsService.CreateSensorLog(ctx, req)
	if err != nil {
		log.Printf("Error creating sensor log: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create sensor log",
			"details": err.Error(),
		})
		return
	}

	log.Printf("Successfully created sensor log: %+v", sensorLog)

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Sensor log created successfully",
		"data":    sensorLog,
	})
}

// GetSensorLog handles GET /api/v1/sensor-logs/:id
func (c *SensorLogsController) GetSensorLog(ctx *gin.Context) {
	idParam := ctx.Param("id")
	log.Printf("Getting sensor log with ID: %s", idParam)

	id, err := uuid.Parse(idParam)
	if err != nil {
		log.Printf("Invalid UUID format: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid UUID format",
			"details": err.Error(),
		})
		return
	}

	sensorLog, err := c.sensorLogsService.GetSensorLog(ctx, id)
	if err != nil {
		log.Printf("Error getting sensor log: %v", err)
		ctx.JSON(http.StatusNotFound, gin.H{
			"error":   "Sensor log not found",
			"details": err.Error(),
		})
		return
	}

	log.Printf("Successfully retrieved sensor log: %+v", sensorLog)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Sensor log retrieved successfully",
		"data":    sensorLog,
	})
}

// GetSensorLogsBySensorID handles GET /api/v1/sensor-logs/sensor/:sensor_id
func (c *SensorLogsController) GetSensorLogsBySensorID(ctx *gin.Context) {
	sensorIDParam := ctx.Param("sensor_id")
	log.Printf("Getting sensor logs for sensor ID: %s", sensorIDParam)

	sensorID, err := uuid.Parse(sensorIDParam)
	if err != nil {
		log.Printf("Invalid sensor UUID format: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid sensor UUID format",
			"details": err.Error(),
		})
		return
	}

	// Parse pagination parameters
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))

	filter := dto.SensorLogsFilter{
		AssetSensorID: &sensorID,
		QueryParams: common.QueryParams{
			Page:     page,
			PageSize: pageSize,
		},
	}

	logs, pagination, err := c.sensorLogsService.ListSensorLogs(ctx, filter)
	if err != nil {
		log.Printf("Error getting sensor logs: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get sensor logs",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":    "Sensor logs retrieved successfully",
		"data":       logs,
		"pagination": pagination,
	})
}

// ListSensorLogs handles GET /api/v1/sensor-logs
func (c *SensorLogsController) ListSensorLogs(ctx *gin.Context) {
	// Parse pagination parameters
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))

	// Parse filter parameters
	filter := dto.SensorLogsFilter{
		QueryParams: common.QueryParams{
			Page:     page,
			PageSize: pageSize,
		},
	}

	// Asset sensor ID filter
	if sensorIDStr := ctx.Query("asset_sensor_id"); sensorIDStr != "" {
		if sensorID, err := uuid.Parse(sensorIDStr); err == nil {
			filter.AssetSensorID = &sensorID
		}
	}

	// Log type filter
	if logType := ctx.Query("log_type"); logType != "" {
		filter.LogType = &logType
	}

	// Log level filter
	if logLevel := ctx.Query("log_level"); logLevel != "" {
		filter.LogLevel = &logLevel
	}

	// Component filter
	if component := ctx.Query("component"); component != "" {
		filter.Component = &component
	}

	// Event type filter
	if eventType := ctx.Query("event_type"); eventType != "" {
		filter.EventType = &eventType
	}

	// Connection status filter
	if connectionStatus := ctx.Query("connection_status"); connectionStatus != "" {
		filter.ConnectionStatus = &connectionStatus
	}

	// Search filter
	if search := ctx.Query("search"); search != "" {
		filter.SearchMessage = &search
	}

	log.Printf("Listing sensor logs with filter: %+v", filter)

	logs, pagination, err := c.sensorLogsService.ListSensorLogs(ctx, filter)
	if err != nil {
		log.Printf("Error listing sensor logs: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to list sensor logs",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":    "Sensor logs retrieved successfully",
		"data":       logs,
		"pagination": pagination,
	})
}

// UpdateSensorLog handles PUT /api/v1/sensor-logs/:id
func (c *SensorLogsController) UpdateSensorLog(ctx *gin.Context) {
	idParam := ctx.Param("id")
	log.Printf("Updating sensor log with ID: %s", idParam)

	id, err := uuid.Parse(idParam)
	if err != nil {
		log.Printf("Invalid UUID format: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid UUID format",
			"details": err.Error(),
		})
		return
	}

	var req dto.UpdateSensorLogsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error binding request: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}

	sensorLog, err := c.sensorLogsService.UpdateSensorLog(ctx, id, req)
	if err != nil {
		log.Printf("Error updating sensor log: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update sensor log",
			"details": err.Error(),
		})
		return
	}

	log.Printf("Successfully updated sensor log: %+v", sensorLog)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Sensor log updated successfully",
		"data":    sensorLog,
	})
}

// DeleteSensorLog handles DELETE /api/v1/sensor-logs/:id
func (c *SensorLogsController) DeleteSensorLog(ctx *gin.Context) {
	idParam := ctx.Param("id")
	log.Printf("Deleting sensor log with ID: %s", idParam)

	id, err := uuid.Parse(idParam)
	if err != nil {
		log.Printf("Invalid UUID format: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid UUID format",
			"details": err.Error(),
		})
		return
	}

	err = c.sensorLogsService.DeleteSensorLog(ctx, id)
	if err != nil {
		log.Printf("Error deleting sensor log: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete sensor log",
			"details": err.Error(),
		})
		return
	}

	log.Printf("Successfully deleted sensor log with ID: %s", id)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Sensor log deleted successfully",
	})
}

// DeleteSensorLogsBySensorID handles DELETE /api/v1/sensor-logs/sensor/:sensor_id
func (c *SensorLogsController) DeleteSensorLogsBySensorID(ctx *gin.Context) {
	sensorIDParam := ctx.Param("sensor_id")
	log.Printf("Deleting sensor logs for sensor ID: %s", sensorIDParam)

	sensorID, err := uuid.Parse(sensorIDParam)
	if err != nil {
		log.Printf("Invalid sensor UUID format: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid sensor UUID format",
			"details": err.Error(),
		})
		return
	}

	deletedCount, err := c.sensorLogsService.DeleteSensorLogsBySensorID(ctx, sensorID)
	if err != nil {
		log.Printf("Error deleting sensor logs: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete sensor logs",
			"details": err.Error(),
		})
		return
	}

	log.Printf("Successfully deleted %d sensor logs for sensor ID: %s", deletedCount, sensorID)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Sensor logs deleted successfully",
		"data": gin.H{
			"deleted_count": deletedCount,
		},
	})
}

// GetConnectionHistory handles GET /api/v1/sensor-logs/sensor/:sensor_id/connection-history
func (c *SensorLogsController) GetConnectionHistory(ctx *gin.Context) {
	sensorIDParam := ctx.Param("sensor_id")
	log.Printf("Getting connection history for sensor ID: %s", sensorIDParam)

	sensorID, err := uuid.Parse(sensorIDParam)
	if err != nil {
		log.Printf("Invalid sensor UUID format: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid sensor UUID format",
			"details": err.Error(),
		})
		return
	}

	// Parse pagination parameters
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))

	params := common.QueryParams{
		Page:     page,
		PageSize: pageSize,
	}

	history, pagination, err := c.sensorLogsService.GetConnectionHistory(ctx, sensorID, params)
	if err != nil {
		log.Printf("Error getting connection history: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get connection history",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":    "Connection history retrieved successfully",
		"data":       history,
		"pagination": pagination,
	})
}

// GetLogStatistics handles GET /api/v1/sensor-logs/statistics
func (c *SensorLogsController) GetLogStatistics(ctx *gin.Context) {
	// Parse filter parameters
	var assetSensorID *uuid.UUID

	// Asset sensor ID filter
	if sensorIDStr := ctx.Query("asset_sensor_id"); sensorIDStr != "" {
		if sensorID, err := uuid.Parse(sensorIDStr); err == nil {
			assetSensorID = &sensorID
		}
	}

	// Parse date range (default to last 30 days)
	now := time.Now()
	startDate := now.AddDate(0, 0, -30) // 30 days ago
	endDate := now

	if startDateStr := ctx.Query("start_date"); startDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", startDateStr); err == nil {
			startDate = parsed
		}
	}

	if endDateStr := ctx.Query("end_date"); endDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", endDateStr); err == nil {
			endDate = parsed
		}
	}

	log.Printf("Getting log statistics for sensor: %v, from: %v to: %v", assetSensorID, startDate, endDate)

	statistics, err := c.sensorLogsService.GetLogStatistics(ctx, assetSensorID, startDate, endDate)
	if err != nil {
		log.Printf("Error getting log statistics: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get log statistics",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Log statistics retrieved successfully",
		"data":    statistics,
	})
}

// GetLogAnalytics handles GET /api/v1/sensor-logs/analytics
func (c *SensorLogsController) GetLogAnalytics(ctx *gin.Context) {
	var req dto.LogAnalyticsRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		log.Printf("Error binding analytics request: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request parameters",
			"details": err.Error(),
		})
		return
	}

	log.Printf("Getting log analytics with request: %+v", req)

	analytics, err := c.sensorLogsService.GetLogAnalytics(ctx, req)
	if err != nil {
		log.Printf("Error getting log analytics: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get log analytics",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Log analytics retrieved successfully",
		"data":    analytics,
	})
}

// CleanupOldLogs handles DELETE /api/v1/sensor-logs/cleanup
func (c *SensorLogsController) CleanupOldLogs(ctx *gin.Context) {
	// Parse days parameter
	daysStr := ctx.DefaultQuery("days", "30")
	days, err := strconv.Atoi(daysStr)
	if err != nil || days <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid days parameter",
			"details": "Days must be a positive integer",
		})
		return
	}

	log.Printf("Cleaning up logs older than %d days", days)

	deletedCount, err := c.sensorLogsService.CleanupOldLogs(ctx, days)
	if err != nil {
		log.Printf("Error cleaning up old logs: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to cleanup old logs",
			"details": err.Error(),
		})
		return
	}

	log.Printf("Successfully cleaned up %d old logs", deletedCount)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Old logs cleaned up successfully",
		"data": gin.H{
			"deleted_count": deletedCount,
			"days":          days,
		},
	})
}
