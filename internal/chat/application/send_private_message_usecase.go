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
	messageIDGenerator    domain.MessageIDGenerator
	eventIDGenerator      event.IDGenerator
	userFinderByNumber    domain.UserFinderByNumber
	messageSaver          domain.PrivateMessageSaver
	unpublishedEventSaver event.UnpublishedEventsSaver
	unitOfWork            kernel.UnitOfWork
}

func NewSendPrivateMessageUseCase(
	messageIDGenerator domain.MessageIDGenerator,
	eventIDGenerator event.IDGenerator,
	userFinderByNumber domain.UserFinderByNumber,
	messageSaver domain.PrivateMessageSaver,
	unpublishedEventSaver event.UnpublishedEventsSaver,
	unitOfWork kernel.UnitOfWork,
) SendPrivateMessageUseCase {
	return &sendPrivateMessageUseCase{
		messageIDGenerator:    messageIDGenerator,
		eventIDGenerator:      eventIDGenerator,
		userFinderByNumber:    userFinderByNumber,
		messageSaver:          messageSaver,
		unpublishedEventSaver: unpublishedEventSaver,
		unitOfWork:            unitOfWork,
	}
}

func (uc *sendPrivateMessageUseCase) Execute(ctx context.Context, input *SendPrivateMessageInput) (*kernel.NoOutput, error) {
	recipient, err := uc.userFinderByNumber.FindByNumber(ctx, input.RecipientNumber)
	if err != nil {
		return nil, err
	}

	message, err := domain.SendPrivateMessage(
		recipient,
		input.SenderID,
		input.Content,
		uc.messageIDGenerator,
		uc.eventIDGenerator,
	)
	if err != nil {
		return nil, err
	}

	if err := uc.unitOfWork.Execute(ctx, func(txCtx context.Context) error {
		if err := uc.messageSaver.SavePrivateMessage(txCtx, message); err != nil {
			return err
		}

		if err := uc.unpublishedEventSaver.Saves(txCtx, message.GetEvents()); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return nil, nil
}
