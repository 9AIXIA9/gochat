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

const SendRoomMessageTopic websocket.Topic = "chat.send_room_message"

type SendRoomMessageData struct {
	SenderID kernel.UserID `json:"-" validate:"required"`
	RoomID   kernel.RoomID `json:"room_id" validate:"required"`
	Content  string        `json:"content" validate:"required,max=1000"`
}

func NewSendRoomMessageHandler(
	uc application.SendRoomMessageUseCase,
	validator websocket.Validator,
) websocket.Handler {
	return websocket.AdaptUsecaseToHandler(
		uc,
		validator,
		func(ctx context.Context, bytes []byte) (*SendRoomMessageData, error) {
			var data SendRoomMessageData
			if err := json.Unmarshal(bytes, &data); err != nil {
				return nil, err
			}
			data.SenderID = utils.GetUserID(ctx)
			return &data, nil
		},
		func(data *SendRoomMessageData) *application.SendRoomMessageInput {
			return &application.SendRoomMessageInput{
				SenderID: data.SenderID,
				RoomID:   data.RoomID,
				Content:  data.Content,
			}
		},
		func(ctx context.Context, err error) *api.Response {
			switch {
			case errors.Is(err, myErrors.ErrNotFound):
				return api.NewResponseWithMessage(api.CodeNotFound, "room is not found")
			case errors.Is(err, myErrors.ErrInvalidLength):
				return api.NewResponseWithMessage(api.CodeInvalidParam, "content is too long")
			case errors.Is(err, domain.ErrNotMember):
				return api.NewResponseWithMessage(api.CodeInvalidParam, "you are not a member of the room")
			default:
				return api.NewResponse(api.CodeServerError)
			}
		},
	)
}
