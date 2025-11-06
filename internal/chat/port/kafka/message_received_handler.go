package kafka

import (
	"context"
	"gochat/internal/chat/domain"
	notificationDomain "gochat/internal/notification/domain"
	"gochat/internal/shared/event"
)

func NewMessageReceivedHandler(eventIDGenerator event.IDGenerator, publisher event.Publisher) event.Handler {
	return func(ctx context.Context, e event.Event) error {
		messageReceivedEvent, err := domain.ToMessageReceivedEvent(e)
		if err != nil {
			return err
		}

		messageNotificationRequestedEvent, err := notificationDomain.NewMessageNotificationRequestedEvent(
			eventIDGenerator.Generate(),
			notificationDomain.MessageID(messageReceivedEvent.MessageID()),
			messageReceivedEvent.Sender(),
			messageReceivedEvent.Recipient(),
			messageReceivedEvent.Content(),
			messageReceivedEvent.SentAt(),
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
