package controller

import (
	"be-lecsens/asset_management/data-layer/config"
	"be-lecsens/asset_management/domain-layer/service"
	"be-lecsens/asset_management/helpers/common"
	"be-lecsens/asset_management/helpers/dto"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// IoTSensorReadingController handles HTTP requests for IoT sensor reading operations
type IoTSensorReadingController struct {
	iotSensorReadingService *service.IoTSensorReadingService
	config                  *config.Config
}

// NewIoTSensorReadingController creates a new IoTSensorReadingController
func NewIoTSensorReadingController(iotSensorReadingService *service.IoTSensorReadingService, cfg *config.Config) *IoTSensorReadingController {
	return &IoTSensorReadingController{
		iotSensorReadingService: iotSensorReadingService,
		config:                  cfg,
	}
}

// CreateReading handles the creation of a new IoT sensor reading
func (c *IoTSensorReadingController) CreateReading(ctx *gin.Context) {
	var request dto.CreateIoTSensorReadingRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Extract tenant context if available
	if tenantID, exists := common.GetTenantID(ctx.Request.Context()); exists {
		// Set tenant ID from context if not provided in request
		if request.TenantID == nil {
			request.TenantID = &tenantID
		}
	}

	reading, err := c.iotSensorReadingService.CreateReading(ctx.Request.Context(), &request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, reading)
}

// GetReading handles retrieving an IoT sensor reading by ID
func (c *IoTSensorReadingController) GetReading(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid reading ID"})
		return
	}

	reading, err := c.iotSensorReadingService.GetReading(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "reading not found"})
		return
	}

	ctx.JSON(http.StatusOK, reading)
}

// ListReadings handles retrieving a list of IoT sensor readings with pagination and filtering
func (c *IoTSensorReadingController) ListReadings(ctx *gin.Context) {
	var queryParams dto.IoTSensorReadingQueryParams

	// Bind query parameters
	if err := ctx.ShouldBindQuery(&queryParams); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set default pagination
	if queryParams.Page == 0 {
		queryParams.Page = 1
	}
	if queryParams.Limit == 0 {
		queryParams.Limit = 10
	}

	// Ensure limit is reasonable
	if queryParams.Limit < 1 || queryParams.Limit > 100 {
		queryParams.Limit = 10
	}

	// Extract tenant context if available and not specified in query
	if tenantID, exists := common.GetTenantID(ctx.Request.Context()); exists {
		if queryParams.TenantID == nil {
			queryParams.TenantID = &tenantID
		}
	}

	readings, err := c.iotSensorReadingService.ListReadings(ctx.Request.Context(), &queryParams)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list readings: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, readings)
}

// ListReadingsByAssetSensor handles retrieving IoT sensor readings for a specific asset sensor
func (c *IoTSensorReadingController) ListReadingsByAssetSensor(ctx *gin.Context) {
	assetSensorID, err := uuid.Parse(ctx.Param("asset_sensor_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid asset sensor ID"})
		return
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	// Ensure page is at least 1
	if page < 1 {
		page = 1
	}

	// Ensure limit is reasonable
	if limit < 1 || limit > 100 {
		limit = 10
	}

	readings, err := c.iotSensorReadingService.ListReadingsByAssetSensor(ctx.Request.Context(), assetSensorID, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list readings by asset sensor: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, readings)
}

// ListReadingsBySensorType handles retrieving IoT sensor readings for a specific sensor type
func (c *IoTSensorReadingController) ListReadingsBySensorType(ctx *gin.Context) {
	sensorTypeID, err := uuid.Parse(ctx.Param("sensor_type_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid sensor type ID"})
		return
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	// Ensure page is at least 1
	if page < 1 {
		page = 1
	}

	// Ensure limit is reasonable
	if limit < 1 || limit > 100 {
		limit = 10
	}

	readings, err := c.iotSensorReadingService.ListReadingsBySensorType(ctx.Request.Context(), sensorTypeID, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list readings by sensor type: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, readings)
}

