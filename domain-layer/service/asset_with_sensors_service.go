package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"be-lecsens/asset_management/data-layer/entity"
	"be-lecsens/asset_management/data-layer/repository"
	"be-lecsens/asset_management/helpers/common"
	"be-lecsens/asset_management/helpers/dto"

	"github.com/google/uuid"
)

// AssetWithSensorsService handles business logic for creating assets with sensors
type AssetWithSensorsService struct {
	assetRepo       repository.AssetRepository
	assetSensorRepo repository.AssetSensorRepository
	assetTypeRepo   *repository.AssetTypeRepository
	locationRepo    *repository.LocationRepository
	sensorTypeRepo  *repository.SensorTypeRepository
}

// NewAssetWithSensorsService creates a new instance of AssetWithSensorsService
func NewAssetWithSensorsService(
	assetRepo repository.AssetRepository,
	assetSensorRepo repository.AssetSensorRepository,
	assetTypeRepo *repository.AssetTypeRepository,
	locationRepo *repository.LocationRepository,
	sensorTypeRepo *repository.SensorTypeRepository,
) *AssetWithSensorsService {
	return &AssetWithSensorsService{
		assetRepo:       assetRepo,
		assetSensorRepo: assetSensorRepo,
		assetTypeRepo:   assetTypeRepo,
		locationRepo:    locationRepo,
		sensorTypeRepo:  sensorTypeRepo,
	}
}

