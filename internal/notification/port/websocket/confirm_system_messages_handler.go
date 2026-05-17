package websocket

import (
	"context"
	"encoding/json"
	"gochat/internal/infrastructure/websocket"
	notificationApp "gochat/internal/notification/application"
	"gochat/internal/shared/kernel"
	"gochat/pkg/ctxutil"
)

// ConfirmSystemMessagesTopic 客户端确认“系统消息已收到/已送达 Delivered”
// payload: {"message_ids": ["m1","m2"]}
const ConfirmSystemMessagesTopic websocket.Topic = "notification.confirm_system_messages"

type ConfirmSystemMessagesData struct {
	UserID     kernel.UserID      `json:"-" validate:"required"`
	MessageIDs []kernel.MessageID `json:"message_ids" validate:"required,min=1,dive,required"`
}

func NewConfirmSystemMessagesHandler(
	uc notificationApp.ConfirmSystemMessagesUseCase,
	validator websocket.Validator,
) websocket.Handler {
	return websocket.AdaptUsecaseToHandler(
		uc,
		validator,
		func(ctx context.Context, bytes []byte) (*ConfirmSystemMessagesData, error) {
			var data ConfirmSystemMessagesData
			if err := json.Unmarshal(bytes, &data); err != nil {
				return nil, err
			}
			data.UserID = ctxutil.UserIDFrom(ctx)
			return &data, nil
		},
		func(data *ConfirmSystemMessagesData) *notificationApp.ConfirmSystemMessagesInput {
			return &notificationApp.ConfirmSystemMessagesInput{
				UserID:     data.UserID,
				MessageIDs: data.MessageIDs,
			}
		},
	)
}
