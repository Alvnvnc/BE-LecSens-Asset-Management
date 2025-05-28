# Duplicate Asset Document Cleanup Guide

This guide provides detailed information about the duplicate asset document cleanup functionality in the LecSens Backend system.

## Overview

The duplicate cleanup feature helps maintain data integrity by identifying and removing duplicate asset documents while preserving the most recent version of each document type per asset. This is crucial for:

- **Storage Optimization**: Reducing unnecessary file storage usage
- **Data Integrity**: Ensuring each asset has only the latest version of each document type
- **Performance**: Improving query performance by reducing redundant data
- **Maintenance**: Keeping the database clean and organized

## How Duplicate Detection Works

### Detection Logic

The system identifies duplicates using the following criteria:

1. **Asset Grouping**: Documents are grouped by `asset_id`
2. **Type Grouping**: Within each asset, documents are further grouped by `document_type`
3. **Chronological Ordering**: Documents are ordered by `created_at` timestamp (DESC)
4. **Duplicate Identification**: Any document that is not the newest in its type group is considered a duplicate

### SQL Logic

The underlying SQL query uses window functions to rank documents:

```sql
WITH RankedDocuments AS (
    SELECT *,
           ROW_NUMBER() OVER(PARTITION BY asset_id, document_type ORDER BY created_at DESC) as row_num
    FROM asset_documents
    WHERE deleted_at IS NULL
)
SELECT * FROM RankedDocuments WHERE row_num > 1
```

## Usage Scenarios

### 1. Regular Maintenance

**Scenario**: Weekly cleanup as part of database maintenance routine.

```bash
# Step 1: Generate report to review duplicates
go run helpers/cmd/cmd.go -action=cleanup-duplicates -dry-run

# Step 2: Review the output and decide if cleanup is needed
# Step 3: Execute cleanup
go run helpers/cmd/cmd.go -action=cleanup-duplicates
```

### 2. Asset-Specific Cleanup

**Scenario**: Clean up duplicates for a specific asset after bulk document upload.

```bash
# Check duplicates for specific asset
go run helpers/cmd/cmd.go -action=cleanup-duplicates -asset-id=abc123e4-e89b-12d3-a456-426614174000 -dry-run

# Clean up that specific asset
go run helpers/cmd/cmd.go -action=cleanup-duplicates -asset-id=abc123e4-e89b-12d3-a456-426614174000
```

go run helpers/cmd/cmd.go helpers/cmd/cleanup_duplicates.go -action=cleanup-duplicates -dry-run
go run helpers/cmd/cmd.go helpers/cmd/cleanup_duplicates.go -action=cleanup-duplicates

### 3. Pre-deployment Cleanup

**Scenario**: Clean up duplicates before deploying to production.

```bash
# Automated check in deployment script
go run helpers/cmd/cmd.go -action=cleanup-duplicates -dry-run > duplicate_report.txt
if [ -s duplicate_report.txt ]; then
    echo "Duplicates found. Manual review required."
    exit 1
fi
```

## Understanding the Output

### Dry Run Report Format

When you run with `-dry-run`, you'll see output like this:

```
🔍 Scanning for duplicate asset documents...
📊 Found duplicates in 3 assets:

================================================================================
Asset ID: 123e4567-e89b-12d3-a456-426614174000
  📄 Total documents: 4
  🗑️  Duplicates to remove: 3
  🏷️  Document types: 
    - manual_v2.pdf [MANUAL] (2025-05-28 10:30:00) - KEEP
    - manual_v1.pdf [MANUAL] (2025-05-28 09:15:00) - DELETE
    - manual_old.pdf [MANUAL] (2025-05-28 08:00:00) - DELETE
    - warranty_new.pdf [WARRANTY] (2025-05-28 11:00:00) - KEEP

Asset ID: 456e7890-e89b-12d3-a456-426614174001
  📄 Total documents: 3
  🗑️  Duplicates to remove: 2
  🏷️  Document types: 
    - spec_sheet_v3.pdf [SPECIFICATION] (2025-05-28 14:30:00) - KEEP
    - spec_sheet_v2.pdf [SPECIFICATION] (2025-05-28 13:15:00) - DELETE
    - spec_sheet_v1.pdf [SPECIFICATION] (2025-05-28 12:00:00) - DELETE

================================================================================
📈 SUMMARY:
  🏢 Assets with duplicates: 2
  📄 Total duplicate documents to remove: 5
```

### Output Explanation

- **KEEP**: The newest document of each type that will be preserved
- **DELETE**: Older documents that will be removed
- **Asset ID**: Unique identifier for each asset
- **Document Type**: Category of the document (MANUAL, WARRANTY, SPECIFICATION, etc.)
- **Timestamp**: When the document was uploaded (creation time)

### Success Messages

After successful cleanup:

```
✅ Successfully cleaned up 5 duplicate documents!
```

For asset-specific cleanup:

```
✅ Successfully cleaned up 3 duplicate documents for asset 123e4567-e89b-12d3-a456-426614174000!
```

### No Duplicates Found

When no duplicates exist:

```
✅ No duplicate documents found!
```

## Safety Measures

### 1. Dry Run Mode

Always use dry run mode first to understand what will be changed:

```bash
go run helpers/cmd/cmd.go -action=cleanup-duplicates -dry-run
```

### 2. Confirmation Prompts

The tool will ask for confirmation before performing destructive operations:

```
❓ Do you want to proceed with deletion? (yes/no):
```

Type `yes` or `y` to proceed, anything else to cancel.

### 3. Preservation Logic

The cleanup process ensures:
- **Most Recent Document**: Always keeps the newest document based on `created_at` timestamp
- **Different Types**: Documents of different types are never considered duplicates
- **Soft Deletes**: Uses soft delete mechanism (sets `deleted_at` timestamp)
- **Cloud Files**: Does not delete actual files from cloud storage (Cloudinary)

