package domain

import "errors"

var (
	ErrEmptyPassword        = errors.New("password is empty")
	ErrInvalidPassword      = errors.New("invalid password")
	ErrAccessTokenGenerated = errors.New("access token already generated")
	ErrRefreshTokenExpired  = errors.New("refresh token expired")
	ErrRefreshLimitExceeded = errors.New("refresh limit exceeded")
)
