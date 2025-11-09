package usecase

import (
	"context"
	"gochat/internal/notification/application"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
)

type UserCreatedUseCase kernel.UseCase[*UserCreatedInput, *kernel.NoOutput]

type UserCreatedInput struct {
	UserID kernel.UserID
	Email  kernel.Email
}

func (r *UserCreatedInput) Validate() error {
	if len(r.UserID) == 0 {
		return myErrors.ErrEmptyInput
	}

	if err := r.Email.Validate(); err != nil {
		return err
	}
	return nil
}

//TODO 后面直接发送邮件

type userCreatedUseCase struct {
	userEmailSaver application.UserEmailSaver
}

func NewUserCreatedUseCase(
	userEmailSaver application.UserEmailSaver,
) UserCreatedUseCase {
	return &userCreatedUseCase{
		userEmailSaver: userEmailSaver,
	}
}

func (uc *userCreatedUseCase) Execute(ctx context.Context, input *UserCreatedInput) (*kernel.NoOutput, error) {
	return nil, uc.userEmailSaver.SaveEmail(ctx, input.UserID, input.Email)
}
