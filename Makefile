# Makefile for Asset Management System

# Variables
APP_NAME=be-lecsens-asset-management
DOCKER_IMAGE=be-lecsens/asset-management
GO_FILES=$(shell find . -name '*.go' -type f)

# Default target
.PHONY: help
help:
	@echo "Available commands:"
	@echo "  build          - Build the Go application"
	@echo "  run            - Run the application locally"
	@echo "  docker-build   - Build Docker image"
	@echo "  docker-run     - Run with Docker Compose"
	@echo "  docker-stop    - Stop Docker Compose services"
	@echo "  docker-clean   - Stop and remove Docker containers/volumes"
	@echo "  migrate        - Run database migrations"
	@echo "  seed           - Run database seeding"
	@echo "  seed-all       - Run complete database seeding"
	@echo "  clean          - Clean build artifacts"
	@echo "  test           - Run tests"

# Build the Go application
.PHONY: build
build:
	@echo "Building application..."
	go build -o bin/main .
	go build -o bin/cmd ./helpers/cmd/cmd.go

# Run the application locally
.PHONY: run
run: build
	@echo "Running application..."
	./bin/main

# Build Docker image
.PHONY: docker-build
docker-build:
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE) .

# Run with Docker Compose
.PHONY: docker-run
docker-run:
	@echo "Starting services with Docker Compose..."
	docker-compose up --build -d

# Stop Docker Compose services
.PHONY: docker-stop
docker-stop:
	@echo "Stopping Docker Compose services..."
	docker-compose down

# Clean Docker containers and volumes
.PHONY: docker-clean
docker-clean:
	@echo "Cleaning Docker containers and volumes..."
	docker-compose down -v
	docker system prune -f

# Run database migrations
.PHONY: migrate
migrate:
	@echo "Running database migrations..."
	go run ./helpers/cmd/cmd.go -action=migrate -force

# Run database seeding
.PHONY: seed
seed:
	@echo "Running database seeding..."
	go run ./helpers/cmd/cmd.go -action=seed -seeder=all -force

# Run complete database seeding with all data
.PHONY: seed-all
seed-all:
	@echo "Running complete database seeding..."
	go run ./helpers/cmd/cmd.go -action=seed -seeder=location -force
	go run ./helpers/cmd/cmd.go -action=seed -seeder=asset-type -force
	go run ./helpers/cmd/cmd.go -action=seed -seeder=sensor-type -force
	go run ./helpers/cmd/cmd.go -action=seed -seeder=measurement-type -force
	go run ./helpers/cmd/cmd.go -action=seed -seeder=measurement-field -force
	go run ./helpers/cmd/cmd.go -action=seed -seeder=asset -force
	go run ./helpers/cmd/cmd.go -action=seed -seeder=asset-sensor -force
	go run ./helpers/cmd/cmd.go -action=seed -seeder=threshold -force
	go run ./helpers/cmd/cmd.go -action=seed -seeder=sensor-status -force
	go run ./helpers/cmd/cmd.go -action=seed -seeder=reading -days=7 -force
	go run ./helpers/cmd/cmd.go -action=seed -seeder=sensor-logs -force
	go run ./helpers/cmd/cmd.go -action=seed -seeder=alert -force

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	go clean

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	go test ./...

# Development setup
.PHONY: dev-setup
dev-setup: docker-run
	@echo "Waiting for services to be ready..."
	sleep 10
	@echo "Development environment is ready!"
	@echo "Application: http://localhost:3122"
	@echo "Database: localhost:5444"

# Production build
.PHONY: prod-build
prod-build:
	@echo "Building for production..."
	ENVIRONMENT=production AUTO_SEED=false docker-compose -f docker-compose.yaml up --build -d

# View logs
.PHONY: logs
logs:
	docker-compose logs -f app

# View database logs
.PHONY: db-logs
db-logs:
	docker-compose logs -f db
