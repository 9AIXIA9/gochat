package application_test

import (
	"gochat/internal/profile/application"
	"gochat/internal/profile/domain/mocks"
	myErrors "gochat/internal/shared/errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestRoomCreatedInput_Validate(t *testing.T) {
	input := &application.RoomCreatedInput{
		RoomID:    fixedRoomID,
		CreatedAt: time.Now().UTC(),
	}

	err := input.Validate()
	require.NoError(t, err)

	inputWithEmptyRoomID := &application.RoomCreatedInput{
		RoomID:    "",
		CreatedAt: time.Now().UTC(),
	}

	err = inputWithEmptyRoomID.Validate()
	require.ErrorIs(t, err, myErrors.ErrEmptyInput)

	inputWithEmptyCreatedAt := &application.RoomCreatedInput{
		RoomID:    fixedRoomID,
		CreatedAt: time.Time{},
	}

	err = inputWithEmptyCreatedAt.Validate()
	require.ErrorIs(t, err, myErrors.ErrEmptyInput)
}

func TestNewRoomCreatedUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRoomSaver := mocks.NewMockRoomSaver(ctrl)
	mockCreator := mocks.NewMockRoomProfileCreator(ctrl)

	useCase, err := application.NewRoomCreatedUseCase(
		mockRoomSaver,
		mockCreator,
	)

	require.NoError(t, err)
	require.NotNil(t, useCase)

	useCaseWithNil, err := application.NewRoomCreatedUseCase(
		nil, nil,
	)

	require.ErrorIs(t, err, myErrors.ErrEmptyPointer)
	require.Nil(t, useCaseWithNil)
}

func TestRoomCreatedUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRoomSaver := mocks.NewMockRoomSaver(ctrl)
	mockCreator := mocks.NewMockRoomProfileCreator(ctrl)

	useCase, err := application.NewRoomCreatedUseCase(
		mockRoomSaver,
		mockCreator,
	)
	require.NoError(t, err)
	require.NotNil(t, useCase)

	// 正常情况
	gomock.InOrder(
		mockRoomSaver.EXPECT().Save(nil, gomock.Any()).Return(nil).Times(1),
		mockCreator.EXPECT().Create(nil, gomock.Any()).Return(nil).Times(1),
	)

	_, err = useCase.Execute(nil, &application.RoomCreatedInput{
		RoomID: fixedRoomID,
	})
	require.NoError(t, err)

	// 重复创建
	gomock.InOrder(
		mockRoomSaver.EXPECT().Save(nil, gomock.Any()).Return(nil).Times(1),
		mockCreator.EXPECT().Create(nil, gomock.Any()).Return(myErrors.ErrDuplicatedKey).Times(1),
	)

	_, err = useCase.Execute(nil, &application.RoomCreatedInput{
		RoomID: fixedRoomID,
	})
	require.NoError(t, err)
}
