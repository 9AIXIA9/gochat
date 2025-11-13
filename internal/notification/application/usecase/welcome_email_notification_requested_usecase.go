package usecase

import (
	"context"
	"gochat/internal/notification/application"
	"gochat/internal/notification/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
)

type WelcomeEmailNotificationRequestedUseCase kernel.UseCase[*WelcomeEmailNotificationRequestedInput, *kernel.NoOutput]

type WelcomeEmailNotificationRequestedInput struct {
	UserID     kernel.UserID
	UserNumber domain.UserNumber
	Email      kernel.Email
}

func (r *WelcomeEmailNotificationRequestedInput) Validate() error {
	if len(r.UserID) == 0 {
		return myErrors.ErrEmptyInput
	}

	if err := r.Email.Validate(); err != nil {
		return err
	}
	return nil
}

type welcomeEmailNotificationRequestedUseCase struct {
	emailNotifier application.WelcomeEmailNotifier
}

func NewWelcomeEmailNotificationRequestedUseCase(emailNotifier application.WelcomeEmailNotifier) WelcomeEmailNotificationRequestedUseCase {
	return &welcomeEmailNotificationRequestedUseCase{emailNotifier: emailNotifier}
}

func (uc *welcomeEmailNotificationRequestedUseCase) Execute(ctx context.Context, input *WelcomeEmailNotificationRequestedInput) (*kernel.NoOutput, error) {
	if err := uc.emailNotifier.NotifyWelcomeEmail(ctx, input.Email, input.UserNumber); err != nil {
		return nil, err
	}
	return nil, nil
}
