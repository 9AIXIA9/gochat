package kafka

import (
	"context"
	"gochat/internal/chat/application/usecase"
	"gochat/internal/chat/domain"
	notificationDomain "gochat/internal/notification/domain"
	"gochat/internal/shared/event"
)

func NewMessageDeliveredHandler(uc usecase.UpdateMessageStateUseCase) event.Handler {
	return func(ctx context.Context, e event.Event) error {
		ev, err := notificationDomain.ToMessageDeliveredEvent(e)
		if err != nil {
			return err
		}
		input := &usecase.UpdateMessageStateInput{
			NewState:    domain.MessageStateDelivered,
			MessageID:   domain.MessageID(ev.MessageID()),
			RecipientID: ev.Recipient(),
		}

		if err := input.Validate(); err != nil {
			return err
		}

		if _, err := uc.Execute(ctx, input); err != nil {
			return err
		}

		return nil
	}
}
