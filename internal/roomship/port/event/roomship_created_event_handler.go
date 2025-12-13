package event

import (
	"gochat/internal/roomship/application"
	"gochat/internal/roomship/domain"
	"gochat/internal/shared/event"
)

func NewRoomshipCreatedEventHandler(uc application.RoomshipCreatedUseCase) event.Handler {
	return event.AdaptUsecaseToHandler(
		uc,
		domain.ToRoomshipCreatedEvent,
		func(createdEvent *domain.RoomshipCreatedEvent) *application.RoomshipCreatedInput {
			return &application.RoomshipCreatedInput{
				RoomshipID: domain.RoomshipID(createdEvent.AggregateID()),
			}
		},
		nil,
		nil,
	)
}
