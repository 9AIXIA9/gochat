package application_test

import (
	"fmt"
	"gochat/internal/friendship/application"
	"gochat/internal/friendship/domain"
	"gochat/internal/friendship/domain/mocks"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const (
	fixedBaseID kernel.OperationID = "base-111"
	fixedLimit                     = 3
)

func TestListFriendRequestsInput_Validate_Success(t *testing.T) {
	input := &application.ListFriendRequestsInput{
		UserID: fixedFromID,
		BaseID: fixedBaseID,
		Limit:  fixedLimit,
	}

	err := input.Validate()
	require.NoError(t, err)

	inputWithNoBaseID := &application.ListFriendRequestsInput{
		UserID: fixedFromID,
		BaseID: "",
		Limit:  fixedLimit,
	}

	err = inputWithNoBaseID.Validate()
	require.NoError(t, err)

	inputWithNoLimit := &application.ListFriendRequestsInput{
		UserID: fixedFromID,
		BaseID: fixedBaseID,
		Limit:  0,
	}

	err = inputWithNoLimit.Validate()
	require.NoError(t, err)

	inputWithNoBaseIDAndLimit := &application.ListFriendRequestsInput{
		UserID: fixedFromID,
		BaseID: "",
		Limit:  0,
	}

	err = inputWithNoBaseIDAndLimit.Validate()
	require.NoError(t, err)
}

func TestListFriendRequestsInput_Validate_EmptyInput(t *testing.T) {
	inputWithEmptyUserID := &application.ListFriendRequestsInput{
		UserID: "",
		BaseID: fixedBaseID,
		Limit:  fixedLimit,
	}

	err := inputWithEmptyUserID.Validate()
	require.ErrorIs(t, err, myErrors.ErrEmptyInput)
}

func TestNewListFriendRequestsUseCase_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	finder := mocks.NewMockFriendRequestsFinderByUserID(ctrl)

	useCase, err := application.NewListFriendRequestsUseCase(
		finder,
	)

	require.NoError(t, err)
	require.NotNil(t, useCase)
}

func TestNewListFriendRequestsUseCase_EmptyPointer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	useCase, err := application.NewListFriendRequestsUseCase(
		nil,
	)

	require.ErrorIs(t, err, myErrors.ErrEmptyPointer)
	require.Nil(t, useCase)
}

func TestListFriendRequestsUseCase_Execute_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	finder := mocks.NewMockFriendRequestsFinderByUserID(ctrl)

	useCase, err := application.NewListFriendRequestsUseCase(
		finder,
	)
	require.NoError(t, err)
	require.NotNil(t, useCase)

	requests := make([]*domain.FriendRequest, fixedLimit)
	for i := 0; i < fixedLimit; i++ {
		requests[i] = domain.LoadFriendRequest(
			kernel.OperationID(fmt.Sprintf("mock-request-id-%d", i+1)),
			fixedFromID,
			fixedToID,
			fmt.Sprintf("Hello %d", i+1),
			domain.StatePending,
			time.Now().UTC(),
		)
	}

	finder.EXPECT().FindsByUserID(gomock.Any(), fixedFromID, fixedLimit, fixedBaseID).Return(requests, nil)

	output, err := useCase.Execute(nil, &application.ListFriendRequestsInput{
		UserID: fixedFromID,
		BaseID: fixedBaseID,
		Limit:  fixedLimit,
	})
	require.NoError(t, err)
	require.NotNil(t, output)

	require.Len(t, output.Requests, 3)
	for i, request := range requests {
		assert.Equal(t, kernel.OperationID(fmt.Sprintf("mock-request-id-%d", i+1)), request.ID())
	}

	requests = make([]*domain.FriendRequest, 10)
	for i := 0; i < 10; i++ {
		requests[i] = domain.LoadFriendRequest(
			kernel.OperationID(fmt.Sprintf("mock-request-id-%d", i+1)),
			fixedFromID,
			fixedToID,
			fmt.Sprintf("Hello %d", i+1),
			domain.StatePending,
			time.Now().UTC(),
		)
	}

	finder.EXPECT().FindsByUserID(gomock.Any(), fixedFromID, 10, fixedBaseID).Return(requests, nil)

	output, err = useCase.Execute(nil, &application.ListFriendRequestsInput{
		UserID: fixedFromID,
		BaseID: fixedBaseID,
		Limit:  0,
	})
	require.NoError(t, err)
	require.NotNil(t, output)

	require.Len(t, output.Requests, 10)
	for i, request := range requests {
		assert.Equal(t, kernel.OperationID(fmt.Sprintf("mock-request-id-%d", i+1)), request.ID())
	}

	// 模拟太大的limit
	requests = make([]*domain.FriendRequest, 100)
	for i := 0; i < 100; i++ {
		requests[i] = domain.LoadFriendRequest(
			kernel.OperationID(fmt.Sprintf("mock-request-id-%d", i+1)),
			fixedFromID,
			fixedToID,
			fmt.Sprintf("Hello %d", i+1),
			domain.StatePending,
			time.Now().UTC(),
		)
	}

	finder.EXPECT().FindsByUserID(gomock.Any(), fixedFromID, 100, fixedBaseID).Return(requests, nil)

	output, err = useCase.Execute(nil, &application.ListFriendRequestsInput{
		UserID: fixedFromID,
		BaseID: fixedBaseID,
		Limit:  1000,
	})
	require.NoError(t, err)
	require.NotNil(t, output)

	require.Len(t, output.Requests, 100)
	for i, request := range requests {
		assert.Equal(t, kernel.OperationID(fmt.Sprintf("mock-request-id-%d", i+1)), request.ID())
	}
}
