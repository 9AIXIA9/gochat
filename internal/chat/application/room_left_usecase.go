package application

import (
	"context"
	"gochat/internal/chat/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
)

type RoomLeftUseCase kernel.UseCase[*RoomLeftInput, *kernel.NoOutput]

type RoomLeftInput struct {
	RoomID kernel.RoomID
	UserID kernel.UserID
}

func (r *RoomLeftInput) Validate() error {
	if len(r.RoomID) == 0 || len(r.UserID) == 0 {
		return myErrors.ErrEmptyInput
	}

	return nil
}

type roomLeftUseCase struct {
	roomFinder  domain.RoomFinderByID
	roomUpdater domain.RoomUpdater
}

func NewRoomLeftUseCase(
	roomFinder domain.RoomFinderByID,
	roomUpdater domain.RoomUpdater,
) RoomLeftUseCase {
	return &roomLeftUseCase{
		roomFinder:  roomFinder,
		roomUpdater: roomUpdater,
	}
}

func (uc *roomLeftUseCase) Execute(ctx context.Context, input *RoomLeftInput) (*kernel.NoOutput, error) {
	room, err := uc.roomFinder.FindByID(ctx, input.RoomID)
	if err != nil {
		return nil, err
	}

	room.DeleteMember(input.UserID)

	return nil, uc.roomUpdater.Update(ctx, room)
}
