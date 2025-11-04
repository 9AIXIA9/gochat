package kafka

import (
	"context"
	"gochat/internal/chat/domain"
	notificationDomain "gochat/internal/notification/domain"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"

	"go.uber.org/zap"
)

func NewMessageSentHandler(eventIDGenerator event.IDGenerator, publisher event.Publisher) event.Handler {
	return func(ctx context.Context, e event.Event) error {
		ev, err := domain.ToMessageSentEvent(e)
		if err != nil {
			return err
		}

		switch ev.Type() {
		case domain.MessageTypePrivate:
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

func handlePrivateMessageNotification(ev *domain.MessageSentEvent, eventIDGenerator event.IDGenerator, publisher event.Publisher) error {
	e, err := notificationDomain.NewMessageNotificationRequestedEvent(
		eventIDGenerator.Generate(),
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
