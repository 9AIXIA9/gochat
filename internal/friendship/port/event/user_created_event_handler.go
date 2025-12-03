package event

import (
	"context"
	"gochat/internal/friendship/application"
	"gochat/internal/friendship/domain"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

func NewUserCreatedEventHandler(uc application.UserCreatedUseCase) event.HandlerFunc {
	return func(ctx context.Context, e event.Event) error {
		ev, err := domain.ToUserCreatedEvent(e)
		if err != nil {
			return err
		}

		input := application.UserCreatedInput{
			UserID: kernel.UserID(ev.AggregateID()),
		}

		if err := input.Validate(); err != nil {
			return err
		}

		if _, err := uc.Execute(ctx, &input); err != nil {
			return err
		}
		return nil
	}
}
