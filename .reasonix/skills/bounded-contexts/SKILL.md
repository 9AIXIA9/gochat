# Bounded Contexts Skill

## Purpose

Define the bounded context map for GoChat backend. Use this skill whenever adding code to determine **which context** it belongs to, how contexts communicate, and what patterns are forbidden.

## What Is a Context in This Project

A **context** is an independently-structured module under `internal/` that:
- Has its own **domain model** (entities, value objects, aggregate roots)
- Has its own **ubiquitous language**
- Exposes **ports** (repository interfaces, service interfaces) to the outside
- Has **adapters** that implement those ports (persistence, transport handlers)
- Communicates with other contexts **only** through defined contracts: events, commands, or shared kernel types

Directories that meet all of these criteria are **contexts**. Directories that don't are **architectural layers** (Clean Architecture port/infrastructure layers) or **shared kernel** — they exist to support contexts, not to model a domain.

---

## Part 1: Domain Contexts

These six contexts model the **IM business domain**. Their ubiquitous language is about users, messages, rooms, friends, and notifications.

| Context | Directory | Aggregate Roots | Responsibility |
|---|---|---|---|
| **authorization** | `internal/authorization/` | `User`, `RefreshToken` | Sign-up, login by email/number, token refresh/parse, auth middleware |
| **profile** | `internal/profile/` | `UserProfile`, `RoomProfile` | User/room profile CRUD, profile queries |
| **friendship** | `internal/friendship/` | `FriendRequest`, `Friendship` | Send/accept/refuse friend requests, list friendships |
| **roomship** | `internal/roomship/` | `Room`, `RoomRequest`, `Roomship` | Create room, join/leave room, member requests, list rooms/members |
| **chat** | `internal/chat/` | `PrivateMessage`, `RoomMessage` | Send/list private messages, send/list room messages |
| **notification** | `internal/notification/` | None (orchestration domain) | Receive `NotificationIntentCreated` events, select notification channel (Gateway push vs email), apply delivery strategy, handle retry/degradation, invoke `GatewayService.PushToUser()` | Generating notification intents (→ other contexts emit the event) |

### Each domain context's NOT-responsible-for

| Context | Does NOT handle |
|---|---|
| **authorization** | User profiles (→ `profile`) |
| **profile** | Login, tokens (→ `authorization`); room membership (→ `roomship`) |
| **friendship** | User identity (→ `authorization`); real-time push (→ `notification`) |
| **roomship** | Room messages (→ `chat`); room metadata display (→ `profile`) |
| **chat** | Friendship validation (→ `friendship`); room membership (→ `roomship`) |
| **notification** | Generating notification intents (→ other contexts emit the event) |

### Aggregate Root Identity Rules

- **authorization.User** has `ID` and `Number` — canonical user identity
- **profile.UserProfile** has the same `UserID` but is a *separate aggregate* — profile data (nickname, avatar, etc.)
- **friendship.Friendship** references `UserID` + `FriendID` — never references profile or auth entities directly
- **roomship.Room** has `RoomID` and `Number` — canonical room identity
- **profile.RoomProfile** has the same `RoomID` — separate aggregate for room metadata

---

## Part 2: Technical Context

**gateway** is a context. Its domain happens to be a technical concern rather than an IM business concern, but structurally it is identical to the six domain contexts above: it has its own domain model, ubiquitous language, ports, and adapters.

| Context | Directory | Aggregate Roots | Ubiquitous Language | Responsibility |
|---|---|---|---|---|
| **gateway** | `internal/gateway/` | `Session`, `Manager` | register, unregister, push, heartbeat, batch, ack, upstream, envelope | WebSocket connection lifecycle, multi-device session pool, real-time push-to-user, upstream command routing |

### Gateway Internal Structure

```
gateway/
├── core/                      # Domain
│   ├── session.go             # Session interface + SessionID
│   ├── manager.go             # Manager: multi-device pool (map[UserID]map[SessionID]Session)
│   ├── envelope.go            # Envelope: wire format {action, client_message_id, payload}
│   ├── ports.go               # SessionIDGenerator interface
│   └── gateway_upstream_handler.go  # UpstreamHandler interface
├── adapter/websocket/         # Port adapter (implements Session with gorilla/websocket)
│   ├── handler.go             # IngressHandler: HTTP→WS upgrade
│   └── session.go             # WSSession: read/write pumps
├── api/local/service.go      # Exposes GatewayService to other contexts
└── infrastructure/uuid/       # SessionIDGenerator implementation
```

