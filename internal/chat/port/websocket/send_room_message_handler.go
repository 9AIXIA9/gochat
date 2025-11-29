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
	RoomNumber kernel.RoomNumber `json:"room_number" validate:"required"`
	Content    string            `json:"content" validate:"required,max=1000"`
}

func NewSendRoomMessageHandler(
	uc application.SendRoomMessageUseCase,
) websocket.HandlerFunc {
	return func(ctx context.Context, data []byte) ([]byte, error) {
		var reqData SendRoomMessageData
		if err := json.Unmarshal(data, &reqData); err != nil {
			return nil, err
		}

		input := &application.SendRoomMessageInput{
			SenderID:   utils.GetUserID(ctx),
			RoomNumber: reqData.RoomNumber,
			Content:    reqData.Content,
		}

		if err := input.Validate(); err != nil {
			return nil, err
		}

		if _, err := uc.Execute(ctx, input); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
