package routes

import (
	"be-lecsens/asset_management/domain-layer/middleware"
	"be-lecsens/asset_management/presentation-layer/controller"

	"github.com/gin-gonic/gin"
)

// SetupAssetTypeRoutes configures all asset type related routes
// All asset type operations require SuperAdmin privileges
func SetupAssetTypeRoutes(
	router *gin.Engine,
	controller *controller.AssetTypeController,
) {
	// Create a route group for SuperAdmin operations
	superAdminApi := router.Group("/api/v1/superadmin")

	// Apply SuperAdminPassthroughMiddleware for all routes
	superAdminApi.Use(middleware.SuperAdminPassthroughMiddleware())
	{
		// Asset type routes for SuperAdmin
		assetTypes := superAdminApi.Group("/asset-types")
		{
			// Read operations
			assetTypes.GET("/:id", controller.GetByID)
			assetTypes.GET("", controller.List)

			// Write operations
			assetTypes.POST("", controller.Create)
			assetTypes.PUT("/:id", controller.Update)
			assetTypes.DELETE("/:id", controller.Delete)
		}
	}
}
