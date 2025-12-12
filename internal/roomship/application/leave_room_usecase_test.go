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

func TestLeaveRoomInput_Validate(t *testing.T) {
	input := &application.LeaveRoomInput{
		UserID: fixedUserID,
		RoomID: fixedRoomID,
	}

	err := input.Validate()
	require.NoError(t, err)

	inputWithNoUserID := &application.LeaveRoomInput{
		UserID: "",
		RoomID: fixedRoomID,
	}

	err = inputWithNoUserID.Validate()
	require.ErrorIs(t, err, myErrors.ErrEmptyInput)

	inputWithNoRoomID := &application.LeaveRoomInput{
		RoomID: "",
	}

	err = inputWithNoRoomID.Validate()
	require.ErrorIs(t, err, myErrors.ErrEmptyInput)
}

func TestNewLeaveRoomUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRoomshipDeleter := mocks.NewMockRoomshipDeleterByUserIDAndRoomID(ctrl)
	mockRoomDeleter := mocks.NewMockRoomDeleterByID(ctrl)
	mockRoomshipFinder := mocks.NewMockRoomshipFinderByUserIDAndRoomID(ctrl)
	mockRoomshipsFinder := mocks.NewMockRoomshipsFinderByRoomID(ctrl)

	useCase, err := application.NewLeaveRoomUseCase(
		mockRoomshipDeleter,
		mockRoomDeleter,
		mockRoomshipFinder,
		mockRoomshipsFinder,
	)

	require.NoError(t, err)
	require.NotNil(t, useCase)

	useCaseWithNil, err := application.NewLeaveRoomUseCase(
		nil, nil, nil, nil,
	)

	require.ErrorIs(t, err, myErrors.ErrEmptyPointer)
	require.Nil(t, useCaseWithNil)
}

func TestLeaveRoomUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRoomshipDeleter := mocks.NewMockRoomshipDeleterByUserIDAndRoomID(ctrl)
	mockRoomDeleter := mocks.NewMockRoomDeleterByID(ctrl)
	mockRoomshipFinder := mocks.NewMockRoomshipFinderByUserIDAndRoomID(ctrl)
	mockRoomshipsFinder := mocks.NewMockRoomshipsFinderByRoomID(ctrl)

	useCase, err := application.NewLeaveRoomUseCase(
		mockRoomshipDeleter,
		mockRoomDeleter,
		mockRoomshipFinder,
		mockRoomshipsFinder,
	)
	require.NoError(t, err)
	require.NotNil(t, useCase)

	// 正常情况
	fixedMemberRoomship := domain.LoadRoomship(
		fixedRoomshipID,
		fixedRoomID,
		fixedUserID,
		domain.MemberRole,
		time.Now().UTC(),
	)
	gomock.InOrder(
		mockRoomshipFinder.EXPECT().FindByUserIDAndRoomID(nil, fixedUserID, fixedRoomID).Return(fixedMemberRoomship, nil).Times(1),
		mockRoomshipDeleter.EXPECT().DeleteByUserIDAndRoomID(nil, fixedUserID, fixedRoomID).Return(nil).Times(1),
		mockRoomshipsFinder.EXPECT().FindsByRoomID(nil, fixedRoomID).Return(getRoomships(10), nil).Times(1),
	)

	_, err = useCase.Execute(nil, &application.LeaveRoomInput{
		UserID: fixedUserID,
		RoomID: fixedRoomID,
	})
	require.NoError(t, err)

	// 房主退出房间
	fixedOwnerRoomship := domain.LoadRoomship(
		fixedRoomshipID,
		fixedRoomID,
		fixedUserID,
		domain.OwnerRole,
		time.Now().UTC(),
	)
	gomock.InOrder(
		mockRoomshipFinder.EXPECT().FindByUserIDAndRoomID(nil, fixedUserID, fixedRoomID).Return(fixedOwnerRoomship, nil).Times(1),
	)

	_, err = useCase.Execute(nil, &application.LeaveRoomInput{
		UserID: fixedUserID,
		RoomID: fixedRoomID,
	})

	// 退出后房间为空
	gomock.InOrder(
		mockRoomshipFinder.EXPECT().FindByUserIDAndRoomID(nil, fixedUserID, fixedRoomID).Return(fixedMemberRoomship, nil).Times(1),
		mockRoomshipDeleter.EXPECT().DeleteByUserIDAndRoomID(nil, fixedUserID, fixedRoomID).Return(nil).Times(1),
		mockRoomshipsFinder.EXPECT().FindsByRoomID(nil, fixedRoomID).Return(nil, nil).Times(1),
		mockRoomDeleter.EXPECT().DeleteByID(nil, fixedRoomID).Return(nil).Times(1),
	)

	_, err = useCase.Execute(nil, &application.LeaveRoomInput{
		UserID: fixedUserID,
		RoomID: fixedRoomID,
	})
	require.NoError(t, err)
}
