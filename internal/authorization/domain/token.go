package domain

import (
	myErrors "gochat/internal/shared/errors"
)

type (
	RefreshToken string
	AccessToken  string
)

func (t AccessToken) Validate() error {
	if len(t) == 0 {
		return myErrors.NewBusiness("access token cannot be empty")
	}
	return nil
}

func (t AccessToken) String() string {
	return string(t)
}

func (t RefreshToken) Validate() error {
	if len(t) == 0 {
		return myErrors.NewBusiness("refresh token cannot be empty")
	}
	return nil
}

func (t RefreshToken) String() string {
	return string(t)
}
