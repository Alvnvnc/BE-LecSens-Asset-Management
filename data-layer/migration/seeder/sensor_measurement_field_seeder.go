package seeder

import (
	"be-lecsens/asset_management/data-layer/repository"
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"
)

type SensorMeasurementFieldSeeder struct {
	db *sql.DB
}

func NewSensorMeasurementFieldSeeder(db *sql.DB) *SensorMeasurementFieldSeeder {
	return &SensorMeasurementFieldSeeder{
		db: db,
	}
}

// Measurement Type IDs (from sensor_measurement_type_seeder.go)
var (
	MeasurementTypeWaterQualityID    = uuid.MustParse("04567890-4444-4444-4444-000000000001") // pH Sensor
	MeasurementTypeTurbidityID       = uuid.MustParse("04567890-4444-4444-4444-000000000002") // Turbidity Sensor
	MeasurementTypeDissolvedOxygenID = uuid.MustParse("04567890-4444-4444-4444-000000000003") // Dissolved Oxygen Sensor
	MeasurementTypeAirQualityID      = uuid.MustParse("04567890-4444-4444-4444-000000000004") // PM2.5 Sensor
	MeasurementTypeCO2ID             = uuid.MustParse("04567890-4444-4444-4444-000000000005") // CO2 Sensor
	MeasurementTypeTemperatureID     = uuid.MustParse("04567890-4444-4444-4444-000000000006") // Temperature Sensor
	MeasurementTypePressureID        = uuid.MustParse("04567890-4444-4444-4444-000000000007") // Pressure Sensor
	MeasurementTypeFlowID            = uuid.MustParse("04567890-4444-4444-4444-000000000008") // Flow Sensor
	MeasurementTypeVibrationID       = uuid.MustParse("04567890-4444-4444-4444-000000000009") // Vibration Sensor
	MeasurementTypePowerID           = uuid.MustParse("04567890-4444-4444-4444-000000000010") // Power Meter
	// Additional measurement types for threshold seeder compatibility
	MeasurementTypeLevelID        = uuid.MustParse("04567890-4444-4444-4444-000000000011") // Level Sensor
	MeasurementTypeGasDetectionID = uuid.MustParse("04567890-4444-4444-4444-000000000012") // Gas Detection Sensor
)

// Predefined field IDs for consistent seeding
var (
	// pH Sensor Fields
	FieldPhValueID       = uuid.MustParse("11111111-1111-1111-1111-000000000001")
	FieldPhTemperatureID = uuid.MustParse("11111111-1111-1111-1111-000000000002")
	FieldPhStatusID      = uuid.MustParse("11111111-1111-1111-1111-000000000003")

	// Turbidity Sensor Fields
	FieldTurbidityValueID   = uuid.MustParse("11111111-1111-1111-1111-000000000004")
	FieldTurbidityUnitID    = uuid.MustParse("11111111-1111-1111-1111-000000000005")
	FieldTurbidityQualityID = uuid.MustParse("11111111-1111-1111-1111-000000000006")

	// PM2.5 Sensor Fields
	FieldPm25ValueID     = uuid.MustParse("11111111-1111-1111-1111-000000000007")
	FieldCo2ValueID      = uuid.MustParse("11111111-1111-1111-1111-000000000008")
	FieldAirQualityID    = uuid.MustParse("11111111-1111-1111-1111-000000000009")
	FieldHumidityValueID = uuid.MustParse("11111111-1111-1111-1111-000000000010")

	// Temperature Sensor Fields
	FieldTemperatureValueID = uuid.MustParse("11111111-1111-1111-1111-000000000011")
	FieldTempAlarmStatusID  = uuid.MustParse("11111111-1111-1111-1111-000000000012")

	// Vibration Sensor Fields
	FieldVibrationAmplitudeID = uuid.MustParse("11111111-1111-1111-1111-000000000013")
	FieldVibrationFrequencyID = uuid.MustParse("11111111-1111-1111-1111-000000000014")
	FieldVibrationLevelID     = uuid.MustParse("11111111-1111-1111-1111-000000000015")

	// Flow Sensor Fields
	FieldFlowRateID    = uuid.MustParse("11111111-1111-1111-1111-000000000016")
	FieldTotalVolumeID = uuid.MustParse("11111111-1111-1111-1111-000000000017")
	FieldFlowStatusID  = uuid.MustParse("11111111-1111-1111-1111-000000000018")

	// Pressure Sensor Fields
	FieldPressureValueID = uuid.MustParse("11111111-1111-1111-1111-000000000019")
	FieldPressureUnitID  = uuid.MustParse("11111111-1111-1111-1111-000000000020")
	FieldPressureAlarmID = uuid.MustParse("11111111-1111-1111-1111-000000000021")

	// Level Sensor Fields
	FieldLevelValueID   = uuid.MustParse("11111111-1111-1111-1111-000000000022")
	FieldLevelPercentID = uuid.MustParse("11111111-1111-1111-1111-000000000023")
	FieldLevelAlarmID   = uuid.MustParse("11111111-1111-1111-1111-000000000024")

	// Power Sensor Fields
	FieldPowerValueID   = uuid.MustParse("11111111-1111-1111-1111-000000000025")
	FieldVoltageValueID = uuid.MustParse("11111111-1111-1111-1111-000000000026")
	FieldCurrentValueID = uuid.MustParse("11111111-1111-1111-1111-000000000027")
	FieldPowerFactorID  = uuid.MustParse("11111111-1111-1111-1111-000000000028")

	// Gas Detection Sensor Fields
	FieldGasConcentrationID = uuid.MustParse("11111111-1111-1111-1111-000000000029")
	FieldGasTypeID          = uuid.MustParse("11111111-1111-1111-1111-000000000030")
	FieldGasAlarmID         = uuid.MustParse("11111111-1111-1111-1111-000000000031")
)

