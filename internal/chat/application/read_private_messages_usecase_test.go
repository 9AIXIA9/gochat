package application_test

import (
	"gochat/internal/chat/application"
	"gochat/internal/chat/domain"
	"gochat/internal/chat/domain/mocks"
	myErrors "gochat/internal/shared/errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestReadPrivateMessagesInput_Validate(t *testing.T) {
	input := &application.ReadPrivateMessagesInput{
		SenderID:    fixedFriendID,
		RecipientID: fixedUserID,
	}

	err := input.Validate()
	require.NoError(t, err)

	inputWithEmptySenderID := &application.ReadPrivateMessagesInput{
		SenderID:    "",
		RecipientID: fixedUserID,
	}

	err = inputWithEmptySenderID.Validate()
	require.ErrorIs(t, err, myErrors.ErrEmptyInput)

	inputWithEmptyRecipientID := &application.ReadPrivateMessagesInput{
		SenderID:    fixedFriendID,
		RecipientID: "",
	}

	err = inputWithEmptyRecipientID.Validate()
	require.ErrorIs(t, err, myErrors.ErrEmptyInput)
}

func TestNewReadPrivateMessagesUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUpdater := mocks.NewMockPrivateMessagesStatesUpdaterByUserID(ctrl)

	useCase, err := application.NewReadPrivateMessagesUseCase(
		mockUpdater,
	)

	require.NoError(t, err)
	require.NotNil(t, useCase)

	useCaseWithNil, err := application.NewReadPrivateMessagesUseCase(
		nil,
	)

	require.ErrorIs(t, err, myErrors.ErrEmptyPointer)
	require.Nil(t, useCaseWithNil)
}

func TestReadPrivateMessagesUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUpdater := mocks.NewMockPrivateMessagesStatesUpdaterByUserID(ctrl)

	useCase, err := application.NewReadPrivateMessagesUseCase(
		mockUpdater,
	)
	require.NoError(t, err)
	require.NotNil(t, useCase)

	// 正常情况
	mockUpdater.EXPECT().UpdatesByUserID(nil, fixedFriendID, fixedUserID, domain.MessageStateRead).Return(nil).Times(1)
	_, err = useCase.Execute(nil, &application.ReadPrivateMessagesInput{
		SenderID:    fixedFriendID,
		RecipientID: fixedUserID,
	})
	require.NoError(t, err)
}
