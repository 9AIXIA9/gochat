# Makefile for GoChat backend on Windows (requires GNU make)
# Use cmd.exe shell semantics for docker commands

SHELL := cmd.exe

PROJECT_NAME := backend
COMPOSE_FILE := docker-compose.yml


# Helper to run docker compose with file and project name
DC := docker compose -f $(COMPOSE_FILE) -p $(PROJECT_NAME)

.PHONY: help
help:
	@echo "Available targets:"
	@echo "  make up           - Build (no cache) and start all services in background"
	@echo "  make up-fast      - Up using build cache"
	@echo "  make down         - Stop and remove containers (keep named volumes)"
	@echo "  make destroy      - Stop and remove containers, images, and named volumes"
	@echo "  make rebuild      - Remove old app image and rebuild with no cache"
	@echo "  make restart      - Restart services using build cache"
	@echo "  make kafka-init   - Create Kafka topics before starting the app"
	@echo "  make db-migrate   - Apply DB migrations using host port"
	@echo "  make logs         - Tail app logs"
	@echo "  make logs-app     - Tail only app logs"
	@echo "  make ps           - Show service status"
	@echo "  make status       - Alias for ps"
	@echo "  make hooks        - Install Git hooks"
	@echo "  make precommit    - Run pre-commit hook manually"
	@echo "  make test-all     - Run all tests with race detector and shuffle"
	@echo "  make lint         - Run golangci-lint if installed"
	@echo "  make swagger      - Generate Swagger docs from annotations"
	@echo "  make swagger-generate - Generate Swagger docs from annotations"
	@echo "  make fix-swagger  - Fix Swagger example fields (example -> x-example)"
	@echo "  make fix-swagger-dry-run - Check what would be fixed (dry run)"
	@echo "  make build        - Build the app binary"
	@echo "  make run          - Run the app locally"
	@echo "  make migrate-status - Show migration status"
	@echo "  make migrate-up   - Apply all pending migrations"
	@echo "  make migrate-down - Roll back the most recent migration"
	@echo "  make migrate-reset - Roll back all migrations"
	@echo "  make migrate-create - Create a new SQL migration (usage: make migrate-create name=add_table)"
	@echo "  make migrate-status-host - Show migration status (host -> 127.0.0.1:$$(MYSQL_HOST_PORT))"
	@echo "  make migrate-up-host - Apply migrations using host port"
	@echo "  make migrate-down-host - Roll back last migration using host port"
	@echo "  make migrate-reset-host - Reset all migrations using host port"
	@echo "  make migrate-verify - Fail if pending migrations (container network DSN)"
	@echo "  make migrate-verify-host - Fail if pending migrations (host 127.0.0.1:$$(MYSQL_HOST_PORT))"

.PHONY: up
up: ## Build (no cache) and up -d
	-$(DC) down
	$(DC) build --no-cache mysql redis kafka jaeger prometheus grafana otel-collector
	$(MAKE) kafka-init
	$(MAKE) db-migrate
	$(DC) build --no-cache app
	$(DC) up -d --remove-orphans mysql redis kafka jaeger prometheus grafana otel-collector app

.PHONY: down
down: ## Stop and remove containers (keep volumes)
	$(DC) down

.PHONY: destroy
destroy: ## Stop and remove containers, images, and named volumes
	-$(DC) down --rmi all -v --remove-orphans
	-@for /f "tokens=1" %%%%i in ('docker images -q gochat-backend 2^>NUL') do docker rmi -f %%%%i 2>NUL || exit /b 0
	-@for /f "tokens=1" %%%%i in ('docker images -q apache/kafka 2^>NUL') do docker rmi -f %%%%i 2>NUL || exit /b 0
	-@for /f "tokens=1" %%%%i in ('docker images -q mysql 2^>NUL') do docker rmi -f %%%%i 2>NUL || exit /b 0
	-@for /f "tokens=1" %%%%i in ('docker images -q redis 2^>NUL') do docker rmi -f %%%%i 2>NUL || exit /b 0

.PHONY: rebuild
rebuild: ## Remove old app image and rebuild without cache
	-@for /f "tokens=1" %%%%i in ('docker images -q gochat-backend 2^>NUL') do docker rmi -f %%%%i 2>NUL || exit /b 0
	$(DC) build --no-cache

.PHONY: logs
logs:
	$(DC) logs -f app

.PHONY: ps
ps:
	$(DC) ps

.PHONY: restart
restart: ## Restart using build cache (no --no-cache)
	$(DC) up -d --build mysql redis kafka jaeger prometheus grafana otel-collector
	$(MAKE) kafka-init
	$(MAKE) db-migrate
	$(DC) up -d --remove-orphans --build app

.PHONY: kafka-init
kafka-init: ## Create Kafka topics using a one-off container
	docker run --rm --network $(PROJECT_NAME)_default -v "$(CURDIR)/deployment/kafka-init/init_topics.sh:/init_topics.sh:ro" -v "$(CURDIR)/deployment/kafka-init/topics.txt:/topics.txt:ro" docker.io/confluentinc/cp-kafka:8.0.0 sh /init_topics.sh

.PHONY: db-migrate
db-migrate: ## Apply DB migrations using host port
	$(MAKE) migrate-up-host

.PHONY: hooks
hooks: ## Configure Git to use the versioned hooks in .githooks
	@git config core.hooksPath .githooks
	@echo Hooks installed to .githooks

.PHONY: precommit
precommit: ## Run pre-commit hook logic locally
	@sh .githooks/pre-commit

.PHONY: test-all
test-all: ## Run all tests with race detector and shuffle
	go test -race -shuffle=on ./...

