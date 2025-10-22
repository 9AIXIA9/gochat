package domain

import (
	"gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
	"time"
)

const (
	refreshTokenMaxRefreshCount  = 7 * 24
	refreshTokenValidityDuration = 7 * 24 * time.Hour
)

type RefreshToken struct {
	string       string
	userID       kernel.UserID
	expiredAt    time.Time
	refreshCount int
}

func NewRefreshToken(string string, userID kernel.UserID, expiredAt time.Time, refreshCount int) *RefreshToken {
	return &RefreshToken{
		string:       string,
		userID:       userID,
		expiredAt:    expiredAt,
		refreshCount: refreshCount,
	}
}

func CreateRefreshToken(string string, userID kernel.UserID) *RefreshToken {
	return &RefreshToken{
		string:       string,
		userID:       userID,
		expiredAt:    time.Now().Add(refreshTokenValidityDuration),
		refreshCount: 0,
	}
}

func (token *RefreshToken) RefreshAccessToken(
	refreshGenerator RandomStringGenerator,
	accessGenerator kernel.AccessTokenGenerator,
) (kernel.AccessToken, error) {
	if err := token.canBeRefreshed(); err != nil {
		return "", err
	}

	newString, err := refreshGenerator.Generate()
	if err != nil {
		return "", err
	}

	token.refreshCount++
	token.string = newString

	return accessGenerator.Generate(token.userID)
}

func (token *RefreshToken) canBeRefreshed() error {
	if time.Now().After(token.expiredAt) {
		return errors.ErrExpired
	}

	if refreshTokenMaxRefreshCount <= token.refreshCount {
		return errors.ErrExceedMaxValue
	}
	return nil
}

func (token *RefreshToken) String() string {
	return token.string
}

func (token *RefreshToken) UserID() kernel.UserID {
	return token.userID
}

func (token *RefreshToken) ExpiredAt() time.Time {
	return token.expiredAt
}

func (token *RefreshToken) RefreshCount() int {
	return token.refreshCount
}
