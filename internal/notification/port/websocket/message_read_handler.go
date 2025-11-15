package websocket

import (
	"context"
	"encoding/json"
	"gochat/internal/infrastructure/websocket"
	"gochat/internal/notification/application/usecase"
	"gochat/internal/notification/domain"
	"gochat/pkg/utils"
)

const MessageReadRequestTopic websocket.RequestTopic = "notification.message_read"

type MessageReadRequestData struct {
	MessageID domain.MessageID `json:"message_id"`
}

func NewMessageReadHandler(
	uc usecase.MessageReadUseCase,
) websocket.HandlerFunc {
	return func(ctx context.Context, request *websocket.Request) (*websocket.Response, error) {
		var reqData MessageReadRequestData
		if err := json.Unmarshal(request.Data, &reqData); err != nil {
			return nil, err
		}

		input := &usecase.MessageReadInput{
			UserID:    utils.GetUserIDFromCtx(ctx),
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
