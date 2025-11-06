package kafka

import (
	"context"
	"gochat/internal/chat/domain"
	notificationDomain "gochat/internal/notification/domain"
	"gochat/internal/shared/event"
)

func NewPrivateMessageReceivedHandler(eventIDGenerator event.IDGenerator, publisher event.Publisher) event.Handler {
	return func(ctx context.Context, e event.Event) error {
		privateMessageReceivedEvent, err := domain.ToPrivateMessageReceivedEvent(e)
		if err != nil {
			return err
		}

		messageNotificationRequestedEvent, err := notificationDomain.NewMessageNotificationRequestedEvent(
			eventIDGenerator.Generate(),
			notificationDomain.MessageID(privateMessageReceivedEvent.MessageID()),
			privateMessageReceivedEvent.Sender(),
			privateMessageReceivedEvent.Recipient(),
			privateMessageReceivedEvent.Content(),
			privateMessageReceivedEvent.SentAt(),
		)
		if err != nil {
			return err
		}

		if err := publisher.Publish(messageNotificationRequestedEvent); err != nil {
			return err
		}

		return err
	}
}
