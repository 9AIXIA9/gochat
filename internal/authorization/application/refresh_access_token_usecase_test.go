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

func TestRefreshAccessTokenInput_Validate(t *testing.T) {
	input := &application.RefreshAccessTokenInput{
		RefreshToken: fixedRefreshToken,
	}

	err := input.Validate()
	require.NoError(t, err)

	inputWithEmptyRefreshToken := &application.RefreshAccessTokenInput{
		RefreshToken: "",
	}

	err = inputWithEmptyRefreshToken.Validate()
	require.ErrorIs(t, err, myErrors.ErrEmptyInput)
}

func TestNewRefreshAccessTokenUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRefreshTokenUpserter := mocks.NewMockRefreshTokenUpserter(ctrl)
	mockRefreshTokenFinder := mocks.NewMockRefreshTokenFinder(ctrl)
	mockAccessTokenGenerator := mocks.NewMockAccessTokenGenerator(ctrl)
	mockRefreshTokenGenerator := mocks.NewMockRefreshTokenGenerator(ctrl)

	useCase, err := application.NewRefreshAccessTokenUseCase(
		mockRefreshTokenUpserter,
		mockRefreshTokenFinder,
		mockAccessTokenGenerator,
		mockRefreshTokenGenerator,
	)

	require.NoError(t, err)
	require.NotNil(t, useCase)

	useCaseWithNil, err := application.NewRefreshAccessTokenUseCase(
		nil, nil, nil, nil,
	)

	require.ErrorIs(t, err, myErrors.ErrEmptyPointer)
	require.Nil(t, useCaseWithNil)
}

func TestRefreshAccessTokenUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRefreshTokenUpserter := mocks.NewMockRefreshTokenUpserter(ctrl)
	mockRefreshTokenFinder := mocks.NewMockRefreshTokenFinder(ctrl)
	mockAccessTokenGenerator := mocks.NewMockAccessTokenGenerator(ctrl)
	mockRefreshTokenGenerator := mocks.NewMockRefreshTokenGenerator(ctrl)

	useCase, err := application.NewRefreshAccessTokenUseCase(
		mockRefreshTokenUpserter,
		mockRefreshTokenFinder,
		mockAccessTokenGenerator,
		mockRefreshTokenGenerator,
	)
	require.NoError(t, err)
	require.NotNil(t, useCase)

	mockRefreshTokenEntity := domain.LoadRefreshToken(
		fixedRefreshToken,
		fixedUserID,
		time.Now().UTC().Add(time.Hour),
		fixedRefreshCount,
	)

	gomock.InOrder(
		mockRefreshTokenFinder.EXPECT().FindByToken(nil, fixedRefreshToken).Return(mockRefreshTokenEntity, nil).Times(1),
		mockRefreshTokenGenerator.EXPECT().Generate().Return(fixedRefreshToken, nil).Times(1),
		mockAccessTokenGenerator.EXPECT().Generate(fixedUserID).Return(fixedAccessToken, nil).Times(1),
		mockRefreshTokenUpserter.EXPECT().Upsert(nil, gomock.Any()).Return(nil).Times(1),
	)

	_, err = useCase.Execute(nil, &application.RefreshAccessTokenInput{
		RefreshToken: fixedRefreshToken,
	})
	require.NoError(t, err)

	// refresh token 不存在
	mockRefreshTokenFinder.EXPECT().FindByToken(nil, fixedRefreshToken).Return(nil, myErrors.ErrNotFound).Times(1)
	_, err = useCase.Execute(nil, &application.RefreshAccessTokenInput{
		RefreshToken: fixedRefreshToken,
	})
	require.ErrorIs(t, err, domain.ErrInvalidRefreshToken)
}
