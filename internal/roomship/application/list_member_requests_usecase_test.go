package application_test

import (
	"gochat/internal/roomship/application"
	"gochat/internal/roomship/domain"
	"gochat/internal/roomship/domain/mocks"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const (
	fixedMemberRequestBaseID   kernel.OperationID = "request-123"
	maxMemberRequestsLimit                        = 100
	defaultMemberRequestsLimit                    = 10
)

func TestListMemberRequestsInput_Validate(t *testing.T) {
	input := &application.ListMemberRequestsInput{
		UserID: fixedUserID,
		BaseID: fixedMemberRequestBaseID,
		Limit:  maxMemberRequestsLimit,
	}

	err := input.Validate()
	require.NoError(t, err)

	inputWithNoBaseID := &application.ListMemberRequestsInput{
		UserID: fixedUserID,
		BaseID: "",
		Limit:  maxMemberRequestsLimit,
	}

	err = inputWithNoBaseID.Validate()
	require.NoError(t, err)

	inputWithNoLimit := &application.ListMemberRequestsInput{
		UserID: fixedUserID,
		BaseID: fixedMemberRequestBaseID,
		Limit:  0,
	}

	err = inputWithNoLimit.Validate()
	require.NoError(t, err)

	inputWithNegativeLimit := &application.ListMemberRequestsInput{
		UserID: fixedUserID,
		BaseID: fixedMemberRequestBaseID,
		Limit:  -10,
	}

	err = inputWithNegativeLimit.Validate()
	require.NoError(t, err)

	inputWithNoBaseIDAndLimit := &application.ListMemberRequestsInput{
		UserID: fixedUserID,
		BaseID: "",
		Limit:  0,
	}

	err = inputWithNoBaseIDAndLimit.Validate()
	require.NoError(t, err)

	inputWithEmptyUserID := &application.ListMemberRequestsInput{
		UserID: "",
		BaseID: fixedMemberRequestBaseID,
		Limit:  maxMemberRequestsLimit,
	}

	err = inputWithEmptyUserID.Validate()
	require.ErrorIs(t, err, myErrors.ErrEmptyInput)
}

func TestNewListMemberRequestsUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	finder := mocks.NewMockMemberRequestsFinderByUserID(ctrl)

	useCase, err := application.NewListMemberRequestsUseCase(
		finder,
	)

	require.NoError(t, err)
	require.NotNil(t, useCase)

	useCaseWithNil, err := application.NewListMemberRequestsUseCase(
		nil,
	)

	require.ErrorIs(t, err, myErrors.ErrEmptyPointer)
	require.Nil(t, useCaseWithNil)
}

func TestListMemberRequestsUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	finder := mocks.NewMockMemberRequestsFinderByUserID(ctrl)

	useCase, err := application.NewListMemberRequestsUseCase(
		finder,
	)
	require.NoError(t, err)
	require.NotNil(t, useCase)

	// 模拟正常情况，limit在范围内
	fixedLimit := maxMemberRequestsLimit - 1
	requests := getMemberRequests(fixedLimit)

	finder.EXPECT().FindsByUserID(gomock.Any(), fixedUserID, fixedLimit, fixedMemberRequestBaseID).Return(requests, nil)

	output, err := useCase.Execute(nil, &application.ListMemberRequestsInput{
		UserID: fixedUserID,
		BaseID: fixedMemberRequestBaseID,
		Limit:  maxMemberRequestsLimit - 1,
	})
	require.NoError(t, err)
	require.NotNil(t, output)

	require.Len(t, output.Requests, fixedLimit)

	// 模拟limit为0
	requests = getMemberRequests(defaultMemberRequestsLimit)

	finder.EXPECT().FindsByUserID(gomock.Any(), fixedUserID, defaultMemberRequestsLimit, fixedMemberRequestBaseID).Return(requests, nil)

	output, err = useCase.Execute(nil, &application.ListMemberRequestsInput{
		UserID: fixedUserID,
		BaseID: fixedMemberRequestBaseID,
		Limit:  0,
	})
	require.NoError(t, err)
	require.NotNil(t, output)

	require.Len(t, output.Requests, defaultMemberRequestsLimit)

	//模拟负数limit
	finder.EXPECT().FindsByUserID(gomock.Any(), fixedUserID, defaultMemberRequestsLimit, fixedMemberRequestBaseID).Return(requests, nil)

	output, err = useCase.Execute(nil, &application.ListMemberRequestsInput{
		UserID: fixedUserID,
		BaseID: fixedMemberRequestBaseID,
		Limit:  -100,
	})
	require.NoError(t, err)
	require.NotNil(t, output)

	require.Len(t, output.Requests, defaultMemberRequestsLimit)

	// 模拟太大的limit
	requests = getMemberRequests(maxMemberRequestsLimit)

	finder.EXPECT().FindsByUserID(gomock.Any(), fixedUserID, maxMemberRequestsLimit, fixedMemberRequestBaseID).Return(requests, nil)

	output, err = useCase.Execute(nil, &application.ListMemberRequestsInput{
		UserID: fixedUserID,
		BaseID: fixedMemberRequestBaseID,
		Limit:  1000,
	})
	require.NoError(t, err)
	require.NotNil(t, output)

	require.Len(t, output.Requests, maxMemberRequestsLimit)
}

func getMemberRequests(count int) []*domain.MemberRequest {
	requests := make([]*domain.MemberRequest, count)
	for i := 0; i < count; i++ {
		requests[i] = domain.LoadMemberRequest(
			fixedOperationID,
			domain.StatePending,
			fixedUserID,
			fixedRoomID,
			fixedContent,
			fixedOperatorID,
			time.Now().UTC(),
			time.Now().UTC(),
		)
	}
	return requests
}
