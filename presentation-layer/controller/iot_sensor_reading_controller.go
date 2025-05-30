package controller

import (
	"be-lecsens/asset_management/domain-layer/service"
	"be-lecsens/asset_management/helpers/common"
	"be-lecsens/asset_management/helpers/dto"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// IoTSensorReadingController handles HTTP requests for IoT sensor reading operations
type IoTSensorReadingController struct {
	iotSensorReadingService *service.IoTSensorReadingService
}

// NewIoTSensorReadingController creates a new IoTSensorReadingController
func NewIoTSensorReadingController(iotSensorReadingService *service.IoTSensorReadingService) *IoTSensorReadingController {
	return &IoTSensorReadingController{
		iotSensorReadingService: iotSensorReadingService,
	}
}

// CreateReading handles POST /api/v1/superadmin/iot-sensor-readings
func (c *IoTSensorReadingController) CreateReading(ctx *gin.Context) {
	var req dto.CreateIoTSensorReadingRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Failed to bind JSON request: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid request body",
		})
		return
	}

	log.Printf("Creating IoT sensor reading with request: %+v", req)

	reading, err := c.iotSensorReadingService.CreateIoTSensorReading(ctx, &req)
	if err != nil {
		log.Printf("Error creating IoT sensor reading: %v", err)
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
			"message": "Failed to create IoT sensor reading",
		})
		return
	}

	log.Printf("Successfully created IoT sensor reading: %+v", reading)

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "IoT sensor reading created successfully",
		"data":    reading,
	})
}

// CreateBatchReading handles POST /api/v1/superadmin/iot-sensor-readings/batch
func (c *IoTSensorReadingController) CreateBatchReading(ctx *gin.Context) {
	var req dto.CreateBatchIoTSensorReadingRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Failed to bind JSON request: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid request body",
		})
		return
	}

	log.Printf("Creating batch IoT sensor readings with %d readings", len(req.Readings))

	readings, err := c.iotSensorReadingService.CreateBatchIoTSensorReading(ctx, &req)
	if err != nil {
		log.Printf("Error creating batch IoT sensor readings: %v", err)
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
			"message": "Failed to create batch IoT sensor readings",
		})
		return
	}

	log.Printf("Successfully created %d IoT sensor readings", len(readings))

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "IoT sensor readings created successfully",
		"data":    readings,
		"count":   len(readings),
	})
}

// GetReading handles GET /api/v1/iot-sensor-readings/:id
func (c *IoTSensorReadingController) GetReading(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid reading ID format",
		})
		return
	}

	reading, err := c.iotSensorReadingService.GetIoTSensorReadingByID(ctx, id)
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
			"message": "Failed to retrieve IoT sensor reading",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "IoT sensor reading retrieved successfully",
		"data":    reading,
	})
}

// ListReadings handles GET /api/v1/iot-sensor-readings
func (c *IoTSensorReadingController) ListReadings(ctx *gin.Context) {
	// Parse pagination parameters
	page := 1
	pageSize := 10

	if pageParam := ctx.Query("page"); pageParam != "" {
		if p, err := strconv.Atoi(pageParam); err == nil && p > 0 {
			page = p
		}
	}

	if pageSizeParam := ctx.Query("page_size"); pageSizeParam != "" {
		if ps, err := strconv.Atoi(pageSizeParam); err == nil && ps > 0 && ps <= 100 {
			pageSize = ps
		}
	}

	req := &dto.IoTSensorReadingListRequest{
		Page:     page,
		PageSize: pageSize,
	}

	response, err := c.iotSensorReadingService.ListIoTSensorReadings(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Failed to list IoT sensor readings",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "IoT sensor readings listed successfully",
		"data":    response,
	})
}

// ListAllReadings handles GET /api/v1/superadmin/iot-sensor-readings (SuperAdmin only)
func (c *IoTSensorReadingController) ListAllReadings(ctx *gin.Context) {
	// Parse pagination parameters
	page := 1
	pageSize := 10

	if pageParam := ctx.Query("page"); pageParam != "" {
		if p, err := strconv.Atoi(pageParam); err == nil && p > 0 {
			page = p
		}
	}

	if pageSizeParam := ctx.Query("page_size"); pageSizeParam != "" {
		if ps, err := strconv.Atoi(pageSizeParam); err == nil && ps > 0 && ps <= 100 {
			pageSize = ps
		}
	}

	req := &dto.IoTSensorReadingListRequest{
		Page:     page,
		PageSize: pageSize,
	}

	response, err := c.iotSensorReadingService.ListIoTSensorReadings(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Failed to list all IoT sensor readings",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "All IoT sensor readings listed successfully",
		"data":    response,
	})
}

