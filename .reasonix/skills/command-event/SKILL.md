# Command & Event Skill

## Purpose

Enforce the command and event envelope patterns, publishing conventions, and the distinction between synchronous commands and asynchronous events.

## Conceptual Distinction

| Aspect | Event | Command |
|---|---|---|
| **Naming** | Past tense: `UserCreated`, `FriendshipCreated` | Imperative: `send_private_message` |
| **Semantics** | Something that *happened* (fact) | Something to *do* (request) |
| **Ownership** | Produced by one context, consumed by N others | Sent by client, handled by one context |
| **Delivery** | Async (at-least-once via Outbox → Kafka) | Sync (request-reply via Kafka) |
| **File location** | `domain/<name>_event.go` | `domain/<name>_command.go` |
| **Interface** | `event.SpecificEvent` | `command.SpecificCommand` |

## Command Envelope

### Interface (`shared/command/command.go`)

```go
type Command interface {
    ID() ID
    AggregateID() kernel.ID
    Action() Action
    OccurredAt() time.Time
    Payload() []byte
    Headers() map[string]string
    AddHeader(key, value string)
    AddHeaders(headers map[string]string)
}

type SpecificCommand interface {
    Command
    kernel.Serializer  // Marshal() + Unmarshal()
}
```

### StandardCommand

Concrete implementation. Always construct via factory, never manually:

```go
cmd := command.NewStandardCommand(
    kernel.UserID(userID),      // aggregate ID
    command.Action("send_private_message"), // action
    payload,                    // JSON bytes
    idGenerator,                // command.IDGenerator
)
cmd.AddHeader("client_message_id", clientMsgID)
cmd.AddHeader("user_id", userID)
```

Extract from any Command: `command.LoadStandardCommandFromCommand(cmd)`.

### Publishing

**Sync command** (blocking, request-reply for WebSocket):

```go
receipt, err := syncPublisher.Publish(ctx, cmd)
```

Used for WebSocket upstream commands (`send_private_message`, `send_room_message`).

**Async receipt** (fire-and-forget with delivery callback):

```go
asyncPublisher.Publish(ctx, cmd, func(receipt Receipt) {
    // Called on delivery confirmation
})
```

### Receipt

```go
type Receipt interface {
    Status() ReceiptStatus  // "succeeded" or "failed"
    Message() string
    Marshal() ([]byte, error)
    Unmarshal(data []byte) error
}
```

Always check receipt status. `StandardReceipt` is the concrete type.

## Event Envelope

### Interface (`shared/event/event.go`)

```go
type Event interface {
    ID() ID
    AggregateID() kernel.ID
    Topic() Topic
    OccurredAt() time.Time
    Payload() []byte
    Headers() map[string]string
    AddHeader(key, value string)
    AddHeaders(headers map[string]string)
}

type SpecificEvent interface {
    Event
    kernel.Serializer
}
```

### StandardEvent

```go
event := event.NewStandardEvent(
    kernel.UserID(userID),      // aggregate ID
    event.Topic("user.created"), // topic
    payload,                     // JSON bytes
    idGenerator,                 // event.IDGenerator
)
```

Extract: `event.LoadStandardEventFromEvent(evt)`.

### Publishing (always Async)

Events are **never** published directly to Kafka. They go through the Outbox:

```go
// In a use case, within a DB transaction:
tx := db.Begin()
repo.Save(ctx, entity)                              // 1. Persist business data
eventRepo.CreateUnpublishedEvents(ctx, events...)   // 2. Write events to outbox (same TX)
tx.Commit()                                          // 3. Commit atomically
```

The Binlog Reader detects the outbox rows and publishes to Kafka. This guarantees at-least-once delivery.

## Event Manager (In-Memory)

`shared/event/manager.go` provides an in-memory event collector for aggregates:

```go
type Manager struct { events []SpecificEvent }

func (m *Manager) RecordEvent(e SpecificEvent)     // Accumulate
func (m *Manager) GetEvents() []SpecificEvent        // Drain (returns + clears)
```

Usage in entity methods:

```go
func (u *User) ChangeNickname(name string) {
    u.Nickname = name
    u.eventManager.RecordEvent(NewUserProfileUpdatedEvent(u.ID, name))
}
```

The use case drains events from the entity and persists them to the outbox.

## Generic Adapters

### Command UseCase Adapter

```go
command.AdaptUsecaseToCommandHandler(
    useCase,             // kernel.UseCase[Input, Output]
    func(cmd SpecificCommand) Input {
        // Unmarshal payload → input
    },
) command.Handler
```

### Event UseCase Adapter

```go
event.AdaptUsecaseToEventHandler(
    useCase,             // kernel.UseCase[Input, Output]
    func(evt SpecificEvent) Input {
        // Unmarshal payload → input
    },
) event.Handler
```

## Kafka Topics & Actions

### Event Topics

| Topic | Publisher | Consumers |
|---|---|---|
| `user.created` | `authorization` | `profile`, `chat`, `friendship`, `roomship` |
| `friendship.created` | `friendship` | `chat` |
| `friend_request.agreed` | `friendship` | `roomship` |
| `room.created` | `roomship` | `profile`, `chat` |
| `roomship.created` | `roomship` | `profile`, `chat` |
| `notification.intent.created` | `chat` | `notification` |

### Command Actions

| Action | Handler | Transport |
|---|---|---|
| `send_private_message` | `chat` | WebSocket → Kafka Sync |
| `send_room_message` | `chat` | WebSocket → Kafka Sync |

## Idempotency

All event consumers MUST be idempotent. Events may be delivered more than once (at-least-once). The Inbox middleware (`internal/delivery/kafka/middleware/inbox_middleware.go`) provides Redis-based deduplication via `SetNX`.

## Common Mistakes

- ❌ Publishing events directly to Kafka instead of through the Outbox
- ❌ Using events for request-reply (use commands)
- ❌ Not checking receipt status for sync commands
- ❌ Non-idempotent event handlers
- ❌ Forgetting to add `client_message_id` header on commands
