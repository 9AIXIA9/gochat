package application_test

import (
	"gochat/internal/chat/application"
	"gochat/internal/chat/domain"
	"gochat/internal/chat/domain/mocks"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestConfirmRoomMessagesInput_Validate(t *testing.T) {
	input := &application.ConfirmRoomMessagesInput{
		UserID:     fixedUserID,
		MessageIDs: []kernel.MessageID{fixedMessageID},
	}

	err := input.Validate()
	require.NoError(t, err)

	inputWithEmpty := &application.ConfirmRoomMessagesInput{
		UserID:     fixedUserID,
		MessageIDs: nil,
	}

	err = inputWithEmpty.Validate()
	require.ErrorIs(t, err, myErrors.ErrEmptyInput)
}

func TestNewConfirmRoomMessagesUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockRoomMessageRepository(ctrl)

	useCase, err := application.NewConfirmRoomMessagesUseCase(
		mockRepo,
	)
	require.NoError(t, err)
	require.NotNil(t, useCase)

	useCaseWithNil, err := application.NewConfirmRoomMessagesUseCase(
		nil,
	)

	require.ErrorIs(t, err, myErrors.ErrEmptyPointer)
	require.Nil(t, useCaseWithNil)
}

func TestConfirmRoomMessagesUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockRoomMessageRepository(ctrl)

	useCase, err := application.NewConfirmRoomMessagesUseCase(
		mockRepo,
	)
	require.NoError(t, err)
	require.NotNil(t, useCase)

	mockRepo.EXPECT().UpdatesByMessageIDs(nil, fixedUserID, []kernel.MessageID{fixedMessageID}, domain.MessageStateDelivered).Return(nil).Times(1)
	_, err = useCase.Execute(nil, &application.ConfirmRoomMessagesInput{
		UserID:     fixedUserID,
		MessageIDs: []kernel.MessageID{fixedMessageID},
	})
	require.NoError(t, err)
}
