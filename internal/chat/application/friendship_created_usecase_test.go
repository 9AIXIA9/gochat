package application_test

import (
	"gochat/internal/chat/application"
	"gochat/internal/chat/domain/mocks"
	myErrors "gochat/internal/shared/errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestFriendshipCreatedInput_Validate(t *testing.T) {
	input := &application.FriendshipCreatedInput{
		ID:      fixedFriendshipID,
		UserID1: fixedUserID,
		UserID2: fixedFriendID,
	}

	err := input.Validate()
	require.NoError(t, err)

	inputWithEmptyFriendshipID := &application.FriendshipCreatedInput{
		ID:      "",
		UserID1: fixedUserID,
		UserID2: fixedFriendID,
	}

	err = inputWithEmptyFriendshipID.Validate()
	require.ErrorIs(t, err, myErrors.ErrEmptyInput)

	inputWithEmptyUserID := &application.FriendshipCreatedInput{
		ID:      fixedFriendshipID,
		UserID1: "",
		UserID2: fixedFriendID,
	}

	err = inputWithEmptyUserID.Validate()
	require.ErrorIs(t, err, myErrors.ErrEmptyInput)

	inputWithEmptyUserID = &application.FriendshipCreatedInput{
		ID:      fixedFriendshipID,
		UserID1: fixedFriendID,
		UserID2: "",
	}

	err = inputWithEmptyUserID.Validate()
	require.ErrorIs(t, err, myErrors.ErrEmptyInput)
}

func TestNewFriendshipCreatedUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFriendshipCreator := mocks.NewMockFriendshipSaver(ctrl)

	useCase, err := application.NewFriendshipCreatedUseCase(
		mockFriendshipCreator,
	)

	require.NoError(t, err)
	require.NotNil(t, useCase)

	useCaseWithNil, err := application.NewFriendshipCreatedUseCase(
		nil,
	)

	require.ErrorIs(t, err, myErrors.ErrEmptyPointer)
	require.Nil(t, useCaseWithNil)
}

func TestFriendshipCreatedUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFriendshipCreator := mocks.NewMockFriendshipSaver(ctrl)

	useCase, err := application.NewFriendshipCreatedUseCase(
		mockFriendshipCreator,
	)
	require.NoError(t, err)
	require.NotNil(t, useCase)

	// 正常情况
	mockFriendshipCreator.EXPECT().Save(nil, gomock.Any()).Return(nil).Times(1)
	_, err = useCase.Execute(nil, &application.FriendshipCreatedInput{
		ID:      fixedFriendshipID,
		UserID1: fixedUserID,
		UserID2: fixedFriendID,
	})
	require.NoError(t, err)
}
