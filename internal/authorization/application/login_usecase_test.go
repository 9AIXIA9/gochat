package application_test

import (
	"gochat/internal/authorization/application"
	"gochat/internal/authorization/domain"
	"gochat/internal/authorization/domain/mocks"
	myErrors "gochat/internal/shared/errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestLoginInput_Validate(t *testing.T) {
	input := &application.LoginInput{
		Number:   fixedUserNumber,
		Password: fixedPassword,
	}

	err := input.Validate()
	require.NoError(t, err)

	inputWithEmptyPassword := &application.LoginInput{
		Number:   fixedUserNumber,
		Password: "",
	}

	err = inputWithEmptyPassword.Validate()
	require.Error(t, err)

	inputWithEmptyNumber := &application.LoginInput{
		Number:   "",
		Password: fixedPassword,
	}

	err = inputWithEmptyNumber.Validate()
	require.ErrorIs(t, err, myErrors.ErrEmptyInput)
}

func TestNewLoginUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockComparator := mocks.NewMockComparator(ctrl)
	mockUserFinder := mocks.NewMockUserFinderByNumber(ctrl)
	mockRefreshTokenUpserter := mocks.NewMockRefreshTokenUpserter(ctrl)
	mockAccessTokenGenerator := mocks.NewMockAccessTokenGenerator(ctrl)
	mockRefreshTokenGenerator := mocks.NewMockRefreshTokenGenerator(ctrl)

	useCase, err := application.NewLoginUseCase(
		mockComparator,
		mockUserFinder,
		mockRefreshTokenUpserter,
		mockAccessTokenGenerator,
		mockRefreshTokenGenerator,
	)

	require.NoError(t, err)
	require.NotNil(t, useCase)

	useCaseWithNil, err := application.NewLoginUseCase(
		nil, nil, nil, nil, nil,
	)

	require.ErrorIs(t, err, myErrors.ErrEmptyPointer)
	require.Nil(t, useCaseWithNil)
}

func TestLoginUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockComparator := mocks.NewMockComparator(ctrl)
	mockUserFinder := mocks.NewMockUserFinderByNumber(ctrl)
	mockRefreshTokenUpserter := mocks.NewMockRefreshTokenUpserter(ctrl)
	mockAccessTokenGenerator := mocks.NewMockAccessTokenGenerator(ctrl)
	mockRefreshTokenGenerator := mocks.NewMockRefreshTokenGenerator(ctrl)

	useCase, err := application.NewLoginUseCase(
		mockComparator,
		mockUserFinder,
		mockRefreshTokenUpserter,
		mockAccessTokenGenerator,
		mockRefreshTokenGenerator,
	)
	require.NoError(t, err)
	require.NotNil(t, useCase)

	mockUser := domain.LoadUser(
		fixedUserID,
		fixedEmail,
		fixedUserNumber,
		fixedEncryptedPassword,
		time.Now().UTC(),
	)

	gomock.InOrder(
		mockUserFinder.EXPECT().FindByNumber(nil, fixedUserNumber).Return(mockUser, nil).Times(1),
		mockComparator.EXPECT().Compare(fixedEncryptedPassword.String(), fixedPassword.String()).Return(nil).Times(1),
		mockRefreshTokenGenerator.EXPECT().Generate().Return(fixedRefreshToken, nil).Times(1),
		mockAccessTokenGenerator.EXPECT().Generate(fixedUserID).Return(fixedAccessToken, nil).Times(1),
		mockRefreshTokenUpserter.EXPECT().Upsert(nil, gomock.Any()).Return(nil).Times(1),
	)

	_, err = useCase.Execute(nil, &application.LoginInput{
		Number:   fixedUserNumber,
		Password: fixedPassword,
	})
	require.NoError(t, err)
}
