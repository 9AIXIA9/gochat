package event

import (
	"gochat/internal/notification/application"
	"gochat/internal/shared/contract"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

func NewNotificationIntentCreatedEventHandler(uc application.NotificationIntentCreatedUseCase) event.Handler {
	return event.AdaptUsecaseToEventHandler(
		uc,
		contract.ToNotificationIntentCreatedEvent,
		func(createdEvent *contract.NotificationIntentCreatedEvent) *application.NotificationIntentCreatedInput {
			return &application.NotificationIntentCreatedInput{
				ID:          kernel.MessageID(createdEvent.AggregateID()),
				RecipientID: createdEvent.RecipientID(),
				RawPayload:  createdEvent.RawPayload(),
			}
		},
	)
}
