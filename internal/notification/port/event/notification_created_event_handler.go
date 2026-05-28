package event

import (
	"gochat/internal/notification/application"
	"gochat/internal/shared/contract"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

func NewNotificationCreatedEventHandler(uc application.NotificationCreatedUseCase) event.Handler {
	return event.AdaptUsecaseToHandler(
		uc,
		contract.ToNotificationCreatedEvent,
		func(createdEvent *contract.NotificationCreatedEvent) *application.NotificationCreatedInput {
			return &application.NotificationCreatedInput{
				ID:          kernel.MessageID(createdEvent.AggregateID()),
				RecipientID: createdEvent.RecipientID(),
				RawPayload:  createdEvent.RawPayload(),
			}
		},
		nil,
		nil,
	)
}
