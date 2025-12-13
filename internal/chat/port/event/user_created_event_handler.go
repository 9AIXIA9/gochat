package event

import (
	"gochat/internal/chat/application"
	"gochat/internal/chat/domain"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

func NewUserCreatedEventHandler(uc application.UserCreatedUseCase) event.Handler {
	return event.AdaptUsecaseToHandler(
		uc,
		domain.ToUserCreatedEvent,
		func(createdEvent *domain.UserCreatedEvent) *application.UserCreatedInput {
			return &application.UserCreatedInput{
				UserID: kernel.UserID(createdEvent.AggregateID()),
			}
		},
		nil,
		nil,
	)
}
