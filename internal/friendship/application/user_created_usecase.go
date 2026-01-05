package application

import (
	"context"
	"gochat/internal/friendship/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
	"gochat/pkg/validate"
)

type UserCreatedUseCase kernel.UseCase[*UserCreatedInput, *kernel.NoOutput]

type UserCreatedInput struct {
	UserID kernel.UserID
}

func (r *UserCreatedInput) Validate() error {
	if len(r.UserID) == 0 {
		return myErrors.WrapBusiness(myErrors.ErrEmptyInput, "user id can't be empty")
	}

	return nil
}

type userCreatedUseCase struct {
	userSaver domain.UserSaver
}

func NewUserCreatedUseCase(
	userSaver domain.UserSaver,
) (UserCreatedUseCase, error) {
	if err := validate.NotNil(userSaver); err != nil {
		return nil, err
	}
	return &userCreatedUseCase{
		userSaver: userSaver,
	}, nil
}

func (uc *userCreatedUseCase) Execute(ctx context.Context, input *UserCreatedInput) (*kernel.NoOutput, error) {
	return nil, uc.userSaver.Save(ctx, domain.LoadUser(input.UserID))
}
