#!/bin/bash

# Test script untuk Flexible IoT Sensor Reading API
# Pastikan server sudah running di localhost:8080

BASE_URL="http://localhost:8080/api/v1/iot-sensor-readings"
AUTH_HEADER="Authorization: Bearer YOUR_JWT_TOKEN_HERE"
CONTENT_TYPE="Content-Type: application/json"

echo "=== Testing Flexible IoT Sensor Reading API ==="
echo ""

# Test 1: Single Reading - Weather Station
echo "1. Testing Single Reading - Weather Station"
curl -X POST "$BASE_URL/flexible" \
  -H "$CONTENT_TYPE" \
  -H "$AUTH_HEADER" \
  -d @weather_station_example.json \
  -w "\nHTTP Status: %{http_code}\n" \
  --silent --show-error
echo ""
echo "---"

# Test 2: Single Reading - Industrial Sensor  
echo "2. Testing Single Reading - Industrial Sensor"
curl -X POST "$BASE_URL/flexible" \
  -H "$CONTENT_TYPE" \
  -H "$AUTH_HEADER" \
  -d @industrial_sensor_example.json \
  -w "\nHTTP Status: %{http_code}\n" \
  --silent --show-error
echo ""
echo "---"

# Test 3: Single Reading - Mixed Data Types
echo "3. Testing Single Reading - Mixed Data Types (Boolean, String, Number)"
curl -X POST "$BASE_URL/flexible" \
  -H "$CONTENT_TYPE" \
  -H "$AUTH_HEADER" \
  -d @mixed_data_types_example.json \
  -w "\nHTTP Status: %{http_code}\n" \
  --silent --show-error
echo ""
echo "---"

# Test 4: Batch Reading
echo "4. Testing Batch Reading - Multiple Sensors"
curl -X POST "$BASE_URL/flexible/batch" \
  -H "$CONTENT_TYPE" \
  -H "$AUTH_HEADER" \
  -d @batch_example.json \
  -w "\nHTTP Status: %{http_code}\n" \
  --silent --show-error
echo ""
echo "---"

# Test 5: Text Parsing - CSV Format
echo "5. Testing Text Parsing - CSV Format"
curl -X POST "$BASE_URL/parse-text" \
  -H "$CONTENT_TYPE" \
  -H "$AUTH_HEADER" \
  -d '{
    "text_data": "temperature:25.3°C,humidity:67%,pressure:1013.2hPa,wind_speed:3.2m/s",
    "sensor_type": "weather_station",
    "asset_sensor_id": "weather-parser-test",
    "sensor_type_id": "weather-station-type",
    "mac_address": "WS:PA:RS:ER:TE:ST"
  }' \
  -w "\nHTTP Status: %{http_code}\n" \
  --silent --show-error
echo ""
echo "---"

# Test 6: Text Parsing - Key-Value Format
echo "6. Testing Text Parsing - Key-Value Format"
curl -X POST "$BASE_URL/parse-text" \
  -H "$CONTENT_TYPE" \
  -H "$AUTH_HEADER" \
  -d '{
    "text_data": "pm25=15.2 pm10=28.1 temp=26.8 hum=72 status=OK",
    "sensor_type": "air_quality",
    "asset_sensor_id": "air-quality-parser",
    "sensor_type_id": "air-quality-type",
    "mac_address": "AQ:PA:RS:ER:99:AA"
  }' \
  -w "\nHTTP Status: %{http_code}\n" \
  --silent --show-error
echo ""
echo "---"

# Test 7: Text Parsing - Industrial Format
echo "7. Testing Text Parsing - Industrial Format"
curl -X POST "$BASE_URL/parse-text" \
  -H "$CONTENT_TYPE" \
  -H "$AUTH_HEADER" \
  -d '{
    "text_data": "VIB:2.3mm/s|TEMP:65.8°C|PWR:12.5kW|RPM:1750|OIL:4.2bar|STATUS:RUNNING",
    "sensor_type": "industrial_monitor",
    "asset_sensor_id": "machine-parser-01",
    "sensor_type_id": "machine-monitor-type",
    "mac_address": "IN:PA:RS:ER:DD:EE"
  }' \
  -w "\nHTTP Status: %{http_code}\n" \
  --silent --show-error
echo ""
echo "---"

echo "=== Testing Completed ==="
echo ""
echo "Note: Make sure to:"
echo "1. Replace YOUR_JWT_TOKEN_HERE with actual JWT token"
echo "2. Ensure the server is running on localhost:8080"
echo "3. Check the database for stored readings and measurement data"
echo "4. Verify the iot_sensor_measurement_data table has the flexible data"
