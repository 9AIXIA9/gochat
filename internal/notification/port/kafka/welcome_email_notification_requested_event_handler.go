package kafka

import (
	"context"
	"gochat/internal/notification/application"
	"gochat/internal/notification/domain"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

func NewWelcomeEmailNotificationRequestedEventHandler(uc application.WelcomeEmailNotificationRequestedUseCase) event.HandlerFunc {
	return func(ctx context.Context, e event.Event) error {
		ev, err := domain.ToWelcomeEmailNotificationRequestedEvent(e)
		if err != nil {
			return err
		}

		input := application.WelcomeEmailNotificationRequestedInput{
			UserID:     kernel.UserID(ev.AggregateID()),
			UserNumber: ev.Number(),
			Email:      ev.Email(),
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
