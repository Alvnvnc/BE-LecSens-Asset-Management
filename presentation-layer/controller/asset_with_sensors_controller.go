package controller

import (
	"net/http"
	"strconv"

	"be-lecsens/asset_management/domain-layer/service"
	"be-lecsens/asset_management/helpers/common"
	"be-lecsens/asset_management/helpers/dto"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AssetWithSensorsController handles HTTP requests for asset with sensors operations
type AssetWithSensorsController struct {
	assetWithSensorsService *service.AssetWithSensorsService
}

// NewAssetWithSensorsController creates a new instance of AssetWithSensorsController
func NewAssetWithSensorsController(assetWithSensorsService *service.AssetWithSensorsService) *AssetWithSensorsController {
	return &AssetWithSensorsController{
		assetWithSensorsService: assetWithSensorsService,
	}
}

// CreateAssetWithSensors creates a new asset with associated sensors
// @Summary Create asset with sensors
// @Description Create a new asset and automatically generate associated sensors based on provided sensor types
// @Tags Asset with Sensors
// @Accept json
// @Produce json
// @Param request body dto.CreateAssetWithSensorsRequest true "Asset with sensors creation request"
// @Success 201 {object} dto.AssetWithSensorsResponse
// @Failure 400 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /api/assets-with-sensors [post]
func (c *AssetWithSensorsController) CreateAssetWithSensors(ctx *gin.Context) {
	var req dto.CreateAssetWithSensorsRequest

	// Bind JSON request
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request format",
			"error":   err.Error(),
		})
		return
	}

	// Validate required fields
	if req.Name == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Asset name is required",
		})
		return
	}

	if req.AssetTypeID == uuid.Nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Asset type ID is required",
		})
		return
	}

	if req.LocationID == uuid.Nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Location ID is required",
		})
		return
	}

	if len(req.SensorTypes) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "At least one sensor type is required",
		})
		return
	}

	// Validate each sensor type request
	for i, sensorReq := range req.SensorTypes {
		if sensorReq.SensorTypeID == uuid.Nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Sensor type ID is required for all sensors",
				"error":   "sensor_types[" + strconv.Itoa(i) + "].sensor_type_id is required",
			})
			return
		}
	}

	// Call service to create asset with sensors
	response, err := c.assetWithSensorsService.CreateAssetWithSensors(ctx.Request.Context(), &req)
	if err != nil {
		// Handle different types of errors
		if validationErr, ok := err.(*common.ValidationError); ok {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": validationErr.Error(),
				"error":   validationErr.Error(),
			})
			return
		}

		if notFoundErr, ok := err.(*common.NotFoundError); ok {
			ctx.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": notFoundErr.Error(),
				"error":   notFoundErr.Error(),
			})
			return
		}

		// Generic server error
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to create asset with sensors",
			"error":   err.Error(),
		})
		return
	}

	// Return success response
	ctx.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Asset with sensors created successfully",
		"data":    response,
	})
}

// GetAssetWithSensors retrieves an asset along with all its sensors
// @Summary Get asset with sensors
// @Description Retrieve an asset and all its associated sensors
// @Tags Asset with Sensors
// @Accept json
// @Produce json
// @Param id path string true "Asset ID"
// @Success 200 {object} dto.AssetWithSensorsResponse
// @Failure 400 {object} common.ErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /api/assets-with-sensors/{id} [get]
func (c *AssetWithSensorsController) GetAssetWithSensors(ctx *gin.Context) {
	// Parse asset ID from URL parameter
	assetIDStr := ctx.Param("id")
	assetID, err := uuid.Parse(assetIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid asset ID format",
			"error":   err.Error(),
		})
		return
	}

	// Call service to get asset with sensors
	response, err := c.assetWithSensorsService.GetAssetWithSensors(ctx.Request.Context(), assetID)
	if err != nil {
		// Handle different types of errors
		if notFoundErr, ok := err.(*common.NotFoundError); ok {
			ctx.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": notFoundErr.Error(),
				"error":   notFoundErr.Error(),
			})
			return
		}

		// Generic server error
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to retrieve asset with sensors",
			"error":   err.Error(),
		})
		return
	}

	// Return success response
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Asset with sensors retrieved successfully",
		"data":    response,
	})
}

