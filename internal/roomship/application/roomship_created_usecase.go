package application

import (
	"context"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
)

type RoomshipCreatedUseCase kernel.UseCase[*RoomshipCreatedInput, *kernel.NoOutput]

type RoomshipCreatedInput struct {
	UserID kernel.UserID
}

func (r *RoomshipCreatedInput) Validate() error {
	if len(r.UserID) == 0 {
		return myErrors.ErrEmptyInput
	}

	return nil
}

type roomshipCreatedUseCase struct {
}

func NewRoomshipCreatedUseCase() RoomshipCreatedUseCase {
	return &roomshipCreatedUseCase{}
}

func (uc *roomshipCreatedUseCase) Execute(ctx context.Context, input *RoomshipCreatedInput) (*kernel.NoOutput, error) {
	//TODO 通知上下文进行通知
	return nil, nil
}
