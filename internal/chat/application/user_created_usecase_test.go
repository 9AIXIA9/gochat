package application_test

import (
	"gochat/internal/chat/application"
	"gochat/internal/chat/domain/mocks"
	myErrors "gochat/internal/shared/errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestUserCreatedInput_Validate(t *testing.T) {
	input := &application.UserCreatedInput{
		UserID: fixedUserID,
	}

	err := input.Validate()
	require.NoError(t, err)

	inputWithEmptyUserID := &application.UserCreatedInput{
		UserID: "",
	}

	err = inputWithEmptyUserID.Validate()
	require.ErrorIs(t, err, myErrors.ErrEmptyInput)
}

func TestNewUserCreatedUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserSaver := mocks.NewMockUserSaver(ctrl)

	useCase, err := application.NewUserCreatedUseCase(
		mockUserSaver,
	)

	require.NoError(t, err)
	require.NotNil(t, useCase)

	useCaseWithNil, err := application.NewUserCreatedUseCase(
		nil,
	)

	require.ErrorIs(t, err, myErrors.ErrEmptyPointer)
	require.Nil(t, useCaseWithNil)
}

func TestUserCreatedUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserSaver := mocks.NewMockUserSaver(ctrl)

	useCase, err := application.NewUserCreatedUseCase(
		mockUserSaver,
	)
	require.NoError(t, err)
	require.NotNil(t, useCase)

	// 正常情况
	mockUserSaver.EXPECT().Save(nil, gomock.Any()).Return(nil).Times(1)
	_, err = useCase.Execute(nil, &application.UserCreatedInput{
		UserID: fixedUserID,
	})
	require.NoError(t, err)
}
