# Makefile for GoChat backend on Windows (requires GNU make)
# Use cmd.exe shell semantics for docker commands

# 使用自定义 env 文件
ENV_FILE := ./config/docker/.env

PROJECT_NAME := backend
COMPOSE_FILE := docker-compose.yml


# Helper to run docker compose with file and project name
DC := docker compose --env-file $(ENV_FILE) -f $(COMPOSE_FILE) -p $(PROJECT_NAME)

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