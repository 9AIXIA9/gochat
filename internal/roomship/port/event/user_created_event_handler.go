package event

import (
	"gochat/internal/roomship/application"
	"gochat/internal/roomship/domain"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

func NewUserCreatedEventHandler(uc application.UserCreatedUseCase) event.Handler {
	return event.AdaptUsecaseToEventHandler(
		uc,
		domain.ToUserCreatedEvent,
		func(createdEvent *domain.UserCreatedEvent) *application.UserCreatedInput {
			return &application.UserCreatedInput{
				UserID: kernel.UserID(createdEvent.AggregateID()),
			}
		},
	)
}
