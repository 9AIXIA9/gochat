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

func TestUpdateUserProfileInput_Validate(t *testing.T) {
	input := &application.UpdateUserProfileInput{
		UserID:      fixedUserID,
		Name:        fixedName,
		Gender:      fixedGender,
		Email:       fixedEmail,
		PhoneNumber: fixedPhoneNumber,
		Address:     fixedAddress,
		Sign:        fixedSign,
	}

	err := input.Validate()
	require.NoError(t, err)

	inputWithEmptyUserID := &application.UpdateUserProfileInput{
		UserID:      "",
		Name:        fixedName,
		Gender:      fixedGender,
		Email:       fixedEmail,
		PhoneNumber: fixedPhoneNumber,
		Address:     fixedAddress,
		Sign:        fixedSign,
	}

	err = inputWithEmptyUserID.Validate()
	require.ErrorIs(t, err, myErrors.ErrEmptyInput)

	inputWithOnlyUserID := &application.UpdateUserProfileInput{
		UserID: fixedUserID,
	}

	err = inputWithOnlyUserID.Validate()
	require.ErrorIs(t, err, myErrors.ErrEmptyInput)
}

func TestNewUpdateUserProfileUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	updater := mocks.NewMockUserProfileUpdater(ctrl)
	finder := mocks.NewMockUserProfileFinder(ctrl)

	useCase, err := application.NewUpdateUserProfileUseCase(
		updater,
		finder,
	)

	require.NoError(t, err)
	require.NotNil(t, useCase)

	useCaseWithNil, err := application.NewUpdateUserProfileUseCase(
		nil, nil,
	)

	require.ErrorIs(t, err, myErrors.ErrEmptyPointer)
	require.Nil(t, useCaseWithNil)
}

func TestUpdateUserProfileUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	updater := mocks.NewMockUserProfileUpdater(ctrl)
	finder := mocks.NewMockUserProfileFinder(ctrl)

	useCase, err := application.NewUpdateUserProfileUseCase(
		updater,
		finder,
	)
	require.NoError(t, err)
	require.NotNil(t, useCase)

	// 正常情况
	gomock.InOrder(
		finder.EXPECT().FindByID(nil, fixedUserID).Return(domain.LoadUserProfile(
			fixedUserID,
			fixedName,
			fixedGender,
			fixedEmail,
			fixedPhoneNumber,
			fixedAddress,
			fixedSign,
			time.Now().UTC(),
		), nil),
		updater.EXPECT().Update(nil, gomock.Any()).Return(nil),
	)

	_, err = useCase.Execute(nil, &application.UpdateUserProfileInput{
		UserID:      fixedUserID,
		Name:        fixedName,
		Gender:      fixedGender,
		Email:       fixedEmail,
		PhoneNumber: fixedPhoneNumber,
		Address:     fixedAddress,
		Sign:        fixedSign,
	})
	require.NoError(t, err)
}
