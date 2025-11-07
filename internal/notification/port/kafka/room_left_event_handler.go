package kafka

import (
	"context"
	"gochat/internal/notification/domain"
	"gochat/internal/shared/event"
)

type RoomLeftEventHandler struct {
	roomMemberDeleter RoomMemberDeleter
}

func NewRoomLeftEventHandler(roomMemberDeleter RoomMemberDeleter) event.Handler {
	return &RoomLeftEventHandler{
		roomMemberDeleter: roomMemberDeleter,
	}
}

func (h *RoomLeftEventHandler) Handle(ctx context.Context, e event.Event) error {
	ev, err := domain.ToRoomLeftEvent(e)
	if err != nil {
		return err
	}

	return h.roomMemberDeleter.DeleteMember(ctx, domain.RoomID(ev.AggregateID()), ev.UserID())
}
