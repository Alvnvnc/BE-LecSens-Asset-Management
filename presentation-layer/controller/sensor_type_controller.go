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

// SensorTypeController handles HTTP requests for sensor type operations
type SensorTypeController struct {
	sensorTypeService *service.SensorTypeService
	config            *config.Config
}

// NewSensorTypeController creates a new SensorTypeController
func NewSensorTypeController(sensorTypeService *service.SensorTypeService, cfg *config.Config) *SensorTypeController {
	return &SensorTypeController{
		sensorTypeService: sensorTypeService,
		config:            cfg,
	}
}

// CreateSensorType handles the creation of a new sensor type
func (c *SensorTypeController) CreateSensorType(ctx *gin.Context) {
	var req dto.CreateSensorTypeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sensorType, err := c.sensorTypeService.CreateSensorType(ctx.Request.Context(), &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, dto.SensorTypeResponse{
		Success: true,
		Data:    sensorType,
	})
}

// GetSensorType handles retrieving a sensor type by ID
func (c *SensorTypeController) GetSensorType(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid sensor type ID"})
		return
	}

	sensorType, err := c.sensorTypeService.GetSensorType(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "sensor type not found"})
		return
	}

	ctx.JSON(http.StatusOK, dto.SensorTypeResponse{
		Success: true,
		Data:    sensorType,
	})
}

// ListSensorTypes handles retrieving a list of sensor types with pagination
func (c *SensorTypeController) ListSensorTypes(ctx *gin.Context) {
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

	sensorTypes, err := c.sensorTypeService.ListSensorTypes(ctx.Request.Context(), page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list sensor types: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, dto.SensorTypeListResponse{
		Success:  true,
		Data:     sensorTypes,
		Page:     page,
		PageSize: limit,
	})
}

// UpdateSensorType handles updating an existing sensor type
func (c *SensorTypeController) UpdateSensorType(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid sensor type ID"})
		return
	}

	var req dto.UpdateSensorTypeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedSensorType, err := c.sensorTypeService.UpdateSensorType(ctx.Request.Context(), id, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, dto.SensorTypeResponse{
		Success: true,
		Data:    updatedSensorType,
	})
}

// DeleteSensorType handles deleting a sensor type
func (c *SensorTypeController) DeleteSensorType(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid sensor type ID"})
		return
	}

	if err := c.sensorTypeService.DeleteSensorType(ctx.Request.Context(), id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}

// GetActiveSensorTypes handles retrieving all active sensor types
func (c *SensorTypeController) GetActiveSensorTypes(ctx *gin.Context) {
	sensorTypes, err := c.sensorTypeService.GetActiveSensorTypes(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get active sensor types: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, dto.SensorTypeListResponse{
		Success: true,
		Data:    sensorTypes,
	})
}

// UpdateSensorTypePartial handles partial updates to a sensor type
func (c *SensorTypeController) UpdateSensorTypePartial(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid sensor type ID"})
		return
	}

	var updateReq map[string]interface{}
	if err := ctx.ShouldBindJSON(&updateReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedSensorType, err := c.sensorTypeService.UpdateSensorTypePartial(ctx.Request.Context(), id, updateReq)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, dto.SensorTypeResponse{
		Success: true,
		Data:    updatedSensorType,
	})
}

// ListAllSensorTypes handles retrieving all sensor types for SuperAdmin (across all tenants)
func (c *SensorTypeController) ListAllSensorTypes(ctx *gin.Context) {
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

	sensorTypes, err := c.sensorTypeService.ListSensorTypes(ctx.Request.Context(), page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list sensor types: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, dto.SensorTypeListResponse{
		Success:  true,
		Data:     sensorTypes,
		Page:     page,
		PageSize: limit,
	})
}
