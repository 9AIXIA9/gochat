# WebSocket Skill

## Purpose

Enforce the WebSocket gateway architecture: session lifecycle, connection management, message routing, and batching.

## Architecture

```
┌─────────────────────────────────────────────────┐
│ gateway/                                        │
│ ├── adapter/websocket/    Gorilla WS adapter    │
│ │   ├── handler.go         HTTP→WS upgrade      │
│ │   └── session.go         Read/write pumps     │
│ ├── api/local/service.go   GatewayService impl  │
│ ├── core/                  Pure domain          │
│ │   ├── manager.go         Multi-device pool    │
│ │   ├── session.go         Session interface    │
│ │   ├── envelope.go        Wire format          │
│ │   └── ports.go           IDGenerator          │
│ └── infrastructure/uuid/   SessionIDGenerator   │
├─────────────────────────────────────────────────┤
│ delivery/gateway/                               │
│ └── upstream_handler.go    Bridge: WS→Command   │
└─────────────────────────────────────────────────┘
```

## Connection Lifecycle

### 1. Handshake (`adapter/websocket/handler.go`)

```
Client → GET /api/v1/ws → Auth middleware → IngressHandler.ServeHTTP → Upgrade
```

- Auth extracts `UserID` from JWT (Bearer header or `?access_token=` query param)
- Gorilla WebSocket upgrader with configurable `CheckOrigin`
- Upgrade happens AFTER auth validation

### 2. Session Creation (`adapter/websocket/session.go`)

```go
session := NewWSSession(sessionID, userID, conn, hub, upstreamHandler)
```

- `send` channel: buffered (256)
- `conn`: gorilla WebSocket connection
- `hub`: core.Manager reference

### 3. Registration (`core/manager.go`)

```go
manager.Register(session)
// Stored as: sessions[userID][sessionID] — supports multi-device
```

### 4. Read Pump (goroutine)

- Reads with 4KB limit per frame
- Pong deadline: 2 minutes
- Every raw message → `upstreamHandler.HandleUpstream(ctx, userID, message)`
- Returns ACK bytes pushed to `send` channel
- On read error → `hub.Unregister(s)` → `conn.Close()`

### 5. Write Pump (goroutine)

- Drains `send` channel
- **Batching**: reads ALL queued messages, joins with `\n`, sends as single frame
- Sends WebSocket pings every 30s
- On write error → close

### 6. Unregister

```
hub.Unregister(s) → remove from sessions[userID][sessionID] → cleanup empty user entries
```

## Envelope Format

Messages use a common JSON envelope (`core/envelope.go`):

```json
{
    "client_message_id": "uuid-v7",
    "action": "send_private_message",
    "payload": {}
}
```

### Upstream (Client → Server)

Actions: `send_private_message`, `send_room_message` (whitelist in `providers_websocket.go`).

### Downstream (Server → Client)

Messages pushed from the server:
- `push_command_ack` — command accepted by Kafka
- `push_command_error` — command failed
- `push_command_receipt` — delivery confirmation
- `push_notification` — real-time notification from other users

## Upstream Handler (`delivery/gateway/upstream_handler.go`)

Bridges WebSocket messages to the command bus:

1. Parse envelope → validate `client_message_id` + `action`
2. Check action against whitelist
3. Create `command.NewStandardCommand(userID, action, payload, idGenerator)`
4. Attach headers: `client_message_id`, `user_id`
5. **Synchronous publish** to Kafka via `command.SyncPublisher.Publish(ctx, cmd)`
6. Return `push_command_ack` or `push_command_error` envelope to WebSocket session

## Push to User (`gateway/api/local/service.go`)

Implements `contract.GatewayService`:

```go
func (s *LocalGatewayService) PushToUser(ctx context.Context, userID kernel.UserID, envelope *core.Envelope) error {
    sessions := s.manager.GetByUserID(userID)
    if len(sessions) == 0 {
        return ErrNotFound  // user offline
    }
    for _, session := range sessions {
        data, _ := json.Marshal(envelope)
        session.Send(data)  // non-blocking push to send channel
    }
    return nil
}
```

## Multi-Device Support

- One user can have multiple WebSocket connections (phone + desktop)
- `Manager` stores `map[UserID]map[SessionID]Session`
- `PushToUser` delivers to ALL sessions for a user
- `Send()` is non-blocking — on full channel, message is dropped (logged as warning)

## Heartbeat

- Server → Client: WebSocket ping every 30s (write pump)
- Client → Server: WebSocket pong (read pump, 2min deadline)
- No custom heartbeat messages — uses standard WebSocket ping/pong

## Batching

The write pump batches messages to reduce syscalls:

```go
// Read ALL queued messages at once, join with \n, send as single frame
```

Messages are newline-delimited within a single WebSocket text frame.

## Adding a New WebSocket Action

1. Define the action string constant
2. Add to allowed actions whitelist in `cmd/api/di/providers_websocket.go`
3. Create a command handler in the target context's `port/command/`
4. Wire the handler into the Kafka consumer in `providers_kafka.go`
5. Run `cd cmd/api/di && go generate`

## Common Mistakes

- ❌ Blocking in the read pump (keep upstream handler fast)
- ❌ Long operations in `HandleUpstream` (use async if needed)
- ❌ Forgetting to add `client_message_id` header (breaks ACK correlation)
- ❌ Direct DB access from gateway (gateway has no business logic)
- ❌ Assuming single-device per user (always iterate sessions)
