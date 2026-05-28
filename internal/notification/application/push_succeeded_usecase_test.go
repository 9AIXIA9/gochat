package application_test

import (
	"gochat/internal/notification/application"
	"gochat/internal/notification/domain"
	"gochat/internal/notification/domain/mocks"
	myErrors "gochat/internal/shared/errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestPushSucceededInput_Validate(t *testing.T) {
	input := &application.PushSucceededInput{
		ID: fixedMessageID,
	}
	err := input.Validate()
	require.NoError(t, err)

	inputWithEmptyNotificationID := &application.PushSucceededInput{
		ID: "",
	}

	err = inputWithEmptyNotificationID.Validate()
	require.ErrorIs(t, err, myErrors.ErrEmptyInput)
}

func TestNewPushSucceededUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUpdater := mocks.NewMockNotificationStateUpdater(ctrl)

	useCase, err := application.NewPushSucceededUseCase(
		mockUpdater,
	)

	require.NoError(t, err)
	require.NotNil(t, useCase)

	useCaseWithNil, err := application.NewPushSucceededUseCase(
		nil,
	)

	require.ErrorIs(t, err, myErrors.ErrEmptyPointer)
	require.Nil(t, useCaseWithNil)
}

func TestPushSucceededUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUpdater := mocks.NewMockNotificationStateUpdater(ctrl)

	useCase, err := application.NewPushSucceededUseCase(
		mockUpdater,
	)
	require.NoError(t, err)
	require.NotNil(t, useCase)

	// 正常情况
	gomock.InOrder(
		mockUpdater.EXPECT().UpdateStateByID(gomock.Any(), fixedMessageID, domain.StateDelivered).Return(nil),
	)
	_, err = useCase.Execute(nil, &application.PushSucceededInput{
		ID: fixedMessageID,
	})
	require.NoError(t, err)
}