// ListReadingsByMacAddress handles retrieving IoT sensor readings for a specific MAC address
func (c *IoTSensorReadingController) ListReadingsByMacAddress(ctx *gin.Context) {
	macAddress := ctx.Param("mac_address")
	if macAddress == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "MAC address is required"})
		return
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	// Ensure page is at least 1
	if page < 1 {
		page = 1
	}

	// Ensure limit is reasonable
	if limit < 1 || limit > 100 {
		limit = 10
	}

	readings, err := c.iotSensorReadingService.ListReadingsByMacAddress(ctx.Request.Context(), macAddress, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list readings by MAC address: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, readings)
}

// ListReadingsByTimeRange handles retrieving IoT sensor readings within a time range
func (c *IoTSensorReadingController) ListReadingsByTimeRange(ctx *gin.Context) {
	startTimeStr := ctx.Query("start_time")
	endTimeStr := ctx.Query("end_time")

	if startTimeStr == "" || endTimeStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "start_time and end_time are required"})
		return
	}

	startTime, err := time.Parse(time.RFC3339, startTimeStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_time format, use RFC3339 (e.g., 2006-01-02T15:04:05Z07:00)"})
		return
	}

	endTime, err := time.Parse(time.RFC3339, endTimeStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_time format, use RFC3339 (e.g., 2006-01-02T15:04:05Z07:00)"})
		return
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	// Ensure page is at least 1
	if page < 1 {
		page = 1
	}

	// Ensure limit is reasonable
	if limit < 1 || limit > 100 {
		limit = 10
	}

	readings, err := c.iotSensorReadingService.ListReadingsByTimeRange(ctx.Request.Context(), startTime, endTime, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list readings by time range: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, readings)
}

// UpdateReading handles updating an existing IoT sensor reading (partial update)
func (c *IoTSensorReadingController) UpdateReading(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid reading ID"})
		return
	}

	var request dto.UpdateIoTSensorReadingRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reading, err := c.iotSensorReadingService.UpdateReading(ctx.Request.Context(), id, &request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, reading)
}

// DeleteReading handles deleting an IoT sensor reading
func (c *IoTSensorReadingController) DeleteReading(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid reading ID"})
		return
	}

	if err := c.iotSensorReadingService.DeleteReading(ctx.Request.Context(), id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}

// GetLatestByMacAddress handles retrieving the latest IoT sensor reading for a MAC address
func (c *IoTSensorReadingController) GetLatestByMacAddress(ctx *gin.Context) {
	macAddress := ctx.Param("mac_address")
	if macAddress == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "MAC address is required"})
		return
	}

	reading, err := c.iotSensorReadingService.GetLatestByMacAddress(ctx.Request.Context(), macAddress)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "reading not found"})
		return
	}

	ctx.JSON(http.StatusOK, reading)
}

// ListAllReadings handles retrieving all readings for SuperAdmin (across all tenants)
func (c *IoTSensorReadingController) ListAllReadings(ctx *gin.Context) {
	// Check if user has SuperAdmin role
	userRole, exists := ctx.Get("user_role")
	if !exists || userRole != "SUPERADMIN" {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "SuperAdmin access required"})
		return
	}

	var queryParams dto.IoTSensorReadingQueryParams

	// Bind query parameters
	if err := ctx.ShouldBindQuery(&queryParams); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set default pagination
	if queryParams.Page == 0 {
		queryParams.Page = 1
	}
	if queryParams.Limit == 0 {
		queryParams.Limit = 10
	}

	// Ensure limit is reasonable
	if queryParams.Limit < 1 || queryParams.Limit > 100 {
		queryParams.Limit = 10
	}

	// For SuperAdmin, don't filter by tenant (leave TenantID as nil to get all)
	queryParams.TenantID = nil

	readings, err := c.iotSensorReadingService.ListReadings(ctx.Request.Context(), &queryParams)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list all readings: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, readings)
}
