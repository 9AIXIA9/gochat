package application

import (
	"context"
	"errors"
	"gochat/internal/profile/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
	"gochat/pkg/validate"
	"time"
)

type UserCreatedUseCase kernel.UseCase[*UserCreatedInput, *kernel.NoOutput]

type UserCreatedInput struct {
	UserID     kernel.UserID
	Email      kernel.Email
	SignedUpAt time.Time
}

func (r *UserCreatedInput) Validate() error {
	if len(r.UserID) == 0 || r.SignedUpAt.IsZero() {
		return myErrors.WrapBusiness(myErrors.ErrEmptyInput, "user id or signed up time is empty")
	}

	if err := r.Email.Validate(); err != nil {
		return err
	}

	return nil
}

type userCreatedUseCase struct {
	userSaver      domain.UserSaver
	profileCreator domain.UserProfileCreator
}

func NewUserCreatedUseCase(
	userSaver domain.UserSaver,
	profileCreator domain.UserProfileCreator,
) (UserCreatedUseCase, error) {
	if err := validate.NotNil(
		userSaver,
		profileCreator,
	); err != nil {
		return nil, err
	}

	return &userCreatedUseCase{
		userSaver:      userSaver,
		profileCreator: profileCreator,
	}, nil
}

func (uc *userCreatedUseCase) Execute(ctx context.Context, input *UserCreatedInput) (*kernel.NoOutput, error) {
	if err := uc.userSaver.Save(ctx, domain.LoadUser(input.UserID)); err != nil {
		return nil, err
	}

	profile := domain.CreateUserProfile(
		input.UserID,
		input.Email,
		input.SignedUpAt,
	)

	if err := uc.profileCreator.Create(ctx, profile); err != nil {
		if errors.Is(err, myErrors.ErrDuplicatedKey) {
			return nil, nil
		}
		return nil, err
	}
	return nil, nil
}
