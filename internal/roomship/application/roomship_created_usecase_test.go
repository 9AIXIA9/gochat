package application_test

import (
	"fmt"
	"gochat/internal/roomship/application"
	"gochat/internal/roomship/domain"
	"gochat/internal/roomship/domain/mocks"
	myErrors "gochat/internal/shared/errors"
	eventMock "gochat/internal/shared/event/mocks"
	"gochat/internal/shared/kernel"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const fixedAllRoomshipsLen = 10

func TestRoomshipCreatedInput_Validate(t *testing.T) {
	input := &application.RoomshipCreatedInput{
		RoomshipID: fixedRoomshipID,
	}

	err := input.Validate()
	require.NoError(t, err)

	inputWithEmptyRoomshipID := &application.RoomshipCreatedInput{
		RoomshipID: "",
	}

	err = inputWithEmptyRoomshipID.Validate()
	require.ErrorIs(t, err, myErrors.ErrEmptyInput)
}

func TestNewRoomshipCreatedUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRoomshipFinderByID := mocks.NewMockRoomshipFinderByID(ctrl)
	mockRoomshipFinderByRoomID := mocks.NewMockRoomshipsFinderByRoomID(ctrl)
	mockIDGenerator := eventMock.NewMockIDGenerator(ctrl)
	mockCreator := eventMock.NewMockUnpublishedEventsCreator(ctrl)

	useCase, err := application.NewRoomshipCreatedUseCase(
		mockRoomshipFinderByID,
		mockRoomshipFinderByRoomID,
		mockIDGenerator,
		mockCreator,
	)

	require.NoError(t, err)
	require.NotNil(t, useCase)

	useCaseWithNil, err := application.NewRoomshipCreatedUseCase(
		nil, nil, nil, nil,
	)

	require.ErrorIs(t, err, myErrors.ErrEmptyPointer)
	require.Nil(t, useCaseWithNil)
}

func TestRoomshipCreatedUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRoomshipFinderByID := mocks.NewMockRoomshipFinderByID(ctrl)
	mockRoomshipFinderByRoomID := mocks.NewMockRoomshipsFinderByRoomID(ctrl)
	mockIDGenerator := eventMock.NewMockIDGenerator(ctrl)
	mockCreator := eventMock.NewMockUnpublishedEventsCreator(ctrl)

	useCase, err := application.NewRoomshipCreatedUseCase(
		mockRoomshipFinderByID,
		mockRoomshipFinderByRoomID,
		mockIDGenerator,
		mockCreator,
	)
	require.NoError(t, err)
	require.NotNil(t, useCase)

	mockRoomship := domain.LoadRoomship(
		fixedRoomshipID,
		fixedRoomID,
		fixedUserID,
		domain.MemberRole,
		time.Now(),
	)

	mockAllRoomships := make([]*domain.Roomship, 0, fixedAllRoomshipsLen)
	for i := 0; i < fixedAllRoomshipsLen-1; i++ {
		mockAllRoomships = append(mockAllRoomships, domain.LoadRoomship(
			domain.RoomshipID(fmt.Sprintf("old-roomship-id-%d", i+1)),
			fixedRoomID,
			kernel.UserID(fmt.Sprintf("old-user-id-%d", i+1)),
			domain.MemberRole,
			time.Now().UTC(),
		))
	}

	mockAllRoomships = append(mockAllRoomships, mockRoomship)

	// 正常情况
	gomock.InOrder(
		mockRoomshipFinderByID.EXPECT().FindByID(nil, fixedRoomshipID).Return(mockRoomship, nil),
		mockRoomshipFinderByRoomID.EXPECT().FindsByRoomID(nil, fixedRoomID).Return(mockAllRoomships, nil),
		mockIDGenerator.EXPECT().Generate().Return(fixedEventID).Times(len(mockAllRoomships)+1),
		mockCreator.EXPECT().CreateUnpublishedEvents(nil, gomock.Any()).Return(nil),
	)

	_, err = useCase.Execute(nil, &application.RoomshipCreatedInput{
		RoomshipID: fixedRoomshipID,
	})
	require.NoError(t, err)
}
