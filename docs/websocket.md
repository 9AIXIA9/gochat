# GoChat WebSocket 接口文档

本文件描述 GoChat 项目的 WebSocket 接口，包括：

- 连接方式与鉴权
- 通用消息格式
- 所有已实现的 WebSocket topic（与 HTTP 路由类似，都有清晰的请求结构）
	- `chat.send_private_message`
	- `chat.send_room_message`

> 说明：本说明严格根据 `internal/chat/port/websocket` 目录下的 handler
> 代码编写，请求字段与后端校验保持一致；响应部分目前均为“无显式回包，有错误则关闭连接或发送错误消息”，你可以根据后续实现再补充服务端推送格式。

---

## 1. WebSocket 连接信息

- **HTTP 升级路由**：`GET /api/v1/ws/`
- **完整 URL 示例**：
	- 本地开发：`ws://localhost:8080/api/v1/ws/`
- **HTTP 方法**：`GET`
- **鉴权方式**：
	- 使用与 HTTP API 相同的 Bearer Token：

	  ```http
	  Authorization: Bearer <access_token>
	  ```

	- 具体携带方式（Header / Cookie）与后端 `AuthorizationMiddleware` 一致，一般推荐 Header。

前端示例（以浏览器为例）：

```js
const ws = new WebSocket("ws://localhost:8080/api/v1/ws/");

ws.onopen = () => {
  console.log("WS connected");
};

ws.onmessage = (event) => {
  const msg = JSON.parse(event.data);
  console.log("WS message", msg);
};

ws.onclose = () => {
  console.log("WS closed");
};
```

> 注意：当前后端从 `context` 中通过 `utils.GetUserID(ctx)` 读取用户 ID，因此必须保证连接握手时已经通过 HTTP Header/Cookie
> 完成鉴权。

---

## 2. 通用消息格式

WebSocket 连接建立后，前端与后端约定使用 JSON 文本进行通信。推荐的统一结构为：

```jsonc
{
  "topic": "chat.send_private_message",
  "payload": {
    // 与 topic 对应的请求体字段
  }
}
```

- `topic`：字符串，表示业务主题，对应后端注册的 `websocket.Topic`：
	- `chat.send_private_message`
	- `chat.send_room_message`
- `payload`：对象，对应各 topic 的请求结构（见下文）。

前端发送时可以统一封装，如：

```js
function sendWS(ws, topic, payload) {
  ws.send(JSON.stringify({ topic, payload }));
}
```

---

## 3. Topic 详情

### 3.1 发送私聊消息：`chat.send_private_message`

- **方向**：前端 → 后端
- **后端 handler**：`internal/chat/port/websocket/send_private_message_handler.go`
- **Topic 常量**：

  ```go
  const SendPrivateMessageTopic websocket.Topic = "chat.send_private_message"
  ```

- **请求 payload 结构**（对应 `SendPrivateMessageData`）：

  ```jsonc
  {
    "recipient_id": 123,          // 必填，对方用户 ID
    "content": "你好，这是消息"   // 必填，最大长度 1000 字符
  }
  ```

	- `recipient_id`：`kernel.UserID`，必填，后端校验 `validate:"required"`
	- `content`：`string`，必填，`validate:"required,max=1000"`

- **服务端逻辑概述**：
	1. 从 payload 反序列化 `SendPrivateMessageData`。
	2. 从 `context` 中取出当前登录用户 ID（`SenderID =ctxutil.UserIDFrom`）。
	3. 构造 `SendPrivateMessageInput`：

	   ```go
	   &application.SendPrivateMessageInput{
		 SenderID:   ctxutil.UserIDFrom,
		 RecipientID: reqData.RecipientID,
		 Content:     reqData.Content,
	   }
	   ```

	4. 调用 `input.Validate()`，再执行 `uc.Execute(ctx, input)`。

- **返回/推送**：
	- 当前 handler 返回 `([]byte, error)` 中的数据部分为 `nil`，不直接通过该调用回包；
	- 实际消息送达到对端通常由 usecase / 其他通道完成（例如通过 HTTP 轮询或后端主动广播）。
	- 前端应以“发送成功/失败”来处理该请求：
		- 无错误时视为发送成功；
		- 若有错误，后端可能通过错误 topic 或关闭连接来表现（建议你后续在服务端统一错误格式）。

- **前端示例**：

  ```js
  sendWS(ws, "chat.send_private_message", {
    recipient_id: 123,
    content: "你好，这是私聊消息"
  });
  ```

---

### 3.2 发送房间消息：`chat.send_room_message`

- **方向**：前端 → 后端
- **后端 handler**：`internal/chat/port/websocket/send_room_message_handler.go`
- **Topic 常量**：

  ```go
  const SendRoomMessageTopic websocket.Topic = "chat.send_room_message"
  ```

- **请求 payload 结构**（对应 `SendRoomMessageData`）：

  ```jsonc
  {
    "room_id": 456,                // 必填，房间 ID
    "content": "大家好，这是房间消息" // 必填，最大长度 1000 字符
  }
  ```

	- `room_id`：`kernel.RoomID`，必填，`validate:"required"`
	- `content`：`string`，必填，`validate:"required,max=1000"`

- **服务端逻辑概述**：

	1. 从 payload 反序列化 `SendRoomMessageData`。
	2. 从 `context` 中取出当前登录用户 ID：`SenderID =ctxutil.UserIDFrom`。
	3. 构造 `SendRoomMessageInput`：

	   ```go
	   &application.SendRoomMessageInput{
		 SenderID:ctxutil.UserIDFrom,
		 RoomID:   reqData.RoomID,
		 Content:  reqData.Content,
	   }
	   ```

	4. 校验并通过 usecase 执行业务逻辑（校验是否为房间成员等）。

- **返回/推送**：

	- 同私聊发送，handler 返回数据为 `nil`，不直接回包；
	- 实际消息广播给房间成员由 usecase 或其他通道完成；
	- 推荐在后端为房间广播另行设计推送消息结构（例如统一的 `topic: "chat.room_message"` 的下行通知）。

- **前端示例**：

  ```js
  sendWS(ws, "chat.send_room_message", {
    room_id: 456,
    content: "大家好，这是房间里的消息"
  });
  ```

---

## 4. 错误处理建议

当前 handler 统一签名：

```go
func(ctx context.Context, data []byte) ([]byte, error)
```

- 如果返回 `error != nil`，上层 WebSocket 框架会根据实现选择：
	- 关闭连接；
	- 或包装成错误消息下发。
- 建议你在 `internal/delivery/websocket/handler` 中统一约定一个错误消息格式，例如：

```jsonc
{
  "topic": "error",
  "payload": {
    "code": 40001,
    "message": "invalid payload: content is empty"
  }
}
```

前端可以统一监听：

```js
ws.onmessage = (event) => {
  const msg = JSON.parse(event.data);
  if (msg.topic === "error") {
    // 统一错误提示
    console.error("WS error", msg.payload.code, msg.payload.message);
    return;
  }

  // 其他业务处理
};
```

---

## 5. 与 HTTP API 的关系

- HTTP API：
	- 登录 / 注册
	- 拉取历史消息
	- 获取资料 / 好友 / 房间列表
- WebSocket：
	- 实时发送/接收消息

推荐前端集成顺序：

1. 使用 HTTP 的 `/authorization/login` 获取 `access_token`。
2. 使用 HTTP 的 `/chat/private`、`/chat/room` 拉取一段历史消息。
3. 建立 WebSocket 连接 `ws://host:port/api/v1/ws/`，携带同样的鉴权信息。
4. 使用本文件定义的 topic+payload 格式，完成实时聊天能力。

