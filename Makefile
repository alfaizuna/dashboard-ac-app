.PHONY: build run test test-unit test-integration clean migrate seed docker-up docker-down docker-build lint

# Variables
BINARY_NAME=dashboard-ac-backend
MAIN_PATH=./cmd/main.go
MIGRATE_PATH=./migrations/migrate.go

# Build the application
build:
	@echo "Building application..."
	go build -o bin/$(BINARY_NAME) $(MAIN_PATH)

# Run the application
run:
	@echo "Running application..."
	go run $(MAIN_PATH)

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run unit tests
test-unit:
	@echo "Running unit tests..."
	go test -v ./tests/unit/...

# Run integration tests
test-integration:
	@echo "Running integration tests..."
	go test -v ./tests/integration/...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf bin/
	rm -f coverage.out coverage.html

# Run database migrations
migrate:
	@echo "Running database migrations..."
	go run $(MIGRATE_PATH)

# Seed database (alias for migrate since it includes seeding)
seed: migrate

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Lint the code
lint:
	@echo "Running linter..."
	golangci-lint run

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Docker commands
docker-build:
	@echo "Building Docker image..."
	docker build -t $(BINARY_NAME) .

docker-up:
	@echo "Starting services with Docker Compose..."
	docker-compose up -d

docker-down:
	@echo "Stopping services..."
	docker-compose down

docker-logs:
	@echo "Showing logs..."
	docker-compose logs -f

# Development workflow
dev: deps fmt lint test build

# Production build
prod-build:
	@echo "Building for production..."
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/$(BINARY_NAME) $(MAIN_PATH)

# Help
help:
	@echo "Available commands:"
	@echo "  build         - Build the application"
	@echo "  run           - Run the application"
	@echo "  test          - Run all tests"
	@echo "  test-unit     - Run unit tests only"
	@echo "  test-integration - Run integration tests only"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  clean         - Clean build artifacts"
	@echo "  migrate       - Run database migrations and seeding"
	@echo "  seed          - Alias for migrate"
	@echo "  deps          - Install dependencies"
	@echo "  lint          - Run linter"
	@echo "  fmt           - Format code"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-up     - Start services with Docker Compose"
	@echo "  docker-down   - Stop services"
	@echo "  docker-logs   - Show Docker logs"
	@echo "  dev           - Run development workflow (deps, fmt, lint, test, build)"
	@echo "  prod-build    - Build for production"
	@echo "  help          - Show this help message"