package controller

import (
	"be-lecsens/asset_management/domain-layer/service"
	"be-lecsens/asset_management/helpers/common"
	"be-lecsens/asset_management/helpers/dto"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// SensorThresholdController handles HTTP requests for sensor threshold operations
type SensorThresholdController struct {
	sensorThresholdService *service.SensorThresholdService
}

// NewSensorThresholdController creates a new SensorThresholdController
func NewSensorThresholdController(sensorThresholdService *service.SensorThresholdService) *SensorThresholdController {
	return &SensorThresholdController{
		sensorThresholdService: sensorThresholdService,
	}
}

// CreateThreshold handles POST /api/v1/superadmin/sensor-thresholds
func (c *SensorThresholdController) CreateThreshold(ctx *gin.Context) {
	var req dto.CreateSensorThresholdRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Failed to bind JSON request: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid request body",
		})
		return
	}
	log.Printf("Creating sensor threshold with request: %+v", req)

	// Extract tenant ID from context for authorization
	tenantID, hasTenantID := common.GetTenantID(ctx.Request.Context())
	if !hasTenantID && !common.IsSuperAdmin(ctx.Request.Context()) {
		log.Printf("No tenant ID found in context for regular user")
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "Tenant context is required",
		})
		return
	}

	threshold, err := (*c.sensorThresholdService).CreateThreshold(ctx.Request.Context(), tenantID, &req)
	if err != nil {
		log.Printf("Error creating sensor threshold: %v", err)
		if common.IsValidationError(err) {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": err.Error(),
			})
			return
		}
		if common.IsNotFoundError(err) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error":   "Not Found",
				"message": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Failed to create sensor threshold",
		})
		return
	}

	log.Printf("Successfully created sensor threshold with ID: %s", threshold.ID)
	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Sensor threshold created successfully",
		"data":    threshold,
	})
}

// GetThreshold handles GET /api/v1/superadmin/sensor-thresholds/:id
func (c *SensorThresholdController) GetThreshold(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		log.Printf("Invalid threshold ID format: %s", idStr)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid threshold ID format",
		})
		return
	}
	log.Printf("Getting sensor threshold with ID: %s", id)

	// Extract tenant ID from context for authorization
	tenantID, hasTenantID := common.GetTenantID(ctx.Request.Context())
	if !hasTenantID && !common.IsSuperAdmin(ctx.Request.Context()) {
		log.Printf("No tenant ID found in context for regular user")
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "Tenant context is required",
		})
		return
	}

	threshold, err := (*c.sensorThresholdService).GetThresholdByID(ctx.Request.Context(), tenantID, id)
	if err != nil {
		log.Printf("Error getting sensor threshold: %v", err)
		if common.IsNotFoundError(err) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error":   "Not Found",
				"message": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Failed to get sensor threshold",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Sensor threshold retrieved successfully",
		"data":    threshold,
	})
}

