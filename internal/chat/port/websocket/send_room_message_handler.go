package websocket

import (
	"context"
	"encoding/json"
	"gochat/internal/chat/application"
	"gochat/internal/infrastructure/websocket"
	"gochat/internal/shared/kernel"
	"gochat/pkg/utils"
)

const SendRoomMessageTopic websocket.Topic = "chat.send_room_message"

type SendRoomMessageData struct {
	SenderID kernel.UserID `json:"-" validate:"required"`
	RoomID   kernel.RoomID `json:"room_id" validate:"required"`
	Content  string        `json:"content" validate:"required,max=1000"`
}

func NewSendRoomMessageHandler(
	uc application.SendRoomMessageUseCase,
	validator websocket.Validator,
) websocket.Handler {
	return websocket.AdaptUsecaseToHandler(
		uc,
		validator,
		func(ctx context.Context, bytes []byte) (*SendRoomMessageData, error) {
			var data SendRoomMessageData
			if err := json.Unmarshal(bytes, &data); err != nil {
				return nil, err
			}
			data.SenderID = utils.GetUserID(ctx)
			return &data, nil
		},
		func(data *SendRoomMessageData) *application.SendRoomMessageInput {
			return &application.SendRoomMessageInput{
				SenderID: data.SenderID,
				RoomID:   data.RoomID,
				Content:  data.Content,
			}
		},
	)
}
