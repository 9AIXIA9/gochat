package websocket

import (
	"context"
	"encoding/json"
	"errors"
	"gochat/internal/chat/application"
	"gochat/internal/chat/domain"
	"gochat/internal/infrastructure/websocket"
	"gochat/internal/shared/api"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
	"gochat/pkg/utils"
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
			data.SenderID = utils.GetUserID(ctx)
			return &data, nil
		},
		func(data *SendPrivateMessageData) *application.SendPrivateMessageInput {
			return &application.SendPrivateMessageInput{
				SenderID:    data.SenderID,
				RecipientID: data.RecipientID,
				Content:     data.Content,
			}
		},
		func(ctx context.Context, err error) *api.Response {
			switch {
			case errors.Is(err, myErrors.ErrNotFound):
				return api.NewResponseWithMessage(api.CodeNotFound, "recipient is not found")
			case errors.Is(err, myErrors.ErrInvalidLength):
				return api.NewResponseWithMessage(api.CodeInvalidParam, "content is too long")
			case errors.Is(err, domain.ErrNotFriends):
				return api.NewResponseWithMessage(api.CodeInvalidParam, "you are not friends with the recipient")
			default:
				return api.NewResponse(api.CodeServerError)
			}
		},
	)
}
