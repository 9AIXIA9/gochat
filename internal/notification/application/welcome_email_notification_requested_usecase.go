package application

import (
	"context"
	"gochat/internal/notification/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
)

type WelcomeEmailNotificationRequestedUseCase kernel.UseCase[*WelcomeEmailNotificationRequestedInput, *kernel.NoOutput]

type WelcomeEmailNotificationRequestedInput struct {
	UserID     kernel.UserID
	UserNumber kernel.UserNumber
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
	emailNotifier domain.WelcomeEmailNotifier
}

func NewWelcomeEmailNotificationRequestedUseCase(emailNotifier domain.WelcomeEmailNotifier) WelcomeEmailNotificationRequestedUseCase {
	return &welcomeEmailNotificationRequestedUseCase{emailNotifier: emailNotifier}
}

func (uc *welcomeEmailNotificationRequestedUseCase) Execute(ctx context.Context, input *WelcomeEmailNotificationRequestedInput) (*kernel.NoOutput, error) {
	if err := uc.emailNotifier.NotifyWelcomeEmail(ctx, input.Email, input.UserNumber); err != nil {
		return nil, err
	}
	return nil, nil
}
