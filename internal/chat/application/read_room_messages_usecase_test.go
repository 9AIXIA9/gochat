package application_test

import (
	"gochat/internal/chat/application"
	"gochat/internal/chat/domain"
	"gochat/internal/chat/domain/mocks"
	myErrors "gochat/internal/shared/errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestReadRoomMessagesInput_Validate(t *testing.T) {
	input := &application.ReadRoomMessagesInput{
		UserID: fixedUserID,
		RoomID: fixedRoomID,
	}

	err := input.Validate()
	require.NoError(t, err)

	inputWithEmptyUserID := &application.ReadRoomMessagesInput{
		UserID: "",
		RoomID: fixedRoomID,
	}

	err = inputWithEmptyUserID.Validate()
	require.Error(t, err)

	inputWithEmptyRoomID := &application.ReadRoomMessagesInput{
		UserID: fixedUserID,
		RoomID: "",
	}

	err = inputWithEmptyRoomID.Validate()
	require.Error(t, err)
}

func TestNewReadRoomMessagesUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUpdater := mocks.NewMockRoomMessagesStatesUpdaterByUserIDAndRoomID(ctrl)

	useCase, err := application.NewReadRoomMessagesUseCase(
		mockUpdater,
	)

	require.NoError(t, err)
	require.NotNil(t, useCase)

	useCaseWithNil, err := application.NewReadRoomMessagesUseCase(
		nil,
	)

	require.ErrorIs(t, err, myErrors.ErrEmptyPointer)
	require.Nil(t, useCaseWithNil)
}

func TestReadRoomMessagesUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUpdater := mocks.NewMockRoomMessagesStatesUpdaterByUserIDAndRoomID(ctrl)

	useCase, err := application.NewReadRoomMessagesUseCase(
		mockUpdater,
	)
	require.NoError(t, err)
	require.NotNil(t, useCase)

	// 正常情况
	mockUpdater.EXPECT().UpdatesByUserIDAndRoomID(nil, fixedUserID, fixedRoomID, domain.MessageStateRead).Return(nil).Times(1)
	_, err = useCase.Execute(nil, &application.ReadRoomMessagesInput{
		UserID: fixedUserID,
		RoomID: fixedRoomID,
	})
	require.NoError(t, err)
}
