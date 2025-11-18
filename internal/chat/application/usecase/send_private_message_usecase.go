package usecase

import (
	"context"
	"gochat/internal/chat/application"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

var _ SendPrivateMessageUseCase = (*sendPrivateMessageUseCase)(nil)

type SendPrivateMessageUseCase kernel.UseCase[*SendPrivateMessageInput, *kernel.NoOutput]

type SendPrivateMessageInput struct {
	SenderID        kernel.UserID
	RecipientNumber kernel.UserNumber
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
	unpublishedEventSaver event.UnpublishedEventsSaver
	unitOfWork            kernel.UnitOfWork
}

func NewSendPrivateMessageUseCase(
	messageIDGenerator application.MessageIDGenerator,
	eventIDGenerator event.IDGenerator,
	userFinder application.UserFinder,
	messageSaver application.PrivateMessageSaver,
	unpublishedEventSaver event.UnpublishedEventsSaver,
	unitOfWork kernel.UnitOfWork,
) SendPrivateMessageUseCase {
	return &sendPrivateMessageUseCase{
		messageIDGenerator:    messageIDGenerator,
		eventIDGenerator:      eventIDGenerator,
		userFinder:            userFinder,
		messageSaver:          messageSaver,
		unpublishedEventSaver: unpublishedEventSaver,
		unitOfWork:            unitOfWork,
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

	if err := uc.unitOfWork.Execute(ctx, func(txCtx context.Context) error {
		if err := uc.messageSaver.SavePrivateMessage(txCtx, user.ID(), message); err != nil {
			return err
		}

		events := user.GetEvents()
		if err := uc.unpublishedEventSaver.Saves(txCtx, events); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return nil, nil
}
