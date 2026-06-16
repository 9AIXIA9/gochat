# AGENTS.md — AI Coding Agent Instructions

> This document provides project conventions for AI coding agents (Reasonix, Claude Code, GitHub Copilot, etc.).
> Detailed specs live under [`.reasonix/skills/`](.reasonix/skills/).

---

## Project Overview

GoChat Backend — Go instant-messaging backend using **Clean Architecture + DDD** with bounded contexts. Multi-transport (HTTP / WebSocket / Kafka), event-driven via Outbox + Binlog + Kafka.

## Quick Command Reference

All commands run from repo root.

| Operation | Command |
|---|---|
| Format | `make fmt` or `go fmt ./...` |
| Test | `make test` or `go test ./...` |
| Race test | `make test-race` or `go test -race ./...` |
| Lint | `make lint` or `golangci-lint run ./...` |
| Build | `make build` or `go build -o gochat-app ./cmd/api` |
| Run | `make run` or `go run ./cmd/api` |
| Tidy | `make tidy` or `go mod tidy` |
| Swagger | `make swagger` |
| Wire DI | `cd cmd/api/di && go generate` |
| Compose up | `make up` |
| Compose down | `make down` |
| DB migrate | `make migrate-up` or `make compose-migrate` |

## Tech Stack

Go 1.24 · Gin · Google Wire · Viper + dotenv · MySQL 8 + GORM · Redis 7 · Kafka (confluent-kafka-go) · Gorilla WebSocket · ulule/limiter · sony/gobreaker · OpenTelemetry · Zap · Swaggo · Goose · testify

## Architecture at a Glance

```
port (transport)        → HTTP / Kafka / WebSocket / Binlog handlers
application (use cases) → UseCase orchestration (no business logic)
domain (business)       → Entities, value objects, repository interfaces, domain errors
infrastructure (adapters)→ DB / cache / message queue implementations
```

**Core rule**: dependency always inward — `port → application → domain ← infrastructure`.

`internal/` is split into **domain contexts** (authorization, chat, friendship, roomship, notification, profile) and **technical contexts** (delivery, gateway, infrastructure, shared, application). See [bounded-contexts](.reasonix/skills/bounded-contexts/SKILL.md).

## Adding a Feature — Checklist

1. Identify the bounded context (or create one) — see [bounded-contexts](.reasonix/skills/bounded-contexts/SKILL.md)
2. Define entity + repository interface + domain errors in `domain/`
3. Write use case + tests in `application/`
4. Implement repository in `infrastructure/persistence/` (model → converter → repository)
5. Write handler in `port/http/` or `port/event/`
6. Add provider in `cmd/api/di/`, wire into the correct Set
7. Run `make swagger`, `make test`, `make lint`

## Skill Index

| Skill | Purpose |
|---|---|
| [bounded-contexts](.reasonix/skills/bounded-contexts/SKILL.md) | Domain vs technical context map, boundaries, cross-context rules |
| [architecture](.reasonix/skills/architecture/SKILL.md) | Clean Architecture layers, config, migrations, project structure |
| [shared-kernel](.reasonix/skills/shared-kernel/SKILL.md) | Generic `UseCase[I,O]`, `Validatable`, value objects, API types, `BusinessError` |
| [dependency-injection](.reasonix/skills/dependency-injection/SKILL.md) | Wire Sets, provider conventions, interface binding |
| [command-event](.reasonix/skills/command-event/SKILL.md) | `StandardCommand`/`StandardEvent`, Receipt, Sync/Async publishers |
| [repository](.reasonix/skills/repository/SKILL.md) | model/converter/repository triad, `LoadXxx()` factory, cursor pagination |
| [middleware](.reasonix/skills/middleware/SKILL.md) | HTTP and Kafka middleware chains, ordering, skip routes |
| [http-handler](.reasonix/skills/http-handler/SKILL.md) | `AdaptUseCaseToHandler` generic adapter, request binding, responses |
| [websocket](.reasonix/skills/websocket/SKILL.md) | Gateway sessions, read/write pumps, batching, heartbeat, upstream routing |
| [testing](.reasonix/skills/testing/SKILL.md) | testify, mockgen, `pkg/httptest`, table-driven, per-layer strategy |
| [observability](.reasonix/skills/observability/SKILL.md) | OTel init/shutdown, metrics lazy-init, Zap, `ctxutil`, `GoSafe` |
