package application

import (
	"context"
	"gochat/internal/chat/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

var _ SendRoomMessageUseCase = (*sendRoomMessageUseCase)(nil)

type SendRoomMessageUseCase kernel.UseCase[*SendRoomMessageInput, *kernel.NoOutput]

type SendRoomMessageInput struct {
	SenderID   kernel.UserID
	RoomNumber kernel.RoomNumber
	Content    string
}

func (i *SendRoomMessageInput) Validate() error {
	if len(i.Content) == 0 {
		return myErrors.ErrEmptyInput
	}
	if len(i.SenderID) == 0 || len(i.RoomNumber) == 0 {
		return myErrors.ErrEmptyInput
	}
	return nil
}

type sendRoomMessageUseCase struct {
	messageIDGenerator    domain.MessageIDGenerator
	eventIDGenerator      event.IDGenerator
	roomFinder            domain.RoomFinder
	messageSaver          domain.RoomMessageSaver
	unpublishedEventSaver event.UnpublishedEventsSaver
	unitOfWork            kernel.UnitOfWork
}

func NewSendRoomMessageUseCase(
	messageIDGenerator domain.MessageIDGenerator,
	eventIDGenerator event.IDGenerator,
	roomFinder domain.RoomFinder,
	messageSaver domain.RoomMessageSaver,
	unpublishedEventSaver event.UnpublishedEventsSaver,
	unitOfWork kernel.UnitOfWork,
) SendRoomMessageUseCase {
	return &sendRoomMessageUseCase{
		messageIDGenerator:    messageIDGenerator,
		eventIDGenerator:      eventIDGenerator,
		roomFinder:            roomFinder,
		messageSaver:          messageSaver,
		unpublishedEventSaver: unpublishedEventSaver,
		unitOfWork:            unitOfWork,
	}
}

func (uc *sendRoomMessageUseCase) Execute(ctx context.Context, input *SendRoomMessageInput) (*kernel.NoOutput, error) {
	room, err := uc.roomFinder.FindByNumber(ctx, input.RoomNumber)
	if err != nil {
		return nil, err
	}

	message, err := domain.SendRoomMessage(
		room,
		input.SenderID,
		input.Content,
		uc.messageIDGenerator,
		uc.eventIDGenerator,
	)
	if err != nil {
		return nil, err
	}

	if err := uc.unitOfWork.Execute(ctx, func(txCtx context.Context) error {
		if err := uc.messageSaver.SaveRoomMessage(txCtx, message); err != nil {
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
