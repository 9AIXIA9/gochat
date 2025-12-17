package application_test

import (
	"gochat/internal/friendship/application"
	"gochat/internal/friendship/domain"
	"gochat/internal/friendship/domain/mocks"
	myErrors "gochat/internal/shared/errors"
	eventMock "gochat/internal/shared/event/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestAgreeFriendRequestInput_Validate(t *testing.T) {
	input := &application.AgreeFriendRequestInput{
		UserID:    fixedUserID,
		RequestID: fixedOperationID,
	}

	err := input.Validate()
	require.NoError(t, err)

	inputWithEmptyUserID := &application.AgreeFriendRequestInput{
		UserID:    "",
		RequestID: fixedOperationID,
	}

	err = inputWithEmptyUserID.Validate()
	require.ErrorIs(t, err, myErrors.ErrEmptyInput)

	inputWithEmptyRequestID := &application.AgreeFriendRequestInput{
		UserID:    fixedUserID,
		RequestID: "",
	}

	err = inputWithEmptyRequestID.Validate()
	require.ErrorIs(t, err, myErrors.ErrEmptyInput)
}

func TestNewAgreeFriendRequestUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFriendRequestFinderByRequestID := mocks.NewMockFriendRequestFinderByID(ctrl)
	mockFriendRequestUpdater := mocks.NewMockFriendRequestUpdater(ctrl)
	mockIDGenerator := eventMock.NewMockIDGenerator(ctrl)

	useCase, err := application.NewAgreeFriendRequestUseCase(
		mockFriendRequestFinderByRequestID, mockFriendRequestUpdater, mockIDGenerator,
	)

	require.NoError(t, err)
	require.NotNil(t, useCase)

	useCaseWithNil, err := application.NewAgreeFriendRequestUseCase(
		nil, nil, nil,
	)

	require.ErrorIs(t, err, myErrors.ErrEmptyPointer)
	require.Nil(t, useCaseWithNil)
}

func TestAgreeFriendRequestUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFriendRequestFinderByRequestID := mocks.NewMockFriendRequestFinderByID(ctrl)
	mockFriendRequestUpdater := mocks.NewMockFriendRequestUpdater(ctrl)
	mockIDGenerator := eventMock.NewMockIDGenerator(ctrl)

	useCase, err := application.NewAgreeFriendRequestUseCase(
		mockFriendRequestFinderByRequestID,
		mockFriendRequestUpdater,
		mockIDGenerator,
	)
	require.NoError(t, err)
	require.NotNil(t, useCase)

	// 正常情况
	mockRequest := domain.LoadFriendRequest(
		fixedOperationID,
		fixedToID,
		fixedUserID,
		fixedContent,
		domain.StatePending,
		time.Now().UTC(),
	)

	gomock.InOrder(
		mockFriendRequestFinderByRequestID.EXPECT().FindByID(nil, fixedOperationID).Return(mockRequest, nil),
		mockIDGenerator.EXPECT().Generate().Return(fixedEventID),
		mockFriendRequestUpdater.EXPECT().Update(nil, gomock.Any()).Return(nil),
	)

	_, err = useCase.Execute(nil, &application.AgreeFriendRequestInput{
		UserID:    fixedUserID,
		RequestID: fixedOperationID,
	})
	require.NoError(t, err)
	assert.Equal(t, domain.StateAgreed, mockRequest.State())

	// 请求不属于用户
	mockRequestWithAnotherUser := domain.LoadFriendRequest(
		fixedOperationID,
		fixedToID,
		"another-user-id",
		fixedContent,
		domain.StatePending,
		time.Now().UTC(),
	)

	mockFriendRequestFinderByRequestID.EXPECT().FindByID(nil, fixedOperationID).Return(mockRequestWithAnotherUser, nil)

	_, err = useCase.Execute(nil, &application.AgreeFriendRequestInput{
		UserID:    fixedUserID,
		RequestID: fixedOperationID,
	})
	require.ErrorIs(t, err, domain.ErrFriendRequestNotForUser)

	// 请求已被处理
	mockRequestWithHandledState := domain.LoadFriendRequest(
		fixedOperationID,
		fixedToID,
		fixedUserID,
		fixedContent,
		domain.StateAgreed,
		time.Now().UTC(),
	)

	mockFriendRequestFinderByRequestID.EXPECT().FindByID(nil, fixedOperationID).Return(mockRequestWithHandledState, nil)

	_, err = useCase.Execute(nil, &application.AgreeFriendRequestInput{
		UserID:    fixedUserID,
		RequestID: fixedOperationID,
	})
	require.ErrorIs(t, err, domain.ErrFriendRequestHasBeenHandled)

	// 未查询到此请求
	mockFriendRequestFinderByRequestID.EXPECT().FindByID(nil, fixedOperationID).Return(nil, myErrors.ErrNotFound).Times(1)
	_, err = useCase.Execute(nil, &application.AgreeFriendRequestInput{
		UserID:    fixedUserID,
		RequestID: fixedOperationID,
	})
	require.ErrorContains(t, err, "request doesn't exist")
}
