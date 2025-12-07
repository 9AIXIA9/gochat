package application_test

import (
	"gochat/internal/notification/application"
	"gochat/internal/notification/domain/mocks"
	myErrors "gochat/internal/shared/errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestWelcomeEmailNotificationRequestedInput_Validate(t *testing.T) {
	input := &application.WelcomeEmailNotificationRequestedInput{
		UserID:     fixedUserID,
		UserNumber: fixedUserNumber,
		Email:      fixedEmail,
	}

	err := input.Validate()
	require.NoError(t, err)

	inputWithEmptyUserID := &application.WelcomeEmailNotificationRequestedInput{
		UserID:     "",
		UserNumber: fixedUserNumber,
		Email:      fixedEmail,
	}

	err = inputWithEmptyUserID.Validate()
	require.ErrorIs(t, err, myErrors.ErrEmptyInput)

	inputWithEmptyUserNumber := &application.WelcomeEmailNotificationRequestedInput{
		UserID:     fixedUserID,
		UserNumber: "",
		Email:      fixedEmail,
	}

	err = inputWithEmptyUserNumber.Validate()
	require.ErrorIs(t, err, myErrors.ErrEmptyInput)

	inputWithEmptyEmail := &application.WelcomeEmailNotificationRequestedInput{
		UserID:     fixedUserID,
		UserNumber: fixedUserNumber,
		Email:      "",
	}

	err = inputWithEmptyEmail.Validate()
	require.ErrorIs(t, err, myErrors.ErrInvalidFormat)
}

func TestNewWelcomeEmailNotificationRequestedUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockNotifier := mocks.NewMockWelcomeEmailNotifier(ctrl)

	useCase, err := application.NewWelcomeEmailNotificationRequestedUseCase(
		mockNotifier,
	)

	require.NoError(t, err)
	require.NotNil(t, useCase)

	useCaseWithNil, err := application.NewWelcomeEmailNotificationRequestedUseCase(
		nil,
	)

	require.ErrorIs(t, err, myErrors.ErrEmptyPointer)
	require.Nil(t, useCaseWithNil)
}

func TestWelcomeEmailNotificationRequestedUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockNotifier := mocks.NewMockWelcomeEmailNotifier(ctrl)

	useCase, err := application.NewWelcomeEmailNotificationRequestedUseCase(
		mockNotifier,
	)
	require.NoError(t, err)
	require.NotNil(t, useCase)

	mockNotifier.EXPECT().NotifyWelcomeEmail(nil, fixedEmail, fixedUserNumber).Return(nil).Times(1)

	_, err = useCase.Execute(nil, &application.WelcomeEmailNotificationRequestedInput{
		UserID:     fixedUserID,
		UserNumber: fixedUserNumber,
		Email:      fixedEmail,
	})
	require.NoError(t, err)
}
