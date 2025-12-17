package application_test

import (
	"gochat/internal/profile/application"
	"gochat/internal/profile/domain/mocks"
	myErrors "gochat/internal/shared/errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestUserCreatedInput_Validate(t *testing.T) {
	input := &application.UserCreatedInput{
		UserID:     fixedUserID,
		Email:      fixedEmail,
		SignedUpAt: time.Now().UTC(),
	}

	err := input.Validate()
	require.NoError(t, err)

	inputWithEmptyUserID := &application.UserCreatedInput{
		UserID:     "",
		Email:      fixedEmail,
		SignedUpAt: time.Now().UTC(),
	}

	err = inputWithEmptyUserID.Validate()
	require.ErrorIs(t, err, myErrors.ErrEmptyInput)

	inputWithEmptyEmail := &application.UserCreatedInput{
		Email:      "",
		UserID:     fixedUserID,
		SignedUpAt: time.Now().UTC(),
	}

	err = inputWithEmptyEmail.Validate()
	require.ErrorIs(t, err, myErrors.ErrInvalidFormat)

	inputWithEmptySignedUpAt := &application.UserCreatedInput{
		SignedUpAt: time.Time{},
		Email:      fixedEmail,
		UserID:     fixedUserID,
	}

	err = inputWithEmptySignedUpAt.Validate()
	require.ErrorIs(t, err, myErrors.ErrEmptyInput)
}

func TestNewUserCreatedUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserSaver := mocks.NewMockUserSaver(ctrl)
	mockProfileCreator := mocks.NewMockUserProfileCreator(ctrl)

	useCase, err := application.NewUserCreatedUseCase(
		mockUserSaver,
		mockProfileCreator,
	)

	require.NoError(t, err)
	require.NotNil(t, useCase)

	useCaseWithNil, err := application.NewUserCreatedUseCase(
		nil, nil,
	)

	require.ErrorIs(t, err, myErrors.ErrEmptyPointer)
	require.Nil(t, useCaseWithNil)
}

func TestUserCreatedUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserSaver := mocks.NewMockUserSaver(ctrl)
	mockProfileCreator := mocks.NewMockUserProfileCreator(ctrl)

	useCase, err := application.NewUserCreatedUseCase(
		mockUserSaver,
		mockProfileCreator,
	)
	require.NoError(t, err)
	require.NotNil(t, useCase)

	// 正常情况
	gomock.InOrder(
		mockUserSaver.EXPECT().Save(nil, gomock.Any()).Return(nil).Times(1),
		mockProfileCreator.EXPECT().Create(nil, gomock.Any()).Return(nil).Times(1),
	)

	_, err = useCase.Execute(nil, &application.UserCreatedInput{
		UserID:     fixedUserID,
		Email:      fixedEmail,
		SignedUpAt: time.Now().UTC(),
	})
	require.NoError(t, err)

	// 重复创建
	gomock.InOrder(
		mockUserSaver.EXPECT().Save(nil, gomock.Any()).Return(nil).Times(1),
		mockProfileCreator.EXPECT().Create(nil, gomock.Any()).Return(myErrors.ErrDuplicatedKey).Times(1),
	)
	_, err = useCase.Execute(nil, &application.UserCreatedInput{
		UserID:     fixedUserID,
		Email:      fixedEmail,
		SignedUpAt: time.Now().UTC(),
	})
	require.NoError(t, err)
}
