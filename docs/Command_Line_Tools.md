# Command Line Tools Documentation

This document provides comprehensive documentation for the command line tools available in the LecSens Backend Asset Management system.

## Overview

The LecSens Backend provides a powerful command line tool (`cmd.go`) that allows developers and administrators to perform various database operations and maintenance tasks. The tool supports multiple actions including database migrations, table management, data seeding, and duplicate asset document cleanup.

## Prerequisites

Before using the command line tools, ensure you have:

1. **Go 1.22 or later** installed
2. **PostgreSQL database** running and accessible
3. **Environment variables** properly configured (see [Environment Setup](#environment-setup))
4. **Database connection** established

## Environment Setup

Create a `.env` file in the project root with the following variables:

```bash
# Database Configuration
DB_HOST=localhost          # or 'db' for Docker
DB_PORT=5432
DB_USER=root              # or your database user
DB_PASSWORD=P@ssw0rd      # your database password
DB_NAME=lecsens

# Application Configuration
ENVIRONMENT=development
PORT=3160

# JWT Configuration
JWT_SECRET_KEY=your-secret-key
JWT_ISSUER=lecsens
JWT_EXPIRES_IN=3600
JWT_DEBUG=true

# External API Configuration
TENANT_API_URL=http://localhost:3000
TENANT_API_KEY=your-tenant-api-key
USER_API_URL=http://localhost:3001
USER_API_KEY=your-user-api-key
USER_AUTH_VALIDATE_TOKEN_ENDPOINT=/auth/validate
USER_AUTH_USER_INFO_ENDPOINT=/auth/user
USER_AUTH_VALIDATE_PERMISSIONS_ENDPOINT=/auth/permissions
USER_AUTH_VALIDATE_SUPERADMIN_ENDPOINT=/auth/superadmin

# Cloudinary Configuration
CLOUDINARY_CLOUD_NAME=your-cloud-name
CLOUDINARY_API_KEY=your-api-key
CLOUDINARY_API_SECRET=your-api-secret
```

## Usage

### Basic Syntax

```bash
go run helpers/cmd/cmd.go -action=<action> [options]
```

### Available Actions

| Action | Description | Required Options | Optional Options |
|--------|-------------|------------------|------------------|
| `drop-table` | Drop a specific database table | `-table=<name>` | `-force` |
| `truncate-table` | Empty a specific table (keep structure) | `-table=<name>` | `-force` |
| `drop-all` | Drop all database tables | None | `-force` |
| `migrate` | Run database migrations | None | None |
| `seed` | Seed location data from CSV | None | `-csv=<path>` |
| `cleanup-duplicates` | Clean up duplicate asset documents | None | `-asset-id=<id>`, `-dry-run` |

### Global Options

| Option | Description | Type | Default |
|--------|-------------|------|---------|
| `-force` | Skip confirmation prompts | boolean | false |
| `-dry-run` | Preview changes without executing (cleanup only) | boolean | false |

## Detailed Action Documentation

### 1. Database Table Management

#### Drop Table
Permanently removes a specific table and all its data.

```bash
# Drop the assets table
go run helpers/cmd/cmd.go -action=drop-table -table=assets

# Drop with force (no confirmation)
go run helpers/cmd/cmd.go -action=drop-table -table=assets -force
```

**⚠️ Warning**: This action is irreversible. All data in the table will be lost.

#### Truncate Table
Removes all data from a table but keeps the table structure intact.

```bash
# Truncate the locations table
go run helpers/cmd/cmd.go -action=truncate-table -table=locations

# Truncate with force (no confirmation)
go run helpers/cmd/cmd.go -action=truncate-table -table=locations -force
```

#### Drop All Tables
Removes all tables from the database schema.

```bash
# Drop all tables (with confirmation)
go run helpers/cmd/cmd.go -action=drop-all

# Drop all tables without confirmation
go run helpers/cmd/cmd.go -action=drop-all -force
```

**⚠️ Warning**: This will completely wipe your database. Use with extreme caution.

### 2. Database Migrations

#### Run Migrations
Executes all pending database migrations to set up or update the database schema.

```bash
# Run all migrations
go run helpers/cmd/cmd.go -action=migrate
```

This command will:
- Create all necessary tables
- Set up indexes and constraints
- Apply any schema updates
- Initialize default data if required

### 3. Data Seeding

#### Location Seeder
Imports location data (cities and regencies) from a CSV file.

```bash
# Use default CSV file (data-layer/migration/seeder/kota_kab.csv)
go run helpers/cmd/cmd.go -action=seed

# Use custom CSV file
go run helpers/cmd/cmd.go -action=seed -csv=path/to/your/locations.csv
```

**CSV Format Requirements**:
The CSV file should contain location data with appropriate columns for city/regency information.

### 4. Duplicate Asset Document Cleanup

#### Overview
The cleanup functionality helps identify and remove duplicate asset documents while preserving the most recent version of each document type per asset.

#### View Duplicates Report
Generate a detailed report of all duplicate documents without making any changes.

```bash
# Show comprehensive duplicates report
go run helpers/cmd/cmd.go -action=cleanup-duplicates -dry-run
```

**Report includes**:
- Assets with duplicate documents
- Document details (filename, type, creation date)
- Which documents would be kept vs. deleted
- Summary statistics

#### Cleanup All Duplicates
Remove duplicate documents from all assets in the system.

```bash
# Preview what would be cleaned (dry run)
go run helpers/cmd/cmd.go -action=cleanup-duplicates -dry-run

# Perform actual cleanup (with confirmation prompt)
go run helpers/cmd/cmd.go -action=cleanup-duplicates
```

**Cleanup Logic**:
- Groups documents by asset ID and document type
- Keeps the newest document (latest `created_at`)
- Deletes all older duplicates
- Preserves document files in cloud storage

#### Cleanup Specific Asset
Remove duplicates for a particular asset only.

```bash
# Preview cleanup for specific asset
go run helpers/cmd/cmd.go -action=cleanup-duplicates -asset-id=123e4567-e89b-12d3-a456-426614174000 -dry-run

# Cleanup specific asset
go run helpers/cmd/cmd.go -action=cleanup-duplicates -asset-id=123e4567-e89b-12d3-a456-426614174000
```

## Safety Features

### Confirmation Prompts
Most destructive operations require user confirmation unless the `-force` flag is used:

```
Are you sure you want to drop table 'assets'? This action cannot be undone. (y/N):
```

### Dry Run Mode
The cleanup operations support dry run mode to preview changes:

```bash
# See what would be cleaned without actually doing it
go run helpers/cmd/cmd.go -action=cleanup-duplicates -dry-run
```

### Database Connection Validation
All commands verify database connectivity before proceeding:

```
Connected to database successfully
```

## Error Handling

### Common Errors and Solutions

#### Environment Variables Missing
```
Required environment variable not set: DB_HOST
```
**Solution**: Ensure all required environment variables are set in your `.env` file.

#### Database Connection Failed
```
Failed to connect to database: dial tcp [::1]:5432: connect: connection refused
```
**Solution**: 
- Verify database is running
- Check database host and port settings
- Ensure database credentials are correct

#### Table Not Found
```
Failed to drop table assets: relation "assets" does not exist
```
**Solution**: The table may have already been dropped or never existed.

#### CSV File Not Found
```
CSV file not found: data-layer/migration/seeder/kota_kab.csv
```
**Solution**: Ensure the CSV file exists at the specified path or provide a valid path with `-csv` option.

#### Invalid Asset ID Format
```
invalid asset ID format: invalid UUID length: 35
```
**Solution**: Provide a valid UUID format for the asset ID.

## Best Practices

### 1. Always Use Dry Run First
Before performing any cleanup operations, use dry run mode to understand what will be changed:

```bash
go run helpers/cmd/cmd.go -action=cleanup-duplicates -dry-run
```

### 2. Backup Before Destructive Operations
Create database backups before using:
- `drop-table`
- `drop-all`
- `cleanup-duplicates` (without dry run)

### 3. Test in Development Environment
Always test commands in a development environment before running in production.

### 4. Monitor Database Size
Regular cleanup of duplicates can help maintain optimal database performance:

```bash
# Weekly cleanup routine
go run helpers/cmd/cmd.go -action=cleanup-duplicates -dry-run  # Review
go run helpers/cmd/cmd.go -action=cleanup-duplicates           # Execute
```

### 5. Use Version Control for Migrations
Ensure all migration files are committed to version control before running migrations.

## Examples

### Complete Database Setup
```bash
# 1. Run migrations to create tables
go run helpers/cmd/cmd.go -action=migrate

# 2. Seed location data
go run helpers/cmd/cmd.go -action=seed

# 3. Check for any existing duplicates
go run helpers/cmd/cmd.go -action=cleanup-duplicates -dry-run
```

### Database Reset
```bash
# 1. Drop all existing tables
go run helpers/cmd/cmd.go -action=drop-all -force

# 2. Run fresh migrations
go run helpers/cmd/cmd.go -action=migrate

# 3. Seed initial data
go run helpers/cmd/cmd.go -action=seed
```

### Maintenance Routine
```bash
# Weekly duplicate cleanup
go run helpers/cmd/cmd.go -action=cleanup-duplicates -dry-run
go run helpers/cmd/cmd.go -action=cleanup-duplicates

# Monthly table maintenance (if needed)
go run helpers/cmd/cmd.go -action=truncate-table -table=temp_data -force
```

## Troubleshooting

### Check Database Connection
```bash
# Test connection with a simple migration check
go run helpers/cmd/cmd.go -action=migrate
```

### Verify Environment Variables
```bash
# Check if variables are loaded
echo $DB_HOST
echo $DB_PORT
echo $DB_NAME
```

### Check Log Output
The command line tool provides detailed logging for all operations. Pay attention to:
- Connection status messages
- Operation progress indicators
- Error details and suggestions
- Success confirmations

## Integration with Development Workflow

### Pre-commit Hooks
Consider adding duplicate cleanup to your development workflow:

```bash
# Check for duplicates before committing
go run helpers/cmd/cmd.go -action=cleanup-duplicates -dry-run
```

### CI/CD Pipeline
Integrate commands into your deployment pipeline:

```bash
# Database migration step
go run helpers/cmd/cmd.go -action=migrate

# Post-deployment cleanup
go run helpers/cmd/cmd.go -action=cleanup-duplicates -dry-run
```

## Support

For issues or questions regarding the command line tools:

1. Check the error messages for specific guidance
2. Verify environment configuration
3. Ensure database connectivity
4. Review this documentation for usage examples
5. Check application logs for additional context

---

**Last Updated**: May 28, 2025  
**Version**: 1.0  
**Compatibility**: Go 1.22+, PostgreSQL 15+
