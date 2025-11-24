package application

import (
	"context"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"gochat/internal/social/domain"
)

type LeaveRoomUseCase kernel.UseCase[*LeaveRoomInput, *kernel.NoOutput]

type LeaveRoomInput struct {
	UserID     kernel.UserID
	RoomNumber kernel.RoomNumber
}

func (r *LeaveRoomInput) Validate() error {
	if len(r.UserID) == 0 || len(r.RoomNumber) == 0 {
		return myErrors.ErrEmptyInput
	}
	return nil
}

type leaveRoomUseCase struct {
	eventIDGenerator event.IDGenerator
	finder           domain.RoomFinderByNumber
	roomUpdater      domain.RoomUpdater
}

func NewLeaveRoomUseCase(
	eventIDGenerator event.IDGenerator,
	finder domain.RoomFinderByNumber,
	roomUpdater domain.RoomUpdater,
) LeaveRoomUseCase {
	return &leaveRoomUseCase{
		eventIDGenerator: eventIDGenerator,
		finder:           finder,
		roomUpdater:      roomUpdater,
	}
}

func (uc *leaveRoomUseCase) Execute(ctx context.Context, input *LeaveRoomInput) (*kernel.NoOutput, error) {
	room, err := uc.finder.FindByNumber(ctx, input.RoomNumber)
	if err != nil {
		return nil, err
	}

	if err := room.DeleteMember(
		input.UserID,
		uc.eventIDGenerator,
	); err != nil {
		return nil, err
	}

	if err := uc.roomUpdater.Update(ctx, room); err != nil {
		return nil, err
	}

	return nil, nil
}
