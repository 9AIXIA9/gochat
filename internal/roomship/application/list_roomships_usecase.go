package application

import (
	"context"
	"gochat/internal/roomship/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
	"gochat/pkg/utils"
)

const (
	defaultRoomshipsLimit = 50
	maxRoomshipsLimit     = 100
)

type ListRoomshipsUseCase kernel.UseCase[*ListRoomshipsInput, *ListRoomshipsOutput]

type ListRoomshipsInput struct {
	UserID kernel.UserID
	BaseID domain.RoomshipID
	Limit  int
}

func (r *ListRoomshipsInput) Validate() error {
	if len(r.UserID) == 0 {
		return myErrors.WrapBusiness(myErrors.ErrEmptyInput, "user id is empty")
	}

	return nil
}

type ListRoomshipsOutput struct {
	Roomships []*domain.Roomship
}

type listRoomshipsUseCase struct {
	roomshipFinder domain.RoomshipsFinderByUserID
}

func NewListRoomshipsUseCase(
	roomshipFinder domain.RoomshipsFinderByUserID,
) (ListRoomshipsUseCase, error) {
	if err := utils.CheckInterfaces(roomshipFinder); err != nil {
		return nil, err
	}
	return &listRoomshipsUseCase{
		roomshipFinder: roomshipFinder,
	}, nil
}

func (uc *listRoomshipsUseCase) Execute(ctx context.Context, input *ListRoomshipsInput) (*ListRoomshipsOutput, error) {
	limit := input.Limit
	switch {
	case limit <= 0:
		limit = defaultRoomshipsLimit
	case limit > maxRoomshipsLimit:
		limit = maxRoomshipsLimit
	}

	roomships, err := uc.roomshipFinder.FindsByUserID(ctx, input.UserID, limit, input.BaseID)
	if err != nil {
		return nil, err
	}

	return &ListRoomshipsOutput{
		Roomships: roomships,
	}, nil
}
