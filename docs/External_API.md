# External API untuk Integrasi Antar Service

Dokumen ini menjelaskan API yang tersedia untuk integrasi antar service, khususnya untuk validasi role user termasuk SuperAdmin.

## Endpoint yang Tersedia

### 1. Validasi JWT Token

**Endpoint:** `POST /api/external/auth/validate-token`

**Request:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Response:**
```json
{
  "isValid": true,
  "userInfo": {
    "userID": "user-uuid",
    "username": "johndoe",
    "email": "john@example.com",
    "role": "Admin",
    "roleID": "role-uuid",
    "status": "active",
    "isActive": true
  }
}
```

### 2. Mendapatkan Info User dari Token

**Endpoint:** `POST /api/external/auth/user-info`

**Header:**
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Response:**
```json
{
  "userID": "user-uuid",
  "email": "john@example.com",
  "userRole": "Admin"
}
```

### 3. Validasi User Permissions

**Endpoint:** `POST /api/external/auth/validate-user-permissions`

**Request:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "tenantID": "tenant-uuid",
  "requiredRole": "Admin",
  "permissions": ["read", "write"]
}
```

**Response:**
```json
{
  "valid": true,
  "hasRolePermission": true,
  "hasTenantAccess": true,
  "userID": "user-uuid",
  "userRole": "Admin"
}
```

### 4. Validasi Role SuperAdmin (Khusus)

**Endpoint:** `POST /api/external/auth/validate-superadmin`

**Request:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

Atau gunakan header:
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Response (Jika SuperAdmin):**
```json
{
  "valid": true,
  "userID": "user-uuid",
  "userRole": "SUPERADMIN",
  "isSuperAdmin": true
}
```

**Response (Jika Bukan SuperAdmin):**
```json
{
  "valid": false,
  "userID": "user-uuid",
  "userRole": "Admin",
  "isSuperAdmin": false
}
```

## Contoh Penggunaan dalam Go

```go
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	UserManagementURL = "http://user-management-service:3120"
	APIKey            = "your-api-key-here"
)

// Validasi apakah user adalah SuperAdmin
func ValidateSuperAdmin(token string) (bool, error) {
	url := UserManagementURL + "/api/external/auth/validate-superadmin"
	
	// Prepare request data
	reqData := map[string]string{"token": token}
	jsonData, err := json.Marshal(reqData)
	if err != nil {
		return false, err
	}
	
	// Create request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return false, err
	}
	
	// Add API Key header
	req.Header.Set("X-API-Key", APIKey)
	req.Header.Set("Content-Type", "application/json")
	
	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	
	// Read response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}
	
	// Parse response
	var result struct {
		Valid        bool   `json:"valid"`
		UserID       string `json:"userID"`
		UserRole     string `json:"userRole"`
		IsSuperAdmin bool   `json:"isSuperAdmin"`
	}
	
	if err := json.Unmarshal(body, &result); err != nil {
		return false, err
	}
	
	return result.Valid, nil
}

func main() {
	// Contoh penggunaan
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
	isSuperAdmin, err := ValidateSuperAdmin(token)
	if err != nil {
		fmt.Printf("Error validating SuperAdmin: %v\n", err)
		return
	}
	
	if isSuperAdmin {
		fmt.Println("User is a SuperAdmin, proceed with privileged operation")
	} else {
		fmt.Println("Access denied: SuperAdmin role required")
	}
}
```

## Catatan Keamanan

1. Selalu menggunakan API Key untuk komunikasi antar service
2. Gunakan HTTPS untuk enkripsi traffic
3. Validasi token sebelum mempercayai data user
4. Jangan menyimpan token JWT dalam cookies atau localStorage pada aplikasi client


pahami itu semua