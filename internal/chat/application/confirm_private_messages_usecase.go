package application

import (
	"context"
	"gochat/internal/chat/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
	"gochat/pkg/validate"
)

type ConfirmPrivateMessagesUseCase kernel.UseCase[*ConfirmPrivateMessagesInput, *kernel.NoOutput]

type ConfirmPrivateMessagesInput struct {
	UserID     kernel.UserID
	MessageIDs []kernel.MessageID
}

func (r *ConfirmPrivateMessagesInput) Validate() error {
	if len(r.UserID) == 0 || len(r.MessageIDs) == 0 {
		return myErrors.WrapBusiness(myErrors.ErrEmptyInput, "user id and message ids can't be empty")
	}

	return nil
}

type confirmPrivateMessagesUseCase struct {
	updater domain.PrivateMessagesStatesUpdaterByMessageIDs
}

func NewConfirmPrivateMessagesUseCase(
	updater domain.PrivateMessagesStatesUpdaterByMessageIDs,
) (ConfirmPrivateMessagesUseCase, error) {
	if err := validate.NotNil(updater); err != nil {
		return nil, err
	}

	return &confirmPrivateMessagesUseCase{
		updater: updater,
	}, nil
}

func (uc *confirmPrivateMessagesUseCase) Execute(ctx context.Context, input *ConfirmPrivateMessagesInput) (*kernel.NoOutput, error) {
	return nil, uc.updater.UpdatesByMessageIDs(ctx, input.UserID, input.MessageIDs, domain.MessageStateDelivered)
}
