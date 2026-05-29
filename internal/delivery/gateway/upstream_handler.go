package gateway

import (
	"context"
	"encoding/json"

	"gochat/internal/gateway/contract"
	"gochat/internal/infrastructure/validator"

	"go.uber.org/zap"

	chatApp "gochat/internal/chat/application"
	"gochat/internal/shared/kernel"
)

// 定义当前所有的 Topic 的常量
const (
	SendPrivateMessageTopic = "chat.send_private_message"
	SendRoomMessageTopic    = "chat.send_room_message"
)

// Request 表示外网长连接上行的标准数据格式
type Request struct {
	Topic   string          `json:"topic" validate:"required"`
	Payload json.RawMessage `json:"payload" validate:"required"`
}

// SendPrivateMessageData 将网关的 JSON 映射到业务输入参数
type SendPrivateMessageData struct {
	SenderID    kernel.UserID `json:"-" validate:"required"`
	RecipientID kernel.UserID `json:"recipient_id" validate:"required"`
	Content     string        `json:"content" validate:"required,max=1000"`
}

// SendRoomMessageData 同理
type SendRoomMessageData struct {
	SenderID kernel.UserID `json:"-" validate:"required"`
	RoomID   kernel.RoomID `json:"room_id" validate:"required"`
	Content  string        `json:"content" validate:"required,max=1000"`
}

// UpstreamRouter 实现了 gateway 上下文暴露的契约 contract.UpstreamHandler
type UpstreamRouter struct {
	validator              *validator.Validator
	sendPrivateMessageCase chatApp.SendPrivateMessageUseCase
	sendRoomMessageCase    chatApp.SendRoomMessageUseCase
}

func NewUpstreamRouter(
	v *validator.Validator,
	privCase chatApp.SendPrivateMessageUseCase,
	roomCase chatApp.SendRoomMessageUseCase,
) contract.UpstreamHandler {
	return &UpstreamRouter{
		validator:              v,
		sendPrivateMessageCase: privCase,
		sendRoomMessageCase:    roomCase,
	}
}

// HandleUpstream 由 Gateway 核心收取到客户端消息时回调
func (r *UpstreamRouter) HandleUpstream(ctx context.Context, userID string, payload []byte) {
	var req Request
	if err := json.Unmarshal(payload, &req); err != nil {
		zap.L().Warn("gateway: payload unmarshal error", zap.Error(err), zap.String("userID", userID))
		return
	}

	uid := kernel.UserID(userID) // 使用类型强转替代 kernel.ParseUserID，因为 kernel 层直接使用 string 别名

	switch req.Topic {
	case SendPrivateMessageTopic:
		var data SendPrivateMessageData
		if err := json.Unmarshal(req.Payload, &data); err != nil {
			zap.L().Warn("gateway: invalid private msg payload", zap.Error(err))
			return
		}
		data.SenderID = uid

		// 校验防参数错误
		if _, err := r.validator.Validate(ctx, &data); err != nil {
			zap.L().Warn("gateway: payload validate error", zap.Error(err))
			return
		}

		_, _ = r.sendPrivateMessageCase.Execute(ctx, &chatApp.SendPrivateMessageInput{
			SenderID:    data.SenderID,
			RecipientID: data.RecipientID,
			Content:     data.Content,
		})

	case SendRoomMessageTopic:
		var data SendRoomMessageData
		if err := json.Unmarshal(req.Payload, &data); err != nil {
			zap.L().Warn("gateway: invalid room msg payload", zap.Error(err))
			return
		}
		data.SenderID = uid

		if _, err := r.validator.Validate(ctx, &data); err != nil {
			zap.L().Warn("gateway: payload validate error", zap.Error(err))
			return
		}

		_, _ = r.sendRoomMessageCase.Execute(ctx, &chatApp.SendRoomMessageInput{
			SenderID: data.SenderID,
			RoomID:   data.RoomID,
			Content:  data.Content,
		})

	default:
		zap.L().Warn("gateway: topic not found", zap.String("topic", req.Topic), zap.String("userID", userID))
	}
}
