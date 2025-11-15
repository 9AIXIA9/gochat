package kafka

import (
	"context"
	"gochat/internal/notification/application/usecase"
	"gochat/internal/notification/domain"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

func NewUndeliveredMessageNotificationRequestedEventHandler(uc usecase.UndeliveredMessageNotificationRequestedUseCase) event.HandlerFunc {
	return func(ctx context.Context, e event.Event) error {
		ev, err := domain.ToUndeliveredMessageNotificationRequestedEvent(e)
		if err != nil {
			return err
		}

		input := usecase.UndeliveredMessageNotificationRequestedInput{
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
