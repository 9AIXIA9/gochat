package application

import (
	"context"
	"gochat/internal/roomship/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

type JoinRoomUseCase kernel.UseCase[*JoinRoomInput, *kernel.NoOutput]

type JoinRoomInput struct {
	UserID     kernel.UserID
	RoomNumber kernel.RoomNumber
	Password   domain.Password
}

func (r *JoinRoomInput) Validate() error {
	if err := r.Password.Validate(); err != nil {
		return err
	}
	if len(r.UserID) == 0 || len(r.RoomNumber) == 0 {
		return myErrors.ErrEmptyInput
	}
	return nil
}

type joinRoomUseCase struct {
	eventIDGenerator event.IDGenerator
	finder           domain.RoomFinderByNumber
	comparator       domain.Comparator
	roomUpdater      domain.RoomUpdater
}

func NewJoinRoomUseCase(
	eventIDGenerator event.IDGenerator,
	finder domain.RoomFinderByNumber,
	comparator domain.Comparator,
	roomUpdater domain.RoomUpdater,
) JoinRoomUseCase {
	return &joinRoomUseCase{
		eventIDGenerator: eventIDGenerator,
		finder:           finder,
		comparator:       comparator,
		roomUpdater:      roomUpdater,
	}
}

func (uc *joinRoomUseCase) Execute(ctx context.Context, input *JoinRoomInput) (*kernel.NoOutput, error) {
	room, err := uc.finder.FindByNumber(ctx, input.RoomNumber)
	if err != nil {
		return nil, err
	}

	if err := room.AddMember(
		input.UserID,
		input.Password,
		uc.comparator,
		uc.eventIDGenerator,
	); err != nil {
		return nil, err
	}

	if err := uc.roomUpdater.Update(ctx, room); err != nil {
		return nil, err
	}

	return nil, nil
}
