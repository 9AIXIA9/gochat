# Middleware Skill

## Purpose

Enforce the correct ordering and usage of middleware chains for HTTP and Kafka transports in GoChat backend.

## HTTP Middleware Chain

### Global Middleware (applied to all routes)

Applied in `cmd/api/di/providers_http.go` on the Gin engine:

```
RequestID → Recover → Logger → CORS
```

| Position | Middleware | File | Purpose |
|---|---|---|---|
| 1 | `RequestID` | `request_id_middleware.go` | Inject unique request ID into context + response header |
| 2 | `Recover` | `recover_middleware.go` | Catch panics, log stack trace, return 500 |
| 3 | `Logger` | `logger_middleware.go` | Log method, path, status, latency, request ID, trace ID, user ID |
| 4 | `CORS` | `CORS_middleware.go` | Cross-origin headers (configurable) |

### Route Group Middleware (API v1)

Applied on the `/api/v1` group:

```
RateLimit → CircuitBreaker(opt) → Timeout → Auth(per-route)
```

| Position | Middleware | File | Purpose |
|---|---|---|---|
| 1 | `RateLimit` | `rate_limit_middleware.go` | Redis-backed rate limiting via ulule/limiter |
| 2 | `CircuitBreaker` | `circuit_break_middleware.go` | Open circuit on downstream failures (optional, config-driven) |
| 3 | `Timeout` | `timeout_middleware.go` | Enforce request timeout from config |
| 4 | `Auth` | (in authorization context) | JWT Bearer token validation (per-route, not global) |

### OTEL Tracing

When OTEL is enabled, an OTEL middleware wraps the entire router for trace propagation.

### Skip Middleware

`skip_middleware.go` provides a `SkipRoute` trie-based skip list. Routes like `/healthz`, `/readyz`, and `/swagger/*` skip rate limiting, circuit breaking, and auth.

## Kafka Consumer Middleware Chain

Applied per consumer in `cmd/api/di/providers_kafka.go`:

```
RateLimit → Timeout → Recover → Inbox → Logger → Trace(if OTEL) → CircuitBreaker(opt)
```

| Position | Middleware | File | Purpose |
|---|---|---|---|
| 1 | `RateLimit` | `rate_limit_middleware.go` | Rate limit message processing |
| 2 | `Timeout` | `timeout_middleware.go` | Enforce per-message timeout |
| 3 | `Recover` | `recover_middleware.go` | Catch panics in handler |
| 4 | `Inbox` | `inbox_middleware.go` | Idempotency via Redis `SetNX` |
| 5 | `Logger` | `logger_middleware.go` | Log message receipt |
| 6 | `Trace` | `trace_middleware.go` | OTEL span per message (when OTEL enabled) |
| 7 | `CircuitBreaker` | `circuit_break_middleware.go` | Open circuit on persistent failures |

### Kafka Error Middleware Chain

Applied to the error handler:

```
Retry → DeadLetter → CircuitBreaker
```

| Position | Middleware | File | Purpose |
|---|---|---|---|
| 1 | `Retry` | `retry_error_middleware.go` | Exponential backoff retry |
| 2 | `DeadLetter` | `dead_letter_error_middleware.go` | Send to DLQ after retries exhausted |
| 3 | `CircuitBreaker` | `circuit_break_middleware.go` | Prevent cascading failures |

## Middleware Pattern

Every middleware follows the same structural pattern:

### HTTP Middleware

```go
func NewXxxMiddleware(deps...) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Before: setup, check, short-circuit
        c.Next()
        // After: logging, metrics
    }
}
```

### Kafka Middleware

```go
func NewXxxMiddleware() Middleware {
    return func(next Handler) Handler {
        return HandlerFunc(func(ctx context.Context, msg *ckafka.Message) error {
            // Before
            err := next.Handle(ctx, msg)
            // After
            return err
        })
    }
}
```

`Middleware` is defined as `func(next Handler) Handler` in `internal/infrastructure/kafka/handler.go`.

## Common Rules

1. **Recover MUST be before any business middleware** — catch panics early
2. **Logger MUST be after Recover** — log the request even if it panics
3. **Trace MUST be outermost possible** — capture full span duration
4. **Inbox MUST be before handler** — idempotency check before processing
5. **Retry/DeadLetter MUST be after handler** — only retry on handler errors
6. **Never skip middleware in the chain** — the order exists for a reason

## Adding New Middleware

1. Create file in `internal/delivery/http/middleware/` or `internal/delivery/kafka/middleware/`
2. Follow the pattern above (struct + constructor + handler func)
3. If HTTP: add to the router chain in `cmd/api/di/providers_http.go`
4. If Kafka: add to the consumer chain in `cmd/api/di/providers_kafka.go`
5. If it needs to be configurable, add config to `config/app_config.go`
6. Run `make test` to verify

## Skip Routes

Routes that bypass middleware (health, readiness, swagger) are configured via `SkipMiddleware` which uses a `collections.Trie` prefix tree. To add a new skip route, add the path to the trie in `providers_http.go`.
