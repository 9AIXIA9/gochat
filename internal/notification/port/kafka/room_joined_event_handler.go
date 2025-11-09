package kafka

import (
	"context"
	"gochat/internal/notification/application/usecase"
	"gochat/internal/notification/domain"
	"gochat/internal/shared/event"
)

func NewRoomJoinedEventHandler(uc usecase.RoomJoinedUseCase) event.HandlerFunc {
	return func(ctx context.Context, e event.Event) error {
		ev, err := domain.ToRoomJoinedEvent(e)
		if err != nil {
			return err
		}

		input := usecase.RoomJoinedInput{
			RoomID: domain.RoomID(e.AggregateID()),
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
