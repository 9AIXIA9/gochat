package usecase

import (
	"context"
	"gochat/internal/chat/application"
	"gochat/internal/chat/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
)

var _ UpdateMessageStateUseCase = (*updateMessageStateUseCase)(nil)

type UpdateMessageStateUseCase kernel.UseCase[*UpdateMessageStateInput, *kernel.NoOutput]

type UpdateMessageStateInput struct {
	NewState    domain.MessageState
	MessageID   domain.MessageID
	RecipientID kernel.UserID
}

func (i *UpdateMessageStateInput) Validate() error {
	if len(i.MessageID) == 0 || len(i.RecipientID) == 0 || len(i.NewState) == 0 {
		return myErrors.ErrEmptyInput
	}
	return nil
}

type updateMessageStateUseCase struct {
	messageStateUpdater application.MessageStateUpdater
}

func NewUpdateMessageStateUseCase(messageStateUpdater application.MessageStateUpdater) UpdateMessageStateUseCase {
	return &updateMessageStateUseCase{messageStateUpdater: messageStateUpdater}
}

func (uc *updateMessageStateUseCase) Execute(ctx context.Context, input *UpdateMessageStateInput) (*kernel.NoOutput, error) {
	return nil, uc.messageStateUpdater.Update(ctx, input.MessageID, input.RecipientID, input.NewState)
}
