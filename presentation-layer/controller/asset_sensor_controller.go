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

// AssetSensorController handles HTTP requests for asset sensor operations
type AssetSensorController struct {
	assetSensorService *service.AssetSensorService
}

// NewAssetSensorController creates a new AssetSensorController
func NewAssetSensorController(assetSensorService *service.AssetSensorService) *AssetSensorController {
	return &AssetSensorController{
		assetSensorService: assetSensorService,
	}
}

// CreateAssetSensor handles POST /api/v1/asset-sensors
func (c *AssetSensorController) CreateAssetSensor(ctx *gin.Context) {
	var req dto.CreateAssetSensorRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Failed to bind JSON request: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid request body",
		})
		return
	}

	log.Printf("Creating asset sensor with request: %+v", req)

	sensor, err := c.assetSensorService.CreateAssetSensor(ctx, &req)
	if err != nil {
		log.Printf("Error creating asset sensor: %v", err)
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
			"message": "Failed to create asset sensor",
		})
		return
	}

	log.Printf("Successfully created asset sensor: %+v", sensor)

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Asset sensor created successfully",
		"data":    sensor,
	})
}

// GetAssetSensor handles GET /api/v1/asset-sensors/:id
func (c *AssetSensorController) GetAssetSensor(ctx *gin.Context) {
	idParam := ctx.Param("id")
	log.Printf("DEBUG: GetAssetSensor called with ID parameter: %s", idParam)

	id, err := uuid.Parse(idParam)
	if err != nil {
		log.Printf("DEBUG: Invalid UUID format for ID: %s, error: %v", idParam, err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid sensor ID format",
		})
		return
	}

	log.Printf("DEBUG: Calling GetCompleteSensorInfo service with parsed UUID: %s", id)

	// Get complete sensor information with details for single sensor endpoint
	completeSensor, err := c.assetSensorService.GetCompleteSensorInfo(ctx, id)
	if err != nil {
		log.Printf("DEBUG: Error from GetCompleteSensorInfo service: %v", err)
		if common.IsNotFoundError(err) {
			log.Printf("DEBUG: Sensor not found, returning 404")
			ctx.JSON(http.StatusNotFound, gin.H{
				"error":   "Not Found",
				"message": err.Error(),
			})
			return
		}
		log.Printf("DEBUG: Internal server error, returning 500")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Failed to retrieve asset sensor",
		})
		return
	}

	log.Printf("DEBUG: Successfully retrieved complete sensor info, building response")
	log.Printf("DEBUG: Sensor ID: %s, Name: %s, SensorType.Version: %s",
		completeSensor.AssetSensor.ID, completeSensor.AssetSensor.Name, completeSensor.SensorType.Version)
	log.Printf("DEBUG: Number of measurement types: %d", len(completeSensor.MeasurementTypes))

	// Create detailed response by merging sensor data with additional info
	response := gin.H{
		"id":                  completeSensor.AssetSensor.ID,
		"tenant_id":           completeSensor.AssetSensor.TenantID,
		"asset_id":            completeSensor.AssetSensor.AssetID,
		"sensor_type_id":      completeSensor.AssetSensor.SensorTypeID,
		"name":                completeSensor.AssetSensor.Name,
		"status":              completeSensor.AssetSensor.Status,
		"configuration":       completeSensor.AssetSensor.Configuration,
		"last_reading_value":  completeSensor.AssetSensor.LastReadingValue,
		"last_reading_time":   completeSensor.AssetSensor.LastReadingTime,
		"last_reading_values": completeSensor.AssetSensor.LastReadingValues,
		"created_at":          completeSensor.AssetSensor.CreatedAt,
		"updated_at":          completeSensor.AssetSensor.UpdatedAt,
		"sensor_type":         completeSensor.SensorType,
		"measurement_types":   completeSensor.MeasurementTypes,
	}

	log.Printf("DEBUG: Response built successfully, returning JSON")

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Asset sensor retrieved successfully",
		"data":    response,
	})
}

