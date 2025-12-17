package application_test

import (
	"gochat/internal/friendship/application"
	"gochat/internal/friendship/domain"
	"gochat/internal/friendship/domain/mocks"
	myErrors "gochat/internal/shared/errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestRefuseFriendRequestInput_Validate(t *testing.T) {
	input := &application.RefuseFriendRequestInput{
		UserID:    fixedUserID,
		RequestID: fixedOperationID,
	}

	err := input.Validate()
	require.NoError(t, err)

	inputWithEmptyUserID := &application.RefuseFriendRequestInput{
		UserID:    "",
		RequestID: fixedOperationID,
	}

	err = inputWithEmptyUserID.Validate()
	require.ErrorIs(t, err, myErrors.ErrEmptyInput)

	inputWithEmptyRequestID := &application.RefuseFriendRequestInput{
		UserID:    fixedUserID,
		RequestID: "",
	}

	err = inputWithEmptyRequestID.Validate()
	require.ErrorIs(t, err, myErrors.ErrEmptyInput)
}

func TestNewRefuseFriendRequestUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFriendRequestFinderByRequestID := mocks.NewMockFriendRequestFinderByID(ctrl)
	mockFriendRequestUpdater := mocks.NewMockFriendRequestUpdater(ctrl)

	useCase, err := application.NewRefuseFriendRequestUseCase(
		mockFriendRequestFinderByRequestID, mockFriendRequestUpdater,
	)

	require.NoError(t, err)
	require.NotNil(t, useCase)

	useCaseWithNil, err := application.NewRefuseFriendRequestUseCase(
		nil, nil,
	)

	require.ErrorIs(t, err, myErrors.ErrEmptyPointer)
	require.Nil(t, useCaseWithNil)
}

func TestRefuseFriendRequestUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFriendRequestFinderByRequestID := mocks.NewMockFriendRequestFinderByID(ctrl)
	mockFriendRequestUpdater := mocks.NewMockFriendRequestUpdater(ctrl)

	useCase, err := application.NewRefuseFriendRequestUseCase(
		mockFriendRequestFinderByRequestID,
		mockFriendRequestUpdater,
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
		mockFriendRequestUpdater.EXPECT().Update(nil, gomock.Any()).Return(nil),
	)

	_, err = useCase.Execute(nil, &application.RefuseFriendRequestInput{
		UserID:    fixedUserID,
		RequestID: fixedOperationID,
	})
	require.NoError(t, err)
	assert.Equal(t, domain.StateRefused, mockRequest.State())

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

	_, err = useCase.Execute(nil, &application.RefuseFriendRequestInput{
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
		domain.StateRefused,
		time.Now().UTC(),
	)

	mockFriendRequestFinderByRequestID.EXPECT().FindByID(nil, fixedOperationID).Return(mockRequestWithHandledState, nil)

	_, err = useCase.Execute(nil, &application.RefuseFriendRequestInput{
		UserID:    fixedUserID,
		RequestID: fixedOperationID,
	})
	require.ErrorIs(t, err, domain.ErrFriendRequestHasBeenHandled)

	// 未查询到此请求
	mockFriendRequestFinderByRequestID.EXPECT().FindByID(nil, fixedOperationID).Return(nil, myErrors.ErrNotFound).Times(1)
	_, err = useCase.Execute(nil, &application.RefuseFriendRequestInput{
		UserID:    fixedUserID,
		RequestID: fixedOperationID,
	})
	require.ErrorContains(t, err, "request doesn't exist")
}
