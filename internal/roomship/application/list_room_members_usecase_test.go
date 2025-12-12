package application_test

import (
	"gochat/internal/roomship/application"
	"gochat/internal/roomship/domain/mocks"
	myErrors "gochat/internal/shared/errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestListRoomMembersInput_Validate(t *testing.T) {
	input := &application.ListRoomMembersInput{
		RoomID: fixedRoomID,
	}

	err := input.Validate()
	require.NoError(t, err)

	inputWithNoRoomID := &application.ListRoomMembersInput{
		RoomID: "",
	}

	err = inputWithNoRoomID.Validate()
	require.ErrorIs(t, err, myErrors.ErrEmptyInput)
}

func TestNewListRoomMembersUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	finder := mocks.NewMockRoomshipsFinderByRoomID(ctrl)

	useCase, err := application.NewListRoomMembersUseCase(
		finder,
	)

	require.NoError(t, err)
	require.NotNil(t, useCase)

	useCaseWithNil, err := application.NewListRoomMembersUseCase(
		nil,
	)

	require.ErrorIs(t, err, myErrors.ErrEmptyPointer)
	require.Nil(t, useCaseWithNil)
}

func TestListRoomMembersUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	finder := mocks.NewMockRoomshipsFinderByRoomID(ctrl)

	useCase, err := application.NewListRoomMembersUseCase(
		finder,
	)
	require.NoError(t, err)
	require.NotNil(t, useCase)

	// 正常情况
	fixedLimit := maxRoomshipsLimit - 1
	roomships := getRoomships(fixedLimit)

	finder.EXPECT().FindsByRoomID(gomock.Any(), fixedRoomID).Return(roomships, nil)

	output, err := useCase.Execute(nil, &application.ListRoomMembersInput{
		RoomID: fixedRoomID,
	})
	require.NoError(t, err)
	require.NotNil(t, output)

	require.Len(t, output.Roomships, fixedLimit)
}
