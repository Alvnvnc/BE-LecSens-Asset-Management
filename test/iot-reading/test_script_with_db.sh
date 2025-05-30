#!/bin/bash

# IoT Sensor Reading Testing Script with Database Scraping
# This script fetches asset sensor data from the database and uses it for testing

# Configuration
SERVER_URL="http://localhost:3160"
JWT_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiZGIzNWIxYmMtN2I3OC00NzkwLWI3ZTEtNmY1NWNjMzc5YjBhIiwicm9sZV9pZCI6ImM2OTkzMGM1LTU1YzAtNDRkYi05Y2M1LTkwNDhkMmMxODFjNCIsInJvbGVfbmFtZSI6IlNVUEVSQURNSU4iLCJleHAiOjE3NDg2MjQwMDgsIm5iZiI6MTc0ODUzNzYwOCwiaWF0IjoxNzQ4NTM3NjA4fQ.Az4PUb7ipuHrFPhoYPC6WQ4vH0XSZwXIY97P2qRx5JY"
TEST_DIR="$(dirname "$0")"

# Database configuration
DB_HOST="localhost"
DB_PORT="5432"
DB_NAME="lecsens_db"
DB_USER="postgres"
DB_PASSWORD="postgres"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
NC='\033[0m' # No Color

echo -e "${BLUE}=== IoT Sensor Reading Testing Script with Database Scraping ===${NC}"
echo ""

# Check if JWT token is set
if [ "$JWT_TOKEN" = "YOUR_JWT_TOKEN_HERE" ]; then
    echo -e "${RED}ERROR: Please set your JWT token in the script${NC}"
    echo "Edit this script and replace YOUR_JWT_TOKEN_HERE with your actual JWT token"
    exit 1
fi

# Check if psql is available
if ! command -v psql &> /dev/null; then
    echo -e "${RED}ERROR: psql is not installed or not in PATH${NC}"
    echo "Please install PostgreSQL client tools"
    exit 1
fi

# Check if jq is available
if ! command -v jq &> /dev/null; then
    echo -e "${RED}ERROR: jq is not installed or not in PATH${NC}"
    echo "Please install jq for JSON processing"
    exit 1
fi

# Function to execute SQL query
execute_sql() {
    local query=$1
    PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -t -c "$query" 2>/dev/null
}

# Function to fetch asset sensors from database
fetch_asset_sensors() {
    echo -e "${PURPLE}Fetching asset sensors from database...${NC}"
    
    local query="
    SELECT 
        json_agg(
            json_build_object(
                'asset_sensor_id', asn.id,
                'asset_sensor_name', asn.name,
                'asset_id', asn.asset_id,
                'asset_name', a.name,
                'location_id', a.location_id,
                'location_name', l.name,
                'sensor_type_id', st.id,
                'sensor_type_name', st.name,
                'status', asn.status
            )
        )
    FROM asset_sensors asn
    JOIN assets a ON asn.asset_id = a.id
    JOIN locations l ON a.location_id = l.id
    JOIN sensor_types st ON st.is_active = true
    WHERE asn.status = 'active'
    LIMIT 10;
    "
    
    local result=$(execute_sql "$query")
    
    if [ $? -eq 0 ] && [ ! -z "$result" ] && [ "$result" != "null" ]; then
        echo "$result"
    else
        echo -e "${RED}Failed to fetch asset sensors from database${NC}"
        echo "Please check your database connection settings"
        return 1
    fi
}

# Function to fetch sensor types
fetch_sensor_types() {
    echo -e "${PURPLE}Fetching sensor types from database...${NC}"
    
    local query="
    SELECT json_agg(
        json_build_object(
            'id', id,
            'name', name,
            'manufacturer', manufacturer,
            'model', model
        )
    )
    FROM sensor_types 
    WHERE is_active = true;
    "
    
    local result=$(execute_sql "$query")
    
    if [ $? -eq 0 ] && [ ! -z "$result" ] && [ "$result" != "null" ]; then
        echo "$result"
    else
        echo -e "${RED}Failed to fetch sensor types from database${NC}"
        return 1
    fi
}

