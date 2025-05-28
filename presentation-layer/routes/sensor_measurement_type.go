package routes

import (
	"be-lecsens/asset_management/domain-layer/middleware"
	"be-lecsens/asset_management/presentation-layer/controller"

	"github.com/gin-gonic/gin"
)

// SetupSensorMeasurementTypeRoutes configures all sensor measurement type related routes
// Combines SuperAdmin routes and regular user routes
func SetupSensorMeasurementTypeRoutes(
	router *gin.Engine,
	controller *controller.SensorMeasurementTypeController,
) {
	// Create a route group for SuperAdmin operations
	superAdminApi := router.Group("/api/v1/superadmin")
	{
		// Apply SuperAdminPassthroughMiddleware for all routes
		superAdminApi.Use(middleware.SuperAdminPassthroughMiddleware())
		{
			// Sensor measurement type routes for SuperAdmin
			sensorMeasurementTypes := superAdminApi.Group("/sensor-measurement-types")
			{
				// Read operations
				sensorMeasurementTypes.GET("/:id", controller.GetSensorMeasurementType)
				sensorMeasurementTypes.GET("", controller.ListAllSensorMeasurementTypes)
				sensorMeasurementTypes.GET("/active", controller.GetActiveSensorMeasurementTypes)
				sensorMeasurementTypes.GET("/by-sensor-type/:sensor_type_id", controller.GetSensorMeasurementTypesBySensorTypeID)

				// Write operations
				sensorMeasurementTypes.POST("", controller.CreateSensorMeasurementType)
				sensorMeasurementTypes.PUT("/:id", controller.UpdateSensorMeasurementType)
				sensorMeasurementTypes.DELETE("/:id", controller.DeleteSensorMeasurementType)
			}
		}
	}

	// Create a route group for regular user operations
	userApi := router.Group("/api/v1/sensor-measurement-types")
	{
		// Apply TenantMiddleware for all routes
		userApi.Use(middleware.TenantMiddleware())
		{
			// Read-only operations for regular users
			userApi.GET("", controller.ListSensorMeasurementTypes)
			userApi.GET("/:id", controller.GetSensorMeasurementType)
			userApi.GET("/active", controller.GetActiveSensorMeasurementTypes)
			userApi.GET("/by-sensor-type/:sensor_type_id", controller.GetSensorMeasurementTypesBySensorTypeID)
		}
	}
}
