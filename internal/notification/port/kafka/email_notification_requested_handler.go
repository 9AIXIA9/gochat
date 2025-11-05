package kafka

import (
	"context"
	"gochat/internal/infrastructure/kafka"
	"gochat/internal/notification/application/usecase"
	"gochat/internal/notification/domain"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

func NewEmailNotificationRequestedHandler(uc usecase.SendEmailUseCase) event.Handler {
	return kafka.AdaptUseCaseToHandler(
		uc,
		func(e event.Event) (*usecase.SendEmailInput, error) {
			ev, err := domain.ToEmailNotificationRequestedEvent(e)
			if err != nil {
				return nil, err
			}

			return &usecase.SendEmailInput{
				Recipient:      kernel.UserID(ev.AggregateID()),
				RecipientEmail: ev.Email(),
				Theme:          ev.Theme(),
				Title:          ev.Title(),
				Content:        ev.Content(),
			}, nil
		},
		func(ctx context.Context, output *kernel.NoOutput) error {
			return nil
		},
		func(ctx context.Context, err error) error {
			return err
		},
	)
}
