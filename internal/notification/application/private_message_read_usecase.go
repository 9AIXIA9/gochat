package application

import (
	"context"
	"gochat/internal/notification/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
)

type PrivateMessageReadUseCase kernel.UseCase[*PrivateMessageReadInput, *kernel.NoOutput]

type PrivateMessageReadInput struct {
	MessageID kernel.MessageID
}

func (r *PrivateMessageReadInput) Validate() error {
	if len(r.MessageID) == 0 {
		return myErrors.ErrEmptyInput
	}
	return nil
}

type privateMessageReadUseCase struct {
	messageFinder domain.PrivateMessageFinder
	messageSaver  domain.PrivateMessageSaver
}

func NewPrivateMessageReadUseCase(
	messageFinder domain.PrivateMessageFinder,
	messageSaver domain.PrivateMessageSaver,
) PrivateMessageReadUseCase {
	return &privateMessageReadUseCase{
		messageFinder: messageFinder,
		messageSaver:  messageSaver,
	}
}

func (uc *privateMessageReadUseCase) Execute(ctx context.Context, input *PrivateMessageReadInput) (*kernel.NoOutput, error) {
	message, err := uc.messageFinder.FindPrivateMessage(ctx, input.MessageID)
	if err != nil {
		return nil, err
	}

	message.Read()

	if err := uc.messageSaver.SavePrivateMessage(ctx, message); err != nil {
		return nil, err
	}

	return nil, nil
}
