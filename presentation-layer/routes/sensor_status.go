package routes

import (
	"be-lecsens/asset_management/domain-layer/middleware"
	"be-lecsens/asset_management/presentation-layer/controller"

	"github.com/gin-gonic/gin"
)

// SetupSensorStatusRoutes configures all sensor status-related routes
func SetupSensorStatusRoutes(router *gin.Engine, sensorStatusController *controller.SensorStatusController) {
	// Group for sensor status routes
	sensorStatusGroup := router.Group("/api/v1/sensor-statuses")
	{
		// Public routes (requires tenant validation from JWT) - READ ONLY for regular users
		sensorStatusGroup.Use(middleware.TenantMiddleware())
		{
			// Get all sensor statuses (for regular users - tenant scoped)
			sensorStatusGroup.GET("", sensorStatusController.GetAllSensorStatuses)
			// Get sensor status by ID
			sensorStatusGroup.GET("/:id", sensorStatusController.GetSensorStatus)
			// Get online sensors
			sensorStatusGroup.GET("/online", sensorStatusController.GetOnlineSensors)
			// Get offline sensors
			sensorStatusGroup.GET("/offline", sensorStatusController.GetOfflineSensors)
			// Get low battery sensors
			sensorStatusGroup.GET("/low-battery", sensorStatusController.GetLowBatterySensors)
			// Get weak signal sensors
			sensorStatusGroup.GET("/weak-signal", sensorStatusController.GetWeakSignalSensors)
			// Get unhealthy sensors
			sensorStatusGroup.GET("/unhealthy", sensorStatusController.GetUnhealthySensors)
			// Get sensor health summary
			sensorStatusGroup.GET("/health-summary", sensorStatusController.GetSensorHealthSummary)
		}

		// Admin routes - use TenantAdmin middleware for role validation
		adminGroup := sensorStatusGroup.Group("")
		adminGroup.Use(middleware.TenantAdminMiddleware())
		{
			// Create new sensor status
			adminGroup.POST("", sensorStatusController.CreateSensorStatus)
			// Update sensor status
			adminGroup.PUT("/:id", sensorStatusController.UpdateSensorStatus)
			// Upsert sensor status (create or update)
			adminGroup.POST("/upsert", sensorStatusController.UpsertSensorStatus)
			// Delete sensor status
			adminGroup.DELETE("/:id", sensorStatusController.DeleteSensorStatus)
		}

		// SuperAdmin only routes - use SuperAdmin middleware for role validation
		superAdminGroup := router.Group("/api/v1/superadmin/sensor-statuses")
		superAdminGroup.Use(middleware.SuperAdminPassthroughMiddleware())
		{
			// Get sensor status by ID (across all tenants)
			superAdminGroup.GET("/:id", sensorStatusController.GetSensorStatus)
			// Create new sensor status (across any tenant)
			superAdminGroup.POST("", sensorStatusController.CreateSensorStatus)
			// Update sensor status (across any tenant)
			superAdminGroup.PUT("/:id", sensorStatusController.UpdateSensorStatus)
			// Upsert sensor status (across any tenant)
			superAdminGroup.POST("/upsert", sensorStatusController.UpsertSensorStatus)
			// Delete sensor status (across any tenant)
			superAdminGroup.DELETE("/:id", sensorStatusController.DeleteSensorStatus)
			// Global online sensors
			superAdminGroup.GET("/online", sensorStatusController.GetOnlineSensors)
			// Global offline sensors
			superAdminGroup.GET("/offline", sensorStatusController.GetOfflineSensors)
			// Global low battery sensors
			superAdminGroup.GET("/low-battery", sensorStatusController.GetLowBatterySensors)
			// Global weak signal sensors
			superAdminGroup.GET("/weak-signal", sensorStatusController.GetWeakSignalSensors)
			// Global unhealthy sensors
			superAdminGroup.GET("/unhealthy", sensorStatusController.GetUnhealthySensors)
			// Global sensor health summary
			superAdminGroup.GET("/health-summary", sensorStatusController.GetSensorHealthSummary)
		}
	}

	// Sensor-specific routes (nested under sensors)
	sensorsGroup := router.Group("/api/v1/sensors")
	{
		// Public routes (requires tenant validation from JWT) - READ ONLY for regular users
		sensorsGroup.Use(middleware.TenantMiddleware())
		{
			// Get current sensor status by sensor ID
			sensorsGroup.GET("/:sensorId/status", sensorStatusController.GetSensorStatusBySensorID)
		}

		// Admin routes - use TenantAdmin middleware for role validation
		adminSensorsGroup := sensorsGroup.Group("")
		adminSensorsGroup.Use(middleware.TenantAdminMiddleware())
		{
			// Update sensor heartbeat
			adminSensorsGroup.PATCH("/:sensorId/heartbeat", sensorStatusController.UpdateHeartbeat)
			// Delete sensor status by sensor ID
			adminSensorsGroup.DELETE("/:sensorId/status", sensorStatusController.DeleteSensorStatusBySensorID)
		}

		// SuperAdmin only routes - use SuperAdmin middleware for role validation
		superAdminSensorsGroup := router.Group("/api/v1/superadmin/sensors")
		superAdminSensorsGroup.Use(middleware.SuperAdminPassthroughMiddleware())
		{
			// Get current sensor status by sensor ID (across all tenants)
			superAdminSensorsGroup.GET("/:sensorId/status", sensorStatusController.GetSensorStatusBySensorID)
			// Update sensor heartbeat (across all tenants)
			superAdminSensorsGroup.PATCH("/:sensorId/heartbeat", sensorStatusController.UpdateHeartbeat)
			// Delete sensor status by sensor ID (across all tenants)
			superAdminSensorsGroup.DELETE("/:sensorId/status", sensorStatusController.DeleteSensorStatusBySensorID)
		}
	}
}
