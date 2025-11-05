package kafka

import (
	"context"
	"gochat/internal/chat/application/usecase"
	"gochat/internal/chat/domain"
	"gochat/internal/infrastructure/kafka"
	notificationDomain "gochat/internal/notification/domain"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

func NewMessageDeliveredHandler(uc usecase.UpdateMessageStateUseCase) event.Handler {
	return kafka.AdaptUseCaseToHandler(
		uc,
		func(e event.Event) (*usecase.UpdateMessageStateInput, error) {
			ev, err := notificationDomain.ToMessageDeliveredEvent(e)
			if err != nil {
				return nil, err
			}
			return &usecase.UpdateMessageStateInput{
				NewState:    domain.MessageStateDelivered,
				MessageID:   domain.MessageID(ev.MessageID()),
				RecipientID: ev.Recipient(),
			}, nil
		},
		func(ctx context.Context, output *kernel.NoOutput) error {
			return nil
		},
		func(ctx context.Context, err error) error {
			return err
		},
	)
}
