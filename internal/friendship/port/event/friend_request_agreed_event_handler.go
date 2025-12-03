package event

import (
	"context"
	"gochat/internal/friendship/application"
	"gochat/internal/friendship/domain"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

func NewFriendRequestAgreedEventHandler(uc application.FriendRequestAgreedUseCase) event.HandlerFunc {
	return func(ctx context.Context, e event.Event) error {
		ev, err := domain.ToFriendRequestAgreedEvent(e)
		if err != nil {
			return err
		}

		input := application.FriendRequestAgreedInput{
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
