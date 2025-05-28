# User-Tenant Context API Documentation

## Overview
The User-Tenant Context API system provides comprehensive management of user-tenant relationships in a multi-tenant environment. This API solves the critical problem where JWT tokens don't contain tenant information, requiring a separate system to manage which tenant a user is currently operating in and validate their access to specific tenants.

## Problem Statement
- JWT tokens contain user authentication but lack tenant context
- Need to track which tenant a user is currently working in
- Require validation of user access to specific tenants
- Support role-based access control within tenant boundaries
- Enable tenant switching functionality for users with multi-tenant access

## Architecture Components

### 1. Data Transfer Objects (DTOs)
**File**: `/helpers/dto/user_tenant.go`
- `UserTenantContextRequest/Response`: Current tenant operations
- `UserTenantsListResponse`: User's accessible tenants
- `TenantUsersResponse`: Admin tenant user management
- `UserTenantAccessValidationRequest/Response`: Access validation
- `SwitchTenantRequest/Response`: Tenant switching operations

### 2. Repository Layer
**File**: `/data-layer/repository/user_tenant_repository.go`
- Database operations for user-tenant relationships
- CRUD operations with proper error handling
- Optimized queries for tenant context retrieval

### 3. Service Layer
**File**: `/domain-layer/service/user_tenant_context_service.go`
- Business logic for tenant management
- Access validation and authorization
- Tenant switching workflow

### 4. Controller Layer
**File**: `/presentation-layer/controller/user_tenant_context_controller.go`
- HTTP request handling
- JWT authentication integration
- Role-based authorization enforcement

### 5. Routes Configuration
**File**: `/presentation-layer/routes/user_tenant_context.go`
- Protected route definitions
- Middleware integration
- Role-based access control

## API Endpoints Documentation

### 1. Get Current Tenant
**Route**: `GET /api/v1/user-tenant/current`
**Authentication**: Required (JWT)
**Authorization**: All authenticated users

**Purpose**: Retrieves the current tenant context for the authenticated user.

**Functionality**:
- Extracts user ID from JWT token context
- Queries user's current active tenant
- Returns tenant information and user role within that tenant

**Response Data**:
```json
{
  "message": "Current tenant retrieved successfully",
  "data": {
    "user_id": "uuid",
    "tenant_id": "uuid",
    "tenant_name": "string",
    "user_role_in_tenant": "string",
    "is_active": true,
    "joined_at": "timestamp"
  }
}
```

**Use Cases**:
- Dashboard initialization to show current tenant context
- Navigation breadcrumb display
- Permission checking for tenant-specific features

---

### 2. Set Current Tenant
**Route**: `PUT /api/v1/user-tenant/current`
**Authentication**: Required (JWT)
**Authorization**: All authenticated users
**Method**: PUT
**Content-Type**: application/json

**Purpose**: Sets or updates the current tenant context for the authenticated user.

**Request Body**:
```json
{
  "tenant_id": "uuid"
}
```

**Functionality**:
- Validates user has access to the specified tenant
- Updates user's current tenant context
- Returns updated tenant information

**Use Cases**:
- Initial tenant selection after login
- Programmatic tenant switching
- Tenant context restoration

---

### 3. Get User Tenants
**Route**: `GET /api/v1/user-tenant/tenants`
**Authentication**: Required (JWT)
**Authorization**: All authenticated users

**Purpose**: Retrieves all tenants that the authenticated user has access to.

**Functionality**:
- Lists all tenant memberships for the user
- Includes role information for each tenant
- Shows active/inactive status

**Response Data**:
```json
{
  "message": "User tenants retrieved successfully",
  "data": {
    "tenants": [
      {
        "tenant_id": "uuid",
        "tenant_name": "string",
        "user_role": "string",
        "is_current": true,
        "is_active": true,
        "joined_at": "timestamp"
      }
    ],
    "total_count": 5
  }
}
```

