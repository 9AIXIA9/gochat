# GoChat 代码演进与路线图（状态版）

更新时间：2026-03-07

本文基于当前仓库实现状态更新，目标是把“已完成能力、未完成缺口、下一步里程碑”放在同一视图中，便于项目管理与面试展示。

## 当前目标与完成度

- 可靠性：**部分完成**（优雅关停主流程已具备，Kafka Consumer 退出等待与 WS 全量下线仍缺失）。
- 可测性：**部分完成**（Application/Domain 测试基础较好，Port/Infrastructure 与 E2E 仍薄弱）。
- 可观测：**部分完成**（HTTP 已接入 OTEL，Kafka/WS 业务级指标和端到端链路仍需补强）。
- 安全与流控：**部分完成**（HTTP 限流已接入，WS/Kafka/Binlog 限流与背压尚未统一）。
- 可运维：**部分完成**（Compose、监控基础、迁移与常用 Make 目标具备，SLO/告警与发布策略未完成）。

## 状态矩阵（按主题）

- 健康检查 `/healthz` + `/readyz`：**已完成**（已在路由注册，并具备依赖探测）。
- 并行依赖就绪检查（MySQL/Redis/Kafka）：**已完成**（`errgroup` 并行检查已实现）。
- 关机时 Readiness 快速摘流：**部分完成**（已有优雅关机流程，但未显式切换 readiness 状态位）。
- Kafka Consumer 关闭等待（WaitGroup/超时）：**未完成**。
- WS Manager 全量主动下线（`CloseAll`）：**未完成**。
- HTTP 限流：**已完成**。
- WS/Kafka/Binlog 统一限流与背压：**未完成**。
- Make 研发命令（`test-all`、`lint`、Swagger、迁移）：**已完成**。
- CI 基础构建+测试+lint：**已完成**。
- CI 质量门禁（`-race -shuffle`、覆盖率阈值、E2E）：**未完成**。
- 统一错误 Envelope（HTTP/WS/Kafka）：**未完成**。
- Swagger 错误模型与限流响应示例：**未完成**。

## 分阶段演进计划（按优先级）

### 短期（1-2 周）- P0 稳定性与一致性

- P0-1：补 Kafka Consumer 退出等待与关停超时，避免消费协程悬挂。
- P0-2：在 WS Manager 增加 `CloseAll()` 并接入关停流程，降低半开连接风险。
- P0-3：增加服务级 readiness 状态位；收到关机信号后立即返回 503，实现快速摘流。
- P0-4：统一错误 Envelope（`code/message/details`）并在 HTTP/WS/Kafka 三通道对齐。
- P0-5：在 Swagger 明确 4xx/5xx 与 429 示例，降低客户端对接歧义。

### 中期（3-6 周）- P1 质量与观测深度

- P1-1：CI 增加 `go test -race -shuffle=on ./...`，建立并发回归门禁。
- P1-2：为 Port/Infrastructure 增加测试，覆盖率目标达到 `>=60%`（先按模块推进）。
- P1-3：引入最小 E2E（注册->登录->发送消息->WS 接收->补偿链路）。
- P1-4：补 Kafka/WS 关键指标（吞吐、失败率、积压、延迟分位），统一命名规范。
- P1-5：引入跨通道限流与背压策略，按用户/IP/连接/分区定义配额。

### 长期（6-12 周）- P2 运维与规模化

- P2-1：定义 SLO/告警规则（API、WS、Kafka），建立错误预算与复盘流程。
- P2-2：补 K8s/Helm 部署样例与发布回滚手册（蓝绿或金丝雀）。
- P2-3：完善密钥轮换、审计日志、配额治理，形成安全基线。
- P2-4：持续推进性能治理（慢查询、热点路径、Kafka 批量参数、WS 广播策略）。

## 度量与验收标准

- 可靠性：
  - 关停 60s 内完成（HTTP/Kafka/WS/OTEL）。
  - 关停后 `/readyz` 在 1s 内进入不可用状态。
- 质量：
  - CI 必过：`go test ./...`、`go test -race -shuffle=on ./...`、`golangci-lint run ./...`。
  - 覆盖率目标：application/domain >=80%，port/infrastructure >=60%。
- 可观测：
  - Kafka/WS 具备吞吐、错误率、P95 延迟、积压指标。
  - 核心链路可按 `request_id/user_id/event_id` 串联日志与 Trace。
- 安全与流控：
  - HTTP/WS/Kafka/Binlog 限流策略统一并可配置。
  - CORS/WS Origin/Cookie 策略在配置层显式化并可审计。

## 即刻可做（本周）

- 在 `internal/infrastructure/websocket/manager.go` 增加 `CloseAll()`。
- 在 `internal/infrastructure/kafka/consumer.go` 增加退出等待机制。
- 在 `cmd/api/main.go` 与 readiness checker 之间增加关停态联动。
- 在 `.github/workflows/go.yml` 增加 race/shuffle 任务与失败门禁。
- 在 `docs/swagger.yaml` 增补统一错误返回与 429 示例。

---

与 `docs/project_evaluation_optimization.md` 的分工：评估文档描述“当前能力与风险”，本文件描述“从当前状态到目标状态的落地路径与验收标准”。
