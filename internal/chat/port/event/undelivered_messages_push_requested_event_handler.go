package event

import (
	"context"
	"gochat/internal/chat/application"
	"gochat/internal/chat/domain"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

func NewUndeliveredMessagesPushRequestedEventHandler(uc application.UndeliveredMessagesPushRequestedUseCase) event.HandlerFunc {
	return func(ctx context.Context, e event.Event) error {
		ev, err := domain.ToUndeliveredMessagesPushRequestedEvent(e)
		if err != nil {
			return err
		}

		input := application.UndeliveredMessagesPushRequestedInput{
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
