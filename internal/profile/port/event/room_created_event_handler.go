package event

import (
	"context"
	"gochat/internal/profile/application"
	"gochat/internal/profile/domain"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

func NewRoomCreatedEventHandler(uc application.RoomCreatedUseCase) event.HandlerFunc {
	return func(ctx context.Context, e event.Event) error {
		ev, err := domain.ToRoomCreatedEvent(e)
		if err != nil {
			return err
		}

		input := application.RoomCreatedInput{
			RoomID:    kernel.RoomID(ev.AggregateID()),
			CreatedAt: ev.CreatedAt(),
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
