package seeder

import (
	"be-lecsens/asset_management/data-layer/entity"
	"be-lecsens/asset_management/data-layer/repository"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

type AssetTypeSeeder struct {
	repo *repository.AssetTypeRepository
}

func NewAssetTypeSeeder(db *sql.DB) *AssetTypeSeeder {
	return &AssetTypeSeeder{
		repo: repository.NewAssetTypeRepository(db),
	}
}

// assetTypeData contains realistic asset type data for industrial monitoring
var assetTypeData = []entity.AssetType{
	{
		ID:          uuid.MustParse("02345678-2222-2222-2222-000000000001"),
		Name:        "Water Quality Monitor",
		Category:    "sensor",
		Description: "System for monitoring water quality parameters including pH, turbidity, dissolved oxygen",
		PropertiesSchema: json.RawMessage(`{
			"type": "object",
			"properties": {
				"flow_rate": {"type": "number", "unit": "L/min", "min": 0, "max": 1000},
				"ph_range": {"type": "object", "properties": {"min": {"type": "number"}, "max": {"type": "number"}}},
				"turbidity_range": {"type": "object", "properties": {"min": {"type": "number"}, "max": {"type": "number"}}},
				"temperature_range": {"type": "object", "properties": {"min": {"type": "number"}, "max": {"type": "number"}}},
				"installation_type": {"type": "string", "enum": ["inline", "bypass", "immersion"]},
				"power_source": {"type": "string", "enum": ["AC", "DC", "Battery", "Solar"]},
				"communication": {"type": "array", "items": {"type": "string", "enum": ["WiFi", "LoRa", "4G", "Ethernet"]}}
			}
		}`),
		CreatedAt: time.Now(),
		UpdatedAt: &[]time.Time{time.Now()}[0],
	},
	{
		ID:          uuid.MustParse("02345678-2222-2222-2222-000000000002"),
		Name:        "Water Flow Meter",
		Category:    "sensor",
		Description: "System for monitoring water flow rates and volumes",
		PropertiesSchema: json.RawMessage(`{
			"type": "object",
			"properties": {
				"flow_range": {"type": "object", "properties": {"min": {"type": "number"}, "max": {"type": "number"}}},
				"pipe_diameter": {"type": "number", "unit": "inch", "min": 0.5, "max": 48},
				"measurement_principle": {"type": "string", "enum": ["ultrasonic", "electromagnetic", "vortex", "turbine"]},
				"accuracy": {"type": "number", "unit": "%", "min": 0.1, "max": 5},
				"pressure_rating": {"type": "number", "unit": "bar", "min": 1, "max": 400},
				"temperature_rating": {"type": "object", "properties": {"min": {"type": "number"}, "max": {"type": "number"}}}
			}
		}`),
		CreatedAt: time.Now(),
		UpdatedAt: &[]time.Time{time.Now()}[0],
	},
	{
		ID:          uuid.MustParse("02345678-2222-2222-2222-000000000003"),
		Name:        "Water Level Sensor",
		Category:    "sensor",
		Description: "System for monitoring water levels in tanks and reservoirs",
		PropertiesSchema: json.RawMessage(`{
			"type": "object",
			"properties": {
				"measurement_range": {"type": "number", "unit": "m", "min": 0.1, "max": 100},
				"technology": {"type": "string", "enum": ["ultrasonic", "radar", "hydrostatic", "capacitive"]},
				"tank_shape": {"type": "string", "enum": ["cylindrical", "rectangular", "spherical", "irregular"]},
				"accuracy": {"type": "number", "unit": "mm", "min": 1, "max": 50},
				"beam_angle": {"type": "number", "unit": "degrees", "min": 3, "max": 30},
				"dead_zone": {"type": "number", "unit": "m", "min": 0.1, "max": 2}
			}
		}`),
		CreatedAt: time.Now(),
		UpdatedAt: &[]time.Time{time.Now()}[0],
	},
	{
		ID:          uuid.MustParse("02345678-2222-2222-2222-000000000004"),
		Name:        "Water Pressure Sensor",
		Category:    "sensor",
		Description: "System for monitoring water pressure in pipelines and systems",
		PropertiesSchema: json.RawMessage(`{
			"type": "object",
			"properties": {
				"pressure_range": {"type": "object", "properties": {"min": {"type": "number"}, "max": {"type": "number"}}},
				"pressure_type": {"type": "string", "enum": ["gauge", "absolute", "differential"]},
				"accuracy": {"type": "number", "unit": "%FS", "min": 0.05, "max": 1},
				"output_signal": {"type": "string", "enum": ["4-20mA", "0-10V", "digital"]},
				"process_connection": {"type": "string", "enum": ["threaded", "flanged", "sanitary"]},
				"wetted_materials": {"type": "array", "items": {"type": "string"}}
			}
		}`),
		CreatedAt: time.Now(),
		UpdatedAt: &[]time.Time{time.Now()}[0],
	},
	{
		ID:          uuid.MustParse("02345678-2222-2222-2222-000000000005"),
		Name:        "Water Sampling Drone",
		Category:    "drone",
		Description: "Autonomous drone for water sampling and monitoring",
		PropertiesSchema: json.RawMessage(`{
			"type": "object",
			"properties": {
				"flight_time": {"type": "number", "unit": "minutes", "min": 15, "max": 60},
				"sampling_capacity": {"type": "number", "unit": "ml", "min": 100, "max": 1000},
				"max_wind_speed": {"type": "number", "unit": "m/s", "min": 0, "max": 15},
				"waterproof_rating": {"type": "string", "enum": ["IP67", "IP68"]},
				"gps_accuracy": {"type": "number", "unit": "m", "min": 0.1, "max": 5},
				"camera_resolution": {"type": "string", "enum": ["1080p", "4K"]}
			}
		}`),
		CreatedAt: time.Now(),
		UpdatedAt: &[]time.Time{time.Now()}[0],
	},
	{
		ID:          uuid.MustParse("02345678-2222-2222-2222-000000000006"),
		Name:        "Water Quality Vehicle",
		Category:    "vehicle",
		Description: "Mobile water quality monitoring vehicle",
		PropertiesSchema: json.RawMessage(`{
			"type": "object",
			"properties": {
				"vehicle_type": {"type": "string", "enum": ["truck", "van", "boat"]},
				"equipment_capacity": {"type": "number", "unit": "kg", "min": 100, "max": 1000},
				"power_supply": {"type": "string", "enum": ["battery", "generator", "hybrid"]},
				"operating_range": {"type": "number", "unit": "km", "min": 50, "max": 500},
				"crew_capacity": {"type": "number", "min": 1, "max": 5},
				"weather_protection": {"type": "boolean"}
			}
		}`),
		CreatedAt: time.Now(),
		UpdatedAt: &[]time.Time{time.Now()}[0],
	},
	{
		ID:          uuid.MustParse("02345678-2222-2222-2222-000000000007"),
		Name:        "Water Treatment Equipment",
		Category:    "other",
		Description: "Water treatment and purification equipment",
		PropertiesSchema: json.RawMessage(`{
			"type": "object",
			"properties": {
				"treatment_capacity": {"type": "number", "unit": "L/hour", "min": 100, "max": 10000},
				"treatment_methods": {"type": "array", "items": {"type": "string", "enum": ["filtration", "uv", "chemical", "reverse_osmosis"]}},
				"power_consumption": {"type": "number", "unit": "kW", "min": 0.5, "max": 50},
				"maintenance_interval": {"type": "number", "unit": "days", "min": 7, "max": 365},
				"automation_level": {"type": "string", "enum": ["manual", "semi_auto", "full_auto"]},
				"monitoring_system": {"type": "boolean"}
			}
		}`),
		CreatedAt: time.Now(),
		UpdatedAt: &[]time.Time{time.Now()}[0],
	},
}

func (s *AssetTypeSeeder) Seed(ctx context.Context) error {
	log.Println("Starting AssetType seeding...")

	// Check if asset types already exist
	existingAssetTypes, err := s.repo.List(ctx, 1, 0)
	if err != nil {
		return fmt.Errorf("failed to check existing asset types: %w", err)
	}

	if len(existingAssetTypes) > 0 {
		log.Printf("Asset types already exist (%d records), skipping asset type seeding", len(existingAssetTypes))
		return nil
	}

	// Insert asset types
	for i, assetType := range assetTypeData {
		err := s.repo.Create(ctx, &assetType)
		if err != nil {
			return fmt.Errorf("failed to create asset type %d (%s): %w", i+1, assetType.Name, err)
		}
		log.Printf("Created asset type: %s (%s)", assetType.Name, assetType.Category)
	}

	log.Printf("Successfully seeded %d asset types", len(assetTypeData))
	return nil
}

func (s *AssetTypeSeeder) GetAssetTypeIDs() []uuid.UUID {
	var ids []uuid.UUID
	for _, assetType := range assetTypeData {
		ids = append(ids, assetType.ID)
	}
	return ids
}
