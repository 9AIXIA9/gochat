package kafka

import (
	"context"
	"gochat/internal/notification/domain"
	"gochat/internal/shared/event"
)

type RoomCreatedEventHandler struct {
	roomIDSaver RoomIDSaver
}

func NewRoomCreatedEventHandler(roomIDSaver RoomIDSaver) event.Handler {
	return &RoomCreatedEventHandler{
		roomIDSaver: roomIDSaver,
	}
}

func (h *RoomCreatedEventHandler) Handle(ctx context.Context, e event.Event) error {
	ev, err := domain.ToRoomCreatedEvent(e)
	if err != nil {
		return err
	}

	return h.roomIDSaver.SaveID(ctx, domain.RoomID(ev.AggregateID()))
}
