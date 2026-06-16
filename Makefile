.PHONY: help build test dev prod stop clean migrate swagger lint fmt

# Default target
help: ## Show this help
	@echo "Anexis Server - Available Commands"
	@echo "==================================="
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

# Build targets
build: ## Build the API server
	@echo "Building API server..."
	cd apps/api && go build -o ../../bin/api-server ./cmd/server
	@echo "Build complete: bin/api-server"

build-docker: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t anexis-api -f apps/api/Dockerfile .

# Development
dev: ## Start development environment with Docker
	@echo "Starting development environment..."
	docker compose up -d
	@echo "Development server running at http://localhost:8080"
	@echo "Swagger docs at http://localhost:8080/swagger/index.html"

dev-local: ## Run API server locally (requires local PostgreSQL)
	@echo "Starting local development server..."
	cd apps/api && go run ./cmd/server

watch: ## Run with hot reload (requires air: go install github.com/air-verse/air@latest)
	@echo "Starting with hot reload..."
	cd apps/api && $(shell go env GOPATH)/bin/air

# Production
prod: ## Start production environment with load balancing
	@echo "Starting production environment..."
	docker compose -f docker-compose.prod.yml up -d
	@echo "Production server running at http://localhost:80"

prod-scale: ## Scale API to N replicas (usage: make prod-scale N=5)
	@echo "Scaling API to $(N) replicas..."
	docker compose -f docker-compose.prod.yml up -d --scale api=$(N)

# Stop/Clean
stop: ## Stop all Docker containers
	@echo "Stopping containers..."
	docker compose down
	docker compose -f docker-compose.prod.yml down 2>/dev/null || true

clean: ## Clean build artifacts and Docker volumes
	@echo "Cleaning up..."
	rm -rf bin/
	docker compose down -v 2>/dev/null || true
	docker compose -f docker-compose.prod.yml down -v 2>/dev/null || true

# Database
migrate: ## Run database migrations
	@echo "Running migrations..."
	@set -a && . ./.env && set +a && atlas migrate apply --env gorm

migrate-new: ## Create new migration (usage: make migrate-new NAME=add_users)
	@echo "Creating migration: $(NAME)"
	@set -a && . ./.env && set +a && atlas migrate diff $(NAME) --env gorm

migrate-status: ## Show migration status
	@set -a && . ./.env && set +a && atlas migrate status --env gorm

# Testing
test: ## Run all tests
	@echo "Running tests..."
	cd packages/database && go test -v ./...
	cd apps/api && go test -v ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	cd apps/api && go test -coverprofile=coverage.out ./...
	cd apps/api && go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: apps/api/coverage.html"

# Code quality
lint: ## Run linter (requires golangci-lint)
	@echo "Running linter..."
	golangci-lint run ./apps/api/... ./packages/database/...

fmt: ## Format code
	@echo "Formatting code..."
	go fmt ./...
	gofmt -s -w .

vet: ## Run go vet
	@echo "Running go vet..."
	cd packages/database && go vet ./...
	cd apps/api && go vet ./...

# Documentation
swagger: ## Generate Swagger documentation
	@echo "Generating Swagger docs..."
	cd apps/api && swag init -g cmd/server/main.go -o docs
	@echo "Swagger docs generated: apps/api/docs/"

swagger-serve: ## Serve Swagger UI locally
	@echo "Open http://localhost:8080/swagger/index.html"

docs: swagger ## Alias for swagger

# Dependencies
deps: ## Download dependencies
	@echo "Downloading dependencies..."
	go work sync
	cd packages/database && go mod tidy
	cd apps/api && go mod tidy

deps-update: ## Update dependencies
	@echo "Updating dependencies..."
	cd packages/database && go get -u ./...
	cd apps/api && go get -u ./...
	$(MAKE) deps

# Utilities
health: ## Check API health
	@curl -s http://localhost:8080/health | jq .

logs: ## Show Docker logs
	docker compose logs -f api

logs-prod: ## Show production Docker logs
	docker compose -f docker-compose.prod.yml logs -f

psql: ## Connect to PostgreSQL
	docker compose exec postgres psql -U postgres -d anexis

# Install tools
install-tools: ## Install development tools
	@echo "Installing development tools..."
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/air-verse/air@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "Tools installed successfully"