// getMeasurementFields returns all measurement fields to be seeded
func (s *SensorMeasurementFieldSeeder) getMeasurementFields() []*repository.SensorMeasurementField {
	now := time.Now()

	return []*repository.SensorMeasurementField{
		// pH Measurement Fields (Water Quality Measurement)
		{
			ID:                      FieldPhValueID,
			SensorMeasurementTypeID: MeasurementTypeWaterQualityID,
			Name:                    "ph_value",
			Label:                   "pH Value",
			Description:             repository.NullStringFromPtr(stringPtr("pH measurement value (0-14 scale)")),
			DataType:                "number",
			Required:                true,
			Unit:                    repository.NullStringFromPtr(stringPtr("pH")),
			Min:                     repository.NullFloat64FromPtr(float64Ptr(0.0)),
			Max:                     repository.NullFloat64FromPtr(float64Ptr(14.0)),
			CreatedAt:               now,
			UpdatedAt:               now,
		},
		{
			ID:                      FieldPhTemperatureID,
			SensorMeasurementTypeID: MeasurementTypeWaterQualityID,
			Name:                    "temperature",
			Label:                   "Water Temperature",
			Description:             repository.NullStringFromPtr(stringPtr("Temperature compensation for pH measurement")),
			DataType:                "number",
			Required:                false,
			Unit:                    repository.NullStringFromPtr(stringPtr("°C")),
			Min:                     repository.NullFloat64FromPtr(float64Ptr(0.0)),
			Max:                     repository.NullFloat64FromPtr(float64Ptr(100.0)),
			CreatedAt:               now,
			UpdatedAt:               now,
		},
		{
			ID:                      FieldPhStatusID,
			SensorMeasurementTypeID: MeasurementTypeWaterQualityID,
			Name:                    "calibration_status",
			Label:                   "Calibration Status",
			Description:             repository.NullStringFromPtr(stringPtr("pH sensor calibration status")),
			DataType:                "string",
			Required:                false,
			Unit:                    repository.NullStringFromPtr(stringPtr("")),
			CreatedAt:               now,
			UpdatedAt:               now,
		},

		// Turbidity Measurement Fields
		{
			ID:                      FieldTurbidityValueID,
			SensorMeasurementTypeID: MeasurementTypeTurbidityID,
			Name:                    "turbidity_value",
			Label:                   "Turbidity Value",
			Description:             repository.NullStringFromPtr(stringPtr("Water turbidity measurement")),
			DataType:                "number",
			Required:                true,
			Unit:                    repository.NullStringFromPtr(stringPtr("NTU")),
			Min:                     repository.NullFloat64FromPtr(float64Ptr(0.0)),
			Max:                     repository.NullFloat64FromPtr(float64Ptr(4000.0)),
			CreatedAt:               now,
			UpdatedAt:               now,
		},
		{
			ID:                      FieldTurbidityUnitID,
			SensorMeasurementTypeID: MeasurementTypeTurbidityID,
			Name:                    "measurement_unit",
			Label:                   "Measurement Unit",
			Description:             repository.NullStringFromPtr(stringPtr("Unit of turbidity measurement")),
			DataType:                "string",
			Required:                false,
			Unit:                    repository.NullStringFromPtr(stringPtr("")),
			CreatedAt:               now,
			UpdatedAt:               now,
		},
		{
			ID:                      FieldTurbidityQualityID,
			SensorMeasurementTypeID: MeasurementTypeTurbidityID,
			Name:                    "water_quality",
			Label:                   "Water Quality Assessment",
			Description:             repository.NullStringFromPtr(stringPtr("Overall water quality based on turbidity")),
			DataType:                "string",
			Required:                false,
			Unit:                    repository.NullStringFromPtr(stringPtr("")),
			CreatedAt:               now,
			UpdatedAt:               now,
		},

		// Air Quality Measurement Fields
		{
			ID:                      FieldPm25ValueID,
			SensorMeasurementTypeID: MeasurementTypeAirQualityID,
			Name:                    "pm25_concentration",
			Label:                   "PM2.5 Concentration",
			Description:             repository.NullStringFromPtr(stringPtr("Fine particulate matter concentration")),
			DataType:                "number",
			Required:                true,
			Unit:                    repository.NullStringFromPtr(stringPtr("μg/m³")),
			Min:                     repository.NullFloat64FromPtr(float64Ptr(0.0)),
			Max:                     repository.NullFloat64FromPtr(float64Ptr(500.0)),
			CreatedAt:               now,
			UpdatedAt:               now,
		},
		{
			ID:                      FieldCo2ValueID,
			SensorMeasurementTypeID: MeasurementTypeAirQualityID,
			Name:                    "co2_concentration",
			Label:                   "CO2 Concentration",
			Description:             repository.NullStringFromPtr(stringPtr("Carbon dioxide concentration")),
			DataType:                "number",
			Required:                false,
			Unit:                    repository.NullStringFromPtr(stringPtr("ppm")),
			Min:                     repository.NullFloat64FromPtr(float64Ptr(300.0)),
			Max:                     repository.NullFloat64FromPtr(float64Ptr(10000.0)),
			CreatedAt:               now,
			UpdatedAt:               now,
		},
		{
			ID:                      FieldAirQualityID,
			SensorMeasurementTypeID: MeasurementTypeAirQualityID,
			Name:                    "air_quality_index",
			Label:                   "Air Quality Index",
			Description:             repository.NullStringFromPtr(stringPtr("Overall air quality assessment")),
			DataType:                "number",
			Required:                false,
			Unit:                    repository.NullStringFromPtr(stringPtr("AQI")),
			Min:                     repository.NullFloat64FromPtr(float64Ptr(0.0)),
			Max:                     repository.NullFloat64FromPtr(float64Ptr(500.0)),
			CreatedAt:               now,
			UpdatedAt:               now,
		},
		{
			ID:                      FieldHumidityValueID,
			SensorMeasurementTypeID: MeasurementTypeCO2ID,
			Name:                    "humidity",
			Label:                   "Relative Humidity",
			Description:             repository.NullStringFromPtr(stringPtr("Ambient relative humidity")),
			DataType:                "number",
			Required:                false,
			Unit:                    repository.NullStringFromPtr(stringPtr("%RH")),
			Min:                     repository.NullFloat64FromPtr(float64Ptr(0.0)),
			Max:                     repository.NullFloat64FromPtr(float64Ptr(100.0)),
			CreatedAt:               now,
			UpdatedAt:               now,
		},

		// Temperature Measurement Fields
		{
			ID:                      FieldTemperatureValueID,
			SensorMeasurementTypeID: MeasurementTypeTemperatureID,
			Name:                    "temperature_value",
			Label:                   "Temperature",
			Description:             repository.NullStringFromPtr(stringPtr("Temperature measurement")),
			DataType:                "number",
			Required:                true,
			Unit:                    repository.NullStringFromPtr(stringPtr("°C")),
			Min:                     repository.NullFloat64FromPtr(float64Ptr(-200.0)),
			Max:                     repository.NullFloat64FromPtr(float64Ptr(850.0)),
			CreatedAt:               now,
			UpdatedAt:               now,
		},
		{
			ID:                      FieldTempAlarmStatusID,
			SensorMeasurementTypeID: MeasurementTypeTemperatureID,
			Name:                    "alarm_status",
			Label:                   "Temperature Alarm Status",
			Description:             repository.NullStringFromPtr(stringPtr("Temperature alarm status indicator")),
			DataType:                "boolean",
			Required:                false,
			Unit:                    repository.NullStringFromPtr(stringPtr("")),
			CreatedAt:               now,
			UpdatedAt:               now,
		},

		// Vibration Measurement Fields
		{
			ID:                      FieldVibrationAmplitudeID,
			SensorMeasurementTypeID: MeasurementTypeVibrationID,
			Name:                    "rms_velocity",
			Label:                   "RMS Velocity",
			Description:             repository.NullStringFromPtr(stringPtr("RMS velocity vibration measurement")),
			DataType:                "number",
			Required:                true,
			Unit:                    repository.NullStringFromPtr(stringPtr("mm/s")),
			Min:                     repository.NullFloat64FromPtr(float64Ptr(0.0)),
			Max:                     repository.NullFloat64FromPtr(float64Ptr(100.0)),
			CreatedAt:               now,
			UpdatedAt:               now,
		},
		{
			ID:                      FieldVibrationFrequencyID,
			SensorMeasurementTypeID: MeasurementTypeVibrationID,
			Name:                    "peak_acceleration",
			Label:                   "Peak Acceleration",
			Description:             repository.NullStringFromPtr(stringPtr("Peak acceleration measurement")),
			DataType:                "number",
			Required:                false,
			Unit:                    repository.NullStringFromPtr(stringPtr("g")),
			Min:                     repository.NullFloat64FromPtr(float64Ptr(0.0)),
			Max:                     repository.NullFloat64FromPtr(float64Ptr(100.0)),
			CreatedAt:               now,
			UpdatedAt:               now,
		},
		{
			ID:                      FieldVibrationLevelID,
			SensorMeasurementTypeID: MeasurementTypeVibrationID,
			Name:                    "frequency",
			Label:                   "Dominant Frequency",
			Description:             repository.NullStringFromPtr(stringPtr("Dominant frequency component")),
			DataType:                "number",
			Required:                false,
			Unit:                    repository.NullStringFromPtr(stringPtr("Hz")),
			Min:                     repository.NullFloat64FromPtr(float64Ptr(0.0)),
			Max:                     repository.NullFloat64FromPtr(float64Ptr(10000.0)),
			CreatedAt:               now,
			UpdatedAt:               now,
		},

		// Flow Measurement Fields
		{
			ID:                      FieldFlowRateID,
			SensorMeasurementTypeID: MeasurementTypeFlowID,
			Name:                    "flow_rate",
			Label:                   "Flow Rate",
			Description:             repository.NullStringFromPtr(stringPtr("Current flow rate measurement")),
			DataType:                "number",
			Required:                true,
			Unit:                    repository.NullStringFromPtr(stringPtr("L/min")),
			Min:                     repository.NullFloat64FromPtr(float64Ptr(0.0)),
			Max:                     repository.NullFloat64FromPtr(float64Ptr(10000.0)),
			CreatedAt:               now,
			UpdatedAt:               now,
		},
		{
			ID:                      FieldTotalVolumeID,
			SensorMeasurementTypeID: MeasurementTypeFlowID,
			Name:                    "totalizer",
			Label:                   "Total Volume",
			Description:             repository.NullStringFromPtr(stringPtr("Cumulative volume measurement")),
			DataType:                "number",
			Required:                false,
			Unit:                    repository.NullStringFromPtr(stringPtr("L")),
			Min:                     repository.NullFloat64FromPtr(float64Ptr(0.0)),
			CreatedAt:               now,
			UpdatedAt:               now,
		},
		{
			ID:                      FieldFlowStatusID,
			SensorMeasurementTypeID: MeasurementTypeFlowID,
			Name:                    "temperature",
			Label:                   "Fluid Temperature",
			Description:             repository.NullStringFromPtr(stringPtr("Fluid temperature during flow measurement")),
			DataType:                "number",
			Required:                false,
			Unit:                    repository.NullStringFromPtr(stringPtr("°C")),
			Min:                     repository.NullFloat64FromPtr(float64Ptr(-20.0)),
			Max:                     repository.NullFloat64FromPtr(float64Ptr(120.0)),
			CreatedAt:               now,
			UpdatedAt:               now,
		},

		// Pressure Measurement Fields
		{
			ID:                      FieldPressureValueID,
			SensorMeasurementTypeID: MeasurementTypePressureID,
			Name:                    "pressure",
			Label:                   "Pressure Value",
			Description:             repository.NullStringFromPtr(stringPtr("Pressure measurement value")),
			DataType:                "number",
			Required:                true,
			Unit:                    repository.NullStringFromPtr(stringPtr("bar")),
			Min:                     repository.NullFloat64FromPtr(float64Ptr(0.0)),
			Max:                     repository.NullFloat64FromPtr(float64Ptr(400.0)),
			CreatedAt:               now,
			UpdatedAt:               now,
		},
		{
			ID:                      FieldPressureUnitID,
			SensorMeasurementTypeID: MeasurementTypePressureID,
			Name:                    "temperature",
			Label:                   "Process Temperature",
			Description:             repository.NullStringFromPtr(stringPtr("Process temperature during pressure measurement")),
			DataType:                "number",
			Required:                false,
			Unit:                    repository.NullStringFromPtr(stringPtr("°C")),
			Min:                     repository.NullFloat64FromPtr(float64Ptr(-40.0)),
			Max:                     repository.NullFloat64FromPtr(float64Ptr(150.0)),
			CreatedAt:               now,
			UpdatedAt:               now,
		},

		// Power Measurement Fields
		{
			ID:                      FieldPowerValueID,
			SensorMeasurementTypeID: MeasurementTypePowerID,
			Name:                    "power",
			Label:                   "Active Power",
			Description:             repository.NullStringFromPtr(stringPtr("Active power consumption")),
			DataType:                "number",
			Required:                true,
			Unit:                    repository.NullStringFromPtr(stringPtr("W")),
			Min:                     repository.NullFloat64FromPtr(float64Ptr(0.0)),
			Max:                     repository.NullFloat64FromPtr(float64Ptr(500000.0)),
			CreatedAt:               now,
			UpdatedAt:               now,
		},
		{
			ID:                      FieldVoltageValueID,
			SensorMeasurementTypeID: MeasurementTypePowerID,
			Name:                    "voltage",
			Label:                   "Voltage",
			Description:             repository.NullStringFromPtr(stringPtr("RMS voltage measurement")),
			DataType:                "number",
			Required:                true,
			Unit:                    repository.NullStringFromPtr(stringPtr("V")),
			Min:                     repository.NullFloat64FromPtr(float64Ptr(0.0)),
			Max:                     repository.NullFloat64FromPtr(float64Ptr(500.0)),
			CreatedAt:               now,
			UpdatedAt:               now,
		},
		{
			ID:                      FieldCurrentValueID,
			SensorMeasurementTypeID: MeasurementTypePowerID,
			Name:                    "current",
			Label:                   "Current",
			Description:             repository.NullStringFromPtr(stringPtr("RMS current measurement")),
			DataType:                "number",
			Required:                true,
			Unit:                    repository.NullStringFromPtr(stringPtr("A")),
			Min:                     repository.NullFloat64FromPtr(float64Ptr(0.0)),
			Max:                     repository.NullFloat64FromPtr(float64Ptr(1000.0)),
			CreatedAt:               now,
			UpdatedAt:               now,
		},
		{
			ID:                      FieldPowerFactorID,
			SensorMeasurementTypeID: MeasurementTypePowerID,
			Name:                    "power_factor",
			Label:                   "Power Factor",
			Description:             repository.NullStringFromPtr(stringPtr("Power factor measurement")),
			DataType:                "number",
			Required:                false,
			Unit:                    repository.NullStringFromPtr(stringPtr("")),
			Min:                     repository.NullFloat64FromPtr(float64Ptr(0.0)),
			Max:                     repository.NullFloat64FromPtr(float64Ptr(1.0)),
			CreatedAt:               now,
			UpdatedAt:               now,
		},
	}
}

