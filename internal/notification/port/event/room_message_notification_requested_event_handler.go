package event

import (
	"context"
	"gochat/internal/notification/application"
	"gochat/internal/notification/domain"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

func NewRoomMessageNotificationRequestedEventHandler(uc application.RoomMessageNotificationRequestedUseCase) event.HandlerFunc {
	return func(ctx context.Context, e event.Event) error {
		ev, err := domain.ToRoomMessageNotificationRequestedEvent(e)
		if err != nil {
			return err
		}

		input := application.RoomMessageNotificationRequestedInput{
			MessageID:    kernel.MessageID(ev.AggregateID()),
			RoomID:       ev.RoomID(),
			RecipientIDs: ev.RecipientIDs(),
			SenderID:     ev.SenderID(),
			Content:      ev.Content(),
			SentAt:       ev.SentAt(),
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
