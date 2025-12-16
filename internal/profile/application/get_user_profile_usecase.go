package application

import (
	"context"
	"gochat/internal/profile/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
	"gochat/pkg/utils"
)

type GetUserProfileUseCase kernel.UseCase[*GetUserProfileInput, *GetUserProfileOutput]

type GetUserProfileInput struct {
	UserID kernel.UserID
}

func (r *GetUserProfileInput) Validate() error {
	if len(r.UserID) == 0 {
		return myErrors.WrapBusiness(myErrors.ErrEmptyInput, "user id is empty")
	}
	return nil
}

type GetUserProfileOutput struct {
	Profile *domain.UserProfile
}

type getUserProfileUseCase struct {
	finder domain.UserProfileFinder
}

func NewGetUserProfileUseCase(
	finder domain.UserProfileFinder,
) (GetUserProfileUseCase, error) {
	if err := utils.CheckInterfaces(
		finder,
	); err != nil {
		return nil, err
	}

	return &getUserProfileUseCase{
		finder: finder,
	}, nil
}

func (uc *getUserProfileUseCase) Execute(ctx context.Context, input *GetUserProfileInput) (*GetUserProfileOutput, error) {
	profile, err := uc.finder.FindByID(ctx, input.UserID)
	if err != nil {
		return nil, err
	}

	return &GetUserProfileOutput{
		Profile: profile,
	}, nil
}
