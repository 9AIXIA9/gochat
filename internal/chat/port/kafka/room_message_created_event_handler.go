package kafka

import (
	"context"
	"gochat/internal/chat/application/usecase"
	"gochat/internal/chat/domain"
	"gochat/internal/shared/event"
)

func NewRoomMessageCreatedEventHandler(uc usecase.RoomMessageCreatedUseCase) event.HandlerFunc {
	return func(ctx context.Context, e event.Event) error {
		ev, err := domain.ToRoomMessageCreatedEvent(e)
		if err != nil {
			return err
		}

		input := usecase.RoomMessageCreatedInput{
			RoomID:    domain.RoomID(ev.AggregateID()),
			MessageID: ev.MessageID(),
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
