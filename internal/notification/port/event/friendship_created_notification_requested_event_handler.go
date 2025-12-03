package event

import (
	"context"
	"gochat/internal/notification/application"
	"gochat/internal/notification/domain"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

func NewFriendshipCreatedNotificationRequestedEventHandler(uc application.FriendshipCreatedNotificationRequestedUseCase) event.HandlerFunc {
	return func(ctx context.Context, e event.Event) error {
		ev, err := domain.ToFriendshipCreatedNotificationRequestedEvent(e)
		if err != nil {
			return err
		}

		input := application.FriendshipCreatedNotificationRequestedInput{
			UserID:    kernel.UserID(ev.AggregateID()),
			FriendID:  ev.FriendID(),
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
