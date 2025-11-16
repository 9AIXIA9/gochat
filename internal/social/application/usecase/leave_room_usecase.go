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
	eventIDGenerator  event.IDGenerator
	finder            application.RoomFinderByNumber
	roomMemberDeleter application.RoomMemberDeleter
	eventSaver        event.UnpublishedSaver
	unitOfWork        kernel.UnitOfWork
}

func NewLeaveRoomUseCase(
	eventIDGenerator event.IDGenerator,
	finder application.RoomFinderByNumber,
	roomMemberDeleter application.RoomMemberDeleter,
	eventSaver event.UnpublishedSaver,
	unitOfWork kernel.UnitOfWork,
) LeaveRoomUseCase {
	return &leaveRoomUseCase{
		eventIDGenerator:  eventIDGenerator,
		finder:            finder,
		roomMemberDeleter: roomMemberDeleter,
		eventSaver:        eventSaver,
		unitOfWork:        unitOfWork,
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

	if err := uc.unitOfWork.Execute(ctx, func(txCtx context.Context) error {
		if err := uc.roomMemberDeleter.DeleteMember(txCtx, room.ID(), input.UserID); err != nil {
			return err
		}

		if err := uc.eventSaver.Saves(txCtx, room.GetEvents()); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return nil, nil
}
