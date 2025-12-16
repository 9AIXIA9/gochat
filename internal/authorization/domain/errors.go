package domain

import "errors"

var (
	ErrEmptyPassword       = errors.New("password is empty")
	ErrInvalidPassword     = errors.New("invalid password")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
	ErrInvalidAccessToken  = errors.New("invalid access token")
)
