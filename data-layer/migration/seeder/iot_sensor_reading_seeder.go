package seeder

import (
	"be-lecsens/asset_management/data-layer/entity"
	"context"
	"database/sql"
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

// Use the actual constants from asset_sensor_seeder.go and sensor type seeder
// No need to redeclare them here as they will be available from the same package

// IoTSensorReadingSeeder handles seeding IoT sensor readings
type IoTSensorReadingSeeder struct {
	db *sql.DB
}

// NewIoTSensorReadingSeeder creates a new IoT sensor reading seeder
func NewIoTSensorReadingSeeder(db *sql.DB) *IoTSensorReadingSeeder {
	return &IoTSensorReadingSeeder{db: db}
}

// SensorReadingConfig defines the configuration for generating sensor readings
type SensorReadingConfig struct {
	AssetSensorID   uuid.UUID
	SensorTypeID    uuid.UUID
	LocationID      uuid.UUID // Fixed location ID from location seeder
	LocationName    string    // Location name for efficiency
	MeasurementType string
	BaseValue       float64
	Variation       float64
	TrendDirection  float64 // Positive for upward trend, negative for downward
	Unit            string
	MinValue        float64
	MaxValue        float64
}

// getSensorReadingConfigs returns configurations for generating realistic sensor data (1 per day)
func (s *IoTSensorReadingSeeder) getSensorReadingConfigs() []SensorReadingConfig {
	return []SensorReadingConfig{
		// Water Quality pH Readings - Jakarta Pusat
		{
			AssetSensorID:   AssetSensorPHProbeID,
			SensorTypeID:    SensorTypePHID,
			LocationID:      uuid.MustParse("01234567-1111-1111-1111-000000000001"), // Jakarta Pusat
			LocationName:    "Jakarta Pusat",
			MeasurementType: "water_quality",
			BaseValue:       7.2,
			Variation:       0.3,
			TrendDirection:  0.01,
			Unit:            "pH",
			MinValue:        6.0,
			MaxValue:        8.5,
		},
		// Water Turbidity Readings - Jakarta Selatan
		{
			AssetSensorID:   AssetSensorTurbidityID,
			SensorTypeID:    SensorTypeTurbidityID,
			LocationID:      uuid.MustParse("01234567-1111-1111-1111-000000000002"), // Jakarta Selatan
			LocationName:    "Jakarta Selatan",
			MeasurementType: "turbidity",
			BaseValue:       5.5,
			Variation:       2.0,
			TrendDirection:  -0.02,
			Unit:            "NTU",
			MinValue:        0.1,
			MaxValue:        50.0,
		},
		// Dissolved Oxygen Readings - Jakarta Utara
		{
			AssetSensorID:   AssetSensorDOProbeID,
			SensorTypeID:    SensorTypeDOID,
			LocationID:      uuid.MustParse("01234567-1111-1111-1111-000000000003"), // Jakarta Utara
			LocationName:    "Jakarta Utara",
			MeasurementType: "dissolved_oxygen",
			BaseValue:       6.5,
			Variation:       0.8,
			TrendDirection:  0.1,
			Unit:            "mg/L",
			MinValue:        4.0,
			MaxValue:        10.0,
		},
		// Water Temperature Readings - Bandung
		{
			AssetSensorID:   AssetSensorColdStorageTempID,
			SensorTypeID:    SensorTypeTemperatureID,
			LocationID:      uuid.MustParse("01234567-1111-1111-1111-000000000004"), // Bandung
			LocationName:    "Bandung",
			MeasurementType: "water_temperature",
			BaseValue:       25.0,
			Variation:       2.0,
			TrendDirection:  0.5,
			Unit:            "°C",
			MinValue:        15.0,
			MaxValue:        35.0,
		},
		// Water Flow Rate Readings - Surabaya
		{
			AssetSensorID:   AssetSensorFlowRateID,
			SensorTypeID:    SensorTypeFlowID,
			LocationID:      uuid.MustParse("01234567-1111-1111-1111-000000000005"), // Surabaya
			LocationName:    "Surabaya",
			MeasurementType: "flow_rate",
			BaseValue:       85.0,
			Variation:       15.0,
			TrendDirection:  0.02,
			Unit:            "L/min",
			MinValue:        20.0,
			MaxValue:        200.0,
		},
		// Water Pressure Readings - Yogyakarta
		{
			AssetSensorID:   AssetSensorPressure1ID,
			SensorTypeID:    SensorTypePressureID,
			LocationID:      uuid.MustParse("01234567-1111-1111-1111-000000000006"), // Yogyakarta
			LocationName:    "Yogyakarta",
			MeasurementType: "water_pressure",
			BaseValue:       4.2,
			Variation:       0.8,
			TrendDirection:  0.01,
			Unit:            "bar",
			MinValue:        1.0,
			MaxValue:        10.0,
		},
		// Water Level Readings - Malang
		{
			AssetSensorID:   AssetSensorLevelTank1ID,
			SensorTypeID:    SensorTypeLevelID,
			LocationID:      uuid.MustParse("01234567-1111-1111-1111-000000000007"), // Malang
			LocationName:    "Malang",
			MeasurementType: "water_level",
			BaseValue:       2.8,
			Variation:       0.5,
			TrendDirection:  -0.05,
			Unit:            "m",
			MinValue:        0.2,
			MaxValue:        6.0,
		},
		// Water Conductivity Readings - Medan
		{
			AssetSensorID:   AssetSensorWaterConductID,
			SensorTypeID:    SensorTypeConductivityID,
			LocationID:      uuid.MustParse("01234567-1111-1111-1111-000000000008"), // Medan
			LocationName:    "Medan",
			MeasurementType: "conductivity",
			BaseValue:       450.0,
			Variation:       50.0,
			TrendDirection:  0.001,
			Unit:            "μS/cm",
			MinValue:        100.0,
			MaxValue:        1000.0,
		},
		// Water Chlorine Readings - Denpasar
		{
			AssetSensorID:   AssetSensorChlorineID,
			SensorTypeID:    SensorTypeChlorineID,
			LocationID:      uuid.MustParse("01234567-1111-1111-1111-000000000009"), // Denpasar
			LocationName:    "Denpasar",
			MeasurementType: "chlorine",
			BaseValue:       2.0,
			Variation:       0.3,
			TrendDirection:  -0.002,
			Unit:            "mg/L",
			MinValue:        0.5,
			MaxValue:        4.0,
		},
		// Water TDS Readings - Makassar
		{
			AssetSensorID:   AssetSensorTDSID,
			SensorTypeID:    SensorTypeTDSID,
			LocationID:      uuid.MustParse("01234567-1111-1111-1111-000000000010"), // Makassar
			LocationName:    "Makassar",
			MeasurementType: "tds",
			BaseValue:       320.0,
			Variation:       45.0,
			TrendDirection:  0.1,
			Unit:            "ppm",
			MinValue:        50.0,
			MaxValue:        800.0,
		},
		// Water ORP Readings - Jakarta Pusat (additional sensor)
		{
			AssetSensorID:   AssetSensorORPID,
			SensorTypeID:    SensorTypeORPID,
			LocationID:      uuid.MustParse("01234567-1111-1111-1111-000000000001"), // Jakarta Pusat
			LocationName:    "Jakarta Pusat",
			MeasurementType: "orp",
			BaseValue:       650.0,
			Variation:       50.0,
			TrendDirection:  0.02,
			Unit:            "mV",
			MinValue:        0.0,
			MaxValue:        1000.0,
		},
	}
}

// generateSensorValue generates a realistic sensor value with trends and variations
func (s *IoTSensorReadingSeeder) generateSensorValue(config SensorReadingConfig, timestamp time.Time, baseTime time.Time) float64 {
	// Calculate time-based trend
	hoursSinceBase := timestamp.Sub(baseTime).Hours()
	trend := config.TrendDirection * hoursSinceBase

	// Add daily cycle (some sensors have daily patterns)
	hour := float64(timestamp.Hour())
	dailyCycle := 0.0

	switch config.MeasurementType {
	case "temperature":
		// Temperature typically peaks in afternoon, lowest at dawn
		dailyCycle = 2.0 * math.Sin((hour-6)*math.Pi/12)
	case "air_quality":
		// Air quality often worse during rush hours
		if hour >= 7 && hour <= 9 || hour >= 17 && hour <= 19 {
			dailyCycle = config.Variation * 0.3
		}
	case "power":
		// Power consumption higher during work hours
		if hour >= 8 && hour <= 18 {
			dailyCycle = config.Variation * 0.4
		}
	}

	// Add random variation
	randomVariation := (rand.Float64() - 0.5) * config.Variation

	// Calculate final value
	value := config.BaseValue + trend + dailyCycle + randomVariation

	// Apply bounds
	if value < config.MinValue {
		value = config.MinValue + rand.Float64()*(config.BaseValue-config.MinValue)*0.1
	}
	if value > config.MaxValue {
		value = config.MaxValue - rand.Float64()*(config.MaxValue-config.BaseValue)*0.1
	}

	return math.Round(value*100) / 100 // Round to 2 decimal places
}

// generateReadingsForConfig generates readings for a specific sensor configuration (1 per day)
func (s *IoTSensorReadingSeeder) generateReadingsForConfig(ctx context.Context, config SensorReadingConfig) error {
	tenantID := GetDefaultTenantID()
	now := time.Now()

	// Start from January 1, 2022
	startDate := time.Date(2022, 1, 1, 12, 0, 0, 0, time.UTC) // 12:00 PM for consistent daily reading time

	// Calculate total days from start date to today
	totalDays := int(now.Sub(startDate).Hours() / 24)

	log.Printf("Generating %d daily readings for sensor %s from %s to %s",
		totalDays, config.MeasurementType, startDate.Format("2006-01-02"), now.Format("2006-01-02"))

	createdCount := 0
	batchSize := 100
	readings := make([]entity.IoTSensorReading, 0, batchSize)

	// Generate one reading per day
	for dayOffset := 0; dayOffset < totalDays; dayOffset++ {
		// Calculate timestamp for this reading (one per day at 12:00 PM)
		readingTime := startDate.AddDate(0, 0, dayOffset)
		if readingTime.After(now) {
			break // Don't generate future readings
		}

		// Generate sensor value
		value := s.generateSensorValue(config, readingTime, startDate)

		reading := entity.IoTSensorReading{
			ID:              uuid.New(),
			TenantID:        tenantID,
			AssetSensorID:   config.AssetSensorID,
			SensorTypeID:    config.SensorTypeID,
			LocationID:      &config.LocationID,   // Add LocationID from config
			LocationName:    &config.LocationName, // Add LocationName from config
			MeasurementType: config.MeasurementType,
			MeasurementUnit: &config.Unit,
			NumericValue:    &value,
			DataSource:      stringPtr("seeder"),
			ReadingTime:     readingTime,
			CreatedAt:       readingTime,
		}

		readings = append(readings, reading)

		// Insert batch when full
		if len(readings) >= batchSize {
			if err := s.insertReadingsBatch(ctx, readings); err != nil {
				log.Printf("Error inserting readings batch: %v", err)
				continue
			}
			createdCount += len(readings)
			readings = readings[:0] // Clear slice but keep capacity
		}

		// Log progress every 100 days
		if dayOffset > 0 && dayOffset%100 == 0 {
			log.Printf("Generated %d/%d daily readings for %s", dayOffset, totalDays, config.MeasurementType)
		}
	}

	// Insert remaining readings
	if len(readings) > 0 {
		if err := s.insertReadingsBatch(ctx, readings); err != nil {
			log.Printf("Error inserting final readings batch: %v", err)
		} else {
			createdCount += len(readings)
		}
	}

	log.Printf("Generated %d daily readings for %s sensor", createdCount, config.MeasurementType)
	return nil
}

// insertReadingsBatch inserts a batch of readings
func (s *IoTSensorReadingSeeder) insertReadingsBatch(ctx context.Context, readings []entity.IoTSensorReading) error {
	if len(readings) == 0 {
		return nil
	}

	query := `
		INSERT INTO iot_sensor_readings (
			id, tenant_id, asset_sensor_id, sensor_type_id, location_id, location_name,
			measurement_type, measurement_unit, numeric_value, 
			data_source, reading_time, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, reading := range readings {
		_, err := stmt.ExecContext(ctx,
			reading.ID,
			reading.TenantID,
			reading.AssetSensorID,
			reading.SensorTypeID,
			reading.LocationID,
			reading.LocationName,
			reading.MeasurementType,
			reading.MeasurementUnit,
			reading.NumericValue,
			reading.DataSource,
			reading.ReadingTime,
			reading.CreatedAt,
		)
		if err != nil {
			// Provide more detailed error information for foreign key violations
			if fmt.Sprintf("%v", err) == "pq: insert or update on table \"iot_sensor_readings\" violates foreign key constraint \"fk_iot_readings_asset_sensor_id\"" {
				return fmt.Errorf("foreign key violation: asset_sensor_id %s does not exist in asset_sensors table", reading.AssetSensorID)
			}
			return fmt.Errorf("failed to insert reading for asset_sensor_id %s: %w", reading.AssetSensorID, err)
		}
	}

	return tx.Commit()
}

// Helper functions for generating realistic measurement data

// clearExistingReadings clears existing readings for testing (optional)
func (s *IoTSensorReadingSeeder) clearExistingReadings(ctx context.Context) error {
	query := `DELETE FROM iot_sensor_readings WHERE tenant_id = $1`
	_, err := s.db.ExecContext(ctx, query, GetDefaultTenantID())
	if err != nil {
		return fmt.Errorf("failed to clear existing readings: %w", err)
	}
	log.Println("Cleared existing IoT sensor readings")
	return nil
}

// verifyAssetSensorsExist checks if all required asset sensors exist in the database
func (s *IoTSensorReadingSeeder) verifyAssetSensorsExist(ctx context.Context) error {
	configs := s.getSensorReadingConfigs()

	log.Println("Verifying asset sensors exist before generating readings...")

	query := `SELECT id FROM asset_sensors WHERE id = $1`

	for _, config := range configs {
		var existingID uuid.UUID
		err := s.db.QueryRowContext(ctx, query, config.AssetSensorID).Scan(&existingID)
		if err != nil {
			if err == sql.ErrNoRows {
				return fmt.Errorf("asset sensor with ID %s does not exist (measurement type: %s)",
					config.AssetSensorID, config.MeasurementType)
			}
			return fmt.Errorf("error checking asset sensor %s: %w", config.AssetSensorID, err)
		}
		log.Printf("✓ Asset sensor %s exists (measurement type: %s)", config.AssetSensorID, config.MeasurementType)
	}

	log.Println("All required asset sensors verified to exist")
	return nil
}

// SeedReadings generates and seeds IoT sensor readings from Jan 1, 2022 to present (1 per day)
func (s *IoTSensorReadingSeeder) SeedReadings(ctx context.Context, clearExisting bool) error {
	log.Println("Starting IoT sensor reading seeding from January 1, 2022 to present day...")

	// First verify all required asset sensors exist
	if err := s.verifyAssetSensorsExist(ctx); err != nil {
		return fmt.Errorf("asset sensor verification failed: %w", err)
	}

	if clearExisting {
		if err := s.clearExistingReadings(ctx); err != nil {
			log.Printf("Warning: Failed to clear existing readings: %v", err)
		}
	}

	configs := s.getSensorReadingConfigs()
	totalConfigs := len(configs)

	for i, config := range configs {
		log.Printf("Processing sensor config %d/%d: %s", i+1, totalConfigs, config.MeasurementType)

		if err := s.generateReadingsForConfig(ctx, config); err != nil {
			log.Printf("Error generating readings for %s: %v", config.MeasurementType, err)
			continue
		}
	}

	log.Println("IoT sensor reading seeding completed successfully!")
	return nil
}

// SeedReadingsLimited generates a limited set of readings for quick testing (recent days only)
func (s *IoTSensorReadingSeeder) SeedReadingsLimited(ctx context.Context, daysBack int) error {
	log.Printf("Starting limited IoT sensor reading seeding (last %d days)...", daysBack)

	// First verify all required asset sensors exist
	if err := s.verifyAssetSensorsExist(ctx); err != nil {
		return fmt.Errorf("asset sensor verification failed: %w", err)
	}

	configs := s.getSensorReadingConfigs()
	now := time.Now()
	tenantID := GetDefaultTenantID()

	// Start from 'daysBack' days ago
	startDate := now.AddDate(0, 0, -daysBack)

	for _, config := range configs {
		log.Printf("Generating readings for %s sensor (last %d days)", config.MeasurementType, daysBack)

		readings := make([]entity.IoTSensorReading, 0, daysBack)

		// Generate one reading per day for the last 'daysBack' days
		for dayOffset := 0; dayOffset < daysBack; dayOffset++ {
			readingTime := startDate.AddDate(0, 0, dayOffset)
			if readingTime.After(now) {
				break
			}

			value := s.generateSensorValue(config, readingTime, startDate)

			reading := entity.IoTSensorReading{
				ID:              uuid.New(),
				TenantID:        tenantID,
				AssetSensorID:   config.AssetSensorID,
				SensorTypeID:    config.SensorTypeID,
				LocationID:      &config.LocationID,   // Add LocationID from config
				LocationName:    &config.LocationName, // Add LocationName from config
				MeasurementType: config.MeasurementType,
				MeasurementUnit: &config.Unit,
				NumericValue:    &value,
				DataSource:      stringPtr("seeder"),
				ReadingTime:     readingTime,
				CreatedAt:       readingTime,
			}

			readings = append(readings, reading)
		}

		if err := s.insertReadingsBatch(ctx, readings); err != nil {
			log.Printf("Error inserting readings for %s: %v", config.MeasurementType, err)
			continue
		}

		log.Printf("Generated %d readings for %s sensor", len(readings), config.MeasurementType)
	}

	log.Println("Limited IoT sensor reading seeding completed!")
	return nil
}
