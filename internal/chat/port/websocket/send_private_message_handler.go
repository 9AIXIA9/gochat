package websocket

import (
	"context"
	"encoding/json"
	"gochat/internal/chat/application"
	"gochat/internal/infrastructure/websocket"
	"gochat/internal/shared/kernel"
	"gochat/pkg/ctxutil"
)

const SendPrivateMessageTopic websocket.Topic = "chat.send_private_message"

type SendPrivateMessageData struct {
	SenderID    kernel.UserID `json:"-" validate:"required"`
	RecipientID kernel.UserID `json:"recipient_id" validate:"required"`
	Content     string        `json:"content" validate:"required,max=1000"`
}

func NewSendPrivateMessageHandler(
	uc application.SendPrivateMessageUseCase,
	validator websocket.Validator,
) websocket.Handler {
	return websocket.AdaptUsecaseToHandler(
		uc,
		validator,
		func(ctx context.Context, bytes []byte) (*SendPrivateMessageData, error) {
			var data SendPrivateMessageData
			if err := json.Unmarshal(bytes, &data); err != nil {
				return nil, err
			}
			data.SenderID = ctxutil.UserIDFrom(ctx)
			return &data, nil
		},
		func(data *SendPrivateMessageData) *application.SendPrivateMessageInput {
			return &application.SendPrivateMessageInput{
				SenderID:    data.SenderID,
				RecipientID: data.RecipientID,
				Content:     data.Content,
			}
		},
	)
}
