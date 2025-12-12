package application_test

import (
	"gochat/internal/profile/application"
	"gochat/internal/profile/domain"
	"gochat/internal/profile/domain/mocks"
	myErrors "gochat/internal/shared/errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetRoomProfileInput_Validate(t *testing.T) {
	input := &application.GetRoomProfileInput{
		RoomID: fixedRoomID,
	}

	err := input.Validate()
	require.NoError(t, err)

	inputWithEmptyRoomID := &application.GetRoomProfileInput{
		RoomID: "",
	}

	err = inputWithEmptyRoomID.Validate()
	require.ErrorIs(t, err, myErrors.ErrEmptyInput)
}

func TestNewGetRoomProfileUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProfileFinder := mocks.NewMockRoomProfileFinder(ctrl)

	useCase, err := application.NewGetRoomProfileUseCase(
		mockProfileFinder,
	)

	require.NoError(t, err)
	require.NotNil(t, useCase)

	useCaseWithNil, err := application.NewGetRoomProfileUseCase(
		nil,
	)

	require.ErrorIs(t, err, myErrors.ErrEmptyPointer)
	require.Nil(t, useCaseWithNil)
}

func TestGetRoomProfileUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProfileFinder := mocks.NewMockRoomProfileFinder(ctrl)

	useCase, err := application.NewGetRoomProfileUseCase(
		mockProfileFinder,
	)

	require.NoError(t, err)
	require.NotNil(t, useCase)

	// 正常情况
	gomock.InOrder(
		mockProfileFinder.EXPECT().FindByID(nil, fixedRoomID).Return(domain.LoadRoomProfile(
			fixedRoomID,
			fixedName,
			fixedSign,
			time.Now().UTC(),
		), nil).Times(1),
	)

	_, err = useCase.Execute(nil, &application.GetRoomProfileInput{
		RoomID: fixedRoomID,
	})
	require.NoError(t, err)
}