// CreateAssetWithSensors creates a new asset and automatically generates associated sensors
func (s *AssetWithSensorsService) CreateAssetWithSensors(ctx context.Context, req *dto.CreateAssetWithSensorsRequest) (*dto.AssetWithSensorsResponse, error) {
	log.Printf("Creating asset with sensors: %+v", req)

	// Validate asset type
	_, err := s.assetTypeRepo.GetByID(ctx, req.AssetTypeID)
	if err != nil {
		log.Printf("Invalid asset type ID: %s, error: %v", req.AssetTypeID, err)
		return nil, common.NewValidationError("invalid asset type: Asset type not found", err)
	}

	// Validate location
	_, err = s.locationRepo.GetByID(ctx, req.LocationID)
	if err != nil {
		log.Printf("Invalid location ID: %s, error: %v", req.LocationID, err)
		return nil, common.NewValidationError("invalid location: Location not found", err)
	}

	// Validate all sensor types exist
	for i, sensorReq := range req.SensorTypes {
		_, err = s.sensorTypeRepo.GetByID(sensorReq.SensorTypeID)
		if err != nil {
			log.Printf("Invalid sensor type ID at index %d: %s, error: %v", i, sensorReq.SensorTypeID, err)
			return nil, common.NewValidationError(fmt.Sprintf("invalid sensor type at index %d: Sensor type not found", i), err)
		}
	}

	// Create the asset first
	now := time.Now()
	asset := &entity.Asset{
		ID:          uuid.New(),
		Name:        req.Name,
		AssetTypeID: req.AssetTypeID,
		LocationID:  req.LocationID,
		Status:      req.Status,
		Properties:  req.Properties,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// Set default status if not provided
	if asset.Status == "" {
		asset.Status = "active"
	}

	log.Printf("Creating asset entity: %+v", asset)

	// Save asset to database
	err = s.assetRepo.Create(ctx, asset)
	if err != nil {
		log.Printf("Failed to create asset: %v", err)
		return nil, fmt.Errorf("failed to create asset: %w", err)
	}

	log.Printf("Asset created successfully with ID: %s", asset.ID)

	// Create asset sensors automatically
	var createdSensors []dto.AssetSensorResponse
	var sensorErrors []string

	for i, sensorReq := range req.SensorTypes {
		log.Printf("Creating sensor %d: %+v", i, sensorReq)

		// Set default sensor status if not provided
		sensorStatus := sensorReq.Status
		if sensorStatus == "" {
			sensorStatus = "active"
		}

		// Create asset sensor entity
		assetSensor := &entity.AssetSensor{
			ID:            uuid.New(),
			TenantID:      asset.TenantID, // Inherit tenant_id from asset
			AssetID:       asset.ID,
			SensorTypeID:  sensorReq.SensorTypeID,
			Name:          sensorReq.Name,
			Status:        sensorStatus,
			Configuration: sensorReq.Configuration,
			CreatedAt:     now,
		}

		log.Printf("Creating asset sensor entity: %+v", assetSensor)

		// Save asset sensor to database
		err = s.assetSensorRepo.Create(ctx, assetSensor)
		if err != nil {
			log.Printf("Failed to create asset sensor %d: %v", i, err)
			sensorErrors = append(sensorErrors, fmt.Sprintf("Failed to create sensor '%s': %v", sensorReq.Name, err))
			continue
		}

		log.Printf("Asset sensor created successfully with ID: %s", assetSensor.ID)

		// Convert to response DTO
		var tenantID uuid.UUID
		if assetSensor.TenantID != nil {
			tenantID = *assetSensor.TenantID
		}

		sensorResponse := dto.AssetSensorResponse{
			ID:            assetSensor.ID,
			TenantID:      tenantID,
			AssetID:       assetSensor.AssetID,
			SensorTypeID:  assetSensor.SensorTypeID,
			Name:          assetSensor.Name,
			Status:        assetSensor.Status,
			Configuration: assetSensor.Configuration,
			CreatedAt:     assetSensor.CreatedAt,
			UpdatedAt:     assetSensor.UpdatedAt,
		}

		createdSensors = append(createdSensors, sensorResponse)
	}

	// Log any sensor creation errors but don't fail the entire operation
	if len(sensorErrors) > 0 {
		log.Printf("Some sensors failed to create: %v", sensorErrors)
	}

	// Convert asset to response DTO
	var assetTenantID *uuid.UUID
	if asset.TenantID != nil {
		assetTenantID = asset.TenantID
	}

	assetResponse := dto.AssetResponse{
		ID:          asset.ID,
		Name:        asset.Name,
		AssetTypeID: asset.AssetTypeID,
		LocationID:  asset.LocationID,
		TenantID:    assetTenantID,
		Status:      asset.Status,
		CreatedAt:   asset.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   asset.UpdatedAt.Format(time.RFC3339),
	}

	// Convert properties to string if present
	if asset.Properties != nil {
		properties := string(asset.Properties)
		assetResponse.Properties = &properties
	}

	response := &dto.AssetWithSensorsResponse{
		Asset:   assetResponse,
		Sensors: createdSensors,
	}

	log.Printf("Asset with sensors created successfully. Asset ID: %s, Sensors created: %d", asset.ID, len(createdSensors))

	return response, nil
}

// GetAssetWithSensors retrieves an asset along with all its sensors
func (s *AssetWithSensorsService) GetAssetWithSensors(ctx context.Context, assetID uuid.UUID) (*dto.AssetWithSensorsResponse, error) {
	// Get the asset
	asset, err := s.assetRepo.GetByID(ctx, assetID)
	if err != nil {
		return nil, fmt.Errorf("failed to get asset: %w", err)
	}
	if asset == nil {
		return nil, common.NewNotFoundError("asset", assetID.String())
	}

	// Get all sensors for this asset
	sensors, err := s.assetSensorRepo.GetByAssetID(ctx, assetID)
	if err != nil {
		return nil, fmt.Errorf("failed to get asset sensors: %w", err)
	}

	// Convert asset to response DTO
	var assetTenantID *uuid.UUID
	if asset.TenantID != nil {
		assetTenantID = asset.TenantID
	}

	assetResponse := dto.AssetResponse{
		ID:          asset.ID,
		Name:        asset.Name,
		AssetTypeID: asset.AssetTypeID,
		LocationID:  asset.LocationID,
		TenantID:    assetTenantID,
		Status:      asset.Status,
		CreatedAt:   asset.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   asset.UpdatedAt.Format(time.RFC3339),
	}

	// Convert properties to string if present
	if asset.Properties != nil {
		properties := string(asset.Properties)
		assetResponse.Properties = &properties
	}

	// Convert sensors to response DTOs
	var sensorResponses []dto.AssetSensorResponse
	for _, sensor := range sensors {
		var tenantID uuid.UUID
		if sensor.AssetSensor.TenantID != nil {
			tenantID = *sensor.AssetSensor.TenantID
		}

		sensorResponse := dto.AssetSensorResponse{
			ID:                sensor.AssetSensor.ID,
			TenantID:          tenantID,
			AssetID:           sensor.AssetSensor.AssetID,
			SensorTypeID:      sensor.AssetSensor.SensorTypeID,
			Name:              sensor.AssetSensor.Name,
			Status:            sensor.AssetSensor.Status,
			Configuration:     sensor.AssetSensor.Configuration,
			LastReadingValue:  sensor.AssetSensor.LastReadingValue,
			LastReadingTime:   sensor.AssetSensor.LastReadingTime,
			LastReadingValues: sensor.AssetSensor.LastReadingValues,
			CreatedAt:         sensor.AssetSensor.CreatedAt,
			UpdatedAt:         sensor.AssetSensor.UpdatedAt,
		}

		sensorResponses = append(sensorResponses, sensorResponse)
	}

	return &dto.AssetWithSensorsResponse{
		Asset:   assetResponse,
		Sensors: sensorResponses,
	}, nil
}
