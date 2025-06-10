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

// SensorStatusController handles HTTP requests for sensor status operations
type SensorStatusController struct {
	service *service.SensorStatusService
}

// NewSensorStatusController creates a new instance of SensorStatusController
func NewSensorStatusController(service *service.SensorStatusService) *SensorStatusController {
	return &SensorStatusController{
		service: service,
	}
}

// CreateSensorStatus creates a new sensor status record
// @Summary Create sensor status
// @Description Create a new sensor status record
// @Tags sensor-status
// @Accept json
// @Produce json
// @Param body body dto.CreateSensorStatusRequest true "Create sensor status request"
// @Success 201 {object} dto.SensorStatusDTO
// @Failure 400 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /sensor-statuses [post]
func (c *SensorStatusController) CreateSensorStatus(ctx *gin.Context) {
	var req dto.CreateSensorStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Failed to bind JSON request: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid request body",
		})
		return
	}

	log.Printf("Creating sensor status with request: %+v", req)

	sensorStatus, err := c.service.CreateSensorStatus(ctx, req)
	if err != nil {
		log.Printf("Error creating sensor status: %v", err)
		if common.IsValidationError(err) {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Failed to create sensor status",
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Sensor status created successfully",
		"data":    sensorStatus,
	})
}

// GetSensorStatus retrieves a sensor status by ID
// @Summary Get sensor status
// @Description Get sensor status by ID
// @Tags sensor-status
// @Produce json
// @Param id path string true "Sensor status ID"
// @Success 200 {object} dto.SensorStatusDTO
// @Failure 404 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /sensor-statuses/{id} [get]
func (c *SensorStatusController) GetSensorStatus(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid sensor status ID format",
		})
		return
	}

	sensorStatus, err := c.service.GetSensorStatus(ctx, id)
	if err != nil {
		if common.IsNotFoundError(err) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error":   "Not Found",
				"message": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Failed to retrieve sensor status",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Sensor status retrieved successfully",
		"data":    sensorStatus,
	})
}

// ListSensorStatuses retrieves sensor statuses with pagination and filtering
// @Summary List sensor statuses
// @Description Get paginated list of sensor statuses with optional filters
// @Tags sensor-status
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Page size (default: 10)"
// @Param sensor_id query string false "Filter by sensor ID"
// @Param status query string false "Filter by status"
// @Success 200 {object} dto.SensorStatusListResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /sensor-statuses [get]
func (c *SensorStatusController) ListSensorStatuses(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	sensorID := ctx.Query("sensor_id")
	status := ctx.Query("status")

	// Parse sensor ID if provided
	var sensorUUID *uuid.UUID
	if sensorID != "" {
		if parsedID, err := uuid.Parse(sensorID); err == nil {
			sensorUUID = &parsedID
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": "Invalid sensor ID format",
			})
			return
		}
	}

	response, err := c.service.ListSensorStatuses(ctx, page, pageSize, sensorUUID, status)
	if err != nil {
		log.Printf("Error listing sensor statuses: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Failed to retrieve sensor statuses",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Sensor statuses retrieved successfully",
		"data":    response,
	})
}

// GetAllSensorStatuses retrieves all sensor statuses with pagination
// @Summary Get all sensor statuses
// @Description Get paginated list of all sensor statuses
// @Tags sensor-status
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Page size (default: 10)"
// @Param sensor_id query string false "Filter by sensor ID"
// @Param status query string false "Filter by status (online, offline, low_battery, weak_signal, unhealthy)"
// @Success 200 {object} dto.SensorStatusListResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /sensor-statuses [get]
func (c *SensorStatusController) GetAllSensorStatuses(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	sensorID := ctx.Query("sensor_id")
	status := ctx.Query("status")

	// Parse sensor ID if provided
	var sensorUUID *uuid.UUID
	if sensorID != "" {
		if parsedID, err := uuid.Parse(sensorID); err == nil {
			sensorUUID = &parsedID
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": "Invalid sensor ID format",
			})
			return
		}
	}

	response, err := c.service.ListSensorStatuses(ctx, page, pageSize, sensorUUID, status)
	if err != nil {
		log.Printf("Error getting all sensor statuses: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Failed to retrieve sensor statuses",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "All sensor statuses retrieved successfully",
		"data":    response,
	})
}

