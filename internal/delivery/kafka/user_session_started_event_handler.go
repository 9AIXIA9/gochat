package kafka

import (
	"context"
	"gochat/internal/application"
	"gochat/internal/infrastructure/websocket"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

func NewUserSessionStartedEventHandler(uc application.UserSessionStartedUseCase) event.HandlerFunc {
	return func(ctx context.Context, e event.Event) error {
		ev, err := websocket.ToUserSessionStartedEvent(e)
		if err != nil {
			return err
		}

		input := application.UserSessionStartedInput{
			UserID: kernel.UserID(ev.AggregateID()),
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
