package application_test

import (
	"gochat/internal/roomship/application"
	"gochat/internal/roomship/domain"
	"gochat/internal/roomship/domain/mocks"
	myErrors "gochat/internal/shared/errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const (
	fixedRoomshipBaseID   domain.RoomshipID = "roomship-123"
	maxRoomshipsLimit                       = 100
	defaultRoomshipsLimit                   = 50
)

func TestListRoomshipsInput_Validate(t *testing.T) {
	input := &application.ListRoomshipsInput{
		UserID: fixedUserID,
		BaseID: fixedRoomshipBaseID,
		Limit:  maxRoomshipsLimit,
	}

	err := input.Validate()
	require.NoError(t, err)

	inputWithNoBaseID := &application.ListRoomshipsInput{
		UserID: fixedUserID,
		BaseID: "",
		Limit:  maxRoomshipsLimit,
	}

	err = inputWithNoBaseID.Validate()
	require.NoError(t, err)

	inputWithNoLimit := &application.ListRoomshipsInput{
		UserID: fixedUserID,
		BaseID: fixedRoomshipBaseID,
		Limit:  0,
	}

	err = inputWithNoLimit.Validate()
	require.NoError(t, err)

	inputWithNegativeLimit := &application.ListRoomshipsInput{
		UserID: fixedUserID,
		BaseID: fixedRoomshipBaseID,
		Limit:  -10,
	}

	err = inputWithNegativeLimit.Validate()
	require.NoError(t, err)

	inputWithNoBaseIDAndLimit := &application.ListRoomshipsInput{
		UserID: fixedUserID,
		BaseID: "",
		Limit:  0,
	}

	err = inputWithNoBaseIDAndLimit.Validate()
	require.NoError(t, err)

	inputWithEmptyUserID := &application.ListRoomshipsInput{
		UserID: "",
		BaseID: fixedRoomshipBaseID,
		Limit:  maxRoomshipsLimit,
	}

	err = inputWithEmptyUserID.Validate()
	require.ErrorIs(t, err, myErrors.ErrEmptyInput)
}

func TestNewListRoomshipsUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	finder := mocks.NewMockRoomshipsFinderByUserID(ctrl)

	useCase, err := application.NewListRoomshipsUseCase(
		finder,
	)

	require.NoError(t, err)
	require.NotNil(t, useCase)

	useCaseWithNil, err := application.NewListRoomshipsUseCase(
		nil,
	)

	require.ErrorIs(t, err, myErrors.ErrEmptyPointer)
	require.Nil(t, useCaseWithNil)
}

func TestListRoomshipsUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	finder := mocks.NewMockRoomshipsFinderByUserID(ctrl)

	useCase, err := application.NewListRoomshipsUseCase(
		finder,
	)
	require.NoError(t, err)
	require.NotNil(t, useCase)

	// 模拟正常情况，limit在范围内
	fixedLimit := maxRoomshipsLimit - 1
	roomships := getRoomships(fixedLimit)

	finder.EXPECT().FindsByUserID(gomock.Any(), fixedUserID, fixedLimit, fixedRoomshipBaseID).Return(roomships, nil)

	output, err := useCase.Execute(nil, &application.ListRoomshipsInput{
		UserID: fixedUserID,
		BaseID: fixedRoomshipBaseID,
		Limit:  maxRoomshipsLimit - 1,
	})
	require.NoError(t, err)
	require.NotNil(t, output)

	require.Len(t, output.Roomships, fixedLimit)

	// 模拟limit为0
	roomships = getRoomships(defaultRoomshipsLimit)

	finder.EXPECT().FindsByUserID(gomock.Any(), fixedUserID, defaultRoomshipsLimit, fixedRoomshipBaseID).Return(roomships, nil)

	output, err = useCase.Execute(nil, &application.ListRoomshipsInput{
		UserID: fixedUserID,
		BaseID: fixedRoomshipBaseID,
		Limit:  0,
	})
	require.NoError(t, err)
	require.NotNil(t, output)

	require.Len(t, output.Roomships, defaultRoomshipsLimit)

	//模拟负数limit
	finder.EXPECT().FindsByUserID(gomock.Any(), fixedUserID, defaultRoomshipsLimit, fixedRoomshipBaseID).Return(roomships, nil)

	output, err = useCase.Execute(nil, &application.ListRoomshipsInput{
		UserID: fixedUserID,
		BaseID: fixedRoomshipBaseID,
		Limit:  -100,
	})
	require.NoError(t, err)
	require.NotNil(t, output)

	require.Len(t, output.Roomships, defaultRoomshipsLimit)

	// 模拟太大的limit
	roomships = getRoomships(maxRoomshipsLimit)

	finder.EXPECT().FindsByUserID(gomock.Any(), fixedUserID, maxRoomshipsLimit, fixedRoomshipBaseID).Return(roomships, nil)

	output, err = useCase.Execute(nil, &application.ListRoomshipsInput{
		UserID: fixedUserID,
		BaseID: fixedRoomshipBaseID,
		Limit:  1000,
	})
	require.NoError(t, err)
	require.NotNil(t, output)

	require.Len(t, output.Roomships, maxRoomshipsLimit)
}

func getRoomships(count int) []*domain.Roomship {
	roomships := make([]*domain.Roomship, count)
	for i := 0; i < count; i++ {
		roomships[i] = domain.LoadRoomship(
			fixedRoomshipID,
			fixedRoomID,
			fixedUserID,
			domain.MemberRole,
			time.Now().UTC(),
		)
	}
	return roomships
}
