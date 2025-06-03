package controller

import (
	"be-lecsens/asset_management/data-layer/entity"
	"be-lecsens/asset_management/domain-layer/service"
	"be-lecsens/asset_management/helpers/common"
	"be-lecsens/asset_management/helpers/dto"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// SensorThresholdController handles HTTP requests for sensor thresholds
type SensorThresholdController struct {
	sensorThresholdService *service.SensorThresholdService
}

// NewSensorThresholdController creates a new sensor threshold controller
func NewSensorThresholdController(sensorThresholdService *service.SensorThresholdService) *SensorThresholdController {
	return &SensorThresholdController{
		sensorThresholdService: sensorThresholdService,
	}
}

// GetSensorThreshold retrieves a sensor threshold by ID
// @Summary Get sensor threshold by ID
// @Description Get a specific sensor threshold by its ID
// @Tags Sensor Thresholds
// @Produce json
// @Param id path string true "Sensor threshold ID"
// @Success 200 {object} dto.SensorThresholdResponse
// @Failure 400 {object} common.ErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /sensor-thresholds/{id} [get]
func (c *SensorThresholdController) GetSensorThreshold(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, common.ErrorResponse{
			Error:   "Invalid ID format",
			Message: "ID must be a valid UUID",
		})
		return
	}

	threshold, err := c.sensorThresholdService.GetSensorThresholdByID(ctx.Request.Context(), id)
	if err != nil {
		if notFoundErr, ok := err.(*common.NotFoundError); ok {
			ctx.JSON(http.StatusNotFound, common.ErrorResponse{
				Error:   "Sensor threshold not found",
				Message: notFoundErr.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, common.ErrorResponse{
			Error:   "Failed to get sensor threshold",
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.FromEntity(threshold))
}

// ListSensorThresholds retrieves paginated sensor thresholds
// @Summary List sensor thresholds
// @Description Get a paginated list of sensor thresholds for a tenant
// @Tags Sensor Thresholds
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of items per page (default: 20, max: 100)"
// @Param asset_sensor_id query string false "Filter by asset sensor ID"
// @Param measurement_type_id query string false "Filter by measurement type ID"
// @Param severity query string false "Filter by severity (warning, critical)"
// @Param is_active query bool false "Filter by active status"
// @Success 200 {object} dto.SensorThresholdListResponse
// @Failure 400 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /sensor-thresholds [get]
func (c *SensorThresholdController) ListSensorThresholds(ctx *gin.Context) {
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
	filter := dto.SensorThresholdFilter{
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

	if measurementTypeIDStr := ctx.Query("measurement_type_id"); measurementTypeIDStr != "" {
		if id, err := uuid.Parse(measurementTypeIDStr); err == nil {
			filter.MeasurementTypeID = &id
		} else {
			ctx.JSON(http.StatusBadRequest, common.ErrorResponse{
				Error:   "Invalid measurement type ID format",
				Message: "Measurement type ID must be a valid UUID",
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

	if isActiveStr := ctx.Query("is_active"); isActiveStr != "" {
		if active, err := strconv.ParseBool(isActiveStr); err == nil {
			filter.IsActive = &active
		} else {
			ctx.JSON(http.StatusBadRequest, common.ErrorResponse{
				Error:   "Invalid is_active value",
				Message: "is_active must be true or false",
			})
			return
		}
	}

	offset := (filter.Page - 1) * filter.Limit

	thresholds, totalCount, err := c.sensorThresholdService.ListSensorThresholds(
		ctx.Request.Context(), tenantUUID, filter.Limit, offset)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.ErrorResponse{
			Error:   "Failed to list sensor thresholds",
			Message: err.Error(),
		})
		return
	}

	totalPages := (totalCount + filter.Limit - 1) / filter.Limit

	response := dto.SensorThresholdListResponse{
		Data: make([]dto.SensorThresholdResponse, len(thresholds)),
		Pagination: dto.PaginationInfo{
			Page:        filter.Page,
			Limit:       filter.Limit,
			TotalItems:  int64(totalCount),
			TotalPages:  totalPages,
			HasNext:     filter.Page < totalPages,
			HasPrevious: filter.Page > 1,
		},
	}

	for i, threshold := range thresholds {
		response.Data[i] = *dto.FromEntity(threshold)
	}

	ctx.JSON(http.StatusOK, response)
}

// ListThresholdsByAssetSensor retrieves thresholds for a specific asset sensor
// @Summary List thresholds by asset sensor
// @Description Get thresholds for a specific asset sensor
// @Tags Sensor Thresholds
// @Produce json
// @Param asset_sensor_id path string true "Asset Sensor ID"
// @Success 200 {array} entity.SensorThreshold
// @Failure 400 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /sensor-thresholds/by-asset-sensor/{asset_sensor_id} [get]
func (c *SensorThresholdController) ListThresholdsByAssetSensor(ctx *gin.Context) {
	assetSensorIDStr := ctx.Param("asset_sensor_id")
	assetSensorID, err := uuid.Parse(assetSensorIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, common.ErrorResponse{
			Error:   "Invalid asset sensor ID format",
			Message: "Asset sensor ID must be a valid UUID",
		})
		return
	}

	thresholds, err := c.sensorThresholdService.GetThresholdsByAssetSensor(ctx.Request.Context(), assetSensorID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.ErrorResponse{
			Error:   "Failed to get sensor thresholds",
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, thresholds)
}

// ListThresholdsByMeasurementType retrieves thresholds for a specific measurement type
// @Summary List thresholds by measurement type
// @Description Get thresholds for a specific measurement type
// @Tags Sensor Thresholds
// @Produce json
// @Param measurement_type_id path string true "Measurement Type ID"
// @Success 200 {array} entity.SensorThreshold
// @Failure 400 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /sensor-thresholds/by-measurement-type/{measurement_type_id} [get]
func (c *SensorThresholdController) ListThresholdsByMeasurementType(ctx *gin.Context) {
	measurementTypeIDStr := ctx.Param("measurement_type_id")
	measurementTypeID, err := uuid.Parse(measurementTypeIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, common.ErrorResponse{
			Error:   "Invalid measurement type ID format",
			Message: "Measurement type ID must be a valid UUID",
		})
		return
	}

	thresholds, err := c.sensorThresholdService.GetThresholdsByMeasurementType(ctx.Request.Context(), measurementTypeID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.ErrorResponse{
			Error:   "Failed to get sensor thresholds",
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, thresholds)
}

// CreateSensorThreshold creates a new sensor threshold
// @Summary Create sensor threshold
// @Description Create a new sensor threshold
// @Tags Sensor Thresholds
// @Accept json
// @Produce json
// @Param threshold body dto.CreateSensorThresholdRequest true "Sensor Threshold"
// @Success 201 {object} dto.SensorThresholdResponse
// @Failure 400 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /sensor-thresholds [post]
func (c *SensorThresholdController) CreateSensorThreshold(ctx *gin.Context) {
	var request dto.CreateSensorThresholdRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, common.ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
		return
	}

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

	threshold := request.ToEntity(tenantUUID)
	createdThreshold, err := c.sensorThresholdService.CreateSensorThreshold(ctx.Request.Context(), threshold)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.ErrorResponse{
			Error:   "Failed to create sensor threshold",
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, dto.FromEntity(createdThreshold))
}

// UpdateSensorThreshold updates an existing sensor threshold
// @Summary Update sensor threshold
// @Description Update an existing sensor threshold
// @Tags Sensor Thresholds
// @Accept json
// @Produce json
// @Param id path string true "Sensor Threshold ID"
// @Param threshold body dto.UpdateSensorThresholdRequest true "Sensor Threshold"
// @Success 200 {object} dto.SensorThresholdResponse
// @Failure 400 {object} common.ErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /sensor-thresholds/{id} [put]
func (c *SensorThresholdController) UpdateSensorThreshold(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, common.ErrorResponse{
			Error:   "Invalid ID format",
			Message: "ID must be a valid UUID",
		})
		return
	}

	var request dto.UpdateSensorThresholdRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, common.ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
		return
	}

	threshold := request.ToEntity(id)
	updatedThreshold, err := c.sensorThresholdService.UpdateSensorThreshold(ctx.Request.Context(), threshold)
	if err != nil {
		if notFoundErr, ok := err.(*common.NotFoundError); ok {
			ctx.JSON(http.StatusNotFound, common.ErrorResponse{
				Error:   "Sensor threshold not found",
				Message: notFoundErr.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, common.ErrorResponse{
			Error:   "Failed to update sensor threshold",
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.FromEntity(updatedThreshold))
}

// DeleteSensorThreshold deletes a sensor threshold
// @Summary Delete sensor threshold
// @Description Delete a sensor threshold by ID
// @Tags Sensor Thresholds
// @Produce json
// @Param id path string true "Sensor Threshold ID"
// @Success 204
// @Failure 400 {object} common.ErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /sensor-thresholds/{id} [delete]
func (c *SensorThresholdController) DeleteSensorThreshold(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, common.ErrorResponse{
			Error:   "Invalid ID format",
			Message: "ID must be a valid UUID",
		})
		return
	}

	err = c.sensorThresholdService.DeleteSensorThreshold(ctx.Request.Context(), id)
	if err != nil {
		if notFoundErr, ok := err.(*common.NotFoundError); ok {
			ctx.JSON(http.StatusNotFound, common.ErrorResponse{
				Error:   "Sensor threshold not found",
				Message: notFoundErr.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, common.ErrorResponse{
			Error:   "Failed to delete sensor threshold",
			Message: err.Error(),
		})
		return
	}

	ctx.Status(http.StatusNoContent)
}

// ActivateSensorThreshold activates a sensor threshold
// @Summary Activate sensor threshold
// @Description Activate a sensor threshold by ID
// @Tags Sensor Thresholds
// @Produce json
// @Param id path string true "Sensor Threshold ID"
// @Success 200 {object} dto.SensorThresholdResponse
// @Failure 400 {object} common.ErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /sensor-thresholds/{id}/activate [post]
func (c *SensorThresholdController) ActivateSensorThreshold(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, common.ErrorResponse{
			Error:   "Invalid ID format",
			Message: "ID must be a valid UUID",
		})
		return
	}

	threshold, err := c.sensorThresholdService.ActivateSensorThreshold(ctx.Request.Context(), id)
	if err != nil {
		if notFoundErr, ok := err.(*common.NotFoundError); ok {
			ctx.JSON(http.StatusNotFound, common.ErrorResponse{
				Error:   "Sensor threshold not found",
				Message: notFoundErr.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, common.ErrorResponse{
			Error:   "Failed to activate sensor threshold",
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.FromEntity(threshold))
}

// DeactivateSensorThreshold deactivates a sensor threshold
// @Summary Deactivate sensor threshold
// @Description Deactivate a sensor threshold by ID
// @Tags Sensor Thresholds
// @Produce json
// @Param id path string true "Sensor Threshold ID"
// @Success 200 {object} dto.SensorThresholdResponse
// @Failure 400 {object} common.ErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /sensor-thresholds/{id}/deactivate [post]
func (c *SensorThresholdController) DeactivateSensorThreshold(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, common.ErrorResponse{
			Error:   "Invalid ID format",
			Message: "ID must be a valid UUID",
		})
		return
	}

	threshold, err := c.sensorThresholdService.DeactivateSensorThreshold(ctx.Request.Context(), id)
	if err != nil {
		if notFoundErr, ok := err.(*common.NotFoundError); ok {
			ctx.JSON(http.StatusNotFound, common.ErrorResponse{
				Error:   "Sensor threshold not found",
				Message: notFoundErr.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, common.ErrorResponse{
			Error:   "Failed to deactivate sensor threshold",
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.FromEntity(threshold))
}

// ListAllSensorThresholds retrieves all sensor thresholds across all tenants (SuperAdmin only)
// @Summary List all sensor thresholds
// @Description Get all sensor thresholds across all tenants (SuperAdmin only)
// @Tags Sensor Thresholds
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of items per page (default: 20, max: 100)"
// @Success 200 {object} dto.SensorThresholdListResponse
// @Failure 400 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /sensor-thresholds [get]
func (c *SensorThresholdController) ListAllSensorThresholds(ctx *gin.Context) {
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

	thresholds, totalCount, err := c.sensorThresholdService.ListAllSensorThresholds(ctx.Request.Context(), limit, offset)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.ErrorResponse{
			Error:   "Failed to list sensor thresholds",
			Message: err.Error(),
		})
		return
	}

	totalPages := (totalCount + limit - 1) / limit

	response := dto.SensorThresholdListResponse{
		Data: make([]dto.SensorThresholdResponse, len(thresholds)),
		Pagination: dto.PaginationInfo{
			Page:        page,
			Limit:       limit,
			TotalItems:  int64(totalCount),
			TotalPages:  totalPages,
			HasNext:     page < totalPages,
			HasPrevious: page > 1,
		},
	}

	for i, threshold := range thresholds {
		response.Data[i] = *dto.FromEntity(threshold)
	}

	ctx.JSON(http.StatusOK, response)
}
