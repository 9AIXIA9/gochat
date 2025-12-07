package application_test

import (
	"gochat/internal/application"
	myErrors "gochat/internal/shared/errors"
	eventMock "gochat/internal/shared/event/mocks"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestUserSessionStartedInput_Validate(t *testing.T) {
	input := &application.UserSessionStartedInput{
		UserID: fixedUserID,
	}

	err := input.Validate()
	require.NoError(t, err)

	inputWithEmptyUserID := &application.UserSessionStartedInput{
		UserID: "",
	}

	err = inputWithEmptyUserID.Validate()
	require.ErrorIs(t, err, myErrors.ErrEmptyInput)
}

func TestNewUserSessionStartedUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIDGenerator := eventMock.NewMockIDGenerator(ctrl)
	mockCreator := eventMock.NewMockUnpublishedEventsCreator(ctrl)

	useCase, err := application.NewUserSessionStartedUseCase(
		mockIDGenerator,
		mockCreator,
	)

	require.NoError(t, err)
	require.NotNil(t, useCase)

	useCaseWithNil, err := application.NewUserSessionStartedUseCase(
		nil, nil,
	)

	require.ErrorIs(t, err, myErrors.ErrEmptyPointer)
	require.Nil(t, useCaseWithNil)
}

func TestUserSessionStartedUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIDGenerator := eventMock.NewMockIDGenerator(ctrl)
	mockCreator := eventMock.NewMockUnpublishedEventsCreator(ctrl)

	useCase, err := application.NewUserSessionStartedUseCase(
		mockIDGenerator,
		mockCreator,
	)
	require.NoError(t, err)
	require.NotNil(t, useCase)

	gomock.InOrder(
		mockIDGenerator.EXPECT().Generate().Return(fixedEventID).Times(2),
		mockCreator.EXPECT().CreateUnpublishedEvents(
			gomock.Any(),
			gomock.Any(),
		).Return(nil),
	)

	_, err = useCase.Execute(nil, &application.UserSessionStartedInput{
		UserID: fixedUserID,
	})
	require.NoError(t, err)
}
