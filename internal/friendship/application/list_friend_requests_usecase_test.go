package application_test

import (
	"gochat/internal/friendship/application"
	"gochat/internal/friendship/domain"
	"gochat/internal/friendship/domain/mocks"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const (
	fixedFriendRequestBaseID   kernel.OperationID = "request-123"
	maxFriendRequestsLimit                        = 100
	defaultFriendRequestsLimit                    = 10
)

func TestListFriendRequestsInput_Validate(t *testing.T) {
	input := &application.ListFriendRequestsInput{
		UserID: fixedUserID,
		BaseID: fixedFriendRequestBaseID,
		Limit:  maxFriendRequestsLimit,
	}

	err := input.Validate()
	require.NoError(t, err)

	inputWithNoBaseID := &application.ListFriendRequestsInput{
		UserID: fixedUserID,
		BaseID: "",
		Limit:  maxFriendRequestsLimit,
	}

	err = inputWithNoBaseID.Validate()
	require.NoError(t, err)

	inputWithNoLimit := &application.ListFriendRequestsInput{
		UserID: fixedUserID,
		BaseID: fixedFriendRequestBaseID,
		Limit:  0,
	}

	err = inputWithNoLimit.Validate()
	require.NoError(t, err)

	inputWithNegativeLimit := &application.ListFriendRequestsInput{
		UserID: fixedUserID,
		BaseID: fixedFriendRequestBaseID,
		Limit:  -10,
	}

	err = inputWithNegativeLimit.Validate()
	require.NoError(t, err)

	inputWithNoBaseIDAndLimit := &application.ListFriendRequestsInput{
		UserID: fixedUserID,
		BaseID: "",
		Limit:  0,
	}

	err = inputWithNoBaseIDAndLimit.Validate()
	require.NoError(t, err)

	inputWithEmptyUserID := &application.ListFriendRequestsInput{
		UserID: "",
		BaseID: fixedFriendRequestBaseID,
		Limit:  maxFriendRequestsLimit,
	}

	err = inputWithEmptyUserID.Validate()
	require.ErrorIs(t, err, myErrors.ErrEmptyInput)
}

func TestNewListFriendRequestsUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	finder := mocks.NewMockFriendRequestsFinderByUserID(ctrl)

	useCase, err := application.NewListFriendRequestsUseCase(
		finder,
	)

	require.NoError(t, err)
	require.NotNil(t, useCase)

	useCaseWithNil, err := application.NewListFriendRequestsUseCase(
		nil,
	)

	require.ErrorIs(t, err, myErrors.ErrEmptyPointer)
	require.Nil(t, useCaseWithNil)
}

func TestListFriendRequestsUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	finder := mocks.NewMockFriendRequestsFinderByUserID(ctrl)

	useCase, err := application.NewListFriendRequestsUseCase(
		finder,
	)
	require.NoError(t, err)
	require.NotNil(t, useCase)

	// 模拟正常情况，limit在范围内
	fixedLimit := maxFriendRequestsLimit - 1
	requests := getFriendRequests(fixedLimit)

	finder.EXPECT().FindsByUserID(gomock.Any(), fixedUserID, fixedLimit, fixedFriendRequestBaseID).Return(requests, nil)

	output, err := useCase.Execute(nil, &application.ListFriendRequestsInput{
		UserID: fixedUserID,
		BaseID: fixedFriendRequestBaseID,
		Limit:  maxFriendRequestsLimit - 1,
	})
	require.NoError(t, err)
	require.NotNil(t, output)

	require.Len(t, output.Requests, fixedLimit)

	// 模拟limit为0
	requests = getFriendRequests(defaultFriendRequestsLimit)

	finder.EXPECT().FindsByUserID(gomock.Any(), fixedUserID, defaultFriendRequestsLimit, fixedFriendRequestBaseID).Return(requests, nil)

	output, err = useCase.Execute(nil, &application.ListFriendRequestsInput{
		UserID: fixedUserID,
		BaseID: fixedFriendRequestBaseID,
		Limit:  0,
	})
	require.NoError(t, err)
	require.NotNil(t, output)

	require.Len(t, output.Requests, defaultFriendRequestsLimit)

	//模拟负数limit
	finder.EXPECT().FindsByUserID(gomock.Any(), fixedUserID, defaultFriendRequestsLimit, fixedFriendRequestBaseID).Return(requests, nil)

	output, err = useCase.Execute(nil, &application.ListFriendRequestsInput{
		UserID: fixedUserID,
		BaseID: fixedFriendRequestBaseID,
		Limit:  -100,
	})
	require.NoError(t, err)
	require.NotNil(t, output)

	require.Len(t, output.Requests, defaultFriendRequestsLimit)

	// 模拟太大的limit
	requests = getFriendRequests(maxFriendRequestsLimit)

	finder.EXPECT().FindsByUserID(gomock.Any(), fixedUserID, maxFriendRequestsLimit, fixedFriendRequestBaseID).Return(requests, nil)

	output, err = useCase.Execute(nil, &application.ListFriendRequestsInput{
		UserID: fixedUserID,
		BaseID: fixedFriendRequestBaseID,
		Limit:  1000,
	})
	require.NoError(t, err)
	require.NotNil(t, output)

	require.Len(t, output.Requests, maxFriendRequestsLimit)
}

func getFriendRequests(count int) []*domain.FriendRequest {
	requests := make([]*domain.FriendRequest, count)
	for i := 0; i < count; i++ {
		requests[i] = domain.LoadFriendRequest(
			fixedOperationID,
			fixedUserID,
			fixedToID,
			fixedContent,
			domain.StatePending,
			time.Now().UTC(),
		)
	}
	return requests
}
