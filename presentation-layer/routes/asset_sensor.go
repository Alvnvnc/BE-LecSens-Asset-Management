package routes

import (
	"be-lecsens/asset_management/domain-layer/middleware"
	"be-lecsens/asset_management/presentation-layer/controller"

	"github.com/gin-gonic/gin"
)

// SetupAssetSensorRoutes configures all asset sensor-related routes
func SetupAssetSensorRoutes(router *gin.Engine, assetSensorController *controller.AssetSensorController) {
	// Group for asset sensor routes
	assetSensorGroup := router.Group("/api/v1/asset-sensors")
	{
		// Public routes (requires tenant validation from JWT) - READ ONLY for regular users
		assetSensorGroup.Use(middleware.TenantMiddleware())
		{
			// List all asset sensors (paginated)
			assetSensorGroup.GET("", assetSensorController.ListAssetSensors)
			// Get asset sensor by ID
			assetSensorGroup.GET("/:id", assetSensorController.GetAssetSensor)
			// Get all sensors for a specific asset
			assetSensorGroup.GET("/asset/:asset_id", assetSensorController.GetAssetSensors)
			// Get active sensors
			assetSensorGroup.GET("/active", assetSensorController.GetActiveSensors)
			// Get sensors by status
			assetSensorGroup.GET("/status/:status", assetSensorController.GetSensorsByStatus)
		}
	}

	// SuperAdmin only routes - use SuperAdmin middleware for role validation
	superAdminGroup := router.Group("/api/v1/superadmin/asset-sensors")
	superAdminGroup.Use(middleware.SuperAdminPassthroughMiddleware())
	{
		// List all asset sensors (for SuperAdmin - across all tenants)
		superAdminGroup.GET("", assetSensorController.ListAssetSensors)
		// Get asset sensor detail by ID (with complete information)
		superAdminGroup.GET("/:id", assetSensorController.GetAssetSensor)
		// Create new asset sensor
		superAdminGroup.POST("", assetSensorController.CreateAssetSensor)
		// Update asset sensor
		superAdminGroup.PUT("/:id", assetSensorController.UpdateAssetSensor)
		// Delete asset sensor
		superAdminGroup.DELETE("/:id", assetSensorController.DeleteAssetSensor)
		// Delete all sensors for an asset
		superAdminGroup.DELETE("/asset/:asset_id", assetSensorController.DeleteAssetSensors)
		// Update sensor reading
		superAdminGroup.PUT("/:id/reading", assetSensorController.UpdateSensorReading)
		// Get all sensors for a specific asset (with detailed information)
		superAdminGroup.GET("/asset/:asset_id", assetSensorController.GetAssetSensorsDetailed)
		// Get active sensors
		superAdminGroup.GET("/active", assetSensorController.GetActiveSensors)
		// Get sensors by status
		superAdminGroup.GET("/status/:status", assetSensorController.GetSensorsByStatus)
	}
}
