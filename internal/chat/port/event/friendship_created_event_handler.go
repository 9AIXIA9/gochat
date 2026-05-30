package event

import (
	"gochat/internal/chat/application"
	"gochat/internal/chat/domain"
	"gochat/internal/shared/event"
)

func NewFriendshipCreatedEventHandler(uc application.FriendshipCreatedUseCase) event.Handler {
	return event.AdaptUsecaseToEventHandler(
		uc,
		domain.ToFriendshipCreatedEvent,
		func(createdEvent *domain.FriendshipCreatedEvent) *application.FriendshipCreatedInput {
			return &application.FriendshipCreatedInput{
				ID:      domain.FriendshipID(createdEvent.AggregateID()),
				UserID1: createdEvent.UserID1(),
				UserID2: createdEvent.UserID2(),
			}
		},
		nil,
		nil,
	)
}
