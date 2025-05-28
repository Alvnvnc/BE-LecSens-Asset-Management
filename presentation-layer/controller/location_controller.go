package controller

import (
	"be-lecsens/asset_management/data-layer/entity"
	"be-lecsens/asset_management/domain-layer/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// LocationController handles location-related HTTP requests
type LocationController struct {
	locationService *service.LocationService
}

// NewLocationController creates a new LocationController
func NewLocationController(locationService *service.LocationService) *LocationController {
	return &LocationController{
		locationService: locationService,
	}
}

// GetLocation handles GET /api/v1/locations/:id
func (c *LocationController) GetLocation(ctx *gin.Context) {
	id := ctx.Param("id")
	locationID, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid location ID format"})
		return
	}

	location, err := c.locationService.GetLocationByID(ctx, locationID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, location)
}

// ListLocations handles GET /api/v1/locations
func (c *LocationController) ListLocations(ctx *gin.Context) {
	// Get pagination parameters
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))

	locations, err := c.locationService.ListLocations(ctx, page, pageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, locations)
}

// UpdateLocation handles PUT /api/v1/locations/:id
// This endpoint requires SuperAdmin access
func (c *LocationController) UpdateLocation(ctx *gin.Context) {
	// Get location ID from URL
	id := ctx.Param("id")
	locationID, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid location ID format"})
		return
	}

	// Parse request body
	var location entity.Location
	if err := ctx.ShouldBindJSON(&location); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ensure the ID in the URL matches the ID in the request body
	location.ID = locationID

	// Update the location
	err = c.locationService.UpdateLocation(ctx, &location)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, location)
}
