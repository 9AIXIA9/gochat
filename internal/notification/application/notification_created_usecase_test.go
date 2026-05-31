package application_test

import (
	"encoding/json"
	"gochat/internal/notification/application"
	"gochat/internal/notification/domain"
	"gochat/internal/notification/domain/mocks"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const (
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
		SendAt:   time.Now().UTC(),
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

	mockService := mocks.NewMockDeliveryService(ctrl)

	useCase, err := application.NewNotificationCreatedUseCase(
		mockService,
	)

	require.NoError(t, err)
	require.NotNil(t, useCase)

	useCaseWithNil, err := application.NewNotificationCreatedUseCase(
		nil,
	)

	require.ErrorIs(t, err, myErrors.ErrEmptyPointer)
	require.Nil(t, useCaseWithNil)
}

func TestNotificationCreatedUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockDeliveryService(ctrl)

	useCase, err := application.NewNotificationCreatedUseCase(
		mockService,
	)
	require.NoError(t, err)
	require.NotNil(t, useCase)

	// 用户在线
	fixedRawPayload, err := json.Marshal(&AliasPayload{
		Content:  fixedContent,
		SenderID: fixedSenderID,
		SendAt:   time.Now().UTC(),
	})
	require.NoError(t, err)

	mockService.EXPECT().Deliver(nil, fixedRecipientID, domain.ActionPushNotification, fixedRawPayload).Return(nil)

	_, err = useCase.Execute(nil, &application.NotificationCreatedInput{
		ID:          fixedMessageID,
		RecipientID: fixedRecipientID,
		RawPayload:  fixedRawPayload,
	})
	require.NoError(t, err)
}
