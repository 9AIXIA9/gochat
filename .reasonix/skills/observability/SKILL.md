# Observability Skill

## Purpose

Enforce observability conventions: OpenTelemetry initialization, metrics recording, structured logging with Zap, and safe goroutine patterns.

## OpenTelemetry

### Initialization (`internal/infrastructure/otel/observability.go`)

```go
shutdown, err := otel.Init(config)
```

- Sets up global `TracerProvider`, `MeterProvider`, and propagator
- Returns a `shutdown` function — MUST be called during graceful shutdown
- Pings endpoint before initializing
- Config is optional — only enabled when `OTEL.Enabled = true`

### Configuration (`otel.Config`)

```go
type Config struct {
    Enabled              bool
    Endpoint             string        // OTel collector endpoint
    ServiceName          string
    Environment          string
    Insecure             bool
    TraceSampleRatio     float64       // [0,1]
    MetricExportInterval time.Duration
    MetricExportTimeout  time.Duration
}
```

### Shutdown

Always defer shutdown in main:

```go
dependencies, err := di.Initialize(conf)
// OTELShutdown is part of Dependencies
// Called during graceful shutdown:
dependencies.OTELShutdown(ctx)
```

## Metrics (`internal/infrastructure/metrics/metrics.go`)

### Pattern: Lazy Init + Function-per-Metric

Each metric group uses a `sync.Once` for lazy initialization:

```go
var (
    wsOnce sync.Once
    wsConnectionDelta otel.Int64UpDownCounter
)

func initWSMetrics() {
    meter := otel.GetMeterProvider().Meter("gochat/websocket")
    wsConnectionDelta, _ = meter.Int64UpDownCounter("ws.connection.delta")
}

func RecordWSConnectionDelta(delta int64) {
    wsOnce.Do(initWSMetrics)
    wsConnectionDelta.Add(context.Background(), delta)
}
```

### Metric Groups

| Group | Metrics | When recorded |
|---|---|---|
| WebSocket | `ws.connection.delta`, `ws.disconnect`, `ws.message.received`, `ws.message.sent`, `ws.send.blocked`, `ws.write.latency`, `ws.read.error`, `ws.write.error`, `ws.session.count` | Connection lifecycle, message flow |
| Kafka | `kafka.produce.latency`, `kafka.produce.error`, `kafka.consume.latency`, `kafka.consume.error`, `kafka.consume.count`, `kafka.consumer.pause` | Produce/consume operations |
| Binlog | `binlog.reader.run` | Reader start/success/failure |
| Outbox | `outbox.dispatcher.run`, `outbox.dispatcher.event_count`, `outbox.dispatcher.error` | Dispatcher sweep cycles |

### Adding a New Metric

1. Add a `sync.Once` + init function for the metric group
2. Define the metric instrument
3. Create a public `RecordXxx()` function
4. Call it from the relevant code path

## Structured Logging (Zap)

### Initialization (`internal/infrastructure/zap/logger.go`)

```go
zaputils.Initialize(config)
```

- JSON encoder for structured log output
- Replaces global zap logger
- Config: `LoggerConfig{ Level string }` with `Validate()`

### Usage

Always use the global zap logger with structured fields:

```go
zap.L().Info("user logged in",
    zap.String("user_id", string(userID)),
    zap.String("method", "email"),
)

zap.L().Error("handler failed",
    zap.String("path", c.FullPath()),
    zap.String("request_id", requestID),
    zap.Error(err),
)
```

### Rules

- Use `zap.L()` (global) — never create local loggers
- Use structured fields (`zap.String`, `zap.Int`, `zap.Error`), not `fmt.Sprintf`
- `Info` for normal operations, `Warn` for recoverable issues, `Error` for failures
- Always include `user_id` and `request_id`/`trace_id` when available

## Context Propagation (`pkg/ctxutil`)

### Store values in context

```go
ctx = ctxutil.WithUserID(ctx, userID)
ctx = ctxutil.WithRequestID(ctx, requestID)
```

### Retrieve values

```go
userID := ctxutil.UserIDFrom(ctx)
requestID := ctxutil.RequestIDFrom(ctx)
spanID, traceID := ctxutil.SpanIDAndTraceIDFrom(ctx)
```

### Headers

```go
headers := ctxutil.HeadersFrom(ctx)  // map[string]string
```

Used for propagating `client_message_id`, `user_id` across Kafka messages.

## Safe Goroutines (`pkg/concurrency`)

### GoSafe

```go
concurrency.GoSafe(func() {
    // any panic is recovered and logged
    doWork()
})
```

Never use raw `go func()` for long-running goroutines. Always wrap with `GoSafe` to prevent unhandled panics from crashing the process.

### Pattern

```go
// ✅ Correct
concurrency.GoSafe(func() {
    d.runLoop()
})

// ❌ Wrong
go func() {
    d.runLoop()  // panic kills the process
}()
```

## Retry with Backoff (`pkg/retry`)

```go
for retry := 0; retry < maxRetries; retry++ {
    if err := doSomething(); err == nil {
        return nil
    }
    if err := retry.Wait(ctx, retry, 100*time.Millisecond); err != nil {
        return err  // context cancelled or deadline exceeded
    }
}
```

`retry.Wait` implements exponential backoff with jitter: `baseDelay * 2^retry + random(baseDelay * 2^retry / 2)`.

## Timeout Check (`shared/timeout`)

```go
if timeout.CheckCtxTimeout(ctx) {
    return ErrTimeout
}
```

Non-blocking context deadline check. Used in loops to avoid unnecessary work when context has expired.
