# Location Management API Documentation

## Overview
API untuk manajemen lokasi dengan akses kontrol berdasarkan role. Role SuperAdmin memiliki akses penuh untuk CRUD operations, sementara public routes tersedia untuk operasi read.

## Endpoints

### Public Endpoints (Read Access)

#### 1. Get All Locations
```
GET /api/v1/locations
```

**Query Parameters:**
- `page` (optional): Nomor halaman (default: 1)
- `page_size` (optional): Jumlah item per halaman (default: 10)

**Response:**
```json
[
  {
    "id": "uuid",
    "region_code": "string",
    "name": "string",
    "description": "string",
    "address": "string",
    "longitude": 0.0,
    "latitude": 0.0,
    "hierarchy_level": 1,
    "is_active": true,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
]
```

#### 2. Get Location by ID
```
GET /api/v1/locations/{id}
```

**Response:**
```json
{
  "id": "uuid",
  "region_code": "string",
  "name": "string",
  "description": "string",
  "address": "string",
  "longitude": 0.0,
  "latitude": 0.0,
  "hierarchy_level": 1,
  "is_active": true,
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

### SuperAdmin Endpoints (Full CRUD Access)

#### 3. Create Location
```
POST /api/v1/superadmin/locations
```

**Headers:**
- `Authorization: Bearer <superadmin_token>`

**Request Body:**
```json
{
  "region_code": "string",
  "name": "string",
  "description": "string",
  "address": "string",
  "longitude": 0.0,
  "latitude": 0.0,
  "hierarchy_level": 1,
  "is_active": true
}
```

**Response (201 Created):**
```json
{
  "id": "uuid",
  "region_code": "string",
  "name": "string",
  "description": "string",
  "address": "string",
  "longitude": 0.0,
  "latitude": 0.0,
  "hierarchy_level": 1,
  "is_active": true,
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

#### 4. Update Location
```
PUT /api/v1/superadmin/locations/{id}
```

**Headers:**
- `Authorization: Bearer <superadmin_token>`

**Request Body:**
```json
{
  "region_code": "string",
  "name": "string",
  "description": "string",
  "address": "string",
  "longitude": 0.0,
  "latitude": 0.0,
  "hierarchy_level": 1,
  "is_active": true
}
```

**Response (200 OK):**
```json
{
  "id": "uuid",
  "region_code": "string",
  "name": "string",
  "description": "string",
  "address": "string",
  "longitude": 0.0,
  "latitude": 0.0,
  "hierarchy_level": 1,
  "is_active": true,
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

#### 5. Delete Location
```
DELETE /api/v1/superadmin/locations/{id}
```

**Headers:**
- `Authorization: Bearer <superadmin_token>`

**Response (200 OK):**
```json
{
  "message": "Location deleted successfully"
}
```

## Error Responses

### 400 Bad Request
```json
{
  "error": "Invalid location ID format"
}
```

### 401 Unauthorized
```json
{
  "error": "Unauthorized"
}
```

### 403 Forbidden
```json
{
  "error": "Forbidden - SuperAdmin access required"
}
```

### 404 Not Found
```json
{
  "error": "location not found: {id}"
}
```

### 500 Internal Server Error
```json
{
  "error": "Failed to create location: {error_details}"
}
```

## Location Entity Fields

- **id**: UUID primary key (auto-generated)
- **region_code**: Code representing the region (province/city/district)
- **name**: Location name (required)
- **description**: Optional description
- **address**: Optional address
- **longitude**: Geographic longitude coordinate
- **latitude**: Geographic latitude coordinate  
- **hierarchy_level**: Level in location hierarchy (1=Province, 2=Kabupaten, 3=Kota, etc.)
- **is_active**: Boolean flag for location status
- **created_at**: Timestamp when location was created
- **updated_at**: Timestamp when location was last updated

## Notes

1. **Partial Updates**: Update endpoint supports partial updates - hanya field yang dikirim yang akan diupdate
2. **Validation**: Field `name` dan `hierarchy_level` adalah required
3. **Coordinates**: Longitude dan Latitude menggunakan format decimal degrees
4. **Hierarchy**: Hierarchy level membantu mengorganisir lokasi dalam struktur hierarkis
5. **Soft Delete**: Saat ini menggunakan hard delete, pertimbangkan implementasi soft delete untuk data integrity
