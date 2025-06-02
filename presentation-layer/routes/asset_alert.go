package routes

import (
	"be-lecsens/asset_management/domain-layer/middleware"
	"be-lecsens/asset_management/presentation-layer/controller"

	"github.com/gin-gonic/gin"
)

func SetupAssetAlertRoutes(router *gin.Engine, assetAlertController *controller.AssetAlertController) {
	// Group for asset alert routes
	assetAlertGroup := router.Group("/api/v1/asset-alerts")
	{
		// Public routes (requires tenant validation from JWT) - READ ONLY for regular users
		assetAlertGroup.Use(middleware.TenantMiddleware())
		{
			// List all asset alerts (paginated)
			assetAlertGroup.GET("", assetAlertController.GetAlerts)
			// Get asset alert by ID
			assetAlertGroup.GET("/:id", assetAlertController.GetAlert)
			// Statistics endpoint
			assetAlertGroup.GET("/statistics", assetAlertController.GetAlertStatistics)
		}
	}

	// Asset-specific alert endpoints (public)
	assetsGroup := router.Group("/api/v1/assets")
	assetsGroup.Use(middleware.TenantMiddleware())
	{
		assetsGroup.GET("/:asset_id/alerts", assetAlertController.GetAssetAlerts)
	}

	// Sensor-specific alert endpoints (public)
	sensorsGroup := router.Group("/api/v1/sensors")
	sensorsGroup.Use(middleware.TenantMiddleware())
	{
		sensorsGroup.GET("/:sensor_id/alerts", assetAlertController.GetSensorAlerts)
	}

	// SuperAdmin only routes - use SuperAdmin middleware for role validation
	superAdminGroup := router.Group("/api/v1/superadmin/asset-alerts")
	superAdminGroup.Use(middleware.SuperAdminPassthroughMiddleware())
	{
		// List all asset alerts (for SuperAdmin - across all tenants)
		superAdminGroup.GET("", assetAlertController.GetAlerts)
		// Get asset alert by ID (for SuperAdmin)
		superAdminGroup.GET("/:id", assetAlertController.GetAlert)
		// Alert resolution endpoints
		superAdminGroup.PUT("/:id/resolve", assetAlertController.ResolveAlert)
		superAdminGroup.PUT("/bulk-resolve", assetAlertController.BulkResolveAlerts)
		// Statistics endpoint (for SuperAdmin)
		superAdminGroup.GET("/statistics", assetAlertController.GetAlertStatistics)
	}

	// SuperAdmin asset-specific alert endpoints
	superAdminAssetsGroup := router.Group("/api/v1/superadmin/assets")
	superAdminAssetsGroup.Use(middleware.SuperAdminPassthroughMiddleware())
	{
		superAdminAssetsGroup.GET("/:asset_id/alerts", assetAlertController.GetAssetAlerts)
	}

	// SuperAdmin sensor-specific alert endpoints
	superAdminSensorsGroup := router.Group("/api/v1/superadmin/sensors")
	superAdminSensorsGroup.Use(middleware.SuperAdminPassthroughMiddleware())
	{
		superAdminSensorsGroup.GET("/:sensor_id/alerts", assetAlertController.GetSensorAlerts)
	}
}
