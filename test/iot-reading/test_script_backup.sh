#!/bin/bash

# IoT Sensor Reading Manual Testing Script
# Make sure to set your JWT token and server URL

# Configuration
SERVER_URL="http://localhost:3160"
JWT_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiZGIzNWIxYmMtN2I3OC00NzkwLWI3ZTEtNmY1NWNjMzc5YjBhIiwicm9sZV9pZCI6ImM2OTkzMGM1LTU1YzAtNDRkYi05Y2M1LTkwNDhkMmMxODFjNCIsInJvbGVfbmFtZSI6IlNVUEVSQURNSU4iLCJleHAiOjE3NDg2MjQwMDgsIm5iZiI6MTc0ODUzNzYwOCwiaWF0IjoxNzQ4NTM3NjA4fQ.Az4PUb7ipuHrFPhoYPC6WQ4vH0XSZwXIY97P2qRx5JY"
TEST_DIR="$(dirname "$0")"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}=== IoT Sensor Reading Manual Testing Script ===${NC}"
echo ""

# Check if JWT token is set
if [ "$JWT_TOKEN" = "YOUR_JWT_TOKEN_HERE" ]; then
    echo -e "${RED}ERROR: Please set your JWT token in the script${NC}"
    echo "Edit this script and replace YOUR_JWT_TOKEN_HERE with your actual JWT token"
    exit 1
fi

# Function to make API call
make_request() {
    local method=$1
    local endpoint=$2
    local data_file=$3
    local description=$4
    
    echo -e "${YELLOW}Testing: $description${NC}"
    echo "Method: $method"
    echo "Endpoint: $endpoint"
    if [ ! -z "$data_file" ]; then
        echo "Data file: $data_file"
    fi
    echo ""
    
    if [ "$method" = "GET" ]; then
        response=$(curl -s -w "\nHTTP_STATUS:%{http_code}" -X GET "$SERVER_URL$endpoint" \
            -H "Authorization: Bearer $JWT_TOKEN" \
            -H "Content-Type: application/json")
    else
        if [ ! -z "$data_file" ] && [ -f "$TEST_DIR/$data_file" ]; then
            response=$(curl -s -w "\nHTTP_STATUS:%{http_code}" -X "$method" "$SERVER_URL$endpoint" \
                -H "Authorization: Bearer $JWT_TOKEN" \
                -H "Content-Type: application/json" \
                -d @"$TEST_DIR/$data_file")
        else
            echo -e "${RED}ERROR: Data file $data_file not found${NC}"
            return 1
        fi
    fi
    
    # Extract HTTP status and body
    http_status=$(echo "$response" | grep "HTTP_STATUS:" | cut -d: -f2)
    body=$(echo "$response" | grep -v "HTTP_STATUS:")
    
    # Color code the status
    if [[ $http_status =~ ^2[0-9][0-9]$ ]]; then
        echo -e "Status: ${GREEN}$http_status${NC}"
    elif [[ $http_status =~ ^4[0-9][0-9]$ ]]; then
        echo -e "Status: ${YELLOW}$http_status${NC}"
    else
        echo -e "Status: ${RED}$http_status${NC}"
    fi
    
    echo "Response:"
    echo "$body" | jq '.' 2>/dev/null || echo "$body"
    echo ""
    echo "----------------------------------------"
    echo ""
}

# Test 1: Get JSON Template
make_request "GET" "/api/v1/superadmin/iot-sensor-readings/template?sensor_type=temperature" "" "Get Temperature JSON Template"

# Test 2: Get Batch JSON Template
make_request "GET" "/api/v1/superadmin/iot-sensor-readings/template/batch?sensor_type=humidity&count=3" "" "Get Batch Humidity JSON Template"

# Test 3: Create Single Temperature Reading
make_request "POST" "/api/v1/superadmin/iot-sensor-readings/from-json" "single_temperature_scraped.json" "Create Single Temperature Reading"

# Test 4: Create Single Humidity Reading
make_request "POST" "/api/v1/superadmin/iot-sensor-readings/from-json" "single_humidity_scraped.json" "Create Single Humidity Reading"

# Test 5: Create Batch Temperature Readings
make_request "POST" "/api/v1/superadmin/iot-sensor-readings/batch-from-json" "batch_scraped.json" "Create Batch Temperature Readings"

# Test 6: Create Mixed Sensor Batch
make_request "POST" "/api/v1/superadmin/iot-sensor-readings/batch-from-json" "batch_mixed_sensors.json" "Create Mixed Sensor Batch"

# Test 7: Create Dummy Temperature Reading
make_request "POST" "/api/v1/superadmin/iot-sensor-readings/dummy" "dummy_temperature_request.json" "Create Dummy Temperature Reading"

# Test 8: Create Multiple Dummy Readings
make_request "POST" "/api/v1/superadmin/iot-sensor-readings/dummy/multiple" "dummy_scraped.json" "Create Multiple Dummy Readings"

# Test 9: Create Simple Reading
make_request "POST" "/api/v1/superadmin/iot-sensor-readings/simple" "simple_reading_scraped.json" "Create Simple Reading"

# Test 10: Test Invalid UUID (Should fail)
make_request "POST" "/api/v1/superadmin/iot-sensor-readings/from-json" "invalid_uuid.json" "Test Invalid UUID (Expected to fail)"

# Test 11: Test Missing Fields (Should fail)
make_request "POST" "/api/v1/superadmin/iot-sensor-readings/from-json" "missing_required_fields.json" "Test Missing Required Fields (Expected to fail)"

# Test 12: Test Invalid Value Type (Should fail)
make_request "POST" "/api/v1/superadmin/iot-sensor-readings/from-json" "invalid_value_type.json" "Test Invalid Value Type (Expected to fail)"

# Test 13: Test High Precision Reading
make_request "POST" "/api/v1/superadmin/iot-sensor-readings/from-json" "high_precision.json" "Test High Precision Reading"

# Test 14: Test Extreme Values
make_request "POST" "/api/v1/superadmin/iot-sensor-readings/batch-from-json" "extreme_values.json" "Test Extreme Values"

echo -e "${GREEN}=== Testing Complete ===${NC}"
echo ""
echo -e "${BLUE}Summary:${NC}"
echo "- Tested all 7 manual input endpoints"
echo "- Tested various sensor types (temperature, humidity, pressure, light, sound, motion)"
echo "- Tested error handling (invalid UUID, missing fields, wrong types)"
echo "- Tested edge cases (extreme values, high precision)"
echo ""
echo -e "${YELLOW}Next Steps:${NC}"
echo "1. Review the responses above"
echo "2. Check server logs for any additional error details"
echo "3. Verify data was correctly stored in the database"
echo "4. Test with different JWT tokens to verify authorization"
