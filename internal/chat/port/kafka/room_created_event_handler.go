package kafka

import (
	"context"
	"gochat/internal/chat/domain"
	"gochat/internal/shared/event"
)

type RoomCreatedEventHandler struct {
	roomNumberSaver RoomNumberSaver
}

func NewRoomCreatedEventHandler(roomNumberSaver RoomNumberSaver) event.Handler {
	return &RoomCreatedEventHandler{
		roomNumberSaver: roomNumberSaver,
	}
}

func (h *RoomCreatedEventHandler) Handle(ctx context.Context, e event.Event) error {
	ev, err := domain.ToRoomCreatedEvent(e)
	if err != nil {
		return err
	}

	return h.roomNumberSaver.SaveNumber(ctx, domain.RoomID(ev.AggregateID()), ev.Number())
}
