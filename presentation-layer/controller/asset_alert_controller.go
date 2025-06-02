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

// AssetAlertController handles HTTP requests for asset alert operations
type AssetAlertController struct {
	assetAlertService *service.AssetAlertService
}

// NewAssetAlertController creates a new AssetAlertController
func NewAssetAlertController(assetAlertService *service.AssetAlertService) *AssetAlertController {
	return &AssetAlertController{
		assetAlertService: assetAlertService,
	}
}

// GetAlert handles GET /api/v1/superadmin/asset-alerts/:id
func (c *AssetAlertController) GetAlert(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		log.Printf("Invalid alert ID format: %s", idStr)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid alert ID format",
		})
		return
	}

	log.Printf("Getting asset alert with ID: %s", id)

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

	alert, err := (*c.assetAlertService).GetAlertByID(ctx.Request.Context(), tenantID, id)
	if err != nil {
		log.Printf("Error getting asset alert: %v", err)
		if common.IsNotFoundError(err) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error":   "Not Found",
				"message": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Failed to get asset alert",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Asset alert retrieved successfully",
		"data":    alert,
	})
}

// GetAlerts handles GET /api/v1/superadmin/asset-alerts
func (c *AssetAlertController) GetAlerts(ctx *gin.Context) {
	// Parse query parameters
	var req dto.GetAssetAlertsRequest

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

	// Parse asset_id
	if assetIDStr := ctx.Query("asset_id"); assetIDStr != "" {
		if assetID, err := uuid.Parse(assetIDStr); err == nil {
			req.AssetID = &assetID
		}
	}

	// Parse sensor_id
	if sensorIDStr := ctx.Query("sensor_id"); sensorIDStr != "" {
		if sensorID, err := uuid.Parse(sensorIDStr); err == nil {
			req.SensorID = &sensorID
		}
	}

	// Parse alert_level
	if alertLevel := ctx.Query("alert_level"); alertLevel != "" {
		req.AlertLevel = &alertLevel
	}

	// Parse status
	if status := ctx.Query("status"); status != "" {
		req.Status = &status
	}

	// Parse is_resolved
	if isResolvedStr := ctx.Query("is_resolved"); isResolvedStr != "" {
		if isResolved, err := strconv.ParseBool(isResolvedStr); err == nil {
			req.IsResolved = &isResolved
		}
	}

	// Parse date range
	if startDateStr := ctx.Query("start_date"); startDateStr != "" {
		req.StartDate = &startDateStr
	}

	if endDateStr := ctx.Query("end_date"); endDateStr != "" {
		req.EndDate = &endDateStr
	}

	log.Printf("Getting asset alerts with request: %+v", req)

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

	// Convert GetAssetAlertsRequest to AssetAlertFilterRequest
	filter := &dto.AssetAlertFilterRequest{
		Page:       req.Page,
		Limit:      req.Limit,
		AssetID:    req.AssetID,
		IsResolved: req.IsResolved,
		StartTime:  nil, // Handle date parsing if needed
		EndTime:    nil, // Handle date parsing if needed
	}

	// Map other fields appropriately
	if req.SensorID != nil {
		filter.AssetSensorID = req.SensorID
	}
	if req.AlertLevel != nil {
		filter.Severity = req.AlertLevel
	}

	result, err := (*c.assetAlertService).ListAlerts(ctx.Request.Context(), tenantID, filter)
	if err != nil {
		log.Printf("Error getting asset alerts: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Failed to get asset alerts",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Asset alerts retrieved successfully",
		"data":    result.Alerts,
		"meta": gin.H{
			"page":        result.Page,
			"limit":       result.Limit,
			"total":       result.Total,
			"total_pages": result.TotalPages,
		},
	})
}

// ResolveAlert handles PUT /api/v1/superadmin/asset-alerts/:id/resolve
func (c *AssetAlertController) ResolveAlert(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		log.Printf("Invalid alert ID format: %s", idStr)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid alert ID format",
		})
		return
	}

	var req dto.ResolveAssetAlertRequest
	if bindErr := ctx.ShouldBindJSON(&req); bindErr != nil {
		log.Printf("Failed to bind JSON request: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid request body",
		})
		return
	}
	log.Printf("Resolving asset alert %s with request: %+v", id, req)

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

	err = (*c.assetAlertService).ResolveAlert(ctx.Request.Context(), tenantID, id)
	if err != nil {
		log.Printf("Error resolving asset alert: %v", err)
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
			"message": "Failed to resolve asset alert",
		})
		return
	}

	log.Printf("Successfully resolved asset alert with ID: %s", id)
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Asset alert resolved successfully",
	})
}

