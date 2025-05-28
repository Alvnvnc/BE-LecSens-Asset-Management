package routes

import (
	"be-lecsens/asset_management/domain-layer/middleware"
	"be-lecsens/asset_management/presentation-layer/controller"

	"github.com/gin-gonic/gin"
)

// SetupSensorMeasurementFieldRoutes sets up all routes for sensor measurement fields
func SetupSensorMeasurementFieldRoutes(router *gin.Engine, controller *controller.SensorMeasurementFieldController) {
	// SuperAdmin routes
	superAdminApi := router.Group("/api/v1/superadmin")
	{
		// Apply SuperAdminPassthroughMiddleware for all routes
		superAdminApi.Use(middleware.SuperAdminPassthroughMiddleware())
		{
			// Sensor measurement field routes for SuperAdmin
			sensorMeasurementFields := superAdminApi.Group("/sensor-measurement-fields")
			{
				// Create new field
				sensorMeasurementFields.POST("", controller.Create)

				// Get all fields
				sensorMeasurementFields.GET("", controller.GetAll)

				// Get field by ID
				sensorMeasurementFields.GET("/:id", controller.GetByID)

				// Update field
				sensorMeasurementFields.PUT("/:id", controller.Update)

				// Delete field
				sensorMeasurementFields.DELETE("/:id", controller.Delete)
			}

			// Routes for fields by measurement type
			measurementTypeFields := superAdminApi.Group("/measurement-type-fields")
			{
				// Get all fields for a measurement type
				measurementTypeFields.GET("/:measurement_type_id", controller.GetByMeasurementTypeID)

				// Get required fields for a measurement type
				measurementTypeFields.GET("/:measurement_type_id/required", controller.GetRequiredFields)
			}
		}
	}

	// Regular user routes
	userApi := router.Group("/api/v1/sensor-measurement-fields")
	{
		// Apply TenantMiddleware for all routes
		userApi.Use(middleware.TenantMiddleware())
		{
			// Get all fields
			userApi.GET("", controller.GetAll)

			// Get field by ID
			userApi.GET("/:id", controller.GetByID)
		}
	}

	// Routes for fields by measurement type (user)
	userMeasurementTypeApi := router.Group("/api/v1/measurement-type-fields")
	{
		// Apply TenantMiddleware for all routes
		userMeasurementTypeApi.Use(middleware.TenantMiddleware())
		{
			// Get all fields for a measurement type
			userMeasurementTypeApi.GET("/:measurement_type_id", controller.GetByMeasurementTypeID)

			// Get required fields for a measurement type
			userMeasurementTypeApi.GET("/:measurement_type_id/required", controller.GetRequiredFields)
		}
	}
}
