package application

import (
	"context"
	"gochat/internal/notification/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
)

type RoomMessageReadUseCase kernel.UseCase[*RoomMessageReadInput, *kernel.NoOutput]

type RoomMessageReadInput struct {
	MessageID kernel.MessageID
	UserID    kernel.UserID
}

func (r *RoomMessageReadInput) Validate() error {
	if len(r.UserID) == 0 || len(r.MessageID) == 0 {
		return myErrors.ErrEmptyInput
	}
	return nil
}

type roomMessageReadUseCase struct {
	messageFinder domain.RoomMessageFinder
	messageSaver  domain.RoomMessageSaver
}

func NewRoomMessageReadUseCase(
	messageFinder domain.RoomMessageFinder,
	messageSaver domain.RoomMessageSaver,
) RoomMessageReadUseCase {
	return &roomMessageReadUseCase{
		messageFinder: messageFinder,
		messageSaver:  messageSaver,
	}
}

func (uc *roomMessageReadUseCase) Execute(ctx context.Context, input *RoomMessageReadInput) (*kernel.NoOutput, error) {
	message, err := uc.messageFinder.FindRoomMessage(ctx, input.MessageID)
	if err != nil {
		return nil, err
	}

	message.Read(input.UserID)

	if err := uc.messageSaver.SaveRoomMessage(ctx, message); err != nil {
		return nil, err
	}

	return nil, nil
}
