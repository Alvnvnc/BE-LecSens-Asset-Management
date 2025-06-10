package routes

import (
	"be-lecsens/asset_management/domain-layer/middleware"
	"be-lecsens/asset_management/presentation-layer/controller"

	"github.com/gin-gonic/gin"
)

// SetupSensorLogsRoutes configures all sensor logs-related routes
func SetupSensorLogsRoutes(router *gin.Engine, sensorLogsController *controller.SensorLogsController) {
	// Group for sensor logs routes
	sensorLogsGroup := router.Group("/api/v1/sensor-logs")
	{
		// Public routes (requires tenant validation from JWT) - READ ONLY for regular users
		sensorLogsGroup.Use(middleware.TenantMiddleware())
		{
			// List sensor logs with filtering and pagination
			sensorLogsGroup.GET("", sensorLogsController.ListSensorLogs)
			// Get sensor log by ID
			sensorLogsGroup.GET("/:id", sensorLogsController.GetSensorLog)
			// Get sensor logs by sensor ID
			sensorLogsGroup.GET("/sensor/:sensorId", sensorLogsController.GetSensorLogsBySensorID)
			// Get connection history analytics (read-only analytics)
			sensorLogsGroup.GET("/analytics/connection-history", sensorLogsController.GetConnectionHistory)
			// Get log statistics (read-only analytics)
			sensorLogsGroup.GET("/analytics/statistics", sensorLogsController.GetLogStatistics)
			// Get log analytics (read-only analytics)
			sensorLogsGroup.GET("/analytics/insights", sensorLogsController.GetLogAnalytics)
		}

		// Admin routes - use TenantAdmin middleware for role validation
		adminGroup := sensorLogsGroup.Group("")
		adminGroup.Use(middleware.TenantAdminMiddleware())
		{
			// Create new sensor log
			adminGroup.POST("", sensorLogsController.CreateSensorLog)
			// Update sensor log
			adminGroup.PUT("/:id", sensorLogsController.UpdateSensorLog)
			// Delete sensor log
			adminGroup.DELETE("/:id", sensorLogsController.DeleteSensorLog)
			// Delete sensor logs by sensor ID
			adminGroup.DELETE("/sensor/:sensorId", sensorLogsController.DeleteSensorLogsBySensorID)
			// Cleanup old logs (maintenance operation)
			adminGroup.DELETE("/cleanup", sensorLogsController.CleanupOldLogs)
		}

		// SuperAdmin only routes - use SuperAdmin middleware for role validation
		superAdminGroup := router.Group("/api/v1/superadmin/sensor-logs")
		superAdminGroup.Use(middleware.SuperAdminPassthroughMiddleware())
		{
			// List all sensor logs (across all tenants)
			superAdminGroup.GET("", sensorLogsController.ListSensorLogs)
			// Get sensor log by ID (with complete information)
			superAdminGroup.GET("/:id", sensorLogsController.GetSensorLog)
			// Create new sensor log (across any tenant)
			superAdminGroup.POST("", sensorLogsController.CreateSensorLog)
			// Update sensor log (across any tenant)
			superAdminGroup.PUT("/:id", sensorLogsController.UpdateSensorLog)
			// Delete sensor log (across any tenant)
			superAdminGroup.DELETE("/:id", sensorLogsController.DeleteSensorLog)
			// Get sensor logs by sensor ID (across any tenant)
			superAdminGroup.GET("/sensor/:sensorId", sensorLogsController.GetSensorLogsBySensorID)
			// Delete sensor logs by sensor ID (across any tenant)
			superAdminGroup.DELETE("/sensor/:sensorId", sensorLogsController.DeleteSensorLogsBySensorID)
			// Global connection history analytics
			superAdminGroup.GET("/analytics/connection-history", sensorLogsController.GetConnectionHistory)
			// Global log statistics
			superAdminGroup.GET("/analytics/statistics", sensorLogsController.GetLogStatistics)
			// Global log analytics
			superAdminGroup.GET("/analytics/insights", sensorLogsController.GetLogAnalytics)
			// Global cleanup of old logs
			superAdminGroup.DELETE("/cleanup", sensorLogsController.CleanupOldLogs)
		}
	}
}
