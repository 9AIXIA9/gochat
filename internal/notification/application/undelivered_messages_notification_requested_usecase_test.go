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

const limit = 100

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

	mockMessages := make([]*domain.SystemMessage, 0, limit-1)
	for i := 0; i < limit-1; i++ {
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
		mockFinder.EXPECT().FindsByState(nil, fixedUserID, domain.MessageStateUndelivered, limit).Return(mockMessages, nil).Times(1),
		mockNotifier.EXPECT().Notify(gomock.Any()).Return(nil).Times(limit-1),
		mockUpdater.EXPECT().Updates(nil, mockMessages).Return(nil).Times(1),
	)

	_, err = useCase.Execute(nil, &application.UndeliveredMessagesNotificationRequestedInput{
		UserID: fixedUserID,
	})
	require.NoError(t, err)

	// 不需要通知的情况
	gomock.InOrder(
		mockFinder.EXPECT().FindsByState(nil, fixedUserID, domain.MessageStateUndelivered, limit).Return([]*domain.SystemMessage{}, nil).Times(1),
	)

	_, err = useCase.Execute(nil, &application.UndeliveredMessagesNotificationRequestedInput{
		UserID: fixedUserID,
	})
	require.NoError(t, err)

	//多次通知的情况
	mockMessages = make([]*domain.SystemMessage, 0, limit)
	for i := 0; i < limit; i++ {
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
		mockFinder.EXPECT().FindsByState(nil, fixedUserID, domain.MessageStateUndelivered, limit).Return(mockMessages, nil).Times(1),
		mockNotifier.EXPECT().Notify(gomock.Any()).Return(nil).Times(limit),
		mockUpdater.EXPECT().Updates(nil, mockMessages).Return(nil).Times(1),
		mockFinder.EXPECT().FindsByState(nil, fixedUserID, domain.MessageStateUndelivered, limit).Return(nil, nil).Times(1),
	)

	_, err = useCase.Execute(nil, &application.UndeliveredMessagesNotificationRequestedInput{
		UserID: fixedUserID,
	})
	require.NoError(t, err)
}
