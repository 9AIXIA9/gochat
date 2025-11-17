package usecase

import (
	"context"
	"gochat/internal/notification/application"
	"gochat/internal/notification/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
)

type MessageReadUseCase kernel.UseCase[*MessageReadInput, *kernel.NoOutput]

type MessageReadInput struct {
	MessageID kernel.MessageID
	UserID    kernel.UserID
}

func (r *MessageReadInput) Validate() error {
	if len(r.UserID) == 0 || len(r.MessageID) == 0 {
		return myErrors.ErrEmptyInput
	}

	return nil
}

type messageReadUseCase struct {
	messageStateUpdater application.MessageStateUpdater
}

func NewMessageReadUseCase(
	messageStateUpdater application.MessageStateUpdater,
) MessageReadUseCase {
	return &messageReadUseCase{
		messageStateUpdater: messageStateUpdater,
	}
}

func (uc *messageReadUseCase) Execute(ctx context.Context, input *MessageReadInput) (*kernel.NoOutput, error) {
	if err := uc.messageStateUpdater.UpdateMessageState(ctx, input.UserID, input.MessageID, domain.MessageStateRead); err != nil {
		return nil, err
	}

	return nil, nil
}
