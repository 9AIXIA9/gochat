package application_test

import (
	"gochat/internal/friendship/application"
	"gochat/internal/friendship/domain"
	"gochat/internal/friendship/domain/mocks"
	myErrors "gochat/internal/shared/errors"
	eventMocks "gochat/internal/shared/event/mocks"
	"gochat/internal/shared/kernel"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestRefuseFriendRequestInput_Validate_Success(t *testing.T) {
	input := &application.RefuseFriendRequestInput{
		UserID:    fixedFromID,
		RequestID: fixedOperationID,
	}

	err := input.Validate()
	require.NoError(t, err)
}

func TestRefuseFriendRequestInput_Validate_EmptyInput(t *testing.T) {
	inputWithEmptyUserID := &application.RefuseFriendRequestInput{
		UserID:    "",
		RequestID: fixedOperationID,
	}

	err := inputWithEmptyUserID.Validate()
	require.ErrorIs(t, err, myErrors.ErrEmptyInput)

	inputWithEmptyRequestID := &application.RefuseFriendRequestInput{
		UserID:    fixedFromID,
		RequestID: "",
	}

	err = inputWithEmptyRequestID.Validate()
	require.ErrorIs(t, err, myErrors.ErrEmptyInput)
}

func TestNewRefuseFriendRequestUseCase_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	finder := mocks.NewMockUserFinderByID(ctrl)
	finder.EXPECT().FindByID(gomock.Any(), fixedFromID).Return(nil, nil).AnyTimes()

	saver := mocks.NewMockUserSaver(ctrl)
	saver.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	idGenerator := eventMocks.NewMockIDGenerator(ctrl)
	idGenerator.EXPECT().Generate().Return(fixedEventID).AnyTimes()

	useCase, err := application.NewRefuseFriendRequestUseCase(
		finder,
		saver,
		idGenerator,
	)

	require.NoError(t, err)
	require.NotNil(t, useCase)
}

func TestNewRefuseFriendRequestUseCase_EmptyPointer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	useCase, err := application.NewRefuseFriendRequestUseCase(
		nil, nil, nil,
	)

	require.ErrorIs(t, err, myErrors.ErrEmptyPointer)
	require.Nil(t, useCase)
}

func TestRefuseFriendRequestUseCase_Execute_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	finder := mocks.NewMockUserFinderByID(ctrl)
	saver := mocks.NewMockUserSaver(ctrl)
	idGenerator := eventMocks.NewMockIDGenerator(ctrl)

	useCase, err := application.NewRefuseFriendRequestUseCase(
		finder,
		saver,
		idGenerator,
	)
	require.NoError(t, err)
	require.NotNil(t, useCase)

	fromUser := domain.LoadUser(fixedFromID, fixedFromNumber, make([]kernel.UserID, 0), make([]*domain.FriendRequest, 0))

	gomock.InOrder(
		finder.EXPECT().FindByID(gomock.Any(), fixedFromID).Return(fromUser, nil),
		saver.EXPECT().Save(gomock.Any(), fromUser).Return(nil),
	)

	_, err = useCase.Execute(nil, &application.RefuseFriendRequestInput{
		UserID:    fixedFromID,
		RequestID: fixedOperationID,
	})
	require.NoError(t, err)
}
