package websocket

import (
	"context"
	"encoding/json"
	"gochat/internal/chat/application"
	"gochat/internal/infrastructure/websocket"
	"gochat/internal/shared/kernel"
	"gochat/pkg/utils"
)

const ReadPrivateMessagesTopic websocket.Topic = "chat.read_private_messages"

type ReadPrivateMessagesData struct {
	SenderID    kernel.UserID `json:"sender_id" validate:"required"`
	RecipientID kernel.UserID `json:"-" validate:"required"`
}

func NewReadPrivateMessagesHandler(
	uc application.ReadPrivateMessagesUseCase,
	validator websocket.Validator,
) websocket.Handler {
	return websocket.AdaptUsecaseToHandler(
		uc,
		validator,
		func(ctx context.Context, bytes []byte) (*ReadPrivateMessagesData, error) {
			var data ReadPrivateMessagesData
			if err := json.Unmarshal(bytes, &data); err != nil {
				return nil, err
			}
			data.RecipientID = utils.GetUserID(ctx)
			return &data, nil
		},
		func(data *ReadPrivateMessagesData) *application.ReadPrivateMessagesInput {
			return &application.ReadPrivateMessagesInput{
				SenderID:    data.SenderID,
				RecipientID: data.RecipientID,
			}
		},
	)
}
