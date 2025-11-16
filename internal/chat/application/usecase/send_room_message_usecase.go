package usecase

import (
	"context"
	"gochat/internal/chat/application"
	"gochat/internal/chat/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

var _ SendRoomMessageUseCase = (*sendRoomMessageUseCase)(nil)

type SendRoomMessageUseCase kernel.UseCase[*SendRoomMessageInput, *kernel.NoOutput]

type SendRoomMessageInput struct {
	SenderID   kernel.UserID
	RoomNumber domain.RoomNumber
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
	messageIDGenerator    application.MessageIDGenerator
	eventIDGenerator      event.IDGenerator
	roomFinder            application.RoomFinder
	messageSaver          application.RoomMessageSaver
	unpublishedEventSaver event.UnpublishedSaver
	unitOfWork            kernel.UnitOfWork
}

func NewSendRoomMessageUseCase(
	messageIDGenerator application.MessageIDGenerator,
	eventIDGenerator event.IDGenerator,
	roomFinder application.RoomFinder,
	messageSaver application.RoomMessageSaver,
	unpublishedEventSaver event.UnpublishedSaver,
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

	message, err := room.ReceiveMessage(uc.messageIDGenerator.Generate(), input.SenderID, input.Content, uc.eventIDGenerator)
	if err != nil {
		return nil, err
	}

	if err := uc.unitOfWork.Execute(ctx, func(txCtx context.Context) error {
		if err := uc.messageSaver.SaveRoomMessage(txCtx, room.ID(), message); err != nil {
			return err
		}

		if err := uc.unpublishedEventSaver.Saves(txCtx, room.GetEvents()); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return nil, nil
}
