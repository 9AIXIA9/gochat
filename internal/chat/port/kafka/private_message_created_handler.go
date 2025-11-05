package kafka

import (
	"context"
	"gochat/internal/chat/domain"
	notificationDomain "gochat/internal/notification/domain"
	"gochat/internal/shared/event"
)

func NewPrivateMessageCreatedHandler(eventIDGenerator event.IDGenerator, publisher event.Publisher) event.Handler {
	return func(ctx context.Context, e event.Event) error {
		privateMessageCreatedEvent, err := domain.ToPrivateMessageCreatedEvent(e)
		if err != nil {
			return err
		}

		ev, err := notificationDomain.NewMessageNotificationRequestedEvent(
			eventIDGenerator.Generate(),
			notificationDomain.MessageID(privateMessageCreatedEvent.MessageID()),
			privateMessageCreatedEvent.Sender(),
			privateMessageCreatedEvent.Recipient(),
			privateMessageCreatedEvent.Content(),
			privateMessageCreatedEvent.SentAt(),
		)
		if err != nil {
			return err
		}

		if err := publisher.Publish(ev); err != nil {
			return err
		}
		return nil
	}
}
