package command

import (
	"context"
	"gochat/internal/chat/application"
	"gochat/internal/chat/domain"
	"gochat/internal/shared/contract"
	"gochat/internal/shared/kernel"

	"gochat/internal/shared/event"
)

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
			_ = gateway.PushToUser(ctx, kernel.UserID(action.AggregateID()), []byte(`{"type":"send_private_message_succeeded"}`))
		},
		func(ctx context.Context, action *domain.SendPrivateMessageCommand, err error) {
			_ = gateway.PushToUser(ctx, kernel.UserID(action.AggregateID()), []byte(`{"type":"send_private_message_failed"}`))
		},
	)
}
