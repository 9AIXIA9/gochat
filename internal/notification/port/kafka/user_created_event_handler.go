package kafka

import (
	"context"
	"gochat/internal/notification/domain"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

type UserCreatedEventHandler struct {
	userEmailSaver UserEmailSaver
}

func NewUserCreatedEventHandler(userEmailSaver UserEmailSaver) event.Handler {
	return &UserCreatedEventHandler{
		userEmailSaver: userEmailSaver,
	}
}

func (h *UserCreatedEventHandler) Handle(ctx context.Context, e event.Event) error {
	ev, err := domain.ToUserCreatedEvent(e)
	if err != nil {
		return err
	}

	return h.userEmailSaver.SaveEmail(ctx, kernel.UserID(ev.AggregateID()), ev.Email())
}
