package kafka

import (
	"context"
	chatDomain "gochat/internal/chat/domain"
	notificationDomain "gochat/internal/notification/domain"
	"gochat/internal/shared/event"
	socialDomain "gochat/internal/social/domain"
)

type RoomCreatedEventHandler struct {
	idGenerator    event.IDGenerator
	publisher      event.Publisher
	roomFinderByID RoomFinderByID
}

func NewRoomCreatedEventHandler(idGenerator event.IDGenerator, publisher event.Publisher, roomFinderByID RoomFinderByID) event.Handler {
	return &RoomCreatedEventHandler{
		idGenerator:    idGenerator,
		publisher:      publisher,
		roomFinderByID: roomFinderByID,
	}
}

func (h *RoomCreatedEventHandler) Handle(ctx context.Context, e event.Event) error {
	ev, err := socialDomain.ToRoomCreatedEvent(e)
	if err != nil {
		return err
	}

	room, err := h.roomFinderByID.FindByID(ctx, socialDomain.RoomID(ev.AggregateID()))
	if err != nil {
		return err
	}

	notificationEv, err := notificationDomain.NewRoomCreatedEvent(h.idGenerator.Generate(), notificationDomain.RoomID(ev.AggregateID()))
	if err != nil {
		return err
	}
	if err := h.publisher.Publish(notificationEv); err != nil {
		return err
	}

	chatEv, err := chatDomain.NewRoomCreatedEvent(h.idGenerator.Generate(), chatDomain.RoomID(ev.AggregateID()), chatDomain.RoomNumber(room.Number()))
	if err != nil {
		return err
	}
	if err := h.publisher.Publish(chatEv); err != nil {
		return err
	}
	return nil
}
