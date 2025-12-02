package application_test

import (
	"gochat/internal/friendship/application"
	"gochat/internal/friendship/domain"
	"gochat/internal/friendship/domain/mocks"
	myErrors "gochat/internal/shared/errors"
	eventMocks "gochat/internal/shared/event/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestFriendRequestAgreedInput_Validate(t *testing.T) {
	input := &application.FriendRequestAgreedInput{
		RequestID: fixedOperationID,
	}

	err := input.Validate()
	require.NoError(t, err)

	inputWithEmptyRequestID := &application.FriendRequestAgreedInput{
		RequestID: "",
	}

	err = inputWithEmptyRequestID.Validate()
	require.ErrorIs(t, err, myErrors.ErrEmptyInput)
}

func TestNewFriendRequestAgreedUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	friendRequestFinderByID := mocks.NewMockFriendRequestFinderByID(ctrl)
	creator := mocks.NewMockFriendshipCreator(ctrl)
	friendshipIDGenerator := mocks.NewMockFriendshipIDGenerator(ctrl)
	idGenerator := eventMocks.NewMockIDGenerator(ctrl)

	useCase, err := application.NewFriendRequestAgreedUseCase(
		friendRequestFinderByID, creator, friendshipIDGenerator, idGenerator,
	)

	require.NoError(t, err)
	require.NotNil(t, useCase)

	useCaseWithNil, err := application.NewFriendRequestAgreedUseCase(
		nil, nil, nil, nil,
	)

	require.ErrorIs(t, err, myErrors.ErrEmptyPointer)
	require.Nil(t, useCaseWithNil)
}

func TestFriendRequestAgreedUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	friendRequestFinderByID := mocks.NewMockFriendRequestFinderByID(ctrl)
	creator := mocks.NewMockFriendshipCreator(ctrl)
	friendshipIDGenerator := mocks.NewMockFriendshipIDGenerator(ctrl)
	idGenerator := eventMocks.NewMockIDGenerator(ctrl)

	useCase, err := application.NewFriendRequestAgreedUseCase(
		friendRequestFinderByID, creator, friendshipIDGenerator, idGenerator,
	)

	require.NoError(t, err)
	require.NotNil(t, useCase)

	mockRequest := domain.LoadFriendRequest(
		fixedOperationID,
		fixedUserID,
		fixedToID,
		fixedContent,
		domain.StateAgreed,
		time.Now(),
	)

	//正常情况
	gomock.InOrder(
		friendRequestFinderByID.EXPECT().FindByID(nil, fixedOperationID).Return(mockRequest, nil),
		friendshipIDGenerator.EXPECT().Generate().Return(fixedFriendID),
		idGenerator.EXPECT().Generate().Return(fixedEventID),
		creator.EXPECT().Create(nil, gomock.Any()).Return(nil),
	)

	_, err = useCase.Execute(nil, &application.FriendRequestAgreedInput{
		RequestID: fixedOperationID,
	})
	require.NoError(t, err)

	//请求状态不是已同意
	mockRequestNotAgreed := domain.LoadFriendRequest(
		fixedOperationID,
		fixedUserID,
		fixedToID,
		fixedContent,
		domain.StatePending,
		time.Now(),
	)

	friendRequestFinderByID.EXPECT().FindByID(nil, fixedOperationID).Return(mockRequestNotAgreed, nil)

	_, err = useCase.Execute(nil, &application.FriendRequestAgreedInput{
		RequestID: fixedOperationID,
	})
	require.ErrorIs(t, err, domain.ErrFriendRequestNotAgreed)
}
