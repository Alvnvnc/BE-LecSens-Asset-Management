# Using External Tenant API

This asset management service is part of a microservice architecture where tenant management is handled by a separate service.

## Tenant Data Flow

1. All tenant information is retrieved from the external tenant API:
   - URL: https://lecsens-iot-api.erplabiim.com/api/external/tenants

2. Tenant identification in requests:
   - All API requests must include the tenant ID in the `X-Tenant-ID` header
   - Alternatively, the tenant ID can be provided via JWT token (preferred method)
   - The tenant ID is validated against the external tenant API
   - JWT authentication is implemented for better security

3. Tenant context in code:
   - The tenant ID is stored in the request context
   - All database operations automatically filter by the tenant ID
   - This ensures proper multi-tenancy data isolation

## API Endpoints

### Public Endpoints (No Tenant Required)
- `GET /api/v1/health` - Health check
- `GET /api/v1/tenants` - List all available tenants

### Protected Endpoints (Tenant Required)
- `GET /api/v1/tenant` - Get current tenant information
- `GET /api/v1/assets` - List assets for current tenant
- `GET /api/v1/assets/:id` - Get asset by ID for current tenant
- `POST /api/v1/assets` - Create new asset for current tenant

## Error Handling

The service implements standardized error handling with appropriate HTTP status codes:

| Error                     | HTTP Status Code | Description                                       |
|---------------------------|-----------------|---------------------------------------------------|
| Tenant not found          | 404             | The specified tenant does not exist               |
| Tenant inactive           | 403             | The tenant exists but is not active               |
| Subscription expired      | 403             | The tenant's subscription has expired             |
| API connection failed     | 503             | Unable to connect to the tenant API               |
| API response invalid      | 503             | Received an invalid response from the tenant API  |
| Unauthorized              | 401             | Missing or invalid authentication credentials     |
| Rate limit exceeded       | 429             | Too many requests to the tenant API               |

## Rate Limiting

The service implements rate limiting for external API calls to prevent overloading the tenant API:

- Default rate: 10 requests per second
- Burst limit: 20 requests
- Rate limits are applied per endpoint and per tenant
- When rate limit is exceeded, a 429 Too Many Requests status is returned

## Caching

To minimize external API calls, tenant data is cached:

- Cache TTL: 5 minutes
- In-memory cache implementation
- Cache is refreshed automatically when TTL expires

## Tenant DTO Structure

```go
type TenantDTO struct {
    ID                    uuid.UUID  `json:"id"`
    Name                  string     `json:"name"`
    Description           string     `json:"description,omitempty"`
    LogoURL               string     `json:"logo_url,omitempty"`
    ContactEmail          string     `json:"contact_email"`
    ContactPhone          string     `json:"contact_phone,omitempty"`
    MaxUsers              int        `json:"max_users"`
    SubscriptionPlan      string     `json:"subscription_plan"`
    SubscriptionStartDate time.Time  `json:"subscription_start_date"`
    SubscriptionEndDate   time.Time  `json:"subscription_end_date"`
    IsActive              bool       `json:"is_active"`
    CreatedAt             time.Time  `json:"created_at"`
    UpdatedAt             *time.Time `json:"updated_at,omitempty"`
}
```
