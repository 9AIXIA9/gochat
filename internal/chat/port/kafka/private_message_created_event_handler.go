package kafka

import (
	"context"
	"gochat/internal/chat/application/usecase"
	"gochat/internal/chat/domain"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

func NewPrivateMessageCreatedEventHandler(uc usecase.PrivateMessageCreatedUseCase) event.HandlerFunc {
	return func(ctx context.Context, e event.Event) error {
		ev, err := domain.ToPrivateMessageCreatedEvent(e)
		if err != nil {
			return err
		}

		input := usecase.PrivateMessageCreatedInput{
			RecipientID: kernel.UserID(ev.AggregateID()),
			MessageID:   ev.MessageID(),
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