// ListReadingsByAssetSensor handles GET /api/v1/iot-sensor-readings/by-asset-sensor/:asset_sensor_id
func (c *IoTSensorReadingController) ListReadingsByAssetSensor(ctx *gin.Context) {
	assetSensorIDParam := ctx.Param("asset_sensor_id")
	assetSensorID, err := uuid.Parse(assetSensorIDParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid asset sensor ID format",
		})
		return
	}

	// Parse limit parameter (instead of pagination)
	limit := 100 // Default limit
	if limitParam := ctx.Query("limit"); limitParam != "" {
		if l, err := strconv.Atoi(limitParam); err == nil && l > 0 && l <= 1000 {
			limit = l
		}
	}

	response, err := c.iotSensorReadingService.GetReadingsByAssetSensor(ctx, assetSensorID, limit)
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
			"message": "Failed to retrieve readings by asset sensor",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "IoT sensor readings by asset sensor retrieved successfully",
		"data":    response,
	})
}

// ListReadingsBySensorType handles GET /api/v1/iot-sensor-readings/by-sensor-type/:sensor_type_id
func (c *IoTSensorReadingController) ListReadingsBySensorType(ctx *gin.Context) {
	sensorTypeIDParam := ctx.Param("sensor_type_id")
	sensorTypeID, err := uuid.Parse(sensorTypeIDParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid sensor type ID format",
		})
		return
	}

	// Parse limit parameter (instead of pagination)
	limit := 100 // Default limit
	if limitParam := ctx.Query("limit"); limitParam != "" {
		if l, err := strconv.Atoi(limitParam); err == nil && l > 0 && l <= 1000 {
			limit = l
		}
	}

	response, err := c.iotSensorReadingService.GetReadingsBySensorType(ctx, sensorTypeID, limit)
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
			"message": "Failed to retrieve readings by sensor type",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "IoT sensor readings by sensor type retrieved successfully",
		"data":    response,
	})
}

// ListReadingsByMacAddress handles GET /api/v1/iot-sensor-readings/by-mac-address/:mac_address
func (c *IoTSensorReadingController) ListReadingsByMacAddress(ctx *gin.Context) {
	macAddress := ctx.Param("mac_address")
	if macAddress == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "MAC address is required",
		})
		return
	}

	// Parse limit parameter (instead of pagination)
	limit := 100 // Default limit
	if limitParam := ctx.Query("limit"); limitParam != "" {
		if l, err := strconv.Atoi(limitParam); err == nil && l > 0 && l <= 1000 {
			limit = l
		}
	}

	response, err := c.iotSensorReadingService.GetReadingsByMacAddress(ctx, macAddress, limit)
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
			"message": "Failed to retrieve readings by MAC address",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "IoT sensor readings by MAC address retrieved successfully",
		"data":    response,
	})
}

// GetLatestByAssetSensor handles GET /api/v1/iot-sensor-readings/latest/by-asset-sensor/:asset_sensor_id
func (c *IoTSensorReadingController) GetLatestByAssetSensor(ctx *gin.Context) {
	assetSensorIDParam := ctx.Param("asset_sensor_id")
	assetSensorID, err := uuid.Parse(assetSensorIDParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid asset sensor ID format",
		})
		return
	}

	reading, err := c.iotSensorReadingService.GetLatestReading(ctx, assetSensorID)
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
			"message": "Failed to retrieve latest reading by asset sensor",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Latest IoT sensor reading by asset sensor retrieved successfully",
		"data":    reading,
	})
}