### 4. Database Transactions

All cleanup operations are performed within database transactions to ensure:
- **Atomicity**: Either all duplicates are cleaned or none
- **Consistency**: Database remains in a valid state
- **Rollback**: Automatic rollback on any error

## Advanced Usage

### Scripting and Automation

Create a shell script for regular maintenance:

```bash
#!/bin/bash
# cleanup_routine.sh

echo "Starting duplicate cleanup routine..."

# Generate report
echo "Generating duplicate report..."
go run helpers/cmd/cmd.go -action=cleanup-duplicates -dry-run > /tmp/duplicate_report.txt

# Check if duplicates exist
if grep -q "📄 Total duplicate documents to remove:" /tmp/duplicate_report.txt; then
    echo "Duplicates found. Proceeding with cleanup..."
    go run helpers/cmd/cmd.go -action=cleanup-duplicates
else
    echo "No duplicates found. Nothing to clean up."
fi

echo "Cleanup routine completed."
```

### Integration with Monitoring

Monitor cleanup operations:

```bash
# Log cleanup results
go run helpers/cmd/cmd.go -action=cleanup-duplicates 2>&1 | tee cleanup_$(date +%Y%m%d_%H%M%S).log
```

### Batch Processing

For multiple specific assets:

```bash
#!/bin/bash
# batch_cleanup.sh

ASSET_IDS=(
    "123e4567-e89b-12d3-a456-426614174000"
    "456e7890-e89b-12d3-a456-426614174001"
    "789e1234-e89b-12d3-a456-426614174002"
)

for asset_id in "${ASSET_IDS[@]}"; do
    echo "Cleaning up asset: $asset_id"
    go run helpers/cmd/cmd.go -action=cleanup-duplicates -asset-id="$asset_id" -dry-run
    read -p "Proceed with cleanup for $asset_id? (y/n): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        go run helpers/cmd/cmd.go -action=cleanup-duplicates -asset-id="$asset_id"
    fi
done
```

## Troubleshooting

### Common Issues

#### 1. Permission Denied
```
Error: failed to cleanup duplicate documents: permission denied
```
**Solution**: Ensure database user has DELETE permissions on asset_documents table.

#### 2. Invalid UUID Format
```
Error: invalid asset ID format: invalid UUID length
```
**Solution**: Verify the asset ID is a valid UUID format (e.g., `123e4567-e89b-12d3-a456-426614174000`).

#### 3. Database Connection Issues
```
Error: failed to connect to database: dial tcp [::1]:5432: connect: connection refused
```
**Solution**: 
- Check if PostgreSQL is running
- Verify database connection parameters
- Ensure `.env` file is properly configured

#### 4. No Assets Found
```
✅ No duplicate documents found for asset 123e4567-e89b-12d3-a456-426614174000
```
**Solution**: This is normal if the asset has no duplicates. Verify the asset ID exists and has documents.

### Debugging Tips

1. **Enable Verbose Logging**: Check application logs for detailed operation information
2. **Database Query**: Manually run the duplicate detection query to verify results
3. **File System**: Ensure adequate disk space for log files
4. **Network**: Verify database network connectivity

### Recovery Procedures

If cleanup was performed incorrectly:

1. **Soft Delete Recovery**: Documents are soft-deleted, so they can be recovered by setting `deleted_at` to NULL
2. **Database Backup**: Restore from backup if hard delete was performed
3. **Cloud Storage**: Files in cloud storage (Cloudinary) are preserved and can be re-linked

## Performance Considerations

### Large Datasets

For databases with many asset documents:

1. **Run During Off-Peak Hours**: Cleanup operations can be I/O intensive
2. **Monitor Database Performance**: Watch for lock contention
3. **Batch Processing**: Consider cleaning specific assets rather than all at once
4. **Index Optimization**: Ensure proper indexes on `asset_id`, `document_type`, and `created_at`

### Resource Usage

- **Memory**: Cleanup operations load document metadata into memory
- **CPU**: Minimal CPU usage for comparison operations
- **I/O**: Database read/write operations for duplicate identification and deletion
- **Network**: Minimal network usage (no cloud file operations)

## Best Practices

### 1. Regular Scheduling

Set up regular cleanup schedules:

```bash
# Add to crontab for weekly cleanup
0 2 * * 0 cd /path/to/project && go run helpers/cmd/cmd.go -action=cleanup-duplicates >/dev/null 2>&1
```

### 2. Monitoring and Alerting

Monitor cleanup results:

```bash
# Send email notification with results
go run helpers/cmd/cmd.go -action=cleanup-duplicates -dry-run | mail -s "Weekly Duplicate Report" admin@example.com
```

### 3. Documentation

Keep records of cleanup operations:

```bash
# Log cleanup activities
echo "$(date): Cleanup completed" >> /var/log/asset_cleanup.log
```

### 4. Testing

Always test in development environment:

```bash
# Development environment test
DB_NAME=lecsens_dev go run helpers/cmd/cmd.go -action=cleanup-duplicates -dry-run
```

## API Integration

The cleanup functionality is also available through REST API endpoints:

- `GET /api/v1/superadmin/asset-documents/duplicates` - Get duplicates report
- `POST /api/v1/superadmin/asset-documents/cleanup-all` - Cleanup all duplicates
- `POST /api/v1/superadmin/asset-documents/cleanup/:assetId` - Cleanup specific asset

These endpoints require SuperAdmin authentication and provide the same functionality as the command line tool.

---

**Last Updated**: May 28, 2025  
**Version**: 1.0  
**Related Documentation**: [Command Line Tools](Command_Line_Tools.md), [External API](External_API.md)
