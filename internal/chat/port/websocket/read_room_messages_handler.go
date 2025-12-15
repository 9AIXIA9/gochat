package websocket

import (
	"context"
	"encoding/json"
	"errors"
	"gochat/internal/chat/application"
	"gochat/internal/infrastructure/websocket"
	"gochat/internal/shared/api"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
	"gochat/pkg/utils"
)

const ReadRoomMessagesTopic websocket.Topic = "chat.read_room_messages"

type ReadRoomMessagesData struct {
	RoomID kernel.RoomID `json:"room_id" validate:"required"`
	UserID kernel.UserID `json:"-" validate:"required"`
}

func NewReadRoomMessagesHandler(
	uc application.ReadRoomMessagesUseCase,
	validator websocket.Validator,
) websocket.Handler {
	return websocket.AdaptUsecaseToHandler(
		uc,
		validator,
		func(ctx context.Context, bytes []byte) (*ReadRoomMessagesData, error) {
			var data ReadRoomMessagesData
			if err := json.Unmarshal(bytes, &data); err != nil {
				return nil, err
			}
			data.UserID = utils.GetUserID(ctx)
			return &data, nil
		},
		func(data *ReadRoomMessagesData) *application.ReadRoomMessagesInput {
			return &application.ReadRoomMessagesInput{
				UserID: data.UserID,
				RoomID: data.RoomID,
			}
		},
		func(ctx context.Context, err error) *api.Response {
			switch {
			case errors.Is(err, myErrors.ErrNotFound):
				return api.NewResponseWithMessage(api.CodeNotFound, "room is not found")
			default:
				return api.NewResponse(api.CodeServerError)
			}
		},
	)
}
