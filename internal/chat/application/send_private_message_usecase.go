package application

import (
	"context"
	"gochat/internal/chat/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
	"gochat/pkg/utils"
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
	exister            domain.FriendshipExisterByUserID
	messageIDGenerator kernel.MessageIDGenerator
	notifier           domain.PrivateMessageNotifier
	messageCreator     domain.PrivateMessageCreator
}

func NewSendPrivateMessageUseCase(
	exister domain.FriendshipExisterByUserID,
	messageIDGenerator kernel.MessageIDGenerator,
	notifier domain.PrivateMessageNotifier,
	messageCreator domain.PrivateMessageCreator,
) (SendPrivateMessageUseCase, error) {
	if err := utils.CheckInterfaces(
		exister,
		messageIDGenerator,
		notifier,
		messageCreator,
	); err != nil {
		return nil, err
	}
	return &sendPrivateMessageUseCase{
		exister:            exister,
		messageIDGenerator: messageIDGenerator,
		notifier:           notifier,
		messageCreator:     messageCreator,
	}, nil
}

func (uc *sendPrivateMessageUseCase) Execute(ctx context.Context, input *SendPrivateMessageInput) (*kernel.NoOutput, error) {
	if input.SenderID != input.RecipientID {
		exist, err := uc.exister.ExistByUserID(ctx, input.SenderID, input.RecipientID)
		if err != nil {
			return nil, err
		}

		if !exist {
			return nil, domain.ErrNotFriends
		}
	}

	message, err := domain.CreatePrivateMessage(
		input.RecipientID,
		input.SenderID,
		input.Content,
		uc.messageIDGenerator,
		uc.notifier,
	)
	if err != nil {
		return nil, err
	}

	if err := uc.messageCreator.Create(ctx, message); err != nil {
		return nil, err
	}

	return nil, nil
}
