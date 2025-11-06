package usecase

import (
	"context"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"gochat/internal/social/application"
	"gochat/internal/social/domain"
)

type LeaveRoomUseCase kernel.UseCase[*LeaveRoomInput, *kernel.NoOutput]

type LeaveRoomInput struct {
	UserID     kernel.UserID
	RoomNumber domain.RoomNumber
}

func (r *LeaveRoomInput) Validate() error {
	if len(r.UserID) == 0 || len(r.RoomNumber) == 0 {
		return myErrors.ErrEmptyInput
	}
	return nil
}

type leaveRoomUseCase struct {
	eventIDGenerator event.IDGenerator
	finder           application.RoomFinder
	roomLeaver       application.RoomLeaver
	eventSaver       event.UnpublishedSaver
}

func NewLeaveRoomUseCase(
	eventIDGenerator event.IDGenerator,
	finder application.RoomFinder,
	roomLeaver application.RoomLeaver,
	eventSaver event.UnpublishedSaver,
) LeaveRoomUseCase {
	return &leaveRoomUseCase{
		eventIDGenerator: eventIDGenerator,
		finder:           finder,
		roomLeaver:       roomLeaver,
		eventSaver:       eventSaver,
	}
}

func (uc *leaveRoomUseCase) Execute(ctx context.Context, input *LeaveRoomInput) (*kernel.NoOutput, error) {
	room, err := uc.finder.FindByNumber(ctx, input.RoomNumber)
	if err != nil {
		return nil, err
	}

	if err := room.Leave(input.UserID, uc.eventIDGenerator); err != nil {
		return nil, err
	}

	if err := uc.roomLeaver.Leave(ctx, room.ID(), input.UserID); err != nil {
		return nil, err
	}

	if err := uc.eventSaver.Saves(ctx, room.GetEvents()); err != nil {
		return nil, err
	}

	return nil, nil
}
