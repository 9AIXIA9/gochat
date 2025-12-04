package event

import (
	"context"
	"gochat/internal/roomship/application"
	"gochat/internal/roomship/domain"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

func NewMemberRequestCreatedEventHandler(uc application.MemberRequestCreatedUseCase) event.HandlerFunc {
	return func(ctx context.Context, e event.Event) error {
		ev, err := domain.ToMemberRequestCreatedEvent(e)
		if err != nil {
			return err
		}

		input := application.MemberRequestCreatedInput{
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
