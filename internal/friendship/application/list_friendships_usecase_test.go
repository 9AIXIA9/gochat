package application_test

import (
	"gochat/internal/friendship/application"
	"gochat/internal/friendship/domain"
	"gochat/internal/friendship/domain/mocks"
	myErrors "gochat/internal/shared/errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const (
	fixedFriendshipBaseID   domain.FriendshipID = "friendship-123"
	maxFriendshipsLimit                         = 100
	defaultFriendshipsLimit                     = 50
)

func TestListFriendshipsInput_Validate(t *testing.T) {
	input := &application.ListFriendshipsInput{
		UserID: fixedUserID,
		BaseID: fixedFriendshipBaseID,
		Limit:  maxFriendshipsLimit,
	}

	err := input.Validate()
	require.NoError(t, err)

	inputWithNoBaseID := &application.ListFriendshipsInput{
		UserID: fixedUserID,
		BaseID: "",
		Limit:  maxFriendshipsLimit,
	}

	err = inputWithNoBaseID.Validate()
	require.NoError(t, err)

	inputWithNoLimit := &application.ListFriendshipsInput{
		UserID: fixedUserID,
		BaseID: fixedFriendshipBaseID,
		Limit:  0,
	}

	err = inputWithNoLimit.Validate()
	require.NoError(t, err)

	inputWithNegativeLimit := &application.ListFriendshipsInput{
		UserID: fixedUserID,
		BaseID: fixedFriendshipBaseID,
		Limit:  -10,
	}

	err = inputWithNegativeLimit.Validate()
	require.NoError(t, err)

	inputWithNoBaseIDAndLimit := &application.ListFriendshipsInput{
		UserID: fixedUserID,
		BaseID: "",
		Limit:  0,
	}

	err = inputWithNoBaseIDAndLimit.Validate()
	require.NoError(t, err)

	inputWithEmptyUserID := &application.ListFriendshipsInput{
		UserID: "",
		BaseID: fixedFriendshipBaseID,
		Limit:  maxFriendshipsLimit,
	}

	err = inputWithEmptyUserID.Validate()
	require.ErrorIs(t, err, myErrors.ErrEmptyInput)
}

func TestNewListFriendshipsUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	finder := mocks.NewMockFriendshipsFinderByUserID(ctrl)

	useCase, err := application.NewListFriendshipsUseCase(
		finder,
	)

	require.NoError(t, err)
	require.NotNil(t, useCase)

	useCaseWithNil, err := application.NewListFriendshipsUseCase(
		nil,
	)

	require.ErrorIs(t, err, myErrors.ErrEmptyPointer)
	require.Nil(t, useCaseWithNil)
}

func TestListFriendshipsUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	finder := mocks.NewMockFriendshipsFinderByUserID(ctrl)

	useCase, err := application.NewListFriendshipsUseCase(
		finder,
	)
	require.NoError(t, err)
	require.NotNil(t, useCase)

	// 模拟正常情况，limit在范围内
	fixedLimit := maxFriendshipsLimit - 1
	friendships := getFriendships(fixedLimit)

	finder.EXPECT().FindsByUserID(gomock.Any(), fixedUserID, fixedLimit, fixedFriendshipBaseID).Return(friendships, nil)

	output, err := useCase.Execute(nil, &application.ListFriendshipsInput{
		UserID: fixedUserID,
		BaseID: fixedFriendshipBaseID,
		Limit:  maxFriendshipsLimit - 1,
	})
	require.NoError(t, err)
	require.NotNil(t, output)

	require.Len(t, output.Friendships, fixedLimit)

	// 模拟limit为0
	friendships = getFriendships(defaultFriendshipsLimit)

	finder.EXPECT().FindsByUserID(gomock.Any(), fixedUserID, defaultFriendshipsLimit, fixedFriendshipBaseID).Return(friendships, nil)

	output, err = useCase.Execute(nil, &application.ListFriendshipsInput{
		UserID: fixedUserID,
		BaseID: fixedFriendshipBaseID,
		Limit:  0,
	})
	require.NoError(t, err)
	require.NotNil(t, output)

	require.Len(t, output.Friendships, defaultFriendshipsLimit)

	//模拟负数limit
	finder.EXPECT().FindsByUserID(gomock.Any(), fixedUserID, defaultFriendshipsLimit, fixedFriendshipBaseID).Return(friendships, nil)

	output, err = useCase.Execute(nil, &application.ListFriendshipsInput{
		UserID: fixedUserID,
		BaseID: fixedFriendshipBaseID,
		Limit:  -100,
	})
	require.NoError(t, err)
	require.NotNil(t, output)

	require.Len(t, output.Friendships, defaultFriendshipsLimit)

	// 模拟太大的limit
	friendships = getFriendships(maxFriendshipsLimit)

	finder.EXPECT().FindsByUserID(gomock.Any(), fixedUserID, maxFriendshipsLimit, fixedFriendshipBaseID).Return(friendships, nil)

	output, err = useCase.Execute(nil, &application.ListFriendshipsInput{
		UserID: fixedUserID,
		BaseID: fixedFriendshipBaseID,
		Limit:  1000,
	})
	require.NoError(t, err)
	require.NotNil(t, output)

	require.Len(t, output.Friendships, maxFriendshipsLimit)
}

func getFriendships(count int) []*domain.Friendship {
	friendships := make([]*domain.Friendship, count)
	for i := 0; i < count; i++ {
		friendships[i] = domain.LoadFriendship(
			fixedFriendshipID,
			fixedUserID,
			fixedToID,
			time.Now().UTC(),
		)
	}
	return friendships
}
