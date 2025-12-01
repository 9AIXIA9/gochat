package application_test

import (
	"gochat/internal/friendship/application"
	"gochat/internal/friendship/domain"
	"gochat/internal/friendship/domain/mocks"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	eventMocks "gochat/internal/shared/event/mocks"
	"gochat/internal/shared/kernel"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const (
	fixedFromID      kernel.UserID      = "from-user-123"
	fixedFromNumber  kernel.UserNumber  = "10086"
	fixedToID        kernel.UserID      = "to-user-456"
	fixedToNumber    kernel.UserNumber  = "10087"
	fixedOperationID kernel.OperationID = "operation-789"
	fixedContent                        = "Let's be friends!"
	fixedEventID     event.ID           = "event-0001"
)

func TestSendFriendRequestInput_Validate_Success(t *testing.T) {
	input := &application.SendFriendRequestInput{
		FromID:   fixedFromID,
		ToNumber: fixedToNumber,
		Content:  fixedContent,
	}

	err := input.Validate()
	require.NoError(t, err)

	inputWithNoContent := &application.SendFriendRequestInput{
		FromID:   fixedFromID,
		ToNumber: fixedToNumber,
	}

	err = inputWithNoContent.Validate()
	require.NoError(t, err)
}

func TestSendFriendRequestInput_Validate_EmptyInput(t *testing.T) {
	inputWithEmptyFromID := &application.SendFriendRequestInput{
		FromID:   "",
		ToNumber: fixedToNumber,
		Content:  fixedContent,
	}

	err := inputWithEmptyFromID.Validate()
	require.ErrorIs(t, err, myErrors.ErrEmptyInput)

	inputWithEmptyToNumber := &application.SendFriendRequestInput{
		FromID:   fixedFromID,
		ToNumber: "",
		Content:  fixedContent,
	}

	err = inputWithEmptyToNumber.Validate()
	require.ErrorIs(t, err, myErrors.ErrEmptyInput)
}

func TestSendFriendRequestInput_Validate_InvalidNumber(t *testing.T) {
	inputWithWrongNumber := &application.SendFriendRequestInput{
		FromID:   fixedFromID,
		ToNumber: "wrong-number!",
		Content:  fixedContent,
	}

	err := inputWithWrongNumber.Validate()
	require.ErrorIs(t, err, myErrors.ErrInvalidNumber)
}

func TestNewSendFriendRequestUseCase_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	finderByID := mocks.NewMockUserFinderByID(ctrl)
	finderByID.EXPECT().FindByID(gomock.Any(), fixedFromID).Return(nil, nil).AnyTimes()

	finderByNumber := mocks.NewMockUserFinderByNumber(ctrl)
	finderByNumber.EXPECT().FindByNumber(gomock.Any(), fixedToNumber).Return(nil, nil).AnyTimes()

	saver := mocks.NewMockUserSaver(ctrl)
	saver.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	idGenerator := eventMocks.NewMockIDGenerator(ctrl)
	idGenerator.EXPECT().Generate().Return(fixedEventID).AnyTimes()

	operationIDGenerator := mocks.NewMockOperationIDGenerator(ctrl)
	operationIDGenerator.EXPECT().Generate().Return(fixedOperationID).AnyTimes()

	useCase, err := application.NewSendFriendRequestUseCase(
		finderByID,
		finderByNumber,
		saver,
		idGenerator,
		operationIDGenerator,
	)

	require.NoError(t, err)
	require.NotNil(t, useCase)
}

func TestNewSendFriendRequestUseCase_EmptyPointer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	useCase, err := application.NewSendFriendRequestUseCase(
		nil, nil, nil, nil, nil,
	)

	require.ErrorIs(t, err, myErrors.ErrEmptyPointer)
	require.Nil(t, useCase)
}

func TestSendFriendRequestUseCase_Execute_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	finderByID := mocks.NewMockUserFinderByID(ctrl)
	finderByNumber := mocks.NewMockUserFinderByNumber(ctrl)
	saver := mocks.NewMockUserSaver(ctrl)
	idGenerator := eventMocks.NewMockIDGenerator(ctrl)
	operationIDGenerator := mocks.NewMockOperationIDGenerator(ctrl)

	useCase, err := application.NewSendFriendRequestUseCase(
		finderByID,
		finderByNumber,
		saver,
		idGenerator,
		operationIDGenerator,
	)
	require.NoError(t, err)
	require.NotNil(t, useCase)

	fromUser := domain.LoadUser(fixedFromID, fixedFromNumber, make([]kernel.UserID, 0), make([]*domain.FriendRequest, 0))
	toUser := domain.LoadUser(fixedToID, fixedToNumber, make([]kernel.UserID, 0), make([]*domain.FriendRequest, 0))

	gomock.InOrder(
		finderByID.EXPECT().FindByID(gomock.Any(), fixedFromID).Return(fromUser, nil),
		finderByNumber.EXPECT().FindByNumber(gomock.Any(), fixedToNumber).Return(toUser, nil),
		operationIDGenerator.EXPECT().Generate().Return(fixedOperationID),
		idGenerator.EXPECT().Generate().Return(fixedEventID),
		saver.EXPECT().Save(gomock.Any(), fromUser).Return(nil),
		saver.EXPECT().Save(gomock.Any(), toUser).Return(nil),
	)

	_, err = useCase.Execute(nil, &application.SendFriendRequestInput{
		FromID:   fixedFromID,
		ToNumber: fixedToNumber,
		Content:  fixedContent,
	})
	require.NoError(t, err)
}
