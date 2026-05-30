package command

import (
	"context"
	"gochat/internal/chat/application"
	"gochat/internal/chat/domain"
	"gochat/internal/gateway/core"
	"gochat/internal/shared/contract"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"

	"gochat/internal/shared/event"

	"go.uber.org/zap"
)

const KeyClientMessageID = "client_message_id"

func NewSendPrivateMessageCommandHandler(
	uc application.SendPrivateMessageUseCase,
	gateway contract.GatewayService,
) event.Handler {
	return event.AdaptUsecaseToCommandHandler(
		uc,
		domain.ToSendPrivateMessageCommand,
		func(ev *domain.SendPrivateMessageCommand) *application.SendPrivateMessageInput {
			return &application.SendPrivateMessageInput{
				SenderID:    kernel.UserID(ev.AggregateID()),
				RecipientID: ev.RecipientID(),
				Content:     ev.Content(),
			}
		},
		func(ctx context.Context, action *domain.SendPrivateMessageCommand, output *kernel.NoOutput) {
			messageID := action.Headers()[KeyClientMessageID]
			envelop := core.NewSuccessEnvelop(kernel.MessageID(messageID), action.Topic(), nil)
			if err := gateway.PushToUser(ctx, kernel.UserID(action.AggregateID()), envelop); err != nil {
				zap.L().Error(
					"failed to push message to sender after sending private message successfully",
					zap.String("client_message_id", messageID),
					zap.Error(err),
				)
				return
			}
		},
		func(ctx context.Context, action *domain.SendPrivateMessageCommand, err error) {
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
					zap.String("client_message_id", messageID),
					zap.Error(err),
				)
				return
			}
		},
	)
}