### Gateway's Ports (Contracts with Other Contexts)

| Port | Interface | Direction |
|---|---|---|
| **Exposed** | `contract.GatewayService.PushToUser(ctx, userID, envelope)` | Other contexts → Gateway |
| **Consumed** | Receipt handler (Kafka consumer) — receives delivery confirmations | Command bus → Gateway |
| **Produced** | Kafka sync commands (`send_private_message`, `send_room_message`) | Gateway → `chat` |

### Why Gateway Is a Context, Not a Layer

Unlike `delivery/` or `infrastructure/`:
- Gateway has its **own domain model**: `Session`, `Manager`, `Envelope` are not just wrappers around a framework
- Gateway has its **own ubiquitous language**: terms like "register", "push", "upstream" have precise meanings within its boundary
- Gateway defines **ports** (`SessionIDGenerator`, `UpstreamHandler`) and has **adapters** that implement them (`adapter/websocket/`)
- Gateway's domain rules (multi-device management, batching strategy, heartbeat intervals, ack correlation) are non-trivial — they *are* the business logic of a real-time connection layer
- Gateway communicates with other contexts only via defined contracts (Kafka commands, `GatewayService` interface)

---

## Part 3: Context Hierarchy & Dependency Direction

Contexts are organized into three tiers. Dependencies flow downward only — a context may depend on tiers below it, never above.

### Tier 1 — Core Domain

These model the **core business** of the IM system. They are the reason the system exists.

| Context | Rationale |
|---|---|
| **chat** | Core value: message sending and retrieval |
| **friendship** | Core value: social graph |
| **roomship** | Core value: group communication |

**Dependency rule**: Core domain contexts may depend on **shared kernel** and **technical domain** only. They MUST NOT depend on supporting domain contexts.

### Tier 2 — Supporting Domain

These **enable** the core domain but are not the primary business differentiator.

| Context | Rationale |
|---|---|
| **authorization** | Supports: user identity for all core contexts |
| **profile** | Supports: display data for users and rooms |
| **notification** | Supports: delivery of notifications triggered by core contexts |

**Dependency rule**: Supporting domain contexts may depend on **shared kernel** and **technical domain** only. They MUST NOT depend on core domain contexts or other supporting domain contexts.

### Tier 3 — Technical Domain

This provides **infrastructure capabilities** with its own domain model.

| Context | Rationale |
|---|---|
| **gateway** | Technical capability: real-time connection management, multi-device push |

**Dependency rule**: Technical domain contexts may depend on **shared kernel** only. They MUST NOT depend on any domain context (core or supporting).

### Dependency Diagram

```
┌──────────────────────────────────────┐
│            Core Domain                │
│  chat    friendship    roomship      │
│                                       │
│  ↓ allowed: shared, technical domain │
│  ✗ forbidden: supporting domain      │
└────────────┬─────────────────────────┘
             │ (events/commands only)
             ↓
┌──────────────────────────────────────┐
│         Supporting Domain             │
│  authorization  profile  notification│
│                                       │
│  ↓ allowed: shared, technical domain │
│  ✗ forbidden: core domain, other     │
│    supporting domains                 │
└────────────┬─────────────────────────┘
             │ (contract interfaces only)
             ↓
┌──────────────────────────────────────┐
│         Technical Domain              │
│            gateway                    │
│                                       │
│  ↓ allowed: shared kernel only       │
│  ✗ forbidden: any domain context     │
└──────────────────────────────────────┘
```

### Cross-Tier Communication

| From → To | Mechanism | Example |
|---|---|---|
| Core → Core | Kafka events | `friendship` emits `FriendshipCreated` → `chat` consumes |
| Core → Supporting | Kafka events | `chat` emits `NotificationIntentCreated` → `notification` consumes |
| Core → Technical | Contract interface | `chat` uses `GatewayService` (not directly — goes through event → notification → gateway) |
| Supporting → Technical | Contract interface | `notification` calls `GatewayService.PushToUser()` |
| Supporting → Core | ❌ FORBIDDEN | `authorization` MUST NOT depend on `chat` |
| Technical → Any domain | ❌ FORBIDDEN | `gateway` MUST NOT import `chat/domain` or `friendship/domain` |

