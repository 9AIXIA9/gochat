package domain

import (
	"gochat/internal/shared/kernel"
	"time"
)

const (
	refreshTokenMaxRefreshCount  = 7 * 24
	refreshTokenValidityDuration = 7 * 24 * time.Hour
	refreshExtendedDuration      = 24 * time.Hour
)

type RefreshTokenEntity struct {
	userID       kernel.UserID
	token        RefreshToken
	expiredAt    time.Time
	refreshCount int
}

func LoadRefreshToken(
	token RefreshToken,
	userID kernel.UserID,
	expiredAt time.Time,
	refreshCount int,
) *RefreshTokenEntity {
	return &RefreshTokenEntity{
		userID:       userID,
		token:        token,
		expiredAt:    expiredAt,
		refreshCount: refreshCount,
	}
}

func CreateRefreshToken(
	userID kernel.UserID,
	generator RefreshTokenGenerator,
) (*RefreshTokenEntity, error) {
	token, err := generator.Generate()
	if err != nil {
		return nil, err
	}
	return &RefreshTokenEntity{
		userID:       userID,
		token:        token,
		expiredAt:    time.Now().UTC().Add(refreshTokenValidityDuration),
		refreshCount: 0,
	}, nil
}

func (t *RefreshTokenEntity) GenerateAccessToken(
	accessTokenGenerator AccessTokenGenerator,
) (AccessToken, error) {
	return accessTokenGenerator.Generate(t.userID)
}

func (t *RefreshTokenEntity) Refresh(
	refreshTokenGenerator RefreshTokenGenerator,
) error {
	if err := t.CanBeRefreshed(); err != nil {
		return err
	}
	newToken, err := refreshTokenGenerator.Generate()
	if err != nil {
		return err
	}
	t.token = newToken
	t.refreshCount++
	t.expiredAt = t.expiredAt.Add(refreshExtendedDuration)
	return nil
}

func (t *RefreshTokenEntity) CanBeRefreshed() error {
	if time.Now().UTC().After(t.expiredAt) || refreshTokenMaxRefreshCount <= t.refreshCount {
		return ErrInvalidRefreshToken
	}
	return nil
}

func (t *RefreshTokenEntity) Token() RefreshToken {
	return t.token
}

func (t *RefreshTokenEntity) UserID() kernel.UserID {
	return t.userID
}

func (t *RefreshTokenEntity) ExpiredAt() time.Time {
	return t.expiredAt
}

func (t *RefreshTokenEntity) RefreshCount() int {
	return t.refreshCount
}