.PHONY: lint
lint: ## Run golangci-lint if installed
	golangci-lint run ./...

# Swagger documentation
.PHONY: swagger
swagger: swagger-generate fix-swagger ## Generate Swagger docs and fix example fields

.PHONY: swagger-generate
swagger-generate: ## Generate Swagger docs from annotations
	swag init -g cmd/api/main.go --parseDependency --parseInternal
	@echo "Swagger docs generated in docs/ directory"

.PHONY: fix-swagger
fix-swagger: ## Fix Swagger example fields (example -> x-example)
	@echo "Fixing Swagger example fields..."
	python scripts/fix_swagger.py docs/docs.go docs/swagger.yaml docs/swagger.json --backup
	@echo "Swagger files fixed"

.PHONY: fix-swagger-dry-run
fix-swagger-dry-run: ## Check what would be fixed (dry run)
	@echo "Checking Swagger files (dry run)..."
	python scripts/fix_swagger.py docs/docs.go docs/swagger.yaml docs/swagger.json --dry-run

# Build & run app locally (without Docker)
APP_MAIN := ./cmd/api/main.go

.PHONY: build
build: ## Build the app binary
	go build -o bin/gochat-app $(APP_MAIN)

.PHONY: run
run: ## Run the app locally (use -config to override)
	go run $(APP_MAIN) -config $(COMPOSE_FILE_DIR)/config/config.yaml

# Compose utilities
COMPOSE_FILE_DIR := .

.PHONY: up-fast
up-fast: ## Up using build cache
	$(DC) up -d --build mysql redis kafka jaeger prometheus grafana otel-collector
	$(MAKE) kafka-init
	$(MAKE) db-migrate
	$(DC) up -d --remove-orphans --build app

.PHONY: logs-app
logs-app: ## Tail only app logs
	$(DC) logs -f app

.PHONY: status
status: ps ## Alias for ps

# Load environment overrides (optional). Default uses app env file.
ENV_FILE ?= .env.development
-include $(ENV_FILE)

# Default DB env (override via environment or .env if desired)
DB_HOST ?= 127.0.0.1
DB_PORT ?= 13306
DB_USERNAME ?= gochat_app
DB_PASSWORD ?= gochat
DB_NAME ?= gochat

# Goose migrations
MIGRATIONS_DIR := db/migrations
DB_DSN := $(DB_USERNAME):$(DB_PASSWORD)@tcp($(DB_HOST):$(DB_PORT))/$(DB_NAME)?parseTime=true&charset=utf8mb4

# Goose executable (override if needed, default assumes goose in PATH)
GOOSE ?= goose

.PHONY: migrate-status
migrate-status: ## Show migration status
	@$(GOOSE) -dir $(MIGRATIONS_DIR) mysql "$(DB_DSN)" status

.PHONY: migrate-up
migrate-up: ## Apply all pending migrations
	@$(GOOSE) -dir $(MIGRATIONS_DIR) mysql "$(DB_DSN)" up

.PHONY: migrate-down
migrate-down: ## Roll back the most recent migration
	@$(GOOSE) -dir $(MIGRATIONS_DIR) mysql "$(DB_DSN)" down

.PHONY: migrate-reset
migrate-reset: ## Roll back all migrations
	@$(GOOSE) -dir $(MIGRATIONS_DIR) mysql "$(DB_DSN)" reset

.PHONY: migrate-create
migrate-create: ## Create a new SQL migration (usage: make migrate-create name=add_table)
	@if not defined name (echo Usage: make migrate-create name=add_table & exit /b 1)
	@$(GOOSE) -dir $(MIGRATIONS_DIR) create $(name) sql

# Host-based DSN for running goose from host machine against published MySQL port
MYSQL_HOST_PORT ?= 13306
DB_DSN_HOST := $(DB_USERNAME):$(DB_PASSWORD)@tcp(127.0.0.1:$(MYSQL_HOST_PORT))/$(DB_NAME)?parseTime=true&charset=utf8mb4

.PHONY: migrate-status-host
migrate-status-host: ## Show migration status (host -> 127.0.0.1:$(MYSQL_HOST_PORT))
	@$(GOOSE) -dir $(MIGRATIONS_DIR) mysql "$(DB_DSN_HOST)" status

.PHONY: migrate-up-host
migrate-up-host: ## Apply migrations using host port
	@$(GOOSE) -dir $(MIGRATIONS_DIR) mysql "$(DB_DSN_HOST)" up

.PHONY: migrate-down-host
migrate-down-host: ## Roll back last migration using host port
	@$(GOOSE) -dir $(MIGRATIONS_DIR) mysql "$(DB_DSN_HOST)" down

.PHONY: migrate-reset-host
migrate-reset-host: ## Reset all migrations using host port
	@$(GOOSE) -dir $(MIGRATIONS_DIR) mysql "$(DB_DSN_HOST)" reset

.PHONY: migrate-verify
migrate-verify: ## Fail if pending migrations (container network DSN)
	@$(GOOSE) -dir $(MIGRATIONS_DIR) mysql "$(DB_DSN)" status | findstr /I "Pending" >nul && (echo Pending migrations detected. Please run make migrate-up. & exit /b 1) || (echo Schema up-to-date.)

.PHONY: migrate-verify-host
migrate-verify-host: ## Fail if pending migrations (host 127.0.0.1:$(MYSQL_HOST_PORT))
	@$(GOOSE) -dir $(MIGRATIONS_DIR) mysql "$(DB_DSN_HOST)" status | findstr /I "Pending" >nul && (echo Pending migrations detected. Please run make migrate-up-host. & exit /b 1) || (echo Schema up-to-date.)