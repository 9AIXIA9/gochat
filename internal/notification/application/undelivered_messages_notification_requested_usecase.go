package application

import (
	"context"
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
	systemMessagesFinderByState  domain.UserSystemMessagesFinderByState
	systemMessagesUpdater        domain.SystemMessagesUpdater
	systemMessageNotifier        domain.SystemMessageNotifier
	privateMessagesFinderByState domain.UserPrivateMessagesFinderByState
	privateMessagesUpdater       domain.PrivateMessagesUpdater
	privateMessageNotifier       domain.PrivateMessageNotifier
	roomMessagesFinderByState    domain.UserRoomMessagesFinderByState
	roomMessagesUpdater          domain.RoomMessagesUpdater
	roomMessageNotifier          domain.RoomMessageNotifier
}

func NewUndeliveredMessagesNotificationRequestedUseCase(
	systemMessagesFinderByState domain.UserSystemMessagesFinderByState,
	systemMessagesUpdater domain.SystemMessagesUpdater,
	systemMessageNotifier domain.SystemMessageNotifier,
	privateMessagesFinderByState domain.UserPrivateMessagesFinderByState,
	privateMessagesUpdater domain.PrivateMessagesUpdater,
	privateMessageNotifier domain.PrivateMessageNotifier,
	roomMessagesFinderByState domain.UserRoomMessagesFinderByState,
	roomMessagesUpdater domain.RoomMessagesUpdater,
	roomMessageNotifier domain.RoomMessageNotifier,
) UndeliveredMessagesNotificationRequestedUseCase {
	return &undeliveredMessagesNotificationRequestedUseCase{
		systemMessagesFinderByState:  systemMessagesFinderByState,
		systemMessagesUpdater:        systemMessagesUpdater,
		systemMessageNotifier:        systemMessageNotifier,
		privateMessagesFinderByState: privateMessagesFinderByState,
		privateMessagesUpdater:       privateMessagesUpdater,
		privateMessageNotifier:       privateMessageNotifier,
		roomMessagesFinderByState:    roomMessagesFinderByState,
		roomMessagesUpdater:          roomMessagesUpdater,
		roomMessageNotifier:          roomMessageNotifier,
	}
}

func (uc *undeliveredMessagesNotificationRequestedUseCase) Execute(ctx context.Context, input *UndeliveredMessagesNotificationRequestedInput) (*kernel.NoOutput, error) {
	systemMessages, err := uc.systemMessagesFinderByState.FindsByState(ctx, input.UserID, domain.MessageStateUndelivered)
	if err != nil {
		return nil, err
	}

	privateMessages, err := uc.privateMessagesFinderByState.FindsByState(ctx, input.UserID, domain.MessageStateUndelivered)
	if err != nil {
		return nil, err
	}

	roomMessages, err := uc.roomMessagesFinderByState.FindsByState(ctx, input.UserID, domain.MessageStateUndelivered)
	if err != nil {
		return nil, err
	}

	for _, message := range systemMessages {
		if err := message.Deliver(uc.systemMessageNotifier); err != nil {
			return nil, err
		}
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

	if len(systemMessages) > 0 {
		if err := uc.systemMessagesUpdater.Updates(ctx, systemMessages); err != nil {
			return nil, err
		}
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
