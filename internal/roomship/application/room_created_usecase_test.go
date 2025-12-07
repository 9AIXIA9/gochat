package application_test

import (
	"gochat/internal/roomship/application"
	"gochat/internal/roomship/domain"
	"gochat/internal/roomship/domain/mocks"
	myErrors "gochat/internal/shared/errors"
	eventMock "gochat/internal/shared/event/mocks"
	"testing"
	"time"

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

	mockRoomshipIDGenerator := mocks.NewMockRoomshipIDGenerator(ctrl)
	mockIDGenerator := eventMock.NewMockIDGenerator(ctrl)
	mockFinder := mocks.NewMockRoomFinderByID(ctrl)
	mockRoomshipCreator := mocks.NewMockRoomshipCreator(ctrl)
	mockEventCreator := eventMock.NewMockUnpublishedEventsCreator(ctrl)

	useCase, err := application.NewRoomCreatedUseCase(
		mockRoomshipIDGenerator,
		mockIDGenerator,
		mockFinder,
		mockRoomshipCreator,
		mockEventCreator,
	)

	require.NoError(t, err)
	require.NotNil(t, useCase)

	useCaseWithNil, err := application.NewRoomCreatedUseCase(
		nil, nil, nil, nil, nil,
	)

	require.ErrorIs(t, err, myErrors.ErrEmptyPointer)
	require.Nil(t, useCaseWithNil)
}

func TestRoomCreatedUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRoomshipIDGenerator := mocks.NewMockRoomshipIDGenerator(ctrl)
	mockIDGenerator := eventMock.NewMockIDGenerator(ctrl)
	mockFinder := mocks.NewMockRoomFinderByID(ctrl)
	mockRoomshipCreator := mocks.NewMockRoomshipCreator(ctrl)
	mockEventCreator := eventMock.NewMockUnpublishedEventsCreator(ctrl)

	useCase, err := application.NewRoomCreatedUseCase(
		mockRoomshipIDGenerator,
		mockIDGenerator,
		mockFinder,
		mockRoomshipCreator,
		mockEventCreator,
	)

	require.NoError(t, err)
	require.NotNil(t, useCase)

	mockRoom := domain.LoadRoom(
		fixedRoomID,
		fixedOwnerID,
		fixedRoomNumber,
		fixedPasswordEncrypted,
		fixedMaxMemberCount,
		time.Now().UTC(),
	)

	gomock.InOrder(
		mockFinder.EXPECT().FindByID(nil, fixedRoomID).Return(mockRoom, nil),
		mockRoomshipIDGenerator.EXPECT().Generate().Return(fixedRoomshipID),
		mockIDGenerator.EXPECT().Generate().Return(fixedEventID),
		mockRoomshipCreator.EXPECT().Create(nil, gomock.Any()).Return(nil),
		mockIDGenerator.EXPECT().Generate().Return(fixedEventID),
		mockEventCreator.EXPECT().CreateUnpublishedEvents(nil, gomock.Any()).Return(nil),
	)

	_, err = useCase.Execute(nil, &application.RoomCreatedInput{
		RoomID: fixedRoomID,
	})
	require.NoError(t, err)
}
