package application

import (
	"context"
	"gochat/internal/chat/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
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
	userNumberSaver domain.UserNumberSaver
}

func NewUserCreatedUseCase(
	userNumberSaver domain.UserNumberSaver,
) UserCreatedUseCase {
	return &userCreatedUseCase{
		userNumberSaver: userNumberSaver,
	}
}

func (uc *userCreatedUseCase) Execute(ctx context.Context, input *UserCreatedInput) (*kernel.NoOutput, error) {
	return nil, uc.userNumberSaver.SaveNumber(ctx, input.UserID, input.UserNumber)
}
