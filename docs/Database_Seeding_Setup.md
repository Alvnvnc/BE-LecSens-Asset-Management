# Database Seeding Setup

This document describes the automatic database seeding setup for the Asset Management System.

## Overview

The application includes an automatic seeding system that populates the database with initial data during the build/startup process. This is particularly useful for development environments and initial deployments.

## Seeding Components

### 1. Auto-Seed Script (`scripts/auto-seed.sh`)
- Automatically runs database migrations
- Seeds data in the correct order (respecting foreign key constraints)
- Configurable via environment variables
- Integrated into Docker startup process

### 2. Command Tool (`helpers/cmd/cmd.go`)
- Flexible CLI tool for database operations
- Supports various seeding options
- Can be used independently for manual seeding

### 3. Seeder Modules
Located in `data-layer/migration/seeder/`:
- `location` - Geographic locations
- `asset-type` - Asset categories and types
- `sensor-type` - Sensor hardware types
- `measurement-type` - Measurement categories
- `measurement-field` - Specific measurement fields
- `asset` - Physical assets
- `asset-sensor` - Sensor assignments to assets
- `threshold` - Alert thresholds
- `sensor-status` - Current sensor status
- `reading` - Historical sensor readings
- `sensor-logs` - System logs
- `alert` - Alert records

## Configuration

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `AUTO_SEED` | `true` | Enable/disable automatic seeding |
| `ENVIRONMENT` | `development` | Environment type (development enables auto-seeding) |
| `DB_HOST` | `db` | Database host |
| `DB_PORT` | `5432` | Database port |
| `DB_USER` | `root` | Database user |
| `DB_PASSWORD` | `P@ssw0rd` | Database password |
| `DB_NAME` | `asset_management` | Database name |

### Docker Compose Configuration

Auto-seeding is enabled by default in `docker-compose.yaml`:

```yaml
environment:
  ENVIRONMENT: development
  AUTO_SEED: true
```

To disable auto-seeding:
```yaml
environment:
  ENVIRONMENT: production
  AUTO_SEED: false
```

## Usage

### Automatic Seeding (Docker)

1. **Development Environment** (with auto-seeding):
   ```bash
   make docker-run
   # or
   docker-compose up --build
   ```

2. **Production Environment** (without auto-seeding):
   ```bash
   make prod-build
   # or
   ENVIRONMENT=production AUTO_SEED=false docker-compose up --build
   ```

### Manual Seeding

1. **Complete Seeding**:
   ```bash
   make seed-all
   ```

2. **Specific Seeder**:
   ```bash
   go run ./helpers/cmd/cmd.go -action=seed -seeder=location -force
   ```

3. **Custom Parameters**:
   ```bash
   # Seed readings for last 30 days
   go run ./helpers/cmd/cmd.go -action=seed -seeder=reading -days=30 -force
   
   # Generate 500 sensor readings
   go run ./helpers/cmd/cmd.go -action=seed -seeder=reading -count=500 -force
   ```

### Available Commands

| Command | Description |
|---------|-------------|
| `make build` | Build the application |
| `make docker-run` | Run with Docker (auto-seed enabled) |
| `make docker-stop` | Stop Docker services |
| `make docker-clean` | Clean Docker containers and volumes |
| `make migrate` | Run database migrations only |
| `make seed` | Run all seeders |
| `make seed-all` | Run seeders in correct order |
| `make dev-setup` | Complete development setup |
| `make logs` | View application logs |
| `make db-logs` | View database logs |

## Seeding Order

The auto-seed script follows this order to respect foreign key constraints:

1. **locations** - Base geographic data
2. **asset-type** - Asset categories
3. **sensor-type** - Sensor hardware types
4. **measurement-type** - Measurement categories
5. **measurement-field** - Specific fields
6. **asset** - Physical assets (depends on locations, asset-types)
7. **asset-sensor** - Sensor assignments (depends on assets, sensor-types)
8. **threshold** - Alert thresholds (depends on asset-sensors)
9. **sensor-status** - Current status (depends on asset-sensors)
10. **reading** - Historical data (depends on asset-sensors)
11. **sensor-logs** - System logs (depends on asset-sensors)
12. **alert** - Alert records (depends on readings, thresholds)

## Troubleshooting

### Common Issues

1. **Database Connection Failed**:
   - Check database is running: `docker-compose logs db`
   - Verify connection parameters in environment variables

2. **Seeding Failed**:
   - Check application logs: `make logs`
   - Verify foreign key constraints are satisfied
   - Ensure proper seeding order

3. **Permission Denied**:
   - Make script executable: `chmod +x scripts/auto-seed.sh`
   - Check Docker container permissions

### Manual Recovery

If auto-seeding fails, you can run manual operations:

```bash
# Reset database
go run ./helpers/cmd/cmd.go -action=drop-all -force
go run ./helpers/cmd/cmd.go -action=migrate -force

# Manual seeding
make seed-all
```

### Logs and Debugging

View detailed logs during startup:
```bash
docker-compose logs -f app
```

Enable debug mode by setting environment variables:
```yaml
environment:
  GIN_MODE: debug
  JWT_DEBUG: true
```

## Production Considerations

1. **Disable Auto-Seeding**: Set `AUTO_SEED=false` for production
2. **Custom Data**: Replace CSV files in `data-layer/migration/seeder/` with production data
3. **Security**: Update default passwords and API keys
4. **Performance**: Consider seeding large datasets offline
5. **Backup**: Always backup production data before running seeders

## Development Workflow

1. **Initial Setup**:
   ```bash
   make dev-setup
   ```

2. **Reset Database**:
   ```bash
   make docker-clean
   make docker-run
   ```

3. **Add New Seeder**:
   - Create seeder in `data-layer/migration/seeder/`
   - Add to `cmd.go` switch statement
   - Update `auto-seed.sh` script
   - Test with `make seed-all`
