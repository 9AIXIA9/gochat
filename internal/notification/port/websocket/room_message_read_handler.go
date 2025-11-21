package websocket

import (
	"context"
	"encoding/json"
	"gochat/internal/infrastructure/websocket"
	"gochat/internal/notification/application"
	"gochat/internal/shared/kernel"
	"gochat/pkg/utils"
)

const RoomMessageReadRequestTopic websocket.RequestTopic = "notification.room_message_read"

type RoomMessageReadRequestData struct {
	MessageID kernel.MessageID `json:"message_id"`
}

func NewRoomMessageReadHandler(
	uc application.RoomMessageReadUseCase,
) websocket.HandlerFunc {
	return func(ctx context.Context, request *websocket.Request) (*websocket.Response, error) {
		var reqData RoomMessageReadRequestData
		if err := json.Unmarshal(request.Data, &reqData); err != nil {
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
