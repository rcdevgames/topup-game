# Makefile for Waw Store Backend

.PHONY: help build run test clean docker-build docker-run docker-stop install dev

# Default target
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Development commands
install: ## Install dependencies
	go mod tidy
	go mod download

dev: ## Run development server with auto-reload (requires air)
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "Installing air for hot reload..."; \
		go install github.com/cosmtrek/air@latest; \
		air; \
	fi

run: ## Run the application
	go run cmd/server/main.go

build: ## Build the application
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/main cmd/server/main.go

test: ## Run tests
	go test -v ./...

test-cover: ## Run tests with coverage
	go test -v -cover ./...

clean: ## Clean build artifacts
	rm -rf bin/
	go clean

# Database commands
db-migrate: ## Run database migrations (automatically done on startup)
	@echo "Migrations run automatically on server startup"

db-seed: ## Seed database with sample data (automatically done on startup)
	@echo "Seeding happens automatically on first startup"

# Docker commands
docker-build: ## Build Docker image
	docker build -t topup-backend .

docker-run: ## Run Docker container
	docker run -p 8080:8080 --env-file .env topup-backend

docker-compose-up: ## Start all services with Docker Compose
	docker-compose up -d

docker-compose-down: ## Stop all services
	docker-compose down

docker-compose-logs: ## View Docker Compose logs
	docker-compose logs -f

# Swagger documentation
swagger: ## Generate Swagger documentation
	@if command -v swag > /dev/null; then \
		swag init -g cmd/server/main.go -o docs; \
	else \
		echo "Installing swag..."; \
		go install github.com/swaggo/swag/cmd/swag@latest; \
		swag init -g cmd/server/main.go -o docs; \
	fi

# Code quality
fmt: ## Format code
	go fmt ./...

vet: ## Run go vet
	go vet ./...

lint: ## Run golangci-lint (requires golangci-lint)
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "Please install golangci-lint: https://golangci-lint.run/usage/install/"; \
	fi

# Environment setup
env: ## Copy .env.example to .env
	@if [ ! -f .env ]; then \
		cp .env.example .env; \
		echo ".env file created from .env.example"; \
		echo "Please edit .env file with your configuration"; \
	else \
		echo ".env file already exists"; \
	fi

# Production build
prod-build: ## Build for production
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o bin/main cmd/server/main.go

# Development setup (run once)
setup: env install swagger ## Setup development environment
	@echo "Development environment setup complete!"
	@echo "1. Edit .env file with your database credentials"
	@echo "2. Make sure PostgreSQL is running"
	@echo "3. Run 'make dev' to start development server"