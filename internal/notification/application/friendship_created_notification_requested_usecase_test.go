package application_test

import (
	"errors"
	"gochat/internal/notification/application"
	"gochat/internal/notification/domain/mocks"
	myErrors "gochat/internal/shared/errors"
	kernelmocks "gochat/internal/shared/kernel/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestFriendshipCreatedNotificationRequestedInput_Validate(t *testing.T) {
	input := &application.FriendshipCreatedNotificationRequestedInput{
		UserID:    fixedUserID,
		FriendID:  fixedFriendID,
		CreatedAt: time.Now().UTC(),
	}

	err := input.Validate()
	require.NoError(t, err)

	inputWithEmptyUserID := &application.FriendshipCreatedNotificationRequestedInput{
		UserID:    "",
		FriendID:  fixedFriendID,
		CreatedAt: time.Now(),
	}

	err = inputWithEmptyUserID.Validate()
	require.ErrorIs(t, err, myErrors.ErrEmptyInput)

	inputWithEmptyFriendID := &application.FriendshipCreatedNotificationRequestedInput{
		UserID:    fixedUserID,
		FriendID:  "",
		CreatedAt: time.Now().UTC(),
	}

	err = inputWithEmptyFriendID.Validate()
	require.ErrorIs(t, err, myErrors.ErrEmptyInput)
}

func TestNewFriendshipCreatedNotificationRequestedUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIDGenerator := kernelmocks.NewMockMessageIDGenerator(ctrl)
	mockCreator := mocks.NewMockSystemMessageCreator(ctrl)
	mockNotifier := mocks.NewMockSystemMessageNotifier(ctrl)

	useCase, err := application.NewFriendshipCreatedNotificationRequestedUseCase(
		mockIDGenerator, mockCreator, mockNotifier,
	)

	require.NoError(t, err)
	require.NotNil(t, useCase)

	useCaseWithNil, err := application.NewFriendshipCreatedNotificationRequestedUseCase(
		nil, nil, nil,
	)

	require.ErrorIs(t, err, myErrors.ErrEmptyPointer)
	require.Nil(t, useCaseWithNil)
}

func TestFriendshipCreatedNotificationRequestedUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIDGenerator := kernelmocks.NewMockMessageIDGenerator(ctrl)
	mockCreator := mocks.NewMockSystemMessageCreator(ctrl)
	mockNotifier := mocks.NewMockSystemMessageNotifier(ctrl)

	useCase, err := application.NewFriendshipCreatedNotificationRequestedUseCase(
		mockIDGenerator, mockCreator, mockNotifier,
	)

	require.NoError(t, err)
	require.NotNil(t, useCase)

	//正常情况
	input := &application.FriendshipCreatedNotificationRequestedInput{
		UserID:    fixedUserID,
		FriendID:  fixedFriendID,
		CreatedAt: time.Now().UTC(),
	}

	gomock.InOrder(
		mockIDGenerator.EXPECT().Generate().Return(fixedMessageID).Times(1),
		mockNotifier.EXPECT().Notify(gomock.Any()).Return(nil).Times(1),
		mockCreator.EXPECT().Create(nil, gomock.Any()).Return(nil).Times(1),
	)

	_, err = useCase.Execute(nil, input)
	require.NoError(t, err)

	//通知失败
	gomock.InOrder(
		mockIDGenerator.EXPECT().Generate().Return(fixedMessageID).Times(1),
		mockNotifier.EXPECT().Notify(gomock.Any()).Return(errors.New("test-notify-failed")).Times(1),
		mockCreator.EXPECT().Create(nil, gomock.Any()).Return(nil).Times(1),
	)

	_, err = useCase.Execute(nil, input)
	require.NoError(t, err)
}
