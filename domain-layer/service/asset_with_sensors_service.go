package service

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"be-lecsens/asset_management/data-layer/entity"
	"be-lecsens/asset_management/data-layer/repository"
	"be-lecsens/asset_management/helpers/common"
	"be-lecsens/asset_management/helpers/dto"

	"github.com/google/uuid"
)

// ContextKey is a custom type for context keys to avoid collisions
type ContextKey int

const (
	// TxContextKey is the key for storing transaction in context
	TxContextKey ContextKey = iota
)

// AssetWithSensorsService handles business logic for creating assets with sensors
type AssetWithSensorsService struct {
	assetRepo       repository.AssetRepository
	assetSensorRepo repository.AssetSensorRepository
	assetTypeRepo   *repository.AssetTypeRepository
	locationRepo    *repository.LocationRepository
	sensorTypeRepo  *repository.SensorTypeRepository
	db              *sql.DB // Add database connection for transaction management
}

// NewAssetWithSensorsService creates a new instance of AssetWithSensorsService
func NewAssetWithSensorsService(
	assetRepo repository.AssetRepository,
	assetSensorRepo repository.AssetSensorRepository,
	assetTypeRepo *repository.AssetTypeRepository,
	locationRepo *repository.LocationRepository,
	sensorTypeRepo *repository.SensorTypeRepository,
	db *sql.DB, // Add database parameter
) *AssetWithSensorsService {
	return &AssetWithSensorsService{
		assetRepo:       assetRepo,
		assetSensorRepo: assetSensorRepo,
		assetTypeRepo:   assetTypeRepo,
		locationRepo:    locationRepo,
		sensorTypeRepo:  sensorTypeRepo,
		db:              db,
	}
}

