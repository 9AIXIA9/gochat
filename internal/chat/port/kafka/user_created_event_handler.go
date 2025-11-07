package kafka

import (
	"context"
	chatDomain "gochat/internal/chat/domain"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

type UserCreatedEventHandler struct {
	userNumberSaver UserNumberSaver
}

func NewUserCreatedEventHandler(userNumberSaver UserNumberSaver) event.Handler {
	return &UserCreatedEventHandler{
		userNumberSaver: userNumberSaver,
	}
}

func (h *UserCreatedEventHandler) Handle(ctx context.Context, e event.Event) error {
	ev, err := chatDomain.ToUserCreatedEvent(e)
	if err != nil {
		return err
	}

	return h.userNumberSaver.SaveNumber(ctx, kernel.UserID(ev.AggregateID()), ev.Number())
}
