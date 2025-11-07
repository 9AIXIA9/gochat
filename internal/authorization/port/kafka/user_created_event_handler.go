package kafka

import (
	"context"
	authorizationDomain "gochat/internal/authorization/domain"
	chatDomain "gochat/internal/chat/domain"
	notificationDomain "gochat/internal/notification/domain"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	socialDomain "gochat/internal/social/domain"
)

type UserCreatedEventHandler struct {
	idGenerator event.IDGenerator
	publisher   event.Publisher
	finder      UserFinderByID
}

func NewUserCreatedEventHandler(idGenerator event.IDGenerator, publisher event.Publisher, finder UserFinderByID) event.Handler {
	return &UserCreatedEventHandler{idGenerator: idGenerator, publisher: publisher, finder: finder}
}

func (h *UserCreatedEventHandler) Handle(ctx context.Context, e event.Event) error {
	ev, err := authorizationDomain.ToUserCreatedEvent(e)
	if err != nil {
		return err
	}

	user, err := h.finder.FindByID(ctx, kernel.UserID(ev.AggregateID()))
	if err != nil {
		return err
	}

	socialEv, err := socialDomain.NewUserCreatedEvent(h.idGenerator.Generate(), user.ID())
	if err != nil {
		return err
	}

	if err := h.publisher.Publish(socialEv); err != nil {
		return err
	}

	chatEv, err := chatDomain.NewUserCreatedEvent(h.idGenerator.Generate(), user.ID(), chatDomain.UserNumber(user.Number()))
	if err != nil {
		return err
	}

	if err := h.publisher.Publish(chatEv); err != nil {
		return err
	}

	notificationEv, err := notificationDomain.NewUserCreatedEvent(h.idGenerator.Generate(), user.ID(), user.Email())
	if err != nil {
		return err
	}

	if err := h.publisher.Publish(notificationEv); err != nil {
		return err
	}
	return nil
}
