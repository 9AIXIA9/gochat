package local

import (
	"context"
	"gochat/internal/infrastructure/metrics"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/timeout"
	"gochat/internal/ws_gateway/contract"
	"gochat/internal/ws_gateway/core"
)

// WSGatewayService 是 WSGatewayService 在单体进程内环境的默认实现。
type WSGatewayService struct {
	manager *core.Manager
}

// NewWSGatewayService 创建本地推送服务
func NewWSGatewayService(manager *core.Manager) contract.WSGatewayService {
	return &WSGatewayService{
		manager: manager,
	}
}

// PushToUser 向本地直接持有连接的用户分发 payload
func (s *WSGatewayService) PushToUser(ctx context.Context, userID string, payload []byte) error {
	if timeout.CheckCtxTimeout(ctx) {
		return myErrors.ErrTimeout
	}

	sessions, ok := s.manager.GetByUserID(userID)
	if !ok {
		// 分布式/微服务架构下：如果这里没拿到说明连到别的 Pod 了，或者是离线。
		metrics.WSMessageOut(ctx, "push_local", "not_found")
		return myErrors.ErrNotFound
	}

	var lastErr error
	for _, session := range sessions {
		if err := session.Send(payload); err != nil {
			lastErr = err
		}
	}

	if lastErr != nil {
		metrics.WSMessageOut(ctx, "push_local", "send_failed")
		return lastErr
	}

	metrics.WSMessageOut(ctx, "push_local", "ok")
	return nil
}

// Broadcast 向多用户批量分发（适用于普通广播场景）
func (s *WSGatewayService) Broadcast(ctx context.Context, userIDs []string, payload []byte) error {
	if timeout.CheckCtxTimeout(ctx) {
		return myErrors.ErrTimeout
	}

	for _, id := range userIDs {
		// 单体应用里，简单地对每个用户去查一下是否由于在本地并发送即可。
		// 如果是大量广播，这里的实现可以做进一步的高性能调优。
		if sessions, ok := s.manager.GetByUserID(id); ok {
			for _, session := range sessions {
				_ = session.Send(payload)
			}
		}
	}

	metrics.WSMessageOut(ctx, "push_broadcast", "ok")
	return nil
}
