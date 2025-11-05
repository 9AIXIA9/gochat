package kafka

import (
	"context"
	"gochat/internal/chat/domain"
	notificationDomain "gochat/internal/notification/domain"
	"gochat/internal/shared/event"

	"go.uber.org/zap"
)

func NewMessageCreatedHandler(eventIDGenerator event.IDGenerator, publisher event.Publisher) event.Handler {
	return func(ctx context.Context, e event.Event) error {
		roomMessageCreatedEvent, err := domain.ToMessageCreatedEvent(e)
		if err != nil {
			return err
		}

		var lastErr error
		for _, recipient := range roomMessageCreatedEvent.Recipients() {
			ev, err := notificationDomain.NewMessageNotificationRequestedEvent(
				eventIDGenerator.Generate(),
				notificationDomain.MessageID(roomMessageCreatedEvent.MessageID()),
				roomMessageCreatedEvent.Sender(),
				recipient,
				roomMessageCreatedEvent.Content(),
				roomMessageCreatedEvent.SentAt(),
			)
			if err != nil {
				lastErr = err
				continue
			}

			if err := publisher.Publish(ev); err != nil {
				lastErr = err
			}
		}

		zap.L().Error(
			"room message created handler encountered errors",
			zap.Error(lastErr),
		)

		return nil
	}
}
