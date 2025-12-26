# GoChat 后端

一个采用 DDD 风格的聊天/好友/房间管理后端，支持 HTTP、WebSocket 与 Kafka 三种交付通道。以 Go 构建，具备完善的依赖注入、持久化与可观测性能力。

## 功能特性
- 授权认证：注册/登录/刷新/解析（JWT + bcrypt）
- 聊天：私聊与房间消息（HTTP + WebSocket + Kafka）
- 好友与房间：好友请求、好友关系、房间成员管理
- 通知：系统消息与未送达消息通知
- 事件驱动：跨上下文的 Kafka 事件处理
- 可观测性：OpenTelemetry Traces、Prometheus 指标、Grafana 仪表盘

## 技术栈
- Go，Gin（HTTP），Gorilla WebSocket（WS）
- Kafka（confluent-kafka-go），MySQL（Gorm），Redis
- Google Wire 进行依赖注入
- Viper + YAML + .env 进行配置管理
- Zap 结构化日志
- Swagger 文档（基于注解）
- Docker 与 docker-compose

## 项目结构
- `cmd/api`：应用入口与依赖注入装配（DI）
- `internal/{context}`：DDD 分层（`application`、`domain`、`infrastructure`、`port`）
- `internal/delivery`：HTTP/Kafka/WebSocket 适配层
- `db/migrations`：SQL 迁移脚本
- `deployment`：OTEL Collector、Prometheus、Grafana、告警与各类初始化脚本
- `docs`：API 文档、评估与演进路线文档

## 使用 Makefile 快速上手

推荐通过 Makefile 执行常用操作（Windows 需安装 GNU make）：

```bash
# 启动（不使用构建缓存）
make up

# 快速启动（使用构建缓存）
make up-fast

# 查看应用日志
make logs-app

# 查看当前服务状态
make status

# 停止并移除容器
make down

# 彻底清理（容器、镜像、卷）
make destroy

# 本地运行测试（race 检测 + 随机化）
make test-all

# 代码静态检查（需安装 golangci-lint）
make lint

# 生成 Swagger 文档（基于注解）
make swagger

# 数据库迁移（容器网络 DSN）
make migrate-status
make migrate-up
make migrate-down
make migrate-reset

# 数据库迁移（宿主机端口 127.0.0.1:13306）
make migrate-status-host
make migrate-up-host
make migrate-down-host
make migrate-reset-host
```

> 备注：仍可使用 docker-compose 作为备选方式，见下文“快速开始”。

## 快速开始

前置条件：已安装 Docker 与 docker-compose

```bash
# 构建并启动（如需调整 compose 文件路径，请先适配）
docker-compose up --build -d

# 查看日志
docker-compose logs -f app

# 停止
docker-compose down
```

默认地址：
- HTTP：`http://localhost:8080`
- BasePath：`/api/v1`
- Swagger：参见 `docs/swagger.yaml` 以及 `cmd/api/main.go` 中的注解

## 健康与就绪
- 健康检查：通过 HTTP Handler 返回状态
- 优雅关机：`cmd/api/main.go` 基于信号触发，依次关闭 HTTP、Kafka Producer/Consumers、Binlog、Email、OTEL

## 开发

```bash
# 运行测试
go test ./...

# 生成 Swagger（基于注解）
go generate ./cmd/api/main.go

# Lint（如已安装 golangci-lint）
golangci-lint run ./...
```

## 配置
- 默认配置：`config/config.yaml`
- 环境变量覆盖：通过 `internal/infrastructure/godotenv` 读取 `.env`

## 文档导航

- 概览与评估
  - 项目评估报告：`docs/project_evaluation_optimization.md`
  - 代码演进与路线图：`docs/code_evolution.md`
  - 项目演进路线图（版本迭代记录）：`docs/roadmap.md`
- API 文档
  - Swagger YAML：`docs/swagger.yaml`
  - Swagger JSON：`docs/swagger.json`
  - 注解入口：`cmd/api/main.go`
- 实时通信
  - WebSocket 说明：`docs/websocket.md`
- 部署与监控
  - docker-compose：`docker-compose.yml`
  - Dockerfile：`Dockerfile`
  - OTEL 采集器：`deployment/otel-collector-config.yaml`
  - Prometheus：`deployment/prometheus.yml`
  - Grafana 数据源与仪表：`deployment/grafana-datasources.yml`、`deployment/grafana-dashboards/gochat-app.json`
  - 告警：`deployment/prometheus-alerts.yml`
- 数据库与初始化
  - SQL 迁移：`db/migrations/`
  - MySQL 初始化：`deployment/grafana-dashboards/mysql-init/init_users.sh`
  - Redis ACL 初始化：`deployment/redis_init/init_users.acl`
- 构建与脚本
  - Makefile：`Makefile`（常用：`make up`、`make up-fast`、`make logs-app`、`make test-all`、`make migrate-up`）
  - 统计代码行数：`scripts/count_golang_code_lines.ps1`

## 许可协议
MIT
