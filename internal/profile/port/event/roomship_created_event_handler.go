package event

import (
	"gochat/internal/profile/application"
	"gochat/internal/profile/domain"
	"gochat/internal/shared/event"
)

func NewRoomshipCreatedEventHandler(uc application.RoomshipCreatedUseCase) event.Handler {
	return event.AdaptUsecaseToEventHandler(
		uc,
		domain.ToRoomshipCreatedEvent,
		func(createdEvent *domain.RoomshipCreatedEvent) *application.RoomshipCreatedInput {
			return &application.RoomshipCreatedInput{
				ID:     domain.RoomshipID(createdEvent.AggregateID()),
				UserID: createdEvent.UserID(),
				RoomID: createdEvent.RoomID(),
				Role:   createdEvent.Role(),
			}
		},
		nil,
		nil,
	)
}
