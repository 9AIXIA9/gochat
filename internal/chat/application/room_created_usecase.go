package application

import (
	"context"
	"gochat/internal/chat/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
	"gochat/pkg/validate"
)

type RoomCreatedUseCase kernel.UseCase[*RoomCreatedInput, *kernel.NoOutput]

type RoomCreatedInput struct {
	RoomID kernel.RoomID
}

func (r *RoomCreatedInput) Validate() error {
	if len(r.RoomID) == 0 {
		return myErrors.WrapBusiness(myErrors.ErrEmptyInput, "room id can't be empty")
	}

	return nil
}

type roomCreatedUseCase struct {
	roomSaver domain.RoomSaver
}

func NewRoomCreatedUseCase(
	roomSaver domain.RoomSaver,
) (RoomCreatedUseCase, error) {
	if err := validate.NotNil(roomSaver); err != nil {
		return nil, err
	}

	return &roomCreatedUseCase{
		roomSaver: roomSaver,
	}, nil
}

func (uc *roomCreatedUseCase) Execute(ctx context.Context, input *RoomCreatedInput) (*kernel.NoOutput, error) {
	if err := uc.roomSaver.Save(ctx, domain.LoadRoom(input.RoomID)); err != nil {
		return nil, err
	}
	return nil, nil
}
