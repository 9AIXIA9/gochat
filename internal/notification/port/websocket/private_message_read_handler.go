package websocket

import (
	"context"
	"encoding/json"
	"gochat/internal/infrastructure/websocket"
	"gochat/internal/notification/application"
	"gochat/internal/shared/kernel"
)

const PrivateMessageReadRequestTopic websocket.RequestTopic = "notification.private_message_read"

type PrivateMessageReadRequestData struct {
	MessageID kernel.MessageID `json:"message_id"`
}

func NewPrivateMessageReadHandler(
	uc application.PrivateMessageReadUseCase,
) websocket.HandlerFunc {
	return func(ctx context.Context, request *websocket.Request) (*websocket.Response, error) {
		var reqData PrivateMessageReadRequestData
		if err := json.Unmarshal(request.Data, &reqData); err != nil {
			return nil, err
		}

		input := &application.PrivateMessageReadInput{
			MessageID: reqData.MessageID,
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
