package application

import (
	"context"
	"gochat/internal/notification/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
)

//TODO 修改 用户已读应该指的是一个用户的所有消息已读而不是单独的一个消息

type ReadPrivateMessageUseCase kernel.UseCase[*ReadPrivateMessageInput, *kernel.NoOutput]

type ReadPrivateMessageInput struct {
	MessageID kernel.MessageID
}

func (r *ReadPrivateMessageInput) Validate() error {
	if len(r.MessageID) == 0 {
		return myErrors.ErrEmptyInput
	}
	return nil
}

type readPrivateMessageUseCase struct {
	messageFinder  domain.PrivateMessageFinderByID
	messageUpdater domain.PrivateMessageUpdater
}

func NewReadPrivateMessageUseCase(
	messageFinder domain.PrivateMessageFinderByID,
	messageUpdater domain.PrivateMessageUpdater,
) ReadPrivateMessageUseCase {
	return &readPrivateMessageUseCase{
		messageFinder:  messageFinder,
		messageUpdater: messageUpdater,
	}
}

func (uc *readPrivateMessageUseCase) Execute(ctx context.Context, input *ReadPrivateMessageInput) (*kernel.NoOutput, error) {
	message, err := uc.messageFinder.FindByID(ctx, input.MessageID)
	if err != nil {
		return nil, err
	}

	message.Read()

	if err := uc.messageUpdater.Update(ctx, message); err != nil {
		return nil, err
	}

	return nil, nil
}
