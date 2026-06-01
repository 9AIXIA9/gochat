# GoChat WebSocket 接入指南

本文档面向前端/客户端工程师，说明如何对接 GoChat 后端的 WebSocket 层：连接地址、鉴权、Envelope（信封）格式、常用
action、心跳与重连策略、以及实现细节与示例代码。

快速要点

- WebSocket URL：GET /api/v1/ws/（在 Gin 路由中注册，且启用鉴权中间件）
- 握手鉴权：推荐使用 HTTP Header：Authorization: Bearer <access_token>；也支持 ?access_token=<access_token> 查询参数。
- Envelope（JSON）：action, client_message_id（可选）, payload（可选）
- 代码中常见的 action 字符串：
	- 上行（客户端发往服务端）：`send_private_message`, `send_room_message`
	- 下行主动推送（服务端发往客户端）：`push_notification`
	- 协议/控制类事件：`push_command_ack`, `push_command_error`, `push_command_receipt`
- 心跳（ping/pong）：服务端周期性发送 ping（客户端需回复 pong）。见 session.go 中的常量：writeWait=5s, pongWait=2m,
  pingPeriod=30s。

---

1. 连接与鉴权

- 示例连接地址（开发环境）：

```
ws://your-host:8080/api/v1/ws/
```

- 握手鉴权示例：

Header:

```
Authorization: Bearer <access_token>
```

或使用查询参数（不推荐放在 URL 日志中）：

```
ws://your-host:8080/api/v1/ws/?access_token=<access_token>
```

注意：鉴权通过中间件在握手阶段校验，未通过会拒绝升级。

2. Envelope（信封）格式

后端定义了统一的 Envelope，用于在 WebSocket 上承载所有业务/协议消息。序列化为 JSON 后的字段：

- `action`（string，必填）：命令或事件名。
- `client_message_id`（string，可选）：客户端生成的请求 ID，用于匹配回执/关联业务。
- `payload`（object / null，可选）：与 action 相关的业务数据（任意 JSON）。

示例（通用）：

```json
{
  "action": "<action_name>",
  "client_message_id": "<uuid-or-client-id>",
  "payload": {}
}
```

实现细节：服务端在内部使用 `core.Envelope`（见 `internal/gateway/core/envelope.go`），客户端应以上述 JSON 作为读写契约。

3. 常见消息类型与示例

3.1 上行：发送私聊消息（示例）

客户端 -> 服务端

```json
{
  "action": "send_private_message",
  "client_message_id": "550e8400-e29b-41d4-a716-446655440000",
  "payload": {
    "recipient_id": "019e7d33-e41f-7d46-9efc-da399b10cc7f",
    "content": "你好"
  }
}
```

说明：

- `action` 使用仓库当前的短名（非 namespaced）。如果团队希望使用 `chat.send_private_message` 之类的 namespaced
  名称，需要在后端/前端统一改名或做映射适配。

3.2 快速回执（ACK）

gateway 在收到并成功将命令同步投递到事件总线后，会返回快速回执：

```json
{
  "action": "send_private_message",
  "client_message_id": "550e8400-e29b-41d4-a716-446655440000",
  "result": "success",
  "payload": null
}
```

- `result` 为 `success` 或 `failed`。注意：ACK 表示命令已成功投递/接收（fast ACK），并不等同于业务最终完成；某些业务最终状态可能通过后续事件或
  receipt 下发。

3.3 服务端主动推送（业务通知）

当前 notification domain 使用统一 `action = "push_notification"`，payload 根据具体场景区分（私聊/群聊/系统通知等）。

私聊推送示例：

```json
{
  "action": "push_notification",
  "payload": {
    "id": "019e7d33-e41f-7d46-9efc-da399b10cc7f",
    "sender_id": "019e7905-1601-784e-90b8-a8c894ac32db",
    "recipient_id": "019e7d33-e41f-7d46-9efc-da399b10cc7f",
    "content": "hello",
    "sent_at": "2026-05-31T12:00:00Z"
  }
}
```

群聊推送示例：

```json
{
  "action": "push_notification",
  "payload": {
    "id": "019e7d33-e438-7c30-befb-ba63e22bc4c3",
    "sender_id": "019e7905-1601-784e-90b8-a8c894ac32db",
    "room_id": "room-019e7d33-...",
    "content": "group hello",
    "sent_at": "2026-05-31T12:00:00Z"
  }
}
```

说明：推送中的 `id` 为服务端生成的业务消息 ID（UUIDv7），客户端应把它视作游标（SyncKey）并持久化以便重连补漏使用。

3.4 协议/运维类事件

- `push_command_ack`：用于某些链路返回命令 ack（带回 `client_message_id`）。
- `push_command_error`：当上行命令校验失败或投递出错时返回（带 `client_message_id` 和 error 信息）。
- `push_command_receipt`：命令处理链路中的回执会被转换为下发事件（见
  `internal/delivery/kafka/handler/receipt_handler.go`）。

4. 客户端实现建议（必须遵循的关键点）

4.1 心跳与超时

- 基于实现，服务端会定期发 ping；客户端必须能够回复 pong（底层实现通常自动处理）。参考 `session.go` 常量：
	- writeWait = 5s
	- pongWait = 2m
	- pingPeriod = 30s

4.2 接收多条消息（Batching）注意

- gateway 的写入端在发送时会尽可能地把当前队列中的消息合并到同一帧并在消息间插入换行符 (`\n`) 以支持
  batch。客户端在接收事件时需能拆分按行解析每条 JSON（benchmark 客户端实现了类似逻辑）。

4.3 重连后的差量拉取（关键）

- 客户端在每次 WebSocket 连接建立成功后，**必须**使用本地保存的最大已接收业务消息 ID（游标，payload.id）调用后端 HTTP
  同步接口拉取断线期间的离线消息。不要依赖服务端通过 WebSocket 自动回补全部历史。
- 目前私聊/群聊的离线拉取是通过 REST API 单独实现（参见 Swagger），客户端应分别调用对应接口（ListPrivateMessages /
  ListRoomMessages），并以游标/分页参数拉取差集。

4.4 发送策略与 ACK 处理

- 发送上行命令时，带上 `client_message_id` 并在应用层实现超时/重试策略。收到 `result: success` 的 ACK 表示命令已被接受并写入处理链路；如果收到
  `push_command_error`，应该标记为失败并给用户提示或重试策略。

5. 简短示例（JavaScript）

```js
const ws = new WebSocket('ws://localhost:8080/api/v1/ws/?access_token=' + token);

ws.onopen = () => {
    // 连接建立后，先用 HTTP 接口拉取离线消息（以本地游标为起点），然后继续处理实时推送
};

ws.onmessage = (evt) => {
    // 单帧内可能包含多条 JSON，用换行分隔
    const parts = evt.data.split('\n');
    for (const p of parts) {
        if (!p) continue;
        const env = JSON.parse(p);
        // 处理 env.action / env.client_message_id / env.payload
    }
};

function sendPrivate(recipientId, content) {
    const env = {
        action: 'send_private_message',
        client_message_id: generateUUIDv7(),
        payload: {recipient_id: recipientId, content},
    };
    ws.send(JSON.stringify(env));
}
```