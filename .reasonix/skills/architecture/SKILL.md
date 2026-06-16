# Architecture Skill

## Purpose

Enforce the layered architecture, configuration conventions, database migration rules, and project structure of GoChat backend.

## Layered Architecture

This project follows **Clean Architecture** with strict inward dependency:

```
────────────────────────────────────────────────────────────
port/              Transport handlers (HTTP, Kafka, WebSocket, Binlog)
    ↓ depends on
────────────────────────────────────────────────────────────
application/       UseCase orchestration — workflow, no business rules
    ↓ depends on
────────────────────────────────────────────────────────────
domain/            Entities, value objects, repository interfaces, domain errors
    ↑ implements
────────────────────────────────────────────────────────────
infrastructure/    Technical implementations (GORM, Redis, Kafka, etc.)
────────────────────────────────────────────────────────────
```

### Layer Rules

| Layer | Allowed imports | Forbidden imports |
|---|---|---|
| **domain/** | standard library, `shared/` | Gin, GORM, Kafka, any framework |
| **application/** | `domain/`, `shared/` | frameworks |
| **infrastructure/** | `domain/`, `shared/`, any framework | other context's `domain/` |
| **port/** | `application/`, `shared/`, transport libs | `domain/` (should go through application) |

### Dependency direction

Always inward: `port → application → domain ← infrastructure`

Domain is the center. It defines *what* (interfaces). Infrastructure implements *how*.

### Cross-Context Interaction Rules

Beyond single-context layering, interactions **between** contexts follow additional constraints:

| Rule | Explanation |
|---|---|
| Domain contexts MUST NOT directly import each other's `domain/`, `application/`, or `infrastructure/` | Use events, commands, or shared kernel contracts instead |
| Domain contexts calling the technical context (`gateway`) MUST only use its exposed application layer (`contract.GatewayService`) | Never access gateway internals (`core/`, `adapter/`, `infrastructure/`) |
| Domain contexts calling infrastructure MUST only use shared ports (`shared/command`, `shared/event`, `shared/contract`) | Never import `internal/infrastructure/kafka` or `internal/infrastructure/gorm` directly |
| NO context may directly access another context's database tables | All data access goes through the owning context's repository interface or published events |

## Project Directory Structure

```
backend/
├── cmd/api/                    # Entry point + DI (Wire)
│   ├── main.go                 # Bootstrap: config → DI → start → signal → shutdown
│   └── di/                     # Wire providers, sets, generated code
├── config/                     # YAML config + Go structs
│   ├── config.yaml             # Base defaults (committed)
│   ├── config.override.yaml    # Environment overrides (gitignored)
│   └── config.local.yaml       # Local dev overrides (gitignored)
├── db/migrations/              # Goose SQL migrations
├── deployment/                 # Grafana dashboards, Prometheus config
├── docs/                       # Swagger output + supplementary docs
├── internal/
│   ├── <domain-context>/       # application/ domain/ infrastructure/ port/ dto/
│   ├── delivery/               # Transport adapters (http/, kafka/, binlog/, gateway/)
│   ├── gateway/                # WebSocket session management
│   ├── infrastructure/         # Technical adapters (gorm/, redis/, kafka/, otel/, …)
│   ├── application/            # Cross-cutting use cases (outbox)
│   └── shared/                 # Shared kernel
├── pkg/                        # Reusable utilities (concurrency, ctxutil, httptest, …)
├── scripts/                    # Dev scripts (fix_swagger.py)
├── docker-compose.yml          # Core services
├── docker-compose.observability.yml
├── Dockerfile                  # Multi-stage build
├── Makefile                    # Dev/ops commands
└── go.mod
```

## Config Conventions

### Three-Tier Overlay

```
config.yaml           ← base defaults (committed to git)
    ↓ overlaid by
config.override.yaml  ← environment overrides (gitignored, optional)
    ↓ overlaid by
config.local.yaml     ← local dev overrides (gitignored, optional)
```

Loaded via Viper with `mapstructure` tags. Entry point: `viper.LoadConfigFile(basePath, overlayPaths...)` → returns `*config.App`.

### Config Struct Pattern

Every config struct follows this pattern (`config/app_config.go`):

```go
type App struct {
    Name    string        `mapstructure:"Name"`
    Timeout time.Duration `mapstructure:"Timeout"`
    Mysql   *gorm.MysqlConfig `mapstructure:"Mysql"`
    // ... nested configs
}

func (c *App) Validate() error {
    if c == nil { return ErrEmptyPointer }
    // Apply defaults for zero-values
    if c.Timeout == 0 { c.Timeout = defaultTimeout }
    // Validate each field
    if c.Name == "" { return fmt.Errorf("App.Name: %w", ErrEmptyInput) }
    // Delegate to nested configs
    if err := c.Mysql.Validate(); err != nil { return fmt.Errorf("App.Mysql: %w", err) }
    return nil
}
```

### Config Files to Never Commit

- `.env` — secrets (passwords, API keys)
- `config/config.override.yaml` — environment-specific
- `config/config.local.yaml` — local dev overrides

Template files (`.env.example`, `config.*.example.yaml`) are committed.

## Database Migrations

### Tool

Goose (`github.com/pressly/goose`). Binary at `goose` or via Docker Compose `migrate` service.

### Naming Convention

```
YYYYMMDDHHMMSS_descriptive_snake_case.sql
```

Example: `20260527075103_refactor_chat_message_tables.sql`

### File Structure

```sql
-- +goose Up
-- SQL statements for forward migration
ALTER TABLE ...;

-- +goose Down
-- SQL statements for rollback
ALTER TABLE ...;
```

Every migration MUST have both `Up` and `Down`. Use `CREATE TABLE IF NOT EXISTS` for idempotency in init migrations.

### Running Migrations

| Scenario | Command |
|---|---|
| Local dev | `make migrate-up` |
| Create new | `make migrate-create name=add_xxx_table` |
| Status check | `make migrate-status` |
| Docker env | `make compose-migrate` |

## Standard Context Directory Structure

Every context under `internal/<context>/` follows:

```
internal/<context>/
├── application/              # UseCases — one file per use case
│   ├── <usecase>.go
│   └── <usecase>_test.go
├── domain/                   # Business model
│   ├── <entity>.go           # Entity with behavior + LoadXxx() factory
│   ├── <entity>_test.go
│   ├── errors.go             # Domain-specific sentinel errors
│   ├── ports.go              # Repository + external service interfaces (optional)
│   ├── <event>_event.go      # Domain events
│   └── mocks/                # Generated mocks (go:generate mockgen)
├── infrastructure/
│   └── persistence/
│       ├── model/            # GORM models
│       │   └── <entity>.go
│       ├── converter/        # Entity ↔ Model (if separate from repository)
│       └── repository/
│           └── <entity>_repository.go
├── port/
│   ├── http/                 # HTTP handlers
│   │   ├── <handler>.go
│   │   └── <handler>_test.go
│   ├── event/                # Kafka event handlers (consumers)
│   │   └── <event>_event_handler.go
│   └── command/              # Kafka command handlers
│       └── <command>_command_handler.go
└── dto/                      # Data Transfer Objects (optional)
    └── <dto>.go
```

## Naming Conventions

| Element | Pattern | Example |
|---|---|---|
| UseCase file | `<verb>_<noun>_usecase.go` | `send_friend_request_usecase.go` |
| UseCase type | `<Verb><Noun>UseCase` | `SendFriendRequestUseCase` |
| Handler file | `<verb>_<noun>_handler.go` | `send_friend_request_handler.go` |
| Handler type | `<Verb><Noun>Handler` | `SendFriendRequestHandler` |
| Repository interface | `<Entity>Repository` | `FriendRequestRepository` |
| Event file | `<noun>_<verb>_event.go` | `friend_request_sent_event.go` |
| Mock file | `mock_<interface>.go` | `mock_friend_request_repository.go` |
| Test file | `<source>_test.go` | `friend_request_test.go` |
| Migration file | `YYYYMMDDHHMMSS_<desc>.sql` | `20260527075103_refactor_chat.sql` |
| Provider function | `provide<Name>` | `provideMysqlConnection` |

## Graceful Shutdown Pattern

`cmd/api/main.go` demonstrates the standard pattern:

1. `signal.NotifyContext` for SIGINT/SIGTERM
2. `<-sigCtx.Done()` blocks until signal
3. `context.WithTimeout` for shutdown deadline
4. Components closed in reverse dependency order:
   - HTTP server → Binlog reader → Outbox dispatcher → Kafka consumers → Kafka publishers → MySQL → Redis → OTel

## Makefile Targets

| Target | Action |
|---|---|
| `make fmt` | `go fmt ./...` |
| `make test` | `go test ./...` |
| `make test-race` | `go test -race ./...` |
| `make lint` | `golangci-lint run ./...` |
| `make build` | `go build -o gochat-app ./cmd/api` |
| `make run` | `go run ./cmd/api` |
| `make tidy` | `go mod tidy` |
| `make swagger` | `swag init -g cmd/api/main.go --parseDependency --parseInternal` |
| `make swagger-fix` | Fix swagger example fields |
| `make up` | Docker Compose: mysql, redis, kafka, migrate, app |
| `make up-obs` | Same + otel-collector, jaeger, prometheus, grafana |
| `make down` | `docker compose down --remove-orphans` |

## Architecture Anti-Patterns — Forbidden Practices

These patterns are **never allowed** in any context. Violations will be rejected in code review.

### ❌ Framework dependencies in domain layer

```go
// ❌ NEVER
import "github.com/gin-gonic/gin"
import "gorm.io/gorm"
import "github.com/confluentinc/confluent-kafka-go/v2/kafka"
```

Domain packages (`<context>/domain/`, `gateway/core/`, `shared/`) import only standard library and `shared/`. No frameworks.

### ❌ Business decisions in transport layer (port/delivery)

```go
// ❌ NEVER — business rule in HTTP handler
func (h *Handler) Handle(c *gin.Context) {
    if user.Age < 18 { return error }  // belongs in domain, not handler
    if amount > limit { return error } // belongs in domain, not handler
}
```

Transport layer only transforms requests → input, calls use case, transforms output → response. Zero business logic.

### ❌ Bypassing the use case layer

```go
// ❌ NEVER — transport layer calling repository directly
func (h *Handler) Handle(c *gin.Context) {
    user, err := h.userRepo.FindByID(ctx, id)  // skip use case!
    room, err := h.roomRepo.FindByID(ctx, id)  // skip use case!
}

// ✅ CORRECT — always go through use case
func (h *Handler) Handle(c *gin.Context) {
    output, err := h.useCase.Execute(ctx, input)
}
```

Transport → UseCase → Repository. Never skip the middle layer. Use cases own transaction boundaries.

### ❌ Anemic domain entities

```go
// ❌ NEVER — entity with no behavior, just data bag
type FriendRequest struct {
    ID     string
    Status string  // modified externally
}

// ✅ CORRECT — entity encapsulates behavior
type FriendRequest struct {
    id     kernel.ID
    status FriendRequestStatus
}

func (fr *FriendRequest) Accept() error {
    if fr.status != Pending { return ErrAlreadyHandled }
    fr.status = Accepted
    return nil
}
```

Domain entities must carry their own business behavior. Don't push logic to services or use cases.

### ❌ Business logic in infrastructure layer

```go
// ❌ NEVER — infrastructure making business decisions
func (r *UserRepository) Save(ctx context.Context, user *User) error {
    if user.Age < 18 { return ErrTooYoung }  // business rule in repository!
    return r.db.Create(r.toModel(user)).Error
}

// ✅ CORRECT — repository only persists
func (r *UserRepository) Save(ctx context.Context, user *User) error {
    return r.db.Create(r.toModel(user)).Error
}
```

Infrastructure implements interfaces. It does not enforce business rules.

