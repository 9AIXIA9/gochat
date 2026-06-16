# Shared Kernel Skill

## Purpose

Enforce correct use of the shared kernel types in `internal/shared/`. These types are the foundation of every use case, entity, and handler in the project.

## Package Map

| Package | Purpose | Key Types |
|---|---|---|
| `shared/kernel` | Generic interfaces + value objects | `UseCase[I,O]`, `Validatable`, `Serializer`, `ID`, `UserID`, `Email`, `PhoneNumber`, `Gender`, `NoInput`, `NoOutput` |
| `shared/api` | HTTP response types | `Response`, `Code`, pre-built response singletons |
| `shared/errors` | Error infrastructure | sentinel errors, `BusinessError` |
| `shared/command` | Command envelope + handler | `Command`, `StandardCommand`, `Handler`, `SyncPublisher`, `Receipt` |
| `shared/event` | Event envelope + handler | `Event`, `StandardEvent`, `Handler`, `AsyncPublisher`, `Manager` |
| `shared/contract` | Cross-context contracts | `GatewayService`, `NotificationIntentCreatedEvent` |
| `shared/collections` | Data structures | `Trie` (prefix tree for route matching) |
| `shared/timeout` | Context utilities | `CheckCtxTimeout` |

## Generic UseCase Interface

Defined in `shared/kernel/generic_usecase.go`:

```go
type UseCase[Input Validatable, Output any] interface {
    Execute(ctx context.Context, input Input) (Output, error)
}
```

This is the contract for **every** use case in the project. Input must satisfy `Validatable`. If no input is needed, use `NoInput`. If no output, use `NoOutput`.

```go
type NoInput struct{}
func (*NoInput) Validate() error { return nil }

type NoOutput struct{}
```

### Usage Pattern

```go
// Define input
type SendFriendRequestInput struct {
    FromUserID kernel.UserID
    ToUserID   kernel.UserID
    Message    string
}
func (i SendFriendRequestInput) Validate() error {
    if i.FromUserID == "" || i.ToUserID == "" { return ErrEmptyInput }
    return nil
}

// Use case implements the generic interface
type SendFriendRequestUseCase struct { ... }
func (uc *SendFriendRequestUseCase) Execute(
    ctx context.Context, input SendFriendRequestInput,
) (SendFriendRequestOutput, error) { ... }

// Compile-time check
var _ kernel.UseCase[SendFriendRequestInput, SendFriendRequestOutput] = (*SendFriendRequestUseCase)(nil)
```

## Validatable Interface

```go
type Validatable interface {
    Validate() error
}
```

Used as the constraint on `UseCase` input. Also used by config structs (`config.App.Validate()`). Every input DTO must implement this.

## Value Objects

All value objects are defined as named string types with `Validate()` methods:

| Type | Underlying | Package | Validation |
|---|---|---|---|
| `ID` | `string` | `kernel` | None (generic) |
| `UserID` | `string` | `kernel` | Type alias of `ID` |
| `RoomID` | `string` | `kernel` | Type alias of `ID` |
| `MessageID` | `string` | `kernel` | Type alias of `ID` |
| `OperationID` | `string` | `kernel` | Type alias of `ID` |
| `Number` | `string` | `kernel` | Digits only |
| `UserNumber` | `string` | `kernel` | Digits only |
| `Email` | `string` | `kernel` | RFC 5322 regex |
| `PhoneNumber` | `string` | `kernel` | Non-empty |
| `Gender` | `int` | `kernel` | One of: 0 (unset), 1 (male), 2 (female), 3 (other) |
| `Address` | `string` | `kernel` | Non-empty |

### Usage

```go
// ✅ Correct — typed value objects
userID := kernel.UserID(uuid.New().String())
email, err := kernel.NewEmail("user@example.com")

// ❌ Wrong — raw strings
userID := "some-uuid"
```

## Serializer Interface

```go
type Serializer interface {
    Marshal() ([]byte, error)
    Unmarshal(data []byte) error
}
```

Used by `SpecificCommand` and `SpecificEvent` for payload serialization. Also used by `StandardCommand`/`StandardEvent`.

## API Response Types

### Response Struct

```go
type Response struct {
    Code    Code   `json:"code"`
    Message string `json:"message"`
    Data    any    `json:"data,omitempty"`
}
```

