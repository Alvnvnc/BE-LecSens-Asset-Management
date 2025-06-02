package routes

import (
	"be-lecsens/asset_management/domain-layer/middleware"
	"be-lecsens/asset_management/presentation-layer/controller"

	"github.com/gin-gonic/gin"
)

func SetupSensorThresholdRoutes(router *gin.Engine, sensorThresholdController *controller.SensorThresholdController) {
	// Group for sensor threshold routes
	sensorThresholdGroup := router.Group("/api/v1/sensor-thresholds")
	{
		// Public routes (requires tenant validation from JWT) - READ ONLY for regular users
		sensorThresholdGroup.Use(middleware.TenantMiddleware())
		{
			// List all sensor thresholds (paginated)
			sensorThresholdGroup.GET("", sensorThresholdController.GetThresholds)
			// Get sensor threshold by ID
			sensorThresholdGroup.GET("/:id", sensorThresholdController.GetThreshold)
			// Threshold checking endpoint
			sensorThresholdGroup.POST("/check", sensorThresholdController.CheckThresholds)
			// Statistics endpoint
			sensorThresholdGroup.GET("/statistics", sensorThresholdController.GetThresholdStatistics)
		}
	}

	// Sensor-specific threshold endpoints (public)
	sensorsGroup := router.Group("/api/v1/sensors")
	sensorsGroup.Use(middleware.TenantMiddleware())
	{
		sensorsGroup.GET("/:sensor_id/thresholds", sensorThresholdController.GetSensorThresholds)
	}

	// SuperAdmin only routes - use SuperAdmin middleware for role validation
	superAdminGroup := router.Group("/api/v1/superadmin/sensor-thresholds")
	superAdminGroup.Use(middleware.SuperAdminPassthroughMiddleware())
	{
		// List all sensor thresholds (for SuperAdmin - across all tenants)
		superAdminGroup.GET("", sensorThresholdController.GetThresholds)
		// Get sensor threshold by ID (for SuperAdmin)
		superAdminGroup.GET("/:id", sensorThresholdController.GetThreshold)
		// Create new sensor threshold
		superAdminGroup.POST("", sensorThresholdController.CreateThreshold)
		// Update sensor threshold
		superAdminGroup.PUT("/:id", sensorThresholdController.UpdateThreshold)
		// Delete sensor threshold
		superAdminGroup.DELETE("/:id", sensorThresholdController.DeleteThreshold)
		// Threshold checking endpoint (for SuperAdmin)
		superAdminGroup.POST("/check", sensorThresholdController.CheckThresholds)
		// Statistics endpoint (for SuperAdmin)
		superAdminGroup.GET("/statistics", sensorThresholdController.GetThresholdStatistics)
	}

	// SuperAdmin sensor-specific threshold endpoints
	superAdminSensorsGroup := router.Group("/api/v1/superadmin/sensors")
	superAdminSensorsGroup.Use(middleware.SuperAdminPassthroughMiddleware())
	{
		superAdminSensorsGroup.GET("/:sensor_id/thresholds", sensorThresholdController.GetSensorThresholds)
	}
}
