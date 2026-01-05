package application

import (
	"context"
	"errors"
	"gochat/internal/profile/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
	"gochat/pkg/validate"
)

type UpdateUserProfileUseCase kernel.UseCase[*UpdateUserProfileInput, *kernel.NoOutput]

type UpdateUserProfileInput struct {
	UserID      kernel.UserID
	Name        string
	Gender      kernel.Gender
	Email       kernel.Email
	PhoneNumber kernel.PhoneNumber
	Address     kernel.Address
	Sign        string
}

func (r *UpdateUserProfileInput) Validate() error {
	if len(r.UserID) == 0 {
		return myErrors.WrapBusiness(myErrors.ErrEmptyInput, "user id is empty")
	}
	if len(r.Email) == 0 &&
		len(r.PhoneNumber) == 0 &&
		len(r.Name) == 0 &&
		len(r.Address) == 0 &&
		len(r.Sign) == 0 &&
		r.Gender == kernel.UnknownGender {
		return myErrors.WrapBusiness(myErrors.ErrEmptyInput, "everything is empty")
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
	if err := validate.NotNil(
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
		if errors.Is(err, myErrors.ErrNotFound) {
			return nil, myErrors.WrapBusiness(err, "user profile not found")
		}
		return nil, err
	}

	if err := uc.updatesProfile(profile, input); err != nil {
		return nil, err
	}

	if err := uc.updater.Update(ctx, profile); err != nil {
		return nil, err
	}

	return nil, nil
}

func (uc *updateUserProfileUseCase) updatesProfile(profile *domain.UserProfile, input *UpdateUserProfileInput) error {
	if input.Name != "" {
		if err := profile.UpdateName(input.Name); err != nil {
			return err
		}
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
	if input.Address != "" {
		if err := profile.UpdateAddress(input.Address); err != nil {
			return err
		}
	}
	if input.Sign != "" {
		if err := profile.UpdateSign(input.Sign); err != nil {
			return err
		}
	}
	return nil
}
