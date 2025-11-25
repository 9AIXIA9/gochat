package domain_test

import (
	"errors"
	"gochat/internal/authorization/domain"
	"gochat/internal/authorization/domain/mocks"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const (
	fixedToken                      domain.RefreshToken = "refresh-token-abc"
	fixedNewToken                   domain.RefreshToken = "new-refresh-token"
	fixedAccessToken                domain.AccessToken  = "access-token-xyz"
	fixedUserIDRT                   kernel.UserID       = "user-123"
	fixedRefreshCount                                   = 3
	validityDuration                                    = 7 * 24 * time.Hour // mirrored business rule
	extendDuration                                      = 24 * time.Hour     // mirrored business rule
	refreshTokenEntityTimeTolerance                     = 200 * time.Millisecond
)

// helper to load a token with given params
func load(fixedToken domain.RefreshToken, uid kernel.UserID, exp time.Time, count int) *domain.RefreshTokenEntity {
	return domain.LoadRefreshToken(fixedToken, uid, exp, count)
}

func TestRefreshTokenEntity_CreateRefreshToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gen := mocks.NewMockRefreshTokenGenerator(ctrl)
	gen.EXPECT().Generate().Return(fixedToken, nil)

	start := time.Now()
	rt, err := domain.CreateRefreshToken(fixedUserIDRT, gen)
	require.NoError(t, err)
	require.NotNil(t, rt)
	assert.Equal(t, fixedUserIDRT, rt.UserID())
	assert.Equal(t, fixedToken, rt.Token())
	assert.Equal(t, 0, rt.RefreshCount())
	// Expiry should be roughly now + validityDuration
	assert.WithinDuration(t, start.Add(validityDuration), rt.ExpiredAt(), refreshTokenEntityTimeTolerance)

	// generator error path
	genErr := mocks.NewMockRefreshTokenGenerator(ctrl)
	genErr.EXPECT().Generate().Return(domain.RefreshToken(""), errors.New("generation error"))
	rt2, err := domain.CreateRefreshToken(fixedUserIDRT, genErr)
	require.Error(t, err)
	require.Nil(t, rt2)
}

func TestRefreshTokenEntity_LoadRefreshToken(t *testing.T) {
	exp := time.Now().Add(time.Hour)
	rt := load(fixedToken, fixedUserIDRT, exp, fixedRefreshCount)
	assert.Equal(t, fixedToken, rt.Token())
	assert.Equal(t, fixedUserIDRT, rt.UserID())
	assert.Equal(t, fixedRefreshCount, rt.RefreshCount())
	assert.Equal(t, exp, rt.ExpiredAt())
}

func TestRefreshTokenEntity_CanBeRefreshed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type tc struct {
		name string
		rt   *domain.RefreshTokenEntity
		err  error
	}

	cases := []tc{
		{
			name: "valid - under max and not expired",
			rt:   load(fixedToken, fixedUserIDRT, time.Now().Add(time.Hour), (7*24)-1),
			err:  nil,
		},
		{
			name: "expired",
			rt:   load(fixedToken, fixedUserIDRT, time.Now().Add(-time.Minute), 0),
			err:  myErrors.ErrExpired,
		},
		{
			name: "exceeded max",
			rt:   load(fixedToken, fixedUserIDRT, time.Now().Add(time.Hour), 7*24),
			err:  myErrors.ErrExceedMaxValue,
		},
	}

	for _, c := range cases {
		got := c.rt.CanBeRefreshed()
		if c.err == nil {
			require.NoError(t, got, c.name)
		} else {
			require.ErrorIs(t, got, c.err, c.name)
		}
	}
}

func TestRefreshTokenEntity_Refresh(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	genSuccess := mocks.NewMockRefreshTokenGenerator(ctrl)
	genSuccess.EXPECT().Generate().Return(fixedNewToken, nil)

	initialExp := time.Now().Add(time.Hour)
	rt := load(fixedToken, fixedUserIDRT, initialExp, 0)
	err := rt.Refresh(genSuccess)
	require.NoError(t, err)
	assert.Equal(t, fixedNewToken, rt.Token())
	assert.Equal(t, 1, rt.RefreshCount())
	assert.WithinDuration(t, initialExp.Add(extendDuration), rt.ExpiredAt(), refreshTokenEntityTimeTolerance)
	// After refresh, generated flag should allow generating access token
	accGen := mocks.NewMockAccessTokenGenerator(ctrl)
	accGen.EXPECT().Generate(fixedUserIDRT).Return(fixedAccessToken, nil)
	accessToken, err := rt.GenerateAccessToken(accGen)
	require.NoError(t, err)
	assert.Equal(t, fixedAccessToken, accessToken)
	// second attempt must fail
	_, err = rt.GenerateAccessToken(accGen)
	require.ErrorIs(t, err, myErrors.ErrAlreadyDone)

	// expired token refresh
	genExpired := mocks.NewMockRefreshTokenGenerator(ctrl)
	genExpired.EXPECT().Generate().Times(0)
	rtExpired := load(fixedToken, fixedUserIDRT, time.Now().Add(-time.Minute), 0)
	err = rtExpired.Refresh(genExpired)
	require.ErrorIs(t, err, myErrors.ErrExpired)

	// exceeded max count
	genExceed := mocks.NewMockRefreshTokenGenerator(ctrl)
	genExceed.EXPECT().Generate().Times(0)
	rtExceed := load(fixedToken, fixedUserIDRT, time.Now().Add(time.Hour), 7*24)
	err = rtExceed.Refresh(genExceed)
	require.ErrorIs(t, err, myErrors.ErrExceedMaxValue)

	// generator error
	genErr := mocks.NewMockRefreshTokenGenerator(ctrl)
	genErr.EXPECT().Generate().Return(domain.RefreshToken(""), errors.New("generation error"))
	rtGenErr := load(fixedToken, fixedUserIDRT, time.Now().Add(time.Hour), 0)
	err = rtGenErr.Refresh(genErr)
	require.Error(t, err)
}

func TestRefreshTokenEntity_GenerateAccessToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// create new token via factory (generated flag false & not expired)
	rtGen := mocks.NewMockRefreshTokenGenerator(ctrl)
	rtGen.EXPECT().Generate().Return(fixedToken, nil)
	rt, err := domain.CreateRefreshToken(fixedUserIDRT, rtGen)
	require.NoError(t, err)

	accGen := mocks.NewMockAccessTokenGenerator(ctrl)
	accGen.EXPECT().Generate(fixedUserIDRT).Return(fixedAccessToken, nil)
	accessToken, err := rt.GenerateAccessToken(accGen)
	require.NoError(t, err)
	assert.Equal(t, fixedAccessToken, accessToken)

	// second attempt should fail
	_, err = rt.GenerateAccessToken(accGen)
	require.ErrorIs(t, err, myErrors.ErrAlreadyDone)

	// expired scenario returns ErrAlreadyDone per implementation
	accGenExpired := mocks.NewMockAccessTokenGenerator(ctrl)
	accGenExpired.EXPECT().Generate(fixedUserIDRT).Times(0)
	rtExpired := load(fixedToken, fixedUserIDRT, time.Now().Add(-time.Hour), 0)
	_, err = rtExpired.GenerateAccessToken(accGenExpired)
	require.ErrorIs(t, err, myErrors.ErrAlreadyDone)
}
