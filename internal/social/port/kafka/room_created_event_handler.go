package kafka

import (
	"context"
	"gochat/internal/shared/event"
	"gochat/internal/social/application/usecase"
	"gochat/internal/social/domain"
)

func NewRoomCreatedEventHandler(uc usecase.RoomCreatedUseCase) event.HandlerFunc {
	return func(ctx context.Context, e event.Event) error {
		ev, err := domain.ToRoomCreatedEvent(e)
		if err != nil {
			return err
		}

		input := usecase.RoomCreatedInput{
			RoomID: domain.RoomID(ev.AggregateID()),
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
