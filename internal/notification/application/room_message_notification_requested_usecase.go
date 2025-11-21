package application

import (
	"context"
	"gochat/internal/notification/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
	"time"
)

type RoomMessageNotificationRequestedUseCase kernel.UseCase[*RoomMessageNotificationRequestedInput, *kernel.NoOutput]

type RoomMessageNotificationRequestedInput struct {
	MessageID    kernel.MessageID
	RoomID       kernel.RoomID
	RecipientIDs []kernel.UserID
	SenderID     kernel.UserID
	Content      string
	SentAt       time.Time
}

func (r *RoomMessageNotificationRequestedInput) Validate() error {
	if len(r.RecipientIDs) == 0 || len(r.RoomID) == 0 || len(r.SenderID) == 0 || len(r.MessageID) == 0 {
		return myErrors.ErrEmptyInput
	}

	if len(r.Content) == 0 {
		return myErrors.ErrEmptyInput
	}

	return nil
}

type roomMessageNotificationRequestedUseCase struct {
	messageSaver    domain.RoomMessageSaver
	messageNotifier domain.RoomMessageNotifier
}

func NewRoomMessageNotificationRequestedUseCase(
	messageNotifier domain.RoomMessageNotifier,
	messageSaver domain.RoomMessageSaver,
) RoomMessageNotificationRequestedUseCase {
	return &roomMessageNotificationRequestedUseCase{
		messageSaver:    messageSaver,
		messageNotifier: messageNotifier,
	}
}

func (uc *roomMessageNotificationRequestedUseCase) Execute(ctx context.Context, input *RoomMessageNotificationRequestedInput) (*kernel.NoOutput, error) {
	message := domain.ReceiveRoomMessage(
		input.MessageID,
		input.SenderID,
		input.RoomID,
		input.RecipientIDs,
		input.Content,
		input.SentAt,
	)

	if err := message.Deliver(
		uc.messageNotifier,
	); err != nil {
		return nil, err
	}

	if err := uc.messageSaver.SaveRoomMessage(ctx, message); err != nil {
		return nil, err
	}
	return nil, nil
}
