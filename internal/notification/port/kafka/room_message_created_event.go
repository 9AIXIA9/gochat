package kafka

import (
	"context"
	"gochat/internal/notification/application/usecase"
	"gochat/internal/notification/domain"
	"gochat/internal/shared/event"
)

func NewRoomMessageCreatedEventHandler(uc usecase.RoomMessageCreatedUseCase) event.HandlerFunc {
	return func(ctx context.Context, e event.Event) error {
		_, err := domain.ToRoomMessageCreatedEvent(e)
		if err != nil {
			return err
		}

		input := usecase.RoomMessageCreatedInput{}

		if err := input.Validate(); err != nil {
			return err
		}

		if _, err := uc.Execute(ctx, &input); err != nil {
			return err
		}
		return nil
	}
}
