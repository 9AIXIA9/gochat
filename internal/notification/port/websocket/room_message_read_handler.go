package websocket

import (
	"context"
	"encoding/json"
	"gochat/internal/infrastructure/websocket"
	"gochat/internal/notification/application"
	"gochat/internal/shared/kernel"
	"gochat/pkg/utils"
)

const RoomMessageReadTopic websocket.Topic = "notification.room_message_read"

type RoomMessageReadData struct {
	MessageID kernel.MessageID `json:"message_id"`
}

func NewRoomMessageReadHandler(
	uc application.RoomMessageReadUseCase,
) websocket.HandlerFunc {
	return func(ctx context.Context, data []byte) ([]byte, error) {
		var reqData RoomMessageReadData
		if err := json.Unmarshal(data, &reqData); err != nil {
			return nil, err
		}

		input := &application.RoomMessageReadInput{
			MessageID: reqData.MessageID,
			UserID:    utils.GetUserIDFromCtx(ctx),
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
