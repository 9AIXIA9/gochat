package websocket

import (
	"context"
	"encoding/json"
	"gochat/internal/chat/application"
	"gochat/internal/infrastructure/websocket"
	"gochat/internal/shared/kernel"
	"gochat/pkg/utils"
)

const SendPrivateMessageTopic websocket.Topic = "chat.send_private_message"

type SendPrivateMessageData struct {
	RecipientNumber kernel.UserNumber `json:"recipient_number" validate:"required"`
	Content         string            `json:"content" validate:"required,max=1000"`
}

func NewSendPrivateMessageHandler(
	uc application.SendPrivateMessageUseCase,
) websocket.HandlerFunc {
	return func(ctx context.Context, data []byte) ([]byte, error) {
		var reqData SendPrivateMessageData
		if err := json.Unmarshal(data, &reqData); err != nil {
			return nil, err
		}

		input := &application.SendPrivateMessageInput{
			SenderID:        utils.GetUserID(ctx),
			RecipientNumber: reqData.RecipientNumber,
			Content:         reqData.Content,
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
