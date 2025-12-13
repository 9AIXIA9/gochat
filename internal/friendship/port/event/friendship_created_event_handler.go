package event

import (
	"gochat/internal/friendship/application"
	"gochat/internal/friendship/domain"
	"gochat/internal/shared/event"
)

func NewFriendshipCreatedEventHandler(uc application.FriendshipCreatedUseCase) event.Handler {
	return event.AdaptUsecaseToHandler(
		uc,
		domain.ToFriendshipCreatedEvent,
		func(createdEvent *domain.FriendshipCreatedEvent) *application.FriendshipCreatedInput {
			return &application.FriendshipCreatedInput{
				FriendshipID: domain.FriendshipID(createdEvent.AggregateID()),
			}
		},
		nil,
		nil,
	)
}
