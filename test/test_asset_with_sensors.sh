#!/bin/bash

# Test script for Asset with Sensors API
# This script tests the new combined asset and sensor creation endpoint

# Configuration
BASE_URL="http://localhost:8080/api/v1"
CONTENT_TYPE="Content-Type: application/json"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}=== Asset with Sensors API Test Script ===${NC}"
echo -e "${BLUE}Testing endpoint: ${BASE_URL}/assets-with-sensors${NC}"
echo ""

# Function to print test results
print_result() {
    local test_name="$1"
    local status_code="$2"
    local expected_code="$3"
    
    if [ "$status_code" -eq "$expected_code" ]; then
        echo -e "${GREEN}✓ $test_name - Status: $status_code (Expected: $expected_code)${NC}"
    else
        echo -e "${RED}✗ $test_name - Status: $status_code (Expected: $expected_code)${NC}"
    fi
}

# Test 1: Create Asset with Sensors - Valid Request
echo -e "${YELLOW}Test 1: Creating asset with sensors (Valid Request)${NC}"
response=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/assets-with-sensors" \
  -H "$CONTENT_TYPE" \
  -d @asset_with_sensors_example.json)

status_code=$(echo "$response" | tail -n1)
response_body=$(echo "$response" | head -n -1)

print_result "Create Asset with Sensors" "$status_code" "201"

if [ "$status_code" -eq "201" ]; then
    echo -e "${GREEN}Response:${NC}"
    echo "$response_body" | jq .
    
    # Extract asset ID for further tests
    ASSET_ID=$(echo "$response_body" | jq -r '.data.asset.id')
    echo -e "${BLUE}Created Asset ID: $ASSET_ID${NC}"
else
    echo -e "${RED}Error Response:${NC}"
    echo "$response_body" | jq .
fi

echo ""
echo "---"
echo ""

# Test 2: Create Asset with Sensors - Invalid Request (Missing Name)
echo -e "${YELLOW}Test 2: Creating asset with sensors (Missing Name)${NC}"
invalid_request='{
  "asset_type_id": "550e8400-e29b-41d4-a716-446655440001",
  "location_id": "550e8400-e29b-41d4-a716-446655440002",
  "sensor_types": [
    {
      "sensor_type_id": "550e8400-e29b-41d4-a716-446655440003",
      "name": "Test Sensor"
    }
  ]
}'

response=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/assets-with-sensors" \
  -H "$CONTENT_TYPE" \
  -d "$invalid_request")

status_code=$(echo "$response" | tail -n1)
response_body=$(echo "$response" | head -n -1)

print_result "Create Asset (Missing Name)" "$status_code" "400"

if [ "$status_code" -eq "400" ]; then
    echo -e "${GREEN}Expected error response:${NC}"
    echo "$response_body" | jq .
else
    echo -e "${RED}Unexpected response:${NC}"
    echo "$response_body" | jq .
fi

echo ""
echo "---"
echo ""

# Test 3: Create Asset with Sensors - Invalid Request (Empty Sensor Types)
echo -e "${YELLOW}Test 3: Creating asset with sensors (Empty Sensor Types)${NC}"
invalid_request2='{
  "name": "Test Asset",
  "asset_type_id": "550e8400-e29b-41d4-a716-446655440001",
  "location_id": "550e8400-e29b-41d4-a716-446655440002",
  "sensor_types": []
}'

response=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/assets-with-sensors" \
  -H "$CONTENT_TYPE" \
  -d "$invalid_request2")

status_code=$(echo "$response" | tail -n1)
response_body=$(echo "$response" | head -n -1)

print_result "Create Asset (Empty Sensors)" "$status_code" "400"

if [ "$status_code" -eq "400" ]; then
    echo -e "${GREEN}Expected error response:${NC}"
    echo "$response_body" | jq .
else
    echo -e "${RED}Unexpected response:${NC}"
    echo "$response_body" | jq .
fi

echo ""
echo "---"
echo ""

# Test 4: Get Asset with Sensors (if asset was created successfully)
if [ ! -z "$ASSET_ID" ] && [ "$ASSET_ID" != "null" ]; then
    echo -e "${YELLOW}Test 4: Getting asset with sensors${NC}"
    
    response=$(curl -s -w "\n%{http_code}" -X GET "$BASE_URL/assets-with-sensors/$ASSET_ID" \
      -H "$CONTENT_TYPE")
    
    status_code=$(echo "$response" | tail -n1)
    response_body=$(echo "$response" | head -n -1)
    
    print_result "Get Asset with Sensors" "$status_code" "200"
    
    if [ "$status_code" -eq "200" ]; then
        echo -e "${GREEN}Response:${NC}"
        echo "$response_body" | jq .
        
        # Count sensors
        sensor_count=$(echo "$response_body" | jq '.data.sensors | length')
        echo -e "${BLUE}Number of sensors found: $sensor_count${NC}"
    else
        echo -e "${RED}Error Response:${NC}"
        echo "$response_body" | jq .
    fi
else
    echo -e "${YELLOW}Test 4: Skipped (No valid asset ID from previous test)${NC}"
fi

echo ""
echo "---"
echo ""

# Test 5: Get Asset with Sensors - Invalid ID
echo -e "${YELLOW}Test 5: Getting asset with sensors (Invalid ID)${NC}"
invalid_id="invalid-uuid-format"

response=$(curl -s -w "\n%{http_code}" -X GET "$BASE_URL/assets-with-sensors/$invalid_id" \
  -H "$CONTENT_TYPE")

status_code=$(echo "$response" | tail -n1)
response_body=$(echo "$response" | head -n -1)

print_result "Get Asset (Invalid ID)" "$status_code" "400"

if [ "$status_code" -eq "400" ]; then
    echo -e "${GREEN}Expected error response:${NC}"
    echo "$response_body" | jq .
else
    echo -e "${RED}Unexpected response:${NC}"
    echo "$response_body" | jq .
fi

echo ""
echo "---"
echo ""

# Test 6: Get Asset with Sensors - Non-existent ID
echo -e "${YELLOW}Test 6: Getting asset with sensors (Non-existent ID)${NC}"
nonexistent_id="550e8400-e29b-41d4-a716-446655440999"

response=$(curl -s -w "\n%{http_code}" -X GET "$BASE_URL/assets-with-sensors/$nonexistent_id" \
  -H "$CONTENT_TYPE")

status_code=$(echo "$response" | tail -n1)
response_body=$(echo "$response" | head -n -1)

print_result "Get Asset (Non-existent)" "$status_code" "404"

if [ "$status_code" -eq "404" ]; then
    echo -e "${GREEN}Expected error response:${NC}"
    echo "$response_body" | jq .
else
    echo -e "${RED}Unexpected response:${NC}"
    echo "$response_body" | jq .
fi

echo ""
echo -e "${BLUE}=== Test Summary ===${NC}"
echo -e "${GREEN}✓ Tests completed${NC}"
echo -e "${BLUE}Note: Make sure the server is running and the database contains valid asset_type_id and location_id${NC}"
echo -e "${BLUE}Note: Update the UUIDs in asset_with_sensors_example.json to match your database${NC}"
echo ""