// UpdateAssetWithSensors updates an existing asset and its sensors
// @Summary Update asset with sensors
// @Description Update an existing asset and manage its associated sensors
// @Tags Asset with Sensors
// @Accept json
// @Produce json
// @Param id path string true "Asset ID"
// @Param request body dto.UpdateAssetWithSensorsRequest true "Asset with sensors update request"
// @Success 200 {object} dto.AssetWithSensorsResponse
// @Failure 400 {object} common.ErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /api/assets-with-sensors/{id} [put]
func (c *AssetWithSensorsController) UpdateAssetWithSensors(ctx *gin.Context) {
	// Parse asset ID from URL parameter
	assetIDStr := ctx.Param("id")
	assetID, err := uuid.Parse(assetIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid asset ID format",
			"error":   err.Error(),
		})
		return
	}

	var req dto.UpdateAssetWithSensorsRequest

	// Bind JSON request
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request format",
			"error":   err.Error(),
		})
		return
	}

	// Call service to update asset with sensors
	response, err := c.assetWithSensorsService.UpdateAssetWithSensors(ctx.Request.Context(), assetID, &req)
	if err != nil {
		// Handle different types of errors
		if validationErr, ok := err.(*common.ValidationError); ok {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": validationErr.Error(),
				"error":   validationErr.Error(),
			})
			return
		}

		if notFoundErr, ok := err.(*common.NotFoundError); ok {
			ctx.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": notFoundErr.Error(),
				"error":   notFoundErr.Error(),
			})
			return
		}

		// Generic server error
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to update asset with sensors",
			"error":   err.Error(),
		})
		return
	}

	// Return success response
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Asset with sensors updated successfully",
		"data":    response,
	})
}

// DeleteAssetWithSensors deletes an asset and all its sensors
// @Summary Delete asset with sensors
// @Description Delete an asset and all its associated sensors
// @Tags Asset with Sensors
// @Accept json
// @Produce json
// @Param id path string true "Asset ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} common.ErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /api/assets-with-sensors/{id} [delete]
func (c *AssetWithSensorsController) DeleteAssetWithSensors(ctx *gin.Context) {
	// Parse asset ID from URL parameter
	assetIDStr := ctx.Param("id")
	assetID, err := uuid.Parse(assetIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid asset ID format",
			"error":   err.Error(),
		})
		return
	}

	// Call service to delete asset with sensors
	err = c.assetWithSensorsService.DeleteAssetWithSensors(ctx.Request.Context(), assetID)
	if err != nil {
		// Handle different types of errors
		if notFoundErr, ok := err.(*common.NotFoundError); ok {
			ctx.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": notFoundErr.Error(),
				"error":   notFoundErr.Error(),
			})
			return
		}

		// Generic server error
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to delete asset with sensors",
			"error":   err.Error(),
		})
		return
	}

	// Return success response
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Asset with sensors deleted successfully",
	})
}

// ListAssetsWithSensors retrieves a paginated list of assets with their sensors
// @Summary List assets with sensors
// @Description Retrieve a paginated list of assets with their associated sensors
// @Tags Asset with Sensors
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Number of items per page" default(10)
// @Success 200 {object} dto.AssetWithSensorsListResponse
// @Failure 400 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /api/assets-with-sensors [get]
func (c *AssetWithSensorsController) ListAssetsWithSensors(ctx *gin.Context) {
	// Parse pagination parameters
	page := 1
	pageSize := 10

	if pageStr := ctx.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if pageSizeStr := ctx.Query("pageSize"); pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 && ps <= 100 {
			pageSize = ps
		}
	}

	// Call service to list assets with sensors
	response, err := c.assetWithSensorsService.ListAssetsWithSensors(ctx.Request.Context(), page, pageSize)
	if err != nil {
		// Generic server error
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to list assets with sensors",
			"error":   err.Error(),
		})
		return
	}

	// Return success response
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Assets with sensors retrieved successfully",
		"data":    response,
	})
}

// BulkCreateAssetsWithSensors creates multiple assets with their sensors in a single operation
// @Summary Bulk create assets with sensors
// @Description Create multiple assets with their sensors in a single atomic operation
// @Tags Asset with Sensors
// @Accept json
// @Produce json
// @Param requests body []dto.CreateAssetWithSensorsRequest true "Array of asset with sensors creation requests"
// @Success 201 {object} dto.BulkCreateAssetWithSensorsResponse
// @Failure 400 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /api/assets-with-sensors/bulk [post]
func (c *AssetWithSensorsController) BulkCreateAssetsWithSensors(ctx *gin.Context) {
	var requests []*dto.CreateAssetWithSensorsRequest

	// Bind JSON request
	if err := ctx.ShouldBindJSON(&requests); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request format",
			"error":   err.Error(),
		})
		return
	}

	// Validate that we have at least one request
	if len(requests) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "At least one asset creation request is required",
		})
		return
	}

	// Call service to bulk create assets with sensors
	response, err := c.assetWithSensorsService.BulkCreateAssetsWithSensors(ctx.Request.Context(), requests)
	if err != nil {
		// For bulk operations, we might have partial success
		// The response will contain both successful and failed creations
		ctx.JSON(http.StatusMultiStatus, gin.H{
			"success": false,
			"message": "Bulk creation completed with some errors",
			"data":    response,
			"error":   err.Error(),
		})
		return
	}

	// Return success response
	ctx.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Bulk asset with sensors creation completed successfully",
		"data":    response,
	})
}
