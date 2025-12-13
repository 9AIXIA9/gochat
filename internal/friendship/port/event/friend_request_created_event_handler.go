package event

import (
	"gochat/internal/friendship/application"
	"gochat/internal/friendship/domain"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

func NewFriendRequestCreatedEventHandler(uc application.FriendRequestCreatedUseCase) event.Handler {
	return event.AdaptUsecaseToHandler(
		uc,
		domain.ToFriendRequestCreatedEvent,
		func(createdEvent *domain.FriendRequestCreatedEvent) *application.FriendRequestCreatedInput {
			return &application.FriendRequestCreatedInput{
				RequestID: kernel.OperationID(createdEvent.AggregateID()),
			}
		},
		nil,
		nil,
	)
}
