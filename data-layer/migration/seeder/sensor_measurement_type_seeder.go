package seeder

import (
	"be-lecsens/asset_management/data-layer/repository"
	"be-lecsens/asset_management/helpers/dto"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

type SensorMeasurementTypeSeeder struct {
	repo repository.SensorMeasurementTypeRepository
}

func NewSensorMeasurementTypeSeeder(db *sql.DB) *SensorMeasurementTypeSeeder {
	return &SensorMeasurementTypeSeeder{
		repo: repository.NewSensorMeasurementTypeRepository(db),
	}
}

// sensorMeasurementTypeData contains realistic measurement type data
var sensorMeasurementTypeData = []dto.SensorMeasurementTypeDTO{
	{
		ID:           uuid.MustParse("04567890-4444-4444-4444-000000000001"),
		SensorTypeID: uuid.MustParse("03456789-3333-3333-3333-000000000001"), // pH Sensor
		Name:         "Water pH Measurement",
		Description:  stringPtr("pH measurement for water quality monitoring"),
		PropertiesSchema: json.RawMessage(`{
			"type": "object",
			"properties": {
				"ph": {
					"type": "number",
					"description": "pH value",
					"minimum": 0,
					"maximum": 14,
					"unit": "pH"
				},
				"temperature": {
					"type": "number", 
					"description": "Temperature compensation",
					"minimum": -10,
					"maximum": 60,
					"unit": "°C"
				}
			},
			"required": ["ph"]
		}`),
		UIConfiguration: json.RawMessage(`{
			"display": {
				"type": "gauge",
				"min": 0,
				"max": 14,
				"ranges": [
					{"min": 0, "max": 6.5, "color": "red", "label": "Acidic"},
					{"min": 6.5, "max": 7.5, "color": "green", "label": "Neutral"},
					{"min": 7.5, "max": 14, "color": "blue", "label": "Basic"}
				]
			}
		}`),
		Version:   "1.0.0",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: &[]time.Time{time.Now()}[0],
	},
	{
		ID:           uuid.MustParse("04567890-4444-4444-4444-000000000002"),
		SensorTypeID: uuid.MustParse("03456789-3333-3333-3333-000000000002"), // Turbidity Sensor
		Name:         "Water Turbidity Measurement",
		Description:  stringPtr("Turbidity measurement for water clarity assessment"),
		PropertiesSchema: json.RawMessage(`{
			"type": "object",
			"properties": {
				"turbidity": {
					"type": "number",
					"description": "Turbidity value",
					"minimum": 0,
					"maximum": 4000,
					"unit": "NTU"
				},
				"temperature": {
					"type": "number",
					"description": "Water temperature",
					"minimum": 0,
					"maximum": 50,
					"unit": "°C"
				}
			},
			"required": ["turbidity"]
		}`),
		UIConfiguration: json.RawMessage(`{
			"display": {
				"type": "line_chart",
				"ranges": [
					{"min": 0, "max": 1, "color": "green", "label": "Excellent"},
					{"min": 1, "max": 5, "color": "yellow", "label": "Good"},
					{"min": 5, "max": 25, "color": "orange", "label": "Fair"},
					{"min": 25, "max": 4000, "color": "red", "label": "Poor"}
				]
			}
		}`),
		Version:   "1.0.0",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: &[]time.Time{time.Now()}[0],
	},
	{
		ID:           uuid.MustParse("04567890-4444-4444-4444-000000000003"),
		SensorTypeID: uuid.MustParse("03456789-3333-3333-3333-000000000003"), // Dissolved Oxygen Sensor
		Name:         "Dissolved Oxygen Measurement",
		Description:  stringPtr("Dissolved oxygen measurement for aquatic environment monitoring"),
		PropertiesSchema: json.RawMessage(`{
			"type": "object",
			"properties": {
				"dissolved_oxygen": {
					"type": "number",
					"description": "Dissolved oxygen concentration",
					"minimum": 0,
					"maximum": 20,
					"unit": "mg/L"
				},
				"saturation": {
					"type": "number",
					"description": "Oxygen saturation percentage",
					"minimum": 0,
					"maximum": 200,
					"unit": "%"
				},
				"temperature": {
					"type": "number",
					"description": "Water temperature",
					"minimum": 0,
					"maximum": 40,
					"unit": "°C"
				}
			},
			"required": ["dissolved_oxygen"]
		}`),
		UIConfiguration: json.RawMessage(`{
			"display": {
				"type": "multi_gauge",
				"primary": "dissolved_oxygen",
				"secondary": "saturation"
			}
		}`),
		Version:   "1.0.0",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: &[]time.Time{time.Now()}[0],
	},
	{
		ID:           uuid.MustParse("04567890-4444-4444-4444-000000000004"),
		SensorTypeID: uuid.MustParse("03456789-3333-3333-3333-000000000004"), // PM2.5 Sensor
		Name:         "Particulate Matter Measurement",
		Description:  stringPtr("PM2.5 and PM10 measurement for air quality monitoring"),
		PropertiesSchema: json.RawMessage(`{
			"type": "object",
			"properties": {
				"pm25": {
					"type": "number",
					"description": "PM2.5 concentration",
					"minimum": 0,
					"maximum": 500,
					"unit": "μg/m³"
				},
				"pm10": {
					"type": "number",
					"description": "PM10 concentration", 
					"minimum": 0,
					"maximum": 500,
					"unit": "μg/m³"
				}
			},
			"required": ["pm25", "pm10"]
		}`),
		UIConfiguration: json.RawMessage(`{
			"display": {
				"type": "bar_chart",
				"categories": ["pm25", "pm10"],
				"thresholds": {
					"pm25": [12, 35, 55, 150],
					"pm10": [54, 154, 254, 354]
				}
			}
		}`),
		Version:   "1.0.0",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: &[]time.Time{time.Now()}[0],
	},
	{
		ID:           uuid.MustParse("04567890-4444-4444-4444-000000000005"),
		SensorTypeID: uuid.MustParse("03456789-3333-3333-3333-000000000005"), // CO2 Sensor
		Name:         "Carbon Dioxide Measurement",
		Description:  stringPtr("CO2 concentration measurement for air quality and indoor environment"),
		PropertiesSchema: json.RawMessage(`{
			"type": "object",
			"properties": {
				"co2": {
					"type": "number",
					"description": "CO2 concentration",
					"minimum": 300,
					"maximum": 10000,
					"unit": "ppm"
				},
				"temperature": {
					"type": "number",
					"description": "Ambient temperature",
					"minimum": -40,
					"maximum": 70,
					"unit": "°C"
				},
				"humidity": {
					"type": "number",
					"description": "Relative humidity",
					"minimum": 0,
					"maximum": 100,
					"unit": "%RH"
				}
			},
			"required": ["co2"]
		}`),
		UIConfiguration: json.RawMessage(`{
			"display": {
				"type": "gauge",
				"min": 300,
				"max": 5000,
				"ranges": [
					{"min": 300, "max": 600, "color": "green", "label": "Good"},
					{"min": 600, "max": 1000, "color": "yellow", "label": "Moderate"},
					{"min": 1000, "max": 2500, "color": "orange", "label": "Poor"},
					{"min": 2500, "max": 5000, "color": "red", "label": "Hazardous"}
				]
			}
		}`),
		Version:   "1.0.0",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: &[]time.Time{time.Now()}[0],
	},
	{
		ID:           uuid.MustParse("04567890-4444-4444-4444-000000000006"),
		SensorTypeID: uuid.MustParse("03456789-3333-3333-3333-000000000006"), // Temperature Sensor
		Name:         "Industrial Temperature Measurement",
		Description:  stringPtr("High-precision temperature measurement for industrial processes"),
		PropertiesSchema: json.RawMessage(`{
			"type": "object",
			"properties": {
				"temperature": {
					"type": "number",
					"description": "Process temperature",
					"minimum": -200,
					"maximum": 850,
					"unit": "°C"
				},
				"accuracy": {
					"type": "number",
					"description": "Measurement accuracy",
					"minimum": 0.1,
					"maximum": 5.0,
					"unit": "°C"
				}
			},
			"required": ["temperature"]
		}`),
		UIConfiguration: json.RawMessage(`{
			"display": {
				"type": "line_chart",
				"trending": true,
				"alerts": true
			}
		}`),
		Version:   "1.0.0",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: &[]time.Time{time.Now()}[0],
	},
	{
		ID:           uuid.MustParse("04567890-4444-4444-4444-000000000007"),
		SensorTypeID: uuid.MustParse("03456789-3333-3333-3333-000000000008"), // Pressure Sensor
		Name:         "Process Pressure Measurement",
		Description:  stringPtr("Pressure measurement for industrial process control"),
		PropertiesSchema: json.RawMessage(`{
			"type": "object",
			"properties": {
				"pressure": {
					"type": "number",
					"description": "Process pressure",
					"minimum": 0,
					"maximum": 400,
					"unit": "bar"
				},
				"temperature": {
					"type": "number",
					"description": "Process temperature",
					"minimum": -40,
					"maximum": 150,
					"unit": "°C"
				}
			},
			"required": ["pressure"]
		}`),
		UIConfiguration: json.RawMessage(`{
			"display": {
				"type": "gauge",
				"min": 0,
				"max": 100,
				"critical_high": 90,
				"critical_low": 5
			}
		}`),
		Version:   "1.0.0",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: &[]time.Time{time.Now()}[0],
	},
	{
		ID:           uuid.MustParse("04567890-4444-4444-4444-000000000008"),
		SensorTypeID: uuid.MustParse("03456789-3333-3333-3333-000000000009"), // Flow Sensor
		Name:         "Fluid Flow Measurement",
		Description:  stringPtr("Flow rate measurement for liquid and gas applications"),
		PropertiesSchema: json.RawMessage(`{
			"type": "object",
			"properties": {
				"flow_rate": {
					"type": "number",
					"description": "Volumetric flow rate",
					"minimum": 0,
					"maximum": 10000,
					"unit": "L/min"
				},
				"totalizer": {
					"type": "number",
					"description": "Cumulative flow total",
					"minimum": 0,
					"unit": "L"
				},
				"temperature": {
					"type": "number",
					"description": "Fluid temperature",
					"minimum": -20,
					"maximum": 120,
					"unit": "°C"
				}
			},
			"required": ["flow_rate"]
		}`),
		UIConfiguration: json.RawMessage(`{
			"display": {
				"type": "combined",
				"primary": {"type": "gauge", "field": "flow_rate"},
				"secondary": {"type": "counter", "field": "totalizer"}
			}
		}`),
		Version:   "1.0.0",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: &[]time.Time{time.Now()}[0],
	},
	{
		ID:           uuid.MustParse("04567890-4444-4444-4444-000000000009"),
		SensorTypeID: uuid.MustParse("03456789-3333-3333-3333-000000000010"), // Vibration Sensor
		Name:         "Machinery Vibration Measurement",
		Description:  stringPtr("Vibration analysis for predictive maintenance"),
		PropertiesSchema: json.RawMessage(`{
			"type": "object",
			"properties": {
				"rms_velocity": {
					"type": "number",
					"description": "RMS velocity",
					"minimum": 0,
					"maximum": 100,
					"unit": "mm/s"
				},
				"peak_acceleration": {
					"type": "number",
					"description": "Peak acceleration",
					"minimum": 0,
					"maximum": 100,
					"unit": "g"
				},
				"frequency": {
					"type": "number",
					"description": "Dominant frequency",
					"minimum": 0,
					"maximum": 10000,
					"unit": "Hz"
				},
				"temperature": {
					"type": "number",
					"description": "Bearing temperature",
					"minimum": -20,
					"maximum": 200,
					"unit": "°C"
				}
			},
			"required": ["rms_velocity"]
		}`),
		UIConfiguration: json.RawMessage(`{
			"display": {
				"type": "spectrum",
				"time_domain": true,
				"frequency_domain": true,
				"alerts": ["rms_velocity", "peak_acceleration"]
			}
		}`),
		Version:   "1.0.0",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: &[]time.Time{time.Now()}[0],
	},
	{
		ID:           uuid.MustParse("04567890-4444-4444-4444-000000000010"),
		SensorTypeID: uuid.MustParse("03456789-3333-3333-3333-000000000012"), // Power Meter
		Name:         "Electrical Power Measurement",
		Description:  stringPtr("Comprehensive electrical power and energy measurement"),
		PropertiesSchema: json.RawMessage(`{
			"type": "object",
			"properties": {
				"voltage": {
					"type": "number",
					"description": "RMS voltage",
					"minimum": 0,
					"maximum": 500,
					"unit": "V"
				},
				"current": {
					"type": "number",
					"description": "RMS current",
					"minimum": 0,
					"maximum": 1000,
					"unit": "A"
				},
				"power": {
					"type": "number",
					"description": "Active power",
					"minimum": 0,
					"maximum": 500000,
					"unit": "W"
				},
				"energy": {
					"type": "number",
					"description": "Energy consumption",
					"minimum": 0,
					"unit": "kWh"
				},
				"power_factor": {
					"type": "number",
					"description": "Power factor",
					"minimum": 0,
					"maximum": 1,
					"unit": ""
				}
			},
			"required": ["voltage", "current", "power"]
		}`),
		UIConfiguration: json.RawMessage(`{
			"display": {
				"type": "dashboard",
				"widgets": [
					{"type": "gauge", "field": "voltage"},
					{"type": "gauge", "field": "current"},
					{"type": "line_chart", "field": "power"},
					{"type": "counter", "field": "energy"}
				]
			}
		}`),
		Version:   "1.0.0",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: &[]time.Time{time.Now()}[0],
	},
	{
		ID:           uuid.MustParse("04567890-4444-4444-4444-000000000011"),
		SensorTypeID: uuid.MustParse("03456789-3333-3333-3333-000000000011"), // Level Sensor
		Name:         "Level Measurement",
		Description:  stringPtr("Liquid level measurement for tanks and containers"),
		PropertiesSchema: json.RawMessage(`{
			"type": "object",
			"properties": {
				"level_value": {
					"type": "number",
					"description": "Liquid level height",
					"minimum": 0,
					"maximum": 10,
					"unit": "m"
				},
				"percentage": {
					"type": "number",
					"description": "Level as percentage of tank capacity",
					"minimum": 0,
					"maximum": 100,
					"unit": "%"
				},
				"volume": {
					"type": "number",
					"description": "Estimated volume",
					"minimum": 0,
					"unit": "L"
				}
			},
			"required": ["level_value"]
		}`),
		UIConfiguration: json.RawMessage(`{
			"display": {
				"type": "tank_gauge",
				"primary_field": "level_value",
				"secondary_field": "percentage",
				"thresholds": {
					"low": 10,
					"high": 90
				}
			}
		}`),
		Version:   "1.0.0",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: &[]time.Time{time.Now()}[0],
	},
	{
		ID:           uuid.MustParse("04567890-4444-4444-4444-000000000012"),
		SensorTypeID: uuid.MustParse("03456789-3333-3333-3333-000000000013"), // Gas Sensor
		Name:         "Gas Detection Measurement",
		Description:  stringPtr("Gas concentration measurement for safety monitoring"),
		PropertiesSchema: json.RawMessage(`{
			"type": "object",
			"properties": {
				"gas_concentration": {
					"type": "number",
					"description": "Gas concentration level",
					"minimum": 0,
					"maximum": 10000,
					"unit": "ppm"
				},
				"temperature": {
					"type": "number",
					"description": "Sensor temperature",
					"minimum": -40,
					"maximum": 85,
					"unit": "°C"
				},
				"humidity": {
					"type": "number",
					"description": "Relative humidity",
					"minimum": 0,
					"maximum": 100,
					"unit": "%"
				}
			},
			"required": ["gas_concentration"]
		}`),
		UIConfiguration: json.RawMessage(`{
			"display": {
				"type": "alert_dashboard",
				"primary_field": "gas_concentration",
				"thresholds": {
					"warning": 50,
					"critical": 100
				},
				"alerts": true
			}
		}`),
		Version:   "1.0.0",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: &[]time.Time{time.Now()}[0],
	},
}

func (s *SensorMeasurementTypeSeeder) Seed(ctx context.Context) error {
	log.Println("Starting SensorMeasurementType seeding...")

	// Check if measurement types already exist
	existingTypes, _, err := s.repo.List(ctx, 0, 1)
	if err != nil {
		return fmt.Errorf("failed to check existing measurement types: %w", err)
	}

	if len(existingTypes) > 0 {
		log.Printf("Sensor measurement types already exist, skipping seeding")
		return nil
	}

	// Insert measurement types
	for i, measurementType := range sensorMeasurementTypeData {
		err := s.repo.Create(ctx, &measurementType)
		if err != nil {
			return fmt.Errorf("failed to create measurement type %d (%s): %w", i+1, measurementType.Name, err)
		}
		log.Printf("Created measurement type: %s", measurementType.Name)
	}

	log.Printf("Successfully seeded %d sensor measurement types", len(sensorMeasurementTypeData))
	return nil
}

func (s *SensorMeasurementTypeSeeder) GetMeasurementTypeIDs() []uuid.UUID {
	var ids []uuid.UUID
	for _, measurementType := range sensorMeasurementTypeData {
		ids = append(ids, measurementType.ID)
	}
	return ids
}

func (s *SensorMeasurementTypeSeeder) GetMeasurementTypeBySensorType(sensorTypeID uuid.UUID) uuid.UUID {
	for _, measurementType := range sensorMeasurementTypeData {
		if measurementType.SensorTypeID == sensorTypeID {
			return measurementType.ID
		}
	}
	return uuid.Nil
}
