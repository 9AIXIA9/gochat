package application

import (
	"context"
	"gochat/internal/chat/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
)

type RoomJoinedUseCase kernel.UseCase[*RoomJoinedInput, *kernel.NoOutput]

type RoomJoinedInput struct {
	RoomID kernel.RoomID
	UserID kernel.UserID
}

func (r *RoomJoinedInput) Validate() error {
	if len(r.RoomID) == 0 || len(r.UserID) == 0 {
		return myErrors.ErrEmptyInput
	}

	return nil
}

type roomJoinedUseCase struct {
	roomFinder  domain.RoomFinderByID
	roomUpdater domain.RoomUpdater
}

func NewRoomJoinedUseCase(
	roomFinder domain.RoomFinderByID,
	roomUpdater domain.RoomUpdater,
) RoomJoinedUseCase {
	return &roomJoinedUseCase{
		roomFinder:  roomFinder,
		roomUpdater: roomUpdater,
	}
}

func (uc *roomJoinedUseCase) Execute(ctx context.Context, input *RoomJoinedInput) (*kernel.NoOutput, error) {
	room, err := uc.roomFinder.FindByID(ctx, input.RoomID)
	if err != nil {
		return nil, err
	}

	room.AddMember(input.UserID)

	return nil, uc.roomUpdater.Update(ctx, room)
}
