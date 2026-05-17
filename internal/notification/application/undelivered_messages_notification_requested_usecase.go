package application

import (
	"context"
	"gochat/internal/notification/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
	"gochat/pkg/validate"
)

const messageCountLimit = 100

type UndeliveredMessagesNotificationRequestedUseCase kernel.UseCase[*UndeliveredMessagesNotificationRequestedInput, *kernel.NoOutput]

type UndeliveredMessagesNotificationRequestedInput struct {
	UserID kernel.UserID
}

func (r *UndeliveredMessagesNotificationRequestedInput) Validate() error {
	if len(r.UserID) == 0 {
		return myErrors.WrapBusiness(myErrors.ErrEmptyInput, "user id is empty")
	}
	return nil
}

type undeliveredMessagesNotificationRequestedUseCase struct {
	systemMessagesFinderByState domain.UserSystemMessagesFinderByState
	systemMessageNotifier       domain.SystemMessageNotifier
}

func NewUndeliveredMessagesNotificationRequestedUseCase(
	systemMessagesFinderByState domain.UserSystemMessagesFinderByState,
	systemMessageNotifier domain.SystemMessageNotifier,
) (UndeliveredMessagesNotificationRequestedUseCase, error) {
	if err := validate.NotNil(
		systemMessagesFinderByState,
		systemMessageNotifier,
	); err != nil {
		return nil, err
	}
	return &undeliveredMessagesNotificationRequestedUseCase{
		systemMessagesFinderByState: systemMessagesFinderByState,
		systemMessageNotifier:       systemMessageNotifier,
	}, nil
}

func (uc *undeliveredMessagesNotificationRequestedUseCase) Execute(ctx context.Context, input *UndeliveredMessagesNotificationRequestedInput) (*kernel.NoOutput, error) {
	systemMessages, err := uc.systemMessagesFinderByState.FindsByState(ctx, input.UserID, domain.MessageStateUndelivered, messageCountLimit)
	if err != nil {
		return nil, err
	}

	if len(systemMessages) == 0 {
		return nil, nil
	}

	for _, message := range systemMessages {
		_ = uc.systemMessageNotifier.Notify(message)
	}
	return nil, nil
}
