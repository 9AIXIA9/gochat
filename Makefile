# Cross-platform Makefile for GoChat backend.
# Keep it focused on the commands you actually use day to day.

APP_NAME ?= gochat-app
APP_CMD ?= ./cmd/api
COMPOSE_PROJECT ?= backend
COMPOSE_FILE ?= docker-compose.yml
OBS_FILE ?= docker-compose.observability.yml

GO ?= go
DOCKER ?= docker
SWAG ?= swag
GOOSE ?= goose
PYTHON ?= python

ifeq ($(OS),Windows_NT)
APP_BIN := $(APP_NAME).exe
else
APP_BIN := $(APP_NAME)
endif

MIGRATIONS_DIR ?= db/migrations
MYSQL_HOST_PORT ?= 13306
DB_USERNAME ?= gochat_app
DB_PASSWORD ?= gochat
DB_NAME ?= gochat
DB_HOST ?= 127.0.0.1
DB_PORT ?= $(MYSQL_HOST_PORT)
DB_DSN := $(DB_USERNAME):$(DB_PASSWORD)@tcp($(DB_HOST):$(DB_PORT))/$(DB_NAME)?parseTime=true&charset=utf8mb4

DC := $(DOCKER) compose -f $(COMPOSE_FILE) -p $(COMPOSE_PROJECT)
DC_OBS := $(DOCKER) compose -f $(COMPOSE_FILE) -f $(OBS_FILE) -p $(COMPOSE_PROJECT)

.PHONY: help
help:
	$(info Available targets:)
	$(info make fmt              - Format Go code)
	$(info make test             - Run all tests)
	$(info make test-race        - Run tests with race detector)
	$(info make tidy             - Run go mod tidy)
	$(info make build            - Build the app binary)
	$(info make run              - Run the app locally)
	$(info make lint             - Run golangci-lint)
	$(info make swagger          - Generate Swagger docs)
	$(info make swagger-fix      - Fix Swagger example fields)
	$(info make swagger-dry-run  - Preview Swagger fixes)
	$(info make up               - Start core services and app)
	$(info make up-obs           - Start core services, observability, and app)
	$(info make compose-migrate  - Run migrations in the Compose network)
	$(info make down             - Stop core services)
	$(info make ps               - Show service status)
	$(info make logs             - Tail app logs)
	$(info make restart          - Restart core services)
	$(info make migrate-status   - Show local host migration status)
	$(info make migrate-up       - Apply local host migrations)
	$(info make migrate-down     - Roll back the latest local migration)
	$(info make migrate-reset    - Roll back all local migrations)
	$(info make migrate-create   - Create a new SQL migration; usage: make migrate-create name=add_table)
	@:

.PHONY: fmt
fmt:
	$(GO) fmt ./...

.PHONY: test
test:
	$(GO) test ./...

.PHONY: test-race
test-race:
	$(GO) test -race ./...

.PHONY: tidy
tidy:
	$(GO) mod tidy

.PHONY: build
build:
	$(GO) build -o $(APP_BIN) $(APP_CMD)

.PHONY: run
run:
	$(GO) run $(APP_CMD)

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: swagger
swagger:
	$(SWAG) init -g cmd/api/main.go --parseDependency --parseInternal
	@echo Swagger docs generated in docs/

.PHONY: swagger-fix
swagger-fix:
	$(PYTHON) scripts/fix_swagger.py docs/docs.go docs/swagger.yaml docs/swagger.json --backup

.PHONY: swagger-dry-run
swagger-dry-run:
	$(PYTHON) scripts/fix_swagger.py docs/docs.go docs/swagger.yaml docs/swagger.json --dry-run

.PHONY: compose-migrate
compose-migrate:
	$(DC) run --rm --build migrate

.PHONY: up
up:
	$(DC) up -d --build mysql redis kafka kafka-init
	$(DC) run --rm --build migrate
	$(DC) up -d --build app

.PHONY: up-obs
up-obs:
	$(DC_OBS) up -d --build mysql redis kafka kafka-init otel-collector jaeger prometheus grafana
	$(DC_OBS) run --rm --build migrate
	$(DC_OBS) up -d --build app

.PHONY: down
down:
	$(DC) down --remove-orphans

.PHONY: ps
ps:
	$(DC) ps

.PHONY: logs
logs:
	$(DC) logs -f app

.PHONY: restart
restart:
	$(MAKE) down
	$(MAKE) up

.PHONY: migrate-status
migrate-status:
	@$(GOOSE) -dir $(MIGRATIONS_DIR) mysql "$(DB_DSN)" status

.PHONY: migrate-up
migrate-up:
	@$(GOOSE) -dir $(MIGRATIONS_DIR) mysql "$(DB_DSN)" up

.PHONY: migrate-down
migrate-down:
	@$(GOOSE) -dir $(MIGRATIONS_DIR) mysql "$(DB_DSN)" down

.PHONY: migrate-reset
migrate-reset:
	@$(GOOSE) -dir $(MIGRATIONS_DIR) mysql "$(DB_DSN)" reset

.PHONY: migrate-create
migrate-create:
	$(if $(strip $(name)),,$(error Usage: make migrate-create name=add_table))
	@$(GOOSE) -dir $(MIGRATIONS_DIR) create $(name) sql
