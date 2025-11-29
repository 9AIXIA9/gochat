package application

import (
	"context"
	"gochat/internal/notification/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
)

//TODO 修改 用户已读应该指的是一个房间的所有消息已读而不是单独的一个消息

type ReadRoomMessageUseCase kernel.UseCase[*ReadRoomMessageInput, *kernel.NoOutput]

type ReadRoomMessageInput struct {
	MessageID kernel.MessageID
	UserID    kernel.UserID
}

func (r *ReadRoomMessageInput) Validate() error {
	if len(r.UserID) == 0 || len(r.MessageID) == 0 {
		return myErrors.ErrEmptyInput
	}
	return nil
}

type readRoomMessageUseCase struct {
	messageFinder  domain.RoomMessageFinderByID
	messageUpdater domain.RoomMessageUpdater
}

func NewReadRoomMessageUseCase(
	messageFinder domain.RoomMessageFinderByID,
	messageUpdater domain.RoomMessageUpdater,
) ReadRoomMessageUseCase {
	return &readRoomMessageUseCase{
		messageFinder:  messageFinder,
		messageUpdater: messageUpdater,
	}
}

func (uc *readRoomMessageUseCase) Execute(ctx context.Context, input *ReadRoomMessageInput) (*kernel.NoOutput, error) {
	message, err := uc.messageFinder.FindRoomMessage(ctx, input.MessageID)
	if err != nil {
		return nil, err
	}

	message.Read(input.UserID)

	if err := uc.messageUpdater.Update(ctx, message); err != nil {
		return nil, err
	}

	return nil, nil
}
