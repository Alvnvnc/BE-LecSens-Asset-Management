package routes

import (
	"be-lecsens/asset_management/domain-layer/middleware"
	"be-lecsens/asset_management/presentation-layer/controller"

	"github.com/gin-gonic/gin"
)

func LocationRoutes(router *gin.Engine, locationController *controller.LocationController) {
	// Group for location routes
	locationGroup := router.Group("/api/v1/locations")
	{
		// Public routes (requires tenant validation from JWT) - READ ONLY for regular users
		locationGroup.Use(middleware.TenantMiddleware())
		{
			// Get location by ID
			locationGroup.GET("/:id", locationController.GetLocation)
			// List all locations (paginated)
			locationGroup.GET("", locationController.ListLocations)
		}
	}

	// SuperAdmin only routes - use SuperAdmin middleware for role validation
	superAdminGroup := router.Group("/api/v1/superadmin/locations")
	superAdminGroup.Use(middleware.SuperAdminPassthroughMiddleware())
	{
		// List all locations (for SuperAdmin - across all tenants)
		superAdminGroup.GET("", locationController.ListLocations)
		// Get location by ID (for SuperAdmin)
		superAdminGroup.GET("/:id", locationController.GetLocation)
		// Create new location
		superAdminGroup.POST("", locationController.CreateLocation)
		// Update location
		superAdminGroup.PUT("/:id", locationController.UpdateLocation)
		// Delete location
		superAdminGroup.DELETE("/:id", locationController.DeleteLocation)
	}
}
