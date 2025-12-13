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
	SenderID kernel.UserID `json:"sender_id" validate:"required"`
}

func NewReadPrivateMessagesHandler(
	uc application.ReadPrivateMessagesUseCase,
) websocket.Handler {
	return websocket.AdaptUsecaseToHandler(
		uc,
		func(ctx context.Context, data []byte) (*application.ReadPrivateMessagesInput, error) {
			var reqData ReadPrivateMessagesData
			if err := json.Unmarshal(data, &reqData); err != nil {
				return nil, err
			}
			return &application.ReadPrivateMessagesInput{
				SenderID:    reqData.SenderID,
				RecipientID: utils.GetUserID(ctx),
			}, nil
		},
		nil,
		nil,
	)
}
