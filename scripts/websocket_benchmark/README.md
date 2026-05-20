# WebSocket Benchmark 使用说明

这个脚本用于对 GoChat 的 WebSocket 服务做压测，主要覆盖三类场景：

1. 连接建立能力
2. 心跳保活能力
3. 消息发送与接收吞吐、端到端延迟

脚本入口是 [websocket.go](websocket.go)，可以直接用 `go run` 运行。

## 运行前准备

先确认后端服务已经启动，并且 WebSocket 地址可访问。默认地址是：

```bash
ws://localhost:8080/api/v1/ws/
```

如果你的前端或网关要求 `Origin`，可以通过 `-origin` 传入，例如：

```bash
-origin http://localhost:5173
```

### 鉴权方式

脚本必须提供以下三种方式中的一种：

1. `-tokens-file`
2. `-auth`
3. `-cookie`

三者不是同时使用的。

### 多用户 token 文件

如果要模拟多用户压测，建议先生成 token 列表。仓库里提供了 [generate_token_pool.ps1](generate_token_pool.ps1)：

```powershell
pwsh .\scripts\websocket_benchmark\generate_token_pool.ps1 -BaseUrl http://localhost:8080/api/v1 -Count 100
```

默认会生成：

- `scripts/websocket_benchmark/tokens.txt`
- `scripts/websocket_benchmark/token_pool.jsonl`
- `scripts/websocket_benchmark/recipients.txt`

PowerShell 版默认是续写模式（append），不会清空已有文件，方便失败后补量重跑。
如果你需要从空文件重新生成，请增加 `-Overwrite`。

其中 `tokens.txt` 每行一个 token，格式可以是 `Bearer <token>` 或裸 token，脚本都会自动处理。
`recipients.txt` 每行一个 `userID`，来自 token 对应的 JWT `UserID` claim，可直接用于 `-recipients-file`。

### 大规模生成（推荐 Go 并发版）

当你要生成 10 万级别 token 时，建议使用 Go 版并发生成器 [token_pool_generator/main.go](token_pool_generator/main.go)。

先用 dry-run 验证参数：

```bash
go run ./scripts/websocket_benchmark/token_pool_generator -dry-run -count 1000 -workers 200
```

再执行真实生成：

```bash
go run ./scripts/websocket_benchmark/token_pool_generator ^
	-base-url http://localhost:8080/api/v1 ^
	-count 100000 ^
	-workers 300 ^
	-max-retries 2 ^
	-retry-backoff 250ms ^
	-timeout 12s ^
	-append=true ^
	-out-file scripts/websocket_benchmark/tokens.txt ^
	-detail-file scripts/websocket_benchmark/token_pool.jsonl ^
	-recipient-file scripts/websocket_benchmark/recipients.txt
```

调参建议：

- `-workers`：并发度，建议从 100 到 300 逐步上调
- `-timeout`：单请求超时，服务端较慢时适当调大
- `-max-retries`：失败重试次数，网络抖动时建议至少 `2`
- `-progress-interval`：进度打印周期，便于观察实时成功/失败数
- `-append`：是否续写输出文件，默认 `true`；设置 `-append=false` 时会覆盖原文件

## 最小可运行示例

### 只压测连接

```bash
go run ./scripts/websocket_benchmark -auth "Bearer <token>" -clients 100 -duration 60s
```

### 连接 + 心跳

```bash
go run ./scripts/websocket_benchmark -auth "Bearer <token>" -clients 500 -duration 120s -ping-interval 5s
```

### 连接 + 发消息

```bash
go run ./scripts/websocket_benchmark -auth "Bearer <token>" -clients 200 -duration 60s -send-interval 2s
```

## 私聊压测

如果要测试私聊消息，使用 `-send-mode private`。

### 单个接收人

```bash
go run ./scripts/websocket_benchmark ^
	-tokens-file .\scripts\websocket_benchmark\tokens.txt ^
	-send-mode private ^
	-recipient-id 10001 ^
	-send-interval 2s ^
	-clients 100 ^
	-duration 60s
```

### 多个接收人轮询

