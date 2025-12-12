package application_test

import (
	"gochat/internal/authorization/application"
	"gochat/internal/authorization/domain"
	"gochat/internal/authorization/domain/mocks"
	myErrors "gochat/internal/shared/errors"
	eventMock "gochat/internal/shared/event/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestUserCreatedInput_Validate(t *testing.T) {
	input := &application.UserCreatedInput{
		UserID: fixedUserID,
	}

	err := input.Validate()
	require.NoError(t, err)

	inputWithEmptyUserID := &application.UserCreatedInput{
		UserID: "",
	}

	err = inputWithEmptyUserID.Validate()
	require.ErrorIs(t, err, myErrors.ErrEmptyInput)
}

func TestNewUserCreatedUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIDGenerator := eventMock.NewMockIDGenerator(ctrl)
	mockCreator := eventMock.NewMockUnpublishedEventsCreator(ctrl)
	mockFinder := mocks.NewMockUserFinderByID(ctrl)

	useCase, err := application.NewUserCreatedUseCase(
		mockIDGenerator,
		mockCreator,
		mockFinder,
	)

	require.NoError(t, err)
	require.NotNil(t, useCase)

	useCaseWithNil, err := application.NewUserCreatedUseCase(
		nil, nil, nil,
	)

	require.ErrorIs(t, err, myErrors.ErrEmptyPointer)
	require.Nil(t, useCaseWithNil)
}

func TestUserCreatedUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIDGenerator := eventMock.NewMockIDGenerator(ctrl)
	mockCreator := eventMock.NewMockUnpublishedEventsCreator(ctrl)
	mockFinder := mocks.NewMockUserFinderByID(ctrl)

	useCase, err := application.NewUserCreatedUseCase(
		mockIDGenerator,
		mockCreator,
		mockFinder,
	)
	require.NoError(t, err)
	require.NotNil(t, useCase)

	mockUser := domain.LoadUser(
		fixedUserID,
		fixedEmail,
		fixedUserNumber,
		fixedEncryptedPassword,
		time.Now().UTC(),
	)

	gomock.InOrder(
		mockFinder.EXPECT().FindByID(nil, fixedUserID).Return(mockUser, nil).Times(1),
		mockIDGenerator.EXPECT().Generate().Return(fixedEventID).Times(5),
		mockCreator.EXPECT().CreateUnpublishedEvents(nil, gomock.Any()).Return(nil).Times(1),
	)

	_, err = useCase.Execute(nil, &application.UserCreatedInput{
		UserID: fixedUserID,
	})
	require.NoError(t, err)
}
