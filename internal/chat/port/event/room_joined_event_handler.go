package event

import (
	"context"
	"gochat/internal/chat/application"
	"gochat/internal/chat/domain"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

func NewRoomJoinedEventHandler(uc application.RoomJoinedUseCase) event.HandlerFunc {
	return func(ctx context.Context, e event.Event) error {
		ev, err := domain.ToRoomJoinedEvent(e)
		if err != nil {
			return err
		}

		input := application.RoomJoinedInput{
			RoomID: kernel.RoomID(e.AggregateID()),
			UserID: ev.UserID(),
		}

		if err := input.Validate(); err != nil {
			return err
		}

		if _, err := uc.Execute(ctx, &input); err != nil {
			return err
		}
		return nil
	}
}
