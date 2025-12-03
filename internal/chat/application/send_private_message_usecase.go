package application

import (
	"context"
	"gochat/internal/chat/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

var _ SendPrivateMessageUseCase = (*sendPrivateMessageUseCase)(nil)

type SendPrivateMessageUseCase kernel.UseCase[*SendPrivateMessageInput, *kernel.NoOutput]

type SendPrivateMessageInput struct {
	SenderID    kernel.UserID
	RecipientID kernel.UserID
	Content     string
}

func (i *SendPrivateMessageInput) Validate() error {
	if len(i.Content) == 0 {
		return myErrors.ErrEmptyInput
	}
	if len(i.SenderID) == 0 || len(i.RecipientID) == 0 {
		return myErrors.ErrEmptyInput
	}
	return nil
}

type sendPrivateMessageUseCase struct {
	messageIDGenerator kernel.MessageIDGenerator
	eventIDGenerator   event.IDGenerator
	messageCreator     domain.PrivateMessageCreator
}

func NewSendPrivateMessageUseCase(
	messageIDGenerator kernel.MessageIDGenerator,
	eventIDGenerator event.IDGenerator,
	messageCreator domain.PrivateMessageCreator,
) SendPrivateMessageUseCase {
	return &sendPrivateMessageUseCase{
		messageIDGenerator: messageIDGenerator,
		eventIDGenerator:   eventIDGenerator,
		messageCreator:     messageCreator,
	}
}

func (uc *sendPrivateMessageUseCase) Execute(ctx context.Context, input *SendPrivateMessageInput) (*kernel.NoOutput, error) {
	// todo 校验用户是否存在

	message, err := domain.CreatePrivateMessage(
		input.RecipientID,
		input.SenderID,
		input.Content,
		uc.messageIDGenerator,
		uc.eventIDGenerator,
	)
	if err != nil {
		return nil, err
	}

	if err := uc.messageCreator.Create(ctx, message); err != nil {
		return nil, err
	}

	return nil, nil
}
