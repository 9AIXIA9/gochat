package usecase

import (
	"context"
	"gochat/internal/notification/application"
	"gochat/internal/notification/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
)

type RoomCreatedUseCase kernel.UseCase[*RoomCreatedInput, *kernel.NoOutput]

type RoomCreatedInput struct {
	RoomID domain.RoomID
}

func (r *RoomCreatedInput) Validate() error {
	if len(r.RoomID) == 0 {
		return myErrors.ErrEmptyInput
	}

	return nil
}

type roomCreatedUseCase struct {
	roomIDSaver application.RoomIDSaver
}

func NewRoomCreatedUseCase(
	roomIDSaver application.RoomIDSaver,
) RoomCreatedUseCase {
	return &roomCreatedUseCase{
		roomIDSaver: roomIDSaver,
	}
}

func (uc *roomCreatedUseCase) Execute(ctx context.Context, input *RoomCreatedInput) (*kernel.NoOutput, error) {
	return nil, uc.roomIDSaver.SaveID(ctx, input.RoomID)
}