# Function to make API call
make_request() {
    local method=$1
    local endpoint=$2
    local data=$3
    local description=$4
    local use_file=$5
    
    echo -e "${YELLOW}Testing: $description${NC}"
    echo "Method: $method"
    echo "Endpoint: $endpoint"
    echo ""
    
    if [ "$method" = "GET" ]; then
        response=$(curl -s -w "\nHTTP_STATUS:%{http_code}" -X GET "$SERVER_URL$endpoint" \
            -H "Authorization: Bearer $JWT_TOKEN" \
            -H "Content-Type: application/json")
    else
        if [ "$use_file" = "true" ]; then
            # Use data as filename
            if [ -f "$TEST_DIR/$data" ]; then
                response=$(curl -s -w "\nHTTP_STATUS:%{http_code}" -X "$method" "$SERVER_URL$endpoint" \
                    -H "Authorization: Bearer $JWT_TOKEN" \
                    -H "Content-Type: application/json" \
                    -d @"$TEST_DIR/$data")
            else
                echo -e "${RED}ERROR: Data file $data not found${NC}"
                return 1
            fi
        else
            # Use data as JSON string
            response=$(curl -s -w "\nHTTP_STATUS:%{http_code}" -X "$method" "$SERVER_URL$endpoint" \
                -H "Authorization: Bearer $JWT_TOKEN" \
                -H "Content-Type: application/json" \
                -d "$data")
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

# Function to generate test data for asset sensor
generate_test_data() {
    local asset_sensor_id=$1
    local sensor_type_name=$2
    local mac_address=$3
    
    local current_time=$(date -u +"%Y-%m-%dT%H:%M:%S.%3NZ")
    
    case "$sensor_type_name" in
        "Temperature Sensor"|"temperature")
            cat <<EOF
{
    "asset_sensor_id": "$asset_sensor_id",
    "mac_address": "$mac_address",
    "measurement_data": {
        "temperature": $(awk -v min=18 -v max=35 'BEGIN{srand(); print min+rand()*(max-min)}'),
        "unit": "Celsius"
    },
    "reading_time": "$current_time"
}
EOF
            ;;
        "Humidity Sensor"|"humidity")
            cat <<EOF
{
    "asset_sensor_id": "$asset_sensor_id",
    "mac_address": "$mac_address",
    "measurement_data": {
        "humidity": $(awk -v min=30 -v max=80 'BEGIN{srand(); print min+rand()*(max-min)}'),
        "unit": "Percent"
    },
    "reading_time": "$current_time"
}
EOF
            ;;
        "Pressure Sensor"|"pressure")
            cat <<EOF
{
    "asset_sensor_id": "$asset_sensor_id",
    "mac_address": "$mac_address",
    "measurement_data": {
        "pressure": $(awk -v min=980 -v max=1050 'BEGIN{srand(); print min+rand()*(max-min)}'),
        "unit": "hPa"
    },
    "reading_time": "$current_time"
}
EOF
            ;;
        *)
            cat <<EOF
{
    "asset_sensor_id": "$asset_sensor_id",
    "mac_address": "$mac_address",
    "measurement_data": {
        "value": $(awk -v min=1 -v max=100 'BEGIN{srand(); print min+rand()*(max-min)}'),
        "unit": "generic"
    },
    "reading_time": "$current_time"
}
EOF
            ;;
    esac
}

# Function to generate batch test data
generate_batch_test_data() {
    local asset_sensors_json=$1
    local count=${2:-3}
    
    echo "["
    
    local sensors_array=$(echo "$asset_sensors_json" | jq -r '.[]')
    local sensor_count=$(echo "$asset_sensors_json" | jq '. | length')
    
    for i in $(seq 1 $count); do
        # Pick a random sensor
        local random_index=$(awk -v max=$sensor_count 'BEGIN{srand(); print int(rand()*max)}')
        local sensor=$(echo "$asset_sensors_json" | jq -r ".[$random_index]")
        
        local asset_sensor_id=$(echo "$sensor" | jq -r '.asset_sensor_id')
        local sensor_type_name=$(echo "$sensor" | jq -r '.sensor_type_name')
        local mac_address="AA:BB:CC:DD:EE:$(printf "%02X" $((RANDOM % 256)))"
        
        generate_test_data "$asset_sensor_id" "$sensor_type_name" "$mac_address"
        
        if [ $i -lt $count ]; then
            echo ","
        fi
    done
    
    echo "]"
}

# Main execution starts here
echo -e "${BLUE}Step 1: Fetching data from database${NC}"

# Fetch asset sensors
asset_sensors_json=$(fetch_asset_sensors)
if [ $? -ne 0 ]; then
    echo -e "${RED}Failed to fetch asset sensors. Exiting.${NC}"
    exit 1
fi

# Fetch sensor types
sensor_types_json=$(fetch_sensor_types)
if [ $? -ne 0 ]; then
    echo -e "${RED}Failed to fetch sensor types. Exiting.${NC}"
    exit 1
fi

echo -e "${GREEN}Database scraping completed!${NC}"
echo ""

# Display fetched data
echo -e "${BLUE}Found Asset Sensors:${NC}"
echo "$asset_sensors_json" | jq '.[] | {asset_sensor_name, asset_name, location_name, sensor_type_name, status}' 2>/dev/null || echo "$asset_sensors_json"
echo ""

echo -e "${BLUE}Found Sensor Types:${NC}"
echo "$sensor_types_json" | jq '.[] | {name, manufacturer, model}' 2>/dev/null || echo "$sensor_types_json"
echo ""

