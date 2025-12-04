package event

import (
	"context"
	"gochat/internal/notification/application"
	"gochat/internal/notification/domain"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

func NewSystemMessageNotificationRequestedEventHandler(uc application.SystemMessageNotificationRequestedUseCase) event.HandlerFunc {
	return func(ctx context.Context, e event.Event) error {
		ev, err := domain.ToSystemMessageNotificationRequestedEvent(e)
		if err != nil {
			return err
		}

		input := application.SystemMessageNotificationRequestedInput{
			RecipientID: kernel.UserID(ev.AggregateID()),
			Content:     ev.Content(),
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
