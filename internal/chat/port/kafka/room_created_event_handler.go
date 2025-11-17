package kafka

import (
	"context"
	"gochat/internal/chat/application/usecase"
	"gochat/internal/chat/domain"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

func NewRoomCreatedEventHandler(uc usecase.RoomCreatedUseCase) event.HandlerFunc {
	return func(ctx context.Context, e event.Event) error {
		ev, err := domain.ToRoomCreatedEvent(e)
		if err != nil {
			return err
		}

		input := usecase.RoomCreatedInput{
			RoomID:     kernel.RoomID(ev.AggregateID()),
			RoomNumber: ev.Number(),
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
