package application_test

import (
	"encoding/json"
	"gochat/internal/notification/application"
	"gochat/internal/notification/domain/mocks"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	eventMock "gochat/internal/shared/event/mocks"
	"gochat/internal/shared/kernel"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const (
	fixedEventID     event.ID         = "event-123"
	fixedMessageID   kernel.MessageID = "notification-123"
	fixedRecipientID kernel.UserID    = "user-456"
	fixedSenderID    kernel.UserID    = "user-789"
	fixedContent     string           = "Hello, this is a notification!"
)

type AliasPayload struct {
	Content  string
	SenderID kernel.UserID
	SendAt   time.Time
}

func TestNotificationCreatedInput_Validate(t *testing.T) {
	fixedRawPayload, err := json.Marshal(&AliasPayload{
		Content:  fixedContent,
		SenderID: fixedSenderID,
		SendAt:   time.Now(),
	})
	require.NoError(t, err)

	input := &application.NotificationCreatedInput{
		ID:          fixedMessageID,
		RecipientID: fixedRecipientID,
		RawPayload:  fixedRawPayload,
	}

	err = input.Validate()
	require.NoError(t, err)

	inputWithEmptyNotificationID := &application.NotificationCreatedInput{
		ID: "",
	}

	err = inputWithEmptyNotificationID.Validate()
	require.ErrorIs(t, err, myErrors.ErrEmptyInput)

	inputWithEmptyRecipientID := &application.NotificationCreatedInput{
		ID: fixedMessageID,
	}

	err = inputWithEmptyRecipientID.Validate()
	require.ErrorIs(t, err, myErrors.ErrEmptyInput)
}

func TestNewNotificationCreatedUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIDGenerator := eventMock.NewMockIDGenerator(ctrl)
	mockCreator := mocks.NewMockNotificationCreator(ctrl)

	useCase, err := application.NewNotificationCreatedUseCase(
		mockIDGenerator,
		mockCreator,
	)

	require.NoError(t, err)
	require.NotNil(t, useCase)

	useCaseWithNil, err := application.NewNotificationCreatedUseCase(
		nil, nil,
	)

	require.ErrorIs(t, err, myErrors.ErrEmptyPointer)
	require.Nil(t, useCaseWithNil)
}

func TestNotificationCreatedUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIDGenerator := eventMock.NewMockIDGenerator(ctrl)
	mockCreator := mocks.NewMockNotificationCreator(ctrl)

	useCase, err := application.NewNotificationCreatedUseCase(
		mockIDGenerator, mockCreator,
	)
	require.NoError(t, err)
	require.NotNil(t, useCase)

	// 正常情况
	fixedRawPayload, err := json.Marshal(&AliasPayload{
		Content:  fixedContent,
		SenderID: fixedSenderID,
		SendAt:   time.Now(),
	})
	require.NoError(t, err)

	gomock.InOrder(
		mockIDGenerator.EXPECT().Generate().Return(fixedEventID).Times(1),
		mockCreator.EXPECT().Create(nil, gomock.Any()).Return(nil).Times(1),
	)
	_, err = useCase.Execute(nil, &application.NotificationCreatedInput{
		ID:          fixedMessageID,
		RecipientID: fixedRecipientID,
		RawPayload:  fixedRawPayload,
	})
	require.NoError(t, err)
}
