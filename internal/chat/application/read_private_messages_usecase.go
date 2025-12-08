package application

import (
	"context"
	"gochat/internal/chat/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
	"gochat/pkg/utils"
)

type ReadPrivateMessagesUseCase kernel.UseCase[*ReadPrivateMessagesInput, *kernel.NoOutput]

type ReadPrivateMessagesInput struct {
	SenderID    kernel.UserID
	RecipientID kernel.UserID
}

func (r *ReadPrivateMessagesInput) Validate() error {
	if len(r.SenderID) == 0 || len(r.RecipientID) == 0 {
		return myErrors.ErrEmptyInput
	}

	return nil
}

type readPrivateMessagesUseCase struct {
	updater domain.PrivateMessagesStatesUpdaterByUserID
}

func NewReadPrivateMessagesUseCase(
	updater domain.PrivateMessagesStatesUpdaterByUserID,
) (ReadPrivateMessagesUseCase, error) {
	if err := utils.CheckInterfaces(updater); err != nil {
		return nil, err
	}

	return &readPrivateMessagesUseCase{
		updater: updater,
	}, nil
}

func (uc *readPrivateMessagesUseCase) Execute(ctx context.Context, input *ReadPrivateMessagesInput) (*kernel.NoOutput, error) {
	return nil, uc.updater.UpdatesByUserID(ctx, input.SenderID, input.RecipientID, domain.MessageStateRead)
}
