package kafka

import (
	"context"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"gochat/internal/social/domain"
)

type UserCreatedEventHandler struct {
	userIDSaver UserIDSaver
}

func NewUserCreatedEventHandler(userIDSaver UserIDSaver) event.Handler {
	return &UserCreatedEventHandler{
		userIDSaver: userIDSaver,
	}
}

func (h *UserCreatedEventHandler) Handle(ctx context.Context, e event.Event) error {
	ev, err := domain.ToUserCreatedEvent(e)
	if err != nil {
		return err
	}

	return h.userIDSaver.SaveID(ctx, kernel.UserID(ev.AggregateID()))
}