// GetThresholds handles GET /api/v1/superadmin/sensor-thresholds
func (c *SensorThresholdController) GetThresholds(ctx *gin.Context) {
	// Parse query parameters
	var req dto.GetSensorThresholdsRequest

	// Parse page
	if pageStr := ctx.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			req.Page = page
		}
	}
	if req.Page == 0 {
		req.Page = 1
	}

	// Parse limit
	if limitStr := ctx.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			req.Limit = limit
		}
	}
	if req.Limit == 0 {
		req.Limit = 10
	}

	// Parse sort_by
	if sortBy := ctx.Query("sort_by"); sortBy != "" {
		req.SortBy = sortBy
	}

	// Parse sort_order
	if sortOrder := ctx.Query("sort_order"); sortOrder != "" {
		req.SortOrder = sortOrder
	}
	// Parse asset_sensor_id
	if assetSensorIDStr := ctx.Query("asset_sensor_id"); assetSensorIDStr != "" {
		if assetSensorID, err := uuid.Parse(assetSensorIDStr); err == nil {
			req.AssetSensorID = &assetSensorID
		}
	}

	// Parse sensor_type_id
	if sensorTypeIDStr := ctx.Query("sensor_type_id"); sensorTypeIDStr != "" {
		if sensorTypeID, err := uuid.Parse(sensorTypeIDStr); err == nil {
			req.SensorTypeID = &sensorTypeID
		}
	}

	// Parse is_active
	if isActiveStr := ctx.Query("is_active"); isActiveStr != "" {
		if isActive, err := strconv.ParseBool(isActiveStr); err == nil {
			req.IsActive = &isActive
		}
	}
	log.Printf("Getting sensor thresholds with request: %+v", req)

	// Extract tenant ID from context for authorization
	tenantID, hasTenantID := common.GetTenantID(ctx.Request.Context())
	if !hasTenantID && !common.IsSuperAdmin(ctx.Request.Context()) {
		log.Printf("No tenant ID found in context for regular user")
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "Tenant context is required",
		})
		return
	}

	// Convert GetSensorThresholdsRequest to SensorThresholdFilterRequest
	filter := &dto.SensorThresholdFilterRequest{
		Page:          req.Page,
		Limit:         req.Limit,
		AssetSensorID: req.AssetSensorID,
		SensorTypeID:  req.SensorTypeID,
		Severity:      req.Severity,
		IsActive:      req.IsActive,
	}

	result, err := (*c.sensorThresholdService).ListThresholds(ctx.Request.Context(), tenantID, filter)
	if err != nil {
		log.Printf("Error getting sensor thresholds: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Failed to get sensor thresholds",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Sensor thresholds retrieved successfully",
		"data":    result.Thresholds,
		"meta": gin.H{
			"page":        result.Page,
			"limit":       result.Limit,
			"total":       result.Total,
			"total_pages": result.TotalPages,
		},
	})
}

// UpdateThreshold handles PUT /api/v1/superadmin/sensor-thresholds/:id
func (c *SensorThresholdController) UpdateThreshold(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		log.Printf("Invalid threshold ID format: %s", idStr)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid threshold ID format",
		})
		return
	}

	var req dto.UpdateSensorThresholdRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Failed to bind JSON request: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid request body",
		})
		return
	}
	log.Printf("Updating sensor threshold %s with request: %+v", id, req)

	// Extract tenant ID from context for authorization
	tenantID, hasTenantID := common.GetTenantID(ctx.Request.Context())
	if !hasTenantID && !common.IsSuperAdmin(ctx.Request.Context()) {
		log.Printf("No tenant ID found in context for regular user")
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "Tenant context is required",
		})
		return
	}

	threshold, err := (*c.sensorThresholdService).UpdateThreshold(ctx.Request.Context(), tenantID, id, &req)
	if err != nil {
		log.Printf("Error updating sensor threshold: %v", err)
		if common.IsValidationError(err) {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": err.Error(),
			})
			return
		}
		if common.IsNotFoundError(err) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error":   "Not Found",
				"message": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Failed to update sensor threshold",
		})
		return
	}

	log.Printf("Successfully updated sensor threshold with ID: %s", threshold.ID)
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Sensor threshold updated successfully",
		"data":    threshold,
	})
}

// DeleteThreshold handles DELETE /api/v1/superadmin/sensor-thresholds/:id
func (c *SensorThresholdController) DeleteThreshold(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		log.Printf("Invalid threshold ID format: %s", idStr)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid threshold ID format",
		})
		return
	}
	log.Printf("Deleting sensor threshold with ID: %s", id)

	// Extract tenant ID from context for authorization
	tenantID, hasTenantID := common.GetTenantID(ctx.Request.Context())
	if !hasTenantID && !common.IsSuperAdmin(ctx.Request.Context()) {
		log.Printf("No tenant ID found in context for regular user")
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "Tenant context is required",
		})
		return
	}

	err = (*c.sensorThresholdService).DeleteThreshold(ctx.Request.Context(), tenantID, id)
	if err != nil {
		log.Printf("Error deleting sensor threshold: %v", err)
		if common.IsNotFoundError(err) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error":   "Not Found",
				"message": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Failed to delete sensor threshold",
		})
		return
	}

	log.Printf("Successfully deleted sensor threshold with ID: %s", id)
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Sensor threshold deleted successfully",
	})
}

