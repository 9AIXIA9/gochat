package websocket

import (
	"context"
	"encoding/json"
	"gochat/internal/chat/application"
	"gochat/internal/infrastructure/websocket"
	"gochat/internal/shared/kernel"
	"gochat/pkg/utils"
)

const ReadRoomMessagesTopic websocket.Topic = "chat.read_room_messages"

type ReadRoomMessagesData struct {
	RoomID  kernel.RoomID `json:"room_id" validate:"required"`
	Content string        `json:"content" validate:"required,max=1000"`
}

func NewReadRoomMessagesHandler(
	uc application.ReadRoomMessagesUseCase,
) websocket.HandlerFunc {
	return func(ctx context.Context, data []byte) ([]byte, error) {
		var reqData ReadRoomMessagesData
		if err := json.Unmarshal(data, &reqData); err != nil {
			return nil, err
		}

		input := &application.ReadRoomMessagesInput{
			UserID: utils.GetUserID(ctx),
			RoomID: reqData.RoomID,
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