// UpdateSensorStatus updates an existing sensor status
// @Summary Update sensor status
// @Description Update an existing sensor status
// @Tags sensor-status
// @Accept json
// @Produce json
// @Param id path string true "Sensor status ID"
// @Param body body dto.UpdateSensorStatusRequest true "Update sensor status request"
// @Success 200 {object} dto.SensorStatusDTO
// @Failure 400 {object} common.ErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /sensor-statuses/{id} [put]
func (c *SensorStatusController) UpdateSensorStatus(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid sensor status ID format",
		})
		return
	}

	var req dto.UpdateSensorStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Failed to bind JSON request: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid request body",
		})
		return
	}

	sensorStatus, err := c.service.UpdateSensorStatus(ctx, id, req)
	if err != nil {
		log.Printf("Error updating sensor status: %v", err)
		if common.IsNotFoundError(err) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error":   "Not Found",
				"message": err.Error(),
			})
			return
		}
		if common.IsValidationError(err) {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Failed to update sensor status",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Sensor status updated successfully",
		"data":    sensorStatus,
	})
}

// DeleteSensorStatus deletes a sensor status by ID
// @Summary Delete sensor status
// @Description Delete sensor status by ID
// @Tags sensor-status
// @Produce json
// @Param id path string true "Sensor status ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /sensor-statuses/{id} [delete]
func (c *SensorStatusController) DeleteSensorStatus(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid sensor status ID format",
		})
		return
	}

	err = c.service.DeleteSensorStatus(ctx, id)
	if err != nil {
		if common.IsNotFoundError(err) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error":   "Not Found",
				"message": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Failed to delete sensor status",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Sensor status deleted successfully",
	})
}

// GetSensorStatusBySensorID retrieves the current status for a specific sensor
// @Summary Get sensor status by sensor ID
// @Description Get current status for a specific sensor
// @Tags sensor-status
// @Produce json
// @Param sensorId path string true "Sensor ID"
// @Success 200 {object} dto.SensorStatusDTO
// @Failure 404 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /sensors/{sensorId}/status [get]
func (c *SensorStatusController) GetSensorStatusBySensorID(ctx *gin.Context) {
	sensorIDParam := ctx.Param("sensorId")
	sensorID, err := uuid.Parse(sensorIDParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid sensor ID format",
		})
		return
	}

	sensorStatus, err := c.service.GetSensorStatusBySensorID(ctx, sensorID)
	if err != nil {
		if common.IsNotFoundError(err) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error":   "Not Found",
				"message": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Failed to retrieve sensor status",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Sensor status retrieved successfully",
		"data":    sensorStatus,
	})
}

// UpdateSensorStatusBySensorID updates the status for a specific sensor
// @Summary Update sensor status by sensor ID
// @Description Update status for a specific sensor
// @Tags sensor-status
// @Accept json
// @Produce json
// @Param sensorId path string true "Sensor ID"
// @Param body body dto.UpdateSensorStatusRequest true "Update sensor status request"
// @Success 200 {object} dto.SensorStatusDTO
// @Failure 400 {object} common.ErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /sensors/{sensorId}/status [put]
func (c *SensorStatusController) UpdateSensorStatusBySensorID(ctx *gin.Context) {
	sensorIDParam := ctx.Param("sensorId")
	sensorID, err := uuid.Parse(sensorIDParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid sensor ID format",
		})
		return
	}

	var req dto.UpdateSensorStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Failed to bind JSON request: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid request body",
		})
		return
	}

	sensorStatus, err := c.service.UpdateSensorStatusBySensorID(ctx, sensorID, &req)
	if err != nil {
		log.Printf("Error updating sensor status: %v", err)
		if common.IsNotFoundError(err) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error":   "Not Found",
				"message": err.Error(),
			})
			return
		}
		if common.IsValidationError(err) {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Failed to update sensor status",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Sensor status updated successfully",
		"data":    sensorStatus,
	})
}

