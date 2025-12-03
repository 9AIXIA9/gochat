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

func TestFriendRequestCreatedNotificationRequestedInput_Validate(t *testing.T) {
	input := &application.FriendRequestCreatedNotificationRequestedInput{
		RequestID: fixedOperationID,
		From:      fixedFriendID,
		To:        fixedUserID,
		SentAt:    time.Now().UTC(),
		Content:   fixedContent,
	}

	err := input.Validate()
	require.NoError(t, err)

	inputWithEmptyRequestID := &application.FriendRequestCreatedNotificationRequestedInput{
		RequestID: "",
		From:      fixedFriendID,
		To:        fixedUserID,
		SentAt:    time.Now().UTC(),
		Content:   fixedContent,
	}

	err = inputWithEmptyRequestID.Validate()
	require.ErrorIs(t, err, myErrors.ErrEmptyInput)

	inputWithEmptyFrom := &application.FriendRequestCreatedNotificationRequestedInput{
		RequestID: fixedOperationID,
		From:      "",
		To:        fixedUserID,
		SentAt:    time.Now().UTC(),
		Content:   fixedContent,
	}

	err = inputWithEmptyFrom.Validate()
	require.ErrorIs(t, err, myErrors.ErrEmptyInput)

	err = inputWithEmptyRequestID.Validate()
	require.ErrorIs(t, err, myErrors.ErrEmptyInput)

	inputWithEmptyTo := &application.FriendRequestCreatedNotificationRequestedInput{
		RequestID: fixedOperationID,
		To:        "",
		From:      fixedFriendID,
		SentAt:    time.Now().UTC(),
		Content:   fixedContent,
	}

	err = inputWithEmptyTo.Validate()
	require.ErrorIs(t, err, myErrors.ErrEmptyInput)
}

func TestNewFriendRequestCreatedNotificationRequestedUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIDGenerator := kernelmocks.NewMockMessageIDGenerator(ctrl)
	mockCreator := mocks.NewMockSystemMessageCreator(ctrl)
	mockNotifier := mocks.NewMockSystemMessageNotifier(ctrl)

	useCase, err := application.NewFriendRequestCreatedNotificationRequestedUseCase(
		mockIDGenerator, mockCreator, mockNotifier,
	)

	require.NoError(t, err)
	require.NotNil(t, useCase)

	useCaseWithNil, err := application.NewFriendRequestCreatedNotificationRequestedUseCase(
		nil, nil, nil,
	)

	require.ErrorIs(t, err, myErrors.ErrEmptyPointer)
	require.Nil(t, useCaseWithNil)
}

func TestFriendRequestCreatedNotificationRequestedUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIDGenerator := kernelmocks.NewMockMessageIDGenerator(ctrl)
	mockCreator := mocks.NewMockSystemMessageCreator(ctrl)
	mockNotifier := mocks.NewMockSystemMessageNotifier(ctrl)

	useCase, err := application.NewFriendRequestCreatedNotificationRequestedUseCase(
		mockIDGenerator, mockCreator, mockNotifier,
	)

	require.NoError(t, err)
	require.NotNil(t, useCase)

	//正常情况
	input := &application.FriendRequestCreatedNotificationRequestedInput{
		RequestID: fixedOperationID,
		From:      fixedFriendID,
		To:        fixedUserID,
		SentAt:    time.Now().UTC(),
		Content:   fixedContent,
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
