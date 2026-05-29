package contract

import "context"

// UpstreamHandler 定义了从长连接接收到业务上行数据时，Push 中心该如何将其透传出去的契约。
// 典型的实现类可能会把消息丢到 Kafka，供 Chat 或 Notification 模块去消费。
type UpstreamHandler interface {
	// HandleUpstream 处理上行消息
	HandleUpstream(ctx context.Context, userID string, payload []byte)
}
