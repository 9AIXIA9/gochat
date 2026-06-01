package event

import (
	"gochat/internal/profile/application"
	"gochat/internal/profile/domain"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

func NewUserCreatedEventHandler(uc application.UserCreatedUseCase) event.Handler {
	return event.AdaptUsecaseToEventHandler(
		uc,
		domain.ToUserCreatedEvent,
		func(createdEvent *domain.UserCreatedEvent) *application.UserCreatedInput {
			return &application.UserCreatedInput{
				UserID:     kernel.UserID(createdEvent.AggregateID()),
				Email:      createdEvent.Email(),
				SignedUpAt: createdEvent.SignedAt(),
			}
		},
	)
}
