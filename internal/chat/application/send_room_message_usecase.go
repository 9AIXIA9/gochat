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
	messageIDGenerator domain.MessageIDGenerator
	eventIDGenerator   event.IDGenerator
	roomFinder         domain.RoomFinderByNumber
	messageCreator     domain.RoomMessageCreator
}

func NewSendRoomMessageUseCase(
	messageIDGenerator domain.MessageIDGenerator,
	eventIDGenerator event.IDGenerator,
	roomFinder domain.RoomFinderByNumber,
	messageCreator domain.RoomMessageCreator,
) SendRoomMessageUseCase {
	return &sendRoomMessageUseCase{
		messageIDGenerator: messageIDGenerator,
		eventIDGenerator:   eventIDGenerator,
		roomFinder:         roomFinder,
		messageCreator:     messageCreator,
	}
}

func (uc *sendRoomMessageUseCase) Execute(ctx context.Context, input *SendRoomMessageInput) (*kernel.NoOutput, error) {
	room, err := uc.roomFinder.FindByNumber(ctx, input.RoomNumber)
	if err != nil {
		return nil, err
	}

	message, err := domain.CreateRoomMessage(
		room,
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