// UpsertSensorStatus creates or updates sensor status for a specific sensor
// @Summary Upsert sensor status by sensor ID
// @Description Create or update status for a specific sensor
// @Tags sensor-status
// @Accept json
// @Produce json
// @Param sensorId path string true "Sensor ID"
// @Param body body dto.UpsertSensorStatusRequest true "Upsert sensor status request"
// @Success 200 {object} dto.SensorStatusDTO
// @Failure 400 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /sensors/{sensorId}/status/upsert [post]
func (c *SensorStatusController) UpsertSensorStatus(ctx *gin.Context) {
	sensorIDParam := ctx.Param("sensorId")
	sensorID, err := uuid.Parse(sensorIDParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid sensor ID format",
		})
		return
	}

	var req dto.UpsertSensorStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Failed to bind JSON request: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid request body",
		})
		return
	}

	// Convert UpsertSensorStatusRequest to CreateSensorStatusRequest
	createReq := req.ToCreateRequest()
	createReq.AssetSensorID = sensorID

	sensorStatus, err := c.service.UpsertSensorStatus(ctx, *createReq)
	if err != nil {
		log.Printf("Error upserting sensor status: %v", err)
		if common.IsValidationError(err) {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Failed to upsert sensor status",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Sensor status upserted successfully",
		"data":    sensorStatus,
	})
}

// RecordHeartbeat updates the last seen timestamp for a sensor
// @Summary Record sensor heartbeat
// @Description Update last seen timestamp for a sensor
// @Tags sensor-status
// @Accept json
// @Produce json
// @Param sensorId path string true "Sensor ID"
// @Success 200 {object} dto.SensorStatusDTO
// @Failure 400 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /sensors/{sensorId}/status/heartbeat [post]
func (c *SensorStatusController) RecordHeartbeat(ctx *gin.Context) {
	sensorIDParam := ctx.Param("sensorId")
	sensorID, err := uuid.Parse(sensorIDParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid sensor ID format",
		})
		return
	}

	sensorStatus, err := c.service.RecordHeartbeat(ctx, sensorID)
	if err != nil {
		log.Printf("Error recording heartbeat: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Failed to record heartbeat",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Heartbeat recorded successfully",
		"data":    sensorStatus,
	})
}

// ListOnlineSensors retrieves all sensors with online status
// @Summary List online sensors
// @Description Get paginated list of sensors with online status
// @Tags sensor-status
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Page size (default: 10)"
// @Success 200 {object} dto.SensorStatusListResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /sensor-statuses/online [get]
func (c *SensorStatusController) ListOnlineSensors(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))

	response, err := c.service.ListSensorStatuses(ctx, page, pageSize, nil, "online")
	if err != nil {
		log.Printf("Error listing online sensors: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Failed to retrieve online sensors",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Online sensors retrieved successfully",
		"data":    response,
	})
}

// GetOnlineSensors is an alias for ListOnlineSensors for route compatibility
func (c *SensorStatusController) GetOnlineSensors(ctx *gin.Context) {
	c.ListOnlineSensors(ctx)
}

// ListOfflineSensors retrieves all sensors with offline status
// @Summary List offline sensors
// @Description Get paginated list of sensors with offline status
// @Tags sensor-status
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Page size (default: 10)"
// @Success 200 {object} dto.SensorStatusListResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /sensor-statuses/offline [get]
func (c *SensorStatusController) ListOfflineSensors(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))

	response, err := c.service.ListSensorStatuses(ctx, page, pageSize, nil, "offline")
	if err != nil {
		log.Printf("Error listing offline sensors: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Failed to retrieve offline sensors",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Offline sensors retrieved successfully",
		"data":    response,
	})
}

// GetOfflineSensors is an alias for ListOfflineSensors for route compatibility
func (c *SensorStatusController) GetOfflineSensors(ctx *gin.Context) {
	c.ListOfflineSensors(ctx)
}

// ListLowBatterySensors retrieves all sensors with low battery status
// @Summary List low battery sensors
// @Description Get paginated list of sensors with low battery status
// @Tags sensor-status
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Page size (default: 10)"
// @Success 200 {object} dto.SensorStatusListResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /sensor-statuses/low-battery [get]
func (c *SensorStatusController) ListLowBatterySensors(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))

	response, err := c.service.ListSensorStatuses(ctx, page, pageSize, nil, "low_battery")
	if err != nil {
		log.Printf("Error listing low battery sensors: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Failed to retrieve low battery sensors",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Low battery sensors retrieved successfully",
		"data":    response,
	})
}

// GetLowBatterySensors is an alias for ListLowBatterySensors for route compatibility
func (c *SensorStatusController) GetLowBatterySensors(ctx *gin.Context) {
	c.ListLowBatterySensors(ctx)
}

