package repository

import (
	"be-lecsens/asset_management/helpers/common"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// BaseRepository provides common functionality for all repositories
// It automatically applies tenant filtering to all operations
type BaseRepository struct {
	DB *sql.DB
}

// NewBaseRepository creates a new BaseRepository
func NewBaseRepository(db *sql.DB) *BaseRepository {
	return &BaseRepository{
		DB: db,
	}
}

// ExecuteWithTenant executes a query with tenant ID from context
// This ensures that all data access is filtered by tenant
func (r *BaseRepository) ExecuteWithTenant(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	tenantID, ok := common.GetTenantID(ctx)
	if !ok {
		return nil, errors.New("tenant ID is required for this operation")
	}

	// Add tenant ID as the first argument
	newArgs := append([]interface{}{tenantID}, args...)

	// The query should have a WHERE tenant_id = $1 clause
	return r.DB.QueryContext(ctx, query, newArgs...)
}

// ExecuteWithOptionalTenant executes a query with tenant ID from context if available
// This is useful for operations that might be executed in a system context without a tenant
func (r *BaseRepository) ExecuteWithOptionalTenant(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	tenantID, ok := common.GetTenantID(ctx)
	if !ok {
		// No tenant filtering, proceed with original query
		return r.DB.QueryContext(ctx, query, args...)
	}

	// Add tenant filtering
	newArgs := append([]interface{}{tenantID}, args...)
	tenantFilteredQuery := fmt.Sprintf("%s AND tenant_id = $1", query)

	return r.DB.QueryContext(ctx, tenantFilteredQuery, newArgs...)
}
