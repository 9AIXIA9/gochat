package websocket

import (
	"context"
	"encoding/json"
	"gochat/internal/chat/application"
	"gochat/internal/infrastructure/websocket"
	"gochat/internal/shared/kernel"
	"gochat/pkg/ctxutil"
)

// ConfirmPrivateMessagesTopic 客户端确认“私聊消息已收到(Delivered)”
// payload: {"message_ids": ["m1","m2"]}
const ConfirmPrivateMessagesTopic websocket.Topic = "chat.confirm_private_messages"

type ConfirmPrivateMessagesData struct {
	UserID     kernel.UserID      `json:"-" validate:"required"`
	MessageIDs []kernel.MessageID `json:"message_ids" validate:"required,min=1,dive,required"`
}

func NewConfirmPrivateMessagesHandler(
	uc application.ConfirmPrivateMessagesUseCase,
	validator websocket.Validator,
) websocket.Handler {
	return websocket.AdaptUsecaseToHandler(
		uc,
		validator,
		func(ctx context.Context, bytes []byte) (*ConfirmPrivateMessagesData, error) {
			var data ConfirmPrivateMessagesData
			if err := json.Unmarshal(bytes, &data); err != nil {
				return nil, err
			}
			data.UserID = ctxutil.UserIDFrom(ctx)
			return &data, nil
		},
		func(data *ConfirmPrivateMessagesData) *application.ConfirmPrivateMessagesInput {
			return &application.ConfirmPrivateMessagesInput{
				UserID:     data.UserID,
				MessageIDs: data.MessageIDs,
			}
		},
	)
}
