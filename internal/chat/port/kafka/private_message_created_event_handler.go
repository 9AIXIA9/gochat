package kafka

import (
	"context"
	chatDomain "gochat/internal/chat/domain"
	notificationDomain "gochat/internal/notification/domain"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

type PrivateMessageCreatedEventHandler struct {
	idGenerator event.IDGenerator
	publisher   event.Publisher
}

func NewPrivateMessageCreatedEventHandler(idGenerator event.IDGenerator, publisher event.Publisher) event.Handler {
	return &PrivateMessageCreatedEventHandler{idGenerator: idGenerator, publisher: publisher}
}

func (h *PrivateMessageCreatedEventHandler) Handle(_ context.Context, e event.Event) error {
	ev, err := chatDomain.ToPrivateMessageCreatedEvent(e)
	if err != nil {
		return err
	}

	notificationEv, err := notificationDomain.NewPrivateMessageCreatedEvent(
		h.idGenerator.Generate(),
		notificationDomain.MessageID(ev.MessageID()),
		kernel.UserID(ev.AggregateID()),
	)
	if err != nil {
		return err
	}

	if err := h.publisher.Publish(notificationEv); err != nil {
		return err
	}
	return nil
}
