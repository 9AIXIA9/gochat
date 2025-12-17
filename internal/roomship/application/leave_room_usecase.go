package application

import (
	"context"
	"errors"
	"gochat/internal/roomship/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
	"gochat/pkg/utils"
)

type LeaveRoomUseCase kernel.UseCase[*LeaveRoomInput, *kernel.NoOutput]

type LeaveRoomInput struct {
	UserID kernel.UserID
	RoomID kernel.RoomID
}

func (r *LeaveRoomInput) Validate() error {
	if len(r.UserID) == 0 || len(r.RoomID) == 0 {
		return myErrors.WrapBusiness(myErrors.ErrEmptyInput, "user id or room id is empty")
	}
	return nil
}

type leaveRoomUseCase struct {
	roomshipDeleter domain.RoomshipDeleterByUserIDAndRoomID
	roomDeleter     domain.RoomDeleterByID
	roomshipFinder  domain.RoomshipFinderByUserIDAndRoomID
	roomshipsFinder domain.RoomshipsFinderByRoomID
}

func NewLeaveRoomUseCase(
	roomshipDeleter domain.RoomshipDeleterByUserIDAndRoomID,
	roomDeleter domain.RoomDeleterByID,
	roomshipFinder domain.RoomshipFinderByUserIDAndRoomID,
	roomshipsFinder domain.RoomshipsFinderByRoomID,
) (LeaveRoomUseCase, error) {
	if err := utils.CheckInterfaces(
		roomshipDeleter, roomshipFinder, roomDeleter, roomshipsFinder,
	); err != nil {
		return nil, err
	}

	return &leaveRoomUseCase{
		roomshipDeleter: roomshipDeleter,
		roomDeleter:     roomDeleter,
		roomshipFinder:  roomshipFinder,
		roomshipsFinder: roomshipsFinder,
	}, nil
}

func (uc *leaveRoomUseCase) Execute(ctx context.Context, input *LeaveRoomInput) (*kernel.NoOutput, error) {
	roomship, err := uc.roomshipFinder.FindByUserIDAndRoomID(ctx, input.UserID, input.RoomID)
	if err != nil {
		if errors.Is(err, myErrors.ErrNotFound) {
			return nil, nil
		}
		return nil, err
	}

	if roomship.IsOwner() {
		return nil, domain.ErrOwnerCantLeave
	}

	if err := uc.roomshipDeleter.DeleteByUserIDAndRoomID(ctx, input.UserID, input.RoomID); err != nil {
		return nil, err
	}

	roomships, err := uc.roomshipsFinder.FindsByRoomID(ctx, input.RoomID)
	if err != nil {
		return nil, err
	}

	if len(roomships) != 0 {
		return nil, nil
	}

	if err := uc.roomDeleter.DeleteByID(ctx, input.RoomID); err != nil {
		return nil, err
	}
	return nil, nil
}
