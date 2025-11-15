package usecase

import (
	"context"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
	"gochat/internal/social/application"
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
	userIDSaver application.UserIDSaver
}

func NewUserCreatedUseCase(
	userIDSaver application.UserIDSaver,
) UserCreatedUseCase {
	return &userCreatedUseCase{
		userIDSaver: userIDSaver,
	}
}

func (uc *userCreatedUseCase) Execute(ctx context.Context, input *UserCreatedInput) (*kernel.NoOutput, error) {
	return nil, uc.userIDSaver.SaveID(ctx, input.UserID)
}