// BulkResolveAlerts handles PUT /api/v1/superadmin/asset-alerts/bulk-resolve
func (c *AssetAlertController) BulkResolveAlerts(ctx *gin.Context) {
	var req dto.BulkResolveAssetAlertsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Failed to bind JSON request: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid request body",
		})
		return
	}

	log.Printf("Bulk resolving asset alerts with request: %+v", req)

	// Extract tenant ID from context for authorization
	tenantID, hasTenantID := common.GetTenantID(ctx.Request.Context())
	if !hasTenantID && !common.IsSuperAdmin(ctx.Request.Context()) {
		log.Printf("No tenant ID found in context for regular user")
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "Tenant context is required",
		})
		return
	} // Convert BulkResolveAssetAlertsRequest to BulkResolveAlertsRequest
	bulkReq := &dto.BulkResolveAlertsRequest{
		AlertIDs: req.AlertIDs,
	}

	result, err := (*c.assetAlertService).BulkResolveAlerts(ctx.Request.Context(), tenantID, bulkReq)
	if err != nil {
		log.Printf("Error bulk resolving asset alerts: %v", err)
		if common.IsValidationError(err) {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Failed to bulk resolve asset alerts",
		})
		return
	}

	log.Printf("Successfully resolved %d asset alerts", result.ResolvedCount)
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Asset alerts resolved successfully",
		"data":    result,
	})
}

// GetAlertStatistics handles GET /api/v1/superadmin/asset-alerts/statistics
func (c *AssetAlertController) GetAlertStatistics(ctx *gin.Context) {
	// Parse query parameters
	var req dto.GetAlertStatisticsRequest

	// Parse asset_id
	if assetIDStr := ctx.Query("asset_id"); assetIDStr != "" {
		if assetID, err := uuid.Parse(assetIDStr); err == nil {
			req.AssetID = &assetID
		}
	}

	// Parse date range
	if startTimeStr := ctx.Query("start_time"); startTimeStr != "" {
		if startTime, err := time.Parse(time.RFC3339, startTimeStr); err == nil {
			req.StartTime = &startTime
		}
	}

	if endTimeStr := ctx.Query("end_time"); endTimeStr != "" {
		if endTime, err := time.Parse(time.RFC3339, endTimeStr); err == nil {
			req.EndTime = &endTime
		}
	}

	// Parse period
	if period := ctx.Query("period"); period != "" {
		req.Period = &period
	}

	log.Printf("Getting alert statistics with request: %+v", req)

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

	// Convert to AlertStatisticsRequest for service call
	serviceReq := &dto.AlertStatisticsRequest{
		AssetID:   req.AssetID,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	}

	stats, err := (*c.assetAlertService).GetAlertStatistics(ctx.Request.Context(), tenantID, serviceReq)
	if err != nil {
		log.Printf("Error getting alert statistics: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Failed to get alert statistics",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Alert statistics retrieved successfully",
		"data":    stats,
	})
}

// GetAssetAlerts handles GET /api/v1/superadmin/assets/:asset_id/alerts
func (c *AssetAlertController) GetAssetAlerts(ctx *gin.Context) {
	assetIDStr := ctx.Param("id")
	assetID, err := uuid.Parse(assetIDStr)
	if err != nil {
		log.Printf("Invalid asset ID format: %s", assetIDStr)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid asset ID format",
		})
		return
	}

	// Parse additional query parameters
	var includeResolved bool
	if includeResolvedStr := ctx.Query("include_resolved"); includeResolvedStr != "" {
		if parsed, parseErr := strconv.ParseBool(includeResolvedStr); parseErr == nil {
			includeResolved = parsed
		}
	}
	log.Printf("Getting alerts for asset ID: %s, include_resolved: %v", assetID, includeResolved)

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

	alerts, err := (*c.assetAlertService).GetAlertsByAsset(ctx.Request.Context(), tenantID, assetID)
	if err != nil {
		log.Printf("Error getting asset alerts: %v", err)
		if common.IsNotFoundError(err) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error":   "Not Found",
				"message": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Failed to get asset alerts",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Asset alerts retrieved successfully",
		"data":    alerts,
	})
}

// GetSensorAlerts handles GET /api/v1/superadmin/sensors/:sensor_id/alerts
func (c *AssetAlertController) GetSensorAlerts(ctx *gin.Context) {
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

	// Parse additional query parameters
	var includeResolved bool
	if includeResolvedStr := ctx.Query("include_resolved"); includeResolvedStr != "" {
		if parsed, parseErr := strconv.ParseBool(includeResolvedStr); parseErr == nil {
			includeResolved = parsed
		}
	}
	log.Printf("Getting alerts for sensor ID: %s, include_resolved: %v", sensorID, includeResolved)

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

	alerts, err := (*c.assetAlertService).GetAlertsByAssetSensor(ctx.Request.Context(), tenantID, sensorID)
	if err != nil {
		log.Printf("Error getting sensor alerts: %v", err)
		if common.IsNotFoundError(err) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error":   "Not Found",
				"message": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Failed to get sensor alerts",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Sensor alerts retrieved successfully",
		"data":    alerts,
	})
}
