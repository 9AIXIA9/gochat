package kafka

import (
	"context"
	"gochat/internal/notification/application/usecase"
	"gochat/internal/notification/domain"
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
			MessageID:   ev.MessageID(),
			RecipientID: kernel.UserID(ev.AggregateID()),
			SenderID:    ev.Sender(),
			Content:     ev.Content(),
			SentAt:      ev.SentAt(),
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
