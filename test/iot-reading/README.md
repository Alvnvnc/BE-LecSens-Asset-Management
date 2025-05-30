# IoT Sensor Reading Testing Scripts

Script-script ini digunakan untuk testing IoT sensor reading dengan data yang diambil langsung dari database.

## File-file yang tersedia:

### 1. `test_script.sh` (Original)
Script testing manual dengan data statis yang sudah disiapkan.

### 2. `test_script_with_db.sh` (Enhanced)
Script testing lengkap yang dapat mengambil data langsung dari database dan melakukan testing secara otomatis.

### 3. `scrape_asset_sensors.sh`
Script terpisah untuk mengambil data asset sensor dari database dan menggenerate file JSON untuk testing.

### 4. `db_config.sh.example`
Template konfigurasi database yang dapat disesuaikan.

## Cara Penggunaan:

### Persiapan:

1. **Install dependencies yang diperlukan:**
   ```bash
   # Ubuntu/Debian
   sudo apt install postgresql-client jq
   
   # macOS
   brew install postgresql jq
   ```

2. **Setup konfigurasi database:**
   ```bash
   cp db_config.sh.example db_config.sh
   # Edit db_config.sh dengan setting database Anda
   ```

3. **Pastikan database dan server berjalan:**
   ```bash
   # Start database (jika menggunakan Docker)
   docker-compose up -d postgres
   
   # Start aplikasi
   go run main.go
   ```

### Opsi 1: Testing dengan Database Scraping Otomatis

Gunakan script yang enhanced untuk testing otomatis dengan data real dari database:

```bash
cd /home/alvn/Documents/playground/kp/be-lecsens/asset_management/test/iot-reading
./test_script_with_db.sh
```

Script ini akan:
- Mengambil data asset sensor dari database
- Menggenerate data test yang realistis
- Melakukan testing semua endpoint
- Menampilkan hasil testing
- Menampilkan data terbaru dari database

### Opsi 2: Generate Test Files + Manual Testing

1. **Generate test files dari database:**
   ```bash
   ./scrape_asset_sensors.sh
   ```
   
   Script ini akan:
   - Mengambil data asset sensor dari database
   - Menggenerate file JSON untuk testing
   - Mengupdate script testing untuk menggunakan file yang baru
   
2. **Jalankan testing manual:**
   ```bash
   ./test_script.sh
   ```

### Opsi 3: Testing dengan Data Statis

Jika Anda ingin menggunakan data yang sudah disiapkan manual:

```bash
./test_script.sh
```

## File yang Dihasilkan:

Setelah menjalankan `scrape_asset_sensors.sh`, file-file berikut akan dibuat:

- `single_temperature_scraped.json` - Data single reading temperature
- `single_humidity_scraped.json` - Data single reading humidity  
- `batch_scraped.json` - Data batch readings
- `simple_reading_scraped.json` - Data simple reading
- `dummy_scraped.json` - Data dummy reading request

## Konfigurasi Database:

Edit file `db_config.sh` (copy dari `db_config.sh.example`) untuk menyesuaikan:

```bash
# Database connection settings
export DB_HOST="localhost"
export DB_PORT="5432"
export DB_NAME="your_database_name"
export DB_USER="your_username"
export DB_PASSWORD="your_password"

# API Testing settings
export SERVER_URL="http://localhost:3160"
export JWT_TOKEN="your_jwt_token_here"
```

## Testing Endpoints:

Script akan melakukan testing pada endpoint-endpoint berikut:

1. `GET /api/v1/superadmin/iot-sensor-readings/template` - Get JSON template
2. `POST /api/v1/superadmin/iot-sensor-readings/from-json` - Create single reading
3. `POST /api/v1/superadmin/iot-sensor-readings/batch-from-json` - Create batch readings
4. `POST /api/v1/superadmin/iot-sensor-readings/simple` - Create simple reading
5. `POST /api/v1/superadmin/iot-sensor-readings/dummy/multiple` - Create dummy readings

## Validasi Automatic Location:

Script akan membantu memvalidasi fitur automatic location yang baru diimplementasi:

- ✅ Location tidak perlu diinput manual lagi
- ✅ Location diambil otomatis dari Asset → Location
- ✅ LocationID dan LocationName dipopulate otomatis
- ✅ Data konsisten antara asset sensor dan reading

## Troubleshooting:

### Database Connection Error:
```bash
# Test koneksi database
PGPASSWORD="your_password" psql -h localhost -p 5432 -U your_user -d your_db -c "SELECT 1;"
```

### JWT Token Expired:
- Login ulang ke aplikasi untuk mendapatkan token baru
- Update JWT_TOKEN di script atau db_config.sh

### Permission Denied:
```bash
chmod +x *.sh
```

### Missing Dependencies:
```bash
# Check if tools are installed
which psql jq curl
```

## Sample Output:

Ketika script berjalan dengan sukses, Anda akan melihat output seperti:

```
=== Asset Sensor Database Scraper ===

Found 5 asset sensors
Asset Sensors Found:
{
  "asset_sensor_name": "Temperature Sensor 1",
  "asset_name": "Production Line A",
  "location_name": "Factory Floor 1",
  "sensor_type_name": "Temperature Sensor",
  "status": "active"
}

Generated: single_temperature_scraped.json
Generated: batch_scraped.json
...

Testing: Create Single Temperature Reading
Method: POST
Status: 201
Response: {
  "success": true,
  "message": "IoT sensor reading created successfully",
  "data": {
    "id": "...",
    "location": "Factory Floor 1"  // ← Automatically populated!
  }
}
```

