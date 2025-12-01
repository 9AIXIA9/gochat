package application

import (
	"context"
	"gochat/internal/friendship/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
	"gochat/pkg/utils"
)

type UserCreatedUseCase kernel.UseCase[*UserCreatedInput, *kernel.NoOutput]

type UserCreatedInput struct {
	UserID     kernel.UserID
	UserNumber kernel.UserNumber
}

func (r *UserCreatedInput) Validate() error {
	if len(r.UserID) == 0 {
		return myErrors.ErrEmptyInput
	}

	if err := r.UserNumber.Validate(); err != nil {
		return err
	}

	return nil
}

type userCreatedUseCase struct {
	userCreator domain.UserCreator
}

func NewUserCreatedUseCase(
	userCreator domain.UserCreator,
) (UserCreatedUseCase, error) {
	if err := utils.CheckInterfaces(userCreator); err != nil {
		return nil, err
	}
	return &userCreatedUseCase{
		userCreator: userCreator,
	}, nil
}

func (uc *userCreatedUseCase) Execute(ctx context.Context, input *UserCreatedInput) (*kernel.NoOutput, error) {
	user := domain.CreateUser(input.UserID, input.UserNumber)
	return nil, uc.userCreator.Create(ctx, user)
}
