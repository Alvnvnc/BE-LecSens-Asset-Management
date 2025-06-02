# Asset with Sensors API Documentation

Fitur ini memungkinkan pembuatan asset beserta sensor-sensor yang terkait secara otomatis dalam satu request API. Hal ini menyederhanakan proses pembuatan asset di frontend karena tidak perlu melakukan multiple API calls.

## Overview

Sebelumnya, untuk membuat asset dengan sensors, frontend harus:
1. Membuat asset terlebih dahulu
2. Untuk setiap sensor type yang diinginkan, membuat asset sensor secara terpisah
3. Mengelola relasi antara asset dan sensors secara manual

Dengan fitur baru ini, semua proses tersebut dapat dilakukan dalam satu API call.

## API Endpoints

### 1. Create Asset with Sensors

**Endpoint:** `POST /api/v1/assets-with-sensors`

**Description:** Membuat asset baru beserta sensors yang terkait secara otomatis.

**Request Body:**
```json
{
  "name": "Temperature Monitoring Station",
  "asset_type_id": "123e4567-e89b-12d3-a456-426614174000",
  "location_id": "123e4567-e89b-12d3-a456-426614174001",
  "status": "active",
  "properties": {
    "description": "Main monitoring station for warehouse temperature",
    "installation_date": "2024-01-15"
  },
  "sensor_types": [
    {
      "sensor_type_id": "123e4567-e89b-12d3-a456-426614174002",
      "name": "Primary Temperature Sensor",
      "status": "active",
      "configuration": {
        "sampling_rate": 60,
        "unit": "celsius",
        "precision": 2
      }
    },
    {
      "sensor_type_id": "123e4567-e89b-12d3-a456-426614174003",
      "name": "Humidity Sensor",
      "status": "active",
      "configuration": {
        "sampling_rate": 300,
        "unit": "percentage",
        "range": "0-100"
      }
    }
  ]
}
```

**Response:**
```json
{
  "success": true,
  "message": "Asset with sensors created successfully",
  "data": {
    "asset": {
      "id": "123e4567-e89b-12d3-a456-426614174004",
      "name": "Temperature Monitoring Station",
      "asset_type_id": "123e4567-e89b-12d3-a456-426614174000",
      "location_id": "123e4567-e89b-12d3-a456-426614174001",
      "tenant_id": "123e4567-e89b-12d3-a456-426614174005",
      "status": "active",
      "properties": "{\"description\":\"Main monitoring station for warehouse temperature\",\"installation_date\":\"2024-01-15\"}",
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T10:30:00Z"
    },
    "sensors": [
      {
        "id": "123e4567-e89b-12d3-a456-426614174006",
        "tenant_id": "123e4567-e89b-12d3-a456-426614174005",
        "asset_id": "123e4567-e89b-12d3-a456-426614174004",
        "sensor_type_id": "123e4567-e89b-12d3-a456-426614174002",
        "name": "Primary Temperature Sensor",
        "status": "active",
        "configuration": "{\"sampling_rate\":60,\"unit\":\"celsius\",\"precision\":2}",
        "created_at": "2024-01-15T10:30:00Z"
      },
      {
        "id": "123e4567-e89b-12d3-a456-426614174007",
        "tenant_id": "123e4567-e89b-12d3-a456-426614174005",
        "asset_id": "123e4567-e89b-12d3-a456-426614174004",
        "sensor_type_id": "123e4567-e89b-12d3-a456-426614174003",
        "name": "Humidity Sensor",
        "status": "active",
        "configuration": "{\"sampling_rate\":300,\"unit\":\"percentage\",\"range\":\"0-100\"}",
        "created_at": "2024-01-15T10:30:00Z"
      }
    ]
  }
}
```

### 2. Get Asset with Sensors

**Endpoint:** `GET /api/v1/assets-with-sensors/{id}`

**Description:** Mengambil data asset beserta semua sensors yang terkait.

**Response:**
```json
{
  "success": true,
  "message": "Asset with sensors retrieved successfully",
  "data": {
    "asset": {
      "id": "123e4567-e89b-12d3-a456-426614174004",
      "name": "Temperature Monitoring Station",
      "asset_type_id": "123e4567-e89b-12d3-a456-426614174000",
      "location_id": "123e4567-e89b-12d3-a456-426614174001",
      "tenant_id": "123e4567-e89b-12d3-a456-426614174005",
      "status": "active",
      "properties": "{\"description\":\"Main monitoring station for warehouse temperature\",\"installation_date\":\"2024-01-15\"}",
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T10:30:00Z"
    },
    "sensors": [
      {
        "id": "123e4567-e89b-12d3-a456-426614174006",
        "tenant_id": "123e4567-e89b-12d3-a456-426614174005",
        "asset_id": "123e4567-e89b-12d3-a456-426614174004",
        "sensor_type_id": "123e4567-e89b-12d3-a456-426614174002",
        "name": "Primary Temperature Sensor",
        "status": "active",
        "configuration": "{\"sampling_rate\":60,\"unit\":\"celsius\",\"precision\":2}",
        "last_reading_value": 23.5,
        "last_reading_time": "2024-01-15T12:00:00Z",
        "created_at": "2024-01-15T10:30:00Z"
      }
    ]
  }
}
```

