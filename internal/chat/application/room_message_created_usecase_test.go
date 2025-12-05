package application_test

import (
	"fmt"
	"gochat/internal/chat/application"
	"gochat/internal/chat/domain"
	"gochat/internal/chat/domain/mocks"
	myErrors "gochat/internal/shared/errors"
	eventMock "gochat/internal/shared/event/mocks"
	"gochat/internal/shared/kernel"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestRoomMessageCreatedInput_Validate(t *testing.T) {
	input := &application.RoomMessageCreatedInput{
		MessageID: fixedMessageID,
	}

	err := input.Validate()
	require.NoError(t, err)

	inputWithEmptyMessageID := &application.RoomMessageCreatedInput{
		MessageID: "",
	}

	err = inputWithEmptyMessageID.Validate()
	require.ErrorIs(t, err, myErrors.ErrEmptyInput)
}

func TestNewRoomMessageCreatedUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEventIDGenerator := eventMock.NewMockIDGenerator(ctrl)
	mockCreator := eventMock.NewMockUnpublishedEventCreator(ctrl)
	mockRoomMessageFinder := mocks.NewMockRoomMessageFinder(ctrl)
	mockRoomshipFinder := mocks.NewMockRoomshipsFinderByRoomID(ctrl)

	useCase, err := application.NewRoomMessageCreatedUseCase(
		mockEventIDGenerator,
		mockRoomMessageFinder,
		mockRoomshipFinder,
		mockCreator,
	)

	require.NoError(t, err)
	require.NotNil(t, useCase)

	useCaseWithNil, err := application.NewRoomMessageCreatedUseCase(
		nil, nil, nil, nil,
	)

	require.ErrorIs(t, err, myErrors.ErrEmptyPointer)
	require.Nil(t, useCaseWithNil)
}

func TestRoomMessageCreatedUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEventIDGenerator := eventMock.NewMockIDGenerator(ctrl)
	mockCreator := eventMock.NewMockUnpublishedEventCreator(ctrl)
	mockRoomMessageFinder := mocks.NewMockRoomMessageFinder(ctrl)
	mockRoomshipFinder := mocks.NewMockRoomshipsFinderByRoomID(ctrl)

	useCase, err := application.NewRoomMessageCreatedUseCase(
		mockEventIDGenerator,
		mockRoomMessageFinder,
		mockRoomshipFinder,
		mockCreator,
	)
	require.NoError(t, err)
	require.NotNil(t, useCase)

	// 正常情况
	mockMessage := domain.LoadRoomMessage(
		fixedMessageID,
		fixedUserID,
		fixedRoomID,
		fixedContent,
		time.Now().UTC(),
	)

	mockRoomships := make([]*domain.Roomship, 0, 10)
	for i := 0; i < 10; i++ {
		mockRoomship := domain.LoadRoomship(
			domain.RoomshipID(fmt.Sprintf("roomship-id-%d", i)),
			kernel.UserID(fmt.Sprintf("user-id-%d", i)),
			fixedRoomID,
		)
		mockRoomships = append(mockRoomships, mockRoomship)
	}

	gomock.InOrder(
		mockRoomMessageFinder.EXPECT().FindRoomMessage(nil, fixedMessageID).Return(mockMessage, nil).Times(1),
		mockRoomshipFinder.EXPECT().FindsByRoomID(nil, fixedRoomID).Return(mockRoomships, nil).Times(1),
		mockEventIDGenerator.EXPECT().Generate().Return(fixedEventID).Times(1),
		mockCreator.EXPECT().CreateUnpublishedEvent(nil, gomock.Any()).Return(nil).Times(1),
	)

	_, err = useCase.Execute(nil, &application.RoomMessageCreatedInput{
		MessageID: fixedMessageID,
	})
	require.NoError(t, err)

	// 房间只有自己不通知
	mockRoomshipsOnlySelf := []*domain.Roomship{
		domain.LoadRoomship(
			fixedRoomshipID,
			fixedUserID,
			fixedRoomID,
		),
	}
	gomock.InOrder(
		mockRoomMessageFinder.EXPECT().FindRoomMessage(nil, fixedMessageID).Return(mockMessage, nil).Times(1),
		mockRoomshipFinder.EXPECT().FindsByRoomID(nil, fixedRoomID).Return(mockRoomshipsOnlySelf, nil).Times(1),
	)

	_, err = useCase.Execute(nil, &application.RoomMessageCreatedInput{
		MessageID: fixedMessageID,
	})
	require.NoError(t, err)
}
