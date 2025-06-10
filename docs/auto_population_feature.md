# IoT Sensor Reading Auto-Population Feature

## Overview

This feature implements auto-population logic for IoT sensor readings to automatically determine `asset_sensor_id` based on `sensor_type_id` when creating sensor readings. This reduces the complexity for API clients who may only know the sensor type but not the specific asset sensor ID.

## Features

### 1. Auto-Population Logic
- **Single Asset Sensor**: If exactly one asset sensor exists for the given sensor type, it will be automatically selected
- **Multiple Asset Sensors**: If multiple asset sensors exist, an error is returned with available options
- **No Asset Sensors**: If no asset sensors exist for the sensor type, an appropriate error is returned
- **Manual Override**: Users can still manually specify `asset_sensor_id` to override auto-population

### 2. Multi-Tenancy Support
- **Regular Users**: Auto-population respects tenant boundaries
- **SuperAdmin**: Can access all tenants or work within a specific tenant context

### 3. New API Endpoints

#### GET `/api/v1/iot-sensor-readings/auto-populate/options`
- **Purpose**: Get available asset sensor options for a sensor type
- **Access**: Public (requires tenant validation)
- **Parameters**: 
  - `sensor_type_id` (query parameter, required)
- **Response**: List of available asset sensors with location information

#### POST `/api/v1/superadmin/iot-sensor-readings/auto-populate`
- **Purpose**: Create IoT sensor reading with auto-population
- **Access**: SuperAdmin only
- **Body**: `CreateIoTSensorReadingWithAutoPopulationRequest`
- **Response**: Created sensor reading with auto-populated data

## Implementation Details

### 1. Data Transfer Objects (DTOs)

#### `CreateIoTSensorReadingWithAutoPopulationRequest`
```go
type CreateIoTSensorReadingWithAutoPopulationRequest struct {
    AssetSensorID *uuid.UUID `json:"asset_sensor_id,omitempty"` // Optional - auto-populated if not provided
    SensorTypeID  uuid.UUID  `json:"sensor_type_id" binding:"required"`
    MacAddress    string     `json:"mac_address" binding:"required"`
    ReadingTime   *time.Time `json:"reading_time,omitempty"`
}
```

#### `AutoPopulationOptionsResponse`
```go
type AutoPopulationOptionsResponse struct {
    SensorTypeID     uuid.UUID                  `json:"sensor_type_id"`
    AvailableOptions []AssetSensorLocationInfo  `json:"available_options"`
    Message          string                     `json:"message"`
}
```

#### `AssetSensorLocationInfo`
```go
type AssetSensorLocationInfo struct {
    AssetSensorID   uuid.UUID `json:"asset_sensor_id"`
    AssetSensorName string    `json:"asset_sensor_name"`
    AssetID         uuid.UUID `json:"asset_id"`
    AssetName       string    `json:"asset_name"`
    LocationID      uuid.UUID `json:"location_id"`
    LocationName    string    `json:"location_name"`
}
```

### 2. Repository Layer

#### `GetAssetSensorsBySensorType` Method
- Queries database to find active asset sensors for a given sensor type
- Joins with assets and locations tables to provide complete information
- Supports multi-tenancy filtering
- Returns asset sensor information with location details

```sql
SELECT 
    asn.id as asset_sensor_id,
    asn.name as asset_sensor_name,
    a.id as asset_id,
    a.name as asset_name,
    l.id as location_id,
    l.name as location_name
FROM asset_sensors asn
JOIN assets a ON asn.asset_id = a.id
JOIN locations l ON a.location_id = l.id
WHERE asn.sensor_type_id = $1
AND asn.status = 'active'
ORDER BY l.name, a.name, asn.name
```

### 3. Service Layer

#### `CreateIoTSensorReadingWithAutoPopulation` Method
- Main auto-population logic
- Handles single/multiple/no asset sensor scenarios
- Validates manual asset_sensor_id if provided
- Creates sensor reading with appropriate asset_sensor_id

#### `GetAutoPopulationOptions` Method
- Retrieves available asset sensors for a sensor type
- Returns formatted response with options and descriptive messages
- Supports different user roles and tenant contexts

### 4. Controller Layer

#### `CreateReadingWithAutoPopulation` Method
- Handles HTTP requests for auto-population creation
- Binds JSON request to DTO
- Calls service layer and handles responses
- Returns appropriate HTTP status codes and error messages

#### `GetAutoPopulationOptions` Method
- Handles HTTP requests for getting auto-population options
- Validates query parameters
- Returns available options or appropriate errors

### 5. Route Configuration

Routes added to `/asset_management/presentation-layer/routes/iot_sensor_reading.go`:

