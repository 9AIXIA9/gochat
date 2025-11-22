package kafka

import (
	"context"
	"gochat/internal/notification/application"
	"gochat/internal/notification/domain"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

func NewUndeliveredMessagesNotificationRequestedEventHandler(uc application.UndeliveredMessagesNotificationRequestedUseCase) event.HandlerFunc {
	return func(ctx context.Context, e event.Event) error {
		ev, err := domain.ToUndeliveredMessagesNotificationRequestedEvent(e)
		if err != nil {
			return err
		}

		input := application.UndeliveredMessagesNotificationRequestedInput{
			UserID: kernel.UserID(ev.AggregateID()),
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
