package kafka

import (
	"context"
	"gochat/internal/chat/domain"
	"gochat/internal/shared/event"
)

type RoomJoinedEventHandler struct {
	roomMemberSaver RoomMemberSaver
}

func NewRoomJoinedEventHandler(roomMemberSaver RoomMemberSaver) event.Handler {
	return &RoomJoinedEventHandler{
		roomMemberSaver: roomMemberSaver,
	}
}

func (h *RoomJoinedEventHandler) Handle(ctx context.Context, e event.Event) error {
	ev, err := domain.ToRoomJoinedEvent(e)
	if err != nil {
		return err
	}

	return h.roomMemberSaver.SaveMember(ctx, domain.RoomID(ev.AggregateID()), ev.UserID())
}
