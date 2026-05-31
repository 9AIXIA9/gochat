package event

import (
	"gochat/internal/roomship/application"
	"gochat/internal/roomship/domain"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

func NewMemberRequestAgreedEventHandler(uc application.MemberRequestAgreedUseCase) event.Handler {
	return event.AdaptUsecaseToEventHandler(
		uc,
		domain.ToMemberRequestAgreedEvent,
		func(agreedEvent *domain.MemberRequestAgreedEvent) *application.MemberRequestAgreedInput {
			return &application.MemberRequestAgreedInput{
				RequestID: kernel.OperationID(agreedEvent.AggregateID()),
			}
		},
	)
}
