package kafka

import (
	"context"
	"gochat/internal/notification/application/usecase"
	"gochat/internal/notification/domain"
	"gochat/internal/shared/event"
)

func NewMessageNotificationRequestedHandler(uc usecase.SendMessageUseCase, publisher event.Publisher, generator event.IDGenerator) event.Handler {
	return func(ctx context.Context, e event.Event) error {
		ev, err := domain.ToMessageNotificationRequestedEvent(e)
		if err != nil {
			return err
		}

		input := &usecase.SendMessageInput{
			MessageID: ev.MessageID(),
			Sender:    ev.Sender(),
			Recipient: ev.Recipient(),
			Content:   ev.Content(),
			SentAt:    ev.SentAt(),
		}

		if err := input.Validate(); err != nil {
			return err
		}

		if _, err := uc.Execute(ctx, input); err != nil {
			return err
		}

		messageDeliveredEvent, err := domain.NewMessageDeliveredEvent(generator.Generate(), input.MessageID, input.Recipient)
		if err != nil {
			return err
		}

		if err := publisher.Publish(messageDeliveredEvent); err != nil {
			return err
		}
		return nil
	}
}
