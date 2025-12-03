package event

import (
	"context"
	"gochat/internal/friendship/application"
	"gochat/internal/friendship/domain"
	"gochat/internal/shared/event"
)

func NewFriendshipCreatedEventHandler(uc application.FriendshipCreatedUseCase) event.HandlerFunc {
	return func(ctx context.Context, e event.Event) error {
		ev, err := domain.ToFriendshipCreatedEvent(e)
		if err != nil {
			return err
		}

		input := application.FriendshipCreatedInput{
			FriendshipID: domain.FriendshipID(ev.AggregateID()),
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
