package application

import (
	"context"
	"gochat/internal/notification/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
	"gochat/pkg/validate"
)

type ConfirmSystemMessagesUseCase kernel.UseCase[*ConfirmSystemMessagesInput, *kernel.NoOutput]

type ConfirmSystemMessagesInput struct {
	UserID     kernel.UserID
	MessageIDs []kernel.MessageID
}

func (r *ConfirmSystemMessagesInput) Validate() error {
	if len(r.UserID) == 0 || len(r.MessageIDs) == 0 {
		return myErrors.WrapBusiness(myErrors.ErrEmptyInput, "user id and message ids can't be empty")
	}
	return nil
}

type confirmSystemMessagesUseCase struct {
	updater domain.SystemMessagesStatesUpdaterByMessageIDs
}

func NewConfirmSystemMessagesUseCase(
	updater domain.SystemMessagesStatesUpdaterByMessageIDs,
) (ConfirmSystemMessagesUseCase, error) {
	if err := validate.NotNil(updater); err != nil {
		return nil, err
	}
	return &confirmSystemMessagesUseCase{updater: updater}, nil
}

func (uc *confirmSystemMessagesUseCase) Execute(ctx context.Context, input *ConfirmSystemMessagesInput) (*kernel.NoOutput, error) {
	return nil, uc.updater.UpdatesByMessageIDs(ctx, input.UserID, input.MessageIDs, domain.MessageStateDelivered)
}
