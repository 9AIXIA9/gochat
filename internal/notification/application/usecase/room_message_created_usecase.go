package usecase

import (
	"context"
	"gochat/internal/shared/kernel"
)

type RoomMessageCreatedUseCase kernel.UseCase[*RoomMessageCreatedInput, *kernel.NoOutput]

type RoomMessageCreatedInput struct {
}

func (r *RoomMessageCreatedInput) Validate() error {
	return nil
}

type roomMessageCreatedUseCase struct {
}

func NewRoomMessageCreatedUseCase() RoomMessageCreatedUseCase {
	return &roomMessageCreatedUseCase{}
}

func (uc *roomMessageCreatedUseCase) Execute(context.Context, *RoomMessageCreatedInput) (*kernel.NoOutput, error) {
	return nil, nil
}
