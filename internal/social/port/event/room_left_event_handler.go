package event

import (
	"context"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"gochat/internal/social/application"
	"gochat/internal/social/domain"
)

func NewRoomLeftEventHandler(uc application.RoomLeftUseCase) event.HandlerFunc {
	return func(ctx context.Context, e event.Event) error {
		ev, err := domain.ToRoomLeftEvent(e)
		if err != nil {
			return err
		}

		input := application.RoomLeftInput{
			RoomID: kernel.RoomID(ev.AggregateID()),
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
