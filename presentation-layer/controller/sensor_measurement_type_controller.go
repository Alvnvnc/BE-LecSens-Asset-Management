package controller

import (
	"be-lecsens/asset_management/data-layer/config"
	"be-lecsens/asset_management/domain-layer/service"
	"be-lecsens/asset_management/helpers/dto"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// SensorMeasurementTypeController handles HTTP requests for sensor measurement type operations
type SensorMeasurementTypeController struct {
	sensorMeasurementTypeService *service.SensorMeasurementTypeService
	config                       *config.Config
}

// NewSensorMeasurementTypeController creates a new SensorMeasurementTypeController
func NewSensorMeasurementTypeController(sensorMeasurementTypeService *service.SensorMeasurementTypeService, cfg *config.Config) *SensorMeasurementTypeController {
	return &SensorMeasurementTypeController{
		sensorMeasurementTypeService: sensorMeasurementTypeService,
		config:                       cfg,
	}
}

// CreateSensorMeasurementType handles the creation of a new sensor measurement type
func (c *SensorMeasurementTypeController) CreateSensorMeasurementType(ctx *gin.Context) {
	var req dto.CreateSensorMeasurementTypeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sensorMeasurementType, err := c.sensorMeasurementTypeService.CreateSensorMeasurementType(ctx.Request.Context(), req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := dto.SensorMeasurementTypeResponse{
		Data:    *sensorMeasurementType,
		Message: "Sensor measurement type created successfully",
	}

	ctx.JSON(http.StatusCreated, response)
}

// GetSensorMeasurementType handles retrieving a sensor measurement type by ID
func (c *SensorMeasurementTypeController) GetSensorMeasurementType(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid sensor measurement type ID"})
		return
	}

	sensorMeasurementType, err := c.sensorMeasurementTypeService.GetSensorMeasurementType(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "sensor measurement type not found"})
		return
	}

	response := dto.SensorMeasurementTypeResponse{
		Data:    *sensorMeasurementType,
		Message: "Sensor measurement type retrieved successfully",
	}

	ctx.JSON(http.StatusOK, response)
}

// ListSensorMeasurementTypes handles retrieving a list of sensor measurement types with pagination
func (c *SensorMeasurementTypeController) ListSensorMeasurementTypes(ctx *gin.Context) {
	// Get pagination parameters
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

	sensorMeasurementTypes, total, err := c.sensorMeasurementTypeService.ListSensorMeasurementTypes(ctx.Request.Context(), page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list sensor measurement types: " + err.Error()})
		return
	}

	response := dto.SensorMeasurementTypeListResponse{
		Data:    sensorMeasurementTypes,
		Total:   total,
		Page:    page,
		Limit:   limit,
		Message: "Sensor measurement types retrieved successfully",
	}

	ctx.JSON(http.StatusOK, response)
}

// UpdateSensorMeasurementType handles updating an existing sensor measurement type
func (c *SensorMeasurementTypeController) UpdateSensorMeasurementType(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid sensor measurement type ID"})
		return
	}

	var req dto.UpdateSensorMeasurementTypeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedSensorMeasurementType, err := c.sensorMeasurementTypeService.UpdateSensorMeasurementType(ctx.Request.Context(), id, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := dto.SensorMeasurementTypeResponse{
		Data:    *updatedSensorMeasurementType,
		Message: "Sensor measurement type updated successfully",
	}

	ctx.JSON(http.StatusOK, response)
}

// DeleteSensorMeasurementType handles deleting a sensor measurement type
func (c *SensorMeasurementTypeController) DeleteSensorMeasurementType(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid sensor measurement type ID"})
		return
	}

	if err := c.sensorMeasurementTypeService.DeleteSensorMeasurementType(ctx.Request.Context(), id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}

// GetActiveSensorMeasurementTypes handles retrieving all active sensor measurement types
func (c *SensorMeasurementTypeController) GetActiveSensorMeasurementTypes(ctx *gin.Context) {
	sensorMeasurementTypes, err := c.sensorMeasurementTypeService.GetActiveSensorMeasurementTypes(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get active sensor measurement types: " + err.Error()})
		return
	}

	response := dto.SensorMeasurementTypeListResponse{
		Data:    sensorMeasurementTypes,
		Message: "Active sensor measurement types retrieved successfully",
	}

	ctx.JSON(http.StatusOK, response)
}

// GetSensorMeasurementTypesBySensorTypeID handles retrieving all measurement types for a specific sensor type
func (c *SensorMeasurementTypeController) GetSensorMeasurementTypesBySensorTypeID(ctx *gin.Context) {
	sensorTypeID, err := uuid.Parse(ctx.Param("sensor_type_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid sensor type ID"})
		return
	}

	sensorMeasurementTypes, err := c.sensorMeasurementTypeService.GetSensorMeasurementTypesBySensorTypeID(ctx.Request.Context(), sensorTypeID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get sensor measurement types: " + err.Error()})
		return
	}

	response := dto.SensorMeasurementTypeListResponse{
		Data:    sensorMeasurementTypes,
		Message: "Sensor measurement types retrieved successfully",
	}

	ctx.JSON(http.StatusOK, response)
}

// ListAllSensorMeasurementTypes handles retrieving all sensor measurement types for SuperAdmin
func (c *SensorMeasurementTypeController) ListAllSensorMeasurementTypes(ctx *gin.Context) {
	// Check if user has SuperAdmin role
	userRole, exists := ctx.Get("user_role")
	if !exists || userRole != "SUPERADMIN" {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "SuperAdmin access required"})
		return
	}

	// Get pagination parameters
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

	sensorMeasurementTypes, total, err := c.sensorMeasurementTypeService.ListSensorMeasurementTypes(ctx.Request.Context(), page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list sensor measurement types: " + err.Error()})
		return
	}

	response := dto.SensorMeasurementTypeListResponse{
		Data:    sensorMeasurementTypes,
		Total:   total,
		Page:    page,
		Limit:   limit,
		Message: "All sensor measurement types retrieved successfully",
	}

	ctx.JSON(http.StatusOK, response)
}
