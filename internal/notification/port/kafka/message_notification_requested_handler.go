package kafka

import (
	"context"
	"gochat/internal/notification/application/usecase"
	"gochat/internal/notification/domain"
	"gochat/internal/shared/event"
)

//TODO 可以做一个通用的事件处理器生成器

func NewMessageNotificationRequestedHandler(uc usecase.SendMessageUseCase) event.Handler {
	return func(ctx context.Context, e event.Event) error {
		ev, err := domain.ToMessageNotificationRequestedEvent(e)
		if err != nil {
			return err
		}

		input := &usecase.SendMessageInput{
			Sender:    ev.Sender(),
			Recipient: ev.Recipient(),
			Content:   ev.Content(),
			SentAt:    ev.SentAt(),
		}

		if err := input.Validate(); err != nil {
			return err
		}

		_, err = uc.Execute(ctx, input)
		return err
	}
}
