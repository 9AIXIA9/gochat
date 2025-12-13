package event

import (
	"gochat/internal/notification/application"
	"gochat/internal/notification/domain"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

func NewUndeliveredMessagesNotificationRequestedEventHandler(uc application.UndeliveredMessagesNotificationRequestedUseCase) event.Handler {
	return event.AdaptUsecaseToHandler(
		uc,
		domain.ToUndeliveredMessagesNotificationRequestedEvent,
		func(requestedEvent *domain.UndeliveredMessagesNotificationRequestedEvent) *application.UndeliveredMessagesNotificationRequestedInput {
			return &application.UndeliveredMessagesNotificationRequestedInput{
				UserID: kernel.UserID(requestedEvent.AggregateID()),
			}
		},
		nil,
		nil,
	)
}
