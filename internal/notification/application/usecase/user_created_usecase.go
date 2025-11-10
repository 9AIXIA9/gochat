package usecase

import (
	"context"
	"gochat/internal/notification/application"
	"gochat/internal/notification/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
)

type UserCreatedUseCase kernel.UseCase[*UserCreatedInput, *kernel.NoOutput]

type UserCreatedInput struct {
	UserID     kernel.UserID
	UserNumber domain.UserNumber
	Email      kernel.Email
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

type userCreatedUseCase struct {
	userIDSaver   application.UserIDSaver
	emailNotifier application.UserCreatedEmailNotifier
}

func NewUserCreatedUseCase(userIDSaver application.UserIDSaver, emailNotifier application.UserCreatedEmailNotifier) UserCreatedUseCase {
	return &userCreatedUseCase{userIDSaver: userIDSaver, emailNotifier: emailNotifier}
}

func (uc *userCreatedUseCase) Execute(ctx context.Context, input *UserCreatedInput) (*kernel.NoOutput, error) {
	if err := uc.userIDSaver.SaveID(ctx, input.UserID); err != nil {
		return nil, err
	}

	if err := uc.emailNotifier.AddUserCreatedEmail(ctx, input.Email, input.UserNumber); err != nil {
		return nil, err
	}
	return nil, nil
}
