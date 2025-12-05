package application

import (
	"context"
	"gochat/internal/chat/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"gochat/pkg/utils"
)

var _ SendRoomMessageUseCase = (*sendRoomMessageUseCase)(nil)

type SendRoomMessageUseCase kernel.UseCase[*SendRoomMessageInput, *kernel.NoOutput]

type SendRoomMessageInput struct {
	SenderID kernel.UserID
	RoomID   kernel.RoomID
	Content  string
}

func (i *SendRoomMessageInput) Validate() error {
	if len(i.Content) == 0 {
		return myErrors.ErrEmptyInput
	}
	if len(i.SenderID) == 0 || len(i.RoomID) == 0 {
		return myErrors.ErrEmptyInput
	}
	return nil
}

type sendRoomMessageUseCase struct {
	messageIDGenerator kernel.MessageIDGenerator
	eventIDGenerator   event.IDGenerator
	roomshipFinder     domain.RoomshipFinderByUserIDAndRoomID
	messageCreator     domain.RoomMessageCreator
}

func NewSendRoomMessageUseCase(
	messageIDGenerator kernel.MessageIDGenerator,
	eventIDGenerator event.IDGenerator,
	roomshipFinder domain.RoomshipFinderByUserIDAndRoomID,
	messageCreator domain.RoomMessageCreator,
) (SendRoomMessageUseCase, error) {
	if err := utils.CheckInterfaces(
		messageIDGenerator,
		eventIDGenerator,
		roomshipFinder,
		messageCreator,
	); err != nil {
		return nil, err
	}

	return &sendRoomMessageUseCase{
		messageIDGenerator: messageIDGenerator,
		eventIDGenerator:   eventIDGenerator,
		roomshipFinder:     roomshipFinder,
		messageCreator:     messageCreator,
	}, nil
}

func (uc *sendRoomMessageUseCase) Execute(ctx context.Context, input *SendRoomMessageInput) (*kernel.NoOutput, error) {
	roomship, err := uc.roomshipFinder.FindByUserIDAndRoomID(ctx, input.RoomID, input.SenderID)
	if err != nil {
		return nil, err
	}

	message, err := domain.CreateRoomMessage(
		roomship.RoomID(),
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
