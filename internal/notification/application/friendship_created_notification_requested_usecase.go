package application

import (
	"context"
	"encoding/json"
	"gochat/internal/notification/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
	"gochat/pkg/utils"
	"time"
)

type FriendshipCreatedNotificationRequestedUseCase kernel.UseCase[*FriendshipCreatedNotificationRequestedInput, *kernel.NoOutput]

type FriendshipCreatedNotificationRequestedInput struct {
	UserID    kernel.UserID
	FriendID  kernel.UserID
	CreatedAt time.Time
}

func (r *FriendshipCreatedNotificationRequestedInput) Validate() error {
	if len(r.UserID) == 0 || len(r.FriendID) == 0 {
		return myErrors.ErrEmptyInput
	}

	return nil
}

type friendshipCreatedNotificationRequestedUseCase struct {
	idGenerator     kernel.MessageIDGenerator
	messageCreator  domain.SystemMessageCreator
	messageNotifier domain.SystemMessageNotifier
}

func NewFriendshipCreatedNotificationRequestedUseCase(
	idGenerator kernel.MessageIDGenerator,
	messageCreator domain.SystemMessageCreator,
	messageNotifier domain.SystemMessageNotifier,
) (FriendshipCreatedNotificationRequestedUseCase, error) {
	if err := utils.CheckInterfaces(idGenerator, messageCreator, messageNotifier); err != nil {
		return nil, err
	}
	return &friendshipCreatedNotificationRequestedUseCase{
		idGenerator:     idGenerator,
		messageCreator:  messageCreator,
		messageNotifier: messageNotifier,
	}, nil
}

func (uc *friendshipCreatedNotificationRequestedUseCase) Execute(ctx context.Context, input *FriendshipCreatedNotificationRequestedInput) (*kernel.NoOutput, error) {
	content, err := uc.buildSystemMessageContent(input.FriendID, input.CreatedAt)
	if err != nil {
		return nil, err
	}

	message, err := domain.CreateSystemMessage(domain.TopicFriendshipCreated, input.UserID, content, uc.idGenerator, uc.messageNotifier)
	if err != nil {
		return nil, err
	}

	if err := uc.messageCreator.Create(ctx, message); err != nil {
		return nil, err
	}

	return nil, nil
}

func (uc *friendshipCreatedNotificationRequestedUseCase) buildSystemMessageContent(
	friendID kernel.UserID,
	createdAt time.Time,
) ([]byte, error) {
	type Alias struct {
		FriendID  kernel.UserID
		CreatedAt time.Time
	}
	return json.Marshal(Alias{
		FriendID:  friendID,
		CreatedAt: createdAt,
	})
}
