package command

import (
	"context"
	"gochat/internal/chat/application"
	"gochat/internal/chat/domain"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

func NewSendRoomMessageCommandHandler(
	uc application.SendRoomMessageUseCase,
	publisher event.SyncPublisher,
	generator event.IDGenerator,
) event.Handler {
	return event.AdaptUsecaseToCommandHandler(
		uc,
		publisher,
		domain.ToSendRoomMessageCommand,
		func(ev *domain.SendRoomMessageCommand) *application.SendRoomMessageInput {
			return &application.SendRoomMessageInput{
				SenderID: kernel.UserID(ev.AggregateID()),
				RoomID:   ev.RoomID(),
				Content:  ev.Content(),
			}
		},
		func(ctx context.Context, action *domain.SendRoomMessageCommand, output *kernel.NoOutput) event.SpecificEvent {
			ev, _ := event.NewStandardCommandSucceedEvent(action, generator)
			return ev
		},
		func(ctx context.Context, action *domain.SendRoomMessageCommand, err error) event.SpecificEvent {
			ev, _ := event.NewStandardCommandFailedEvent(action, err, generator)
			return ev
		},
	)
}
