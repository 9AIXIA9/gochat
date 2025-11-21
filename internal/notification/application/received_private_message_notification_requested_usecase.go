package application

import (
	"context"
	"gochat/internal/notification/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
)

//TODO 不太对

type ReceivedPrivateMessageNotificationRequestedUseCase kernel.UseCase[*ReceivedPrivateMessageNotificationRequestedInput, *kernel.NoOutput]

type ReceivedPrivateMessageNotificationRequestedInput struct {
	UserID kernel.UserID
}

func (r *ReceivedPrivateMessageNotificationRequestedInput) Validate() error {
	if len(r.UserID) == 0 {
		return myErrors.ErrEmptyInput
	}

	return nil
}

type receivedPrivateMessageNotificationRequestedUseCase struct {
	privateMessageFinder   domain.ReceivedPrivateMessageFinder
	privateMessageSaver    domain.PrivateMessageSaver
	privateMessageNotifier domain.PrivateMessageNotifier
}

func NewReceivedPrivateMessageNotificationRequestedUseCase(
	privateMessageFinder domain.ReceivedPrivateMessageFinder,
	privateMessageSaver domain.PrivateMessageSaver,
	privateMessageNotifier domain.PrivateMessageNotifier,
) ReceivedPrivateMessageNotificationRequestedUseCase {
	return &receivedPrivateMessageNotificationRequestedUseCase{
		privateMessageFinder:   privateMessageFinder,
		privateMessageSaver:    privateMessageSaver,
		privateMessageNotifier: privateMessageNotifier,
	}
}

func (uc *receivedPrivateMessageNotificationRequestedUseCase) Execute(ctx context.Context, input *ReceivedPrivateMessageNotificationRequestedInput) (*kernel.NoOutput, error) {
	privateMessages, err := uc.privateMessageFinder.FindReceivedPrivateMessage(ctx, input.UserID)
	if err != nil {
		return nil, err
	}

	for _, msg := range privateMessages {
		if err := msg.Deliver(uc.privateMessageNotifier); err != nil {
			return nil, err
		}
	}
	return nil, nil
}
