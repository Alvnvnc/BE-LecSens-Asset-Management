package controller

import (
	"be-lecsens/asset_management/data-layer/entity"
	"be-lecsens/asset_management/domain-layer/service"
	"be-lecsens/asset_management/helpers/common"
	"be-lecsens/asset_management/helpers/dto"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AssetAlertController handles HTTP requests for asset alerts
type AssetAlertController struct {
	assetAlertService *service.AssetAlertService
}

// NewAssetAlertController creates a new asset alert controller
func NewAssetAlertController(assetAlertService *service.AssetAlertService) *AssetAlertController {
	return &AssetAlertController{
		assetAlertService: assetAlertService,
	}
}

// GetAssetAlert retrieves an asset alert by ID
// @Summary Get asset alert by ID
// @Description Get a specific asset alert by its ID
// @Tags Asset Alerts
// @Produce json
// @Param id path string true "Asset alert ID"
// @Success 200 {object} dto.AssetAlertResponse
// @Failure 400 {object} common.ErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /asset-alerts/{id} [get]
func (c *AssetAlertController) GetAssetAlert(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, common.ErrorResponse{
			Error:   "Invalid ID format",
			Message: "ID must be a valid UUID",
		})
		return
	}

	alert, err := c.assetAlertService.GetAssetAlertByID(ctx.Request.Context(), id)
	if err != nil {
		if notFoundErr, ok := err.(*common.NotFoundError); ok {
			ctx.JSON(http.StatusNotFound, common.ErrorResponse{
				Error:   "Asset alert not found",
				Message: notFoundErr.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, common.ErrorResponse{
			Error:   "Failed to get asset alert",
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, alert)
}

// ListAssetAlerts retrieves paginated asset alerts
// @Summary List asset alerts
// @Description Get a paginated list of asset alerts for a tenant
// @Tags Asset Alerts
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of items per page (default: 20, max: 100)"
// @Param asset_id query string false "Filter by asset ID"
// @Param asset_sensor_id query string false "Filter by asset sensor ID"
// @Param severity query string false "Filter by severity (warning, critical)"
// @Param is_resolved query bool false "Filter by resolution status"
// @Param from_time query string false "Filter alerts from this time (RFC3339 format)"
// @Param to_time query string false "Filter alerts until this time (RFC3339 format)"
// @Success 200 {object} dto.AssetAlertListResponse
// @Failure 400 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /asset-alerts [get]
func (c *AssetAlertController) ListAssetAlerts(ctx *gin.Context) {
	// Get tenant ID from context
	tenantID, exists := ctx.Get("tenant_id")
	if !exists {
		ctx.JSON(http.StatusBadRequest, common.ErrorResponse{
			Error:   "Tenant ID not found",
			Message: "Tenant ID is required",
		})
		return
	}

	// Convert tenant ID to UUID
	tenantUUID, ok := tenantID.(uuid.UUID)
	if !ok {
		tenantUUIDStr, ok := tenantID.(string)
		if !ok {
			ctx.JSON(http.StatusBadRequest, common.ErrorResponse{
				Error:   "Invalid tenant ID format",
				Message: "Tenant ID must be a valid UUID",
			})
			return
		}
		var err error
		tenantUUID, err = uuid.Parse(tenantUUIDStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, common.ErrorResponse{
				Error:   "Invalid tenant ID format",
				Message: "Tenant ID must be a valid UUID",
			})
			return
		}
	}

	// Parse pagination parameters
	filter := dto.AssetAlertFilter{
		Page:  1,
		Limit: 20,
	}

	if pageStr := ctx.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			filter.Page = p
		}
	}

	if limitStr := ctx.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			filter.Limit = l
		}
	}

	// Parse filtering parameters
	if assetIDStr := ctx.Query("asset_id"); assetIDStr != "" {
		if id, err := uuid.Parse(assetIDStr); err == nil {
			filter.AssetID = &id
		} else {
			ctx.JSON(http.StatusBadRequest, common.ErrorResponse{
				Error:   "Invalid asset ID format",
				Message: "Asset ID must be a valid UUID",
			})
			return
		}
	}

	if assetSensorIDStr := ctx.Query("asset_sensor_id"); assetSensorIDStr != "" {
		if id, err := uuid.Parse(assetSensorIDStr); err == nil {
			filter.AssetSensorID = &id
		} else {
			ctx.JSON(http.StatusBadRequest, common.ErrorResponse{
				Error:   "Invalid asset sensor ID format",
				Message: "Asset sensor ID must be a valid UUID",
			})
			return
		}
	}

	if severityStr := ctx.Query("severity"); severityStr != "" {
		s := entity.ThresholdSeverity(severityStr)
		if s == entity.ThresholdSeverityWarning || s == entity.ThresholdSeverityCritical {
			filter.Severity = &s
		} else {
			ctx.JSON(http.StatusBadRequest, common.ErrorResponse{
				Error:   "Invalid severity",
				Message: "Severity must be 'warning' or 'critical'",
			})
			return
		}
	}

	if isResolvedStr := ctx.Query("is_resolved"); isResolvedStr != "" {
		if resolved, err := strconv.ParseBool(isResolvedStr); err == nil {
			filter.IsResolved = &resolved
		} else {
			ctx.JSON(http.StatusBadRequest, common.ErrorResponse{
				Error:   "Invalid is_resolved value",
				Message: "is_resolved must be true or false",
			})
			return
		}
	}

	if fromTimeStr := ctx.Query("from_time"); fromTimeStr != "" {
		if t, err := time.Parse(time.RFC3339, fromTimeStr); err == nil {
			filter.FromTime = &t
		} else {
			ctx.JSON(http.StatusBadRequest, common.ErrorResponse{
				Error:   "Invalid from_time format",
				Message: "from_time must be in RFC3339 format",
			})
			return
		}
	}

	if toTimeStr := ctx.Query("to_time"); toTimeStr != "" {
		if t, err := time.Parse(time.RFC3339, toTimeStr); err == nil {
			filter.ToTime = &t
		} else {
			ctx.JSON(http.StatusBadRequest, common.ErrorResponse{
				Error:   "Invalid to_time format",
				Message: "to_time must be in RFC3339 format",
			})
			return
		}
	}

	response, err := c.assetAlertService.ListAssetAlerts(ctx.Request.Context(), tenantUUID, filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.ErrorResponse{
			Error:   "Failed to list asset alerts",
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// ResolveAssetAlert resolves an asset alert
// @Summary Resolve asset alert
// @Description Mark an asset alert as resolved
// @Tags Asset Alerts
// @Produce json
// @Param id path string true "Asset alert ID"
// @Success 200 {object} dto.AssetAlertResponse
// @Failure 400 {object} common.ErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /asset-alerts/{id}/resolve [patch]
func (c *AssetAlertController) ResolveAssetAlert(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, common.ErrorResponse{
			Error:   "Invalid ID format",
			Message: "ID must be a valid UUID",
		})
		return
	}

	alert, err := c.assetAlertService.ResolveAssetAlert(ctx.Request.Context(), id)
	if err != nil {
		if notFoundErr, ok := err.(*common.NotFoundError); ok {
			ctx.JSON(http.StatusNotFound, common.ErrorResponse{
				Error:   "Asset alert not found",
				Message: notFoundErr.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, common.ErrorResponse{
			Error:   "Failed to resolve asset alert",
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, alert)
}

// ResolveMultipleAssetAlerts resolves multiple asset alerts
// @Summary Resolve multiple asset alerts
// @Description Mark multiple asset alerts as resolved
// @Tags Asset Alerts
// @Accept json
// @Produce json
// @Param request body dto.ResolveMultipleAlertsRequest true "Alert IDs to resolve"
// @Success 200 {object} dto.ResolveMultipleAlertsResponse
// @Failure 400 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /asset-alerts/resolve-multiple [patch]
func (c *AssetAlertController) ResolveMultipleAssetAlerts(ctx *gin.Context) {
	var request dto.ResolveMultipleAlertsRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, common.ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
		return
	}

	if len(request.AlertIDs) == 0 {
		ctx.JSON(http.StatusBadRequest, common.ErrorResponse{
			Error:   "No alert IDs provided",
			Message: "At least one alert ID is required",
		})
		return
	}

	if len(request.AlertIDs) > 100 {
		ctx.JSON(http.StatusBadRequest, common.ErrorResponse{
			Error:   "Too many alert IDs",
			Message: "Maximum 100 alert IDs allowed per request",
		})
		return
	}

	response, err := c.assetAlertService.ResolveMultipleAssetAlerts(ctx.Request.Context(), request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.ErrorResponse{
			Error:   "Failed to resolve alerts",
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// GetAlertStatistics retrieves alert statistics for a tenant
// @Summary Get alert statistics
// @Description Get alert statistics for monitoring dashboard
// @Tags Asset Alerts
// @Produce json
// @Param asset_id query string false "Filter by asset ID"
// @Param from_time query string false "Statistics from this time (RFC3339 format)"
// @Param to_time query string false "Statistics until this time (RFC3339 format)"
// @Success 200 {object} dto.AssetAlertStatisticsResponse
// @Failure 400 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /asset-alerts/statistics [get]
func (c *AssetAlertController) GetAlertStatistics(ctx *gin.Context) {
	// Get tenant ID from context
	tenantID, exists := ctx.Get("tenant_id")
	if !exists {
		ctx.JSON(http.StatusBadRequest, common.ErrorResponse{
			Error:   "Tenant ID not found",
			Message: "Tenant ID is required",
		})
		return
	}

	// Convert tenant ID to UUID
	tenantUUID, ok := tenantID.(uuid.UUID)
	if !ok {
		tenantUUIDStr, ok := tenantID.(string)
		if !ok {
			ctx.JSON(http.StatusBadRequest, common.ErrorResponse{
				Error:   "Invalid tenant ID format",
				Message: "Tenant ID must be a valid UUID",
			})
			return
		}
		var err error
		tenantUUID, err = uuid.Parse(tenantUUIDStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, common.ErrorResponse{
				Error:   "Invalid tenant ID format",
				Message: "Tenant ID must be a valid UUID",
			})
			return
		}
	}

	// Parse filtering parameters
	filter := dto.AssetAlertFilter{}

	if assetIDStr := ctx.Query("asset_id"); assetIDStr != "" {
		if id, err := uuid.Parse(assetIDStr); err == nil {
			filter.AssetID = &id
		} else {
			ctx.JSON(http.StatusBadRequest, common.ErrorResponse{
				Error:   "Invalid asset ID format",
				Message: "Asset ID must be a valid UUID",
			})
			return
		}
	}

	if fromTimeStr := ctx.Query("from_time"); fromTimeStr != "" {
		if t, err := time.Parse(time.RFC3339, fromTimeStr); err == nil {
			filter.FromTime = &t
		} else {
			ctx.JSON(http.StatusBadRequest, common.ErrorResponse{
				Error:   "Invalid from_time format",
				Message: "from_time must be in RFC3339 format",
			})
			return
		}
	}

	if toTimeStr := ctx.Query("to_time"); toTimeStr != "" {
		if t, err := time.Parse(time.RFC3339, toTimeStr); err == nil {
			filter.ToTime = &t
		} else {
			ctx.JSON(http.StatusBadRequest, common.ErrorResponse{
				Error:   "Invalid to_time format",
				Message: "to_time must be in RFC3339 format",
			})
			return
		}
	}

	statistics, err := c.assetAlertService.GetAlertStatistics(ctx.Request.Context(), tenantUUID, filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.ErrorResponse{
			Error:   "Failed to get alert statistics",
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, statistics)
}

// DeleteAssetAlert deletes an asset alert
// @Summary Delete asset alert
// @Description Delete an asset alert by ID
// @Tags Asset Alerts
// @Produce json
// @Param id path string true "Asset alert ID"
// @Success 204
// @Failure 400 {object} common.ErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /asset-alerts/{id} [delete]
func (c *AssetAlertController) DeleteAssetAlert(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, common.ErrorResponse{
			Error:   "Invalid ID format",
			Message: "ID must be a valid UUID",
		})
		return
	}

	err = c.assetAlertService.DeleteAssetAlert(ctx.Request.Context(), id)
	if err != nil {
		if notFoundErr, ok := err.(*common.NotFoundError); ok {
			ctx.JSON(http.StatusNotFound, common.ErrorResponse{
				Error:   "Asset alert not found",
				Message: notFoundErr.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, common.ErrorResponse{
			Error:   "Failed to delete asset alert",
			Message: err.Error(),
		})
		return
	}

	ctx.Status(http.StatusNoContent)
}

// ListAlertsByAsset retrieves alerts for a specific asset
// @Summary List alerts by asset
// @Description Get alerts for a specific asset
// @Tags Asset Alerts
// @Produce json
// @Param asset_id path string true "Asset ID"
// @Success 200 {array} entity.AssetAlert
// @Failure 400 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /asset-alerts/by-asset/{asset_id} [get]
func (c *AssetAlertController) ListAlertsByAsset(ctx *gin.Context) {
	assetIDStr := ctx.Param("asset_id")
	assetID, err := uuid.Parse(assetIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, common.ErrorResponse{
			Error:   "Invalid asset ID format",
			Message: "Asset ID must be a valid UUID",
		})
		return
	}

	alerts, err := c.assetAlertService.GetAssetAlertsByAsset(ctx.Request.Context(), assetID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.ErrorResponse{
			Error:   "Failed to get asset alerts",
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, alerts)
}

// ListAlertsByAssetSensor retrieves alerts for a specific asset sensor
// @Summary List alerts by asset sensor
// @Description Get alerts for a specific asset sensor
// @Tags Asset Alerts
// @Produce json
// @Param asset_sensor_id path string true "Asset Sensor ID"
// @Success 200 {array} entity.AssetAlert
// @Failure 400 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /asset-alerts/by-asset-sensor/{asset_sensor_id} [get]
func (c *AssetAlertController) ListAlertsByAssetSensor(ctx *gin.Context) {
	assetSensorIDStr := ctx.Param("asset_sensor_id")
	assetSensorID, err := uuid.Parse(assetSensorIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, common.ErrorResponse{
			Error:   "Invalid asset sensor ID format",
			Message: "Asset sensor ID must be a valid UUID",
		})
		return
	}

	alerts, err := c.assetAlertService.GetAssetAlertsByAssetSensor(ctx.Request.Context(), assetSensorID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.ErrorResponse{
			Error:   "Failed to get asset alerts",
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, alerts)
}

// ListAlertsByMeasurementType retrieves alerts for a specific measurement type
// @Summary List alerts by measurement type
// @Description Get alerts for a specific measurement type
// @Tags Asset Alerts
// @Produce json
// @Param measurement_type_id path string true "Measurement Type ID"
// @Success 200 {array} entity.AssetAlert
// @Failure 400 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /asset-alerts/by-measurement-type/{measurement_type_id} [get]
func (c *AssetAlertController) ListAlertsByMeasurementType(ctx *gin.Context) {
	measurementTypeIDStr := ctx.Param("measurement_type_id")
	measurementTypeID, err := uuid.Parse(measurementTypeIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, common.ErrorResponse{
			Error:   "Invalid measurement type ID format",
			Message: "Measurement type ID must be a valid UUID",
		})
		return
	}

	alerts, err := c.assetAlertService.GetAlertsByMeasurementType(ctx.Request.Context(), measurementTypeID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.ErrorResponse{
			Error:   "Failed to get asset alerts",
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, alerts)
}

// DeleteMultipleAssetAlerts deletes multiple asset alerts
// @Summary Delete multiple asset alerts
// @Description Delete multiple asset alerts by their IDs
// @Tags Asset Alerts
// @Accept json
// @Produce json
// @Param request body dto.DeleteMultipleAlertsRequest true "Alert IDs to delete"
// @Success 200 {object} dto.DeleteMultipleAlertsResponse
// @Failure 400 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /asset-alerts/delete-multiple [delete]
func (c *AssetAlertController) DeleteMultipleAssetAlerts(ctx *gin.Context) {
	var request dto.DeleteMultipleAlertsRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, common.ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
		return
	}

	if len(request.AlertIDs) == 0 {
		ctx.JSON(http.StatusBadRequest, common.ErrorResponse{
			Error:   "No alert IDs provided",
			Message: "At least one alert ID is required",
		})
		return
	}

	if len(request.AlertIDs) > 100 {
		ctx.JSON(http.StatusBadRequest, common.ErrorResponse{
			Error:   "Too many alert IDs",
			Message: "Maximum 100 alert IDs allowed per request",
		})
		return
	}

	response, err := c.assetAlertService.DeleteMultipleAssetAlerts(ctx.Request.Context(), request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.ErrorResponse{
			Error:   "Failed to delete alerts",
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// ListAllAssetAlerts retrieves all asset alerts across all tenants (SuperAdmin only)
// @Summary List all asset alerts
// @Description Get all asset alerts across all tenants (SuperAdmin only)
// @Tags Asset Alerts
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of items per page (default: 20, max: 100)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /asset-alerts [get]
func (c *AssetAlertController) ListAllAssetAlerts(ctx *gin.Context) {
	// Parse pagination parameters
	page := 1
	if pageStr := ctx.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	limit := 20
	if limitStr := ctx.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	offset := (page - 1) * limit

	alerts, totalCount, err := c.assetAlertService.ListAllAssetAlerts(ctx.Request.Context(), limit, offset)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.ErrorResponse{
			Error:   "Failed to list asset alerts",
			Message: err.Error(),
		})
		return
	}

	totalPages := (totalCount + limit - 1) / limit

	response := map[string]interface{}{
		"data": alerts,
		"pagination": map[string]interface{}{
			"page":        page,
			"limit":       limit,
			"total_items": totalCount,
			"total_pages": totalPages,
			"has_next":    page < totalPages,
			"has_prev":    page > 1,
		},
	}

	ctx.JSON(http.StatusOK, response)
}

// GetGlobalAlertStatistics retrieves alert statistics across all tenants (SuperAdmin only)
// @Summary Get global alert statistics
// @Description Get alert statistics across all tenants (SuperAdmin only)
// @Tags Asset Alerts
// @Produce json
// @Success 200 {object} dto.AssetAlertStatisticsResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /asset-alerts/statistics [get]
func (c *AssetAlertController) GetGlobalAlertStatistics(ctx *gin.Context) {
	statistics, err := c.assetAlertService.GetGlobalAlertStatistics(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.ErrorResponse{
			Error:   "Failed to get global alert statistics",
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, statistics)
}
