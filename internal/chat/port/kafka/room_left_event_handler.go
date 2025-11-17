package kafka

import (
	"context"
	"gochat/internal/chat/application/usecase"
	"gochat/internal/chat/domain"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

func NewRoomLeftEventHandler(uc usecase.RoomLeftUseCase) event.HandlerFunc {
	return func(ctx context.Context, e event.Event) error {
		ev, err := domain.ToRoomLeftEvent(e)
		if err != nil {
			return err
		}

		input := usecase.RoomLeftInput{
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
