package kafka

import (
	"context"
	"gochat/internal/notification/application"
	"gochat/internal/notification/domain"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

func NewReceivedPrivateMessageNotificationRequestedEventHandler(uc application.ReceivedPrivateMessageNotificationRequestedUseCase) event.HandlerFunc {
	return func(ctx context.Context, e event.Event) error {
		ev, err := domain.ToReceivedPrivateMessageNotificationRequestedEvent(e)
		if err != nil {
			return err
		}

		input := application.ReceivedPrivateMessageNotificationRequestedInput{
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