// fieldExists checks if a measurement field exists
func (s *SensorMeasurementFieldSeeder) fieldExists(ctx context.Context, id uuid.UUID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM sensor_measurement_fields WHERE id = $1)`
	var exists bool
	err := s.db.QueryRowContext(ctx, query, id).Scan(&exists)
	return exists, err
}

// insertField inserts a measurement field into the database
func (s *SensorMeasurementFieldSeeder) insertField(ctx context.Context, field *repository.SensorMeasurementField) error {
	query := `
		INSERT INTO sensor_measurement_fields (
			id, sensor_measurement_type_id, name, label, description, 
			data_type, required, unit, min, max, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`

	_, err := s.db.ExecContext(ctx, query,
		field.ID,
		field.SensorMeasurementTypeID,
		field.Name,
		field.Label,
		field.Description,
		field.DataType,
		field.Required,
		field.Unit,
		field.Min,
		field.Max,
		field.CreatedAt,
		field.UpdatedAt,
	)

	return err
}

// Seed seeds all measurement fields
func (s *SensorMeasurementFieldSeeder) Seed(ctx context.Context) error {
	log.Println("Starting sensor measurement field seeding...")

	fields := s.getMeasurementFields()
	createdCount := 0
	skippedCount := 0

	for _, field := range fields {
		exists, err := s.fieldExists(ctx, field.ID)
		if err != nil {
			log.Printf("Error checking field existence for %s: %v", field.Name, err)
			continue
		}

		if exists {
			log.Printf("Field %s already exists, skipping", field.Name)
			skippedCount++
			continue
		}

		err = s.insertField(ctx, field)
		if err != nil {
			log.Printf("Error creating field %s: %v", field.Name, err)
			continue
		}

		log.Printf("Created measurement field: %s (%s)", field.Label, field.Name)
		createdCount++
	}

	log.Printf("Sensor measurement field seeding completed. Created: %d, Skipped: %d", createdCount, skippedCount)
	return nil
}

// Helper functions
func float64Ptr(f float64) *float64 {
	return &f
}

// GetFieldIDByName returns the UUID for a field by its name and measurement type
func GetFieldIDByName(measurementTypeID uuid.UUID, fieldName string) uuid.UUID {
	// Create a mapping for quick lookups
	fieldMap := map[string]map[string]uuid.UUID{
		MeasurementTypeWaterQualityID.String(): {
			"ph_value":           FieldPhValueID,
			"temperature":        FieldPhTemperatureID,
			"calibration_status": FieldPhStatusID,
		},
		MeasurementTypeTurbidityID.String(): {
			"turbidity_value":  FieldTurbidityValueID,
			"measurement_unit": FieldTurbidityUnitID,
			"water_quality":    FieldTurbidityQualityID,
		},
		MeasurementTypeAirQualityID.String(): {
			"pm25_concentration": FieldPm25ValueID,
			"co2_concentration":  FieldCo2ValueID,
			"air_quality_index":  FieldAirQualityID,
		},
		MeasurementTypeCO2ID.String(): {
			"humidity": FieldHumidityValueID,
		},
		MeasurementTypeTemperatureID.String(): {
			"temperature_value": FieldTemperatureValueID,
			"alarm_status":      FieldTempAlarmStatusID,
		},
		MeasurementTypeVibrationID.String(): {
			"rms_velocity":      FieldVibrationAmplitudeID,
			"peak_acceleration": FieldVibrationFrequencyID,
			"frequency":         FieldVibrationLevelID,
		},
		MeasurementTypeFlowID.String(): {
			"flow_rate":   FieldFlowRateID,
			"totalizer":   FieldTotalVolumeID,
			"temperature": FieldFlowStatusID,
		},
		MeasurementTypePressureID.String(): {
			"pressure":    FieldPressureValueID,
			"temperature": FieldPressureUnitID,
		},
		MeasurementTypePowerID.String(): {
			"power":        FieldPowerValueID,
			"voltage":      FieldVoltageValueID,
			"current":      FieldCurrentValueID,
			"power_factor": FieldPowerFactorID,
		},
	}

	if typeFields, exists := fieldMap[measurementTypeID.String()]; exists {
		if fieldID, found := typeFields[fieldName]; found {
			return fieldID
		}
	}

	return uuid.Nil
}
