package kafka

import (
	"context"
	chatDomain "gochat/internal/chat/domain"
	notificationDomain "gochat/internal/notification/domain"
	"gochat/internal/shared/event"
	socialDomain "gochat/internal/social/domain"
)

type RoomJoinedEventHandler struct {
	idGenerator event.IDGenerator
	publisher   event.Publisher
}

func NewRoomJoinedEventHandler(idGenerator event.IDGenerator, publisher event.Publisher) event.Handler {
	return &RoomJoinedEventHandler{idGenerator: idGenerator, publisher: publisher}
}

func (h *RoomJoinedEventHandler) Handle(_ context.Context, e event.Event) error {
	ev, err := socialDomain.ToRoomJoinedEvent(e)
	if err != nil {
		return err
	}

	notificationEv, err := notificationDomain.NewRoomJoinedEvent(h.idGenerator.Generate(), notificationDomain.RoomID(ev.AggregateID()), ev.UserID())
	if err != nil {
		return err
	}
	if err := h.publisher.Publish(notificationEv); err != nil {
		return err
	}

	chatEv, err := chatDomain.NewRoomJoinedEvent(h.idGenerator.Generate(), chatDomain.RoomID(ev.AggregateID()), ev.UserID())
	if err != nil {
		return err
	}
	if err := h.publisher.Publish(chatEv); err != nil {
		return err
	}
	return nil
}
