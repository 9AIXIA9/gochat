package usecase

import (
	"context"
	"errors"
	"gochat/internal/notification/application"
	"gochat/internal/notification/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
)

type UndeliveredMessageNotificationRequestedUseCase kernel.UseCase[*UndeliveredMessageNotificationRequestedInput, *kernel.NoOutput]

type UndeliveredMessageNotificationRequestedInput struct {
	UserID kernel.UserID
}

func (r *UndeliveredMessageNotificationRequestedInput) Validate() error {
	if len(r.UserID) == 0 {
		return myErrors.ErrEmptyInput
	}

	return nil
}

type undeliveredMessageNotificationRequestedUseCase struct {
	messageFinder        application.MessageFinder
	messageNotifier      application.MessageNotifier
	messageStatesUpdater application.MessageStatesUpdater
}

func NewUndeliveredMessageNotificationRequestedUseCase(
	messageFinder application.MessageFinder,
	messageNotifier application.MessageNotifier,
	messageStatesUpdater application.MessageStatesUpdater,
) UndeliveredMessageNotificationRequestedUseCase {
	return &undeliveredMessageNotificationRequestedUseCase{
		messageFinder:        messageFinder,
		messageNotifier:      messageNotifier,
		messageStatesUpdater: messageStatesUpdater,
	}
}

func (uc *undeliveredMessageNotificationRequestedUseCase) Execute(ctx context.Context, input *UndeliveredMessageNotificationRequestedInput) (*kernel.NoOutput, error) {
	messages, err := uc.messageFinder.FindMessagesByUserID(ctx, input.UserID)
	if err != nil {
		return nil, err
	}

	deliveredIDs := make([]kernel.MessageID, 0, len(messages))
	for _, message := range messages {
		if err := uc.messageNotifier.Notify(input.UserID, message); err != nil {
			if errors.Is(err, myErrors.ErrNotFound) {
				return nil, nil
			}
			return nil, err
		}
		deliveredIDs = append(deliveredIDs, message.ID())
	}
	if err := uc.messageStatesUpdater.UpdateMessageStates(ctx, input.UserID, deliveredIDs, domain.MessageStateDelivered); err != nil {
		return nil, err
	}

	return nil, nil
}
