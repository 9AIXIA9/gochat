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
	RecipientID kernel.UserID `json:"recipient_id" validate:"required"`
	Content     string        `json:"content" validate:"required,max=1000"`
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
			SenderID:    utils.GetUserID(ctx),
			RecipientID: reqData.RecipientID,
			Content:     reqData.Content,
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
