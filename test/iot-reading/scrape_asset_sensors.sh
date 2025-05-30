#!/bin/bash

# Database Asset Sensor Scraper
# This script fetches asset sensor data and generates test JSON files

# Database configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=alvn
DB_PASSWORD=alvn12345
DB_NAME=asset_management

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
NC='\033[0m' # No Color

TEST_DIR="$(dirname "$0")"
SCRAPED_DATA_DIR="$TEST_DIR/scraped_data"

echo -e "${BLUE}=== Asset Sensor Database Scraper ===${NC}"
echo ""

# Create scraped_data directory if it doesn't exist
if [ ! -d "$SCRAPED_DATA_DIR" ]; then
    mkdir -p "$SCRAPED_DATA_DIR"
    echo -e "${GREEN}Created directory: $SCRAPED_DATA_DIR${NC}"
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
    local result=$(PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -t -c "$query" 2>/dev/null)
    echo "$result"
}

# Function to fetch asset sensors from database
fetch_asset_sensors() {
    echo -e "${PURPLE}Fetching asset sensors from database...${NC}" >&2
    
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
                'manufacturer', st.manufacturer,
                'model', st.model,
                'status', asn.status,
                'mac_address', COALESCE(asn.configuration->>'mac_address', 'AA:BB:CC:DD:EE:' || LPAD(UPPER(TO_HEX(EXTRACT(epoch FROM NOW())::int % 256)), 2, '0'))
            )
        )
    FROM asset_sensors asn
    JOIN assets a ON asn.asset_id = a.id
    JOIN locations l ON a.location_id = l.id
    JOIN sensor_types st ON asn.sensor_type_id = st.id
    WHERE asn.status = 'active'
    LIMIT 10;
    "
    
    local result=$(execute_sql "$query")
    
    if [ $? -eq 0 ] && [ ! -z "$result" ] && [ "$result" != "null" ]; then
        echo "$result" | tr -d '\n' | sed 's/^[ \t]*//;s/[ \t]*$//'
    else
        echo -e "${RED}Failed to fetch asset sensors from database${NC}" >&2
        echo -e "Please check your database connection settings" >&2
        return 1
    fi
}

# Function to generate temperature reading
generate_temperature_reading() {
    local asset_sensor_id=$1
    local mac_address=$2
    local current_time=$(date -u +"%Y-%m-%dT%H:%M:%S.%3NZ")
    local temp=$(awk -v min=18 -v max=35 'BEGIN{srand(); printf "%.2f", min+rand()*(max-min)}')
    
    cat > "$SCRAPED_DATA_DIR/single_temperature_scraped.json" <<EOF
{
    "asset_sensor_id": "$asset_sensor_id",
    "mac_address": "$mac_address",
    "measurement_data": {
        "temperature": $temp,
        "unit": "Celsius"
    },
    "reading_time": "$current_time"
}
EOF
}

# Function to generate humidity reading
generate_humidity_reading() {
    local asset_sensor_id=$1
    local mac_address=$2
    local current_time=$(date -u +"%Y-%m-%dT%H:%M:%S.%3NZ")
    local humidity=$(awk -v min=30 -v max=80 'BEGIN{srand(); printf "%.2f", min+rand()*(max-min)}')
    
    cat > "$SCRAPED_DATA_DIR/single_humidity_scraped.json" <<EOF
{
    "asset_sensor_id": "$asset_sensor_id",
    "mac_address": "$mac_address",
    "measurement_data": {
        "humidity": $humidity,
        "unit": "Percent"
    },
    "reading_time": "$current_time"
}
EOF
}

# Function to generate batch readings
generate_batch_readings() {
    local asset_sensors_json=$1
    local count=${2:-3}
    local current_time=$(date -u +"%Y-%m-%dT%H:%M:%S.%3NZ")
    
    echo "[" > "$SCRAPED_DATA_DIR/batch_scraped.json"
    
    local sensor_count=$(echo "$asset_sensors_json" | jq '. | length')
    
    for i in $(seq 1 $count); do
        local random_index=$(awk -v max=$sensor_count 'BEGIN{srand(); print int(rand()*max)}')
        local sensor=$(echo "$asset_sensors_json" | jq -r ".[$random_index]")
        
        local asset_sensor_id=$(echo "$sensor" | jq -r '.asset_sensor_id')
        local mac_address=$(echo "$sensor" | jq -r '.mac_address')
        local temp=$(awk -v min=18 -v max=35 'BEGIN{srand(); printf "%.2f", min+rand()*(max-min)}')
        
        cat >> "$SCRAPED_DATA_DIR/batch_scraped.json" <<EOF
    {
        "asset_sensor_id": "$asset_sensor_id",
        "mac_address": "$mac_address",
        "measurement_data": {
            "temperature": $temp,
            "unit": "Celsius"
        },
        "reading_time": "$current_time"
    }
EOF
        
        if [ $i -lt $count ]; then
            echo "," >> "$SCRAPED_DATA_DIR/batch_scraped.json"
        else
            echo "" >> "$SCRAPED_DATA_DIR/batch_scraped.json"
        fi
    done
    
    echo "]" >> "$SCRAPED_DATA_DIR/batch_scraped.json"
}