**Use Cases**:
- Tenant selector dropdown in UI
- Multi-tenant dashboard overview
- Access audit trails

---

### 4. Switch Tenant
**Route**: `POST /api/v1/user-tenant/switch`
**Authentication**: Required (JWT)
**Authorization**: All authenticated users
**Method**: POST
**Content-Type**: application/json

**Purpose**: Switches the user's current active tenant context.

**Request Body**:
```json
{
  "tenant_id": "uuid"
}
```

**Functionality**:
- Validates user access to target tenant
- Updates current tenant context
- Logs tenant switch activity
- Returns new tenant context

**Response Data**:
```json
{
  "message": "Tenant switched successfully",
  "data": {
    "previous_tenant_id": "uuid",
    "new_tenant_id": "uuid",
    "new_tenant_name": "string",
    "user_role": "string",
    "switched_at": "timestamp"
  }
}
```

**Use Cases**:
- User-initiated tenant switching via UI
- Automated tenant switching based on URL routing
- Context switching in multi-tenant workflows

---

### 5. Validate User Tenant Access
**Route**: `POST /api/v1/user-tenant/validate-access`
**Authentication**: Required (JWT)
**Authorization**: All authenticated users
**Method**: POST
**Content-Type**: application/json

**Purpose**: Validates if the authenticated user has access to a specific tenant.

**Request Body**:
```json
{
  "tenant_id": "uuid"
}
```

**Functionality**:
- Checks user membership in specified tenant
- Validates active status of the membership
- Returns access permissions and role information

**Response Data**:
```json
{
  "message": "Access validation completed",
  "data": {
    "user_id": "uuid",
    "tenant_id": "uuid",
    "has_access": true,
    "user_role": "string",
    "is_active": true,
    "permissions": ["read", "write", "admin"]
  }
}
```

**Use Cases**:
- Pre-flight checks before tenant-specific operations
- API gateway authorization decisions
- Dynamic permission enforcement

---

### 6. Get Tenant Users (Admin Only)
**Route**: `GET /api/v1/user-tenant/users`
**Authentication**: Required (JWT)
**Authorization**: Admin, SuperAdmin only
**Query Parameters**: 
- `page` (optional): Page number, default 1
- `limit` (optional): Items per page, default 10, max 100

**Purpose**: Retrieves all users within the admin's current tenant.

**Functionality**:
- Validates admin role permissions
- Gets admin's current tenant context
- Returns paginated list of users in that tenant
- Includes user roles and status information

**Response Data**:
```json
{
  "message": "Tenant users retrieved successfully",
  "data": {
    "users": [
      {
        "user_id": "uuid",
        "username": "string",
        "email": "string",
        "role_in_tenant": "string",
        "is_active": true,
        "joined_at": "timestamp",
        "last_active": "timestamp"
      }
    ],
    "pagination": {
      "current_page": 1,
      "total_pages": 5,
      "total_users": 45,
      "per_page": 10
    },
    "tenant_info": {
      "tenant_id": "uuid",
      "tenant_name": "string"
    }
  }
}
```

**Use Cases**:
- Admin user management dashboard
- Tenant user audit and reporting
- Role assignment and management
- User activity monitoring within tenant

---

### 7. Get Users by Tenant ID (SuperAdmin Only)
**Route**: `GET /api/v1/user-tenant/tenants/{tenantId}/users`
**Authentication**: Required (JWT)
**Authorization**: SuperAdmin only
**Path Parameters**:
- `tenantId`: UUID of the target tenant
**Query Parameters**: 
- `page` (optional): Page number, default 1
- `limit` (optional): Items per page, default 10, max 100

**Purpose**: Allows SuperAdmins to view users in any specific tenant across the system.

**Functionality**:
- Validates SuperAdmin role permissions
- Retrieves users for specified tenant ID
- Returns comprehensive user information
- Supports cross-tenant user management

