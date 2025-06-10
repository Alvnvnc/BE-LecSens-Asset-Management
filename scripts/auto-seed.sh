#!/bin/bash

# Auto-seeding script for the Asset Management system
# This script runs database migration and seeding automatically

echo "Starting auto-seeding process..."

# Set environment variables
export ENVIRONMENT=${ENVIRONMENT:-development}
export AUTO_SEED=${AUTO_SEED:-false}
export FORCE_RESEED=${FORCE_RESEED:-false}
export DB_HOST=${DB_HOST:-db}
export DB_PORT=${DB_PORT:-5432}
export DB_USER=${DB_USER:-root}
export DB_PASSWORD=${DB_PASSWORD:-P@ssw0rd}
export DB_NAME=${DB_NAME:-asset_management}

# Wait for database to be ready
echo "Waiting for database to be ready..."
until pg_isready -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME; do
  echo "Database not ready, waiting..."
  sleep 3
done

echo "Database is ready!"

# Function to run seeding process
run_seeding() {
    echo "Starting seeding process..."
    
    # Seed in the correct order (respecting foreign key constraints)
    echo "Seeding locations..."
    ./cmd -action=seed -seeder=location
    sleep 2  # Wait for transaction to commit
    
    echo "Seeding asset types..."
    ./cmd -action=seed -seeder=asset-type
    sleep 2  # Wait for transaction to commit
    
    echo "Seeding sensor types..."
    ./cmd -action=seed -seeder=sensor-type
    sleep 2  # Wait for transaction to commit
    
    echo "Seeding measurement types..."
    ./cmd -action=seed -seeder=measurement-type
    sleep 2  # Wait for transaction to commit
    
    echo "Seeding measurement fields..."
    ./cmd -action=seed -seeder=measurement-field
    sleep 2  # Wait for transaction to commit
    
    echo "Seeding assets..."
    ./cmd -action=seed -seeder=asset
    sleep 3  # Wait longer for asset seeding to complete
    
    echo "Seeding asset sensors..."
    ./cmd -action=seed -seeder=asset-sensor
    sleep 5  # Wait longer for asset sensor seeding to complete
    
    echo "Seeding sensor thresholds..."
    ./cmd -action=seed -seeder=threshold
    sleep 2  # Wait for transaction to commit
    
    echo "Seeding sensor status..."
    ./cmd -action=seed -seeder=sensor-status
    sleep 2  # Wait for transaction to commit
    
    echo "Seeding sensor readings (from Jan 1, 2022 to present)..."
    ./cmd -action=seed -seeder=reading
    sleep 2  # Wait for transaction to commit
    
    echo "Seeding sensor logs..."
    ./cmd -action=seed -seeder=sensor-logs
    sleep 2  # Wait for transaction to commit
    
    echo "Seeding asset alerts..."
    ./cmd -action=seed -seeder=alert
    
    echo "All seeding completed successfully!"
}

# Run migrations
echo "Running migrations..."
./cmd -action=migrate -force

# Check if seeding should be performed
if [ "$AUTO_SEED" = "true" ] || [ "$ENVIRONMENT" = "development" ]; then
    echo "Checking if database already has data..."
    
    # Check if tables exist and have data
    TABLE_COUNT=$(PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -t -c "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public';" 2>/dev/null | tr -d ' ')
    
    if [ "$TABLE_COUNT" -gt 0 ] && [ "$FORCE_RESEED" != "true" ]; then
        # Check if location table has data (as an indicator)
        LOCATION_COUNT=$(PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -t -c "SELECT COUNT(*) FROM locations;" 2>/dev/null | tr -d ' ')
        
        if [ "$LOCATION_COUNT" -gt 0 ]; then
            echo "Database already contains data ($LOCATION_COUNT locations found). Skipping seeding..."
            echo "Set FORCE_RESEED=true to force re-seeding."
        else
            echo "Tables exist but no data found. Running seeding..."
            run_seeding
        fi
    else
        if [ "$FORCE_RESEED" = "true" ]; then
            echo "FORCE_RESEED=true. Running full re-seeding..."
            # Clean up existing data first
            echo "Cleaning up existing data..."
            ./cmd -action=seed -seeder=location -force-reseed
            sleep 2
        else
            echo "No tables found. Running migrations and seeding..."
        fi
        
        run_seeding
    fi
else
    echo "Auto-seeding skipped (AUTO_SEED=$AUTO_SEED, ENVIRONMENT=$ENVIRONMENT)"
fi

echo "Starting application..."
exec ./main