# Function to generate simple reading request
generate_simple_reading() {
    local asset_sensors_json=$1
    local sensor_types_json=$2
    
    local first_sensor=$(echo "$asset_sensors_json" | jq -r '.[0]')
    local asset_sensor_id=$(echo "$first_sensor" | jq -r '.asset_sensor_id')
    local sensor_type_id=$(echo "$first_sensor" | jq -r '.sensor_type_id')
    local mac_address=$(echo "$first_sensor" | jq -r '.mac_address')
    
    cat > "$SCRAPED_DATA_DIR/simple_reading_scraped.json" <<EOF
{
    "asset_sensor_id": "$asset_sensor_id",
    "sensor_type_id": "$sensor_type_id",
    "mac_address": "$mac_address",
    "measurement_data": {
        "value": 25.5,
        "unit": "test"
    }
}
EOF
}

# Function to generate dummy reading request
generate_dummy_reading() {
    local sensor_types_json=$1
    
    local sensor_type_id=$(echo "$sensor_types_json" | jq -r '.[0].sensor_type_id')
    
    cat > "$SCRAPED_DATA_DIR/dummy_scraped.json" <<EOF
{
    "sensor_type_id": "$sensor_type_id",
    "count": 3
}
EOF
}

# Function to update main test script
update_test_script() {
    local original_script="$TEST_DIR/test_script.sh"
    local backup_script="$TEST_DIR/test_script_backup.sh"
    
    if [ -f "$original_script" ]; then
        # Create backup
        cp "$original_script" "$backup_script"
        echo -e "${GREEN}Created backup: test_script_backup.sh${NC}"
        
        # Update the script to use scraped files
        sed -i 's/single_temperature\.json/single_temperature_scraped.json/g' "$original_script"
        sed -i 's/single_humidity\.json/single_humidity_scraped.json/g' "$original_script"
        sed -i 's/batch_temperature\.json/batch_scraped.json/g' "$original_script"
        sed -i 's/simple_reading_request\.json/simple_reading_scraped.json/g' "$original_script"
        sed -i 's/dummy_multiple_request\.json/dummy_scraped.json/g' "$original_script"
        
        echo -e "${GREEN}Updated test_script.sh to use scraped data files${NC}"
    fi
}

# Main execution
echo -e "${BLUE}Step 1: Fetching asset sensors from database${NC}"

# Fetch asset sensors
asset_sensors_json=$(fetch_asset_sensors)
if [ $? -ne 0 ]; then
    echo -e "${RED}Failed to fetch asset sensors. Exiting.${NC}"
    exit 1
fi

# Check if we have any asset sensors
sensor_count=$(echo "$asset_sensors_json" | jq '. | length' 2>/dev/null || echo "0")
if [ "$sensor_count" = "0" ] || [ "$sensor_count" = "null" ]; then
    echo -e "${RED}No active asset sensors found in database.${NC}"
    echo -e "${YELLOW}Please add some asset sensors first or check your database connection.${NC}"
    exit 1
fi

echo -e "${GREEN}Found $sensor_count asset sensors${NC}"
echo ""

# Display fetched data
echo -e "${BLUE}Asset Sensors Found:${NC}"
echo "$asset_sensors_json" | jq '.[] | {asset_sensor_name, asset_name, location_name, sensor_type_name, status}' 2>/dev/null || echo "$asset_sensors_json"
echo ""

echo -e "${BLUE}Step 2: Generating test data files${NC}"

# Generate test files
first_sensor=$(echo "$asset_sensors_json" | jq -r '.[0]')
asset_sensor_id=$(echo "$first_sensor" | jq -r '.asset_sensor_id')
mac_address=$(echo "$first_sensor" | jq -r '.mac_address')

# Generate various test files
generate_temperature_reading "$asset_sensor_id" "$mac_address"
echo -e "${GREEN}Generated: single_temperature_scraped.json${NC}"

generate_humidity_reading "$asset_sensor_id" "$mac_address"
echo -e "${GREEN}Generated: single_humidity_scraped.json${NC}"

generate_batch_readings "$asset_sensors_json" 3
echo -e "${GREEN}Generated: batch_scraped.json${NC}"

generate_simple_reading "$asset_sensors_json" ""
echo -e "${GREEN}Generated: simple_reading_scraped.json${NC}"

generate_dummy_reading "$asset_sensors_json"
echo -e "${GREEN}Generated: dummy_scraped.json${NC}"

echo ""
echo -e "${BLUE}Step 3: Updating test script${NC}"
update_test_script

echo ""
echo -e "${GREEN}=== Scraping Complete ===${NC}"
echo -e "${BLUE}Summary:${NC}"
echo "- Scraped $sensor_count asset sensors from database"
echo "- Generated 5 test data files with real asset sensor IDs"
echo "- Updated test_script.sh to use scraped data"
echo ""
echo -e "${YELLOW}Generated Files (in scraped_data/ directory):${NC}"
echo "- single_temperature_scraped.json"
echo "- single_humidity_scraped.json"
echo "- batch_scraped.json"
echo "- simple_reading_scraped.json"
echo "- dummy_scraped.json"
echo ""
echo -e "${YELLOW}Next Steps:${NC}"
echo "1. Run ./test_script.sh to test with real database data"
echo "2. Check that location information is automatically populated"
echo "3. Verify readings are created successfully"

# Show sample of generated file
echo ""
echo -e "${BLUE}Sample generated file (single_temperature_scraped.json):${NC}"
cat "$SCRAPED_DATA_DIR/single_temperature_scraped.json" | jq '.' 2>/dev/null || cat "$SCRAPED_DATA_DIR/single_temperature_scraped.json"