// ListReadingsByTimeRange handles GET /api/v1/iot-sensor-readings/by-time-range
func (c *IoTSensorReadingController) ListReadingsByTimeRange(ctx *gin.Context) {
	// Parse query parameters
	assetSensorIDParam := ctx.Query("asset_sensor_id")
	startTimeParam := ctx.Query("start_time")
	endTimeParam := ctx.Query("end_time")

	// Validate required parameters
	if assetSensorIDParam == "" || startTimeParam == "" || endTimeParam == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "asset_sensor_id, start_time, and end_time are required",
		})
		return
	}

	assetSensorID, err := uuid.Parse(assetSensorIDParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid asset sensor ID format",
		})
		return
	}

	startTime, err := time.Parse(time.RFC3339, startTimeParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid start_time format. Use RFC3339 format (e.g., 2006-01-02T15:04:05Z07:00)",
		})
		return
	}

	endTime, err := time.Parse(time.RFC3339, endTimeParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid end_time format. Use RFC3339 format (e.g., 2006-01-02T15:04:05Z07:00)",
		})
		return
	}

	// Parse pagination parameters (limit for time range)
	pageSize := 100 // Default limit for time range queries
	if pageSizeParam := ctx.Query("page_size"); pageSizeParam != "" {
		if ps, err := strconv.Atoi(pageSizeParam); err == nil && ps > 0 && ps <= 1000 {
			pageSize = ps
		}
	}

	req := &dto.GetReadingsInTimeRangeRequest{
		AssetSensorID: &assetSensorID,
		FromTime:      startTime,
		ToTime:        endTime,
		Limit:         pageSize,
	}

	response, err := c.iotSensorReadingService.GetReadingsInTimeRange(ctx, req)
	if err != nil {
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
			"message": "Failed to retrieve readings by time range",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "IoT sensor readings by time range retrieved successfully",
		"data":    response,
	})
}

// GetAggregatedData handles GET /api/v1/iot-sensor-readings/aggregated
func (c *IoTSensorReadingController) GetAggregatedData(ctx *gin.Context) {
	// Parse query parameters
	assetSensorIDParam := ctx.Query("asset_sensor_id")
	startTimeParam := ctx.Query("start_time")
	endTimeParam := ctx.Query("end_time")
	intervalParam := ctx.Query("interval")

	// Validate required parameters
	if assetSensorIDParam == "" || startTimeParam == "" || endTimeParam == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "asset_sensor_id, start_time, and end_time are required",
		})
		return
	}

	assetSensorID, err := uuid.Parse(assetSensorIDParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid asset sensor ID format",
		})
		return
	}

	startTime, err := time.Parse(time.RFC3339, startTimeParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid start_time format. Use RFC3339 format (e.g., 2006-01-02T15:04:05Z07:00)",
		})
		return
	}

	endTime, err := time.Parse(time.RFC3339, endTimeParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid end_time format. Use RFC3339 format (e.g., 2006-01-02T15:04:05Z07:00)",
		})
		return
	}

	// Set defaults for optional parameters
	intervalStr := "hour" // Default to hour
	if intervalParam != "" {
		intervalStr = intervalParam
	}

	req := &dto.GetAggregatedDataRequest{
		AssetSensorID: &assetSensorID,
		FromTime:      startTime,
		ToTime:        endTime,
		Interval:      intervalStr,
		AggregateBy:   []string{}, // Could be populated from query params if needed
	}

	response, err := c.iotSensorReadingService.GetAggregatedData(ctx, req)
	if err != nil {
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
			"message": "Failed to retrieve aggregated data",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Aggregated IoT sensor data retrieved successfully",
		"data":    response,
	})
}

// UpdateReading handles PUT /api/v1/superadmin/iot-sensor-readings/:id
func (c *IoTSensorReadingController) UpdateReading(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid reading ID format",
		})
		return
	}

	var req dto.UpdateIoTSensorReadingRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid request body",
		})
		return
	}

	reading, err := c.iotSensorReadingService.UpdateIoTSensorReading(ctx, id, &req)
	if err != nil {
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
			"message": "Failed to update IoT sensor reading",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "IoT sensor reading updated successfully",
		"data":    reading,
	})
}

// DeleteReading handles DELETE /api/v1/superadmin/iot-sensor-readings/:id
func (c *IoTSensorReadingController) DeleteReading(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid reading ID format",
		})
		return
	}

	err = c.iotSensorReadingService.DeleteIoTSensorReading(ctx, id)
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
			"message": "Failed to delete IoT sensor reading",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "IoT sensor reading deleted successfully",
	})
}

