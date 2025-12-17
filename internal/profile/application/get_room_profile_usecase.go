package application

import (
	"context"
	"errors"
	"gochat/internal/profile/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
	"gochat/pkg/utils"
)

type GetRoomProfileUseCase kernel.UseCase[*GetRoomProfileInput, *GetRoomProfileOutput]

type GetRoomProfileInput struct {
	RoomID kernel.RoomID
}

func (r *GetRoomProfileInput) Validate() error {
	if len(r.RoomID) == 0 {
		return myErrors.WrapBusiness(myErrors.ErrEmptyInput, "room id is empty")
	}
	return nil
}

type GetRoomProfileOutput struct {
	Profile *domain.RoomProfile
}

type getRoomProfileUseCase struct {
	finder domain.RoomProfileFinder
}

func NewGetRoomProfileUseCase(
	finder domain.RoomProfileFinder,
) (GetRoomProfileUseCase, error) {
	if err := utils.CheckInterfaces(
		finder,
	); err != nil {
		return nil, err
	}

	return &getRoomProfileUseCase{
		finder: finder,
	}, nil
}

func (uc *getRoomProfileUseCase) Execute(ctx context.Context, input *GetRoomProfileInput) (*GetRoomProfileOutput, error) {
	profile, err := uc.finder.FindByID(ctx, input.RoomID)
	if err != nil {
		if errors.Is(err, myErrors.ErrNotFound) {
			return nil, myErrors.WrapBusiness(err, "room profile not found")
		}
		return nil, err
	}

	return &GetRoomProfileOutput{
		Profile: profile,
	}, nil
}
