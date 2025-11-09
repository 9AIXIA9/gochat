package kafka

import (
	"context"
	"gochat/internal/notification/application/usecase"
	"gochat/internal/notification/domain"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

func NewUserCreatedEventHandler(uc usecase.UserCreatedUseCase) event.HandlerFunc {
	return func(ctx context.Context, e event.Event) error {
		ev, err := domain.ToUserCreatedEvent(e)
		if err != nil {
			return err
		}

		input := usecase.UserCreatedInput{
			UserID: kernel.UserID(ev.AggregateID()),
			Email:  ev.Email(),
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