// ValidateAndCreateReading handles POST /api/v1/superadmin/iot-sensor-readings/validate
func (c *IoTSensorReadingController) ValidateAndCreateReading(ctx *gin.Context) {
	var req dto.CreateIoTSensorReadingRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Failed to bind JSON request: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid request body",
		})
		return
	}

	log.Printf("Validating and creating IoT sensor reading with request: %+v", req)

	validateReq := &dto.ValidateAndCreateRequest{
		CreateIoTSensorReadingRequest: req,
		ValidateSchema:                true,
	}

	reading, err := c.iotSensorReadingService.ValidateAndCreateReading(ctx, validateReq)
	if err != nil {
		log.Printf("Error validating and creating IoT sensor reading: %v", err)
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
			"message": "Failed to validate and create IoT sensor reading",
		})
		return
	}

	log.Printf("Successfully validated and created IoT sensor reading: %+v", reading)

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "IoT sensor reading validated and created successfully",
		"data":    reading,
	})
}

// Manual JSON Input Controllers

// CreateFromJSON handles POST /api/v1/superadmin/iot-sensor-readings/from-json
func (c *IoTSensorReadingController) CreateFromJSON(ctx *gin.Context) {
	var request struct {
		JSONData string `json:"json_data" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "json_data field is required",
		})
		return
	}

	reading, err := c.iotSensorReadingService.CreateFromJSONString(ctx, request.JSONData)
	if err != nil {
		if common.IsValidationError(err) {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Failed to create reading from JSON",
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "IoT sensor reading created from JSON successfully",
		"data":    reading,
	})
}

// CreateBatchFromJSON handles POST /api/v1/superadmin/iot-sensor-readings/batch-from-json
func (c *IoTSensorReadingController) CreateBatchFromJSON(ctx *gin.Context) {
	var request struct {
		JSONData string `json:"json_data" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "json_data field is required",
		})
		return
	}

	readings, err := c.iotSensorReadingService.CreateBatchFromJSONString(ctx, request.JSONData)
	if err != nil {
		if common.IsValidationError(err) {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Failed to create batch readings from JSON",
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "IoT sensor readings created from JSON successfully",
		"data":    readings,
		"count":   len(readings),
	})
}

// CreateDummyReading handles POST /api/v1/superadmin/iot-sensor-readings/dummy
func (c *IoTSensorReadingController) CreateDummyReading(ctx *gin.Context) {
	var request struct {
		SensorType string `json:"sensor_type" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "sensor_type field is required",
		})
		return
	}

	reading, err := c.iotSensorReadingService.CreateDummyReading(ctx, request.SensorType)
	if err != nil {
		if common.IsValidationError(err) {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Failed to create dummy reading",
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Dummy IoT sensor reading created successfully",
		"data":    reading,
	})
}

// CreateMultipleDummyReadings handles POST /api/v1/superadmin/iot-sensor-readings/dummy/multiple
func (c *IoTSensorReadingController) CreateMultipleDummyReadings(ctx *gin.Context) {
	var request struct {
		Count       int      `json:"count" binding:"required,min=1,max=100"`
		SensorTypes []string `json:"sensor_types" binding:"required,min=1"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "count and sensor_types fields are required",
		})
		return
	}

	readings, err := c.iotSensorReadingService.CreateMultipleDummyReadings(ctx, request.Count, request.SensorTypes)
	if err != nil {
		if common.IsValidationError(err) {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Failed to create multiple dummy readings",
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Multiple dummy IoT sensor readings created successfully",
		"data":    readings,
		"count":   len(readings),
	})
}

// GetJSONTemplate handles GET /api/v1/superadmin/iot-sensor-readings/template
func (c *IoTSensorReadingController) GetJSONTemplate(ctx *gin.Context) {
	sensorType := ctx.Query("sensor_type")
	if sensorType == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "sensor_type query parameter is required",
		})
		return
	}

	template := map[string]interface{}{
		"asset_sensor_id": "00000000-0000-0000-0000-000000000000", // Replace with actual UUID
		"sensor_type_id":  "00000000-0000-0000-0000-000000000000", // Replace with actual UUID
		"mac_address":     fmt.Sprintf("%s_sensor_001", sensorType),
		"reading_time":    time.Now().Format(time.RFC3339),
	}

	jsonBytes, err := json.MarshalIndent(template, "", "  ")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Failed to generate JSON template",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":     "JSON template generated successfully",
		"sensor_type": sensorType,
		"template":    string(jsonBytes),
	})
}

