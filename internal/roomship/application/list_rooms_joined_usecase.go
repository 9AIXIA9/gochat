package application

import (
	"context"
	"gochat/internal/roomship/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
	"gochat/pkg/utils"
)

const (
	defaultRoomsJoinedListLimit = 10
	maxRoomsJoinedListLimit     = 20
)

type ListRoomsJoinedUseCase kernel.UseCase[*ListRoomsJoinedInput, *ListRoomsJoinedOutput]

type ListRoomsJoinedInput struct {
	UserID kernel.UserID
	BaseID domain.RoomshipID
	Limit  int
}

func (r *ListRoomsJoinedInput) Validate() error {
	if len(r.UserID) == 0 {
		return myErrors.ErrEmptyInput
	}

	return nil
}

type ListRoomsJoinedOutput struct {
	Rooms []*domain.Room
}

type listRoomsJoinedUseCase struct {
	roomshipFinder domain.RoomshipsFinderByUserID
	roomsFinder    domain.RoomsFinderByIDs
}

func NewListRoomsJoinedUseCase(
	roomshipFinder domain.RoomshipsFinderByUserID,
	roomsFinder domain.RoomsFinderByIDs,
) (ListRoomsJoinedUseCase, error) {
	if err := utils.CheckInterfaces(roomshipFinder); err != nil {
		return nil, err
	}
	return &listRoomsJoinedUseCase{
		roomshipFinder: roomshipFinder,
		roomsFinder:    roomsFinder,
	}, nil
}

func (uc *listRoomsJoinedUseCase) Execute(ctx context.Context, input *ListRoomsJoinedInput) (*ListRoomsJoinedOutput, error) {
	limit := input.Limit
	switch {
	case limit <= 0:
		limit = defaultRoomsJoinedListLimit
	case limit > maxRoomsJoinedListLimit:
		limit = maxRoomsJoinedListLimit
	}

	roomships, err := uc.roomshipFinder.FindsByUserID(ctx, input.UserID, limit, input.BaseID)
	if err != nil {
		return nil, err
	}

	roomIDs := make([]kernel.RoomID, 0, len(roomships))
	for _, roomship := range roomships {
		roomIDs = append(roomIDs, roomship.RoomID())
	}

	rooms, err := uc.roomsFinder.FindsByIDs(ctx, roomIDs)
	if err != nil {
		return nil, err
	}

	return &ListRoomsJoinedOutput{
		rooms,
	}, nil
}