// CreateAssetWithSensors creates a new asset and automatically generates associated sensors
// This operation is atomic - either all succeed or all fail
func (s *AssetWithSensorsService) CreateAssetWithSensors(ctx context.Context, req *dto.CreateAssetWithSensorsRequest) (*dto.AssetWithSensorsResponse, error) {
	log.Printf("Creating asset with sensors: %+v", req)

	// Start database transaction for atomic operation
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // Rollback if not committed

	// Create context with transaction for repository operations
	txCtx := context.WithValue(ctx, TxContextKey, tx)

	// Validate asset type
	_, err = s.assetTypeRepo.GetByID(txCtx, req.AssetTypeID)
	if err != nil {
		log.Printf("Invalid asset type ID: %s, error: %v", req.AssetTypeID, err)
		return nil, common.NewValidationError("invalid asset type: Asset type not found", err)
	}

	// Validate location
	_, err = s.locationRepo.GetByID(txCtx, req.LocationID)
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
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// Handle Properties field conversion
	if req.Properties != nil {
		asset.Properties = req.Properties
	}

	// Set default status if not provided
	if asset.Status == "" {
		asset.Status = "active"
	}

	log.Printf("Creating asset entity: %+v", asset)

	// Save asset to database using transaction context
	err = s.assetRepo.Create(txCtx, asset)
	if err != nil {
		log.Printf("Failed to create asset: %v", err)
		return nil, fmt.Errorf("failed to create asset: %w", err)
	}

	log.Printf("Asset created successfully with ID: %s", asset.ID)

	// Create asset sensors automatically
	var createdSensors []dto.AssetSensorResponse

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
			Name:          asset.Name, // Use asset name as sensor name
			Status:        sensorStatus,
			Configuration: sensorReq.Configuration,
			CreatedAt:     now,
		}

		log.Printf("Creating asset sensor entity: %+v", assetSensor)

		// Save asset sensor to database using transaction context
		err = s.assetSensorRepo.Create(txCtx, assetSensor)
		if err != nil {
			log.Printf("Failed to create asset sensor %d: %v", i, err)
			// Return error to rollback transaction
			return nil, fmt.Errorf("failed to create sensor '%s': %w", asset.Name, err)
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

	// Commit transaction - all operations succeeded
	if err := tx.Commit(); err != nil {
		log.Printf("Failed to commit transaction: %v", err)
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
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
	log.Printf("Getting asset with sensors for asset ID: %s", assetID)

	// Get the asset
	asset, err := s.assetRepo.GetByID(ctx, assetID)
	if err != nil {
		log.Printf("Failed to get asset: %v", err)
		return nil, fmt.Errorf("failed to get asset: %w", err)
	}
	if asset == nil {
		return nil, common.NewNotFoundError("asset", assetID.String())
	}

	// Get all sensors for this asset
	sensors, err := s.assetSensorRepo.GetByAssetID(ctx, assetID)
	if err != nil {
		log.Printf("Failed to get asset sensors: %v", err)
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

	response := &dto.AssetWithSensorsResponse{
		Asset:   assetResponse,
		Sensors: sensorResponses,
	}

	log.Printf("Successfully retrieved asset with %d sensors", len(sensorResponses))
	return response, nil
}

// UpdateAssetWithSensors updates an asset and can add/remove/update sensors
func (s *AssetWithSensorsService) UpdateAssetWithSensors(ctx context.Context, assetID uuid.UUID, req *dto.UpdateAssetWithSensorsRequest) (*dto.AssetWithSensorsResponse, error) {
	log.Printf("Updating asset with sensors for asset ID: %s", assetID)

	// Start database transaction for atomic operation
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // Rollback if not committed

	// Create context with transaction for repository operations
	txCtx := context.WithValue(ctx, TxContextKey, tx)

	// Get existing asset
	existingAsset, err := s.assetRepo.GetByID(txCtx, assetID)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing asset: %w", err)
	}
	if existingAsset == nil {
		return nil, common.NewNotFoundError("asset", assetID.String())
	}

	// Update asset fields if provided
	updatedAsset := *existingAsset
	now := time.Now()

	if req.Name != nil {
		updatedAsset.Name = *req.Name
	}
	if req.AssetTypeID != nil {
		// Validate asset type
		_, err = s.assetTypeRepo.GetByID(txCtx, *req.AssetTypeID)
		if err != nil {
			return nil, common.NewValidationError("invalid asset type: Asset type not found", err)
		}
		updatedAsset.AssetTypeID = *req.AssetTypeID
	}
	if req.LocationID != nil {
		// Validate location
		_, err = s.locationRepo.GetByID(txCtx, *req.LocationID)
		if err != nil {
			return nil, common.NewValidationError("invalid location: Location not found", err)
		}
		updatedAsset.LocationID = *req.LocationID
	}
	if req.Status != nil {
		updatedAsset.Status = *req.Status
	}
	if req.Properties != nil {
		updatedAsset.Properties = req.Properties
	}

	updatedAsset.UpdatedAt = now

	// Update asset in database
	err = s.assetRepo.Update(txCtx, &updatedAsset)
	if err != nil {
		return nil, fmt.Errorf("failed to update asset: %w", err)
	}

	// Handle sensor updates if provided
	var sensorResponses []dto.AssetSensorResponse
	if req.SensorTypes != nil {
		// Get existing sensors
		existingSensors, err := s.assetSensorRepo.GetByAssetID(txCtx, assetID)
		if err != nil {
			return nil, fmt.Errorf("failed to get existing sensors: %w", err)
		}

		// Create a map of existing sensors by sensor type for easier lookup
		existingSensorMap := make(map[uuid.UUID]*repository.AssetSensorWithDetails)
		for _, sensor := range existingSensors {
			existingSensorMap[sensor.AssetSensor.SensorTypeID] = sensor
		}

		// Process each requested sensor type
		for _, sensorReq := range req.SensorTypes {
			// Validate sensor type
			_, err = s.sensorTypeRepo.GetByID(sensorReq.SensorTypeID)
			if err != nil {
				return nil, common.NewValidationError("invalid sensor type: Sensor type not found", err)
			}

			if existingSensor, exists := existingSensorMap[sensorReq.SensorTypeID]; exists {
				// Update existing sensor
				updatedSensor := *existingSensor.AssetSensor
				if sensorReq.Status != "" {
					updatedSensor.Status = sensorReq.Status
				}
				if sensorReq.Configuration != nil {
					updatedSensor.Configuration = sensorReq.Configuration
				}
				updatedSensor.UpdatedAt = &now

				err = s.assetSensorRepo.Update(txCtx, &updatedSensor)
				if err != nil {
					return nil, fmt.Errorf("failed to update sensor: %w", err)
				}

				// Convert to response DTO
				var tenantID uuid.UUID
				if updatedSensor.TenantID != nil {
					tenantID = *updatedSensor.TenantID
				}

				sensorResponse := dto.AssetSensorResponse{
					ID:            updatedSensor.ID,
					TenantID:      tenantID,
					AssetID:       updatedSensor.AssetID,
					SensorTypeID:  updatedSensor.SensorTypeID,
					Name:          updatedSensor.Name,
					Status:        updatedSensor.Status,
					Configuration: updatedSensor.Configuration,
					CreatedAt:     updatedSensor.CreatedAt,
					UpdatedAt:     updatedSensor.UpdatedAt,
				}
				sensorResponses = append(sensorResponses, sensorResponse)

				// Remove from map to track which sensors were processed
				delete(existingSensorMap, sensorReq.SensorTypeID)
			} else {
				// Create new sensor
				sensorStatus := sensorReq.Status
				if sensorStatus == "" {
					sensorStatus = "active"
				}

				newSensor := &entity.AssetSensor{
					ID:            uuid.New(),
					TenantID:      updatedAsset.TenantID,
					AssetID:       assetID,
					SensorTypeID:  sensorReq.SensorTypeID,
					Name:          updatedAsset.Name,
					Status:        sensorStatus,
					Configuration: sensorReq.Configuration,
					CreatedAt:     now,
				}

				err = s.assetSensorRepo.Create(txCtx, newSensor)
				if err != nil {
					return nil, fmt.Errorf("failed to create new sensor: %w", err)
				}

				// Convert to response DTO
				var tenantID uuid.UUID
				if newSensor.TenantID != nil {
					tenantID = *newSensor.TenantID
				}

				sensorResponse := dto.AssetSensorResponse{
					ID:            newSensor.ID,
					TenantID:      tenantID,
					AssetID:       newSensor.AssetID,
					SensorTypeID:  newSensor.SensorTypeID,
					Name:          newSensor.Name,
					Status:        newSensor.Status,
					Configuration: newSensor.Configuration,
					CreatedAt:     newSensor.CreatedAt,
					UpdatedAt:     newSensor.UpdatedAt,
				}
				sensorResponses = append(sensorResponses, sensorResponse)
			}
		}

		// Remove sensors that were not in the request (optional - you might want to keep them)
		// Uncomment below if you want to remove sensors not specified in the update request
		/*
			for _, remainingSensor := range existingSensorMap {
				err = s.assetSensorRepo.Delete(txCtx, remainingSensor.AssetSensor.ID)
				if err != nil {
					return nil, fmt.Errorf("failed to delete sensor: %w", err)
				}
			}
		*/
	} else {
		// If no sensor updates requested, just return existing sensors
		existingSensors, err := s.assetSensorRepo.GetByAssetID(txCtx, assetID)
		if err != nil {
			return nil, fmt.Errorf("failed to get existing sensors: %w", err)
		}

		for _, sensor := range existingSensors {
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
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Convert asset to response DTO
	var assetTenantID *uuid.UUID
	if updatedAsset.TenantID != nil {
		assetTenantID = updatedAsset.TenantID
	}

	assetResponse := dto.AssetResponse{
		ID:          updatedAsset.ID,
		Name:        updatedAsset.Name,
		AssetTypeID: updatedAsset.AssetTypeID,
		LocationID:  updatedAsset.LocationID,
		TenantID:    assetTenantID,
		Status:      updatedAsset.Status,
		CreatedAt:   updatedAsset.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   updatedAsset.UpdatedAt.Format(time.RFC3339),
	}

	// Convert properties to string if present
	if updatedAsset.Properties != nil {
		properties := string(updatedAsset.Properties)
		assetResponse.Properties = &properties
	}

	response := &dto.AssetWithSensorsResponse{
		Asset:   assetResponse,
		Sensors: sensorResponses,
	}

	log.Printf("Asset with sensors updated successfully. Asset ID: %s, Sensors: %d", assetID, len(sensorResponses))
	return response, nil
}

// DeleteAssetWithSensors deletes an asset and all its associated sensors
func (s *AssetWithSensorsService) DeleteAssetWithSensors(ctx context.Context, assetID uuid.UUID) error {
	log.Printf("Deleting asset with sensors for asset ID: %s", assetID)

	// Start database transaction for atomic operation
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // Rollback if not committed

	// Create context with transaction for repository operations
	txCtx := context.WithValue(ctx, TxContextKey, tx)

	// Check if asset exists
	asset, err := s.assetRepo.GetByID(txCtx, assetID)
	if err != nil {
		return fmt.Errorf("failed to get asset: %w", err)
	}
	if asset == nil {
		return common.NewNotFoundError("asset", assetID.String())
	}

	// Get all sensors for this asset
	sensors, err := s.assetSensorRepo.GetByAssetID(txCtx, assetID)
	if err != nil {
		return fmt.Errorf("failed to get asset sensors: %w", err)
	}

	// Delete all sensors first (due to foreign key constraints)
	for _, sensor := range sensors {
		err = s.assetSensorRepo.Delete(txCtx, sensor.AssetSensor.ID)
		if err != nil {
			return fmt.Errorf("failed to delete sensor %s: %w", sensor.AssetSensor.ID, err)
		}
		log.Printf("Deleted sensor: %s", sensor.AssetSensor.ID)
	}

	// Delete the asset
	err = s.assetRepo.Delete(txCtx, assetID)
	if err != nil {
		return fmt.Errorf("failed to delete asset: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Printf("Asset and %d sensors deleted successfully. Asset ID: %s", len(sensors), assetID)
	return nil
}

// ListAssetsWithSensors retrieves a paginated list of assets with their sensors
func (s *AssetWithSensorsService) ListAssetsWithSensors(ctx context.Context, page, pageSize int) (*dto.AssetWithSensorsListResponse, error) {
	log.Printf("Listing assets with sensors - page: %d, pageSize: %d", page, pageSize)

	// Get paginated assets
	assetList, err := s.assetRepo.List(ctx, nil, page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("failed to get assets: %w", err)
	}

	var assetsWithSensors []dto.AssetWithSensorsResponse
	for _, asset := range assetList {
		// Get sensors for each asset
		sensors, err := s.assetSensorRepo.GetByAssetID(ctx, asset.ID)
		if err != nil {
			log.Printf("Warning: failed to get sensors for asset %s: %v", asset.ID, err)
			// Continue with empty sensors list instead of failing
			sensors = []*repository.AssetSensorWithDetails{}
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

		assetWithSensors := dto.AssetWithSensorsResponse{
			Asset:   assetResponse,
			Sensors: sensorResponses,
		}

		assetsWithSensors = append(assetsWithSensors, assetWithSensors)
	}

	response := &dto.AssetWithSensorsListResponse{
		Assets:     assetsWithSensors,
		Page:       page,
		Limit:      pageSize,
		Total:      int64(len(assetsWithSensors)),
		TotalPages: (len(assetsWithSensors) + pageSize - 1) / pageSize,
	}

	log.Printf("Successfully retrieved %d assets with sensors", len(assetsWithSensors))
	return response, nil
}

// GetAssetSensorsCount returns the count of sensors for a specific asset
func (s *AssetWithSensorsService) GetAssetSensorsCount(ctx context.Context, assetID uuid.UUID) (int, error) {
	log.Printf("Getting sensor count for asset ID: %s", assetID)

	// Check if asset exists
	asset, err := s.assetRepo.GetByID(ctx, assetID)
	if err != nil {
		return 0, fmt.Errorf("failed to get asset: %w", err)
	}
	if asset == nil {
		return 0, common.NewNotFoundError("asset", assetID.String())
	}

	// Get sensors for this asset
	sensors, err := s.assetSensorRepo.GetByAssetID(ctx, assetID)
	if err != nil {
		return 0, fmt.Errorf("failed to get asset sensors: %w", err)
	}

	count := len(sensors)
	log.Printf("Asset %s has %d sensors", assetID, count)
	return count, nil
}

// BulkCreateAssetsWithSensors creates multiple assets with their sensors in a single transaction
func (s *AssetWithSensorsService) BulkCreateAssetsWithSensors(ctx context.Context, requests []*dto.CreateAssetWithSensorsRequest) (*dto.BulkCreateAssetWithSensorsResponse, error) {
	log.Printf("Bulk creating %d assets with sensors", len(requests))

	// Start database transaction for atomic operation
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // Rollback if not committed

	// Create context with transaction for repository operations
	txCtx := context.WithValue(ctx, TxContextKey, tx)

	var results []dto.CreateAssetWithSensorsResult
	var errors []string

	for i, req := range requests {
		log.Printf("Processing asset %d: %s", i+1, req.Name)

		// Create single asset with sensors using existing method logic
		result, err := s.createSingleAssetWithSensors(txCtx, req)
		if err != nil {
			errorMsg := fmt.Sprintf("Failed to create asset %d (%s): %v", i+1, req.Name, err)
			errors = append(errors, errorMsg)
			log.Print(errorMsg)
			continue
		}

		results = append(results, *result)
	}

	// If any errors occurred, rollback transaction
	if len(errors) > 0 {
		return &dto.BulkCreateAssetWithSensorsResponse{
			Results:      results,
			SuccessCount: len(results),
			ErrorCount:   len(errors),
			Errors:       errors,
		}, fmt.Errorf("bulk create failed with %d errors", len(errors))
	}

	// Commit transaction if all succeeded
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	response := &dto.BulkCreateAssetWithSensorsResponse{
		Results:      results,
		SuccessCount: len(results),
		ErrorCount:   0,
		Errors:       []string{},
	}

	log.Printf("Bulk create completed successfully: %d assets created", len(results))
	return response, nil
}

// createSingleAssetWithSensors is a helper method for creating a single asset with sensors within a transaction
func (s *AssetWithSensorsService) createSingleAssetWithSensors(ctx context.Context, req *dto.CreateAssetWithSensorsRequest) (*dto.CreateAssetWithSensorsResult, error) {
	// Validate asset type
	_, err := s.assetTypeRepo.GetByID(ctx, req.AssetTypeID)
	if err != nil {
		return nil, common.NewValidationError("invalid asset type: Asset type not found", err)
	}

	// Validate location
	_, err = s.locationRepo.GetByID(ctx, req.LocationID)
	if err != nil {
		return nil, common.NewValidationError("invalid location: Location not found", err)
	}

	// Validate all sensor types exist
	for i, sensorReq := range req.SensorTypes {
		_, err = s.sensorTypeRepo.GetByID(sensorReq.SensorTypeID)
		if err != nil {
			return nil, common.NewValidationError(fmt.Sprintf("invalid sensor type at index %d: Sensor type not found", i), err)
		}
	}

	// Create the asset
	now := time.Now()
	asset := &entity.Asset{
		ID:          uuid.New(),
		Name:        req.Name,
		AssetTypeID: req.AssetTypeID,
		LocationID:  req.LocationID,
		Status:      req.Status,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// Handle Properties field conversion
	if req.Properties != nil {
		asset.Properties = req.Properties
	}

	// Set default status if not provided
	if asset.Status == "" {
		asset.Status = "active"
	}

	// Save asset to database
	err = s.assetRepo.Create(ctx, asset)
	if err != nil {
		return nil, fmt.Errorf("failed to create asset: %w", err)
	}

	// Create asset sensors
	var sensorIDs []uuid.UUID
	for _, sensorReq := range req.SensorTypes {
		sensorStatus := sensorReq.Status
		if sensorStatus == "" {
			sensorStatus = "active"
		}

		assetSensor := &entity.AssetSensor{
			ID:            uuid.New(),
			TenantID:      asset.TenantID,
			AssetID:       asset.ID,
			SensorTypeID:  sensorReq.SensorTypeID,
			Name:          asset.Name,
			Status:        sensorStatus,
			Configuration: sensorReq.Configuration,
			CreatedAt:     now,
		}

		err = s.assetSensorRepo.Create(ctx, assetSensor)
		if err != nil {
			return nil, fmt.Errorf("failed to create sensor: %w", err)
		}

		sensorIDs = append(sensorIDs, assetSensor.ID)
	}

	result := &dto.CreateAssetWithSensorsResult{
		AssetID:   asset.ID,
		SensorIDs: sensorIDs,
		CreatedAt: now,
		TenantID:  asset.TenantID,
		Errors:    []string{},
	}

	return result, nil
}
