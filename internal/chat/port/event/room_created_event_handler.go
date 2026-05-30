package event

import (
	"gochat/internal/chat/application"
	"gochat/internal/chat/domain"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

func NewRoomCreatedEventHandler(uc application.RoomCreatedUseCase) event.Handler {
	return event.AdaptUsecaseToEventHandler(
		uc,
		domain.ToRoomCreatedEvent,
		func(createdEvent *domain.RoomCreatedEvent) *application.RoomCreatedInput {
			return &application.RoomCreatedInput{
				RoomID: kernel.RoomID(createdEvent.AggregateID()),
			}
		},
		nil,
		nil,
	)
}
