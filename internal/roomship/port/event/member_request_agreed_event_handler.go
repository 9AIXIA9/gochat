package event

import (
	"context"
	"gochat/internal/roomship/application"
	"gochat/internal/roomship/domain"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

func NewMemberRequestAgreedEventHandler(uc application.MemberRequestAgreedUseCase) event.HandlerFunc {
	return func(ctx context.Context, e event.Event) error {
		ev, err := domain.ToMemberRequestAgreedEvent(e)
		if err != nil {
			return err
		}

		input := application.MemberRequestAgreedInput{
			RequestID: kernel.OperationID(ev.AggregateID()),
		}

		if err := input.Validate(); err != nil {
			return err
		}

		if _, err := uc.Execute(ctx, &input); err != nil {
			return err
		}
		return nil
	}
}