## Request Validation

### Required Fields
- `name`: Nama asset (string, required)
- `asset_type_id`: ID tipe asset (UUID, required)
- `location_id`: ID lokasi (UUID, required)
- `sensor_types`: Array sensor types (required, minimal 1 item)

### Sensor Types Validation
Setiap item dalam `sensor_types` harus memiliki:
- `sensor_type_id`: ID tipe sensor (UUID, required)
- `name`: Nama sensor (string, required)
- `status`: Status sensor (string, optional, default: "active")
- `configuration`: Konfigurasi sensor (JSON, optional)

## Business Logic

### Tenant ID Inheritance
- Asset sensors akan otomatis mewarisi `tenant_id` dari asset yang dibuat
- Jika asset tidak memiliki `tenant_id`, maka asset sensors juga akan memiliki `tenant_id` null
- Hal ini memastikan konsistensi data dalam sistem multi-tenant

### Error Handling
- Jika pembuatan asset gagal, seluruh operasi akan dibatalkan
- Jika pembuatan beberapa sensors gagal, asset tetap akan dibuat dan sensors yang berhasil akan disimpan
- Error pada sensors individual akan dicatat dalam log tetapi tidak menggagalkan seluruh operasi

### Default Values
- Asset status default: "active"
- Sensor status default: "active"
- Timestamps akan diset otomatis saat pembuatan

## Error Responses

### 400 Bad Request
```json
{
  "success": false,
  "message": "Invalid request format",
  "error": "Key: 'CreateAssetWithSensorsRequest.Name' Error:Field validation for 'Name' failed on the 'required' tag"
}
```

### 404 Not Found
```json
{
  "success": false,
  "message": "Asset type not found",
  "error": "invalid asset type"
}
```

### 500 Internal Server Error
```json
{
  "success": false,
  "message": "Failed to create asset with sensors",
  "error": "database connection failed"
}
```

## Usage Examples

### Frontend Integration

```javascript
// Example: Create asset with multiple sensors
const createAssetWithSensors = async (assetData) => {
  try {
    const response = await fetch('/api/v1/assets-with-sensors', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`
      },
      body: JSON.stringify({
        name: assetData.name,
        asset_type_id: assetData.assetTypeId,
        location_id: assetData.locationId,
        status: assetData.status || 'active',
        properties: assetData.properties,
        sensor_types: assetData.sensorTypes.map(sensor => ({
          sensor_type_id: sensor.typeId,
          name: sensor.name,
          status: sensor.status || 'active',
          configuration: sensor.configuration
        }))
      })
    });
    
    const result = await response.json();
    
    if (result.success) {
      console.log('Asset created with ID:', result.data.asset.id);
      console.log('Sensors created:', result.data.sensors.length);
      return result.data;
    } else {
      throw new Error(result.message);
    }
  } catch (error) {
    console.error('Failed to create asset with sensors:', error);
    throw error;
  }
};
```

## Benefits

1. **Simplified Frontend**: Hanya perlu satu API call untuk membuat asset dengan sensors
2. **Atomic Operations**: Pembuatan asset dan sensors dalam satu transaksi
3. **Consistent Data**: Tenant ID dan timestamps konsisten antara asset dan sensors
4. **Error Handling**: Penanganan error yang lebih baik dengan partial success
5. **Performance**: Mengurangi jumlah round-trips ke server

## Migration from Existing API

Jika sudah menggunakan API terpisah untuk asset dan asset sensors, Anda dapat:

1. **Gradual Migration**: Gunakan API baru untuk asset baru, pertahankan API lama untuk asset existing
2. **Batch Update**: Buat script untuk mengkonversi data existing menggunakan API baru
3. **Hybrid Approach**: Gunakan API baru untuk create, API lama untuk update individual sensors

## Notes

- API ini kompatibel dengan sistem tenant yang ada
- Validasi sensor types dilakukan sebelum pembuatan asset
- Log detail tersedia untuk debugging dan monitoring
- Response format konsisten dengan API existing lainnya