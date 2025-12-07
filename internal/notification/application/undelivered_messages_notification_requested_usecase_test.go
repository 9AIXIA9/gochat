package application_test

import (
	"gochat/internal/notification/application"
	"gochat/internal/notification/domain"
	"gochat/internal/notification/domain/mocks"
	myErrors "gochat/internal/shared/errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const fixedMessageLen = 10

func TestUndeliveredMessagesNotificationRequestedInput_Validate(t *testing.T) {
	input := &application.UndeliveredMessagesNotificationRequestedInput{
		UserID: fixedUserID,
	}

	err := input.Validate()
	require.NoError(t, err)

	inputWithEmptyUserID := &application.UndeliveredMessagesNotificationRequestedInput{
		UserID: "",
	}

	err = inputWithEmptyUserID.Validate()
	require.ErrorIs(t, err, myErrors.ErrEmptyInput)
}

func TestNewUndeliveredMessagesNotificationRequestedUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFinder := mocks.NewMockUserSystemMessagesFinderByState(ctrl)
	mockUpdater := mocks.NewMockSystemMessagesUpdater(ctrl)
	mockNotifier := mocks.NewMockSystemMessageNotifier(ctrl)

	useCase, err := application.NewUndeliveredMessagesNotificationRequestedUseCase(
		mockFinder, mockUpdater, mockNotifier,
	)

	require.NoError(t, err)
	require.NotNil(t, useCase)

	useCaseWithNil, err := application.NewUndeliveredMessagesNotificationRequestedUseCase(
		nil, nil, nil,
	)

	require.ErrorIs(t, err, myErrors.ErrEmptyPointer)
	require.Nil(t, useCaseWithNil)
}

func TestUndeliveredMessagesNotificationRequestedUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFinder := mocks.NewMockUserSystemMessagesFinderByState(ctrl)
	mockUpdater := mocks.NewMockSystemMessagesUpdater(ctrl)
	mockNotifier := mocks.NewMockSystemMessageNotifier(ctrl)

	useCase, err := application.NewUndeliveredMessagesNotificationRequestedUseCase(
		mockFinder, mockUpdater, mockNotifier,
	)
	require.NoError(t, err)
	require.NotNil(t, useCase)

	mockMessages := make([]*domain.SystemMessage, 0, fixedMessageLen)
	for i := 0; i < fixedMessageLen; i++ {
		message := domain.LoadSystemMessage(
			fixedMessageID,
			fixedUserID,
			domain.MessageStateUndelivered,
			"Test content",
			time.Now().UTC(),
		)
		mockMessages = append(mockMessages, message)
	}

	gomock.InOrder(
		mockFinder.EXPECT().FindsByState(nil, fixedUserID, domain.MessageStateUndelivered).Return(mockMessages, nil),
		mockNotifier.EXPECT().Notify(gomock.Any()).Times(fixedMessageLen).Return(nil),
		mockUpdater.EXPECT().Updates(nil, mockMessages).Return(nil),
	)

	_, err = useCase.Execute(nil, &application.UndeliveredMessagesNotificationRequestedInput{
		UserID: fixedUserID,
	})
	require.NoError(t, err)
}
