package application

import (
	"context"
	"gochat/internal/profile/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
	"gochat/pkg/utils"
)

type RoomshipCreatedUseCase kernel.UseCase[*RoomshipCreatedInput, *kernel.NoOutput]

type RoomshipCreatedInput struct {
	ID     domain.RoomshipID
	UserID kernel.UserID
	RoomID kernel.RoomID
	Role   domain.RoomshipRole
}

func (r *RoomshipCreatedInput) Validate() error {
	if len(r.ID) == 0 || len(r.UserID) == 0 || len(r.RoomID) == 0 || len(r.Role) == 0 {
		return myErrors.WrapBusiness(myErrors.ErrEmptyInput, "roomship id, user id, room id or role is empty")
	}

	return nil
}

type roomshipCreatedUseCase struct {
	roomshipSaver domain.RoomshipSaver
}

func NewRoomshipCreatedUseCase(
	roomshipSaver domain.RoomshipSaver,
) (RoomshipCreatedUseCase, error) {
	if err := utils.CheckInterfaces(roomshipSaver); err != nil {
		return nil, err
	}

	return &roomshipCreatedUseCase{
		roomshipSaver: roomshipSaver,
	}, nil
}

func (uc *roomshipCreatedUseCase) Execute(ctx context.Context, input *RoomshipCreatedInput) (*kernel.NoOutput, error) {
	roomship := domain.LoadRoomship(input.ID, input.UserID, input.RoomID, input.Role)

	if err := uc.roomshipSaver.Save(ctx, roomship); err != nil {
		return nil, err
	}
	return nil, nil
}
