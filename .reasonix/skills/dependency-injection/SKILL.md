# Dependency Injection Skill

## Purpose

Enforce Google Wire conventions for GoChat backend dependency injection.

## Critical Rule

**`wire_gen.go` is auto-generated. NEVER edit it manually.** Changes go in `wire.go` + `providers_*.go`, then run `cd cmd/api/di && go generate`.

## Key Files

| File | Role | Editable |
|---|---|---|
| `cmd/api/di/wire.go` | Wire build definition (`//go:build wireinject`) | ✅ |
| `cmd/api/di/wire_gen.go` | Generated init code | ❌ NEVER |
| `cmd/api/di/providers_*.go` | Provider functions | ✅ |
| `cmd/api/di/dependencies.go` | `Dependencies` struct + `BuildDependencies` | ✅ |

## The 9 Wire Sets

```go
wire.Build(
    InfraSet,        // DB, Redis, Kafka clients, OTel, Logger, Validator, Limiters, Breaker, ID generators
    RepoSet,         // All repository implementations bound to domain interfaces
    KafkaSet,        // Kafka producer + 15 consumers + error handler
    BinlogSet,       // Canal binlog reader + outbox dispatcher
    UseCaseHTTPSet,  // Use cases triggered by HTTP requests (19 use cases)
    UseCaseKafkaSet, // Use cases triggered by Kafka events (15 use cases)
    WebsocketSet,    // Gateway Manager, UpstreamHandler, IngressHandler, GatewayService
    HTTPSet,         // Gin router with all routes + HTTP server wrapper
    BuildDependencies, // Assembles Dependencies struct for lifecycle management
)
```

## Provider Files

| File | Contents |
|---|---|
| `providers_infra.go` | Infrastructure: MySQL, Redis, Kafka publishers, OTel, Logger, Validator, HTTP/Kafka limiters, breaker, all ID generators, type aliases, interface binds |
| `providers_repository.go` | All repository implementations bound to domain interfaces via `wire.Bind` |
| `providers_kafka.go` | Kafka producer, 15 consumers across 5 domains, error retry handler |
| `providers_usecase.go` | UseCaseHTTPSet (19) + UseCaseKafkaSet (15) |
| `providers_http.go` | Gin router with full route tree, HTTP server, readiness checker |
| `providers_websocket.go` | Gateway Manager, UpstreamHandler, allowed actions, GatewayService, IngressHandler |
| `providers_binlog_reader.go` | Canal binlog reader + handler |

## Provider Function Rules

### Naming

```go
func provide<Name>(deps...) (<Type>, error)
```

Prefix always `provide`. CamelCase matching the return type.

### Error Handling

Providers that can fail MUST return `(T, error)`. Wire propagates the error — if any provider fails, `Initialize()` returns that error.

### Interface Binding

Always bind concrete implementations to domain interfaces:

```go
wire.Bind(new(domain.UserRepository), new(*repository.UserRepository))
```

The interface is in `domain/`, the implementation in `infrastructure/persistence/repository/`.

### Constructor Injection

Providers receive dependencies as parameters. Wire resolves them from the graph:

```go
func provideSignUpUseCase(
    eventIDGen event.IDGenerator,
    userIDGen kernel.UserIDGenerator,
    hasher authorizationApp.Hasher,
    userRepo authorizationDomain.UserRepository,
) *application.SignUpUseCase { ... }
```

## Kafka Consumer Structure

`provideKafkaConsumers` builds **15 consumers** across contexts:

| Context | Consumers |
|---|---|
| `profile` | UserCreated, RoomCreated, RoomshipCreated (3) |
| `chat` | UserCreated, RoomCreated, RoomshipCreated, FriendshipCreated, send_private_message command, send_room_message command (6) |
| `gateway` | Receipt handler (1 — handles 4 receipt patterns) |
| `roomship` | UserCreated, RoomCreated, MemberRequestAgreed (3) |
| `friendship` | UserCreated, FriendRequestAgreed (2) |
| `notification` | NotificationIntentCreated (1) |

Each consumer's middleware stack:

```
RateLimit → Timeout → Recover → Inbox(idempotency) → Logger → Trace(if OTEL) → CircuitBreaker(optional)
```

## Dependencies Struct

`cmd/api/di/dependencies.go` holds all long-lived components for startup/shutdown:

```go
type Dependencies struct {
    HttpServer                  *gin.Server
    MysqlDB                     *gorm.DB
    RedisClient                 *redis.Client
    KafkaAsyncEventPublisher    *kafka.EventAsyncPublisher
    CommandReceiptAsyncPublisher *kafka.CommandReceiptAsyncPublisher
    KafkaSyncCommandPublisher   *kafka.CommandSyncPublisher
    OutboxDispatcher            *outbox.Dispatcher
    KafkaConsumers              []*kafka.Consumer
    BinlogReader                *canal.BinlogReader
    OTELShutdown                func(context.Context) error
}
```

## Adding a New Dependency

1. **Write the provider** in the appropriate `providers_*.go`:

```go
func provideNewRepository(db *gorm.DB) *repository.NewRepository {
    return repository.NewNewRepository(db)
}
```

2. **Add to the correct Wire Set** in the same file:

```go
var RepoSet = wire.NewSet(
    // ... existing ...
    provideNewRepository,
    wire.Bind(new(domain.NewRepository), new(*repository.NewRepository)),
)
```

3. **If it needs lifecycle management**, add it to `Dependencies` and wire it in `BuildDependencies`.

4. **Regenerate**: `cd cmd/api/di && go generate`

5. **Verify**: `make build`

## Common Errors

| Error | Cause | Fix |
|---|---|---|
| `no provider found for <type>` | Missing provider | Add provider function returning that type |
| `multiple providers for <type>` | Two providers return same type | Use distinct type or wire set separation |
| Cycle detected | Circular dependency | Introduce interface or restructure |
| `wire_gen.go` edited | Manual edit lost on regenerate | Always edit source files only |
