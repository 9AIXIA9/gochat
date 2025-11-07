package kafka

import (
	"context"
	"fmt"
	"gochat/internal/notification/domain"
	"gochat/internal/shared/event"
)

type RoomMessageCreatedEventHandler struct {
}

func NewRoomMessageCreatedEventHandler() event.Handler {
	return &RoomMessageCreatedEventHandler{}
}

func (h *RoomMessageCreatedEventHandler) Handle(ctx context.Context, e event.Event) error {
	ev, err := domain.ToRoomMessageCreatedEvent(e)
	if err != nil {
		return err
	}

	fmt.Println("Handled RoomMessageCreatedEvent:", ev)

	return nil
}
