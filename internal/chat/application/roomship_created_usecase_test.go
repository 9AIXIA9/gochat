package application_test

import (
	"gochat/internal/chat/application"
	"gochat/internal/chat/domain/mocks"
	myErrors "gochat/internal/shared/errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestRoomshipCreatedInput_Validate(t *testing.T) {
	input := &application.RoomshipCreatedInput{
		ID:     fixedRoomshipID,
		UserID: fixedUserID,
		RoomID: fixedRoomID,
	}

	err := input.Validate()
	require.NoError(t, err)

	inputWithEmptyRoomshipID := &application.RoomshipCreatedInput{
		ID:     "",
		UserID: fixedUserID,
		RoomID: fixedRoomID,
	}

	err = inputWithEmptyRoomshipID.Validate()
	require.ErrorIs(t, err, myErrors.ErrEmptyInput)

	inputWithEmptyUserID := &application.RoomshipCreatedInput{
		ID:     fixedRoomshipID,
		UserID: "",
		RoomID: fixedRoomID,
	}

	err = inputWithEmptyUserID.Validate()
	require.ErrorIs(t, err, myErrors.ErrEmptyInput)

	inputWithEmptyRoomID := &application.RoomshipCreatedInput{
		ID:     fixedRoomshipID,
		RoomID: "",
		UserID: fixedUserID,
	}

	err = inputWithEmptyRoomID.Validate()
	require.ErrorIs(t, err, myErrors.ErrEmptyInput)
}

func TestNewRoomshipCreatedUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRoomshipSaver := mocks.NewMockRoomshipSaver(ctrl)

	useCase, err := application.NewRoomshipCreatedUseCase(
		mockRoomshipSaver,
	)

	require.NoError(t, err)
	require.NotNil(t, useCase)

	useCaseWithNil, err := application.NewRoomshipCreatedUseCase(
		nil,
	)

	require.ErrorIs(t, err, myErrors.ErrEmptyPointer)
	require.Nil(t, useCaseWithNil)
}

func TestRoomshipCreatedUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRoomshipSaver := mocks.NewMockRoomshipSaver(ctrl)

	useCase, err := application.NewRoomshipCreatedUseCase(
		mockRoomshipSaver,
	)
	require.NoError(t, err)
	require.NotNil(t, useCase)

	// 正常情况
	mockRoomshipSaver.EXPECT().Save(nil, gomock.Any()).Return(nil).Times(1)
	_, err = useCase.Execute(nil, &application.RoomshipCreatedInput{
		ID:     fixedRoomshipID,
		UserID: fixedUserID,
		RoomID: fixedRoomID,
	})
	require.NoError(t, err)
}
