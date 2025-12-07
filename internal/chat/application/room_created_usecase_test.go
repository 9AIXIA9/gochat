package application_test

import (
	"gochat/internal/chat/application"
	"gochat/internal/chat/domain/mocks"
	myErrors "gochat/internal/shared/errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestRoomCreatedInput_Validate(t *testing.T) {
	input := &application.RoomCreatedInput{
		RoomID: fixedRoomID,
	}

	err := input.Validate()
	require.NoError(t, err)

	inputWithEmptyRoomID := &application.RoomCreatedInput{
		RoomID: "",
	}

	err = inputWithEmptyRoomID.Validate()
	require.ErrorIs(t, err, myErrors.ErrEmptyInput)
}

func TestNewRoomCreatedUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRoomSaver := mocks.NewMockRoomSaver(ctrl)

	useCase, err := application.NewRoomCreatedUseCase(
		mockRoomSaver,
	)

	require.NoError(t, err)
	require.NotNil(t, useCase)

	useCaseWithNil, err := application.NewRoomCreatedUseCase(
		nil,
	)

	require.ErrorIs(t, err, myErrors.ErrEmptyPointer)
	require.Nil(t, useCaseWithNil)
}

func TestRoomCreatedUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRoomSaver := mocks.NewMockRoomSaver(ctrl)

	useCase, err := application.NewRoomCreatedUseCase(
		mockRoomSaver,
	)
	require.NoError(t, err)
	require.NotNil(t, useCase)

	// 正常情况
	mockRoomSaver.EXPECT().Save(nil, gomock.Any()).Return(nil).Times(1)
	_, err = useCase.Execute(nil, &application.RoomCreatedInput{
		RoomID: fixedRoomID,
	})
	require.NoError(t, err)
}
