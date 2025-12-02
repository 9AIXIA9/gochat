package event

import (
	"context"
	"gochat/internal/friendship/application"
	"gochat/internal/friendship/domain"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

func NewFriendshipCreatedEventHandler(uc application.FriendshipCreatedUseCase) event.HandlerFunc {
	return func(ctx context.Context, e event.Event) error {
		ev, err := domain.ToFriendshipCreatedEvent(e)
		if err != nil {
			return err
		}

		input := application.FriendshipCreatedInput{
			RequestID: kernel.OperationID(ev.AggregateID()),
			UserID1:   ev.UserID1(),
			UserID2:   ev.UserID2(),
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