**Response Data**:
```json
{
  "message": "Tenant users retrieved successfully",
  "data": {
    "users": [
      {
        "user_id": "uuid",
        "username": "string",
        "email": "string",
        "role_in_tenant": "string",
        "is_active": true,
        "joined_at": "timestamp",
        "last_active": "timestamp"
      }
    ],
    "pagination": {
      "current_page": 1,
      "total_pages": 3,
      "total_users": 28,
      "per_page": 10
    },
    "tenant_info": {
      "tenant_id": "uuid",
      "tenant_name": "string"
    }
  }
}
```

**Use Cases**:
- System-wide tenant administration
- Cross-tenant user analytics
- Global user management operations
- Tenant health monitoring

## Authentication & Authorization Flow

### JWT Integration
1. **Token Extraction**: All endpoints extract user information from JWT via middleware
2. **Context Population**: User ID and role are stored in Gin context (`ctx.Get("userID")`, `ctx.Get("userRole")`)
3. **Permission Validation**: Role-based checks before accessing protected endpoints

### Role Hierarchy
- **User**: Access to own tenant context and switching
- **Admin**: Access to manage users within their current tenant
- **SuperAdmin**: Access to manage users across all tenants

## Error Handling

### Common Error Responses
- `400 Bad Request`: Invalid request format or parameters
- `401 Unauthorized`: Missing or invalid JWT token
- `403 Forbidden`: Insufficient permissions for the operation
- `404 Not Found`: Requested resource not found
- `500 Internal Server Error`: Server-side processing errors

### Error Response Format
```json
{
  "error": "Error description",
  "details": "Additional error details"
}
```

## Integration Points

### Middleware Dependencies
- **Authentication Middleware**: JWT token validation
- **Role Middleware**: Role extraction and validation
- **Audit Middleware**: Activity logging

### Database Tables
- `user_tenants`: User-tenant relationship mapping
- `tenants`: Tenant information
- `users`: User account data
- `roles`: Role definitions

## Usage Scenarios

### 1. User Login Flow
1. User authenticates and receives JWT
2. Call `GET /api/v1/user-tenant/tenants` to get available tenants
3. If multiple tenants, show tenant selector
4. Call `PUT /api/v1/user-tenant/current` to set active tenant
5. Proceed with tenant-specific operations

### 2. Tenant Switching
1. User selects different tenant from dropdown
2. Call `POST /api/v1/user-tenant/switch` with new tenant ID
3. Update UI context and navigation
4. Refresh tenant-specific data

### 3. Admin User Management
1. Admin accesses user management section
2. Call `GET /api/v1/user-tenant/users` to list tenant users
3. Display user list with roles and status
4. Enable user role modification and status updates

### 4. Access Validation
1. Before sensitive operations, call `POST /api/v1/user-tenant/validate-access`
2. Check returned permissions
3. Allow or deny operation based on validation result

## Performance Considerations

### Caching Strategy
- Cache user tenant lists for short periods
- Cache current tenant context in user session
- Implement Redis caching for frequently accessed data

### Pagination
- Default page size of 10 items
- Maximum page size of 100 items
- Efficient database queries with proper indexing

### Database Optimization
- Indexed foreign keys on user_id and tenant_id
- Composite indexes for common query patterns
- Connection pooling for concurrent requests

## Security Features

### Input Validation
- UUID format validation for all ID parameters
- Request body schema validation
- Query parameter sanitization

### Access Control
- Role-based endpoint protection
- Tenant boundary enforcement
- User session validation

### Audit Logging
- All tenant switch operations logged
- Admin actions on users recorded
- Failed access attempts tracked

## Future Enhancements

### Planned Features
1. Bulk user operations for admins
2. Tenant usage analytics
3. User invitation system
4. Advanced role permissions
5. Tenant-specific configuration management

### API Versioning
- Current version: v1
- Backward compatibility maintenance
- Deprecation notices for old endpoints

This documentation provides a comprehensive guide for understanding and implementing the User-Tenant Context API system, enabling effective management of users within multi-tenant environments.
