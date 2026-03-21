# WebSocket 压测脚本使用说明

本文只介绍 `scripts/websocket_benchmark/websocket.go` 的使用方式。

## 1. 脚本用途

该脚本用于压测 WebSocket 服务，主要看三类指标：

- 连接能力：连接成功率、连接时延、峰值在线连接。
- 消息链路：发送成功率、接收速率、端到端时延（p50/p95/p99）。
- 稳定性：意外断开数、读错误、ping 失败。

## 2. 前置准备

- 可用的 WebSocket 地址，例如：`ws://localhost:8080/api/v1/ws/`。
- 鉴权信息（必须三选一）：
  - `-tokens-file`（推荐，多用户压测）
  - `-auth`（单个 Authorization）
  - `-cookie`（单个 Cookie）
- 若使用 `-tokens-file`，文件每行一个 token，支持：
  - `<token>`
  - `Bearer <token>`

示例：`scripts/benchmark/tokens.txt`

```txt
Bearer eyJhbGciOi...
Bearer eyJhbGciOi...
```

> 注意：`-clients` 大于 token 行数会直接报错。

## 3. 常用参数速查

- `-url`：WebSocket 地址。
- `-origin`：握手 Origin。
- `-clients`：总连接数。
- `-duration`：压测时长。
- `-connect-rate`：每秒新建连接数。
- `-connect-timeout`：连接超时。
- `-ping-interval`：心跳间隔，`0s` 表示关闭。
- `-send-interval`：发送间隔，`0s` 表示不主动发消息。
- `-payload`：`raw` 模式发送内容模板（支持占位符）。
- `-latency-field`：从接收消息中提取延迟时间戳的字段名。
- `-report-interval`：进度日志打印间隔。
- `-tokens-file`：token 文件路径。

支持的占位符：`{{ts_unix_nano}}`、`{{client_id}}`、`{{seq}}`。

## 4. 你的三种测试命令（可直接复用）

以下命令为 Linux/Unix shell 写法（`nohup` 后台跑）。

### 4.1 连接能力测试（5000 连接，5 分钟）

```bash
nohup go run ./scripts/benchmark/websocket.go -url ws://localhost:8080/api/v1/ws/ -origin http://localhost:5173 -clients 5000 -duration 300s -connect-rate 500 -tokens-file ./scripts/benchmark/tokens.txt -ping-interval 10s -report-interval 10s > test_5000.log 2>&1 &
```

重点看：

- `connection success rate`
- `avg connect latency` 与 `connect latency p50/p95/p99`
- `active peak connections`、`retained connection rate`

### 4.2 消息时延测试（10000 连接，1 分钟）

```bash
nohup go run ./scripts/benchmark/websocket.go -url ws://localhost:8080/api/v1/ws/ -origin http://localhost:5173 -clients 10000 -duration 60s -connect-rate 1000 -tokens-file ./scripts/benchmark/tokens.txt -ping-interval 0s -send-interval 5s -payload '{"type":"benchmark","ts":"{{ts_unix_nano}}"}' -latency-field ts -report-interval 5s > latency_10000.log 2>&1 &
```

重点看：

- `message send success rate`
- `message send/recv throughput(tps)`
- `avg end-to-end message latency`
- `message latency p50/p95/p99`

### 4.3 快速连通性检查（1000 连接，10 秒）

```bash
nohup go run ./scripts/benchmark/websocket.go -url ws://localhost:8080/api/v1/ws/ -origin http://localhost:5173 -clients 1000 -duration 10s -connect-rate 500 -tokens-file ./scripts/benchmark/tokens.txt -report-interval 5s > connect_latency.log 2>&1 &
```

重点看：

- 是否能快速拉起连接
- 是否存在明显连接失败或异常关闭

## 5. 日志查看

```bash
tail -f test_5000.log
tail -f latency_10000.log
tail -f connect_latency.log
```

结束后可在日志末尾看到 `benchmark complete` 汇总指标。

## 6. 常见问题

- 报错 `provide one of -tokens-file, -auth, or -cookie`
  - 原因：没有提供鉴权。
- 报错 `tokens-file has X tokens but clients=Y`
  - 原因：token 数量不足。
- 报错 `tokens-file cannot be used together with -auth or -cookie`
  - 原因：鉴权参数冲突。
- 时延始终不统计（`latency matched messages` 很低）
  - 原因：`-latency-field` 与实际消息字段不一致，或返回消息不带时间戳字段。

## 7. 路径说明

你当前命令使用的是 `scripts/benchmark/websocket.go`。
如果你的仓库目录是 `scripts/websocket_benchmark/websocket.go`，把命令中的路径替换为实际路径即可。
