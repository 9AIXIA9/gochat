# Makefile for GoChat backend on Windows (requires GNU make)
# Use cmd.exe shell semantics for docker commands

PROJECT_NAME := backend
COMPOSE_FILE := docker-compose.yml


# Helper to run docker compose with file and project name
DC := docker compose -f $(COMPOSE_FILE) -p $(PROJECT_NAME)

.PHONY: help
help:
	@echo Available targets:
	@echo   make up           - Build (no cache) and start all services in background
	@echo   make down         - Stop and remove containers (keep named volumes)
	@echo   make destroy      - Stop and remove containers, images, and named volumes
	@echo   make rebuild      - Remove old app image and rebuild with no cache
	@echo   make restart      - Restart services using build cache
	@echo   make logs         - Tail app logs
	@echo   make ps           - Show service status
	@echo   make hooks        - Install Git hooks
	@echo   make precommit    - Run pre-commit hook manually

.PHONY: up
up: ## Build (no cache) and up -d
	-$(DC) down
	$(DC) build --no-cache
	$(DC) up -d

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
	$(DC) up -d --build

.PHONY: hooks
hooks: ## Configure Git to use the versioned hooks in .githooks
	@git config core.hooksPath .githooks
	@echo Hooks installed to .githooks

.PHONY: precommit
precommit: ## Run pre-commit hook logic locally
	@sh .githooks/pre-commit

# Goose migrations
MIGRATIONS_DIR := db/migrations
DB_DSN := mysql://$(DB_USERNAME):$(DB_PASSWORD)@tcp($(DB_HOST):$(DB_PORT))/$(DB_NAME)?parseTime=true&charset=utf8mb4

.PHONY: migrate-status
migrate-status: ## Show migration status
	@go run github.com/pressly/goose/v3/cmd/goose -dir $(MIGRATIONS_DIR) mysql "$(DB_DSN)" status

.PHONY: migrate-up
migrate-up: ## Apply all pending migrations
	@go run github.com/pressly/goose/v3/cmd/goose -dir $(MIGRATIONS_DIR) mysql "$(DB_DSN)" up

.PHONY: migrate-down
migrate-down: ## Roll back the most recent migration
	@go run github.com/pressly/goose/v3/cmd/goose -dir $(MIGRATIONS_DIR) mysql "$(DB_DSN)" down

.PHONY: migrate-reset
migrate-reset: ## Roll back all migrations
	@go run github.com/pressly/goose/v3/cmd/goose -dir $(MIGRATIONS_DIR) mysql "$(DB_DSN)" reset

.PHONY: migrate-create
migrate-create: ## Create a new SQL migration (usage: make migrate-create name=add_table)
	@if not defined name (echo Usage: make migrate-create name=add_table & exit /b 1)
	@go run github.com/pressly/goose/v3/cmd/goose -dir $(MIGRATIONS_DIR) create $(name) sql

# Default DB env (override via environment or .env if desired)
DB_HOST ?= mysql
DB_PORT ?= 3306
DB_USERNAME ?= gochat_app
DB_PASSWORD ?= gochat
DB_NAME ?= gochat

# Host-based DSN for running goose from host machine against published MySQL port
MYSQL_HOST_PORT ?= 13306
DB_DSN_HOST := mysql://$(DB_USERNAME):$(DB_PASSWORD)@tcp(127.0.0.1:$(MYSQL_HOST_PORT))/$(DB_NAME)?parseTime=true&charset=utf8mb4

.PHONY: migrate-status-host
migrate-status-host: ## Show migration status (host -> 127.0.0.1:$(MYSQL_HOST_PORT))
	@go run github.com/pressly/goose/v3/cmd/goose -dir $(MIGRATIONS_DIR) mysql "$(DB_DSN_HOST)" status

.PHONY: migrate-up-host
migrate-up-host: ## Apply migrations using host port
	@go run github.com/pressly/goose/v3/cmd/goose -dir $(MIGRATIONS_DIR) mysql "$(DB_DSN_HOST)" up

.PHONY: migrate-down-host
migrate-down-host: ## Roll back last migration using host port
	@go run github.com/pressly/goose/v3/cmd/goose -dir $(MIGRATIONS_DIR) mysql "$(DB_DSN_HOST)" down

.PHONY: migrate-reset-host
migrate-reset-host: ## Reset all migrations using host port
	@go run github.com/pressly/goose/v3/cmd/goose -dir $(MIGRATIONS_DIR) mysql "$(DB_DSN_HOST)" reset

.PHONY: migrate-verify
migrate-verify: ## Fail if pending migrations (container network DSN)
	@go run github.com/pressly/goose/v3/cmd/goose -dir $(MIGRATIONS_DIR) mysql "$(DB_DSN)" status | findstr /I "Pending" >nul && (echo Pending migrations detected. Please run make migrate-up. & exit /b 1) || (echo Schema up-to-date.)

.PHONY: migrate-verify-host
migrate-verify-host: ## Fail if pending migrations (host 127.0.0.1:$(MYSQL_HOST_PORT))
	@go run github.com/pressly/goose/v3/cmd/goose -dir $(MIGRATIONS_DIR) mysql "$(DB_DSN_HOST)" status | findstr /I "Pending" >nul && (echo Pending migrations detected. Please run make migrate-up-host. & exit /b 1) || (echo Schema up-to-date.)
