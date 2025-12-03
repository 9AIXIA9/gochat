package event

import (
	"context"
	"gochat/internal/friendship/application"
	"gochat/internal/friendship/domain"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

func NewFriendRequestCreatedEventHandler(uc application.FriendRequestCreatedUseCase) event.HandlerFunc {
	return func(ctx context.Context, e event.Event) error {
		ev, err := domain.ToFriendRequestCreatedEvent(e)
		if err != nil {
			return err
		}

		input := application.FriendRequestCreatedInput{
			RequestID: kernel.OperationID(ev.AggregateID()),
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
