.PHONY: help run build test clean migrate-up migrate-down dev

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

run: ## Run the application
	go run cmd/api/main.go

build: ## Build the application
	go build -o bin/server cmd/api/main.go

dev: ## Run in development mode with auto-reload (requires air)
	air

test: ## Run tests
	go test -v ./...

clean: ## Clean build artifacts
	rm -rf bin/

install-deps: ## Install Go dependencies
	go mod download
	go mod tidy

migrate-up: ## Run database migrations
	psql -U postgres -d perjalanan_dinas -f migrations/001_initial_schema.sql
	psql -U postgres -d perjalanan_dinas -f migrations/002_seed_data.sql

migrate-down: ## Drop all tables (WARNING: This will delete all data)
	@echo "WARNING: This will drop all tables and delete all data!"
	@read -p "Are you sure? [y/N] " -n 1 -r; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		psql -U postgres -d perjalanan_dinas -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"; \
	fi

db-reset: migrate-down migrate-up ## Reset database (drop and recreate)

setup: install-deps ## Initial setup
	@echo "Setting up backend..."
	@if [ ! -f .env ]; then cp .env.example .env; echo "Created .env file. Please update with your configuration."; fi
	@echo "Setup complete! Run 'make migrate-up' to setup database."

docker-build: ## Build Docker image
	docker build -t perjalanan-dinas-backend .

docker-run: ## Run in Docker container
	docker run -p 8080:8080 --env-file .env perjalanan-dinas-backend
