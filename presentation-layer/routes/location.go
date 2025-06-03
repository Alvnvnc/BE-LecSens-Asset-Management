package routes

import (
	"be-lecsens/asset_management/domain-layer/middleware"
	"be-lecsens/asset_management/presentation-layer/controller"

	"github.com/gin-gonic/gin"
)

// SetupLocationRoutes configures all location related routes
func SetupLocationRoutes(
	router *gin.Engine,
	locationController *controller.LocationController,
) {
	// Public routes for reading locations
	api := router.Group("/api/v1")
	{
		locations := api.Group("/locations")
		{
			// Read operations - public access
			locations.GET("", locationController.ListLocations)
			locations.GET("/:id", locationController.GetLocation)
		}
	}

	// Protected routes for managing locations (SuperAdmin only)
	superAdminApi := router.Group("/api/v1/superadmin")
	superAdminApi.Use(middleware.SuperAdminPassthroughMiddleware())
	{
		locations := superAdminApi.Group("/locations")
		{
			// CRUD operations - SuperAdmin only
			locations.POST("", locationController.CreateLocation)
			locations.PUT("/:id", locationController.UpdateLocation)
			locations.DELETE("/:id", locationController.DeleteLocation)
		}
	}
}
