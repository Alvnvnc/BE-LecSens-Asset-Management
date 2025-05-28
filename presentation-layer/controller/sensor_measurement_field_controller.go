package controller

import (
	"be-lecsens/asset_management/domain-layer/service"
	"be-lecsens/asset_management/helpers/dto"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// SensorMeasurementFieldController handles HTTP requests for sensor measurement fields
type SensorMeasurementFieldController struct {
	service *service.SensorMeasurementFieldService
}

// NewSensorMeasurementFieldController creates a new instance of SensorMeasurementFieldController
func NewSensorMeasurementFieldController(service *service.SensorMeasurementFieldService) *SensorMeasurementFieldController {
	return &SensorMeasurementFieldController{
		service: service,
	}
}

// GetAll handles retrieving all sensor measurement fields
func (c *SensorMeasurementFieldController) GetAll(ctx *gin.Context) {
	// Get pagination parameters
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	fields, err := c.service.GetAll(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve sensor measurement fields",
			"details": err.Error(),
		})
		return
	}

	// Calculate pagination
	start := (page - 1) * limit
	end := start + limit
	if start >= len(fields) {
		start = len(fields)
	}
	if end > len(fields) {
		end = len(fields)
	}

	ctx.JSON(http.StatusOK, dto.SensorMeasurementFieldListResponse{
		Data:    fields[start:end],
		Total:   len(fields),
		Page:    page,
		Limit:   limit,
		Message: "Sensor measurement fields retrieved successfully",
	})
}

// Create handles the creation of a new sensor measurement field
func (c *SensorMeasurementFieldController) Create(ctx *gin.Context) {
	var req dto.CreateSensorMeasurementFieldRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	field, err := c.service.Create(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create sensor measurement field",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, dto.SensorMeasurementFieldResponse{
		Data:    *field,
		Message: "Sensor measurement field created successfully",
	})
}

// GetByID handles retrieving a sensor measurement field by its ID
func (c *SensorMeasurementFieldController) GetByID(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid ID format",
			"details": err.Error(),
		})
		return
	}

	field, err := c.service.GetByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error":   "Sensor measurement field not found",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.SensorMeasurementFieldResponse{
		Data:    *field,
		Message: "Sensor measurement field retrieved successfully",
	})
}

// GetByMeasurementTypeID handles retrieving all fields for a measurement type
func (c *SensorMeasurementFieldController) GetByMeasurementTypeID(ctx *gin.Context) {
	measurementTypeID, err := uuid.Parse(ctx.Param("measurement_type_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid measurement type ID format",
			"details": err.Error(),
		})
		return
	}

	// Get pagination parameters
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	fields, err := c.service.GetByMeasurementTypeID(ctx, measurementTypeID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve sensor measurement fields",
			"details": err.Error(),
		})
		return
	}

	// Calculate pagination
	start := (page - 1) * limit
	end := start + limit
	if start >= len(fields) {
		start = len(fields)
	}
	if end > len(fields) {
		end = len(fields)
	}

	ctx.JSON(http.StatusOK, dto.SensorMeasurementFieldListResponse{
		Data:    fields[start:end],
		Total:   len(fields),
		Page:    page,
		Limit:   limit,
		Message: "Sensor measurement fields retrieved successfully",
	})
}

// Update handles updating an existing sensor measurement field
func (c *SensorMeasurementFieldController) Update(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid ID format",
			"details": err.Error(),
		})
		return
	}

	var req dto.UpdateSensorMeasurementFieldRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	field, err := c.service.Update(ctx, id, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update sensor measurement field",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.SensorMeasurementFieldResponse{
		Data:    *field,
		Message: "Sensor measurement field updated successfully",
	})
}

// Delete handles deleting a sensor measurement field
func (c *SensorMeasurementFieldController) Delete(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid ID format",
			"details": err.Error(),
		})
		return
	}

	if err := c.service.Delete(ctx, id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete sensor measurement field",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Sensor measurement field deleted successfully",
	})
}

// GetRequiredFields handles retrieving all required fields for a measurement type
func (c *SensorMeasurementFieldController) GetRequiredFields(ctx *gin.Context) {
	measurementTypeID, err := uuid.Parse(ctx.Param("measurement_type_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid measurement type ID format",
			"details": err.Error(),
		})
		return
	}

	fields, err := c.service.GetRequiredFields(ctx, measurementTypeID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve required sensor measurement fields",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.SensorMeasurementFieldListResponse{
		Data:    fields,
		Total:   len(fields),
		Message: "Required sensor measurement fields retrieved successfully",
	})
}
