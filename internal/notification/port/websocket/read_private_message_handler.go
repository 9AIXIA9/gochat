package websocket

import (
	"context"
	"encoding/json"
	"gochat/internal/infrastructure/websocket"
	"gochat/internal/notification/application"
	"gochat/internal/shared/kernel"
)

const ReadPrivateMessageTopic websocket.Topic = "notification.read_private_message"

type ReadPrivateMessageData struct {
	MessageID kernel.MessageID `json:"message_id"`
}

func NewReadPrivateMessageHandler(
	uc application.ReadPrivateMessageUseCase,
) websocket.HandlerFunc {
	return func(ctx context.Context, data []byte) ([]byte, error) {
		var reqData ReadPrivateMessageData
		if err := json.Unmarshal(data, &reqData); err != nil {
			return nil, err
		}

		input := &application.ReadPrivateMessageInput{
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
