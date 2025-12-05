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

func TestFriendshipCreatedInput_Validate(t *testing.T) {
	input := &application.FriendshipCreatedInput{
		FriendshipID: fixedFriendshipID,
	}

	err := input.Validate()
	require.NoError(t, err)

	inputWithEmptyFriendshipID := &application.FriendshipCreatedInput{
		FriendshipID: "",
	}

	err = inputWithEmptyFriendshipID.Validate()
	require.ErrorIs(t, err, myErrors.ErrEmptyInput)
}

func TestNewFriendshipCreatedUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIDGenerator := eventMock.NewMockIDGenerator(ctrl)
	mockCreator := eventMock.NewMockUnpublishedEventsCreator(ctrl)
	mockFinder := mocks.NewMockFriendshipFinderByID(ctrl)

	useCase, err := application.NewFriendshipCreatedUseCase(
		mockFinder, mockCreator, mockIDGenerator,
	)

	require.NoError(t, err)
	require.NotNil(t, useCase)

	useCaseWithNil, err := application.NewFriendshipCreatedUseCase(
		nil, nil, nil,
	)

	require.ErrorIs(t, err, myErrors.ErrEmptyPointer)
	require.Nil(t, useCaseWithNil)
}

func TestFriendshipCreatedUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIDGenerator := eventMock.NewMockIDGenerator(ctrl)
	mockCreator := eventMock.NewMockUnpublishedEventsCreator(ctrl)
	mockFinder := mocks.NewMockFriendshipFinderByID(ctrl)

	useCase, err := application.NewFriendshipCreatedUseCase(
		mockFinder, mockCreator, mockIDGenerator,
	)

	require.NoError(t, err)
	require.NotNil(t, useCase)

	// 正常情况
	mockFriendship := domain.LoadFriendship(
		fixedFriendshipID,
		fixedUserID,
		fixedToID,
		time.Now().UTC(),
	)

	gomock.InOrder(
		mockFinder.EXPECT().FindByID(nil, fixedFriendshipID).Return(mockFriendship, nil),
		mockIDGenerator.EXPECT().Generate().Return(fixedEventID),
		mockIDGenerator.EXPECT().Generate().Return(fixedEventID),
		mockIDGenerator.EXPECT().Generate().Return(fixedEventID),
		mockCreator.EXPECT().CreateUnpublishedEvents(nil, gomock.Any()).Return(nil),
	)

	_, err = useCase.Execute(nil, &application.FriendshipCreatedInput{
		FriendshipID: fixedFriendshipID,
	})
	require.NoError(t, err)
}