// ListWeakSignalSensors retrieves all sensors with weak signal status
// @Summary List weak signal sensors
// @Description Get paginated list of sensors with weak signal status
// @Tags sensor-status
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Page size (default: 10)"
// @Success 200 {object} dto.SensorStatusListResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /sensor-statuses/weak-signal [get]
func (c *SensorStatusController) ListWeakSignalSensors(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))

	response, err := c.service.ListSensorStatuses(ctx, page, pageSize, nil, "weak_signal")
	if err != nil {
		log.Printf("Error listing weak signal sensors: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Failed to retrieve weak signal sensors",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Weak signal sensors retrieved successfully",
		"data":    response,
	})
}

// GetWeakSignalSensors is an alias for ListWeakSignalSensors for route compatibility
func (c *SensorStatusController) GetWeakSignalSensors(ctx *gin.Context) {
	c.ListWeakSignalSensors(ctx)
}

// ListUnhealthySensors retrieves all sensors with unhealthy status
// @Summary List unhealthy sensors
// @Description Get paginated list of sensors with unhealthy status
// @Tags sensor-status
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Page size (default: 10)"
// @Success 200 {object} dto.SensorStatusListResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /sensor-statuses/unhealthy [get]
func (c *SensorStatusController) ListUnhealthySensors(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))

	response, err := c.service.ListSensorStatuses(ctx, page, pageSize, nil, "unhealthy")
	if err != nil {
		log.Printf("Error listing unhealthy sensors: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Failed to retrieve unhealthy sensors",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Unhealthy sensors retrieved successfully",
		"data":    response,
	})
}

// GetUnhealthySensors is an alias for ListUnhealthySensors for route compatibility
func (c *SensorStatusController) GetUnhealthySensors(ctx *gin.Context) {
	c.ListUnhealthySensors(ctx)
}

// GetHealthSummary retrieves health summary for all sensors
// @Summary Get sensor health summary
// @Description Get health summary statistics for all sensors
// @Tags sensor-status
// @Produce json
// @Success 200 {object} dto.SensorHealthSummaryResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /sensor-statuses/health/summary [get]
func (c *SensorStatusController) GetHealthSummary(ctx *gin.Context) {
	summary, err := c.service.GetHealthSummary(ctx)
	if err != nil {
		log.Printf("Error getting health summary: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Failed to retrieve health summary",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Health summary retrieved successfully",
		"data":    summary,
	})
}

// GetSensorHealthSummary is an alias for GetHealthSummary for route compatibility
func (c *SensorStatusController) GetSensorHealthSummary(ctx *gin.Context) {
	c.GetHealthSummary(ctx)
}

// GetHealthAnalytics retrieves detailed health analytics
// @Summary Get sensor health analytics
// @Description Get detailed health analytics for sensors
// @Tags sensor-status
// @Produce json
// @Param timeframe query string false "Time frame for analytics (24h, 7d, 30d)"
// @Success 200 {object} dto.SensorHealthAnalyticsResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /sensor-statuses/health/analytics [get]
func (c *SensorStatusController) GetHealthAnalytics(ctx *gin.Context) {
	timeframe := ctx.DefaultQuery("timeframe", "24h")

	analytics, err := c.service.GetHealthAnalytics(ctx, timeframe)
	if err != nil {
		log.Printf("Error getting health analytics: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Failed to retrieve health analytics",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Health analytics retrieved successfully",
		"data":    analytics,
	})
}

// DeleteSensorStatusBySensorID deletes sensor status by sensor ID
// @Summary Delete sensor status by sensor ID
// @Description Delete sensor status record by asset sensor ID
// @Tags sensor-status
// @Param sensorId path string true "Asset Sensor ID"
// @Success 200 {object} common.SuccessResponse
// @Failure 400 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /sensors/{sensorId}/status [delete]
func (c *SensorStatusController) DeleteSensorStatusBySensorID(ctx *gin.Context) {
	sensorIDParam := ctx.Param("sensorId")
	sensorID, err := uuid.Parse(sensorIDParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid sensor ID format",
		})
		return
	}

	err = c.service.DeleteSensorStatusBySensorID(ctx, sensorID)
	if err != nil {
		log.Printf("Error deleting sensor status by sensor ID: %v", err)
		if common.IsValidationError(err) {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Failed to delete sensor status",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Sensor status deleted successfully",
	})
}

// UpdateHeartbeat is an alias for RecordHeartbeat for route compatibility
func (c *SensorStatusController) UpdateHeartbeat(ctx *gin.Context) {
	c.RecordHeartbeat(ctx)
}
