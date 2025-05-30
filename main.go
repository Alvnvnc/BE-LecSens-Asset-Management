package main

import (
	"be-lecsens/asset_management/data-layer/cloudinary"
	"be-lecsens/asset_management/data-layer/config"
	"be-lecsens/asset_management/data-layer/migration"
	"be-lecsens/asset_management/data-layer/repository"
	"be-lecsens/asset_management/domain-layer/middleware"
	"be-lecsens/asset_management/domain-layer/service"
	"be-lecsens/asset_management/presentation-layer/controller"
	"be-lecsens/asset_management/presentation-layer/routes"
	"database/sql"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// Load .env file if it exists
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found or error loading .env file. Using environment variables.")
	}

	// Load configuration
	cfg := config.Load()

	// Run database migrations
	err = migration.MigrateDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to run database migrations: %v", err)
	}

	// Set up database connection
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.Name)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test database connection
	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Connected to database successfully")

	// Initialize Cloudinary service
	cloudinaryService, err := cloudinary.NewCloudinaryService(&cloudinary.CloudinaryConfig{
		CloudName: cfg.Cloudinary.CloudName,
		APIKey:    cfg.Cloudinary.APIKey,
		APISecret: cfg.Cloudinary.APISecret,
	})
	if err != nil {
		log.Fatalf("Failed to initialize Cloudinary service: %v", err)
	}

	// Initialize repositories
	assetRepo := repository.NewAssetRepository(db)
	assetTypeRepo := repository.NewAssetTypeRepository(db)
	locationRepo := repository.NewLocationRepository(db)
	assetDocumentRepo := repository.NewAssetDocumentRepository(db)
	assetSensorRepo := repository.NewAssetSensorRepository(db)
	sensorTypeRepo := repository.NewSensorTypeRepository(db)
	sensorMeasurementFieldRepo := repository.NewSensorMeasurementFieldRepository(db)
	sensorMeasurementTypeRepo := repository.NewSensorMeasurementTypeRepository(db)
	iotSensorReadingRepo := repository.NewIoTSensorReadingRepository(db)

	// Initialize services
	log.Println("Initializing services")
	assetService := service.NewAssetService(assetRepo, assetTypeRepo, locationRepo)
	assetTypeService := service.NewAssetTypeService(assetTypeRepo)
	locationService := service.NewLocationService(locationRepo)
	assetDocumentService := service.NewAssetDocumentService(assetDocumentRepo, assetRepo, cloudinaryService)
	assetSensorService := service.NewAssetSensorService(assetSensorRepo, assetRepo)
	sensorTypeService := service.NewSensorTypeService(sensorTypeRepo)
	sensorMeasurementFieldService := service.NewSensorMeasurementFieldService(sensorMeasurementFieldRepo)
	sensorMeasurementTypeService := service.NewSensorMeasurementTypeService(sensorMeasurementTypeRepo)
	iotSensorReadingService := service.NewIoTSensorReadingService(iotSensorReadingRepo, assetSensorRepo, sensorTypeRepo, assetRepo, locationRepo)

	// Initialize controllers
	assetController := controller.NewAssetController(assetService, cfg)
	assetTypeController := controller.NewAssetTypeController(assetTypeService)
	locationController := controller.NewLocationController(locationService)
	assetDocumentController := controller.NewAssetDocumentController(assetDocumentService, cfg)
	assetSensorController := controller.NewAssetSensorController(assetSensorService)
	sensorTypeController := controller.NewSensorTypeController(sensorTypeService, cfg)
	sensorMeasurementFieldController := controller.NewSensorMeasurementFieldController(sensorMeasurementFieldService)
	sensorMeasurementTypeController := controller.NewSensorMeasurementTypeController(sensorMeasurementTypeService, cfg)
	iotSensorReadingController := controller.NewIoTSensorReadingController(iotSensorReadingService)

	// Initialize JWT config
	jwtConfig := middleware.JWTConfig{
		SecretKey: cfg.JWT.SecretKey,
	}

	// Set up Gin router
	router := gin.Default()

	// Add custom middleware
	router.Use(middleware.LoggingMiddleware())

	// Configure routes
	routes.SetupRoutes(
		router,
		assetController,
		assetTypeController,
		locationController,
		assetDocumentController,
		assetSensorController,
		sensorTypeController,
		sensorMeasurementFieldController,
		sensorMeasurementTypeController,
		iotSensorReadingController,
		jwtConfig,
	)

	// Start the server
	log.Printf("Server running on port %s", cfg.Server.Port)
	router.Run(":" + cfg.Server.Port)
}
