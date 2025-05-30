# Contoh Penggunaan Flexible IoT Sensor Reading API

## 1. Single Reading - Air Quality Sensor

### Endpoint: `POST /api/v1/iot-sensor-readings/flexible`

```bash
curl -X POST http://localhost:8080/api/v1/iot-sensor-readings/flexible \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "asset_sensor_id": "d906316f-c0bc-44c2-846a-ce5ee6363498",
    "sensor_type_id": "5c5b5461-e8d6-4c88-93fd-4b04019669bf",
    "mac_address": "AA:BB:CC:DD:EE:A2",
    "raw_value": {
      "unit": "μg/m³",
      "label": "Raw Value",
      "value": 45.6
    },
    "temperature": {
      "unit": "°C",
      "label": "Temperature",
      "value": 25.3
    },
    "humidity": {
      "unit": "%",
      "label": "Humidity",
      "value": 67.8
    },
    "pm25": {
      "unit": "μg/m³",
      "label": "PM2.5",
      "value": 12.4
    }
  }'
```

## 2. Batch Reading - Multiple Sensors

### Endpoint: `POST /api/v1/iot-sensor-readings/flexible/batch`

```bash
curl -X POST http://localhost:8080/api/v1/iot-sensor-readings/flexible/batch \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d @batch_example.json
```

## 3. Text Parsing - CSV Format

### Endpoint: `POST /api/v1/iot-sensor-readings/parse-text`

```bash
curl -X POST http://localhost:8080/api/v1/iot-sensor-readings/parse-text \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "text_data": "temperature:25.3°C,humidity:67%,pressure:1013.2hPa",
    "sensor_type": "weather_station",
    "asset_sensor_id": "weather-01",
    "sensor_type_id": "weather-station-type",
    "mac_address": "WS:AA:BB:CC:DD:EE"
  }'
```

## 4. Text Parsing - Key-Value Format

```bash
curl -X POST http://localhost:8080/api/v1/iot-sensor-readings/parse-text \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "text_data": "pm25=15.2 pm10=28.1 temp=26.8 hum=72",
    "sensor_type": "air_quality",
    "asset_sensor_id": "air-quality-01",
    "sensor_type_id": "air-quality-type",
    "mac_address": "AQ:66:77:88:99:AA"
  }'
```

## 5. Text Parsing - Industrial Format

```bash
curl -X POST http://localhost:8080/api/v1/iot-sensor-readings/parse-text \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "text_data": "VIB:2.3mm/s|TEMP:65.8°C|PWR:12.5kW|RPM:1750|OIL:4.2bar|STATUS:RUNNING",
    "sensor_type": "industrial_monitor",
    "asset_sensor_id": "machine-01",
    "sensor_type_id": "machine-monitor-type",
    "mac_address": "IN:AA:BB:CC:DD:EE"
  }'
```

## Array Examples untuk Testing

### 1. Simple Array (array_example_simple.json)
Array dengan 3 sensor readings sederhana dengan parameter temperature, humidity, dan PM2.5.

### 2. Batch Array (batch_example.json) 
Array untuk batch processing dengan multiple sensor locations.

### 3. Weather Station (weather_station_example.json)
Single reading dengan banyak parameter cuaca.

### 4. Industrial Sensor (industrial_sensor_example.json)
Single reading untuk monitoring mesin industri.

## Format Data yang Didukung

### 1. Measurement Value Structure
```json
{
  "field_name": {
    "unit": "unit_of_measurement",
    "label": "Human readable label",
    "value": actual_value
  }
}
```

### 2. Supported Value Types
- `number` (int, float): 25.3, 1750, 12
- `string`: "RUNNING", "OK", "ERROR"
- `boolean`: true, false

### 3. Text Parsing Formats
- **CSV**: `field1:value1unit,field2:value2unit`
- **Key-Value**: `field1=value1 field2=value2`
- **Pipe Separated**: `FIELD1:value1unit|FIELD2:value2unit`
- **JSON String**: `{"field1":value1,"field2":value2}`

## Response Format

### Success Response
```json
{
  "id": "uuid",
  "tenant_id": "uuid",
  "asset_sensor_id": "uuid",
  "sensor_type_id": "uuid", 
  "mac_address": "string",
  "location": "string",
  "reading_time": "2025-05-30T08:00:00Z",
  "created_at": "2025-05-30T08:00:00Z",
  "measurement_data": {
    "temperature": {
      "unit": "°C",
      "label": "Temperature",
      "value": 25.3
    },
    "humidity": {
      "unit": "%",
      "label": "Humidity", 
      "value": 67.8
    }
  }
}
```

### Batch Response
```json
{
  "created_readings": [
    // Array of IoTSensorReadingResponse
  ],
  "total_created": 3,
  "errors": []
}
```

### Text Parsing Response
```json
{
  "parsed_json": {
    "temperature": {
      "unit": "°C",
      "label": "Temperature",
      "value": 25.3
    }
  },
  "success": true,
  "message": "Successfully parsed text data",
  "warnings": [],
  "suggested_fields": {
    "temp": "temperature",
    "hum": "humidity"
  }
}
```

## Testing Steps

1. **Test Basic Flexible Reading**
   - Gunakan weather_station_example.json
   - POST ke endpoint flexible
   - Verifikasi response dan database

2. **Test Batch Processing**
   - Gunakan batch_example.json  
   - POST ke endpoint batch
   - Verifikasi semua readings tersimpan

3. **Test Text Parsing**
   - Test berbagai format text
   - Verifikasi parsing berhasil
   - Gunakan hasil untuk create reading

4. **Test Different Data Types**
   - Test numeric, string, boolean values
   - Verifikasi storage di measurement_data table
