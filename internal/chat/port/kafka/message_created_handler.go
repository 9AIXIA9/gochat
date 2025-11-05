package kafka

import (
	"context"
	"gochat/internal/chat/domain"
	notificationDomain "gochat/internal/notification/domain"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"

	"go.uber.org/zap"
)

func NewMessageCreatedHandler(eventIDGenerator event.IDGenerator, publisher event.Publisher) event.Handler {
	return func(ctx context.Context, e event.Event) error {
		ev, err := domain.ToMessageCreatedEvent(e)
		if err != nil {
			return err
		}

		switch ev.Type() {
		case domain.PrivateType:
			return handlePrivateMessageNotification(ev, eventIDGenerator, publisher)
		default:
			zap.L().Error(
				"unsupported message type for notification",
				zap.String("type", string(ev.Type())),
			)
			return nil
		}
	}
}

func handlePrivateMessageNotification(ev *domain.MessageCreatedEvent, eventIDGenerator event.IDGenerator, publisher event.Publisher) error {
	e, err := notificationDomain.NewMessageNotificationRequestedEvent(
		eventIDGenerator.Generate(),
		notificationDomain.MessageID(ev.MessageID()),
		ev.Sender(),
		kernel.UserID(ev.Recipient()),
		ev.Content(),
		ev.SentAt(),
	)
	if err != nil {
		return err
	}

	if err := publisher.Publish(e); err != nil {
		return err
	}
	return nil
}