// GetAssetSensors handles GET /api/v1/asset-sensors/asset/:asset_id
func (c *AssetSensorController) GetAssetSensors(ctx *gin.Context) {
	assetIDParam := ctx.Param("asset_id")
	assetID, err := uuid.Parse(assetIDParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid asset ID format",
		})
		return
	}

	sensors, err := c.assetSensorService.GetAssetSensors(ctx, assetID)
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
			"message": "Failed to retrieve asset sensors",
		})
		return
	}

	// Return only basic sensor information for list view
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Asset sensors retrieved successfully",
		"data":    sensors,
	})
}

// ListAssetSensors handles GET /api/v1/asset-sensors
func (c *AssetSensorController) ListAssetSensors(ctx *gin.Context) {
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

	response, err := c.assetSensorService.ListAssetSensors(ctx, page, pageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Failed to list asset sensors",
		})
		return
	}

	// Return only basic sensor information for list view
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Asset sensors listed successfully",
		"data": gin.H{
			"sensors":     response.Sensors,
			"page":        response.Page,
			"limit":       response.Limit,
			"total":       response.Total,
			"total_pages": response.TotalPages,
		},
	})
}

// UpdateAssetSensor handles PUT /api/v1/asset-sensors/:id
func (c *AssetSensorController) UpdateAssetSensor(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid sensor ID format",
		})
		return
	}

	var req dto.UpdateAssetSensorRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid request body",
		})
		return
	}

	sensor, err := c.assetSensorService.UpdateAssetSensor(ctx, id, &req)
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
			"message": "Failed to update asset sensor",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Asset sensor updated successfully",
		"data":    sensor,
	})
}

// DeleteAssetSensor handles DELETE /api/v1/asset-sensors/:id
func (c *AssetSensorController) DeleteAssetSensor(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid sensor ID format",
		})
		return
	}

	err = c.assetSensorService.DeleteAssetSensor(ctx, id)
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
			"message": "Failed to delete asset sensor",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Asset sensor deleted successfully",
	})
}

// DeleteAssetSensors handles DELETE /api/v1/asset-sensors/asset/:asset_id
func (c *AssetSensorController) DeleteAssetSensors(ctx *gin.Context) {
	assetIDParam := ctx.Param("asset_id")
	assetID, err := uuid.Parse(assetIDParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid asset ID format",
		})
		return
	}

	err = c.assetSensorService.DeleteAssetSensors(ctx, assetID)
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
			"message": "Failed to delete asset sensors",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Asset sensors deleted successfully",
	})
}

// UpdateSensorReading handles PUT /api/v1/asset-sensors/:id/reading
func (c *AssetSensorController) UpdateSensorReading(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid sensor ID format",
		})
		return
	}

	var req dto.UpdateSensorReadingRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid request body",
		})
		return
	}

	err = c.assetSensorService.UpdateSensorReading(ctx, id, &req)
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
			"message": "Failed to update sensor reading",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Sensor reading updated successfully",
	})
}

// GetActiveSensors handles GET /api/v1/asset-sensors/active
func (c *AssetSensorController) GetActiveSensors(ctx *gin.Context) {
	sensors, err := c.assetSensorService.GetActiveSensors(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Failed to get active sensors",
		})
		return
	}

	// Return only basic sensor information for list view
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Active sensors retrieved successfully",
		"data":    sensors,
	})
}

// GetSensorsByStatus handles GET /api/v1/asset-sensors/status/:status
func (c *AssetSensorController) GetSensorsByStatus(ctx *gin.Context) {
	status := ctx.Param("status")
	if status == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Status is required",
		})
		return
	}

	sensors, err := c.assetSensorService.GetSensorsByStatus(ctx, status)
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
			"message": "Failed to get sensors by status",
		})
		return
	}

	// Return only basic sensor information for list view
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Sensors retrieved successfully",
		"data":    sensors,
	})
}
