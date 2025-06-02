package controller

import (
	"be-lecsens/asset_management/domain-layer/service"
	"be-lecsens/asset_management/helpers/dto"
	"log"
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

func (c *LocationController) CreateLocation(ctx *gin.Context) {
	log.Printf("Location Controller: CreateLocation called")

	// Parse request body
	var req dto.CreateLocationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Location Controller: Invalid request body: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request format",
			"error":   err.Error(),
		})
		return
	}

	log.Printf("Location Controller: Create request: %+v", req)

	// Create location
	response, err := c.locationService.CreateLocation(ctx, req)
	if err != nil {
		log.Printf("Location Controller: Error creating location: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to create location",
			"error":   err.Error(),
		})
		return
	}

	log.Printf("Location Controller: Successfully created location: %s", response.ID)
	ctx.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Location created successfully",
		"data":    response,
	})
}

func (c *LocationController) UpdateLocation(ctx *gin.Context) {
	log.Printf("Location Controller: UpdateLocation called")

	// Get location ID from URL parameter
	idParam := ctx.Param("id")
	locationID, err := uuid.Parse(idParam)
	if err != nil {
		log.Printf("Location Controller: Invalid location ID: %s", idParam)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid location ID format",
			"error":   err.Error(),
		})
		return
	}

	// Parse request body
	var req dto.UpdateLocationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Location Controller: Invalid request body: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request format",
			"error":   err.Error(),
		})
		return
	}

	log.Printf("Location Controller: Update request for location %s: %+v", locationID, req)

	// Update location
	response, err := c.locationService.UpdateLocation(ctx, locationID, req)
	if err != nil {
		log.Printf("Location Controller: Error updating location: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to update location",
			"error":   err.Error(),
		})
		return
	}

	log.Printf("Location Controller: Successfully updated location: %s", locationID)
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Location updated successfully",
		"data":    response,
	})
}

func (c *LocationController) DeleteLocation(ctx *gin.Context) {
	log.Printf("Location Controller: DeleteLocation called")

	// Get location ID from URL parameter
	idParam := ctx.Param("id")
	locationID, err := uuid.Parse(idParam)
	if err != nil {
		log.Printf("Location Controller: Invalid location ID: %s", idParam)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid location ID format",
			"error":   err.Error(),
		})
		return
	}

	log.Printf("Location Controller: Delete request for location %s", locationID)

	// Delete location
	err = c.locationService.DeleteLocation(ctx, locationID)
	if err != nil {
		log.Printf("Location Controller: Error deleting location: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to delete location",
			"error":   err.Error(),
		})
		return
	}

	log.Printf("Location Controller: Successfully deleted location: %s", locationID)
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Location deleted successfully",
	})
}
