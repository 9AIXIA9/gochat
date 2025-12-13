package event

import (
	"gochat/internal/chat/application"
	"gochat/internal/chat/domain"
	"gochat/internal/shared/event"
)

func NewRoomshipCreatedEventHandler(uc application.RoomshipCreatedUseCase) event.Handler {
	return event.AdaptUsecaseToHandler(
		uc,
		domain.ToRoomshipCreatedEvent,
		func(createdEvent *domain.RoomshipCreatedEvent) *application.RoomshipCreatedInput {
			return &application.RoomshipCreatedInput{
				ID:     domain.RoomshipID(createdEvent.AggregateID()),
				UserID: createdEvent.UserID(),
				RoomID: createdEvent.RoomID(),
			}
		},
		nil,
		nil,
	)
}
