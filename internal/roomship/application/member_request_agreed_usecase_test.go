package application_test

import (
	"gochat/internal/roomship/application"
	"gochat/internal/roomship/domain/mocks"
	myErrors "gochat/internal/shared/errors"
	eventMock "gochat/internal/shared/event/mocks"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestMemberRequestAgreedInput_Validate(t *testing.T) {
	input := &application.MemberRequestAgreedInput{
		UserID: fixedUserID,
		RoomID: fixedRoomID,
	}

	err := input.Validate()
	require.NoError(t, err)

	inputWithEmptyRoomID := &application.MemberRequestAgreedInput{
		UserID: fixedUserID,
		RoomID: "",
	}

	err = inputWithEmptyRoomID.Validate()
	require.ErrorIs(t, err, myErrors.ErrEmptyInput)

	inputWithEmptyUserID := &application.MemberRequestAgreedInput{
		UserID: "",
		RoomID: fixedRoomID,
	}

	err = inputWithEmptyUserID.Validate()
	require.ErrorIs(t, err, myErrors.ErrEmptyInput)
}

func TestNewMemberRequestAgreedUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRoomshipIDGenerator := mocks.NewMockRoomshipIDGenerator(ctrl)
	mockIDGenerator := eventMock.NewMockIDGenerator(ctrl)
	mockCreator := mocks.NewMockRoomshipCreator(ctrl)

	useCase, err := application.NewMemberRequestAgreedUseCase(
		mockRoomshipIDGenerator,
		mockIDGenerator,
		mockCreator,
	)

	require.NoError(t, err)
	require.NotNil(t, useCase)

	useCaseWithNil, err := application.NewMemberRequestAgreedUseCase(
		nil, nil, nil,
	)

	require.ErrorIs(t, err, myErrors.ErrEmptyPointer)
	require.Nil(t, useCaseWithNil)
}

func TestMemberRequestAgreedUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRoomshipIDGenerator := mocks.NewMockRoomshipIDGenerator(ctrl)
	mockIDGenerator := eventMock.NewMockIDGenerator(ctrl)
	mockCreator := mocks.NewMockRoomshipCreator(ctrl)

	useCase, err := application.NewMemberRequestAgreedUseCase(
		mockRoomshipIDGenerator,
		mockIDGenerator,
		mockCreator,
	)

	require.NoError(t, err)
	require.NotNil(t, useCase)

	gomock.InOrder(
		mockRoomshipIDGenerator.EXPECT().Generate().Return(fixedRoomshipID),
		mockIDGenerator.EXPECT().Generate().Return(fixedEventID),
		mockCreator.EXPECT().Create(nil, gomock.Any()).Return(nil),
	)

	_, err = useCase.Execute(nil, &application.MemberRequestAgreedInput{
		UserID: fixedUserID,
		RoomID: fixedRoomID,
	})
	require.NoError(t, err)
}
