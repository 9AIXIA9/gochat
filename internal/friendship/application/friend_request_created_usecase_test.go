package application_test

import (
	"gochat/internal/friendship/application"
	"gochat/internal/friendship/domain"
	"gochat/internal/friendship/domain/mocks"
	myErrors "gochat/internal/shared/errors"
	eventMock "gochat/internal/shared/event/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestFriendRequestCreatedInput_Validate(t *testing.T) {
	input := &application.FriendRequestCreatedInput{
		RequestID: fixedOperationID,
	}

	err := input.Validate()
	require.NoError(t, err)

	inputWithEmptyRequestID := &application.FriendRequestCreatedInput{
		RequestID: "",
	}

	err = inputWithEmptyRequestID.Validate()
	require.ErrorIs(t, err, myErrors.ErrEmptyInput)
}

func TestNewFriendRequestCreatedUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIDGenerator := eventMock.NewMockIDGenerator(ctrl)
	mockCreator := eventMock.NewMockUnpublishedEventsCreator(ctrl)
	mockFinder := mocks.NewMockFriendRequestFinderByID(ctrl)

	useCase, err := application.NewFriendRequestCreatedUseCase(
		mockFinder, mockCreator, mockIDGenerator,
	)

	require.NoError(t, err)
	require.NotNil(t, useCase)

	useCaseWithNil, err := application.NewFriendRequestCreatedUseCase(
		nil, nil, nil,
	)

	require.ErrorIs(t, err, myErrors.ErrEmptyPointer)
	require.Nil(t, useCaseWithNil)
}

func TestFriendRequestCreatedUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIDGenerator := eventMock.NewMockIDGenerator(ctrl)
	mockCreator := eventMock.NewMockUnpublishedEventsCreator(ctrl)
	mockFinder := mocks.NewMockFriendRequestFinderByID(ctrl)

	useCase, err := application.NewFriendRequestCreatedUseCase(
		mockFinder, mockCreator, mockIDGenerator,
	)

	require.NoError(t, err)
	require.NotNil(t, useCase)

	// 正常情况
	mockFriendRequest := domain.LoadFriendRequest(
		fixedOperationID,
		fixedUserID,
		fixedToID,
		fixedContent,
		domain.StatePending,
		time.Now().UTC(),
	)

	gomock.InOrder(
		mockFinder.EXPECT().FindByID(nil, fixedOperationID).Return(mockFriendRequest, nil),
		mockIDGenerator.EXPECT().Generate().Return(fixedEventID),
		mockCreator.EXPECT().CreateUnpublishedEvents(nil, gomock.Any()).Return(nil),
	)

	_, err = useCase.Execute(nil, &application.FriendRequestCreatedInput{
		RequestID: fixedOperationID,
	})
	require.NoError(t, err)

	// 不是待处理状态
	mockFriendRequestNotPending := domain.LoadFriendRequest(
		fixedOperationID,
		fixedUserID,
		fixedToID,
		fixedContent,
		domain.StateAgreed,
		time.Now().UTC(),
	)

	gomock.InOrder(
		mockFinder.EXPECT().FindByID(nil, fixedOperationID).Return(mockFriendRequestNotPending, nil),
	)

	_, err = useCase.Execute(nil, &application.FriendRequestCreatedInput{
		RequestID: fixedOperationID,
	})
	require.NoError(t, err)

	// 添加自己为好友
	mockFriendRequestAddYourself := domain.LoadFriendRequest(
		fixedOperationID,
		fixedUserID,
		fixedUserID,
		fixedContent,
		domain.StatePending,
		time.Now().UTC(),
	)

	gomock.InOrder(
		mockFinder.EXPECT().FindByID(nil, fixedOperationID).Return(mockFriendRequestAddYourself, nil),
	)

	_, err = useCase.Execute(nil, &application.FriendRequestCreatedInput{
		RequestID: fixedOperationID,
	})
	require.ErrorIs(t, err, domain.ErrAddYourselfAsFriend)
}
