package websocket

import (
	"context"
	"encoding/json"
	"gochat/internal/infrastructure/websocket"
	"gochat/internal/notification/application"
	"gochat/internal/shared/kernel"
)

//TODO 改为行为而不是事件

const PrivateMessageReadTopic websocket.Topic = "notification.private_message_read"

type PrivateMessageReadData struct {
	MessageID kernel.MessageID `json:"message_id"`
}

func NewPrivateMessageReadHandler(
	uc application.PrivateMessageReadUseCase,
) websocket.HandlerFunc {
	return func(ctx context.Context, data []byte) ([]byte, error) {
		var reqData PrivateMessageReadData
		if err := json.Unmarshal(data, &reqData); err != nil {
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
