package application

import (
	"context"
	"gochat/internal/notification/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
)

//TODO 添加系统消息通知未送达处理逻辑

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
	userPrivateMessagesFinderByState domain.UserPrivateMessagesFinderByState
	privateMessagesUpdater           domain.PrivateMessagesUpdater
	privateMessageNotifier           domain.PrivateMessageNotifier
	userRoomMessagesFinderByState    domain.UserRoomMessagesFinderByState
	roomMessagesUpdater              domain.RoomMessagesUpdater
	roomMessageNotifier              domain.RoomMessageNotifier
}

func NewUndeliveredMessagesNotificationRequestedUseCase(
	userPrivateMessagesFinderByState domain.UserPrivateMessagesFinderByState,
	privateMessagesUpdater domain.PrivateMessagesUpdater,
	privateMessageNotifier domain.PrivateMessageNotifier,
	userRoomMessagesFinderByState domain.UserRoomMessagesFinderByState,
	roomMessagesUpdater domain.RoomMessagesUpdater,
	roomMessageNotifier domain.RoomMessageNotifier,
) UndeliveredMessagesNotificationRequestedUseCase {
	return &undeliveredMessagesNotificationRequestedUseCase{
		userPrivateMessagesFinderByState: userPrivateMessagesFinderByState,
		privateMessagesUpdater:           privateMessagesUpdater,
		privateMessageNotifier:           privateMessageNotifier,
		userRoomMessagesFinderByState:    userRoomMessagesFinderByState,
		roomMessagesUpdater:              roomMessagesUpdater,
		roomMessageNotifier:              roomMessageNotifier,
	}
}

func (uc *undeliveredMessagesNotificationRequestedUseCase) Execute(ctx context.Context, input *UndeliveredMessagesNotificationRequestedInput) (*kernel.NoOutput, error) {
	privateMessages, err := uc.userPrivateMessagesFinderByState.FindsByState(ctx, input.UserID, domain.MessageStateUndelivered)
	if err != nil {
		return nil, err
	}

	roomMessages, err := uc.userRoomMessagesFinderByState.FindsByState(ctx, input.UserID, domain.MessageStateUndelivered)
	if err != nil {
		return nil, err
	}

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
		if err := uc.privateMessagesUpdater.Updates(ctx, privateMessages); err != nil {
			return nil, err
		}
	}

	if len(roomMessages) > 0 {
		if err := uc.roomMessagesUpdater.Updates(ctx, roomMessages); err != nil {
			return nil, err
		}
	}

	return nil, nil
}
