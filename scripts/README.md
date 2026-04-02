# Scripts 使用说明

本文档介绍 `scripts/` 目录下所有脚本的用途、使用方式和常见场景。

## 目录概览

- `fix_swagger.py`
- `websocket_benchmark/websocket.go`
- `websocket_benchmark/generate_token_pool.ps1`
- `websocket_benchmark/token_pool_generator/main.go`

## 使用约定

- 建议在仓库根目录执行所有命令（即 `backend/` 目录）。
- 若脚本依赖服务接口（注册、登录、WebSocket），请先启动后端服务。

---

## 1) fix_swagger.py

### 作用

修复 Swagger/OpenAPI 文档中参数对象的 `example` 字段，将其替换为 `x-example`，避免工具链兼容问题。

### 适用文件

- `docs.go`
- `swagger*.go`
- `*swagger*.json|yaml|yml`
- `*openapi*.json|yaml|yml`

### 常用命令

1. 处理单文件

```bash
python scripts/fix_swagger.py docs/docs.go
```

2. 只预览，不修改文件

```bash
python scripts/fix_swagger.py docs/docs.go --dry-run
```

3. 修改前先备份

```bash
python scripts/fix_swagger.py docs/docs.go --backup
```

4. 批量扫描目录并处理

```bash
python scripts/fix_swagger.py . --batch
```

5. 输出到指定目录

```bash
python scripts/fix_swagger.py docs/docs.go --output tmp/swagger_fixed
```

### 参数说明

- `--dry-run`：模拟运行，不落盘
- `--backup`：修改前生成备份文件（`.bak`）
- `--batch`：目录批量模式
- `--output` / `-o`：指定输出目录

---

## 2) websocket_benchmark/websocket.go

### 作用

WebSocket 压测入口，支持：

- 大量连接建立
- 心跳保活（Ping）
- 文本消息发送
- 私聊协议消息发送
- 端到端延迟统计
- Chaos 随机断连演练

### 最小示例

```bash
go run ./scripts/websocket_benchmark -auth "Bearer <token>" -clients 100 -duration 60s
```

### 典型示例

1. 连接 + 心跳

```bash
go run ./scripts/websocket_benchmark -auth "Bearer <token>" -clients 500 -duration 120s -ping-interval 5s
```

2. 发送 raw 消息

```bash
go run ./scripts/websocket_benchmark -auth "Bearer <token>" -clients 200 -duration 60s -send-interval 2s -send-mode raw
```

3. 发送 private 消息（接收人列表）

```bash
go run ./scripts/websocket_benchmark -tokens-file ./scripts/websocket_benchmark/tokens.txt -send-mode private -recipients-file ./scripts/websocket_benchmark/recipients.txt -send-interval 2s -clients 200 -duration 60s
```

### 鉴权输入（三选一）

- `-tokens-file`
- `-auth`
- `-cookie`

### 关键参数

- `-url`：WebSocket 地址，默认 `ws://localhost:8080/api/v1/ws/`
- `-clients`：总连接数
- `-duration`：测试时长
- `-connect-rate`：每秒建连数（`0` 表示无爬坡）
- `-connect-timeout`：握手超时
- `-ping-interval`：心跳间隔（`0` 禁用）
- `-send-interval`：发消息间隔（`0` 禁用）
- `-send-mode`：`raw` / `private`
- `-payload`：raw 模式消息模板
- `-content`：private 模式内容模板
- `-latency-field`：延迟时间戳字段名（默认 `bench_ts_ns`）
- `-report-interval`：进度打印间隔
- `-recipient-id`：private 模式单接收人
- `-recipients-file`：private 模式接收人列表
- `-chaos-drop-ratio`：随机断连比例
- `-chaos-drop-after`：随机断连触发延迟

### 注意事项

- `-send-mode private` 必须配 `-recipient-id` 或 `-recipients-file`
- `-tokens-file` 与 `-auth` / `-cookie` 不能同时使用
- `-tokens-file` 行数必须不少于 `-clients`

---

## 3) websocket_benchmark/generate_token_pool.ps1

### 作用

PowerShell 串行版 token 生成器：循环调用注册与登录接口，产出 token 文件。

### 适用场景

- 小规模账号（如 100~5000）
- Windows 环境快速脚本化

### 命令示例

```powershell
pwsh .\scripts\websocket_benchmark\generate_token_pool.ps1 -BaseUrl http://localhost:8080/api/v1 -Count 1000
```

### 关键参数

- `-BaseUrl`：API 根地址
- `-Count`：生成数量
- `-StartIndex`：起始序号
- `-Password`：统一密码
- `-EmailPrefix` / `-EmailDomain`：邮箱生成规则
- `-OutFile`：token 输出路径
- `-DetailFile`：详情 JSONL 输出路径
- `-ThrottleMs`：每次请求间隔（毫秒）
- `-DryRun`：仅生成模拟 token
- `-Overwrite`：覆盖输出文件（默认不覆盖，按续写模式追加）

### 输出

- `scripts/websocket_benchmark/tokens.txt`
- `scripts/websocket_benchmark/token_pool.jsonl`

---

## 4) websocket_benchmark/token_pool_generator/main.go

### 作用

Go 并发版 token 生成器，面向大规模账号（例如 10 万）场景。

相比 PowerShell 版，支持：

- Worker 并发
- HTTP 连接复用
- 自动重试与退避
- 周期性进度输出

### 建议场景

- 1 万到 10 万级 token 生成
- 需要控制吞吐与重试策略

### 命令示例

1. dry-run 验证参数

```bash
go run ./scripts/websocket_benchmark/token_pool_generator -dry-run -count 1000 -workers 200
```

2. 真实生成（10 万示例）

```bash
go run ./scripts/websocket_benchmark/token_pool_generator -base-url http://localhost:8080/api/v1 -count 100000 -workers 300 -max-retries 2 -retry-backoff 250ms -timeout 12s -out-file scripts/websocket_benchmark/tokens.txt -detail-file scripts/websocket_benchmark/token_pool.jsonl
```

### 关键参数

- `-base-url`：API 根地址
- `-count`：生成数量
- `-start-index`：起始序号
- `-password`：统一密码
- `-email-prefix` / `-email-domain`：邮箱生成规则
- `-workers`：并发 worker 数
- `-timeout`：单请求超时
- `-max-retries`：失败重试次数
- `-retry-backoff`：重试基础退避
- `-progress-interval`：进度输出间隔
- `-dry-run`：模拟生成
- `-append`：是否续写输出文件（默认 `true`）
- `-out-file` / `-detail-file`：输出路径

### 调参建议

- 先从 `-workers 100` 起步，观察后端 CPU / DB / Redis，再逐步升到 200 或 300
- 失败率高时优先增大 `-timeout`，其次降低 `-workers`
- 生产环境或共享环境压测前，先小规模试跑（如 `-count 1000`）

---

## 推荐流程

1. 用 `token_pool_generator` 生成多用户 token（大规模优先 Go 版）
2. 用 `websocket.go` 执行连接与消息压测
3. 如需修复 Swagger 文档字段，再运行 `fix_swagger.py`

如果你只做 WebSocket 压测，建议从 `scripts/websocket_benchmark/README.md` 开始。
