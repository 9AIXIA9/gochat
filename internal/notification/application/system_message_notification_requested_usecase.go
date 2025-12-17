package application

import (
	"context"
	"gochat/internal/notification/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
	"gochat/pkg/utils"
)

type SystemMessageNotificationRequestedUseCase kernel.UseCase[*SystemMessageNotificationRequestedInput, *kernel.NoOutput]

type SystemMessageNotificationRequestedInput struct {
	RecipientID kernel.UserID
	Content     string
}

func (r *SystemMessageNotificationRequestedInput) Validate() error {
	if len(r.RecipientID) == 0 {
		return myErrors.WrapBusiness(myErrors.ErrEmptyInput, "recipient id is empty")
	}

	if len(r.Content) == 0 {
		return myErrors.WrapBusiness(myErrors.ErrEmptyInput, "content is empty")
	}

	return nil
}

type systemMessageNotificationRequestedUseCase struct {
	messageCreator  domain.SystemMessageCreator
	messageNotifier domain.SystemMessageNotifier
	idGenerator     kernel.MessageIDGenerator
}

func NewSystemMessageNotificationRequestedUseCase(
	messageCreator domain.SystemMessageCreator,
	messageNotifier domain.SystemMessageNotifier,
	idGenerator kernel.MessageIDGenerator,
) (SystemMessageNotificationRequestedUseCase, error) {
	if err := utils.CheckInterfaces(
		messageCreator,
		messageNotifier,
		idGenerator,
	); err != nil {
		return nil, err
	}
	return &systemMessageNotificationRequestedUseCase{
		messageCreator:  messageCreator,
		messageNotifier: messageNotifier,
		idGenerator:     idGenerator,
	}, nil
}

func (uc *systemMessageNotificationRequestedUseCase) Execute(ctx context.Context, input *SystemMessageNotificationRequestedInput) (*kernel.NoOutput, error) {
	message, err := domain.CreateSystemMessage(
		input.RecipientID,
		input.Content,
		uc.idGenerator,
		uc.messageNotifier,
	)
	if err != nil {
		return nil, err
	}

	if err := uc.messageCreator.Create(ctx, message); err != nil {
		return nil, err
	}
	return nil, nil
}
