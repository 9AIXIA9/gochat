package domain_test

import (
	"gochat/internal/authorization/domain"
	"gochat/internal/authorization/domain/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestLoadRefreshToken(t *testing.T) {
	token := domain.LoadRefreshToken(
		fixedRefreshToken,
		fixedUserID,
		time.Now().UTC(),
		fixedRefreshCount,
	)
	require.NotNil(t, token)
	assert.Equal(t, fixedRefreshToken, token.Token())
	assert.Equal(t, fixedUserID, token.UserID())
	assert.Equal(t, fixedRefreshCount, token.RefreshCount())
	assert.WithinDuration(t, time.Now().UTC(), token.ExpiredAt(), timeTolerance)
}

func TestCreateRefreshToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGenerator := mocks.NewMockRefreshTokenGenerator(ctrl)

	mockGenerator.EXPECT().Generate().Return(fixedRefreshToken, nil).Times(1)
	token, err := domain.CreateRefreshToken(
		fixedUserID,
		mockGenerator,
	)
	require.NoError(t, err)
	require.NotNil(t, token)
	assert.Equal(t, fixedRefreshToken, token.Token())
	assert.Equal(t, fixedUserID, token.UserID())
	assert.Equal(t, 0, token.RefreshCount())
	assert.WithinDuration(t, time.Now().UTC().Add(validityDuration), token.ExpiredAt(), timeTolerance)
}

func TestRefreshTokenEntity_GenerateAccessToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRefreshTokenGenerator := mocks.NewMockRefreshTokenGenerator(ctrl)
	mockAccessTokenGenerator := mocks.NewMockAccessTokenGenerator(ctrl)

	//正常情况
	mockRefreshTokenGenerator.EXPECT().Generate().Return(fixedRefreshToken, nil).Times(1)
	refreshToken, err := domain.CreateRefreshToken(
		fixedUserID,
		mockRefreshTokenGenerator,
	)
	require.NoError(t, err)
	require.NotNil(t, refreshToken)

	mockAccessTokenGenerator.EXPECT().Generate(fixedUserID).Return(fixedAccessToken, nil).Times(1)
	accessToken, err := refreshToken.GenerateAccessToken(mockAccessTokenGenerator)
	require.NoError(t, err)
	assert.Equal(t, fixedAccessToken, accessToken)
}

func TestRefreshTokenEntity_Refresh(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRefreshTokenGenerator := mocks.NewMockRefreshTokenGenerator(ctrl)

	//正常情况
	mockRefreshTokenGenerator.EXPECT().Generate().Return(fixedRefreshToken, nil).Times(1)
	refreshToken, err := domain.CreateRefreshToken(
		fixedUserID,
		mockRefreshTokenGenerator,
	)
	require.NoError(t, err)
	require.NotNil(t, refreshToken)

	mockRefreshTokenGenerator.EXPECT().Generate().Return(fixedRefreshToken, nil).Times(1)
	err = refreshToken.Refresh(mockRefreshTokenGenerator)
	require.NoError(t, err)
	assert.Equal(t, fixedRefreshToken, refreshToken.Token())
	assert.Equal(t, 1, refreshToken.RefreshCount())
	assert.WithinDuration(t, time.Now().UTC().Add(validityDuration).Add(extendDuration), refreshToken.ExpiredAt(), timeTolerance)

	//过期失败
	expiredToken := domain.LoadRefreshToken(
		fixedRefreshToken,
		fixedUserID,
		time.Now().UTC().Add(-time.Hour),
		fixedRefreshCount,
	)
	err = expiredToken.Refresh(mockRefreshTokenGenerator)
	require.Error(t, err)

	//超过刷新次数失败
	limitToken := domain.LoadRefreshToken(
		fixedRefreshToken,
		fixedUserID,
		time.Now().UTC(),
		maxRefreshCount,
	)
	err = limitToken.Refresh(mockRefreshTokenGenerator)
	require.Error(t, err)
}

func TestRefreshTokenEntity_CanBeRefreshed(t *testing.T) {
	//正常情况
	token := domain.LoadRefreshToken(
		fixedRefreshToken,
		fixedUserID,
		time.Now().UTC().Add(time.Hour),
		fixedRefreshCount,
	)
	err := token.CanBeRefreshed()
	require.NoError(t, err)

	//过期失败
	expiredToken := domain.LoadRefreshToken(
		fixedRefreshToken,
		fixedUserID,
		time.Now().UTC().Add(-time.Hour),
		fixedRefreshCount,
	)
	err = expiredToken.CanBeRefreshed()
	require.Error(t, err)

	//超过刷新次数失败
	limitToken := domain.LoadRefreshToken(
		fixedRefreshToken,
		fixedUserID,
		time.Now().UTC().Add(time.Hour),
		maxRefreshCount,
	)
	err = limitToken.CanBeRefreshed()
	require.Error(t, err)
}