---

## Part 4: Architectural Layers (Not Contexts)

These four directories are **not contexts**. They are Clean Architecture layers + Shared Kernel pattern. They support contexts but have no domain model of their own.

### `delivery/` — Transport Layer (Clean Architecture "port" aggregation)

`internal/delivery/` aggregates transport handlers that belong to *other* contexts' port layers. It's a packaging convenience, not a context.

```
delivery/
├── http/           # HTTP middleware + health/ready/404 handlers
├── kafka/          # Kafka middleware + receipt handler
├── binlog/         # Binlog reader (triggers outbox → Kafka)
└── gateway/        # UpstreamRouter: bridges WebSocket → command bus (belongs to gateway context logically)
```

**Why not a context**: No domain model. No ports of its own. Just a collection of technical adapters implementing ports defined by contexts.

### `infrastructure/` — Infrastructure Layer (Clean Architecture "infrastructure")

`internal/infrastructure/` contains technical implementations of interfaces defined by contexts' `domain/` packages.

```
infrastructure/
├── gorm/           # MySQL connector + translator
├── redis/          # Redis connector + generic Repository[Model,Domain]
├── kafka/          # Kafka producer/consumer/router + adapters
├── gin/            # HTTP server + AdaptUseCaseToHandler
├── otel/           # OpenTelemetry init
├── zap/            # Logger init
├── viper/          # Config loader
├── bcrypt/         # Hasher
├── breaker/        # Circuit breaker factory
├── ulule/          # Rate limiter factory
├── canal/          # MySQL binlog reader
├── validator/      # go-playground/validator wrapper
├── uuid/           # UUID v7 ID generators
├── persistence/    # Outbox event/DeadLetter models + repository
├── godotenv/       # .env loader
└── metrics/        # OTEL metric instruments
```

**Why not a context**: Pure technical implementations. Implements interfaces defined elsewhere. No domain logic, no ubiquitous language, no ports of its own.

### `shared/` — Shared Kernel (DDD pattern, not a context)

`internal/shared/` is a **Shared Kernel** — types and contracts used by multiple contexts. It's a DDD collaboration pattern, not a bounded context.

```
shared/
├── kernel/     # UseCase[I,O], Validatable, Serializer, ID types, Email, PhoneNumber, etc.
├── api/        # Response, Code enum, pre-built response singletons
├── errors/     # Sentinel errors, BusinessError
├── command/    # Command/StandardCommand/Handler/SyncPublisher interfaces
├── event/      # Event/StandardEvent/Handler/AsyncPublisher/Manager interfaces
├── contract/   # GatewayService interface, NotificationIntentCreatedEvent (shared events)
├── collections/# Trie (prefix tree)
└── timeout/    # Context timeout checker
```

**Why not a context**: No domain model. Types here are *shared by* contexts, not *owned by* any single one. Adding a type to `shared/` is a cross-context design decision, not a feature addition to a context.

**Rule**: Only add a type to `shared/` when **two or more** contexts need it. If only one context uses it, keep it in that context's `domain/`.

### `application/` — Cross-Cutting Use Case

`internal/application/` has a single file: `unpublished_events_created_usecase.go` — the outbox processor.

**Why not a context**: It's one use case that reads from the outbox table and publishes to Kafka. No domain model, no ports, no ubiquitous language. Just a cross-cutting glue component.

---

## Cross-Context Communication Map

### Events (async, 1→N)

| Publisher | Event | Consumers |
|---|---|---|
| `authorization` | `UserCreatedEvent` | `profile`, `chat`, `friendship`, `roomship` |
| `friendship` | `FriendshipCreatedEvent` | `chat` |
| `friendship` | `FriendRequestAgreedEvent` | `roomship` |
| `roomship` | `RoomCreatedEvent` | `profile`, `chat` |
| `roomship` | `RoomshipCreatedEvent` | `profile`, `chat` |
| `chat` | `NotificationIntentCreatedEvent` | `notification` |

### Commands (sync, 1→1, request-reply)

| Sender | Command Action | Handler |
|---|---|---|
| `gateway` | `send_private_message` | `chat` |
| `gateway` | `send_room_message` | `chat` |

### Contract Interfaces (shared/kernel)