```go
// Public route for getting options
iotSensorReadingGroup.GET("/auto-populate/options", iotSensorReadingController.GetAutoPopulationOptions)

// SuperAdmin route for auto-population creation
superAdminGroup.POST("/auto-populate", iotSensorReadingController.CreateReadingWithAutoPopulation)
```

## Usage Examples

### 1. Get Auto-Population Options

```bash
curl -X GET "http://localhost:8080/api/v1/iot-sensor-readings/auto-populate/options?sensor_type_id=550e8400-e29b-41d4-a716-446655440000" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json"
```

**Response (Single Option):**
```json
{
  "message": "Auto-population options retrieved successfully",
  "data": {
    "sensor_type_id": "550e8400-e29b-41d4-a716-446655440000",
    "available_options": [
      {
        "asset_sensor_id": "660e8400-e29b-41d4-a716-446655440001",
        "asset_sensor_name": "Temperature Sensor 1",
        "asset_id": "770e8400-e29b-41d4-a716-446655440002",
        "asset_name": "Boiler Unit A",
        "location_id": "880e8400-e29b-41d4-a716-446655440003",
        "location_name": "Factory Floor 1"
      }
    ],
    "message": "1 asset sensor found for this sensor type"
  }
}
```

### 2. Create Reading with Auto-Population

```bash
curl -X POST "http://localhost:8080/api/v1/superadmin/iot-sensor-readings/auto-populate" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "sensor_type_id": "550e8400-e29b-41d4-a716-446655440000",
    "mac_address": "AA:BB:CC:DD:EE:FF"
  }'
```

**Response (Success):**
```json
{
  "message": "IoT sensor reading created successfully with auto-population",
  "data": {
    "id": "990e8400-e29b-41d4-a716-446655440004",
    "tenant_id": "aa0e8400-e29b-41d4-a716-446655440005",
    "asset_sensor_id": "660e8400-e29b-41d4-a716-446655440001",
    "sensor_type_id": "550e8400-e29b-41d4-a716-446655440000",
    "mac_address": "AA:BB:CC:DD:EE:FF",
    "reading_time": "2025-06-10T10:30:00Z",
    "created_at": "2025-06-10T10:30:00Z"
  }
}
```

### 3. Manual Asset Sensor ID

```bash
curl -X POST "http://localhost:8080/api/v1/superadmin/iot-sensor-readings/auto-populate" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "asset_sensor_id": "660e8400-e29b-41d4-a716-446655440001",
    "sensor_type_id": "550e8400-e29b-41d4-a716-446655440000",
    "mac_address": "AA:BB:CC:DD:EE:FF"
  }'
```

## Error Scenarios

### 1. Multiple Asset Sensors Available
```json
{
  "error": "Bad Request",
  "message": "Multiple asset sensors found for this sensor type. Please specify asset_sensor_id. Available options: Temperature Sensor 1, Temperature Sensor 2"
}
```

### 2. No Asset Sensors Found
```json
{
  "error": "Not Found",
  "message": "No asset sensors found for this sensor type"
}
```

### 3. Invalid Sensor Type ID
```json
{
  "error": "Bad Request",
  "message": "Invalid sensor_type_id format"
}
```

### 4. Asset Sensor Mismatch
```json
{
  "error": "Bad Request",
  "message": "Provided asset_sensor_id does not belong to the specified sensor_type_id"
}
```

## Benefits

1. **Simplified API Usage**: Clients can create readings with just sensor type ID
2. **Reduced Complexity**: No need to query for asset sensor IDs separately
3. **Error Prevention**: Clear error messages when multiple options exist
4. **Flexibility**: Manual override still available when needed
5. **Multi-Tenancy**: Respects tenant boundaries and permissions
6. **Performance**: Efficient database queries with proper joins

## Testing

Run the test script to see example usage:
```bash
./test/test_auto_population.sh
```

This script provides comprehensive examples of all endpoints and scenarios.

## Files Modified

1. **DTOs**: `/helpers/dto/iot_sensor_reading_dto.go`
2. **Repository**: `/data-layer/repository/iot_sensor_reading_repository.go`
3. **Service**: `/domain-layer/service/iot_sensor_reading_service.go`
4. **Controller**: `/presentation-layer/controller/iot_sensor_reading_controller.go`
5. **Routes**: `/presentation-layer/routes/iot_sensor_reading.go`
6. **Test**: `/test/test_auto_population.sh`
7. **Documentation**: `/docs/auto_population_feature.md`

## Future Enhancements

1. **Caching**: Add caching for asset sensor lookups to improve performance
2. **Bulk Operations**: Extend auto-population to batch reading creation
3. **Smart Selection**: Add logic to prefer certain asset sensors based on criteria
4. **Audit Trail**: Log auto-population decisions for auditing purposes
5. **Configuration**: Allow configuration of auto-population behavior per tenant
