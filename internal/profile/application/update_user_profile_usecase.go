package application

import (
	"context"
	"gochat/internal/profile/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
	"gochat/pkg/utils"
)

type UpdateUserProfileUseCase kernel.UseCase[*UpdateUserProfileInput, *kernel.NoOutput]

type UpdateUserProfileInput struct {
	UserID      kernel.UserID
	Name        string
	Gender      kernel.Gender
	Email       kernel.Email
	PhoneNumber kernel.PhoneNumber
	Address     *kernel.Address
	Sign        string
}

func (r *UpdateUserProfileInput) Validate() error {
	if len(r.UserID) == 0 {
		return myErrors.ErrEmptyInput
	}

	if len(r.Email) == 0 &&
		len(r.PhoneNumber) == 0 &&
		len(r.Name) == 0 &&
		r.Address == nil &&
		r.Gender == kernel.UnknownGender &&
		len(r.Sign) == 0 {
		return myErrors.ErrEmptyInput
	}

	return nil
}

type updateUserProfileUseCase struct {
	updater domain.UserProfileUpdater
	finder  domain.UserProfileFinder
}

func NewUpdateUserProfileUseCase(
	updater domain.UserProfileUpdater,
	finder domain.UserProfileFinder,
) (UpdateUserProfileUseCase, error) {
	if err := utils.CheckInterfaces(
		updater,
		finder,
	); err != nil {
		return nil, err
	}

	return &updateUserProfileUseCase{
		updater: updater,
		finder:  finder,
	}, nil
}

func (uc *updateUserProfileUseCase) Execute(ctx context.Context, input *UpdateUserProfileInput) (*kernel.NoOutput, error) {
	profile, err := uc.finder.FindByID(ctx, input.UserID)
	if err != nil {
		return nil, err
	}

	uc.updatesProfile(profile, input)

	if err := uc.updater.Update(ctx, profile); err != nil {
		return nil, err
	}

	return nil, nil
}

func (uc *updateUserProfileUseCase) updatesProfile(profile *domain.UserProfile, input *UpdateUserProfileInput) {
	if input.Name != "" {
		profile.UpdateName(input.Name)
	}
	if input.Gender != kernel.UnknownGender {
		profile.UpdateGender(input.Gender)
	}
	if input.Email != "" {
		profile.UpdateEmail(input.Email)
	}
	if input.PhoneNumber != "" {
		profile.UpdatePhoneNumber(input.PhoneNumber)
	}
	if input.Address != nil {
		profile.UpdateAddress(input.Address)
	}
	if input.Sign != "" {
		profile.UpdateSign(input.Sign)
	}
}
