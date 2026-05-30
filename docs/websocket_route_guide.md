# GoChat WebSocket 与消息同步机制接入指南

这份文档旨在帮助前端/客户端开发者快速接入 GoChat 项目的 WebSocket 通信。本系统采用了**「推拉结合 (Push-Pull)」**的现代 IM 架构，实时通信依赖 WebSocket 极速推送 (Push)，而离线/断网补偿依赖 HTTP 接口 (Pull)。

## 1. 连接地址与鉴权

- **连接地址**：ws://<host>:<port>/api/v1/ws/
- **鉴权方式**：
  - 握手时在 Header 中携带 Authorization: Bearer <access_token>
  - 亦支持通过 Query 参数传递 ?access_token=<access_token>
- **注意事项**：未带 Token、Token 格式错误或无效时，连接会被直接拒绝。

## 2. 核心通信协议 (Envelope)

WebSocket 通信数据被封装在统一的 Envelope (信封) 格式中，系统不再直接处理未经信封包裹的裸业务数据。

### 2.1 客户端上行请求 (UpstreamEnvelope)
客户端下发给服务端的 JSON 指令必须包含以下结构：
`json
{
  "action": "chat.send_private_message", 
  "client_message_id": "req-12345",   // 客户端生成的唯一请求凭证，用于配对回执
  "payload": {
    "recipient_id": "u_1002",
    "content": "你好"
  }
}
`

### 2.2 服务端下行响应 (DownstreamEnvelop)
服务端下发的数据包含两种类型：**操作回执 (ACK)** 和 **主动推送 (Push)**。

**情况 A：操作回执 (ACK)**
主要用于确认客户端方才调用的指令是否落库成功。
**注意：回执的 payload 固定返回 null，客户端仅需通过 client_message_id 和 esult 判断请求是否成功。**
`json
{
  "action": "chat.send_private_message",
  "client_message_id": "req-12345",  // 原样退回客户端的请求凭证
  "result": "success",               // 值为 "success" 或者 "failed"
  "payload": null
}
`

**情况 B：主动推送 (Push)**
当客户端收到新消息或系统事件时，服务端会主动下发。此时下行数据不包含 client_message_id 和 esult。
推送的固定格式为 ction 搭配强关联的 payload，payload 内部的字段根据不同的 ction 会有所区分（例如私聊推 ecipient_id，群聊推 oom_id）。
`json
{
  "action": "chat.push_private_message",
  "payload": {
        // ...根据action有所不同的具体内容
  }
}
`

## 3. 消息机制与可靠性保障 (必读)

为了保证消息**绝对不丢失**，本系统采用游标同步补漏机制：

1. **真实 ID 作为游标**：消息下发推送时的 payload.id 是消息在数据库生成时的绝对唯一 ID（UUIDv7），它具有**时间单调递增性**。客户端需要将收到的各类消息的 id 缓存在本地作为游标 (SyncKey)。
2. **重连主动拉取机制**：由于网络断开或 APP 处于后台极易丢失 WebSocket 推送帧，**客户端在 WebSocket 每次重新连上后，严禁依赖服务端积压全量历史推送**。客户端应当在 WebSocket 握手成功后，立即用本地存储的最新的游标 id，主动调用后端的 HTTP 同步接口补齐差集数据。
3. **独立拉取接口**：目前私聊消息和群聊消息属于不同的业务隔离流，前端需要分别调用对应的 HTTP 接口来完成离线拉取过程：
   - **离线私聊消息拉取接口 (HTTP)**：调用 ListPrivateMessages 接口。
   - **离线群聊消息拉取接口 (HTTP)**：调用对应的 Room 控制器 HTTP 接口。
   *(具体的 HTTP Method、Path 与分页参数请参见详细的 HTTP RESTful API Swagger 接口定义)*

## 4. WebSocket 指令与推送载荷清单

### 4.1 发送私聊消息 (chat.send_private_message)
- **请求 (Upstream)**
  `json
  {
    "action": "chat.send_private_message",
    "client_message_id": "唯一的uuid或自增id",
    "payload": {
      "recipient_id": "目标用户ID",
      "content": "消息内容"
    }
  }
  `
- **操作回执 (Downstream)**
  `json
  {
    "action": "chat.send_private_message",
    "client_message_id": "上行传入的相同ID",
    "result": "success", 
    "payload": null
  }
  `

### 4.2 发送群聊消息 (chat.send_room_message)
- **请求 (Upstream)**
  `json
  {
    "action": "chat.send_room_message",
    "client_message_id": "唯一的uuid或自增id",
    "payload": {
      "room_id": "目标房间ID",
      "content": "群消息内容"
    }
  }
  `
- **操作回执 (Downstream)**
  `json
  {
    "action": "chat.send_room_message",
    "client_message_id": "上行传入的相同ID",
    "result": "success", 
    "payload": null
  }
  `

### 4.3 实时推送事件下发

此类事件属于纯下行通道，无需客户端发起请求，由服务器根据路由系统主动推送下发。

#### 4.3.1 私聊消息下发推送 (chat.push_private_message)
- **下发载荷 (Downstream)**
  `json
  {
    "action": "chat.push_private_message",
    "payload": {
      "id": "服务端生成的绝对消息UUIDv7",
      "sender_id": "发送方UUID",
      "recipient_id": "接受方用户UUID",
      "content": "私聊消息文本内容",
      "sent_at": "服务端消息产生时间"
    }
  }
  `

#### 4.3.2 群聊消息下发推送 (chat.push_room_message)
- **下发载荷 (Downstream)**
  `json
  {
    "action": "chat.push_room_message",
    "payload": {
      "id": "服务端生成的绝对消息UUIDv7",
      "sender_id": "发送方UUID",
      "room_id": "发生聊天的群组/房间UUID",
      "content": "群组消息文本内容",
      "sent_at": "服务端消息产生时间"
    }
  }
  `

## 5. 心跳与防死链策略

- **Ping/Pong**：服务端会持续定期发送底层 ping 帧检测链路，客户端必须确保其 WebSocket 底层引擎支持并自动对 ping 提供相应的 pong 回复。
- **状态维护**：若心跳超时或长时间未收到 pong 返回，服务端将视为死链接并强制切断物理连接。发生截断后，客户端执行前文描述的重连 + HTTP API 差集同步拉取逻辑即可。
