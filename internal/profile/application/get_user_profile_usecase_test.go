package application_test

import (
	"gochat/internal/profile/application"
	"gochat/internal/profile/domain"
	"gochat/internal/profile/domain/mocks"
	myErrors "gochat/internal/shared/errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetUserProfileInput_Validate(t *testing.T) {
	input := &application.GetUserProfileInput{
		UserID: fixedUserID,
	}

	err := input.Validate()
	require.NoError(t, err)

	inputWithEmptyUserID := &application.GetUserProfileInput{
		UserID: "",
	}

	err = inputWithEmptyUserID.Validate()
	require.ErrorIs(t, err, myErrors.ErrEmptyInput)
}

func TestNewGetUserProfileUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProfileFinder := mocks.NewMockUserProfileFinder(ctrl)

	useCase, err := application.NewGetUserProfileUseCase(
		mockProfileFinder,
	)

	require.NoError(t, err)
	require.NotNil(t, useCase)

	useCaseWithNil, err := application.NewGetUserProfileUseCase(
		nil,
	)

	require.ErrorIs(t, err, myErrors.ErrEmptyPointer)
	require.Nil(t, useCaseWithNil)
}

func TestGetUserProfileUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProfileFinder := mocks.NewMockUserProfileFinder(ctrl)

	useCase, err := application.NewGetUserProfileUseCase(
		mockProfileFinder,
	)

	require.NoError(t, err)
	require.NotNil(t, useCase)

	// 正常情况
	gomock.InOrder(
		mockProfileFinder.EXPECT().FindByID(nil, fixedUserID).Return(domain.LoadUserProfile(
			fixedUserID,
			fixedName,
			fixedGender,
			fixedEmail,
			fixedPhoneNumber,
			fixedAddress,
			fixedSign,
			time.Now().UTC(),
		), nil).Times(1),
	)

	_, err = useCase.Execute(nil, &application.GetUserProfileInput{
		UserID: fixedUserID,
	})
	require.NoError(t, err)
}
