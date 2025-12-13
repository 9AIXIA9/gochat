package event

import (
	"gochat/internal/notification/application"
	"gochat/internal/notification/domain"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

func NewSystemMessageNotificationRequestedEventHandler(uc application.SystemMessageNotificationRequestedUseCase) event.Handler {
	return event.AdaptUsecaseToHandler(
		uc,
		domain.ToSystemMessageNotificationRequestedEvent,
		func(requestedEvent *domain.SystemMessageNotificationRequestedEvent) *application.SystemMessageNotificationRequestedInput {
			return &application.SystemMessageNotificationRequestedInput{
				RecipientID: kernel.UserID(requestedEvent.AggregateID()),
				Content:     requestedEvent.Content(),
			}
		},
		nil,
		nil,
	)
}
