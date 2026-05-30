package command

import (
	"context"
	"gochat/internal/chat/application"
	"gochat/internal/chat/domain"
	"gochat/internal/gateway/core"
	"gochat/internal/shared/contract"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"

	"go.uber.org/zap"
)

func NewSendRoomMessageCommandHandler(
	uc application.SendRoomMessageUseCase,
	gateway contract.GatewayService,
) event.Handler {
	return event.AdaptUsecaseToCommandHandler(
		uc,
		domain.ToSendRoomMessageCommand,
		func(ev *domain.SendRoomMessageCommand) *application.SendRoomMessageInput {
			return &application.SendRoomMessageInput{
				SenderID: kernel.UserID(ev.AggregateID()),
				RoomID:   ev.RoomID(),
				Content:  ev.Content(),
			}
		},
		func(ctx context.Context, action *domain.SendRoomMessageCommand, output *kernel.NoOutput) {
			messageID := action.Headers()[KeyClientMessageID]
			envelop := core.NewSuccessEnvelop(kernel.MessageID(messageID), action.Topic(), nil)
			if err := gateway.PushToUser(ctx, kernel.UserID(action.AggregateID()), envelop); err != nil {
				zap.L().Error(
					"failed to push message to sender after sending private message successfully",
					zap.String("message_id", messageID),
					zap.Error(err),
				)
				return
			}
		},
		func(ctx context.Context, action *domain.SendRoomMessageCommand, err error) {
			messageID := action.Headers()[KeyClientMessageID]
			errorMessage := ""
			if myErrors.IsBusinessError(err) {
				errorMessage = err.Error()
			} else {
				errorMessage = myErrors.ErrServerBusy.Error()
			}
			envelop := core.NewFailedEnvelop(kernel.MessageID(messageID), action.Topic(), errorMessage)
			if err := gateway.PushToUser(ctx, kernel.UserID(action.AggregateID()), envelop); err != nil {
				zap.L().Error(
					"failed to push message to sender after sending private message successfully",
					zap.String("message_id", messageID),
					zap.Error(err),
				)
				return
			}
		},
	)
}