// CheckThresholds handles POST /api/v1/superadmin/sensor-thresholds/check
func (c *SensorThresholdController) CheckThresholds(ctx *gin.Context) {
	var req dto.CheckThresholdsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Failed to bind JSON request: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid request body",
		})
		return
	}
	log.Printf("Checking thresholds for sensor reading: %+v", req)

	// Extract tenant ID from context for authorization
	tenantID, hasTenantID := common.GetTenantID(ctx.Request.Context())
	if !hasTenantID && !common.IsSuperAdmin(ctx.Request.Context()) {
		log.Printf("No tenant ID found in context for regular user")
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "Tenant context is required",
		})
		return
	}

	result, err := (*c.sensorThresholdService).CheckThresholds(ctx.Request.Context(), tenantID, req.AssetSensorID, req.MeasurementField, req.Value)
	if err != nil {
		log.Printf("Error checking thresholds: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Failed to check thresholds",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Thresholds checked successfully",
		"data":    result,
	})
}

// GetThresholdStatistics handles GET /api/v1/superadmin/sensor-thresholds/statistics
func (c *SensorThresholdController) GetThresholdStatistics(ctx *gin.Context) {
	// Parse query parameters
	var req dto.GetThresholdStatisticsRequest

	// Parse asset_sensor_id
	if assetSensorIDStr := ctx.Query("asset_sensor_id"); assetSensorIDStr != "" {
		if assetSensorID, err := uuid.Parse(assetSensorIDStr); err == nil {
			req.AssetSensorID = &assetSensorID
		}
	}

	// Parse sensor_type_id
	if sensorTypeIDStr := ctx.Query("sensor_type_id"); sensorTypeIDStr != "" {
		if sensorTypeID, err := uuid.Parse(sensorTypeIDStr); err == nil {
			req.SensorTypeID = &sensorTypeID
		}
	}

	// Parse date range
	if startDateStr := ctx.Query("start_date"); startDateStr != "" {
		req.StartDate = &startDateStr
	}

	if endDateStr := ctx.Query("end_date"); endDateStr != "" {
		req.EndDate = &endDateStr
	}

	log.Printf("Getting threshold statistics with request: %+v", req)

	// Extract tenant ID from context for authorization
	tenantID, hasTenantID := common.GetTenantID(ctx.Request.Context())
	if !hasTenantID && !common.IsSuperAdmin(ctx.Request.Context()) {
		log.Printf("No tenant ID found in context for regular user")
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "Tenant context is required",
		})
		return
	}

	stats, err := (*c.sensorThresholdService).GetThresholdStatistics(ctx.Request.Context(), tenantID)
	if err != nil {
		log.Printf("Error getting threshold statistics: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Failed to get threshold statistics",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Threshold statistics retrieved successfully",
		"data":    stats,
	})
}

// GetSensorThresholds handles GET /api/v1/superadmin/sensors/:sensor_id/thresholds
func (c *SensorThresholdController) GetSensorThresholds(ctx *gin.Context) {
	sensorIDStr := ctx.Param("sensor_id")
	sensorID, err := uuid.Parse(sensorIDStr)
	if err != nil {
		log.Printf("Invalid sensor ID format: %s", sensorIDStr)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid sensor ID format",
		})
		return
	}
	log.Printf("Getting thresholds for sensor ID: %s", sensorID)

	// Extract tenant ID from context for authorization
	tenantID, hasTenantID := common.GetTenantID(ctx.Request.Context())
	if !hasTenantID && !common.IsSuperAdmin(ctx.Request.Context()) {
		log.Printf("No tenant ID found in context for regular user")
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "Tenant context is required",
		})
		return
	}

	thresholds, err := (*c.sensorThresholdService).GetThresholdsByAssetSensor(ctx.Request.Context(), tenantID, sensorID)
	if err != nil {
		log.Printf("Error getting sensor thresholds: %v", err)
		if common.IsNotFoundError(err) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error":   "Not Found",
				"message": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Failed to get sensor thresholds",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Sensor thresholds retrieved successfully",
		"data":    thresholds,
	})
}
