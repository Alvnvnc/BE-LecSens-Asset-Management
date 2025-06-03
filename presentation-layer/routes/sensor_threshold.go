package routes

import (
	"be-lecsens/asset_management/domain-layer/middleware"
	"be-lecsens/asset_management/presentation-layer/controller"

	"github.com/gin-gonic/gin"
)

func SetupSensorThresholdRoutes(router *gin.Engine, sensorThresholdController *controller.SensorThresholdController) {
	// Create a route group for SuperAdmin operations (following sensor_type.go pattern)
	superAdminApi := router.Group("/api/v1/superadmin")
	{
		// Apply SuperAdminPassthroughMiddleware for all routes
		superAdminApi.Use(middleware.SuperAdminPassthroughMiddleware())
		{
			// Sensor threshold routes for SuperAdmin
			sensorThresholds := superAdminApi.Group("/sensor-thresholds")
			{
				// Read operations
				sensorThresholds.GET("/:id", sensorThresholdController.GetThreshold)
				sensorThresholds.GET("", sensorThresholdController.GetThresholds)
				sensorThresholds.GET("/statistics", sensorThresholdController.GetThresholdStatistics)

				// Write operations
				sensorThresholds.POST("", sensorThresholdController.CreateThreshold)
				sensorThresholds.PUT("/:id", sensorThresholdController.UpdateThreshold)
				sensorThresholds.DELETE("/:id", sensorThresholdController.DeleteThreshold)
				sensorThresholds.POST("/check", sensorThresholdController.CheckThresholds)
			}

			// SuperAdmin sensor-specific threshold endpoints
			sensors := superAdminApi.Group("/sensors")
			{
				sensors.GET("/:sensor_id/thresholds", sensorThresholdController.GetSensorThresholds)
			}
		}
	}

	// Create a route group for regular user operations
	userApi := router.Group("/api/v1/sensor-thresholds")
	{
		// Apply TenantMiddleware for all routes
		userApi.Use(middleware.TenantMiddleware())
		{
			// Read-only operations for regular users
			userApi.GET("", sensorThresholdController.GetThresholds)
			userApi.GET("/:id", sensorThresholdController.GetThreshold)
			userApi.POST("/check", sensorThresholdController.CheckThresholds)
			userApi.GET("/statistics", sensorThresholdController.GetThresholdStatistics)
		}
	}

	// Sensor-specific threshold endpoints for regular users
	sensorsGroup := router.Group("/api/v1/sensors")
	sensorsGroup.Use(middleware.TenantMiddleware())
	{
		sensorsGroup.GET("/:sensor_id/thresholds", sensorThresholdController.GetSensorThresholds)
	}
}