# Check if we have any asset sensors
sensor_count=$(echo "$asset_sensors_json" | jq '. | length' 2>/dev/null || echo "0")
if [ "$sensor_count" = "0" ] || [ "$sensor_count" = "null" ]; then
    echo -e "${RED}No active asset sensors found in database. Please add some asset sensors first.${NC}"
    exit 1
fi

echo -e "${BLUE}Step 2: Testing API endpoints with real data${NC}"
echo ""

# Test 1: Get JSON Template for first sensor type
first_sensor_type=$(echo "$sensor_types_json" | jq -r '.[0].name // "temperature"')
make_request "GET" "/api/v1/superadmin/iot-sensor-readings/template?sensor_type=$first_sensor_type" "" "Get JSON Template for $first_sensor_type"

# Test 2: Create readings using real asset sensor data
echo -e "${BLUE}Step 3: Creating test readings with scraped data${NC}"

# Get first 3 asset sensors for testing
for i in $(seq 0 2); do
    sensor=$(echo "$asset_sensors_json" | jq -r ".[$i] // empty")
    if [ ! -z "$sensor" ] && [ "$sensor" != "null" ]; then
        asset_sensor_id=$(echo "$sensor" | jq -r '.asset_sensor_id')
        asset_sensor_name=$(echo "$sensor" | jq -r '.asset_sensor_name')
        sensor_type_name=$(echo "$sensor" | jq -r '.sensor_type_name // "temperature"')
        mac_address="AA:BB:CC:DD:EE:$(printf "%02X" $((i + 1)))"
        
        # Generate test data
        test_data=$(generate_test_data "$asset_sensor_id" "$sensor_type_name" "$mac_address")
        
        make_request "POST" "/api/v1/superadmin/iot-sensor-readings/from-json" "$test_data" "Create Reading for $asset_sensor_name ($sensor_type_name)" "false"
    fi
done

# Test 3: Create batch readings
echo -e "${BLUE}Step 4: Creating batch readings${NC}"
batch_data=$(generate_batch_test_data "$asset_sensors_json" 3)
make_request "POST" "/api/v1/superadmin/iot-sensor-readings/batch-from-json" "$batch_data" "Create Batch Readings" "false"

# Test 4: Create simple reading using first asset sensor
first_sensor=$(echo "$asset_sensors_json" | jq -r '.[0]')
if [ ! -z "$first_sensor" ] && [ "$first_sensor" != "null" ]; then
    asset_sensor_id=$(echo "$first_sensor" | jq -r '.asset_sensor_id')
    sensor_type_id=$(echo "$sensor_types_json" | jq -r '.[0].id')
    
    simple_data="{
        \"asset_sensor_id\": \"$asset_sensor_id\",
        \"sensor_type_id\": \"$sensor_type_id\",
        \"mac_address\": \"AA:BB:CC:DD:EE:FF\",
        \"measurement_data\": {
            \"value\": 25.5,
            \"unit\": \"test\"
        }
    }"
    
    make_request "POST" "/api/v1/superadmin/iot-sensor-readings/simple" "$simple_data" "Create Simple Reading" "false"
fi

# Test 5: Create dummy reading
if [ ! -z "$first_sensor" ] && [ "$first_sensor" != "null" ]; then
    sensor_type_id=$(echo "$sensor_types_json" | jq -r '.[0].id')
    
    dummy_data="{
        \"sensor_type_id\": \"$sensor_type_id\",
        \"count\": 2
    }"
    
    make_request "POST" "/api/v1/superadmin/iot-sensor-readings/dummy/multiple" "$dummy_data" "Create Multiple Dummy Readings" "false"
fi

echo -e "${GREEN}=== Testing Complete ===${NC}"
echo ""
echo -e "${BLUE}Summary:${NC}"
echo "- Scraped $sensor_count asset sensors from database"
echo "- Tested readings creation with real asset sensor IDs"
echo "- Created individual and batch readings"
echo "- Used dynamic MAC addresses and realistic sensor data"
echo ""
echo -e "${YELLOW}Next Steps:${NC}"
echo "1. Check the database to verify readings were created"
echo "2. Verify location information was automatically populated"
echo "3. Test retrieval endpoints to fetch the created readings"
echo "4. Review server logs for any errors"

# Optional: Show recent readings
echo ""
echo -e "${PURPLE}Recent readings in database:${NC}"
recent_query="
SELECT 
    isr.id,
    isr.mac_address,
    isr.location as location_name,
    asn.name as asset_sensor_name,
    isr.reading_time,
    isr.measurement_data::text
FROM iot_sensor_readings isr
JOIN asset_sensors asn ON isr.asset_sensor_id = asn.id
ORDER BY isr.created_at DESC
LIMIT 5;
"

PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -c "$recent_query" 2>/dev/null || echo "Could not fetch recent readings"
