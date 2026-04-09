# GoChat WebSocket 路由使用指南

这份文档旨在帮助前端开发者快速理解和使用 GoChat 项目的 WebSocket 路由接口。我们将介绍连接方式、消息格式、可用接口以及一些调试建议。

## 1. 连接地址与鉴权

- WebSocket URL：`ws://<host>:<port>/api/v1/ws/`
- 例如本地：`ws://localhost:8080/api/v1/ws/`
- 鉴权方式：必须在握手请求头中携带 `Authorization: Bearer <access_token>`

未带 Token、Token 格式错误或无效时，握手会被拒绝。

## 2. 收发消息通用格式

客户端发送（请求格式）：

~~~json
{
  "topic": "chat.send_private_message",
  "payload": {
    "recipient_id": "u_1002",
    "content": "你好" 
  }
}
~~~

服务端返回（统一响应包在 body 中）：

~~~json
{
  "topic": "chat.send_private_message",
  "body": {
    "code": 200,
    "message": "success"
  }
}
~~~

说明：

- `topic`：路由标识，决定走哪个 WebSocket handler。
- `payload`：请求参数。
- `body`：服务端响应体，结构与 HTTP 统一：`code/message/data`。

## 3. 全部 WebSocket 接口清单

项目当前一共 4 个 WebSocket 相关 topic：

- 客户端可主动调用（2 个）
- 服务端主动推送（2 个）

### 3.1 chat.send_private_message（客户端发送私聊消息）

请求：

~~~json
{
  "topic": "chat.send_private_message",
  "payload": {
    "recipient_id": "u_1002",
    "content": "你好，今晚有空吗？"
  }
}
~~~

字段：

- `recipient_id`：接收方用户 ID，必填。
- `content`：消息内容，必填，最大 1000 字符。

成功响应：

~~~json
{
  "topic": "chat.send_private_message",
  "body": {
    "code": 200,
    "message": "success"
  }
}
~~~

常见失败：

- 参数缺失/格式错误：`code=401`（invalid parameter）
- 非好友发送私聊：`code=200`，`message` 为业务错误描述（例如 not friends）
- 服务内部异常：`code=500`

### 3.2 chat.send_room_message（客户端发送群聊消息）

请求：

~~~json
{
  "topic": "chat.send_room_message",
  "payload": {
    "room_id": "r_9001",
    "content": "大家晚上好"
  }
}
~~~

字段：

- `room_id`：房间 ID，必填。
- `content`：消息内容，必填，最大 1000 字符。

成功响应：

~~~json
{
  "topic": "chat.send_room_message",
  "body": {
    "code": 200,
    "message": "success"
  }
}
~~~

常见失败：

- 参数缺失/格式错误：`code=401`
- 不是房间成员、房间不存在：`code=200`，`message` 为业务错误描述
- 服务内部异常：`code=500`

### 3.3 chat.notify_private_message（服务端私聊消息推送）

这是服务端主动发给在线接收方用户的通知，客户端不能主动调用该 topic。

推送示例：

~~~json
{
  "topic": "chat.notify_private_message",
  "body": {
    "id": "m_123",
    "sender_id": "u_1001",
    "state": "sent",
    "content": "你好",
    "sent_at": "2026-04-09T10:00:00Z"
  }
}
~~~

### 3.4 chat.notify_room_message（服务端群聊消息推送）

这是服务端主动广播给房间在线成员（不含发送者）的通知，客户端不能主动调用该 topic。

推送示例：

~~~json
{
  "topic": "chat.notify_room_message",
  "body": {
    "id": "m_456",
    "sender_id": "u_1001",
    "room_id": "r_9001",
    "states": {
      "u_1002": "sent",
      "u_1003": "sent"
    },
    "content": "开会了",
    "sent_at": "2026-04-09T10:01:00Z"
  }
}
~~~

## 4. 未定义 topic 的行为

如果你发送了不存在的 topic（例如 `chat.xxx`），服务端会返回：

~~~json
{
  "topic": "chat.xxx",
  "body": {
    "code": 403,
    "message": "not found"
  }
}
~~~

## 5. 心跳与连接稳定性

- 服务端会定期发送 ping 帧。
- 客户端应正确响应 pong（大部分 WebSocket 客户端库会自动处理）。
- 若长时间无 pong，服务端会断开连接并清理会话。

## 6. 调试建议

- 先用 HTTP 登录拿到 Access Token，再连 WebSocket。
- 首次联调建议先发 `chat.send_private_message`，更容易验证。
- 收到 `code=500` 时优先查看后端日志定位具体错误。