| Interface | Defined In | Implemented By | Used By |
|---|---|---|---|
| `GatewayService.PushToUser()` | `shared/contract` | `gateway/api/local/service.go` | `notification` |

---

## Forbidden Patterns

### ❌ Direct import of another context's domain

```go
// ❌ NEVER — friendship domain importing chat domain
import "gochat/internal/chat/domain"
```

Always use events, commands, or shared kernel contracts instead.

### ❌ Treating authorization.User as profile.UserProfile

They are *different aggregates* in different contexts. Authorization owns identity + credentials. Profile owns display data. If `chat` needs user info, it defines its own value object populated from `UserCreatedEvent`.

### ❌ Context logic in architectural layers

```go
// ❌ NEVER — domain logic in delivery layer
func (h *SomeHandler) Handle(c *gin.Context) {
    if user.Age < 18 { return error }  // belongs in domain, not delivery
}
```

`delivery/` only transforms requests and forwards to application layer. `infrastructure/` only implements interfaces. `shared/` only provides shared types. None of them make business decisions.

### ❌ Gateway doing business logic

Gateway's domain is *connection management*. It should NOT contain IM business rules:

```go
// ❌ NEVER — IM business rule in gateway
func (m *Manager) PushToUser(userID, friendID UserID, msg string) {
    if !areFriends(userID, friendID) { return }  // belongs in friendship context!
}
```

Gateway pushes whatever `notification` tells it to push — it doesn't decide *who* should receive *what*.

### ❌ Domain importing frameworks

```go
// ❌ NEVER — framework import in domain
import "github.com/gin-gonic/gin"
import "gorm.io/gorm"
```

All domain packages (`<context>/domain/`, `gateway/core/`, `shared/`) are pure Go. Only `shared/` and standard library.

### ❌ Skipping the outbox for events

Events must go through `event.Repository.CreateUnpublishedEvents()` in the same DB transaction as the business data change. Never publish to Kafka directly from a use case.

### ❌ Defining single-context types in shared/

`shared/` is for types used by multiple contexts. If only one context uses a type, define it in that context's `domain/`.

### ❌ Business rules in the technical domain

```go
// ❌ NEVER — gateway making IM business decisions
func (m *Manager) PushToUser(userID, friendID kernel.UserID, msg string) error {
    if !areFriends(userID, friendID) { return ErrNotFriends }  // belongs in friendship context!
    if !isRoomMember(userID, roomID) { return ErrNotMember }    // belongs in roomship context!
    // ...
}
```

`gateway` pushes whatever it's told to push. It does not decide *who* should receive *what*. Channel selection and delivery policy belong to `notification`. Friendship/membership rules belong to their respective contexts.

### ❌ Reverse dependency across tiers

```go
// ❌ NEVER — technical context importing a domain context
import "gochat/internal/chat/domain"           // gateway MUST NOT import chat

// ❌ NEVER — supporting domain importing core domain
import "gochat/internal/chat/domain"           // authorization MUST NOT import chat

// ❌ NEVER — core domain importing supporting domain
import "gochat/internal/authorization/domain"  // chat MUST NOT import authorization
```

Dependency flows downward: Core → Supporting → Technical → Shared Kernel. Never reverse.

### ❌ Accessing another context's database tables

```go
// ❌ NEVER — chat context querying friendship tables directly
db.Table("friendships").Where("user_id = ?", userID).Find(&friendships)

// ❌ NEVER — profile context querying authorization tables directly
db.Table("users").Where("id = ?", userID).First(&user)
```

Every context owns its database tables. Other contexts access data only through published events (async) or contract interfaces (sync), never through direct DB queries.

---

## Adding a New Context

Whether domain or technical, the structure is the same:

1. Create `internal/<name>/` with standard structure: `core/` or `domain/`, `adapter/` or `infrastructure/persistence/`, `api/` or `port/`
2. Define entities and aggregate roots in the domain layer
3. Define ports (interfaces) in the domain layer
4. Define domain errors in `errors.go`
5. Define domain events in `<name>_event.go`
6. Implement infrastructure adapters
7. Create use cases in `application/`
8. Create transport handlers in `port/http/` or `port/event/`
9. Add providers to `cmd/api/di/`
10. If producing events consumed by other contexts, register consumers in `cmd/api/di/providers_kafka.go`
11. Document the context here — add to Part 1 (domain) or Part 2 (technical) above
