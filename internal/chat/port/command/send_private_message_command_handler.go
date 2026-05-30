package command

import (
	"context"
	"gochat/internal/chat/application"
	"gochat/internal/chat/domain"
	"gochat/internal/shared/kernel"

	"gochat/internal/shared/event"
)

func NewSendPrivateMessageCommandHandler(
	uc application.SendPrivateMessageUseCase,
	publisher event.SyncPublisher,
	generator event.IDGenerator,
) event.Handler {
	return event.AdaptUsecaseToCommandHandler(
		uc,
		publisher,
		domain.ToSendPrivateMessageCommand,
		func(ev *domain.SendPrivateMessageCommand) *application.SendPrivateMessageInput {
			return &application.SendPrivateMessageInput{
				SenderID:    kernel.UserID(ev.AggregateID()),
				RecipientID: ev.RecipientID(),
				Content:     ev.Content(),
			}
		},
		func(ctx context.Context, action *domain.SendPrivateMessageCommand, output *kernel.NoOutput) event.SpecificEvent {
			ev, _ := event.NewStandardCommandSucceedEvent(action, generator)
			return ev
		},
		func(ctx context.Context, action *domain.SendPrivateMessageCommand, err error) event.SpecificEvent {
			ev, _ := event.NewStandardCommandFailedEvent(action, err, generator)
			return ev
		},
	)
}