```bash
go run ./scripts/websocket_benchmark ^
	-tokens-file .\scripts\websocket_benchmark\tokens.txt ^
	-send-mode private ^
	-recipients-file .\scripts\websocket_benchmark\recipients.txt ^
	-send-interval 2s ^
	-clients 100 ^
	-duration 60s
```

注意：`-send-mode private` 时，必须提供 `-recipient-id` 或 `-recipients-file`。

## 好友预热（Paired 私聊前置）

Paired 私聊压测前，发送端和接收端需要先建立好友关系。可以使用预热脚本：

```powershell
go run .\scripts\websocket_benchmark\friendship_warmup `
  -base-url http://127.0.0.1:8080/api/v1 `
  -tokens-file .\scripts\websocket_benchmark\tokens.txt `
  -recipients-file .\scripts\websocket_benchmark\recipients.txt `
  -pairs 50 `
  -workers 20
```

默认配对规则与 paired mode 一致：前 N 行为接收端，后 N 行为发送端，发送端 i 对应接收端 i。

## 常用参数

- `-url`：WebSocket 地址，默认 `ws://localhost:8080/api/v1/ws/`
- `-clients`：总连接数
- `-duration`：压测总时长
- `-connect-rate`：每秒创建多少个连接，`0` 表示一次性发起
- `-connect-timeout`：单次握手超时
- `-ping-interval`：心跳间隔，`0` 表示关闭
- `-send-interval`：发消息间隔，`0` 表示不发消息
- `-send-mode`：`raw` 或 `private`
- `-payload`：`raw` 模式发送的文本内容模板
- `-content`：`private` 模式里消息内容模板
- `-latency-field`：从收到的消息中提取时间戳字段名，默认 `bench_ts_ns`
- `-report-interval`：压测期间的周期性进度输出，`0` 表示关闭
- `-origin`：握手时附带的 `Origin`
- `-auth`：统一的 `Authorization` 值，脚本会自动补 `Bearer`
- `-cookie`：统一的 `Cookie` 值
- `-tokens-file`：每行一个 token，用于多用户压测
- `-recipient-id`：private 模式下的单个接收人
- `-recipients-file`：private 模式下的接收人列表
- `-chaos-drop-ratio`：随机断开连接的比例，范围 `0~1`
- `-chaos-drop-after`：连接后多久开始触发随机断开

## 模板占位符

`-payload` 和 `-content` 支持以下占位符：

- `{{ts_unix_nano}}`
- `{{client_id}}`
- `{{seq}}`

例如：

```bash
-payload '{"type":"benchmark","bench_ts_ns":"{{ts_unix_nano}}","bench_client":"{{client_id}}","bench_seq":"{{seq}}","content":"ping"}'
```

`private` 模式下，脚本会自动组装成：

```json
{
	"topic": "chat.send_private_message",
	"payload": {
		"recipient_id": "...",
		"content": "..."
	}
}
```

## 延迟统计

脚本会尝试从收到的消息里提取时间戳，并计算端到端延迟。

默认字段名是 `bench_ts_ns`，如果服务端返回的字段不同，可以通过 `-latency-field` 修改。

支持两种格式：

1. JSON 对象中的同名字段
2. 文本中的键值对，例如 `bench_ts_ns=123456789`

## 结果输出

结束后会输出一组汇总指标，包括：

- 连接成功率
- 峰值在线数
- 保留连接数
- 消息发送/接收吞吐
- 消息发送成功率
- 端到端延迟的平均值和 P50/P95/P99
- 心跳成功/失败数

如果开启了 `-report-interval`，压测过程中还会周期性打印进度。

## 使用注意

- `-tokens-file` 和 `-auth` / `-cookie` 不能同时使用
- `-tokens-file` 中的 token 数量必须不少于 `-clients`
- `-send-mode private` 必须配合 `-recipient-id` 或 `-recipients-file`
- 如果设置了 `-chaos-drop-ratio`，必须同时设置大于 `0` 的 `-chaos-drop-after`

## 建议的压测流程

1. 先用 10 到 50 个连接确认鉴权和 WebSocket 地址正确
2. 再逐步提高 `-clients` 和 `-connect-rate`
3. 最后开启 `-send-interval`、`-send-mode private` 或 `-chaos-drop-ratio` 做更完整的场景测试
