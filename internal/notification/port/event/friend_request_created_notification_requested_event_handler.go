package event

import (
	"context"
	"gochat/internal/notification/application"
	"gochat/internal/notification/domain"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

func NewFriendRequestCreatedNotificationRequestedEventHandler(uc application.FriendRequestCreatedNotificationRequestedUseCase) event.HandlerFunc {
	return func(ctx context.Context, e event.Event) error {
		ev, err := domain.ToFriendRequestCreatedNotificationRequestedEvent(e)
		if err != nil {
			return err
		}

		input := application.FriendRequestCreatedNotificationRequestedInput{
			RequestID: kernel.OperationID(ev.AggregateID()),
			From:      ev.From(),
			To:        ev.To(),
			SentAt:    ev.SentAt(),
			Content:   ev.Content(),
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