## Notes:

- Script akan menggunakan data asset sensor yang statusnya 'active'
- MAC address akan di-generate secara otomatis jika tidak ada
- Reading time menggunakan timestamp saat ini
- Measurement data di-generate dengan nilai random yang realistis
- Location akan otomatis ter-populate dari relasi Asset → Location

## Available Endpoints

### 6. Get Batch JSON Template
**GET** `/api/v1/superadmin/iot-sensor-readings/template/batch?sensor_type=temperature&count=3`

### 7. Create Simple Reading
**POST** `/api/v1/superadmin/iot-sensor-readings/simple`

## Test Files

### Single Reading Tests
- `single_temperature.json` - Valid temperature sensor reading
- `single_humidity.json` - Valid humidity sensor reading
- `single_pressure.json` - Valid pressure sensor reading
- `single_light.json` - Valid light sensor reading
- `single_sound.json` - Valid sound sensor reading
- `single_motion.json` - Valid motion sensor reading
- `high_precision.json` - High precision reading with detailed metadata

### Batch Reading Tests
- `batch_temperature.json` - Multiple temperature readings from same sensor
- `batch_humidity.json` - Multiple humidity readings from same sensor
- `batch_mixed_sensors.json` - Mixed sensor types in one batch
- `extreme_values.json` - Extreme temperature and humidity values

### Invalid Data Tests
- `invalid_uuid.json` - Invalid UUID format
- `missing_required_fields.json` - Missing or empty required fields
- `invalid_value_type.json` - Invalid data type for reading_value
- `batch_partial_invalid.json` - Batch with some valid and some invalid readings

### Request Body Tests
- `dummy_temperature_request.json` - Request for dummy temperature readings
- `dummy_multiple_request.json` - Request for multiple dummy readings
- `simple_reading_request.json` - Simple reading creation request

## Testing Commands

### Using curl (Linux/macOS):

#### 1. Create Single Reading
```bash
curl -X POST http://localhost:8080/api/v1/superadmin/iot-sensor-readings/from-json \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d @test/iot-reading/single_temperature.json
```

#### 2. Create Batch Readings
```bash
curl -X POST http://localhost:8080/api/v1/superadmin/iot-sensor-readings/batch-from-json \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d @test/iot-reading/batch_temperature.json
```

#### 3. Create Dummy Reading
```bash
curl -X POST http://localhost:8080/api/v1/superadmin/iot-sensor-readings/dummy \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d @test/iot-reading/dummy_temperature_request.json
```

#### 4. Create Multiple Dummy Readings
```bash
curl -X POST http://localhost:8080/api/v1/superadmin/iot-sensor-readings/dummy/multiple \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d @test/iot-reading/dummy_multiple_request.json
```

#### 5. Get JSON Template
```bash
curl -X GET "http://localhost:8080/api/v1/superadmin/iot-sensor-readings/template?sensor_type=temperature" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

#### 6. Get Batch JSON Template
```bash
curl -X GET "http://localhost:8080/api/v1/superadmin/iot-sensor-readings/template/batch?sensor_type=temperature&count=3" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

#### 7. Create Simple Reading
```bash
curl -X POST http://localhost:8080/api/v1/superadmin/iot-sensor-readings/simple \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d @test/iot-reading/simple_reading_request.json
```

### Testing Invalid Data
Test error handling by using invalid data files:

```bash
# Test invalid UUID
curl -X POST http://localhost:8080/api/v1/superadmin/iot-sensor-readings/from-json \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d @test/iot-reading/invalid_uuid.json

# Test missing fields
curl -X POST http://localhost:8080/api/v1/superadmin/iot-sensor-readings/from-json \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d @test/iot-reading/missing_required_fields.json
```

## Expected Responses

### Success Response (201 Created)
```json
{
  "status": "success",
  "message": "IoT sensor reading created successfully",
  "data": {
    "id": "uuid-here",
    "sensor_id": "sensor-uuid",
    "sensor_type": "temperature",
    "reading_value": 25.5,
    "unit": "°C",
    "location": "Server Room A",
    "timestamp": "2025-05-29T10:30:00Z",
    "metadata": {...},
    "created_at": "2025-05-29T10:30:00Z",
    "updated_at": "2025-05-29T10:30:00Z"
  }
}
```

### Batch Success Response (201 Created)
```json
{
  "status": "success",
  "message": "3 IoT sensor readings created successfully",
  "data": [
    {...},
    {...},
    {...}
  ]
}
```

### Error Response (400 Bad Request)
```json
{
  "status": "error",
  "message": "Invalid UUID format",
  "error": "uuid: invalid UUID length: 12"
}
```

## Notes

1. All endpoints require SuperAdmin authorization
2. Replace `YOUR_JWT_TOKEN` with a valid JWT token
3. Make sure the server is running on `localhost:8080` (adjust port if different)
4. The server should have proper CORS configuration for browser testing
5. Check server logs for detailed error information if requests fail

## Supported Sensor Types

- temperature
- humidity
- pressure
- light
- sound
- motion
- vibration
- gas
- proximity
- acceleration

## Required Fields

- `sensor_id` (UUID format)
- `sensor_type` (string)
- `reading_value` (number)
- `unit` (string)
- `location` (string)
- `timestamp` (ISO 8601 format)
