package application_test

import (
	"gochat/internal/notification/application"
	"gochat/internal/notification/domain/mocks"
	myErrors "gochat/internal/shared/errors"
	kernelmocks "gochat/internal/shared/kernel/mocks"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestSystemMessageNotificationRequestedInput_Validate(t *testing.T) {
	input := &application.SystemMessageNotificationRequestedInput{
		RecipientID: fixedUserID,
		Content:     fixedContent,
	}

	err := input.Validate()
	require.NoError(t, err)

	inputWithEmptyRecipientID := &application.SystemMessageNotificationRequestedInput{
		RecipientID: "",
		Content:     fixedContent,
	}

	err = inputWithEmptyRecipientID.Validate()
	require.ErrorIs(t, err, myErrors.ErrEmptyInput)

	inputWithEmptyContent := &application.SystemMessageNotificationRequestedInput{
		RecipientID: fixedUserID,
		Content:     "",
	}
	err = inputWithEmptyContent.Validate()
	require.ErrorIs(t, err, myErrors.ErrEmptyInput)
}

func TestNewSystemMessageNotificationRequestedUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMessageCreator := mocks.NewMockSystemMessageCreator(ctrl)
	mockMessageNotifier := mocks.NewMockSystemMessageNotifier(ctrl)
	mockIDGenerator := kernelmocks.NewMockMessageIDGenerator(ctrl)

	useCase, err := application.NewSystemMessageNotificationRequestedUseCase(
		mockMessageCreator,
		mockMessageNotifier,
		mockIDGenerator,
	)

	require.NoError(t, err)
	require.NotNil(t, useCase)

	useCaseWithNil, err := application.NewSystemMessageNotificationRequestedUseCase(
		nil, nil, nil,
	)

	require.ErrorIs(t, err, myErrors.ErrEmptyPointer)
	require.Nil(t, useCaseWithNil)
}

func TestSystemMessageNotificationRequestedUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMessageCreator := mocks.NewMockSystemMessageCreator(ctrl)
	mockMessageNotifier := mocks.NewMockSystemMessageNotifier(ctrl)
	mockIDGenerator := kernelmocks.NewMockMessageIDGenerator(ctrl)

	useCase, err := application.NewSystemMessageNotificationRequestedUseCase(
		mockMessageCreator,
		mockMessageNotifier,
		mockIDGenerator,
	)
	require.NoError(t, err)
	require.NotNil(t, useCase)

	gomock.InOrder(
		mockIDGenerator.EXPECT().Generate().Return(fixedMessageID).Times(1),
		mockMessageNotifier.EXPECT().Notify(gomock.Any()).Times(1),
		mockMessageCreator.EXPECT().Create(nil, gomock.Any()).Return(nil).Times(1),
	)

	_, err = useCase.Execute(nil, &application.SystemMessageNotificationRequestedInput{
		RecipientID: fixedUserID,
		Content:     fixedContent,
	})
	require.NoError(t, err)
}
