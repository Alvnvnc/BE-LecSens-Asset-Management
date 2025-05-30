package dto

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// CreateIoTSensorReadingRequest represents the request to create a new IoT sensor reading
type CreateIoTSensorReadingRequest struct {
	AssetSensorID uuid.UUID  `json:"asset_sensor_id" binding:"required" validate:"required"`
	SensorTypeID  uuid.UUID  `json:"sensor_type_id" binding:"required" validate:"required"`
	MacAddress    string     `json:"mac_address" binding:"required" validate:"required"`
	ReadingTime   *time.Time `json:"reading_time,omitempty"` // Optional, defaults to current time
}

// CreateBatchIoTSensorReadingRequest represents the request to create multiple IoT sensor readings
type CreateBatchIoTSensorReadingRequest struct {
	Readings []CreateIoTSensorReadingRequest `json:"readings" binding:"required" validate:"required,min=1,max=1000"`
}

// UpdateIoTSensorReadingRequest represents the request to update an existing IoT sensor reading
type UpdateIoTSensorReadingRequest struct {
	MacAddress  *string    `json:"mac_address,omitempty"`
	ReadingTime *time.Time `json:"reading_time,omitempty"`
}

// IoTSensorReadingResponse represents the response structure for IoT sensor reading operations
type IoTSensorReadingResponse struct {
	ID              uuid.UUID                   `json:"id"`
	TenantID        uuid.UUID                   `json:"tenant_id"`
	AssetSensorID   uuid.UUID                   `json:"asset_sensor_id"`
	SensorTypeID    uuid.UUID                   `json:"sensor_type_id"`
	MacAddress      string                      `json:"mac_address"`
	Location        string                      `json:"location"`
	ReadingTime     time.Time                   `json:"reading_time"`
	CreatedAt       time.Time                   `json:"created_at"`
	UpdatedAt       *time.Time                  `json:"updated_at,omitempty"`
	MeasurementData map[string]MeasurementValue `json:"measurement_data,omitempty"`
	Message         string                      `json:"message,omitempty"`
	Warnings        []string                    `json:"warnings,omitempty"`
}

// IoTSensorReadingWithDetailsResponse represents the response with detailed related information
type IoTSensorReadingWithDetailsResponse struct {
	*IoTSensorReadingResponse
	AssetSensor struct {
		ID            uuid.UUID       `json:"id"`
		AssetID       uuid.UUID       `json:"asset_id"`
		Name          string          `json:"name"`
		Status        string          `json:"status"`
		Configuration json.RawMessage `json:"configuration"`
	} `json:"asset_sensor"`
	SensorType struct {
		ID           uuid.UUID `json:"id"`
		Name         string    `json:"name"`
		Description  string    `json:"description"`
		Manufacturer string    `json:"manufacturer"`
		Model        string    `json:"model"`
		Version      string    `json:"version"`
		IsActive     bool      `json:"is_active"`
	} `json:"sensor_type"`
	MeasurementTypes []struct {
		ID               uuid.UUID       `json:"id"`
		Name             string          `json:"name"`
		Description      string          `json:"description"`
		PropertiesSchema json.RawMessage `json:"properties_schema"`
		UIConfiguration  json.RawMessage `json:"ui_configuration"`
		Version          string          `json:"version"`
		IsActive         bool            `json:"is_active"`
		Fields           []struct {
			ID          uuid.UUID `json:"id"`
			Name        string    `json:"name"`
			Label       string    `json:"label"`
			Description *string   `json:"description"`
			DataType    string    `json:"data_type"`
			Required    bool      `json:"required"`
			Unit        *string   `json:"unit"`
			Min         *float64  `json:"min"`
			Max         *float64  `json:"max"`
		} `json:"fields"`
	} `json:"measurement_types"`
}

// IoTSensorReadingListResponse represents the response for listing IoT sensor readings with pagination
type IoTSensorReadingListResponse struct {
	Readings   []IoTSensorReadingWithDetailsResponse `json:"readings"`
	Page       int                                   `json:"page"`
	Limit      int                                   `json:"limit"`
	Total      int64                                 `json:"total"`
	TotalPages int                                   `json:"total_pages"`
}

