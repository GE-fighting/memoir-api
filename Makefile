# Memoir API Makefile

.PHONY: help build run migrate-up migrate-down migrate-status clean test

# Default target
help:
	@echo "Available commands:"
	@echo "  build        - Build the API server"
	@echo "  run          - Run the API server"
	@echo "  migrate-up   - Run database migrations"
	@echo "  migrate-down - Rollback database migrations"
	@echo "  migrate-status - Check migration status"
	@echo "  clean        - Clean build artifacts"
	@echo "  test         - Run tests"

# Build the API server
build:
	@echo "Building API server..."
	go build -o bin/api cmd/api/main.go

# Run the API server
run:
	@echo "Starting API server..."
	go run cmd/api/main.go

# Database migrations
migrate-up:
	@echo "Running database migrations..."
	go run cmd/migrate/main.go -action=up

migrate-down:
	@echo "Rolling back database migrations..."
	go run cmd/migrate/main.go -action=down

migrate-status:
	@echo "Checking migration status..."
	go run cmd/migrate/main.go -action=status

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/

# Run tests
test:
	@echo "Running tests..."
	go test ./...

# Build migration tool
build-migrate:
	@echo "Building migration tool..."
	go build -o bin/migrate cmd/migrate/main.go