### Business Code

```go
const (
    CodeSuccess       Code = 200
    CodeBusinessError Code = 300
    CodeServerError   Code = 500
    CodeTimeout       Code = 400 + iota  // 400
    CodeInvalidParam                     // 401
    CodeUnauthorized                     // 402
    CodeNotFound                         // 403
    CodeServiceUnavailable               // 404
    CodeLimitExceeded                    // 405
)
```

Each code has `ToHTTPCode()` and `DefaultMessage()`. Pre-built response singletons exist in `api/consts.go`.

### Response Helpers

Use `ginutils.Response*` functions (from `internal/infrastructure/gin/response.go`), not raw `gin.Context.JSON`:

```go
ginutils.Response(ctx, api.CodeSuccess)
ginutils.ResponseSuccessWithData(ctx, output)
ginutils.ResponseWithMessage(ctx, api.CodeBusinessError, "friend request already sent")
```

## BusinessError

Defined in `shared/errors/business_error.go`:

```go
type BusinessError struct {
    Message string
    Cause   error
}
```

### Creating

```go
// With formatted message
err := errors.NewBusiness("friend request from %s to %s: %w", fromID, toID, domain.ErrAlreadySent)

// Wrapping an existing error
err := errors.WrapBusiness(domain.ErrAlreadySent, "failed to send friend request")
```

### Checking

```go
if errors.IsBusinessError(err) {
    // This is a business error — return to client with message
}
```

### Sentinel Errors

Organized in `shared/errors/errors.go` in three categories:

| Category | Variable prefix | Examples |
|---|---|---|
| System | `Err*` | `ErrServerBusy`, `ErrWrongEventTopic`, `ErrServiceUnavailable` |
| Logic | `Err*` | `ErrEmptyInput`, `ErrEmptyPointer`, `ErrInvalidNumber`, `ErrInvalidLength` |
| Database | `Err*` | `ErrNotFound`, `ErrDuplicatedKey`, `ErrForeignKeyViolated`, `ErrAlreadyExists` |

## Command & Event Base Types

See [command-event skill](../command-event/SKILL.md) for the full envelope patterns. In brief:

- `command.Command` — interface for all commands (ID, AggregateID, Action, OccurredAt, Payload, Headers)
- `command.StandardCommand` — concrete implementation with `NewStandardCommand()` and `LoadStandardCommandFromCommand()`
- `event.Event` — interface for all events (ID, AggregateID, Topic, OccurredAt, Payload, Headers)
- `event.StandardEvent` — concrete implementation with `NewStandardEvent()` and `LoadStandardEventFromEvent()`
- `event.Manager` — in-memory event collector: `RecordEvent()` + `GetEvents()` (drain pattern)

## Contract Interface

`shared/contract/gateway_service.go`:

```go
type GatewayService interface {
    PushToUser(ctx context.Context, userID kernel.UserID, envelop *core.Envelope) error
}
```

Used by `notification` context to push real-time messages to WebSocket clients. Implemented by `gateway/api/local/service.go`.

## Generic Adapters

### AdaptUsecaseToCommandHandler

Converts a `kernel.UseCase[Input, Output]` into a `command.Handler`:

```go
func AdaptUsecaseToCommandHandler[
    SpecificCmd command.SpecificCommand,
    Input kernel.Validatable,
    Output any,
](
    useCase kernel.UseCase[Input, Output],
    convert func(SpecificCmd) Input,
) command.Handler
```

### AdaptUsecaseToEventHandler

Converts a `kernel.UseCase[Input, Output]` into an `event.Handler`:

```go
func AdaptUsecaseToEventHandler[
    SpecificEvt event.SpecificEvent,
    Input kernel.Validatable,
    Output any,
](
    useCase kernel.UseCase[Input, Output],
    convert func(SpecificEvt) Input,
) event.Handler
```

## When to Add to Shared Kernel

Add a type to `shared/` only when **two or more domain contexts** need it. If only one context uses it, keep it in that context's `domain/`.

```go
// ✅ Good — used by multiple contexts
// shared/kernel/id.go
type UserID string

// ❌ Bad — only used by friendship
// shared/kernel/friend_request_status.go  ← belongs in friendship/domain/
```