// IoTSensorReadingListRequest represents parameters for listing IoT sensor readings
type IoTSensorReadingListRequest struct {
	AssetSensorID *uuid.UUID `json:"asset_sensor_id,omitempty"`
	SensorTypeID  *uuid.UUID `json:"sensor_type_id,omitempty"`
	MacAddress    *string    `json:"mac_address,omitempty"`
	FromTime      *time.Time `json:"from_time,omitempty"`
	ToTime        *time.Time `json:"to_time,omitempty"`
	Page          int        `json:"page"`
	PageSize      int        `json:"page_size"`
}

// GetReadingsInTimeRangeRequest represents request for time-range queries
type GetReadingsInTimeRangeRequest struct {
	AssetSensorID *uuid.UUID `json:"asset_sensor_id,omitempty"`
	SensorTypeID  *uuid.UUID `json:"sensor_type_id,omitempty"`
	FromTime      time.Time  `json:"from_time" binding:"required" validate:"required"`
	ToTime        time.Time  `json:"to_time" binding:"required" validate:"required"`
	Limit         int        `json:"limit,omitempty"` // Optional limit, defaults to 1000
}

// GetAggregatedDataRequest represents request for aggregated analytics data
type GetAggregatedDataRequest struct {
	AssetSensorID *uuid.UUID `json:"asset_sensor_id,omitempty"`
	SensorTypeID  *uuid.UUID `json:"sensor_type_id,omitempty"`
	FromTime      time.Time  `json:"from_time" binding:"required" validate:"required"`
	ToTime        time.Time  `json:"to_time" binding:"required" validate:"required"`
	Interval      string     `json:"interval,omitempty"`     // hour, day, week, month - defaults to "hour"
	AggregateBy   []string   `json:"aggregate_by,omitempty"` // Fields to aggregate from measurement_data
}

// AggregatedDataPoint represents a single aggregated data point
type AggregatedDataPoint struct {
	Time     time.Time              `json:"time"`
	Count    int64                  `json:"count"`
	Averages map[string]float64     `json:"averages,omitempty"`
	Sums     map[string]float64     `json:"sums,omitempty"`
	Mins     map[string]float64     `json:"mins,omitempty"`
	Maxs     map[string]float64     `json:"maxs,omitempty"`
	Data     map[string]interface{} `json:"data,omitempty"` // Additional aggregated data
}

// GetAggregatedDataResponse represents response for aggregated analytics data
type GetAggregatedDataResponse struct {
	DataPoints  []AggregatedDataPoint `json:"data_points"`
	TotalCount  int64                 `json:"total_count"`
	FromTime    time.Time             `json:"from_time"`
	ToTime      time.Time             `json:"to_time"`
	Interval    string                `json:"interval"`
	AggregateBy []string              `json:"aggregate_by"`
	RequestedAt time.Time             `json:"requested_at"`
}

// ValidateAndCreateRequest represents request for validating measurement data against schemas
type ValidateAndCreateRequest struct {
	CreateIoTSensorReadingRequest
	ValidateSchema bool `json:"validate_schema,omitempty"` // Whether to validate against measurement type schemas
}

// ValidationError represents schema validation error details
type ValidationError struct {
	Field   string      `json:"field"`
	Message string      `json:"message"`
	Value   interface{} `json:"value,omitempty"`
}

// ValidateAndCreateResponse represents response with validation details
type ValidateAndCreateResponse struct {
	*IoTSensorReadingResponse
	ValidationErrors []ValidationError `json:"validation_errors,omitempty"`
	IsValid          bool              `json:"is_valid"`
}

// FlexibleIoTSensorReadingRequest represents a flexible request structure for creating IoT sensor readings
// with dynamic measurement data that automatically maps to database tables
type FlexibleIoTSensorReadingRequest struct {
	AssetSensorID   uuid.UUID                   `json:"asset_sensor_id" binding:"required" validate:"required"`
	SensorTypeID    uuid.UUID                   `json:"sensor_type_id" binding:"required" validate:"required"`
	MacAddress      string                      `json:"mac_address" binding:"required" validate:"required"`
	ReadingTime     *time.Time                  `json:"reading_time,omitempty"` // Optional, defaults to current time
	MeasurementData map[string]MeasurementValue `json:"-"`                      // Will be populated from other fields
	RawJSON         map[string]interface{}      `json:"-"`                      // Store the raw JSON for processing
}

// MeasurementValue represents a single measurement value with metadata
type MeasurementValue struct {
	Label string      `json:"label"` // e.g., "Temperature", "Raw Value"
	Unit  string      `json:"unit"`  // e.g., "°C", "μg/m³"
	Value interface{} `json:"value"` // The actual measurement value (can be float64, int, string, etc.)
}

