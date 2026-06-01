package core

import (
	"context"
	"gochat/internal/shared/kernel"
)

// UpstreamHandler 定义了从长连接接收到业务上行数据时，gateway 中心该如何将其透传出去的契约
type UpstreamHandler interface {
	// HandleUpstream 处理上行消息
	HandleUpstream(ctx context.Context, userID kernel.UserID, payload []byte) []byte
}
