# Asset Management Service - User Management Integration

## Overview

This document outlines the integration between the Asset Management service and the User Management service. The integration enables the Asset Management service to:

1. Validate JWT tokens
2. Verify user permissions
3. Check tenant access and subscription status
4. Apply business rules based on tenant subscription limits

## Authentication & Authorization Flow

1. Client sends request to Asset Management service with JWT token
2. Asset Management validates token with User Management service
3. User Management returns token validity and user information
4. Asset Management checks tenant access with User Management service
5. If validation succeeds, Asset Management processes the request
6. Business rules are applied based on tenant subscription limits

## Integration Components

### 1. User Service

The `UserService` handles interactions with the User Management API:

- **ValidateToken**: Validates JWT tokens against the User Management service
- **ValidateUserPermissions**: Checks if a user has specific permissions
- **ValidateUserAccess**: Checks if a user has access to a tenant
- **GetUser**: Retrieves user information
- **GetUserTenants**: Lists all tenants a user has access to

### 2. Tenant Service

The `TenantService` handles interactions with the Tenant Management API:

- **GetTenant**: Retrieves tenant information
- **ValidateTenantAccess**: Checks if tenant is active with valid subscription
- **GetTenantLimits**: Retrieves subscription-based limits for a tenant
- **GetTenantSubscription**: Gets subscription details for a tenant
- **ValidateUserTenantAccess**: Validates if a user has access to a tenant

### 3. Middleware

- **JWTMiddleware**: Validates JWT tokens against the User Management service
- **TenantMiddleware**: Verifies tenant existence and user access to tenant

## Configuration

Set the following environment variables to configure the integration:

```
USER_API_URL=https://user-management-service.com/api
USER_API_KEY=your-user-api-key
TENANT_API_URL=https://tenant-management-service.com/api
TENANT_API_KEY=your-tenant-api-key
```

## Error Handling

The integration handles various error scenarios:

- Invalid or expired tokens
- User not found or inactive
- Tenant not found or inactive
- User doesn't have access to tenant
- Tenant subscription expired
- Rate limiting exceeded

## Data Transfer Objects

The integration uses the following DTOs:

- **TokenValidationRequest/Response**: For token validation
- **PermissionValidationRequest/Response**: For permission checking
- **UserDTO**: User information
- **TenantDTO**: Tenant information
- **TenantLimitsResponse**: Subscription-based limits
- **TenantSubscriptionResponse**: Subscription details
- **UserTenantAccessResponse**: User-tenant access validation

## Implementation Examples

### Token Validation

```go
validationResponse, err := userService.ValidateToken(ctx, tokenString)
if err != nil || !validationResponse.Valid {
    // Handle invalid token
}
```

### Permission Checking

```go
permissionResponse, err := userService.ValidateUserPermissions(
    ctx, 
    tokenString, 
    tenantID, 
    "ADMIN", 
    []string{"asset:create", "asset:update"}
)
if err != nil || !permissionResponse.Valid {
    // Handle permission denied
}
```

### Tenant Limits

```go
limits, err := tenantService.GetTenantLimits(ctx, tenantID)
if err != nil {
    // Handle error
}

if currentAssetCount >= limits.Limits.MaxAssets {
    // Handle limit exceeded
}
```