// GetBatchJSONTemplate handles GET /api/v1/superadmin/iot-sensor-readings/template/batch
func (c *IoTSensorReadingController) GetBatchJSONTemplate(ctx *gin.Context) {
	sensorTypesParam := ctx.Query("sensor_types")
	countParam := ctx.Query("count")

	if sensorTypesParam == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "sensor_types query parameter is required",
		})
		return
	}

	sensorTypes := strings.Split(sensorTypesParam, ",")
	count := 2 // Default count

	if countParam != "" {
		if c, err := strconv.Atoi(countParam); err == nil && c > 0 && c <= 50 {
			count = c
		}
	}

	template, err := c.iotSensorReadingService.GetBatchJSONTemplate(sensorTypes, count)
	if err != nil {
		if common.IsValidationError(err) {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Failed to generate batch JSON template",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":      "Batch JSON template generated successfully",
		"sensor_types": sensorTypes,
		"count":        count,
		"template":     template,
	})
}

// CreateSimpleReading handles POST /api/v1/superadmin/iot-sensor-readings/simple
func (c *IoTSensorReadingController) CreateSimpleReading(ctx *gin.Context) {
	var request struct {
		AssetSensorID string    `json:"asset_sensor_id" binding:"required"`
		SensorTypeID  string    `json:"sensor_type_id" binding:"required"`
		MacAddress    string    `json:"mac_address" binding:"required"`
		ReadingTime   time.Time `json:"reading_time,omitempty"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid request body",
		})
		return
	}

	assetSensorID, err := uuid.Parse(request.AssetSensorID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid asset_sensor_id format",
		})
		return
	}

	sensorTypeID, err := uuid.Parse(request.SensorTypeID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid sensor_type_id format",
		})
		return
	}

	req := &dto.CreateIoTSensorReadingRequest{
		AssetSensorID: assetSensorID,
		SensorTypeID:  sensorTypeID,
		MacAddress:    request.MacAddress,
		ReadingTime:   &request.ReadingTime,
	}

	reading, err := c.iotSensorReadingService.CreateIoTSensorReading(ctx, req)
	if err != nil {
		if common.IsValidationError(err) {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Failed to create simple reading",
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Simple IoT sensor reading created successfully",
		"data":    reading,
	})
}

// Flexible JSON Controllers

// CreateFlexibleReading handles POST /api/v1/superadmin/iot-sensor-readings/flexible
func (c *IoTSensorReadingController) CreateFlexibleReading(ctx *gin.Context) {
	var req dto.FlexibleIoTSensorReadingRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Failed to bind flexible JSON request: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid request body",
		})
		return
	}

	log.Printf("Creating flexible IoT sensor reading with request: %+v", req)

	reading, err := c.iotSensorReadingService.CreateFlexibleIoTSensorReading(ctx, &req)
	if err != nil {
		log.Printf("Error creating flexible IoT sensor reading: %v", err)
		if common.IsValidationError(err) {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Failed to create flexible reading",
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Flexible IoT sensor reading created successfully",
		"data":    reading,
	})
}

// CreateBulkFlexibleReadings handles POST /api/v1/superadmin/iot-sensor-readings/flexible/bulk
func (c *IoTSensorReadingController) CreateBulkFlexibleReadings(ctx *gin.Context) {
	var req dto.FlexibleBatchIoTSensorReadingRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Failed to bind bulk flexible JSON request: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid request body",
		})
		return
	}

	log.Printf("Creating bulk flexible IoT sensor readings with %d readings", len(req.Readings))

	// Convert to slice of pointers
	var requests []*dto.FlexibleIoTSensorReadingRequest
	for i := range req.Readings {
		requests = append(requests, &req.Readings[i])
	}

	readings, err := c.iotSensorReadingService.CreateBulkFlexibleIoTSensorReadings(ctx, requests)
	if err != nil {
		log.Printf("Error creating bulk flexible IoT sensor readings: %v", err)
		if common.IsValidationError(err) {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Failed to create bulk flexible readings",
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": fmt.Sprintf("Successfully created %d flexible IoT sensor readings", len(readings)),
		"data":    readings,
		"count":   len(readings),
	})
}

// ParseTextToFlexible handles POST /api/v1/superadmin/iot-sensor-readings/parse-text
func (c *IoTSensorReadingController) ParseTextToFlexible(ctx *gin.Context) {
	var req dto.TextToJSONRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Failed to bind text parsing request: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid request body",
		})
		return
	}

	log.Printf("Parsing text to flexible IoT sensor reading")

	response, err := c.iotSensorReadingService.ParseTextToFlexibleReading(ctx, &req)
	if err != nil {
		log.Printf("Error parsing text to flexible reading: %v", err)
		if common.IsValidationError(err) {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Failed to parse text data",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Text data parsed successfully",
		"data":    response,
	})
}

// GetFlexibleReading handles GET /api/v1/superadmin/iot-sensor-readings/flexible/:id
func (c *IoTSensorReadingController) GetFlexibleReading(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid ID format",
		})
		return
	}

	reading, err := c.iotSensorReadingService.GetFlexibleIoTSensorReading(ctx, id)
	if err != nil {
		log.Printf("Error getting flexible IoT sensor reading: %v", err)
		if common.IsValidationError(err) || common.IsNotFoundError(err) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error":   "Not Found",
				"message": "Flexible IoT sensor reading not found",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Failed to get flexible reading",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Flexible IoT sensor reading retrieved successfully",
		"data":    reading,
	})
}

// CreateFromRawJSON handles POST /api/v1/superadmin/iot-sensor-readings/from-raw-json
func (c *IoTSensorReadingController) CreateFromRawJSON(ctx *gin.Context) {
	// Accept any JSON payload
	var rawJSON json.RawMessage
	if err := ctx.ShouldBindJSON(&rawJSON); err != nil {
		log.Printf("Failed to bind raw JSON request: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid JSON payload",
		})
		return
	}

	log.Printf("Creating IoT sensor reading from raw JSON")

	// Parse the raw JSON to flexible request
	var flexibleReq dto.FlexibleIoTSensorReadingRequest
	if err := json.Unmarshal(rawJSON, &flexibleReq); err != nil {
		log.Printf("Failed to parse raw JSON to flexible request: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Failed to parse JSON structure",
		})
		return
	}

	reading, err := c.iotSensorReadingService.CreateFlexibleIoTSensorReading(ctx, &flexibleReq)
	if err != nil {
		log.Printf("Error creating IoT sensor reading from raw JSON: %v", err)
		if common.IsValidationError(err) {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Failed to create reading from raw JSON",
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "IoT sensor reading created from raw JSON successfully",
		"data":    reading,
	})
}

// CreateFromArrayJSON handles POST /api/v1/superadmin/iot-sensor-readings/from-array
func (c *IoTSensorReadingController) CreateFromArrayJSON(ctx *gin.Context) {
	// Accept array of JSON payloads
	var rawJSONArray []json.RawMessage
	if err := ctx.ShouldBindJSON(&rawJSONArray); err != nil {
		log.Printf("Failed to bind raw JSON array request: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid JSON array payload",
		})
		return
	}

	if len(rawJSONArray) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Array cannot be empty",
		})
		return
	}

	if len(rawJSONArray) > 1000 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Maximum 1000 readings allowed per batch",
		})
		return
	}

	log.Printf("Creating %d IoT sensor readings from JSON array", len(rawJSONArray))

	var requests []*dto.FlexibleIoTSensorReadingRequest
	var errors []string

	// Parse each JSON in the array
	for i, rawJSON := range rawJSONArray {
		var flexibleReq dto.FlexibleIoTSensorReadingRequest
		if err := json.Unmarshal(rawJSON, &flexibleReq); err != nil {
			log.Printf("Failed to parse JSON at index %d: %v", i, err)
			errors = append(errors, fmt.Sprintf("index %d: %v", i, err))
			continue
		}
		requests = append(requests, &flexibleReq)
	}

	if len(errors) > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Some JSON items could not be parsed",
			"errors":  errors,
		})
		return
	}

	readings, err := c.iotSensorReadingService.CreateBulkFlexibleIoTSensorReadings(ctx, requests)
	if err != nil {
		log.Printf("Error creating IoT sensor readings from JSON array: %v", err)
		if common.IsValidationError(err) {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Failed to create readings from JSON array",
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": fmt.Sprintf("Successfully created %d IoT sensor readings from JSON array", len(readings)),
		"data":    readings,
		"count":   len(readings),
	})
}
