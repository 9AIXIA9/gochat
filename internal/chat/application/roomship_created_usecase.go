package application

import (
	"context"
	"gochat/internal/chat/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
	"gochat/pkg/utils"
)

type RoomshipCreatedUseCase kernel.UseCase[*RoomshipCreatedInput, *kernel.NoOutput]

type RoomshipCreatedInput struct {
	ID     domain.RoomshipID
	UserID kernel.UserID
	RoomID kernel.RoomID
}

func (r *RoomshipCreatedInput) Validate() error {
	if len(r.ID) == 0 || len(r.UserID) == 0 || len(r.RoomID) == 0 {
		return myErrors.ErrEmptyInput
	}

	return nil
}

type roomshipCreatedUseCase struct {
	roomshipCreator domain.RoomshipCreator
}

func NewRoomshipCreatedUseCase(
	roomshipCreator domain.RoomshipCreator,
) (RoomshipCreatedUseCase, error) {
	if err := utils.CheckInterfaces(roomshipCreator); err != nil {
		return nil, err
	}

	return &roomshipCreatedUseCase{
		roomshipCreator: roomshipCreator,
	}, nil
}

func (uc *roomshipCreatedUseCase) Execute(ctx context.Context, input *RoomshipCreatedInput) (*kernel.NoOutput, error) {
	roomship := domain.CreateRoomship(input.ID, input.UserID, input.RoomID)

	return nil, uc.roomshipCreator.Create(ctx, roomship)
}
