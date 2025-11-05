package kafka

import (
	"context"
	"gochat/internal/infrastructure/kafka"
	"gochat/internal/notification/application/usecase"
	"gochat/internal/notification/domain"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

//TODO 可以做一个通用的事件处理器生成器

func NewMessageNotificationRequestedHandler(uc usecase.SendMessageUseCase) event.Handler {
	return kafka.AdaptUseCaseToHandler(
		uc,
		func(e event.Event) (*usecase.SendMessageInput, error) {
			ev, err := domain.ToMessageNotificationRequestedEvent(e)
			if err != nil {
				return nil, err
			}

			return &usecase.SendMessageInput{
				MessageID: ev.MessageID(),
				Sender:    ev.Sender(),
				Recipient: ev.Recipient(),
				Content:   ev.Content(),
				SentAt:    ev.SentAt(),
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
