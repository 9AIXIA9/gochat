package kafka

import (
	"context"
	"gochat/internal/shared/event"
	"gochat/internal/social/application/usecase"
	"gochat/internal/social/domain"
)

func NewRoomLeftEventHandler(uc usecase.RoomLeftUseCase) event.HandlerFunc {
	return func(ctx context.Context, e event.Event) error {
		ev, err := domain.ToRoomLeftEvent(e)
		if err != nil {
			return err
		}

		input := usecase.RoomLeftInput{
			RoomID: domain.RoomID(ev.AggregateID()),
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
