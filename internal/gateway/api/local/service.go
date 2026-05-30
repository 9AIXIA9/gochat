package local

import (
	"context"
	"encoding/json"
	"gochat/internal/gateway/core"
	"gochat/internal/infrastructure/metrics"
	"gochat/internal/shared/contract"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
	"gochat/internal/shared/timeout"
)

// gatewayService 是 contract.GatewayService 在单体进程内环境的默认实现
type gatewayService struct {
	manager *core.Manager
}

// NewLocalGatewayService 创建本地网关服务
func NewLocalGatewayService(manager *core.Manager) contract.GatewayService {
	return &gatewayService{
		manager: manager,
	}
}

func (s *gatewayService) PushToUser(ctx context.Context, userID kernel.UserID, envelop *core.DownstreamEnvelop) error {
	if timeout.CheckCtxTimeout(ctx) {
		return myErrors.ErrTimeout
	}

	payload, err := json.Marshal(envelop)
	if err != nil {
		return err
	}

	sessions, ok := s.manager.GetByUserID(userID)
	if !ok {
		metrics.WSMessageOut(ctx, "gateway_local", "not_found")
		return myErrors.ErrNotFound
	}

	var lastErr error
	for _, session := range sessions {
		if err := session.Send(payload); err != nil {
			lastErr = err
		}
	}

	if lastErr != nil {
		metrics.WSMessageOut(ctx, "gateway_local", "send_failed")
		return lastErr
	}

	metrics.WSMessageOut(ctx, "gateway_local", "ok")
	return nil
}
