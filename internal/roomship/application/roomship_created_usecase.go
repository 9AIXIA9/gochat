package application

import (
	"context"
	"gochat/internal/roomship/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
)

type RoomshipCreatedUseCase kernel.UseCase[*RoomshipCreatedInput, *kernel.NoOutput]

type RoomshipCreatedInput struct {
	RoomshipID domain.RoomshipID
}

func (r *RoomshipCreatedInput) Validate() error {
	if len(r.RoomshipID) == 0 {
		return myErrors.ErrEmptyInput
	}

	return nil
}

type roomshipCreatedUseCase struct {
}

func NewRoomshipCreatedUseCase() (RoomshipCreatedUseCase, error) {
	return &roomshipCreatedUseCase{}, nil
}

func (uc *roomshipCreatedUseCase) Execute(ctx context.Context, input *RoomshipCreatedInput) (*kernel.NoOutput, error) {
	//TODO 通知上下文进行通知
	return nil, nil
}
