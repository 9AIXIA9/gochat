package command

import (
	"context"
	"gochat/internal/chat/application"
	"gochat/internal/chat/domain"
	"gochat/internal/gateway/core"
	"gochat/internal/shared/contract"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"

	"gochat/internal/shared/command"

	"go.uber.org/zap"
)

const KeyClientMessageID = "client_message_id"

func NewSendPrivateMessageCommandHandler(
	uc application.SendPrivateMessageUseCase,
	gateway contract.GatewayService,
) command.Handler {
	return command.AdaptUsecaseToCommandHandler(
		uc,
		domain.ToSendPrivateMessageCommand,
		func(com *domain.SendPrivateMessageCommand) *application.SendPrivateMessageInput {
			return &application.SendPrivateMessageInput{
				SenderID:    kernel.UserID(com.AggregateID()),
				RecipientID: com.RecipientID(),
				Content:     com.Content(),
			}
		},
		func(ctx context.Context, action *domain.SendPrivateMessageCommand, output *kernel.NoOutput) {
			messageID := action.Headers()[KeyClientMessageID]
			envelop := core.NewActionSucceededDownstreamEnvelop(kernel.MessageID(messageID), action.Action(), nil)
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
			envelop := core.NewActionFailedDownstreamEnvelop(kernel.MessageID(messageID), action.Action(), errorMessage)
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
