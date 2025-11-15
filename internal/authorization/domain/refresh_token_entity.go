package domain

import (
	"gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
	"time"
)

const (
	refreshTokenMaxRefreshCount  = 7 * 24
	refreshTokenValidityDuration = 7 * 24 * time.Hour
	refreshExtendedDuration      = 24 * time.Hour
)

type RefreshTokenEntity struct {
	token        RefreshToken
	userID       kernel.UserID
	expiredAt    time.Time
	refreshCount int
}

func NewRefreshToken(token RefreshToken, userID kernel.UserID, expiredAt time.Time, refreshCount int) *RefreshTokenEntity {
	return &RefreshTokenEntity{
		token:        token,
		userID:       userID,
		expiredAt:    expiredAt,
		refreshCount: refreshCount,
	}
}

func CreateRefreshToken(token RefreshToken, userID kernel.UserID) *RefreshTokenEntity {
	return NewRefreshToken(token, userID, time.Now().Add(refreshTokenValidityDuration), 0)
}

func (r *RefreshTokenEntity) CanBeRefreshed() error {
	if time.Now().After(r.expiredAt) {
		return errors.ErrExpired
	}

	if refreshTokenMaxRefreshCount <= r.refreshCount {
		return errors.ErrExceedMaxValue
	}
	return nil
}

func (r *RefreshTokenEntity) RefreshAccessToken(newToken RefreshToken) {
	r.token = newToken
	r.refreshCount++
	r.expiredAt.Add(refreshExtendedDuration)
}

func (r *RefreshTokenEntity) Token() RefreshToken {
	return r.token
}

func (r *RefreshTokenEntity) UserID() kernel.UserID {
	return r.userID
}

func (r *RefreshTokenEntity) ExpiredAt() time.Time {
	return r.expiredAt
}

func (r *RefreshTokenEntity) RefreshCount() int {
	return r.refreshCount
}
