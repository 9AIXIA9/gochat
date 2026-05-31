package gateway

import (
	"context"
	"encoding/json"
	"gochat/internal/gateway/core"
	myErrors "gochat/internal/shared/errors"

	"gochat/internal/shared/command"
	"gochat/internal/shared/kernel"

	"go.uber.org/zap"
)

const (
	KeyClientMessageID = "client_message_id"
	KeyUserID          = "user_id"
)

// UpstreamRouter 实现了 gateway 上下文暴露的契约 contract.UpstreamHandler
type UpstreamRouter struct {
	idGenerator    command.IDGenerator
	publisher      command.SyncPublisher
	allowedActions map[command.Action]struct{}
}

func NewUpstreamRouter(
	publisher command.SyncPublisher,
	idGenerator command.IDGenerator,
	allowedActions map[command.Action]struct{},
) core.UpstreamHandler {
	return &UpstreamRouter{
		publisher:      publisher,
		idGenerator:    idGenerator,
		allowedActions: allowedActions,
	}
}

// HandleUpstream 处理上行消息，并返回应发给客户端的回执字节
func (r *UpstreamRouter) HandleUpstream(ctx context.Context, userID kernel.UserID, payload []byte) []byte {
	var env core.UpstreamEnvelope
	if err := json.Unmarshal(payload, &env); err != nil {
		zap.L().Warn("gateway: invalid envelope format", zap.Error(err), zap.String("userID", userID.String()))
		return nil // 无效报文无法确认 clientMessageID，直接抛弃或关闭连接
	}

	// 包未带 ID 也无法回复
	if env.ClientMessageID == "" || env.Action == "" {
		zap.L().Warn("gateway: envelope missing clientMessageID or action", zap.String("userID", userID.String()))
		return nil
	}

	// 安全查验：Action 必须在白名单映射表中
	_, ok := r.allowedActions[(env.Action)]
	if !ok {
		zap.L().Warn("gateway: unauthorized or unknown action", zap.String("action", string(env.Action)), zap.String("userID", userID.String()))
		// 直接快速阻断
		ackBytes, _ := NewAckError(env.ClientMessageID, myErrors.ErrWrongCommandAction)
		return ackBytes
	}

	// 将 UpstreamEnvelope 装裱为一个领域事件发送，使用映射后的内部真实 Action
	ev := command.NewStandardCommand(kernel.ID(userID), env.Action, env.Payload, r.idGenerator)
	ev.AddHeaders(map[string]string{
		KeyClientMessageID: env.ClientMessageID.String(),
		KeyUserID:          userID.String(),
	})

	// 阻塞投递给事件总线，等待结果
	err := r.publisher.Publish(ctx, ev)
	if err != nil {
		zap.L().Error("gateway: sync publish failed", zap.Error(err), zap.String("action", string(env.Action)))

		// 投递失败的回执 (Fast ACK = error)
		ackBytes, err := NewAckError(env.ClientMessageID, myErrors.ErrServerBusy)
		if err != nil {
			zap.L().Error("gateway: create ack error failed", zap.Error(err), zap.String("clientMessageID", env.ClientMessageID.String()))
			return nil
		}
		return ackBytes
	}
	data, err := NewAckReceived(env.ClientMessageID)
	if err != nil {
		zap.L().Error("gateway: create ack received failed", zap.Error(err), zap.String("clientMessageID", env.ClientMessageID.String()))
		return nil
	}
	return data
}
