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

type FriendRequestCreatedNotificationRequestedUseCase kernel.UseCase[*FriendRequestCreatedNotificationRequestedInput, *kernel.NoOutput]

type FriendRequestCreatedNotificationRequestedInput struct {
	RequestID kernel.OperationID
	From      kernel.UserID
	To        kernel.UserID
	SentAt    time.Time
	Content   string
}

func (r *FriendRequestCreatedNotificationRequestedInput) Validate() error {
	if len(r.RequestID) == 0 || len(r.From) == 0 || len(r.To) == 0 {
		return myErrors.ErrEmptyInput
	}

	return nil
}

type friendRequestCreatedNotificationRequestedUseCase struct {
	idGenerator     kernel.MessageIDGenerator
	messageCreator  domain.SystemMessageCreator
	messageNotifier domain.SystemMessageNotifier
}

func NewFriendRequestCreatedNotificationRequestedUseCase(
	idGenerator kernel.MessageIDGenerator,
	messageCreator domain.SystemMessageCreator,
	messageNotifier domain.SystemMessageNotifier,
) (FriendRequestCreatedNotificationRequestedUseCase, error) {
	if err := utils.CheckInterfaces(idGenerator, messageCreator, messageNotifier); err != nil {
		return nil, err
	}
	return &friendRequestCreatedNotificationRequestedUseCase{
		idGenerator:     idGenerator,
		messageCreator:  messageCreator,
		messageNotifier: messageNotifier,
	}, nil
}

func (uc *friendRequestCreatedNotificationRequestedUseCase) Execute(ctx context.Context, input *FriendRequestCreatedNotificationRequestedInput) (*kernel.NoOutput, error) {
	content, err := uc.buildSystemMessageContent(input.RequestID, input.From, input.SentAt, input.Content)
	if err != nil {
		return nil, err
	}

	message, err := domain.CreateSystemMessage(domain.TopicFriendRequestCreated, input.To, content, uc.idGenerator, uc.messageNotifier)
	if err != nil {
		return nil, err
	}

	if err := uc.messageCreator.Create(ctx, message); err != nil {
		return nil, err
	}

	return nil, nil
}

func (uc *friendRequestCreatedNotificationRequestedUseCase) buildSystemMessageContent(
	requestID kernel.OperationID,
	from kernel.UserID,
	sentAt time.Time,
	content string,
) ([]byte, error) {
	type Alias struct {
		RequestID kernel.OperationID
		From      kernel.UserID
		SentAt    time.Time
		Content   string
	}
	return json.Marshal(Alias{
		RequestID: requestID,
		From:      from,
		SentAt:    sentAt,
		Content:   content,
	})
}
