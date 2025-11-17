package usecase

import (
	"context"
	"gochat/internal/chat/application"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
)

type RoomCreatedUseCase kernel.UseCase[*RoomCreatedInput, *kernel.NoOutput]

type RoomCreatedInput struct {
	RoomID     kernel.RoomID
	RoomNumber kernel.RoomNumber
}

func (r *RoomCreatedInput) Validate() error {
	if len(r.RoomID) == 0 {
		return myErrors.ErrEmptyInput
	}

	if err := r.RoomNumber.Validate(); err != nil {
		return err
	}
	return nil
}

type roomCreatedUseCase struct {
	roomNumberSaver application.RoomNumberSaver
}

func NewRoomCreatedUseCase(
	roomNumberSaver application.RoomNumberSaver,
) RoomCreatedUseCase {
	return &roomCreatedUseCase{
		roomNumberSaver: roomNumberSaver,
	}
}

func (uc *roomCreatedUseCase) Execute(ctx context.Context, input *RoomCreatedInput) (*kernel.NoOutput, error) {
	return nil, uc.roomNumberSaver.SaveNumber(ctx, input.RoomID, input.RoomNumber)
}
