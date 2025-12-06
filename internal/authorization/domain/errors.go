package domain

import "errors"

var (
	ErrEmptyPassword        = errors.New("password is empty")
	ErrInvalidPassword      = errors.New("invalid password")
	ErrInvalidAccessToken   = errors.New("invalid access token")
	ErrAccessTokenGenerated = errors.New("access token already generated")
	ErrRefreshTokenExpired  = errors.New("refresh token expired")
	ErrRefreshLimitExceeded = errors.New("refresh limit exceeded")
)
