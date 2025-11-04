package kafka

import (
	"context"
	"gochat/internal/notification/application/usecase"
	"gochat/internal/notification/domain"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

func NewEmailNotificationRequestedHandler(uc usecase.SendEmailUseCase) event.Handler {
	return func(ctx context.Context, e event.Event) error {
		ev, err := domain.ToEmailNotificationRequestedEvent(e)
		if err != nil {
			return err
		}

		input := &usecase.SendEmailInput{
			Recipient:      kernel.UserID(ev.AggregateID()),
			RecipientEmail: ev.Email(),
			Theme:          ev.Theme(),
			Title:          ev.Title(),
			Content:        ev.Content(),
		}

		if err := input.Validate(); err != nil {
			return err
		}
		_, err = uc.Execute(ctx, input)
		return err
	}
}
