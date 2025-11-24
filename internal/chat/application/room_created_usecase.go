package application

import (
	"context"
	"gochat/internal/chat/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
)

type RoomCreatedUseCase kernel.UseCase[*RoomCreatedInput, *kernel.NoOutput]

type RoomCreatedInput struct {
	OwnerID    kernel.UserID
	RoomID     kernel.RoomID
	RoomNumber kernel.RoomNumber
}

func (r *RoomCreatedInput) Validate() error {
	if len(r.RoomID) == 0 || len(r.OwnerID) == 0 {
		return myErrors.ErrEmptyInput
	}

	if err := r.RoomNumber.Validate(); err != nil {
		return err
	}
	return nil
}

type roomCreatedUseCase struct {
	roomCreator domain.RoomCreator
}

func NewRoomCreatedUseCase(
	roomCreator domain.RoomCreator,
) RoomCreatedUseCase {
	return &roomCreatedUseCase{
		roomCreator: roomCreator,
	}
}

func (uc *roomCreatedUseCase) Execute(ctx context.Context, input *RoomCreatedInput) (*kernel.NoOutput, error) {
	room := domain.CreateRoom(input.RoomID, input.RoomNumber, input.OwnerID)

	if err := uc.roomCreator.Create(ctx, room); err != nil {
		return nil, err
	}
	return nil, nil
}
