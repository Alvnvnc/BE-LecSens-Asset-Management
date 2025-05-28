package routes

import (
	"be-lecsens/asset_management/domain-layer/middleware"
	"be-lecsens/asset_management/presentation-layer/controller"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all the routes for the application
func SetupRoutes(
	router *gin.Engine,
	assetController *controller.AssetController,
	assetTypeController *controller.AssetTypeController,
	locationController *controller.LocationController,
	assetDocumentController *controller.AssetDocumentController,
	jwtConfig middleware.JWTConfig,
) {

	// Public routes (no tenant required)
	public := router.Group("/api/v1")
	{
		public.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})
	}

	// Setup Asset Type routes
	SetupAssetTypeRoutes(router, assetTypeController)

	// Setup Location routes
	SetupLocationRoutes(router, locationController)

	// Setup Asset routes
	SetupAssetRoutes(router, assetController)

	// Setup Asset Document routes
	SetupAssetDocumentRoutes(router, assetDocumentController)

}
