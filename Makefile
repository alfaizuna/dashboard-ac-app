# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=main
BINARY_UNIX=$(BINARY_NAME)_unix

# Docker parameters
DOCKER_COMPOSE=docker-compose
DOCKER_COMPOSE_PROD=docker-compose -f docker-compose.prod.yml

.PHONY: all build clean test coverage deps run docker-build docker-up docker-down docker-logs help

all: test build

# Build the application
build:
	$(GOBUILD) -o $(BINARY_NAME) -v ./cmd/main.go

# Clean build files
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)

# Run tests
test:
	$(GOTEST) -v ./...

# Run tests with coverage
coverage:
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out

# Download dependencies
deps:
	$(GOMOD) download
	$(GOMOD) tidy

# Run the application locally
run:
	$(GOCMD) run ./cmd/main.go

# Build for Linux
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v ./cmd/main.go

# Docker commands
docker-build:
	docker build -t dashboard-ac-backend .

docker-up:
	$(DOCKER_COMPOSE) up -d

docker-down:
	$(DOCKER_COMPOSE) down

docker-logs:
	$(DOCKER_COMPOSE) logs -f

docker-restart:
	$(DOCKER_COMPOSE) restart

docker-clean:
	$(DOCKER_COMPOSE) down -v --remove-orphans
	docker system prune -f

# Production Docker commands
docker-prod-up:
	$(DOCKER_COMPOSE_PROD) up -d

docker-prod-down:
	$(DOCKER_COMPOSE_PROD) down

docker-prod-logs:
	$(DOCKER_COMPOSE_PROD) logs -f

# Database commands
db-migrate:
	$(GOCMD) run ./migrations/migrate.go

db-reset:
	$(DOCKER_COMPOSE) exec postgres psql -U postgres -d dashboard_ac_dev -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"

# Development workflow
dev: deps docker-up
	@echo "Development environment is running!"
	@echo "API: http://localhost:8088"
	@echo "Drizzle Studio: http://localhost:4983"

# Linting
lint:
	golangci-lint run

# Format code
fmt:
	$(GOCMD) fmt ./...

# Help
help:
	@echo "Available commands:"
	@echo "  build         - Build the application"
	@echo "  clean         - Clean build files"
	@echo "  test          - Run tests"
	@echo "  coverage      - Run tests with coverage"
	@echo "  deps          - Download dependencies"
	@echo "  run           - Run the application locally"
	@echo "  build-linux   - Build for Linux"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-up     - Start Docker containers"
	@echo "  docker-down   - Stop Docker containers"
	@echo "  docker-logs   - View Docker logs"
	@echo "  docker-clean  - Clean Docker containers and volumes"
	@echo "  dev           - Setup development environment"
	@echo "  lint          - Run linter"
	@echo "  fmt           - Format code"
	@echo "  help          - Show this help message"