package kafka

import (
	"context"
	"gochat/internal/notification/application"
	"gochat/internal/notification/domain"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

func NewPrivateMessageNotificationRequestedEventHandler(uc application.PrivateMessageNotificationRequestedUseCase) event.HandlerFunc {
	return func(ctx context.Context, e event.Event) error {
		ev, err := domain.ToPrivateMessageNotificationRequestedEvent(e)
		if err != nil {
			return err
		}

		input := application.PrivateMessageNotificationRequestedInput{
			MessageID:   kernel.MessageID(ev.AggregateID()),
			RecipientID: ev.RecipientID(),
			SenderID:    ev.SenderID(),
			Content:     ev.Content(),
			SentAt:      ev.SentAt(),
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
