package contract

import (
	"context"
)

// WSGatewayService 是对外部业务系统暴露的推送接口。
type WSGatewayService interface {
	// PushToUser 向指定用户的所有在线设备推送数据
	PushToUser(ctx context.Context, userID string, payload []byte) error

	// Broadcast 向一批用户广播数据
	Broadcast(ctx context.Context, userIDs []string, payload []byte) error
}

// RouterProvider 为 Push 模块提供节点路由等实现
type RouterProvider interface {
	// GetNodeForUser 查找对应用户所在的 Pod/Node 地址，如果在这个节点自己，则返回 self 等标记
	GetNodeForUser(ctx context.Context, userID string) (string, error)
}
