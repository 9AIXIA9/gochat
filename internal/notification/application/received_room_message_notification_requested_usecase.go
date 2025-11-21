package application

import (
	"context"
	"gochat/internal/notification/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
)

type ReceivedRoomMessageNotificationRequestedUseCase kernel.UseCase[*ReceivedRoomMessageNotificationRequestedInput, *kernel.NoOutput]

type ReceivedRoomMessageNotificationRequestedInput struct {
	UserID kernel.UserID
}

func (r *ReceivedRoomMessageNotificationRequestedInput) Validate() error {
	if len(r.UserID) == 0 {
		return myErrors.ErrEmptyInput
	}

	return nil
}

type receivedRoomMessageNotificationRequestedUseCase struct {
	roomMessageFinder   domain.ReceivedRoomMessageFinder
	roomMessageSaver    domain.RoomMessageSaver
	roomMessageNotifier domain.RoomMessageNotifier
}

func NewReceivedRoomMessageNotificationRequestedUseCase(
	roomMessageFinder domain.ReceivedRoomMessageFinder,
	roomMessageSaver domain.RoomMessageSaver,
	roomMessageNotifier domain.RoomMessageNotifier,
) ReceivedRoomMessageNotificationRequestedUseCase {
	return &receivedRoomMessageNotificationRequestedUseCase{
		roomMessageFinder:   roomMessageFinder,
		roomMessageSaver:    roomMessageSaver,
		roomMessageNotifier: roomMessageNotifier,
	}
}

func (uc *receivedRoomMessageNotificationRequestedUseCase) Execute(ctx context.Context, input *ReceivedRoomMessageNotificationRequestedInput) (*kernel.NoOutput, error) {
	roomMessages, err := uc.roomMessageFinder.FindReceivedRoomMessage(ctx, input.UserID)
	if err != nil {
		return nil, err
	}

	//TODO 这只是单个人上线 不需要通知所有人

	for _, msg := range roomMessages {
		if err := msg.Deliver(uc.roomMessageNotifier); err != nil {
			return nil, err
		}
	}

	return nil, nil
}
