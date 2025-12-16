package domain

import (
	myErrors "gochat/internal/shared/errors"
)

var (
	ErrInvalidPassword     = myErrors.NewBusiness("invalid password")
	ErrInvalidRefreshToken = myErrors.NewBusiness("invalid refresh token")
	ErrInvalidAccessToken  = myErrors.NewBusiness("invalid access token")
)
