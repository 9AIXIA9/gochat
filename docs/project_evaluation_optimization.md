# GoChat 后端项目评估与优化建议（现状版）

更新时间：2026-03-07

本报告基于仓库当前实现进行再评估，聚焦三件事：
1) 当前完成度；2) 关键风险；3) 下一步优化优先级。

## 评估依据（代码证据）

- 路由与健康检查：`cmd/api/di/providers_http.go`
- 启停与关停流程：`cmd/api/main.go`
- Kafka 消费与生产：`internal/infrastructure/kafka/consumer.go`、`internal/infrastructure/kafka/event_publisher.go`
- WebSocket 连接管理：`internal/infrastructure/websocket/manager.go`
- 工程命令与迁移：`Makefile`
- CI 流水线：`.github/workflows/go.yml`
- 文档与接口：`docs/swagger.yaml`

## 总体评估结论

- 架构成熟度：**较高**（DDD 分层与 DI 落地完整，多上下文边界清晰）。
- 交付能力：**较高**（HTTP/WS/Kafka/Binlog 同时具备，工程化基础扎实）。
- 工程质量：**中等偏上**（测试与 CI 有基础，但并发门禁、E2E 和覆盖率目标未完全闭环）。
- 运行可靠性：**中等**（关机流程完整，但 Consumer/WS 关停细节存在尾部风险）。
- 安全与治理：**中等**（HTTP 限流已落地，跨通道流控、错误契约与审计治理仍有缺口）。

## 完成度盘点（已完成 / 部分完成 / 未完成）

- 已完成
  - DDD + 端口适配器 + Wire 注入体系。
  - `/healthz`、`/readyz` 路由与依赖就绪探测。
  - HTTP 限流中间件接入业务路由。
  - 优雅关停主流程（HTTP、Kafka Producer/Consumer、Binlog、OTEL）。
  - Make 常用研发命令（测试、lint、Swagger、迁移）。
  - CI 的构建、常规测试、lint。

- 部分完成
  - Readiness 与关停联动（缺少显式“摘流状态位”）。
  - 可观测性（HTTP 较完整，Kafka/WS 业务指标深度不足）。
  - 测试结构（应用层较好，端口/基础设施层覆盖不足）。

- 未完成
  - Kafka Consumer 退出等待与超时控制（防止 goroutine 残留）。
  - WS Manager 的 `CloseAll()` 与全量连接回收。
  - 跨通道统一错误 Envelope（HTTP/WS/Kafka）。
  - WS/Kafka/Binlog 统一限流与背压策略。
  - CI 并发门禁（`-race -shuffle`）、覆盖率阈值、E2E 流程。
  - Swagger 中统一错误模型与 429 等典型错误响应示例。

## 关键风险与影响（按严重度）

- P0 高风险
  - 关停尾部风险：Consumer 与 WS 缺少“可等待关闭”的完整闭环，可能导致进程退出时资源未完全释放。
  - 契约不一致：跨通道错误语义未统一，客户端重试策略与排障路径容易分叉。

- P1 中风险
  - 可观测盲区：Kafka/WS 的队列积压、失败分布、延迟分位不可见，容量规划难度较高。
  - 质量门禁缺口：CI 缺 `-race -shuffle` 与 E2E，隐藏并发问题和集成回归。

- P2 可优化项
  - 安全治理深度：缺少统一配额、密钥轮换、审计日志等长期治理能力。
  - 交付体系：缺少面向 K8s 的标准化发布与 SLO 告警体系。

## 优化建议（可执行）

- 第一阶段（1-2 周，P0）
  - 为 `Consumer` 增加 WaitGroup/退出信号/关停超时。
  - 为 `Manager` 增加 `CloseAll()` 并在关机流程调用。
  - 增加 readiness 关停态（收到信号后立刻返回 503）。
  - 定义并落地统一错误 Envelope（先覆盖最核心链路）。

- 第二阶段（3-6 周，P1）
  - CI 增加 `go test -race -shuffle=on ./...`。
  - 建立最小 E2E（注册->登录->发消息->WS 接收->补偿）。
  - 按模块补 Port/Infrastructure 测试，逐步达到覆盖率目标。
  - 增加 Kafka/WS 指标与 Trace 标签，形成可观测闭环。

- 第三阶段（6-12 周，P2）
  - 定义 API/WS/Kafka 的 SLO 与告警规则。
  - 建立 K8s/Helm 发布与回滚模板。
  - 推进配额治理、密钥轮换、审计日志与性能基线。

## 建议的度量指标（用于复盘）

- 稳定性：部署后 7 天内异常重启次数、关停超时率、未清理连接数。
- 质量：CI 失败类型分布（lint/test/race/e2e）、回归缺陷数、覆盖率趋势。
- 观测：Kafka lag、WS 在线连接、P95 延迟、错误率。
- 效率：PR 到上线周期、缺陷修复 lead time、发布回滚次数。

## 结语

当前项目已经具备“可用且可扩展”的工程基础，下一步重点不是继续堆功能，而是优先补齐“稳定性闭环 + 契约一致性 + 质量门禁”三件核心能力。完成这些后，再推进 SLO 与规模化运维，性价比最高。

---

对应演进路径见 `docs/code_evolution.md`。
