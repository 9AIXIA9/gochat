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
	systemMessagesFinderByState domain.UserSystemMessagesFinderByState
	systemMessagesUpdater       domain.SystemMessagesUpdater
	systemMessageNotifier       domain.SystemMessageNotifier
}

func NewUndeliveredMessagesNotificationRequestedUseCase(
	systemMessagesFinderByState domain.UserSystemMessagesFinderByState,
	systemMessagesUpdater domain.SystemMessagesUpdater,
	systemMessageNotifier domain.SystemMessageNotifier,
) UndeliveredMessagesNotificationRequestedUseCase {
	return &undeliveredMessagesNotificationRequestedUseCase{
		systemMessagesFinderByState: systemMessagesFinderByState,
		systemMessagesUpdater:       systemMessagesUpdater,
		systemMessageNotifier:       systemMessageNotifier,
	}
}

func (uc *undeliveredMessagesNotificationRequestedUseCase) Execute(ctx context.Context, input *UndeliveredMessagesNotificationRequestedInput) (*kernel.NoOutput, error) {
	systemMessages, err := uc.systemMessagesFinderByState.FindsByState(ctx, input.UserID, domain.MessageStateUndelivered)
	if err != nil {
		return nil, err
	}

	for _, message := range systemMessages {
		if err := message.Deliver(uc.systemMessageNotifier); err != nil {
			return nil, err
		}
	}

	if len(systemMessages) > 0 {
		if err := uc.systemMessagesUpdater.Updates(ctx, systemMessages); err != nil {
			return nil, err
		}
	}

	return nil, nil
}