// FlexibleBatchIoTSensorReadingRequest represents a flexible batch request
type FlexibleBatchIoTSensorReadingRequest struct {
	Readings []FlexibleIoTSensorReadingRequest `json:"readings" binding:"required" validate:"required,min=1,max=1000"`
}

// UnmarshalJSON custom unmarshaler for FlexibleIoTSensorReadingRequest
func (f *FlexibleIoTSensorReadingRequest) UnmarshalJSON(data []byte) error {
	// First, unmarshal into a raw map to capture all fields
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// Store the raw JSON for later processing
	f.RawJSON = raw

	// Extract standard fields
	if assetSensorID, ok := raw["asset_sensor_id"].(string); ok {
		if id, err := uuid.Parse(assetSensorID); err == nil {
			f.AssetSensorID = id
		}
	}

	if sensorTypeID, ok := raw["sensor_type_id"].(string); ok {
		if id, err := uuid.Parse(sensorTypeID); err == nil {
			f.SensorTypeID = id
		}
	}

	if macAddress, ok := raw["mac_address"].(string); ok {
		f.MacAddress = macAddress
	}

	if readingTimeStr, ok := raw["reading_time"].(string); ok {
		if t, err := time.Parse(time.RFC3339, readingTimeStr); err == nil {
			f.ReadingTime = &t
		}
	}

	// Extract measurement data
	f.MeasurementData = make(map[string]MeasurementValue)

	// First, check if there's an explicit "measurement_data" field
	if measurementDataRaw, ok := raw["measurement_data"].(map[string]interface{}); ok {
		for key, value := range measurementDataRaw {
			if valueMap, ok := value.(map[string]interface{}); ok {
				measurement := MeasurementValue{}

				if label, ok := valueMap["label"].(string); ok {
					measurement.Label = label
				}
				if unit, ok := valueMap["unit"].(string); ok {
					measurement.Unit = unit
				}
				if val, ok := valueMap["value"]; ok {
					measurement.Value = val
				}

				// Only add if we have at least a value
				if measurement.Value != nil {
					f.MeasurementData[key] = measurement
				}
			}
		}
	} else {
		// Fallback: check for top-level measurement fields (backwards compatibility)
		standardFields := map[string]bool{
			"asset_sensor_id":  true,
			"sensor_type_id":   true,
			"mac_address":      true,
			"reading_time":     true,
			"measurement_data": true,
		}

		for key, value := range raw {
			if standardFields[key] {
				continue
			}

			// Check if this field looks like measurement data (has unit, label, value)
			if valueMap, ok := value.(map[string]interface{}); ok {
				measurement := MeasurementValue{}

				if label, ok := valueMap["label"].(string); ok {
					measurement.Label = label
				}
				if unit, ok := valueMap["unit"].(string); ok {
					measurement.Unit = unit
				}
				if val, ok := valueMap["value"]; ok {
					measurement.Value = val
				}

				// Only add if we have at least a value
				if measurement.Value != nil {
					f.MeasurementData[key] = measurement
				}
			}
		}
	}

	return nil
}

// FlexibleIoTSensorReadingResponse represents the response for flexible IoT sensor reading operations
type FlexibleIoTSensorReadingResponse struct {
	*IoTSensorReadingResponse
	MeasurementData map[string]MeasurementValue `json:"measurement_data,omitempty"`
}

// TextToJSONRequest represents a request to convert text/raw data to JSON format
type TextToJSONRequest struct {
	TextData      string `json:"text_data" binding:"required"`
	SensorType    string `json:"sensor_type,omitempty"`     // Optional hint for parsing
	AssetSensorID string `json:"asset_sensor_id,omitempty"` // Optional pre-fill
	SensorTypeID  string `json:"sensor_type_id,omitempty"`  // Optional pre-fill
	MacAddress    string `json:"mac_address,omitempty"`     // Optional pre-fill
}

// TextToJSONResponse represents the response from text to JSON conversion
type TextToJSONResponse struct {
	ParsedJSON     map[string]interface{} `json:"parsed_json"`
	Success        bool                   `json:"success"`
	Message        string                 `json:"message"`
	Warnings       []string               `json:"warnings,omitempty"`
	SuggestedField map[string]string      `json:"suggested_fields,omitempty"` // Field mapping suggestions
}
