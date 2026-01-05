package application

import (
	"context"
	"gochat/internal/roomship/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
	"gochat/pkg/validate"
)

type ListRoomMembersUseCase kernel.UseCase[*ListRoomMembersInput, *ListRoomMembersOutput]

type ListRoomMembersInput struct {
	RoomID kernel.RoomID
}

func (r *ListRoomMembersInput) Validate() error {
	if len(r.RoomID) == 0 {
		return myErrors.WrapBusiness(myErrors.ErrEmptyInput, "room id is empty")
	}

	return nil
}

type ListRoomMembersOutput struct {
	Roomships []*domain.Roomship
}

type listRoomMembersUseCase struct {
	roomshipFinder domain.RoomshipsFinderByRoomID
}

func NewListRoomMembersUseCase(
	roomshipFinder domain.RoomshipsFinderByRoomID,
) (ListRoomMembersUseCase, error) {
	if err := validate.NotNil(roomshipFinder); err != nil {
		return nil, err
	}
	return &listRoomMembersUseCase{
		roomshipFinder: roomshipFinder,
	}, nil
}

func (uc *listRoomMembersUseCase) Execute(ctx context.Context, input *ListRoomMembersInput) (*ListRoomMembersOutput, error) {
	roomships, err := uc.roomshipFinder.FindsByRoomID(ctx, input.RoomID)
	if err != nil {
		return nil, err
	}

	return &ListRoomMembersOutput{
		Roomships: roomships,
	}, nil
}
