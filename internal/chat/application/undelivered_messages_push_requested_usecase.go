package application

import (
	"context"
	"gochat/internal/chat/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
	"gochat/pkg/utils"
)

//TODO 消息太多的情况下还是没有优化，后续考虑分批次推送

type UndeliveredMessagesPushRequestedUseCase kernel.UseCase[*UndeliveredMessagesPushRequestedInput, *kernel.NoOutput]

type UndeliveredMessagesPushRequestedInput struct {
	UserID kernel.UserID
}

func (r *UndeliveredMessagesPushRequestedInput) Validate() error {
	if len(r.UserID) == 0 {
		return myErrors.ErrEmptyInput
	}

	return nil
}

type undeliveredMessagesPushRequestedUseCase struct {
	privateMessagesFinder  domain.PrivateMessagesFinderByRecipientIDAndState
	privateMessageNotifier domain.PrivateMessageNotifier
	privateMessagesUpdater domain.PrivateMessagesUpdater
	roomMessagesFinder     domain.RoomMessagesFinderByRecipientIDAndState
	roomMessageNotifier    domain.RoomMessageNotifier
	roomMessagesUpdater    domain.RoomMessagesUpdater
}

func NewUndeliveredMessagesPushRequestedUseCase(
	privateMessagesFinder domain.PrivateMessagesFinderByRecipientIDAndState,
	privateMessageNotifier domain.PrivateMessageNotifier,
	privateMessagesUpdater domain.PrivateMessagesUpdater,
	roomMessagesFinder domain.RoomMessagesFinderByRecipientIDAndState,
	roomMessageNotifier domain.RoomMessageNotifier,
	roomMessagesUpdater domain.RoomMessagesUpdater,
) (UndeliveredMessagesPushRequestedUseCase, error) {
	if err := utils.CheckInterfaces(
		privateMessagesFinder,
		privateMessageNotifier,
		privateMessagesUpdater,
		roomMessagesFinder,
		roomMessageNotifier,
		roomMessagesUpdater,
	); err != nil {
		return nil, err
	}

	return &undeliveredMessagesPushRequestedUseCase{
		privateMessagesFinder:  privateMessagesFinder,
		privateMessageNotifier: privateMessageNotifier,
		privateMessagesUpdater: privateMessagesUpdater,
		roomMessagesFinder:     roomMessagesFinder,
		roomMessageNotifier:    roomMessageNotifier,
		roomMessagesUpdater:    roomMessagesUpdater,
	}, nil
}

func (uc *undeliveredMessagesPushRequestedUseCase) Execute(ctx context.Context, input *UndeliveredMessagesPushRequestedInput) (*kernel.NoOutput, error) {
	privateMessages, err := uc.privateMessagesFinder.FindPrivateMessagesByRecipientIDAndState(ctx, input.UserID, domain.MessageStateUndelivered)
	if err != nil {
		return nil, err
	}

	roomMessages, err := uc.roomMessagesFinder.FindRoomMessagesByRecipientIDAndState(ctx, input.UserID, domain.MessageStateUndelivered)
	if err != nil {
		return nil, err
	}

	if len(privateMessages) > 0 {
		for _, message := range privateMessages {
			if err := message.Deliver(uc.privateMessageNotifier); err != nil {
				return nil, err
			}
		}

		if err := uc.privateMessagesUpdater.Updates(ctx, privateMessages); err != nil {
			return nil, err
		}
	}

	if len(roomMessages) > 0 {
		for _, message := range roomMessages {
			if err := message.Deliver(input.UserID, uc.roomMessageNotifier); err != nil {
				return nil, err
			}
		}

		if err := uc.roomMessagesUpdater.Updates(ctx, roomMessages); err != nil {
			return nil, err
		}
	}

	return nil, nil
}
