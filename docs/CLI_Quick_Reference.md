# Quick Reference - Command Line Tools

## Basic Commands

### Database Management
```bash
# Run migrations
go run helpers/cmd/cmd.go -action=migrate

# Drop specific table
go run helpers/cmd/cmd.go -action=drop-table -table=assets

# Truncate table (empty data, keep structure)
go run helpers/cmd/cmd.go -action=truncate-table -table=locations

# Drop all tables (DANGEROUS!)
go run helpers/cmd/cmd.go -action=drop-all -force
```

### Data Operations
```bash
# Seed location data
go run helpers/cmd/cmd.go -action=seed

# Custom CSV file
go run helpers/cmd/cmd.go -action=seed -csv=path/to/file.csv
```

### Duplicate Cleanup
```bash
# Preview all duplicates (safe)
go run helpers/cmd/cmd.go -action=cleanup-duplicates -dry-run

# Clean all duplicates
go run helpers/cmd/cmd.go -action=cleanup-duplicates

# Preview specific asset
go run helpers/cmd/cmd.go -action=cleanup-duplicates -asset-id=UUID -dry-run

# Clean specific asset
go run helpers/cmd/cmd.go -action=cleanup-duplicates -asset-id=UUID
```

## Quick Setup

### Environment Variables (.env)
```bash
DB_HOST=localhost
DB_PORT=5432
DB_USER=root
DB_PASSWORD=P@ssw0rd
DB_NAME=lecsens
ENVIRONMENT=development
```

### First Time Setup
```bash
# 1. Create database and run migrations
go run helpers/cmd/cmd.go -action=migrate

# 2. Seed initial data
go run helpers/cmd/cmd.go -action=seed

# 3. Check for any issues
go run helpers/cmd/cmd.go -action=cleanup-duplicates -dry-run
```

## Common Use Cases

### Development Reset
```bash
go run helpers/cmd/cmd.go -action=drop-all -force
go run helpers/cmd/cmd.go -action=migrate
go run helpers/cmd/cmd.go -action=seed
```

### Weekly Maintenance
```bash
go run helpers/cmd/cmd.go -action=cleanup-duplicates -dry-run
go run helpers/cmd/cmd.go -action=cleanup-duplicates
```

### Emergency Table Reset
```bash
go run helpers/cmd/cmd.go -action=truncate-table -table=asset_documents -force
```

## Safety Options

| Option | Description |
|--------|-------------|
| `-dry-run` | Preview without changes (cleanup only) |
| `-force` | Skip confirmations |

## Exit Codes

- `0` - Success
- `1` - Error (check logs for details)

## Need Help?

```bash
# Show help
go run helpers/cmd/cmd.go

# Or just run without parameters
go run helpers/cmd/cmd.go -action=
```

---

**📚 Full Documentation**: [Command_Line_Tools.md](Command_Line_Tools.md)  
**🧹 Cleanup Guide**: [Duplicate_Cleanup_Guide.md](Duplicate_Cleanup_Guide.md)
