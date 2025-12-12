package event

import (
	"context"
	"gochat/internal/profile/application"
	"gochat/internal/profile/domain"
	"gochat/internal/shared/event"
)

func NewRoomshipCreatedEventHandler(uc application.RoomshipCreatedUseCase) event.HandlerFunc {
	return func(ctx context.Context, e event.Event) error {
		ev, err := domain.ToRoomshipCreatedEvent(e)
		if err != nil {
			return err
		}

		input := application.RoomshipCreatedInput{
			ID:     domain.RoomshipID(ev.AggregateID()),
			UserID: ev.UserID(),
			RoomID: ev.RoomID(),
			Role:   ev.Role(),
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
