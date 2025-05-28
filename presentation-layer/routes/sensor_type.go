package routes

import (
	"be-lecsens/asset_management/domain-layer/middleware"
	"be-lecsens/asset_management/presentation-layer/controller"

	"github.com/gin-gonic/gin"
)

// SetupSensorTypeRoutes configures all sensor type related routes
// Combines SuperAdmin routes and regular user routes
func SetupSensorTypeRoutes(
	router *gin.Engine,
	controller *controller.SensorTypeController,
) {
	// Create a route group for SuperAdmin operations
	superAdminApi := router.Group("/api/v1/superadmin")
	{
		// Apply SuperAdminPassthroughMiddleware for all routes
		superAdminApi.Use(middleware.SuperAdminPassthroughMiddleware())
		{
			// Sensor type routes for SuperAdmin
			sensorTypes := superAdminApi.Group("/sensor-types")
			{
				// Read operations
				sensorTypes.GET("/:id", controller.GetSensorType)
				sensorTypes.GET("", controller.ListAllSensorTypes)
				sensorTypes.GET("/active", controller.GetActiveSensorTypes)

				// Write operations
				sensorTypes.POST("", controller.CreateSensorType)
				sensorTypes.PUT("/:id", controller.UpdateSensorType)
				sensorTypes.PATCH("/:id", controller.UpdateSensorTypePartial)
				sensorTypes.DELETE("/:id", controller.DeleteSensorType)
			}
		}
	}

	// Create a route group for regular user operations
	userApi := router.Group("/api/v1/sensor-types")
	{
		// Apply TenantMiddleware for all routes
		userApi.Use(middleware.TenantMiddleware())
		{
			// Read-only operations for regular users
			userApi.GET("", controller.ListSensorTypes)
			userApi.GET("/:id", controller.GetSensorType)
			userApi.GET("/active", controller.GetActiveSensorTypes)
		}
	}
}
