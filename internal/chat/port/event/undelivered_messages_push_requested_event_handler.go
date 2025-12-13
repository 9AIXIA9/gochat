package event

import (
	"gochat/internal/chat/application"
	"gochat/internal/chat/domain"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

func NewUndeliveredMessagesPushRequestedEventHandler(uc application.UndeliveredMessagesPushRequestedUseCase) event.Handler {
	return event.AdaptUsecaseToHandler(
		uc,
		domain.ToUndeliveredMessagesPushRequestedEvent,
		func(requestedEvent *domain.UndeliveredMessagesPushRequestedEvent) *application.UndeliveredMessagesPushRequestedInput {
			return &application.UndeliveredMessagesPushRequestedInput{
				UserID: kernel.UserID(requestedEvent.AggregateID()),
			}
		},
		nil,
		nil,
	)
}
