package event

import (
	"gochat/internal/notification/application"
	"gochat/internal/shared/contract"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

func NewPushSucceededEventHandler(uc application.PushSucceededUseCase) event.Handler {
	return event.AdaptUsecaseToEventHandler(
		uc,
		contract.ToPushSucceededEvent,
		func(createdEvent *contract.PushSucceededEvent) *application.PushSucceededInput {
			return &application.PushSucceededInput{
				ID: kernel.MessageID(createdEvent.AggregateID()),
			}
		},
		nil,
		nil,
	)
}
