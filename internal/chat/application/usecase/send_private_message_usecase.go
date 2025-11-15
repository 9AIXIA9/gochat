package usecase

import (
	"context"
	"gochat/internal/chat/application"
	"gochat/internal/chat/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

//TODO 一致性保証（事件和用户创建）

var _ SendPrivateMessageUseCase = (*sendPrivateMessageUseCase)(nil)

type SendPrivateMessageUseCase kernel.UseCase[*SendPrivateMessageInput, *kernel.NoOutput]

type SendPrivateMessageInput struct {
	SenderID        kernel.UserID
	RecipientNumber domain.UserNumber
	Content         string
}

func (i *SendPrivateMessageInput) Validate() error {
	if len(i.Content) == 0 {
		return myErrors.ErrEmptyInput
	}
	if len(i.SenderID) == 0 || len(i.RecipientNumber) == 0 {
		return myErrors.ErrEmptyInput
	}
	return nil
}

type sendPrivateMessageUseCase struct {
	messageIDGenerator    application.MessageIDGenerator
	eventIDGenerator      event.IDGenerator
	userFinder            application.UserFinder
	messageSaver          application.PrivateMessageSaver
	unpublishedEventSaver event.UnpublishedSaver
}

func NewSendPrivateMessageUseCase(
	messageIDGenerator application.MessageIDGenerator,
	eventIDGenerator event.IDGenerator,
	userFinder application.UserFinder,
	messageSaver application.PrivateMessageSaver,
	unpublishedEventSaver event.UnpublishedSaver,
) SendPrivateMessageUseCase {
	return &sendPrivateMessageUseCase{
		messageIDGenerator:    messageIDGenerator,
		eventIDGenerator:      eventIDGenerator,
		userFinder:            userFinder,
		messageSaver:          messageSaver,
		unpublishedEventSaver: unpublishedEventSaver,
	}
}

func (uc *sendPrivateMessageUseCase) Execute(ctx context.Context, input *SendPrivateMessageInput) (*kernel.NoOutput, error) {
	user, err := uc.userFinder.FindByNumber(ctx, input.RecipientNumber)
	if err != nil {
		return nil, err
	}

	message, err := user.ReceiveMessage(uc.messageIDGenerator.Generate(), input.SenderID, input.Content, uc.eventIDGenerator)
	if err != nil {
		return nil, err
	}

	if err := uc.messageSaver.SavePrivateMessage(ctx, user.ID(), message); err != nil {
		return nil, err
	}

	events := user.GetEvents()
	if err := uc.unpublishedEventSaver.Saves(ctx, events); err != nil {
		return nil, err
	}
	return nil, nil
}
