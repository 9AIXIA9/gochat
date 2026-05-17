package application

import (
	"context"
	"gochat/internal/chat/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
	"gochat/pkg/validate"
)

const messageCountLimit = 100

type UndeliveredMessagesPushRequestedUseCase kernel.UseCase[*UndeliveredMessagesPushRequestedInput, *kernel.NoOutput]

type UndeliveredMessagesPushRequestedInput struct {
	UserID kernel.UserID
}

func (r *UndeliveredMessagesPushRequestedInput) Validate() error {
	if len(r.UserID) == 0 {
		return myErrors.WrapBusiness(myErrors.ErrEmptyInput, "user id can't be empty")
	}

	return nil
}

type undeliveredMessagesPushRequestedUseCase struct {
	privateMessagesFinder  domain.PrivateMessagesFinderByRecipientIDAndState
	privateMessageNotifier domain.PrivateMessageNotifier
	roomMessagesFinder     domain.RoomMessagesFinderByRecipientIDAndState
	roomMessageNotifier    domain.RoomMessageNotifier
}

func NewUndeliveredMessagesPushRequestedUseCase(
	privateMessagesFinder domain.PrivateMessagesFinderByRecipientIDAndState,
	privateMessageNotifier domain.PrivateMessageNotifier,
	roomMessagesFinder domain.RoomMessagesFinderByRecipientIDAndState,
	roomMessageNotifier domain.RoomMessageNotifier,
) (UndeliveredMessagesPushRequestedUseCase, error) {
	if err := validate.NotNil(
		privateMessagesFinder,
		privateMessageNotifier,
		roomMessagesFinder,
		roomMessageNotifier,
	); err != nil {
		return nil, err
	}

	return &undeliveredMessagesPushRequestedUseCase{
		privateMessagesFinder:  privateMessagesFinder,
		privateMessageNotifier: privateMessageNotifier,
		roomMessagesFinder:     roomMessagesFinder,
		roomMessageNotifier:    roomMessageNotifier,
	}, nil
}

func (uc *undeliveredMessagesPushRequestedUseCase) Execute(ctx context.Context, input *UndeliveredMessagesPushRequestedInput) (*kernel.NoOutput, error) {
	privateMessages, err := uc.privateMessagesFinder.FindPrivateMessagesByRecipientIDAndState(ctx, input.UserID, domain.MessageStateUndelivered, messageCountLimit)
	if err != nil {
		return nil, err
	}

	if len(privateMessages) > 0 {
		for _, message := range privateMessages {
			_ = uc.privateMessageNotifier.Notify(message)
		}
	}

	roomMessages, err := uc.roomMessagesFinder.FindRoomMessagesByRecipientIDAndState(ctx, input.UserID, domain.MessageStateUndelivered, messageCountLimit)
	if err != nil {
		return nil, err
	}

	if len(roomMessages) > 0 {
		for _, message := range roomMessages {
			_, _ = uc.roomMessageNotifier.Notify(message, message.UndeliveredRecipientIDs())

		}
	}

	return nil, nil
}
