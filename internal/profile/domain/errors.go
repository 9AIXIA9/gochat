package domain

import (
	myErrors "gochat/internal/shared/errors"
)

var (
	ErrNoPermission = myErrors.NewBusiness("no permission")
)
