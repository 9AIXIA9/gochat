package application

import (
	"context"
	"fmt"
	"gochat/internal/notification/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
)

type UndeliveredMessagesNotificationRequestedUseCase kernel.UseCase[*UndeliveredMessagesNotificationRequestedInput, *kernel.NoOutput]

type UndeliveredMessagesNotificationRequestedInput struct {
	UserID kernel.UserID
}

func (r *UndeliveredMessagesNotificationRequestedInput) Validate() error {
	if len(r.UserID) == 0 {
		return myErrors.ErrEmptyInput
	}
	return nil
}

type undeliveredMessagesNotificationRequestedUseCase struct {
	undeliveredPrivateMessageFinder domain.UndeliveredPrivateMessageFinder
	privateMessagesSaver            domain.PrivateMessagesSaver
	privateMessageNotifier          domain.PrivateMessageNotifier
	undeliveredRoomMessageFinder    domain.UndeliveredRoomMessageFinder
	roomMessagesSaver               domain.RoomMessagesSaver
	roomMessageNotifier             domain.RoomMessageNotifier
}

func NewUndeliveredMessagesNotificationRequestedUseCase(
	undeliveredPrivateMessageFinder domain.UndeliveredPrivateMessageFinder,
	privateMessagesSaver domain.PrivateMessagesSaver,
	privateMessageNotifier domain.PrivateMessageNotifier,
	undeliveredRoomMessageFinder domain.UndeliveredRoomMessageFinder,
	roomMessagesSaver domain.RoomMessagesSaver,
	roomMessageNotifier domain.RoomMessageNotifier,
) UndeliveredMessagesNotificationRequestedUseCase {
	return &undeliveredMessagesNotificationRequestedUseCase{
		undeliveredPrivateMessageFinder: undeliveredPrivateMessageFinder,
		privateMessagesSaver:            privateMessagesSaver,
		privateMessageNotifier:          privateMessageNotifier,
		undeliveredRoomMessageFinder:    undeliveredRoomMessageFinder,
		roomMessagesSaver:               roomMessagesSaver,
		roomMessageNotifier:             roomMessageNotifier,
	}
}

func (uc *undeliveredMessagesNotificationRequestedUseCase) Execute(ctx context.Context, input *UndeliveredMessagesNotificationRequestedInput) (*kernel.NoOutput, error) {
	privateMessages, err := uc.undeliveredPrivateMessageFinder.FindUndeliveredPrivateMessages(ctx, input.UserID)
	if err != nil {
		return nil, err
	}

	fmt.Println("Undelivered private messages found:", len(privateMessages))

	roomMessages, err := uc.undeliveredRoomMessageFinder.FindUndeliveredRoomMessages(ctx, input.UserID)
	if err != nil {
		return nil, err
	}

	fmt.Println("Undelivered room messages found:", len(roomMessages))

	if err := domain.NotifyUndeliveredMessages(
		input.UserID,
		privateMessages,
		roomMessages,
		uc.privateMessageNotifier,
		uc.roomMessageNotifier,
	); err != nil {
		return nil, err
	}

	if len(privateMessages) > 0 {
		if err := uc.privateMessagesSaver.SavePrivateMessages(ctx, privateMessages); err != nil {
			return nil, err
		}
	}

	if len(roomMessages) > 0 {
		if err := uc.roomMessagesSaver.SaveRoomMessages(ctx, roomMessages); err != nil {
			return nil, err
		}
	}

	return nil, nil
}
