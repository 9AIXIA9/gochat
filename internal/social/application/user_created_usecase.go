package application

import (
	"context"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
	"gochat/internal/social/domain"
)

type UserCreatedUseCase kernel.UseCase[*UserCreatedInput, *kernel.NoOutput]

type UserCreatedInput struct {
	UserID kernel.UserID
}

func (r *UserCreatedInput) Validate() error {
	if len(r.UserID) == 0 {
		return myErrors.ErrEmptyInput
	}

	return nil
}

type userCreatedUseCase struct {
	userIDSaver domain.UserIDSaver
}

func NewUserCreatedUseCase(
	userIDSaver domain.UserIDSaver,
) UserCreatedUseCase {
	return &userCreatedUseCase{
		userIDSaver: userIDSaver,
	}
}

func (uc *userCreatedUseCase) Execute(ctx context.Context, input *UserCreatedInput) (*kernel.NoOutput, error) {
	return nil, uc.userIDSaver.SaveID(ctx, input.UserID)
}
