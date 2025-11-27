package websocket

import (
	"context"
	"encoding/json"
	"gochat/internal/infrastructure/websocket"
	"gochat/internal/notification/application"
	"gochat/internal/shared/kernel"
	"gochat/pkg/utils"
)

const ReadRoomMessageTopic websocket.Topic = "notification.read_room_message"

type ReadRoomMessageData struct {
	MessageID kernel.MessageID `json:"message_id"`
}

func NewReadRoomMessageHandler(
	uc application.ReadRoomMessageUseCase,
) websocket.HandlerFunc {
	return func(ctx context.Context, data []byte) ([]byte, error) {
		var reqData ReadRoomMessageData
		if err := json.Unmarshal(data, &reqData); err != nil {
			return nil, err
		}

		input := &application.ReadRoomMessageInput{
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
