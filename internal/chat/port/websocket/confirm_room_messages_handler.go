package websocket

import (
	"context"
	"encoding/json"
	"gochat/internal/chat/application"
	"gochat/internal/infrastructure/websocket"
	"gochat/internal/shared/kernel"
	"gochat/pkg/ctxutil"
)

// ConfirmRoomMessagesTopic 客户端确认“群聊消息已收到(Delivered)”
// payload: {"message_ids": ["m1","m2"]}
const ConfirmRoomMessagesTopic websocket.Topic = "chat.confirm_room_messages"

type ConfirmRoomMessagesData struct {
	UserID     kernel.UserID      `json:"-" validate:"required"`
	MessageIDs []kernel.MessageID `json:"message_ids" validate:"required,min=1,dive,required"`
}

func NewConfirmRoomMessagesHandler(
	uc application.ConfirmRoomMessagesUseCase,
	validator websocket.Validator,
) websocket.Handler {
	return websocket.AdaptUsecaseToHandler(
		uc,
		validator,
		func(ctx context.Context, bytes []byte) (*ConfirmRoomMessagesData, error) {
			var data ConfirmRoomMessagesData
			if err := json.Unmarshal(bytes, &data); err != nil {
				return nil, err
			}
			data.UserID = ctxutil.UserIDFrom(ctx)
			return &data, nil
		},
		func(data *ConfirmRoomMessagesData) *application.ConfirmRoomMessagesInput {
			return &application.ConfirmRoomMessagesInput{
				UserID:     data.UserID,
				MessageIDs: data.MessageIDs,
			}
		},
	)
}
