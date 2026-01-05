# GoChat 代码演进与路线图

更新时间：2025-12-26

本文描述 GoChat 后端未来的发展方向与阶段性目标，作为简历中“工程演进能力”的展示，配合《项目评估报告》使用。

## 目标概览

- 可靠：具备完善的优雅关机、健康/就绪探针、重试/退避/幂等、容量可视化。
- 可测：单元/集成/E2E 测试体系完善，覆盖 Port/Infrastructure 层，CI 门禁稳健。
- 可观测：端到端 Trace/Metric/Log 一致性与可追踪性，Kafka/WS 关键指标完善。
- 安全：多通道限流/背压、CORS/WS Origin、Cookie 策略统一、密钥管理与轮换。
- 可运维：可移植的部署（Compose→K8s）、告警与 SLO、发布策略（蓝绿/金丝雀）。

## 分阶段演进计划

### 短期（1–2 周）

- 健康/就绪：将现有 HTTP 健康检查拆分为 `/healthz` 与 `/readyz`；并行依赖检查（MySQL/Redis/Kafka/Canal）带短超时，旁路限流/鉴权中间件；收到关机信号时
  `/readyz` 立即返回 503。
- 错误模型：统一 HTTP/WS/Kafka 的错误 Envelope（code/message/details），在 Swagger 中给出示例；WS 回执字段与 Kafka 错误处理标准化。
- 观测增强：为 Kafka Consumer/Producer、WS Handler、核心用例增加 Span/Metric，统一命名规范；增加处理耗时、失败率、队列长度等指标。
- 测试补齐（起步）：
	- HTTP Handler 与中间件（含限流）使用 `httptest`；
	- WebSocket 编解码、心跳、基础限流单测；
	- Kafka 消费处理器的无外部依赖单测；
	- Gorm/Redis 适配器最小集成测试（事务回滚/Test DB）。

### 中期（3–6 周）

- 安全与流控：将 HTTP 的限流扩展至 WS/Kafka/Binlog（按连接/用户/IP 消息速率、按键/分区并发与节流、Binlog 拉取节流与背压）；统一
  CORS/WS Origin 白名单与 Cookie 策略；抽象限流接口便于替换。
- 可靠性：Kafka 消费重试/退避（复用 `pkg/utils/backoff.go`）、幂等键策略；为 Consumer 增加退出等待（WaitGroup）与关停超时；WS
  Manager 增加 CloseAll 主动下线。
- 集成与 CI：docker-compose 驱动的 E2E（注册→登录→消息→WS 接收→未送达补偿），在 CI（GitHub Actions）执行；引入 `golangci-lint`、
  `go test -race -shuffle` 门禁；覆盖率门槛（application/domain ≥80%，port/infrastructure ≥60%）。
- 文档与开发者体验：在 `docs/` 增加上下文架构/事件流图、Quick Start与 FAQ；Makefile 增设 `test-all`、`lint`、
  `generate-swagger`、`run-migrations`。

### 长期（6–12 周）

- 部署策略：提供 Helm Chart 与 K8s 部署样例；支持蓝绿/金丝雀发布与回滚指南；配置化健康探针与资源限额。
- 告警与 SLO：完善 `deployment/prometheus-alerts.yml`，定义 API/WS/Kafka 的 SLO 与告警规则；接入错误预算与容量规划流程。
- 安全与合规：令牌轮换、密钥管理（Vault/KMS）、审计日志（敏感操作），对外 API 的速率与配额管理。
- 性能优化：DB 索引与慢查询治理、WebSocket 扩展与广播优化、Kafka 批量与吞吐调优、热点用例基准测试。

## 关键设计与度量

- Trace/Metric 规范：
	- Trace 名称：`<context>.<usecase|handler>.<operation>`；常用属性（user_id, room_id, event_id, topic/partition,
	  ws_client_id）。
	- Metric：HTTP/WS/Kafka 的延迟/错误率/吞吐；Kafka 消费 lag、commit 成功率；WS 并发/背压/丢弃数；DB 查询耗时与命中率。
- 健康检查策略：
	- Liveness 仅本进程探活；Readiness 带依赖检查与短超时；关机时 Readiness 置 503；对公网部署的访问控制与信息最小化。
- 测试门禁与阈值：
	- 覆盖率：application/domain ≥80%，port/infrastructure ≥60% 起步；
	- 质量门禁：`golangci-lint` 必过，`go test -race -shuffle` 必过；
	- 性能基线：热点用例与仓储查询提供基准，监控回归。

## 小型即刻改进清单（Low-risk quick wins）

- 在 WS Manager 增加 `CloseAll()` 并接入关机流程；
- Kafka Consumer 增加 WaitGroup/关闭等待与超时；
- 健康检查拆分 `/healthz` 与 `/readyz` 并添加依赖并行检查；
- Makefile 增加 `test-all`、`lint`、`generate-swagger`、`run-migrations`；
- 在 Swagger 中补充统一错误模型示例与限流说明。

——

本文与《项目评估报告》相互配合：评估报告呈现当前状态与长/短板，本文件给出未来演进与度量标准，便于在简历与面试中体现“可持续建设能力”。
