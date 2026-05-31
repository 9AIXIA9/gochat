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

func TestLoginByEmailInput_Validate(t *testing.T) {
	input := &application.LoginByEmailInput{
		Email:    fixedEmail,
		Password: fixedPassword,
	}

	err := input.Validate()
	require.NoError(t, err)

	inputWithEmptyPassword := &application.LoginByEmailInput{
		Email:    fixedEmail,
		Password: "",
	}

	err = inputWithEmptyPassword.Validate()
	require.ErrorIs(t, err, myErrors.ErrInvalidLength)

	inputWithEmptyEmail := &application.LoginByEmailInput{
		Email:    "",
		Password: fixedPassword,
	}

	err = inputWithEmptyEmail.Validate()
	require.ErrorIs(t, err, myErrors.ErrInvalidFormat)
}

func TestNewLoginByEmailUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockComparator := mocks.NewMockComparator(ctrl)
	mockUserFinder := mocks.NewMockUserFinderByEmail(ctrl)
	mockRefreshTokenUpserter := mocks.NewMockRefreshTokenUpserter(ctrl)
	mockAccessTokenGenerator := mocks.NewMockAccessTokenGenerator(ctrl)
	mockRefreshTokenGenerator := mocks.NewMockRefreshTokenGenerator(ctrl)

	useCase, err := application.NewLoginByEmailUseCase(
		mockComparator,
		mockUserFinder,
		mockRefreshTokenUpserter,
		mockAccessTokenGenerator,
		mockRefreshTokenGenerator,
	)

	require.NoError(t, err)
	require.NotNil(t, useCase)

	useCaseWithNil, err := application.NewLoginByEmailUseCase(
		nil, nil, nil, nil, nil,
	)

	require.ErrorIs(t, err, myErrors.ErrEmptyPointer)
	require.Nil(t, useCaseWithNil)
}

func TestLoginByEmailUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockComparator := mocks.NewMockComparator(ctrl)
	mockUserFinder := mocks.NewMockUserFinderByEmail(ctrl)
	mockRefreshTokenUpserter := mocks.NewMockRefreshTokenUpserter(ctrl)
	mockAccessTokenGenerator := mocks.NewMockAccessTokenGenerator(ctrl)
	mockRefreshTokenGenerator := mocks.NewMockRefreshTokenGenerator(ctrl)

	useCase, err := application.NewLoginByEmailUseCase(
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
		mockUserFinder.EXPECT().FindByEmail(nil, fixedEmail).Return(mockUser, nil).Times(1),
		mockComparator.EXPECT().Compare(fixedEncryptedPassword.String(), fixedPassword.String()).Return(nil).Times(1),
		mockRefreshTokenGenerator.EXPECT().Generate().Return(fixedRefreshToken, nil).Times(1),
		mockAccessTokenGenerator.EXPECT().Generate(fixedUserID).Return(fixedAccessToken, nil).Times(1),
		mockRefreshTokenUpserter.EXPECT().Upsert(nil, gomock.Any()).Return(nil).Times(1),
	)

	_, err = useCase.Execute(nil, &application.LoginByEmailInput{
		Email:    fixedEmail,
		Password: fixedPassword,
	})
	require.NoError(t, err)

	// 用户不存在
	mockUserFinder.EXPECT().FindByEmail(nil, fixedEmail).Return(nil, myErrors.ErrNotFound).Times(1)
	_, err = useCase.Execute(nil, &application.LoginByEmailInput{
		Email:    fixedEmail,
		Password: fixedPassword,
	})
	require.ErrorIs(t, err, domain.ErrInvalidPassword)
}
