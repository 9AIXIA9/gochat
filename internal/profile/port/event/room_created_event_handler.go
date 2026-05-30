package event

import (
	"gochat/internal/profile/application"
	"gochat/internal/profile/domain"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

func NewRoomCreatedEventHandler(uc application.RoomCreatedUseCase) event.Handler {
	return event.AdaptUsecaseToEventHandler(
		uc,
		domain.ToRoomCreatedEvent,
		func(createdEvent *domain.RoomCreatedEvent) *application.RoomCreatedInput {
			return &application.RoomCreatedInput{
				RoomID:    kernel.RoomID(createdEvent.AggregateID()),
				CreatedAt: createdEvent.CreatedAt(),
			}
		},
		nil,
		nil,
	)
}
