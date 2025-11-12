package kafka

import (
	"context"
	"gochat/internal/notification/application/usecase"
	"gochat/internal/notification/domain"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

func NewUserConnectedEventHandler(uc usecase.UserConnectedUseCase) event.HandlerFunc {
	return func(ctx context.Context, e event.Event) error {
		ev, err := domain.ToUserConnectedEvent(e)
		if err != nil {
			return err
		}

		input := usecase.UserConnectedInput{
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
