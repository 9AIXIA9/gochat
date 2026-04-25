# GoChat WebSocket 压测性能报告

## 1. 报告信息

- 报告日期：2026-04-25
- 测试对象：GoChat WebSocket 服务（`/api/v1/ws/`）
- 测试工具：`scripts/websocket_benchmark/websocket.go`
- 原始日志目录：`docs/test_report/websocket_stress_test/raw_logs/linux_4core_16g/`

## 2. 测试目标

- 验证高并发连接建立与长连接保持稳定性。
- 验证高并发下心跳与消息收发吞吐能力。
- 识别当前压测口径下的瓶颈与后续优化方向。

## 3. 测试环境与口径

### 3.1 环境

- 环境标识：`linux_4core_16g`
- 目标地址：`ws://172.16.0.12:808x/api/v1/ws/`
- 并发线程上限：`max_threads=50000`

### 3.2 场景设计

场景 A：连接与保活稳定性测试（10 端口并行）

- 端口范围：`8080` - `8089`
- 单端口连接数：`26000`
- 总连接数：`260000`
- 测试时长：`10m`
- 建连速率：`500/s`
- 心跳：关闭（`ping_interval=0s`）
- 消息发送：关闭（`send_interval=0s`）

场景 B：消息吞吐与心跳测试（单端口）

- 端口：`8080`
- 连接数：`10000`
- 测试时长：`60s`
- 建连速率：`500/s`
- 心跳：开启（`ping_interval=5s`）
- 消息发送：开启（`send_interval=1s`）

## 4. 测试结果

### 4.1 场景 A：连接稳定性（260000 总连接）

汇总日志：`websocket_connection_stress_test/websocket_stress_summary.txt`

| 指标               | 结果                  |
|------------------|---------------------|
| 总尝试连接数           | 260000              |
| 总成功连接数           | 260000              |
| 总连接成功率           | 100.00%             |
| 峰值在线连接数          | 260000              |
| 保留连接数            | 260000              |
| 保留连接率            | 100.00%             |
| 异常关闭连接           | 0                   |
| 读错误（read errors） | 0                   |
| 任务完成状态           | 10/10 端口全部 `NORMAL` |

连接时延（各端口范围）

- 平均建连时延：`0.57ms` - `0.77ms`
- P50：`0.49ms` - `0.60ms`
- P95：`0.58ms` - `0.88ms`
- P99：`0.97ms` - `6.92ms`

结论：在当前环境和口径下，系统可稳定承载 26 万长连接，未出现连接失败扩散、异常断开或读错误。

### 4.2 场景 B：消息吞吐与心跳（10000 连接）

日志：`websocket_message_stress_test.txt`

| 指标                 | 结果                       |
|--------------------|--------------------------|
| 尝试连接数              | 10000                    |
| 成功连接数              | 10000                    |
| 连接成功率              | 100.00%                  |
| 保留连接数              | 10000                    |
| 保留连接率              | 100.00%                  |
| 异常关闭连接             | 0                        |
| 平均建连时延             | 0.58ms                   |
| 建连时延 P50/P95/P99   | 0.53ms / 0.64ms / 0.77ms |
| 发送成功消息数            | 494222                   |
| 发送失败消息数            | 0                        |
| 消息发送成功率            | 100.00%                  |
| 接收消息数              | 494222                   |
| 投递率（received/sent） | 100.00%                  |
| 发送吞吐（TPS）          | 8197.40                  |
| 接收吞吐（TPS）          | 8197.40                  |
| 心跳成功/失败            | 94815 / 0                |
| 读错误（read errors）   | 0                        |

## 5. 综合结论

- 连接能力：通过。系统在 26 万并发连接压测下表现稳定，连接成功率和保活率均为 100%。
- 消息能力：通过。在 1 万连接、每秒持续发送条件下，消息收发吞吐约 8.2k TPS，发送成功率与投递率均为 100%。
- 稳定性：通过。两个场景均未出现异常断连放大、读错误、心跳失败。

---

原始数据来源：

-

`docs/test_report/websocket_stress_test/raw_logs/linux_4core_16g/websocket_connection_stress_test/websocket_stress_summary.txt`
-
`docs/test_report/websocket_stress_test/raw_logs/linux_4core_16g/websocket_connection_stress_test/websocket_stress_report_8080.txt`
-
`docs/test_report/websocket_stress_test/raw_logs/linux_4core_16g/websocket_connection_stress_test/websocket_stress_report_8081.txt`
-
`docs/test_report/websocket_stress_test/raw_logs/linux_4core_16g/websocket_connection_stress_test/websocket_stress_report_8082.txt`
-
`docs/test_report/websocket_stress_test/raw_logs/linux_4core_16g/websocket_connection_stress_test/websocket_stress_report_8083.txt`
-
`docs/test_report/websocket_stress_test/raw_logs/linux_4core_16g/websocket_connection_stress_test/websocket_stress_report_8084.txt`
-
`docs/test_report/websocket_stress_test/raw_logs/linux_4core_16g/websocket_connection_stress_test/websocket_stress_report_8085.txt`
-
`docs/test_report/websocket_stress_test/raw_logs/linux_4core_16g/websocket_connection_stress_test/websocket_stress_report_8086.txt`
-
`docs/test_report/websocket_stress_test/raw_logs/linux_4core_16g/websocket_connection_stress_test/websocket_stress_report_8087.txt`
-
`docs/test_report/websocket_stress_test/raw_logs/linux_4core_16g/websocket_connection_stress_test/websocket_stress_report_8088.txt`
-
`docs/test_report/websocket_stress_test/raw_logs/linux_4core_16g/websocket_connection_stress_test/websocket_stress_report_8089.txt`

- `docs/test_report/websocket_stress_test/raw_logs/linux_4core_16g/websocket_message_stress_test.txt`
