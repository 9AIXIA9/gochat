package event

import (
	"gochat/internal/roomship/application"
	"gochat/internal/roomship/domain"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

func NewMemberRequestCreatedEventHandler(uc application.MemberRequestCreatedUseCase) event.Handler {
	return event.AdaptUsecaseToHandler(
		uc,
		domain.ToMemberRequestCreatedEvent,
		func(createdEvent *domain.MemberRequestCreatedEvent) *application.MemberRequestCreatedInput {
			return &application.MemberRequestCreatedInput{
				RequestID: kernel.OperationID(createdEvent.AggregateID()),
			}
		},
		nil,
		nil,
	)
}
