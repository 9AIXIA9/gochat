package event

import (
	"gochat/internal/friendship/application"
	"gochat/internal/friendship/domain"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

func NewFriendRequestAgreedEventHandler(uc application.FriendRequestAgreedUseCase) event.Handler {
	return event.AdaptUsecaseToEventHandler(
		uc,
		domain.ToFriendRequestAgreedEvent,
		func(agreedEvent *domain.FriendRequestAgreedEvent) *application.FriendRequestAgreedInput {
			return &application.FriendRequestAgreedInput{
				RequestID: kernel.OperationID(agreedEvent.AggregateID()),
			}
		},
		nil,
		nil,
	)
}
