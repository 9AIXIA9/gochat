package kafka

import (
	"context"
	"fmt"
	"gochat/internal/notification/domain"
	"gochat/internal/shared/event"
)

type PrivateMessageCreatedEventHandler struct {
}

func NewPrivateMessageCreatedEventHandler() event.Handler {
	return &PrivateMessageCreatedEventHandler{}
}

func (h *PrivateMessageCreatedEventHandler) Handle(ctx context.Context, e event.Event) error {
	ev, err := domain.ToPrivateMessageCreatedEvent(e)
	if err != nil {
		return err
	}

	fmt.Println("Handled PrivateMessageCreatedEvent:", ev)

	return nil
}
