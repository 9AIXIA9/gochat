package command

import (
	"context"
	"gochat/internal/chat/application"
	"gochat/internal/chat/domain"
	"gochat/internal/gateway/core"
	"gochat/internal/shared/command"
	"gochat/internal/shared/contract"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"

	"go.uber.org/zap"
)

func NewSendRoomMessageCommandHandler(
	uc application.SendRoomMessageUseCase,
	gateway contract.GatewayService,
) command.Handler {
	return command.AdaptUsecaseToCommandHandler(
		uc,
		domain.ToSendRoomMessageCommand,
		func(com *domain.SendRoomMessageCommand) *application.SendRoomMessageInput {
			return &application.SendRoomMessageInput{
				SenderID: kernel.UserID(com.AggregateID()),
				RoomID:   com.RoomID(),
				Content:  com.Content(),
			}
		},
		func(ctx context.Context, action *domain.SendRoomMessageCommand, output *kernel.NoOutput) {
			messageID := action.Headers()[KeyClientMessageID]
			envelop := core.NewActionSucceededDownstreamEnvelop(kernel.MessageID(messageID), action.Action(), nil)
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
			envelop := core.NewActionFailedDownstreamEnvelop(kernel.MessageID(messageID), action.Action(), errorMessage)
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
