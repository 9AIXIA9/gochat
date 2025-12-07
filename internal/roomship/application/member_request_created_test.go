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

func TestMemberRequestCreatedInput_Validate(t *testing.T) {
	input := &application.MemberRequestCreatedInput{
		RequestID: fixedOperationID,
	}

	err := input.Validate()
	require.NoError(t, err)

	inputWithEmptyRequestID := &application.MemberRequestCreatedInput{
		RequestID: "",
	}

	err = inputWithEmptyRequestID.Validate()
	require.ErrorIs(t, err, myErrors.ErrEmptyInput)
}

func TestNewMemberRequestCreatedUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRequestFinder := mocks.NewMockMemberRequestFinderByID(ctrl)
	mockRoomshipFinder := mocks.NewMockRoomshipsFinderByRoomIDAndRole(ctrl)
	mockIDGenerator := eventMock.NewMockIDGenerator(ctrl)
	mockCreator := eventMock.NewMockUnpublishedEventsCreator(ctrl)

	useCase, err := application.NewMemberRequestCreatedUseCase(
		mockRequestFinder,
		mockRoomshipFinder,
		mockIDGenerator,
		mockCreator,
	)

	require.NoError(t, err)
	require.NotNil(t, useCase)

	useCaseWithNil, err := application.NewMemberRequestCreatedUseCase(
		nil, nil, nil, nil,
	)

	require.ErrorIs(t, err, myErrors.ErrEmptyPointer)
	require.Nil(t, useCaseWithNil)
}

func TestMemberRequestCreatedUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRequestFinder := mocks.NewMockMemberRequestFinderByID(ctrl)
	mockRoomshipFinder := mocks.NewMockRoomshipsFinderByRoomIDAndRole(ctrl)
	mockIDGenerator := eventMock.NewMockIDGenerator(ctrl)
	mockCreator := eventMock.NewMockUnpublishedEventsCreator(ctrl)

	useCase, err := application.NewMemberRequestCreatedUseCase(
		mockRequestFinder,
		mockRoomshipFinder,
		mockIDGenerator,
		mockCreator,
	)
	require.NoError(t, err)
	require.NotNil(t, useCase)

	mockRequest := domain.LoadMemberRequest(
		fixedOperationID,
		domain.StatePending,
		fixedUserID,
		fixedRoomID,
		fixedContent,
		"",
		time.Now().UTC(),
		time.Now().UTC(),
	)

	mockOwnerRoomship := make([]*domain.Roomship, 10)
	for i := 0; i < 10; i++ {
		mockOwnerRoomship[i] = domain.LoadRoomship(
			domain.RoomshipID(fmt.Sprintf("owner-roomship-id-%d", i+1)),
			fixedRoomID,
			kernel.UserID(fmt.Sprintf("owner-user-id-%d", i+1)),
			domain.MemberRole,
			time.Now().UTC(),
		)
	}

	gomock.InOrder(
		mockRequestFinder.EXPECT().FindByID(nil, fixedOperationID).Return(mockRequest, nil),
		mockRoomshipFinder.EXPECT().FindsByRoomIDAndRole(nil, fixedRoomID, domain.OwnerRole).Return(mockOwnerRoomship, nil),
		mockIDGenerator.EXPECT().Generate().Return(fixedEventID).Times(len(mockOwnerRoomship)),
		mockCreator.EXPECT().CreateUnpublishedEvents(nil, gomock.Any()).Return(nil),
	)

	_, err = useCase.Execute(nil, &application.MemberRequestCreatedInput{
		RequestID: fixedOperationID,
	})
	require.NoError(t, err)
}
