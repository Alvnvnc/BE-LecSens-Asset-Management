package seeder

import (
	"be-lecsens/asset_management/data-layer/repository"
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

type SensorTypeSeeder struct {
	repo *repository.SensorTypeRepository
}

func NewSensorTypeSeeder(db *sql.DB) *SensorTypeSeeder {
	return &SensorTypeSeeder{
		repo: repository.NewSensorTypeRepository(db),
	}
}

// sensorTypeData contains realistic sensor type data
var sensorTypeData = []repository.SensorType{
	{
		ID:           uuid.MustParse("03456789-3333-3333-3333-000000000001"),
		Name:         "pH Sensor",
		Description:  "Sensor for measuring pH levels in water and chemical solutions",
		Manufacturer: "Hach",
		Model:        "PHC101",
		Version:      "2.1.0",
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    &[]time.Time{time.Now()}[0],
	},
	{
		ID:           uuid.MustParse("03456789-3333-3333-3333-000000000002"),
		Name:         "Turbidity Sensor",
		Description:  "Optical sensor for measuring water turbidity and suspended particles",
		Manufacturer: "Hach",
		Model:        "TU5300sc",
		Version:      "1.5.0",
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    &[]time.Time{time.Now()}[0],
	},
	{
		ID:           uuid.MustParse("03456789-3333-3333-3333-000000000003"),
		Name:         "Dissolved Oxygen Sensor",
		Description:  "Electrochemical sensor for measuring dissolved oxygen in water",
		Manufacturer: "YSI",
		Model:        "ProODO",
		Version:      "3.0.0",
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    &[]time.Time{time.Now()}[0],
	},
	{
		ID:           uuid.MustParse("03456789-3333-3333-3333-000000000004"),
		Name:         "PM2.5 Sensor",
		Description:  "Laser scattering sensor for measuring fine particulate matter PM2.5",
		Manufacturer: "Sensirion",
		Model:        "SPS30",
		Version:      "1.2.0",
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    &[]time.Time{time.Now()}[0],
	},
	{
		ID:           uuid.MustParse("03456789-3333-3333-3333-000000000005"),
		Name:         "CO2 Sensor",
		Description:  "NDIR sensor for measuring carbon dioxide concentration",
		Manufacturer: "Sensirion",
		Model:        "SCD30",
		Version:      "2.0.0",
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    &[]time.Time{time.Now()}[0],
	},
	{
		ID:           uuid.MustParse("03456789-3333-3333-3333-000000000006"),
		Name:         "Temperature Sensor",
		Description:  "High-precision RTD temperature sensor for industrial applications",
		Manufacturer: "Endress+Hauser",
		Model:        "TMT181",
		Version:      "4.1.0",
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    &[]time.Time{time.Now()}[0],
	},
	{
		ID:           uuid.MustParse("03456789-3333-3333-3333-000000000007"),
		Name:         "Humidity Sensor",
		Description:  "Capacitive humidity sensor with temperature compensation",
		Manufacturer: "Sensirion",
		Model:        "SHT35",
		Version:      "1.8.0",
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    &[]time.Time{time.Now()}[0],
	},
	{
		ID:           uuid.MustParse("03456789-3333-3333-3333-000000000008"),
		Name:         "Pressure Sensor",
		Description:  "Piezoresistive pressure sensor for liquid and gas applications",
		Manufacturer: "Endress+Hauser",
		Model:        "Cerabar PMC21",
		Version:      "3.2.0",
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    &[]time.Time{time.Now()}[0],
	},
	{
		ID:           uuid.MustParse("03456789-3333-3333-3333-000000000009"),
		Name:         "Flow Sensor",
		Description:  "Ultrasonic flow sensor for non-invasive flow measurement",
		Manufacturer: "Endress+Hauser",
		Model:        "Prosonic Flow 93T",
		Version:      "2.5.0",
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    &[]time.Time{time.Now()}[0],
	},
	{
		ID:           uuid.MustParse("03456789-3333-3333-3333-000000000010"),
		Name:         "Vibration Sensor",
		Description:  "MEMS accelerometer for machinery vibration monitoring",
		Manufacturer: "SKF",
		Model:        "CMSS 2200",
		Version:      "1.6.0",
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    &[]time.Time{time.Now()}[0],
	},
	{
		ID:           uuid.MustParse("03456789-3333-3333-3333-000000000011"),
		Name:         "Level Sensor",
		Description:  "Radar level sensor for continuous level measurement",
		Manufacturer: "Endress+Hauser",
		Model:        "Micropilot FMR20",
		Version:      "2.8.0",
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    &[]time.Time{time.Now()}[0],
	},
	{
		ID:           uuid.MustParse("03456789-3333-3333-3333-000000000012"),
		Name:         "Power Meter",
		Description:  "Three-phase power quality analyzer and energy meter",
		Manufacturer: "Schneider Electric",
		Model:        "PM8000",
		Version:      "3.5.0",
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    &[]time.Time{time.Now()}[0],
	},
	{
		ID:           uuid.MustParse("03456789-3333-3333-3333-000000000013"),
		Name:         "Gas Sensor",
		Description:  "Multi-gas sensor for detecting toxic and combustible gases",
		Manufacturer: "Honeywell",
		Model:        "MIDAS-E-H2S",
		Version:      "1.9.0",
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    &[]time.Time{time.Now()}[0],
	},
	{
		ID:           uuid.MustParse("03456789-3333-3333-3333-000000000014"),
		Name:         "Conductivity Sensor",
		Description:  "Inductive conductivity sensor for water quality monitoring",
		Manufacturer: "Endress+Hauser",
		Model:        "Indumax CLS15D",
		Version:      "2.3.0",
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    &[]time.Time{time.Now()}[0],
	},
	{
		ID:           uuid.MustParse("03456789-3333-3333-3333-000000000015"),
		Name:         "Noise Sensor",
		Description:  "Sound level meter for environmental noise monitoring",
		Manufacturer: "Bruel & Kjaer",
		Model:        "Type 2250",
		Version:      "1.4.0",
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    &[]time.Time{time.Now()}[0],
	},
	{
		ID:           uuid.MustParse("03456789-3333-3333-3333-000000000016"),
		Name:         "Chlorine Sensor",
		Description:  "Electrochemical sensor for measuring chlorine levels in water",
		Manufacturer: "Hach",
		Model:        "CL17",
		Version:      "2.0.0",
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    &[]time.Time{time.Now()}[0],
	},
	{
		ID:           uuid.MustParse("03456789-3333-3333-3333-000000000017"),
		Name:         "TDS Sensor",
		Description:  "Total Dissolved Solids sensor for water quality monitoring",
		Manufacturer: "Hanna Instruments",
		Model:        "HI98301",
		Version:      "1.5.0",
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    &[]time.Time{time.Now()}[0],
	},
	{
		ID:           uuid.MustParse("03456789-3333-3333-3333-000000000018"),
		Name:         "ORP Sensor",
		Description:  "Oxidation-Reduction Potential sensor for water treatment monitoring",
		Manufacturer: "Hach",
		Model:        "MQ45",
		Version:      "3.1.0",
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    &[]time.Time{time.Now()}[0],
	},
}

func (s *SensorTypeSeeder) Seed(ctx context.Context) error {
	log.Println("Starting SensorType seeding...")

	// Check if sensor types already exist
	existingSensorTypes, err := s.repo.GetAll()
	if err != nil {
		return fmt.Errorf("failed to check existing sensor types: %w", err)
	}

	if len(existingSensorTypes) > 0 {
		log.Printf("Sensor types already exist (%d records), skipping sensor type seeding", len(existingSensorTypes))
		return nil
	}

	// Insert sensor types
	for i, sensorType := range sensorTypeData {
		err := s.repo.Create(&sensorType)
		if err != nil {
			return fmt.Errorf("failed to create sensor type %d (%s): %w", i+1, sensorType.Name, err)
		}
		log.Printf("Created sensor type: %s (%s %s)", sensorType.Name, sensorType.Manufacturer, sensorType.Model)
	}

	log.Printf("Successfully seeded %d sensor types", len(sensorTypeData))
	return nil
}

func (s *SensorTypeSeeder) GetSensorTypeIDs() []uuid.UUID {
	var ids []uuid.UUID
	for _, sensorType := range sensorTypeData {
		ids = append(ids, sensorType.ID)
	}
	return ids
}

func (s *SensorTypeSeeder) GetSensorTypeByCategory(category string) []uuid.UUID {
	var ids []uuid.UUID
	categoryMap := map[string][]string{
		"water":         {"pH Sensor", "Turbidity Sensor", "Dissolved Oxygen Sensor", "Conductivity Sensor"},
		"air":           {"PM2.5 Sensor", "CO2 Sensor", "Gas Sensor", "Noise Sensor"},
		"industrial":    {"Temperature Sensor", "Pressure Sensor", "Flow Sensor", "Vibration Sensor", "Level Sensor"},
		"electrical":    {"Power Meter"},
		"environmental": {"Temperature Sensor", "Humidity Sensor"},
	}

	if sensorNames, exists := categoryMap[category]; exists {
		for _, sensorType := range sensorTypeData {
			for _, name := range sensorNames {
				if sensorType.Name == name {
					ids = append(ids, sensorType.ID)
					break
				}
			}
		}
	}
	return ids
}
