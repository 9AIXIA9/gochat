package kafka

import (
	"context"
	chatDomain "gochat/internal/chat/domain"
	notificationDomain "gochat/internal/notification/domain"
	"gochat/internal/shared/event"
)

type RoomMessageCreatedEventHandler struct {
	idGenerator event.IDGenerator
	publisher   event.Publisher
}

func NewRoomMessageCreatedEventHandler(idGenerator event.IDGenerator, publisher event.Publisher) event.Handler {
	return &RoomMessageCreatedEventHandler{idGenerator: idGenerator, publisher: publisher}
}

func (h *RoomMessageCreatedEventHandler) Handle(_ context.Context, e event.Event) error {
	ev, err := chatDomain.ToRoomMessageCreatedEvent(e)
	if err != nil {
		return err
	}

	notificationEv, err := notificationDomain.NewRoomMessageCreatedEvent(
		h.idGenerator.Generate(),
		notificationDomain.MessageID(ev.MessageID()),
		notificationDomain.RoomID(ev.AggregateID()),
	)
	if err != nil {
		return err
	}

	if err := h.publisher.Publish(notificationEv); err != nil {
		return err
	}
	return nil
}
