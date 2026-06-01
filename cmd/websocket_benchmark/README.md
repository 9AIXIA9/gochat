# WebSocket 压测脚本说明

本目录包含 GoChat 的 WebSocket 压测与预热脚本，用于：

- 批量创建/登录测试账号，并导出 token / user_id
- 预热好友关系与房间关系
- 对 WebSocket 连接进行并发压测、消息发送压测、ACK/端到端延迟采样

- WebSocket 地址：`/api/v1/ws/`
- 上行 action：`send_private_message`、`send_room_message`
- 下行推送 action：`push_notification`
- 鉴权：`Authorization: Bearer <access_token>`
- 登录接口：`/auth/login/user_number`（`account_manager` 使用该路由登录并导出 token）

---

## 目录

- `account_manager/`：批量注册登录账号，并导出压测所需的 token / recipients
- `friendship_warmup/`：预热好友关系，生成好友配对文件
- `room_warmup/`：预热房间与成员关系，生成房间与拓扑文件
- `ws_runner/`：WebSocket 压测主程序

---

## 1. 准备压测数据

建议先在仓库根目录执行账号准备，再做关系预热，最后跑 WebSocket 压测。

### 1.1 账号准备：`account_manager`

作用：
- 批量注册测试账号
- 登录获取 token
- 导出用户 ID 列表供后续脚本使用

示例（PowerShell）：

```powershell
cd E:\App\GoChat\backend
cd .\cmd\websocket_benchmark\account_manager

go run . `
  --base-url http://localhost:8080/api/v1 `
  --count 100 `
  --workers 50 `
  --state-file ..\..\..\websocket_benchmark_data\state\accounts.json `
  --tokens-file ..\..\..\websocket_benchmark_data\output\tokens.txt `
  --recipients-file ..\..\..\websocket_benchmark_data\output\recipients.txt
```

### 1.2 好友预热：`friendship_warmup`

作用：
- 生成好友配对关系
- 输出 `friendship_pairs.jsonl`，供 `ws_runner` 的 private 模式使用

示例：

```powershell
cd E:\App\GoChat\backend
cd .\cmd\websocket_benchmark\friendship_warmup

go run . `
  --base-url http://localhost:8080/api/v1 `
  --tokens-file ..\..\..\websocket_benchmark_data\output\tokens.txt `
  --recipients-file ..\..\..\websocket_benchmark_data\output\recipients.txt `
  --out-map-file ..\..\..\websocket_benchmark_data\output\friendship_pairs.jsonl
```

### 1.3 房间预热：`room_warmup`

作用：
- 批量创建房间
- 邀请成员入群
- 输出 `rooms.txt` 和 `room_topology.jsonl`

示例：

```powershell
cd E:\App\GoChat\backend
cd .\cmd\websocket_benchmark\room_warmup

go run . `
  --base-url http://localhost:8080/api/v1 `
  --tokens-file ..\..\..\websocket_benchmark_data\output\tokens.txt `
  --recipients-file ..\..\..\websocket_benchmark_data\output\recipients.txt `
  --rooms-file ..\..\..\websocket_benchmark_data\output\rooms.txt `
  --out-map-file ..\..\..\websocket_benchmark_data\output\room_topology.jsonl
```

---

## 2. WebSocket 压测主程序：`ws_runner`

### 2.1 作用

`ws_runner` 会：
- 并发建立 WebSocket 连接
- 按配置发送私聊/群聊消息
- 统计连接成功率、发送成功率、ACK 延迟、端到端延迟
- 将结果输出为 JSONL，便于后续分析

### 2.2 输入文件

`ws_runner` 默认读取以下文件：

- `websocket_benchmark_data/output/tokens.txt`
- `websocket_benchmark_data/output/recipients.txt`
- `websocket_benchmark_data/output/rooms.txt`
- `websocket_benchmark_data/output/friendship_pairs.jsonl`（可选，private/mixed 模式推荐）

### 2.3 常用运行示例

Private 模式：

```powershell
cd E:\App\GoChat\backend
cd .\cmd\websocket_benchmark\ws_runner

go run . `
  --url ws://localhost:8080/api/v1/ws/ `
  --tokens-file ..\..\..\websocket_benchmark_data\output\tokens.txt `
  --recipients-file ..\..\..\websocket_benchmark_data\output\recipients.txt `
  --pairs-file ..\..\..\websocket_benchmark_data\output\friendship_pairs.jsonl `
  --send-mode private `
  --clients 200 `
  --duration 60s `
  --send-interval 1500ms
```

Room 模式：

```powershell
cd E:\App\GoChat\backend
cd .\cmd\websocket_benchmark\ws_runner

go run . `
  --url ws://localhost:8080/api/v1/ws/ `
  --tokens-file ..\..\..\websocket_benchmark_data\output\tokens.txt `
  --recipients-file ..\..\..\websocket_benchmark_data\output\recipients.txt `
  --rooms-file ..\..\..\websocket_benchmark_data\output\rooms.txt `
  --send-mode room `
  --clients 200 `
  --duration 60s `
  --send-interval 1500ms
```

Mixed 模式：

```powershell
cd E:\App\GoChat\backend
cd .\cmd\websocket_benchmark\ws_runner

go run . `
  --url ws://localhost:8080/api/v1/ws/ `
  --tokens-file ..\..\..\websocket_benchmark_data\output\tokens.txt `
  --recipients-file ..\..\..\websocket_benchmark_data\output\recipients.txt `
  --rooms-file ..\..\..\websocket_benchmark_data\output\rooms.txt `
  --pairs-file ..\..\..\websocket_benchmark_data\output\friendship_pairs.jsonl `
  --send-mode mixed `
  --room-ratio 0.5 `
  --clients 200 `
  --duration 60s
```

### 2.4 重要参数说明

- `--url`：WebSocket 地址，默认 `ws://localhost:8080/api/v1/ws/`
- `--send-mode`：发送模式，支持 `private` / `room` / `mixed`
- `--clients`：并发连接数
- `--duration`：压测时长
- `--send-interval`：每个客户端发送间隔
- `--ack-timeout`：等待 ACK 的超时时间
- `--origin`：可选 Origin 头
- `--json-out-file`：总体指标输出 JSONL
- `--latency-out-file`：ACK 延迟采样输出 JSONL
- `--e2e-latency-out-file`：端到端延迟采样输出 JSONL

### 2.5 协议对齐说明

`ws_runner` 已按当前后端协议对齐：

- 上行命令使用 `send_private_message` / `send_room_message`
- 主动推送识别 `push_notification`
- ACK / Error 使用 `push_command_ack` / `push_command_error`（并兼容旧的 result 风格）
- 发送时会附带 `client_message_id`
- 读取时支持处理单帧中多条 JSON（按换行分隔）
---

## 3. 常见排查

- 连接失败：确认 `--url`、`Authorization`、服务端是否已启动
- 登录失败：确认认证服务使用的是 `/auth/login/user_number`（不是旧的 `/auth/login`）
- ACK 不匹配：确认客户端生成的 `client_message_id` 是否唯一，且服务端是否返回了同一 ID
- 推送统计为空：确认测试用户是否确实收到了 `push_notification`
- 数据文件缺失：先运行 `account_manager`，再运行 `friendship_warmup` / `room_warmup`，最后运行 `ws_runner`

---

## 4. 推荐执行顺序

1. `account_manager`
2. `friendship_warmup`
3. `room_warmup`
4. `ws_runner`

这个顺序可以保证压测所需的 token、recipients、rooms、pairs 文件都已生成。

