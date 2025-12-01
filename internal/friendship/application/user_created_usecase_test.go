package application_test

import (
	"gochat/internal/friendship/application"
	"gochat/internal/friendship/domain"
	"gochat/internal/friendship/domain/mocks"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const (
	fixedUserID     kernel.UserID     = "mock-user-123"
	fixedUserNumber kernel.UserNumber = "123456"
)

func TestUserCreatedInput_Validate_Success(t *testing.T) {
	input := &application.UserCreatedInput{
		UserID:     fixedUserID,
		UserNumber: fixedUserNumber,
	}

	err := input.Validate()
	require.NoError(t, err)
}

func TestUserCreatedInput_Validate_EmptyInput(t *testing.T) {
	inputWithEmptyUserID := &application.UserCreatedInput{
		UserID:     "",
		UserNumber: fixedUserNumber,
	}

	err := inputWithEmptyUserID.Validate()
	require.ErrorIs(t, err, myErrors.ErrEmptyInput)

	inputWithEmptyUserNumber := &application.UserCreatedInput{
		UserID:     fixedUserID,
		UserNumber: "",
	}

	err = inputWithEmptyUserNumber.Validate()
	require.ErrorIs(t, err, myErrors.ErrEmptyInput)
}

func TestUserCreatedInput_Validate_InvalidNumber(t *testing.T) {
	inputWithInvalidNumber := &application.UserCreatedInput{
		UserID:     fixedUserID,
		UserNumber: "wrong-number",
	}

	err := inputWithInvalidNumber.Validate()
	require.ErrorIs(t, err, myErrors.ErrInvalidNumber)
}

func TestNewUserCreatedUseCase_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	creator := mocks.NewMockUserCreator(ctrl)

	useCase, err := application.NewUserCreatedUseCase(
		creator,
	)

	require.NoError(t, err)
	require.NotNil(t, useCase)
}

func TestNewUserCreatedUseCase_EmptyPointer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	useCase, err := application.NewUserCreatedUseCase(
		nil,
	)

	require.ErrorIs(t, err, myErrors.ErrEmptyPointer)
	require.Nil(t, useCase)
}

func TestUserCreatedUseCase_Execute_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	creator := mocks.NewMockUserCreator(ctrl)

	useCase, err := application.NewUserCreatedUseCase(
		creator,
	)
	require.NoError(t, err)
	require.NotNil(t, useCase)

	user := domain.CreateUser(fixedUserID, fixedUserNumber)

	creator.EXPECT().Create(gomock.Any(), user).Return(nil).Times(1)

	_, err = useCase.Execute(nil, &application.UserCreatedInput{
		UserID:     fixedUserID,
		UserNumber: fixedUserNumber,
	})
	require.NoError(t, err)
}
