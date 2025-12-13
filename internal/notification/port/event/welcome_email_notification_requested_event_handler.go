package event

import (
	"gochat/internal/notification/application"
	"gochat/internal/notification/domain"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

func NewWelcomeEmailNotificationRequestedEventHandler(uc application.WelcomeEmailNotificationRequestedUseCase) event.Handler {
	return event.AdaptUsecaseToHandler(
		uc,
		domain.ToWelcomeEmailNotificationRequestedEvent,
		func(requestedEvent *domain.WelcomeEmailRequestedNotificationEvent) *application.WelcomeEmailNotificationRequestedInput {
			return &application.WelcomeEmailNotificationRequestedInput{
				UserID:     kernel.UserID(requestedEvent.AggregateID()),
				UserNumber: requestedEvent.Number(),
				Email:      requestedEvent.Email(),
			}
		},
		nil,
		nil,
	)
}
